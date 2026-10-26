/**
 * @brief Go callbacks invoked by the scripting engine
 *
 * @file script_callback.go
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
*/
import "C"
import (
	"fmt"
	"runtime"
	"runtime/cgo"
	"strings"
	"unsafe"
)

// scriptCallbackFrame links existing ATMIBuf objects to provider slots.
// Its handle is published on the pinned callback thread; C retains no Go pointer.
type scriptCallbackFrame struct {
	previous      *scriptCallbackFrame
	buffers       [3]*ATMIBuf
	nativeSlots   [3]**C.char
	nativeLengths [3]*C.long
}

// syncScriptCallback publishes buffer pointers and their stored lengths before
// nested execution and afterward, so provider views follow buffer reallocations.
// @return Nil after writeback, otherwise TPEINVAL for an invalid buffer length
func (ac *ATMICtx) syncScriptCallback() ATMIError {
	handle := C.go_script_frame_get()
	if 0 == handle {
		return nil
	}
	frame := cgo.Handle(handle).Value().(*scriptCallbackFrame)
	var err ATMIError
	for ; nil != frame; frame = frame.previous {
		for i, buffer := range frame.buffers {
			var ptr *C.char
			var length C.long
			if nil != buffer && nil != buffer.C_ptr {
				ptr, length = buffer.C_ptr, buffer.C_len
				if 0 > length {
					length = 0
					err = NewCustomATMIError(TPEINVAL, "negative scripting buffer length")
				}
			}
			if nil != frame.nativeSlots[i] {
				*frame.nativeSlots[i] = ptr
			}
			if nil != frame.nativeLengths[i] {
				*frame.nativeLengths[i] = length
			}
		}
	}
	return err
}

