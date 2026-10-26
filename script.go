/**
 * @brief Embedded scripting VMs and owned compiler output
 *
 * @file script.go
 */
/* -----------------------------------------------------------------------------
 * Enduro/X Middleware Platform for Distributed Transaction Processing
 * Copyright (C) 2009-2016, ATR Baltic, Ltd. All Rights Reserved.
 * Copyright (C) 2017-2026, Mavimax, Ltd. All Rights Reserved.
 * This software is released under one of the following licenses:
 * LGPL or Mavimax's license for commercial use.
 * See LICENSE file for full text.
 *
 * C (as designed by Dennis Ritchie and later authors) language code is licensed
 * under Enduro/X Modified GNU Affero General Public License, version 3.
 * See LICENSE_C file for full text.
 * -----------------------------------------------------------------------------
 * LGPL license:
 *
 * This program is free software; you can redistribute it and/or modify it under
 * the terms of the GNU Lesser General Public License, version 3 as published
 * by the Free Software Foundation;
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT ANY
 * WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
 * PARTICULAR PURPOSE. See the GNU Lesser General Public License, version 3
 * for more details.
 *
 * You should have received a copy of the Lesser General Public License along
 * with this program; if not, write to the Free Software Foundation, Inc.,
 * 59 Temple Place, Suite 330, Boston, MA 02111-1307 USA
 *
 * -----------------------------------------------------------------------------
 * A commercial use license is available from Mavimax, Ltd
 * contact@mavimax.com
 * -----------------------------------------------------------------------------
 */
package atmi

/*
#cgo pkg-config: atmisrvinteg

#include "script_bridge.h"
#include <stdlib.h>
#include <string.h>
*/
import "C"
import (
	"fmt"
	"runtime"
	"runtime/cgo"
	"strings"
	"unsafe"
)

// TpScrCfg holds one native scripting configuration; nil selects provider defaults.
// The zero value is ready to use. Copies of an initialized value share the native record; use a new
// value for independent state. Access a record serially.
type TpScrCfg struct {
	c_ctl *C.tpscr_cfg_t
}

// native returns the single backing C record, allocating it on first use.
// @return Native record shared by this wrapper and its copies
func (ctl *TpScrCfg) native() *C.tpscr_cfg_t {
	if nil == ctl.c_ctl {
		ctl.c_ctl = new(C.tpscr_cfg_t)
		ctl.c_ctl.version = C.NDRX_SCR_CFG_VERSION_1
	}
	return ctl.c_ctl
}

// GetVersion reads scripting configuration version from the native record.
// @return Current scripting configuration version
func (ctl *TpScrCfg) GetVersion() int {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return 0
	}
	return int(ctl.native().version)
}

// SetVersion writes scripting configuration version to the native record.
// @param value New scripting configuration version
// @return Nil on success, otherwise TPEINVAL without changing the field
func (ctl *TpScrCfg) SetVersion(value int) ATMIError {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return NewCustomATMIError(TPEINVAL, "TpScrCfg is nil")
	}
	if int(C.int(value)) != value {
		return NewCustomATMIError(TPEINVAL, "Version exceeds native int")
	}
	ctl.native().version = C.int(value)
	return nil
}

// GetInitString reads provider initialization string from the native record.
// @return Current provider initialization string
func (ctl *TpScrCfg) GetInitString() string {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return ""
	}
	return C.GoString(&ctl.native().initstring[0])
}

// SetInitString writes provider initialization string to the native record.
// @param value New provider initialization string
// @return Nil on success, otherwise TPEINVAL without changing the field
func (ctl *TpScrCfg) SetInitString(value string) ATMIError {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return NewCustomATMIError(TPEINVAL, "TpScrCfg is nil")
	}
	return copyCText(ctl.native().initstring[:], value, "InitString")
}

// TpScrVm is an opaque scripting VM owned by the ATMICtx that created it.
// All VM operations use that context, so a VM follows the context's rule of
// one goroutine at a time; use a VM per context for parallel work. Copies share
// the same native resource. Release it with TpScrUninit, or FreeATMICtx of the
// owning context. Registered callbacks retain their Go functions and user data
// until removed or closed.
type TpScrVm struct{ state *scriptVMState }

type scriptVMState struct {
	ac        *ATMICtx // owning context used by every operation
	ptr       *C.tpscr_vm_t
	callbacks map[string]*scriptRegistration
	retained  []*scriptRegistration
	active    int // native operations which may invoke host callbacks
}

// TpScrBytecode owns a compiler allocation of the context that compiled it.
// Release it with TpScrFree, or FreeATMICtx of the owning context. Copies share
// the same allocation; Bytes returns an independent Go copy.
type TpScrBytecode struct{ state *scriptBytecodeState }

type scriptBytecodeState struct {
	ac     *ATMICtx // owning context used to release the allocation
	ptr    *C.char
	length C.long
}

// Bytes copies compiled bytecode into Go storage, or returns nil after release.
// @return Independent Go copy of compiled bytecode, or nil if the allocation
//
//	is nil or released
func (code *TpScrBytecode) Bytes() []byte {
	if nil == code || nil == code.state || nil == code.state.ptr {
		return nil
	}
	return append([]byte(nil), unsafe.Slice((*byte)(unsafe.Pointer(code.state.ptr)),
		int(code.state.length))...)
}

// TpScrCallback processes ordinary typed buffers on the invoking context.
// Buffer pointers and lengths are maintained by the binding, as for TpCall.
// Input wrappers are borrowed and must not be retained after the callback.
// Returning a buffer transfers its native allocation back to the provider.
// @param ctx Borrowed invoking context; do not terminate or free it
// @param vm Invoking VM; supports nested execution but not active shutdown
// @param name Callback registration name
// @param param Parameter buffer, or nil
// @param input Input buffer, or nil
// @param user User data supplied to TpScrRegCB
// @param flags Execution flags supplied by the provider
// @return Output buffer and error; output may alias an input, including on failure
// Return TpScrError to preserve an application-specific ScriptCode.
type TpScrCallback func(ctx *ATMICtx, vm *TpScrVm, name string,
	param TypedBuffer, input TypedBuffer, user interface{},
	flags int64) (TypedBuffer, ATMIError)

type scriptRegistration struct {
	handle   cgo.Handle
	token    unsafe.Pointer
	vm       *TpScrVm
	callback TpScrCallback
	user     interface{}
}

// TpScrError preserves both the ATMI failure and the scripting engine result.
type TpScrError struct {
	ATMICode   int
	ScriptCode int
	Text       string
}

// Error returns the ATMI and scripting failure details.
// @return Formatted ATMI code, scripting code and diagnostic text
func (e *TpScrError) Error() string {
	return fmt.Sprintf("%d: %s (script code %d)", e.ATMICode, e.Text, e.ScriptCode)
}

// Code returns the ATMI error code.
// @return ATMI failure code; the separate ScriptCode field holds the provider
//
//	result
func (e *TpScrError) Code() int {
	return e.ATMICode
}

// Message returns the scripting engine error text.
// @return Scripting engine diagnostic text
func (e *TpScrError) Message() string {
	return e.Text
}

// newScriptError captures VM errors before any cleanup can overwrite them.
// The ATMI error is kept by the context and the script error by the VM, so
// the capture may run on any OS thread. Initialization and shutdown failures,
// which have no usable VM, are captured in C on the failing thread instead.
// @param vm Live provider VM
// @return ATMIError retaining native and scripting diagnostics
func (ac *ATMICtx) newScriptError(vm *C.tpscr_vm_t) ATMIError {
	native := ac.NewATMIError()
	result := &TpScrError{ATMICode: native.Code(),
		ScriptCode: int(C.Otpscrerrno(&ac.c_ctx, vm)),
		Text:       C.GoString(C.Otpscrerror(&ac.c_ctx, vm))}
	if result.Text == "" {
		result.Text = native.Message()
	}
	return result
}