// goScriptCallback bridges native buffer slots to temporary Go wrappers.
// Native code owns the argument storage; no Go pointer is retained by C.
// The callback receives the VM's owning ATMICtx: VM operations always run on
// it, so resources created in the callback are registered with the real owner.
// An invocation from another context (e.g. a script of another VM calling this
// VM's callback) is rejected, so XATMI work never runs under a foreign context.
// @param context Native context of the invocation
// @param nativeVM Invoking native scripting VM
// @param name Native callback name
// @param param Borrowed native parameter pointer slot
// @param paramLength Parameter byte-length slot
// @param input Borrowed native input pointer slot
// @param inputLength Input byte-length slot
// @param output Native output pointer slot receiving ownership
// @param outputLength Output byte-length slot
// @param handle runtime/cgo.Handle stored in native user data
// @param flags Native execution flags
// @return SUCCEED on callback success, otherwise FAIL after buffer writeback
//
//export goScriptCallback
func goScriptCallback(context C.TPCONTEXT_T, nativeVM *C.tpscr_vm_t,
	name *C.char, param **C.char, paramLength *C.long, input **C.char,
	inputLength *C.long, output **C.char, outputLength *C.long,
	handle C.uintptr_t, flags C.long) (result C.int) {
	registration := cgo.Handle(handle).Value().(*scriptRegistration)
	ac := registration.vm.state.ac
	if context != ac.c_ctx {
		invoking := context
		text := C.CString("scripting callback invoked from another ATMI context")
		C.Otpscrseterror(&invoking, nativeVM, C.int(TPEPROTO), text)
		C.free(unsafe.Pointer(text))
		return FAIL
	}
	cSlots := [3]**C.char{param, input, output}
	cLengths := [3]*C.long{paramLength, inputLength, outputLength}
	var buffers [3]*ATMIBuf
	borrowed := make(map[*C.char]*ATMIBuf)
	var returned TypedBuffer
	var callbackErr ATMIError
	var frame *scriptCallbackFrame
	var frameHandle cgo.Handle
	var previousFrame C.uintptr_t
	result = FAIL

	defer func() {
		if recovered := recover(); nil != recovered {
			callbackErr = NewCustomATMIError(TPESYSTEM,
				fmt.Sprintf("panic in scripting callback: %v", recovered))
			result = FAIL
		}
		resultBuffer := nativeBuffer(returned)
		previousOutput := buffers[2]
		if nil != previousOutput && previousOutput != resultBuffer &&
			nil != previousOutput.C_ptr {
			survives := false
			for _, buffer := range []*ATMIBuf{buffers[0], buffers[1], resultBuffer} {
				if nil != buffer && buffer.C_ptr == previousOutput.C_ptr {
					survives = true
				}
			}
			if !survives {
				C.Otpfree(&ac.c_ctx, previousOutput.C_ptr)
				previousOutput.C_ptr, previousOutput.C_len = nil, 0
			}
		}
		if nil != frame {
			frame.buffers[2] = resultBuffer
		}
		if err := ac.syncScriptCallback(); nil != err {
			callbackErr = err
			result = FAIL
		}
		if nil == cSlots[2] && nil != resultBuffer && nil != resultBuffer.C_ptr {
			survives := false
			for _, buffer := range buffers[:2] {
				if nil != buffer && buffer.C_ptr == resultBuffer.C_ptr {
					survives = true
				}
			}
			if !survives {
				C.Otpfree(&ac.c_ctx, resultBuffer.C_ptr)
			}
		}
		transferred := make(map[*ATMIBuf]bool)
		if nil != resultBuffer {
			transferred[resultBuffer] = true
		}
		for _, buffer := range borrowed {
			transferred[buffer] = true
		}
		for buffer := range transferred {
			retained := false
			if nil != frame {
				for parent := frame.previous; nil != parent; parent = parent.previous {
					for _, outer := range parent.buffers {
						if buffer == outer {
							retained = true
						}
					}
				}
			}
			if retained {
				continue
			}
			runtime.SetFinalizer(buffer, nil)
			buffer.HaveFinalizer = false
			buffer.C_ptr, buffer.C_len = nil, 0
		}
		if nil != callbackErr {
			text := C.CString(strings.ReplaceAll(callbackErr.Message(), "\x00", "\\0"))
			code := callbackErr.Code()
			if scriptErr, ok := callbackErr.(*TpScrError); ok {
				code = scriptErr.ScriptCode
			}
			C.Otpscrseterror(&ac.c_ctx, nativeVM, C.int(code), text)
			C.free(unsafe.Pointer(text))
		}
		if 0 != frameHandle {
			C.go_script_frame_swap(previousFrame)
			frameHandle.Delete()
		}
		runtime.KeepAlive(ac)
	}()

	for i, cSlot := range cSlots {
		if nil == cSlot || nil == *cSlot {
			continue
		}
		ptr := *cSlot
		buffer := borrowed[ptr]
		if nil == buffer {
			buffer = &ATMIBuf{C_ptr: ptr, Ctx: ac}
			borrowed[ptr] = buffer
		}
		if nil != cLengths[i] {
			buffer.C_len = *cLengths[i]
		}
		buffers[i] = buffer
	}
	frame = &scriptCallbackFrame{buffers: buffers,
		nativeSlots: cSlots, nativeLengths: cLengths}
	previousFrame = C.go_script_frame_get()
	if 0 != previousFrame {
		frame.previous = cgo.Handle(previousFrame).Value().(*scriptCallbackFrame)
	}
	frameHandle = cgo.NewHandle(frame)
	C.go_script_frame_swap(C.uintptr_t(frameHandle))
	var paramBuffer, inputBuffer TypedBuffer
	if nil != buffers[0] {
		paramBuffer = buffers[0]
	}
	if nil != buffers[1] {
		inputBuffer = buffers[1]
	}
	registration.vm.state.active++
	defer func() { registration.vm.state.active-- }()
	returned, callbackErr = registration.callback(ac, registration.vm,
		C.GoString(name), paramBuffer, inputBuffer, registration.user, int64(flags))
	if nil != callbackErr {
		return FAIL
	}
	return SUCCEED
}

/* vim: set ts=4 sw=4 et smartindent: */