// TpScrInit creates a VM using cfg and the scripting plugin selected by the
// environment. Initialization errors include the plugin's error details.
// @param cfg Optional TpScrCfg configured with SetVersion and SetInitString; nil uses
//
//	configured plugin defaults
//
// @param flags Native scripting initialization flags; normally 0
// @return Owned VM and nil error, or nil and ATMIError with initialization
//
//	diagnostics; release with TpScrUninit
func (ac *ATMICtx) TpScrInit(cfg *TpScrCfg, flags int64) (*TpScrVm, ATMIError) {
	defer runtime.KeepAlive(ac)
	defer runtime.KeepAlive(cfg)
	if err := ac.checkContext(); nil != err {
		return nil, err
	}
	if ac.scriptClosing {
		return nil, NewCustomATMIError(TPEPROTO, "ATMI context is being freed")
	}
	var config *C.tpscr_cfg_t
	if nil != cfg {
		config = cfg.native()
	}
	var nerr C.go_script_err_t
	ptr := C.go_tpscrinit(&ac.c_ctx, config, C.long(flags), &nerr)
	if nil == ptr {
		result := &TpScrError{ATMICode: int(nerr.atmi_code),
			ScriptCode: int(nerr.script_code), Text: C.GoString(nerr.script_text)}
		if result.Text == "" {
			result.Text = C.GoString(nerr.atmi_text)
		}
		C.free(unsafe.Pointer(nerr.atmi_text))
		C.free(unsafe.Pointer(nerr.script_text))
		return nil, result
	}
	state := &scriptVMState{ac: ac, ptr: ptr,
		callbacks: make(map[string]*scriptRegistration)}
	if nil == ac.scriptVMs {
		ac.scriptVMs = make(map[*scriptVMState]struct{})
	}
	ac.scriptVMs[state] = struct{}{}
	return &TpScrVm{state: state}, nil
}

// TpScrUninit destroys the VM and releases all retained callback handles on success.
// No execution or callback may be active while the VM is being destroyed.
// The handle is invalidated even on failure; callback storage is then retained.
// @param flags Native scripting shutdown flags; normally 0
// @return Nil on success, otherwise ATMIError; a native shutdown attempt
//
//	invalidates the VM even on failure
func (vm *TpScrVm) TpScrUninit(flags int64) ATMIError {
	if nil == vm || nil == vm.state || nil == vm.state.ptr {
		return NewCustomATMIError(TPEINVAL, "script VM is nil or released")
	}
	ac := vm.state.ac
	defer runtime.KeepAlive(ac)
	if err := ac.checkContext(); nil != err {
		return err
	}
	if vm.state.active != 0 {
		return NewCustomATMIError(TPEPROTO, "cannot destroy a scripting VM during a native operation")
	}
	state := vm.state
	state.active++
	defer func() { state.active-- }()
	// Shutdown callbacks see the VM as released, never the native VM in teardown.
	ptr := state.ptr
	state.ptr = nil
	var nerr C.go_script_err_t
	ret := C.go_tpscruninit(&ac.c_ctx, ptr, C.long(flags), &nerr)
	delete(ac.scriptVMs, state)
	if ret != SUCCEED {
		// A failing provider may already have consumed the VM. Keep userdata
		// alive, but never reuse or query the possibly freed native handle.
		result := &TpScrError{ATMICode: int(nerr.atmi_code),
			ScriptCode: int(nerr.script_code), Text: C.GoString(nerr.script_text)}
		if result.Text == "" {
			result.Text = C.GoString(nerr.atmi_text)
		}
		C.free(unsafe.Pointer(nerr.atmi_text))
		C.free(unsafe.Pointer(nerr.script_text))
		return result
	}
	for _, registration := range vm.state.callbacks {
		registration.handle.Delete()
		C.free(registration.token)
	}
	for _, registration := range vm.state.retained {
		registration.handle.Delete()
		C.free(registration.token)
	}
	vm.state.callbacks = nil
	vm.state.retained = nil
	return nil
}

// TpScrComp compiles source and returns an owned native bytecode allocation.
// Use Bytes to obtain a Go copy, and TpScrFree to release the allocation.
// @param source Script source text, passed with its byte length
// @param flags Native compiler flags; normally 0
// @return Owned compiler output and nil error, or nil and ATMIError; release
//
//	compiler storage with TpScrFree
func (vm *TpScrVm) TpScrComp(source string,
	flags int64) (*TpScrBytecode, ATMIError) {
	if nil == vm || nil == vm.state || nil == vm.state.ptr {
		return nil, NewCustomATMIError(TPEINVAL, "script VM is nil or released")
	}
	ac := vm.state.ac
	defer runtime.KeepAlive(ac)
	if err := ac.checkContext(); nil != err {
		return nil, err
	}
	if ac.scriptClosing {
		return nil, NewCustomATMIError(TPEPROTO, "ATMI context is being freed")
	}
	if err := ac.syncScriptCallback(); nil != err {
		return nil, err
	}
	state := vm.state
	state.active++
	defer func() { state.active-- }()
	src := C.CString(source)
	defer C.free(unsafe.Pointer(src))
	var ptr *C.char
	var length C.long
	ret := C.Otpscrcomp(&ac.c_ctx, vm.state.ptr, src, C.long(len(source)),
		&ptr, &length, C.long(flags))
	if ret != SUCCEED {
		err := ac.newScriptError(vm.state.ptr)
		if nil != ptr {
			C.Otpscrfree(&ac.c_ctx, ptr)
		}
		return nil, err
	}
	if nil == ptr || length < 0 || C.long(int(length)) != length {
		if nil != ptr {
			C.Otpscrfree(&ac.c_ctx, ptr)
		}
		return nil, NewCustomATMIError(TPESYSTEM, "invalid compiler bytecode result")
	}
	code := &scriptBytecodeState{ac: ac, ptr: ptr, length: length}
	if nil == ac.scriptCodes {
		ac.scriptCodes = make(map[*scriptBytecodeState]struct{})
	}
	ac.scriptCodes[code] = struct{}{}
	return &TpScrBytecode{state: code}, nil
}

// TpScrFree releases compiler output with the scripting allocator, using the
// context that compiled it. Nil and already released bytecode need no action.
// @return Nil on success, otherwise ATMIError
func (code *TpScrBytecode) TpScrFree() ATMIError {
	if nil == code || nil == code.state || nil == code.state.ptr {
		return nil
	}
	ac := code.state.ac
	defer runtime.KeepAlive(ac)
	if err := ac.checkContext(); nil != err {
		return err
	}
	if C.Otpscrfree(&ac.c_ctx, code.state.ptr) != SUCCEED {
		return ac.NewATMIError()
	}
	code.state.ptr = nil
	code.state.length = 0
	delete(ac.scriptCodes, code.state)
	return nil
}

// TpScrLoad loads bytecode under name. Go owns bytes and may reuse it afterward.
// flags accepts the native NDRX_TPSCR_* loading options.
// @param name NUL-free script or package name
// @param bytes Compiled bytecode copied into temporary native storage for this
//
//	call
//
// @param flags NDRX_TPSCR_PACKAGE, NDRX_TPSCR_REPLACE, NDRX_TPSCR_FLAT or 0 as
//
//	supported by the provider
//
// @return Nil on success, otherwise ATMIError
func (vm *TpScrVm) TpScrLoad(name string, bytes []byte,
	flags int64) ATMIError {
	if nil == vm || nil == vm.state || nil == vm.state.ptr {
		return NewCustomATMIError(TPEINVAL, "script VM is nil or released")
	}
	ac := vm.state.ac
	defer runtime.KeepAlive(ac)
	if err := ac.checkContext(); nil != err {
		return err
	}
	if strings.IndexByte(name, 0) >= 0 {
		return NewCustomATMIError(TPEINVAL, "script name contains NUL")
	}
	if err := ac.syncScriptCallback(); nil != err {
		return err
	}
	state := vm.state
	state.active++
	defer func() { state.active-- }()
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	// A non-NULL address also represents an explicitly empty bytecode block.
	data := C.malloc(C.size_t(len(bytes) + 1))
	if nil == data {
		return NewCustomATMIError(TPEOS, "bytecode allocation failed")
	}
	defer C.free(data)
	if len(bytes) != 0 {
		C.memcpy(data, unsafe.Pointer(&bytes[0]), C.size_t(len(bytes)))
	}
	if C.Otpscrload(&ac.c_ctx, vm.state.ptr, cName, (*C.char)(data),
		C.long(len(bytes)), C.long(flags)) != SUCCEED {
		return ac.newScriptError(vm.state.ptr)
	}
	return nil
}

// TpScrLoadStr compiles and loads source under name with NDRX_TPSCR_* flags.
// @param name NUL-free script or package name
// @param source Script source text, passed with its byte length
// @param flags NDRX_TPSCR_PACKAGE, NDRX_TPSCR_REPLACE, NDRX_TPSCR_FLAT or 0 as
//
//	supported by the provider
//
// @return Nil on success, otherwise ATMIError
func (vm *TpScrVm) TpScrLoadStr(name, source string,
	flags int64) ATMIError {
	if nil == vm || nil == vm.state || nil == vm.state.ptr {
		return NewCustomATMIError(TPEINVAL, "script VM is nil or released")
	}
	ac := vm.state.ac
	defer runtime.KeepAlive(ac)
	if err := ac.checkContext(); nil != err {
		return err
	}
	if strings.IndexByte(name, 0) >= 0 {
		return NewCustomATMIError(TPEINVAL, "script name contains NUL")
	}
	if err := ac.syncScriptCallback(); nil != err {
		return err
	}
	state := vm.state
	state.active++
	defer func() { state.active-- }()
	cName, src := C.CString(name), C.CString(source)
	defer C.free(unsafe.Pointer(cName))
	defer C.free(unsafe.Pointer(src))
	if C.Otpscrloadstr(&ac.c_ctx, vm.state.ptr, cName, src,
		C.long(len(source)), C.long(flags)) != SUCCEED {
		return ac.newScriptError(vm.state.ptr)
	}
	return nil
}

// TpScrUnload removes the loaded script identified by name.
// @param name Loaded script or package name
// @param flags Native unload flags; normally 0
// @return Nil on success, otherwise ATMIError
func (vm *TpScrVm) TpScrUnload(name string, flags int64) ATMIError {
	if nil == vm || nil == vm.state || nil == vm.state.ptr {
		return NewCustomATMIError(TPEINVAL, "script VM is nil or released")
	}
	ac := vm.state.ac
	defer runtime.KeepAlive(ac)
	if err := ac.checkContext(); nil != err {
		return err
	}
	if strings.IndexByte(name, 0) >= 0 {
		return NewCustomATMIError(TPEINVAL, "script name contains NUL")
	}
	if err := ac.syncScriptCallback(); nil != err {
		return err
	}
	state := vm.state
	state.active++
	defer func() { state.active-- }()
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	if C.Otpscrunload(&ac.c_ctx, vm.state.ptr, cName, C.long(flags)) != SUCCEED {
		return ac.newScriptError(vm.state.ptr)
	}
	return nil
}

// TpScrRegCB registers callback and user data under name. C retains only a
// C-allocated handle token; it never retains a Go pointer. Use the VM serially.
// @param name NUL-free callback registration name
// @param callback Go callback invoked synchronously with borrowed context and
//
//	mutable buffer slots
//
// @param user Go user data retained until unregister or VM shutdown
// @param flags Native registration flags; normally 0
// @return Nil on success, otherwise ATMIError; failed native registration may
//
//	retain callback storage until VM shutdown
func (vm *TpScrVm) TpScrRegCB(name string, callback TpScrCallback,
	user interface{}, flags int64) ATMIError {
	if nil == vm || nil == vm.state || nil == vm.state.ptr || nil == callback {
		return NewCustomATMIError(TPEINVAL, "script VM or callback is nil/released")
	}
	ac := vm.state.ac
	defer runtime.KeepAlive(ac)
	if err := ac.checkContext(); nil != err {
		return err
	}
	if strings.IndexByte(name, 0) >= 0 {
		return NewCustomATMIError(TPEINVAL, "callback name contains NUL")
	}
	if err := ac.syncScriptCallback(); nil != err {
		return err
	}
	state := vm.state
	state.active++
	defer func() { state.active-- }()
	registration := &scriptRegistration{vm: vm, callback: callback, user: user}
	registration.token = C.malloc(C.sizeof_uintptr_t)
	if nil == registration.token {
		return NewCustomATMIError(TPEOS, "callback token allocation failed")
	}
	registration.handle = cgo.NewHandle(registration)
	*(*C.uintptr_t)(registration.token) = C.uintptr_t(registration.handle)
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	if C.go_tpscrregcb(&ac.c_ctx, vm.state.ptr, cName, registration.token,
		C.long(flags)) != SUCCEED {
		// A plugin may have installed the callback before reporting failure.
		// Retain its token until VM destruction to avoid a dangling user value.
		vm.state.retained = append(vm.state.retained, registration)
		return ac.newScriptError(vm.state.ptr)
	}
	previous := vm.state.callbacks[name]
	vm.state.callbacks[name] = registration
	if nil != previous {
		previous.handle.Delete()
		C.free(previous.token)
	}
	return nil
}

// TpScrUnregCB unregisters name and releases its Go callback and user data.
// @param name Callback registration name to remove
// @param flags Native unregister flags; normally 0
// @return Nil on success, otherwise ATMIError; successful removal releases
//
//	callback and user data
func (vm *TpScrVm) TpScrUnregCB(name string, flags int64) ATMIError {
	if nil == vm || nil == vm.state || nil == vm.state.ptr {
		return NewCustomATMIError(TPEINVAL, "script VM is nil or released")
	}
	ac := vm.state.ac
	defer runtime.KeepAlive(ac)
	if err := ac.checkContext(); nil != err {
		return err
	}
	if strings.IndexByte(name, 0) >= 0 {
		return NewCustomATMIError(TPEINVAL, "callback name contains NUL")
	}
	if err := ac.syncScriptCallback(); nil != err {
		return err
	}
	state := vm.state
	state.active++
	defer func() { state.active-- }()
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	if C.Otpscrunregcb(&ac.c_ctx, vm.state.ptr, cName, C.long(flags)) != SUCCEED {
		return ac.newScriptError(vm.state.ptr)
	}
	if registration := vm.state.callbacks[name]; nil != registration {
		delete(vm.state.callbacks, name)
		registration.handle.Delete()
		C.free(registration.token)
	}
	return nil
}

// TpScrExec executes a script using the existing typed-buffer representation.
// Native input pointer and length changes update the supplied buffers directly,
// as for TpCall and TpGetRply. Returned aliases reuse the same ATMIBuf owner.
// @param name Loaded script or function name
// @param param Parameter buffer, or nil
// @param input Input buffer, or nil
// @param flags Native execution flags forwarded to the script
// @return Output buffer and error; inspect the output even on failure
func (vm *TpScrVm) TpScrExec(name string,
	param TypedBuffer, input TypedBuffer, flags int64) (TypedBuffer, ATMIError) {
	defer runtime.KeepAlive(param)
	defer runtime.KeepAlive(input)
	if nil == vm || nil == vm.state || nil == vm.state.ptr {
		return nil, NewCustomATMIError(TPEINVAL, "script VM is nil or released")
	}
	ac := vm.state.ac
	defer runtime.KeepAlive(ac)
	if err := ac.checkContext(); nil != err {
		return nil, err
	}
	if strings.IndexByte(name, 0) >= 0 {
		return nil, NewCustomATMIError(TPEINVAL, "script name contains NUL")
	}
	if err := ac.syncScriptCallback(); nil != err {
		return nil, err
	}
	state := vm.state
	state.active++
	defer func() { state.active-- }()
	pbuf, ibuf := nativeBuffer(param), nativeBuffer(input)
	var output *C.char
	var olen C.long
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	// Direct field arguments let cgo restrict pinning to the native fields.
	var ret C.int
	if nil != pbuf && nil != ibuf {
		ret = C.Otpscrexec(&ac.c_ctx, state.ptr, cName,
			&pbuf.C_ptr, &pbuf.C_len, &ibuf.C_ptr, &ibuf.C_len,
			&output, &olen, C.long(flags))
	} else if nil != pbuf {
		ret = C.Otpscrexec(&ac.c_ctx, state.ptr, cName,
			&pbuf.C_ptr, &pbuf.C_len, nil, nil, &output, &olen, C.long(flags))
	} else if nil != ibuf {
		ret = C.Otpscrexec(&ac.c_ctx, state.ptr, cName,
			nil, nil, &ibuf.C_ptr, &ibuf.C_len, &output, &olen, C.long(flags))
	} else {
		ret = C.Otpscrexec(&ac.c_ctx, state.ptr, cName,
			nil, nil, nil, nil, &output, &olen, C.long(flags))
	}
	var err ATMIError
	if SUCCEED != ret {
		err = ac.newScriptError(state.ptr)
	}
	var result *ATMIBuf
	if nil != output {
		if nil != pbuf && output == pbuf.C_ptr {
			result = pbuf
		} else if nil != ibuf && output == ibuf.C_ptr {
			result = ibuf
		} else {
			// A nested callback may return an outer borrowed buffer.
			handle := C.go_script_frame_get()
			if 0 != handle {
				frame := cgo.Handle(handle).Value().(*scriptCallbackFrame)
				for ; nil != frame && nil == result; frame = frame.previous {
					for _, buffer := range frame.buffers {
						if nil != buffer && output == buffer.C_ptr {
							result = buffer
							break
						}
					}
				}
			}
			if nil == result {
				result = &ATMIBuf{Ctx: ac, HaveFinalizer: true}
				runtime.SetFinalizer(result, tpfree)
			}
		}
		result.C_ptr, result.C_len = output, olen
	}
	for _, buffer := range []TypedBuffer{param, input} {
		if view, ok := buffer.(*TypedVIEW); ok && nil != view &&
			nil != view.Buf && nil != view.Buf.C_ptr {
			var kind string
			_, typeErr := ac.TpTypes(view.Buf, &kind, &view.view)
			if nil == err {
				err = typeErr
			}
		}
	}
	if syncErr := ac.syncScriptCallback(); nil == err {
		err = syncErr
	}
	if nil == result {
		return nil, err
	}
	return result, err
}

// TpScrError copies the current error message of a live VM.
// Initialization failures are already captured in the error from TpScrInit.
// @return Copied provider diagnostic string and nil error, or empty string and
//
//	ATMIError
func (vm *TpScrVm) TpScrError() (string, ATMIError) {
	if nil == vm || nil == vm.state || nil == vm.state.ptr {
		return "", NewCustomATMIError(TPEINVAL, "script VM is nil or released")
	}
	ac := vm.state.ac
	defer runtime.KeepAlive(ac)
	if err := ac.checkContext(); nil != err {
		return "", err
	}
	return C.GoString(C.Otpscrerror(&ac.c_ctx, vm.state.ptr)), nil
}

// TpScrErrno returns the current scripting error code of a live VM.
// @return Current provider error code and nil error, or 0 and ATMIError
func (vm *TpScrVm) TpScrErrno() (int, ATMIError) {
	if nil == vm || nil == vm.state || nil == vm.state.ptr {
		return 0, NewCustomATMIError(TPEINVAL, "script VM is nil or released")
	}
	ac := vm.state.ac
	defer runtime.KeepAlive(ac)
	if err := ac.checkContext(); nil != err {
		return 0, err
	}
	return int(C.Otpscrerrno(&ac.c_ctx, vm.state.ptr)), nil
}

// TpScrSetError sets the VM's script code/message, including from a host callback.
// @param code Provider or application scripting error code
// @param message NUL-free diagnostic text
// @return Nil on success, otherwise ATMIError; setting an error does not
//
//	itself fail a callback
func (vm *TpScrVm) TpScrSetError(code int, message string) ATMIError {
	if nil == vm || nil == vm.state || nil == vm.state.ptr {
		return NewCustomATMIError(TPEINVAL, "script VM is nil or released")
	}
	ac := vm.state.ac
	defer runtime.KeepAlive(ac)
	if err := ac.checkContext(); nil != err {
		return err
	}
	if int(C.int(code)) != code || strings.IndexByte(message, 0) >= 0 {
		return NewCustomATMIError(TPEINVAL, "script code exceeds native int or message contains NUL")
	}
	text := C.CString(message)
	defer C.free(unsafe.Pointer(text))
	if C.Otpscrseterror(&ac.c_ctx, vm.state.ptr, C.int(code), text) != SUCCEED {
		return ac.NewATMIError()
	}
	return nil
}

/* vim: set ts=4 sw=4 et smartindent: */
