/**
 * @brief Scripting lifecycle, callbacks and buffer ownership
 *
 * @file client.go
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
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	atmi "github.com/endurox-dev/endurox-go/v9"
	"github.com/endurox-dev/endurox-go/v9/tests/testutil"
	"ubftab"
)

var scripting = flag.Bool("scripting", false,
	"enable scripting tests with separately licensed endurox-python")

// Check provider reentry, aliases and script-specific errors.
// @return Error if scripting or cleanup fails
func test_scripting() (ret error) {
	duration, durationErr := testutil.Duration()
	if nil != durationErr {
		return durationErr
	}
	var started time.Time
	var i int64
	ctx, err := atmi.NewATMICtx()
	if nil != err {
		return err
	}
	defer ctx.FreeATMICtx()
	var cfg atmi.TpScrCfg
	if err := cfg.SetVersion(atmi.NDRX_SCR_CFG_VERSION_1); nil != err {
		return err
	}
	if err := cfg.SetInitString(atmi.NDRX_SCR_INITSTRING_PYTHON); nil != err {
		return err
	}
	vm, err := ctx.TpScrInit(&cfg, 0)
	if nil != err {
		return err
	}
	defer func() {
		if nil != vm {
			if err := vm.TpScrUninit(0); nil != err && nil == ret {
				ret = err
			}
		}
	}()
	// Compiled code is portable only between VMs of the same format
	if !strings.HasPrefix(cfg.GetCompFormat(), "python:") || len(cfg.GetCompFormat()) <= len("python:") {
		return fmt.Errorf("compiled code format [%s]", cfg.GetCompFormat())
	}
	var same atmi.TpScrCfg
	other, err := ctx.TpScrInit(&same, 0)
	if nil != err {
		return err
	}
	if err := other.TpScrUninit(0); nil != err {
		return err
	}
	if same.GetCompFormat() != cfg.GetCompFormat() {
		return fmt.Errorf("compiled code format differs: [%s] [%s]",
			same.GetCompFormat(), cfg.GetCompFormat())
	}
	if err := vm.TpScrSetError(1234, "manual error"); nil != err {
		return err
	}
	if code, err := vm.TpScrErrno(); nil != err || code != 1234 {
		return fmt.Errorf("manual script code: %d, %v", code, err)
	}
	if message, err := vm.TpScrError(); nil != err || message != "manual error" {
		return fmt.Errorf("manual script message: %q, %v", message, err)
	}
	if err := vm.TpScrLoadStr("go_contract_pkg", "VALUE=b'package'\n", atmi.NDRX_TPSCR_PACKAGE); nil != err {
		return err
	}
	if err := vm.TpScrLoadStr("go_contract_pkg.child", `from go_contract_pkg import VALUE
def main(name, param, incoming, flags):
    return {'data': VALUE}
`, 0); nil != err {
		return err
	}
	for _, expected := range []string{"package", "replacement"} {
		output, err := vm.TpScrExec("go_contract_pkg.child", nil, nil, 0)
		if nil != err {
			return err
		}
		value, err := ctx.CastToCarray(output.GetBuf())
		if nil != err || string(value.GetBytes()) != expected {
			return fmt.Errorf("package/replacement result: %v", err)
		}
		ctx.TpFree(output.GetBuf())
		if expected == "package" {
			if err := vm.TpScrLoadStr("go_contract_pkg.child", `def main(name, param, incoming, flags):
    return {'data': b'replacement'}
`, atmi.NDRX_TPSCR_REPLACE); nil != err {
				return err
			}
		}
	}
	if err := vm.TpScrUnload("go_contract_pkg", 0); nil != err {
		return err
	}
	if _, err := vm.TpScrExec("go_contract_pkg.child", nil, nil, 0); nil == err {
		return fmt.Errorf("unloaded package remains executable")
	}

	code, err := vm.TpScrComp(`def main(name, param, incoming, flags):
    return incoming
`, 0)
	if nil != err {
		return err
	}
	saved := append([]byte(nil), code...)
	if len(saved) == 0 {
		return fmt.Errorf("compiler returned empty bytecode")
	}
	if err := vm.TpScrLoad("echo", saved, atmi.NDRX_TPSCR_FLAT); nil != err {
		return err
	}

	err = vm.TpScrRegCB("relocate", func(call *atmi.ATMICtx, vm *atmi.TpScrVm,
		name string, param atmi.TypedBuffer, input atmi.TypedBuffer,
		user interface{}, flags int64) (atmi.TypedBuffer, atmi.ATMIError) {
		if err := input.GetBuf().TpRealloc(8 * 1024 * 1024); nil != err {
			return nil, err
		}
		return input, nil
	}, nil, 0)
	if nil != err {
		return err
	}
	inspect := fmt.Sprintf(`import builtins
def main(name, param, incoming, flags):
    incoming['data'][%d] = b'x' * (2*1024*1024)
    rc, output = relocate(param, incoming, flags)
    assert builtins._endurox_test_held['data'][%d][0] == 'grown'
    return rc, output
`, ubftab.T_CARRAY_FLD, ubftab.T_STRING_FLD)
	if err := vm.TpScrLoadStr("inspect", inspect, atmi.NDRX_TPSCR_FLAT); nil != err {
		return err
	}
	var held atmi.TypedBuffer
	err = vm.TpScrRegCB("borrow", func(call *atmi.ATMICtx, vm *atmi.TpScrVm,
		name string, param atmi.TypedBuffer, input atmi.TypedBuffer,
		user interface{}, flags int64) (atmi.TypedBuffer, atmi.ATMIError) {
		return held, nil
	}, nil, 0)
	if nil != err {
		return err
	}
	if err := vm.TpScrLoadStr("borrowed", `def main(name, param, incoming, flags):
    return borrow(None, None, flags)
`, atmi.NDRX_TPSCR_FLAT); nil != err {
		return err
	}
	err = vm.TpScrRegCB("grow", func(call *atmi.ATMICtx, vm *atmi.TpScrVm,
		name string, param atmi.TypedBuffer, input atmi.TypedBuffer,
		user interface{}, flags int64) (atmi.TypedBuffer, atmi.ATMIError) {
		if err := vm.TpScrUninit(0); nil == err || err.Code() != atmi.TPEPROTO {
			return nil, atmi.NewCustomATMIError(atmi.TPESYSTEM, "active VM destruction was not rejected")
		}
		held = input
		defer func() { held = nil }()
		alias, err := vm.TpScrExec("borrowed", nil, nil, 0)
		if nil != err {
			return nil, err
		}
		if nil == alias || alias.GetBuf() != input.GetBuf() {
			return nil, atmi.NewCustomATMIError(atmi.TPESYSTEM,
				"outer callback alias has multiple Go owners")
		}
		ubf, err := call.CastToUBF(input.GetBuf())
		if nil != err {
			return nil, err
		}
		if err := ubf.TpRealloc(131072); nil != err {
			return nil, err
		}
		if err := ubf.BChg(ubftab.T_STRING_FLD, 0, "grown"); nil != err {
			return nil, atmi.NewCustomATMIError(atmi.TPESYSTEM, err.Error())
		}
		return vm.TpScrExec("inspect", param, input, 0)
	}, nil, 0)
	if nil != err {
		return err
	}
	source := fmt.Sprintf(`import builtins
def main(name, param, incoming, flags):
    builtins._endurox_test_held = incoming
    try:
        rc, output = grow(incoming, incoming, flags)
        if rc != 0:
            return rc
        assert incoming['data'][%d][0] == 'grown'
        return output
    finally:
        del builtins._endurox_test_held
`, ubftab.T_STRING_FLD)
	if err := vm.TpScrLoadStr("nested", source, atmi.NDRX_TPSCR_FLAT); nil != err {
		return err
	}
	started = time.Now()
	for i = 0; i < 1 || time.Since(started) < duration; i++ {
		ubf, err := ctx.NewUBF(128)
		if nil != err {
			return err
		}
		input := ubf.GetBuf()
		param := input
		output, err := vm.TpScrExec("nested", param, input, 0)
		if nil != err {
			return err
		}
		if param != input || input != output.GetBuf() {
			return fmt.Errorf("aliased allocation has multiple Go owners")
		}
		got, err := ctx.CastToUBF(output.GetBuf())
		if nil != err {
			return err
		}
		value, uerr := got.BGetString(ubftab.T_STRING_FLD, 0)
		if nil != uerr || value != "grown" {
			return fmt.Errorf("nested callback writeback: %q, %v", value, uerr)
		}
		ctx.TpFree(output.GetBuf())
	}
	fmt.Printf("test_scripting loop 1: %d iterations in %.3fs\n", i, time.Since(started).Seconds())
	if err := vm.TpScrUnregCB("grow", 0); nil != err {
		return err
	}

	err = vm.TpScrRegCB("reject", func(call *atmi.ATMICtx, vm *atmi.TpScrVm,
		name string, param atmi.TypedBuffer, input atmi.TypedBuffer,
		user interface{}, flags int64) (atmi.TypedBuffer, atmi.ATMIError) {
		return nil, &atmi.TpScrError{ATMICode: atmi.TPESVCFAIL, ScriptCode: 7777, Text: "callback rejected"}
	}, nil, 0)
	if nil != err {
		return err
	}
	source = `def main(name, param, incoming, flags):
    rc, output = reject(None, None, flags)
    return {'rc':rc, 'out':output,
            'screrrno':tpscrerrno(), 'error':tpscrerror()}
`
	if err := vm.TpScrLoadStr("failure", source, atmi.NDRX_TPSCR_FLAT); nil != err {
		return err
	}
	_, err = vm.TpScrExec("failure", nil, nil, 0)
	scriptErr, ok := err.(*atmi.TpScrError)
	if !ok || scriptErr.ScriptCode != 7777 || !strings.Contains(scriptErr.Text, "callback rejected") {
		return fmt.Errorf("custom script error lost: %v", err)
	}
	if err := vm.TpScrUnregCB("reject", 0); nil != err {
		return err
	}
	err = vm.TpScrRegCB("handoff", func(call *atmi.ATMICtx, vm *atmi.TpScrVm,
		name string, param atmi.TypedBuffer, input atmi.TypedBuffer,
		user interface{}, flags int64) (atmi.TypedBuffer, atmi.ATMIError) {
		return vm.TpScrExec("echo", param, input, 0)
	}, nil, 0)
	if nil != err {
		return err
	}
	if err := vm.TpScrLoadStr("handoff_script", `def main(name, param, incoming, flags):
    return handoff(param, incoming, flags)
`, atmi.NDRX_TPSCR_FLAT); nil != err {
		return err
	}
	// An idle VM is used from another goroutine; no OS thread locking is needed.
	done := make(chan error, 1)
	go func() {
		var result error
		defer func() { done <- result }()
		started := time.Now()
		var count int64
		for count = 0; count < 1 || time.Since(started) < duration; count++ {
			data, err := ctx.NewCarray([]byte("handoff"))
			if nil != err {
				result = err
				return
			}
			output, err := vm.TpScrExec("handoff_script", nil, data, 0)
			if nil != err {
				result = err
			} else if nil == output {
				result = fmt.Errorf("thread handoff returned no buffer")
			} else {
				value, castErr := ctx.CastToCarray(output.GetBuf())
				if nil != castErr || string(value.GetBytes()) != "handoff" {
					result = fmt.Errorf("thread handoff result: %v", castErr)
				}
			}
			if nil != output && output.GetBuf() != data.GetBuf() {
				ctx.TpFree(output.GetBuf())
			}
			ctx.TpFree(data.GetBuf())
			if nil != result {
				return
			}
		}
		fmt.Printf("Thread handoff: %d iterations in %.3fs\n", count, time.Since(started).Seconds())
		result = vm.TpScrUninit(0)
		vm = nil
	}()
	return <-done
}

// Check that bytecode survives its original context and can be loaded elsewhere.
// @return Error if a released VM remains usable
func test_context_release() error {
	duration, durationErr := testutil.Duration()
	if nil != durationErr {
		return durationErr
	}
	var i int64
	started := time.Now()
	for i = 0; i < 1 || time.Since(started) < duration; i++ {
		ctx, err := atmi.NewATMICtx()
		if nil != err {
			return err
		}
		vm, err := ctx.TpScrInit(nil, 0)
		if nil != err {
			ctx.FreeATMICtx()
			return err
		}
		code, err := vm.TpScrComp(`def main(name, param, incoming, flags):
    return incoming
`, 0)
		if nil != err {
			ctx.FreeATMICtx()
			return err
		}
		ctx.FreeATMICtx()
		if _, err := vm.TpScrExec("echo", nil, nil, 0); nil == err || err.Code() != atmi.TPEINVAL {
			return fmt.Errorf("released VM exec: expected TPEINVAL, got %v", err)
		}
		if err := vm.TpScrUninit(0); nil == err || err.Code() != atmi.TPEINVAL {
			return fmt.Errorf("released VM uninit: expected TPEINVAL, got %v", err)
		}
		next, err := atmi.NewATMICtx()
		if nil != err {
			return err
		}
		again, err := next.TpScrInit(nil, 0)
		if nil != err {
			next.FreeATMICtx()
			return err
		}
		if err := again.TpScrLoad("echo", code, atmi.NDRX_TPSCR_FLAT); nil != err {
			next.FreeATMICtx()
			return err
		}
		output, err := again.TpScrExec("echo", nil, nil, 0)
		next.FreeATMICtx()
		if nil != err || nil != output {
			return fmt.Errorf("bytecode after context release: %v", err)
		}
	}
	fmt.Printf("Context release: %d iterations in %.3fs\n", i, time.Since(started).Seconds())
	return nil
}

// Check context ownership from a scripting callback: freeing the owner while
// its VM executes is ignored, and a VM created through the callback context
// belongs to the owner and is released with it.
// @return Error if a lifecycle rule is violated
func test_callback_lifecycle() error {
	duration, durationErr := testutil.Duration()
	if nil != durationErr {
		return durationErr
	}
	var i int64
	started := time.Now()
	for i = 0; i < 1 || time.Since(started) < duration; i++ {
		ctx, err := atmi.NewATMICtx()
		if nil != err {
			return err
		}
		vm, err := ctx.TpScrInit(nil, 0)
		if nil != err {
			ctx.FreeATMICtx()
			return err
		}
		helper, err := ctx.TpScrInit(nil, 0)
		if nil != err {
			ctx.FreeATMICtx()
			return err
		}
		var child *atmi.TpScrVm
		err = vm.TpScrRegCB("lifecycle", func(call *atmi.ATMICtx, vm *atmi.TpScrVm,
			name string, param atmi.TypedBuffer, input atmi.TypedBuffer,
			user interface{}, flags int64) (atmi.TypedBuffer, atmi.ATMIError) {
			if call != ctx {
				return nil, atmi.NewCustomATMIError(atmi.TPESYSTEM,
					"callback context is not the VM owner")
			}
			call.FreeATMICtx()
			if err := helper.TpScrSetError(1, "idle helper"); nil != err {
				return nil, atmi.NewCustomATMIError(atmi.TPESYSTEM,
					"idle VM released by an ignored context free")
			}
			var err atmi.ATMIError
			child, err = call.TpScrInit(nil, 0)
			return nil, err
		}, nil, 0)
		if nil != err {
			ctx.FreeATMICtx()
			return err
		}
		if err := vm.TpScrLoadStr("lifecycle_script", `def main(name, param, incoming, flags):
    rc, output = lifecycle(None, None, flags)
    return rc
`, atmi.NDRX_TPSCR_FLAT); nil != err {
			ctx.FreeATMICtx()
			return err
		}
		if _, err := vm.TpScrExec("lifecycle_script", nil, nil, 0); nil != err {
			ctx.FreeATMICtx()
			return fmt.Errorf("callback lifecycle exec: %v", err)
		}
		if err := vm.TpScrSetError(1, "still open"); nil != err {
			ctx.FreeATMICtx()
			return fmt.Errorf("VM unusable after ignored free: %v", err)
		}
		ctx.FreeATMICtx()
		if _, err := child.TpScrExec("lifecycle_script", nil, nil, 0); nil == err || err.Code() != atmi.TPEINVAL {
			return fmt.Errorf("callback-created VM survived owner free: %v", err)
		}
		if _, err := vm.TpScrExec("lifecycle_script", nil, nil, 0); nil == err || err.Code() != atmi.TPEINVAL {
			return fmt.Errorf("VM survived owner free: %v", err)
		}
		if err := helper.TpScrSetError(1, "closed"); nil == err || err.Code() != atmi.TPEINVAL {
			return fmt.Errorf("helper VM survived owner free: %v", err)
		}
	}
	fmt.Printf("Callback lifecycle: %d iterations in %.3fs\n", i, time.Since(started).Seconds())
	return nil
}

// Check callbacks run by VM shutdown (a script finalizer): during
// FreeATMICtx no new VM can be created on the closing context and a nested
// free is ignored; during a plain TpScrUninit a VM created by the callback
// belongs to the owner and is released with it.
// @return Error if a shutdown callback breaks a lifecycle rule
func test_shutdown_callback() error {
	duration, durationErr := testutil.Duration()
	if nil != durationErr {
		return durationErr
	}
	hook := `import weakref
class Holder:
    pass
holder = Holder()
weakref.finalize(holder, shutdown_cb, None, None, 0)
`
	var i int64
	started := time.Now()
	for i = 0; i < 1 || time.Since(started) < duration; i++ {
		for _, closing := range []bool{true, false} {
			ctx, err := atmi.NewATMICtx()
			if nil != err {
				return err
			}
			vm, err := ctx.TpScrInit(nil, 0)
			if nil != err {
				ctx.FreeATMICtx()
				return err
			}
			ran := false
			var child *atmi.TpScrVm
			var childErr atmi.ATMIError
			var nestedErr atmi.ATMIError
			err = vm.TpScrRegCB("shutdown_cb", func(call *atmi.ATMICtx, cbvm *atmi.TpScrVm,
				name string, param atmi.TypedBuffer, input atmi.TypedBuffer,
				user interface{}, flags int64) (atmi.TypedBuffer, atmi.ATMIError) {
				ran = true
				call.FreeATMICtx()
				nestedErr = cbvm.TpScrSetError(1, "in shutdown")
				child, childErr = call.TpScrInit(nil, 0)
				return nil, nil
			}, nil, 0)
			if nil == err {
				err = vm.TpScrLoadStr("go_shutdown_hook", hook, atmi.NDRX_TPSCR_PACKAGE)
			}
			if nil != err {
				ctx.FreeATMICtx()
				return err
			}
			if closing {
				ctx.FreeATMICtx()
			} else if err := vm.TpScrUninit(0); nil != err {
				ctx.FreeATMICtx()
				return err
			}
			if !ran {
				ctx.FreeATMICtx()
				return fmt.Errorf("shutdown callback did not run (closing %t)", closing)
			}
			if nil == nestedErr || nestedErr.Code() != atmi.TPEINVAL {
				ctx.FreeATMICtx()
				return fmt.Errorf("VM in shutdown remained usable: %v", nestedErr)
			}
			if closing {
				if nil == childErr || childErr.Code() != atmi.TPEPROTO {
					return fmt.Errorf("VM created on a closing context: %v", childErr)
				}
				if _, err := vm.TpScrExec("x", nil, nil, 0); nil == err || err.Code() != atmi.TPEINVAL {
					return fmt.Errorf("VM survived context free: %v", err)
				}
				continue
			}
			if nil != childErr {
				ctx.FreeATMICtx()
				return fmt.Errorf("shutdown callback could not create a VM: %v", childErr)
			}
			if err := child.TpScrSetError(1, "child open"); nil != err {
				ctx.FreeATMICtx()
				return fmt.Errorf("callback-created VM unusable: %v", err)
			}
			ctx.FreeATMICtx()
			if err := child.TpScrSetError(1, "closed"); nil == err || err.Code() != atmi.TPEINVAL {
				return fmt.Errorf("callback-created VM survived owner free: %v", err)
			}
		}
	}
	fmt.Printf("Shutdown callback: %d iterations in %.3fs\n", i, time.Since(started).Seconds())
	return nil
}

// Check that a VM callback reached from another context's script is rejected,
// so its XATMI work never runs under the foreign context.
// @return Error if the callback runs for another context
func test_foreign_callback() error {
	duration, durationErr := testutil.Duration()
	if nil != durationErr {
		return durationErr
	}
	var i int64
	started := time.Now()
	for i = 0; i < 1 || time.Since(started) < duration; i++ {
		ctxA, err := atmi.NewATMICtx()
		if nil != err {
			return err
		}
		ctxB, err := atmi.NewATMICtx()
		if nil != err {
			ctxA.FreeATMICtx()
			return err
		}
		vmA, err := ctxA.TpScrInit(nil, 0)
		if nil != err {
			ctxA.FreeATMICtx()
			ctxB.FreeATMICtx()
			return err
		}
		ran := false
		err = vmA.TpScrRegCB("owner_cb", func(call *atmi.ATMICtx, vm *atmi.TpScrVm,
			name string, param atmi.TypedBuffer, input atmi.TypedBuffer,
			user interface{}, flags int64) (atmi.TypedBuffer, atmi.ATMIError) {
			ran = true
			return nil, nil
		}, nil, 0)
		if nil == err {
			err = vmA.TpScrLoadStr("go_shared_owner", "VALUE = 1\n", atmi.NDRX_TPSCR_PACKAGE)
		}
		vmB, errB := ctxB.TpScrInit(nil, 0)
		if nil == err {
			err = errB
		}
		if nil == err {
			err = vmB.TpScrLoadStr("foreign", `def main(name, param, incoming, flags):
    import go_shared_owner
    rc, output = go_shared_owner.owner_cb(None, None, flags)
    return rc
`, atmi.NDRX_TPSCR_FLAT)
		}
		if nil != err {
			ctxA.FreeATMICtx()
			ctxB.FreeATMICtx()
			return err
		}
		_, execErr := vmB.TpScrExec("foreign", nil, nil, 0)
		ctxB.FreeATMICtx()
		ctxA.FreeATMICtx()
		if nil == execErr || ran {
			return fmt.Errorf("foreign context callback was not rejected (ran %t): %v", ran, execErr)
		}
	}
	fmt.Printf("Foreign callback: %d iterations in %.3fs\n", i, time.Since(started).Seconds())
	return nil
}

// Run one context and VM per goroutine in parallel.
// @return First worker error
func test_parallel_vms() error {
	duration, durationErr := testutil.Duration()
	if nil != durationErr {
		return durationErr
	}
	const workers = 4
	done := make(chan error, workers)
	started := time.Now()
	for w := 0; w < workers; w++ {
		go func(w int) {
			ctx, err := atmi.NewATMICtx()
			if nil != err {
				done <- err
				return
			}
			defer ctx.FreeATMICtx()
			vm, err := ctx.TpScrInit(nil, 0)
			if nil != err {
				done <- err
				return
			}
			defer vm.TpScrUninit(0)
			if err := vm.TpScrLoadStr("worker", `def main(name, param, incoming, flags):
    return {'buftype': 'CARRAY', 'data': incoming['data'] + str(flags).encode()}
`, atmi.NDRX_TPSCR_FLAT); nil != err {
				done <- err
				return
			}
			var i int64
			for i = 0; i < 1 || time.Since(started) < duration; i++ {
				data, err := ctx.NewCarray([]byte(fmt.Sprintf("w%d:", w)))
				if nil != err {
					done <- err
					return
				}
				output, err := vm.TpScrExec("worker", nil, data, int64(w))
				if nil != err {
					done <- err
					return
				}
				got, err := ctx.CastToCarray(output.GetBuf())
				if nil != err {
					done <- err
					return
				}
				if want := fmt.Sprintf("w%d:%d", w, w); want != string(got.GetBytes()) {
					done <- fmt.Errorf("worker %d: got %q want %q", w, got.GetBytes(), want)
					return
				}
				if output.GetBuf() != data.GetBuf() {
					ctx.TpFree(output.GetBuf())
				}
				ctx.TpFree(data.GetBuf())
			}
			fmt.Printf("Parallel VM worker %d: %d iterations in %.3fs\n", w, i,
				time.Since(started).Seconds())
			done <- nil
		}(w)
	}
	var first error
	for w := 0; w < workers; w++ {
		if err := <-done; nil != err && nil == first {
			first = err
		}
	}
	return first
}

// Check scripting callbacks invoked from an XATMI service process.
// @return Error if the server output differs from the callback result
func test_server() error {
	duration, durationErr := testutil.Duration()
	if nil != durationErr {
		return durationErr
	}
	var started time.Time
	var i int64
	ac, err := atmi.NewATMICtx()
	if nil != err {
		return err
	}
	defer ac.FreeATMICtx()
	started = time.Now()
	for i = 0; i < 1 || time.Since(started) < duration; i++ {
		data, err := ac.NewCarray([]byte("request"))
		if nil != err {
			return err
		}
		if _, err := ac.TpCall("GOSCRIPT", data, 0); nil != err {
			return err
		}
		if "script-server-reply" != string(data.GetBytes()) {
			return fmt.Errorf("server callback returned %q", data.GetBytes())
		}
		ac.TpFree(data.GetBuf())
	}
	fmt.Printf("test_server loop 2: %d iterations in %.3fs\n", i, time.Since(started).Seconds())
	return nil
}

// Run the scripting lifecycle, callbacks and buffer ownership integration case.
func main() {
	flag.Parse()
	if 0 != flag.NArg() {
		fmt.Fprintln(os.Stderr, "Usage: client [-scripting]")
		os.Exit(2)
	}
	if !*scripting {
		fmt.Println("SKIPPED: scripting requires -scripting and separately licensed endurox-python")
		return
	}
	atmi.RuntimeInit()
	fmt.Println("*** Scripting lifecycle, callbacks and buffer ownership (duration from NDRX_TEST_DURATATION; default 10s)")
	if err := test_scripting(); nil != err {
		fmt.Fprintf(os.Stderr, "TESTERROR: 16_scripting: %v\n", err)
		os.Exit(1)
	}
	if err := test_context_release(); nil != err {
		fmt.Fprintf(os.Stderr, "TESTERROR: 16_scripting context release: %v\n", err)
		os.Exit(1)
	}
	if err := test_callback_lifecycle(); nil != err {
		fmt.Fprintf(os.Stderr, "TESTERROR: 16_scripting callback lifecycle: %v\n", err)
		os.Exit(1)
	}
	if err := test_shutdown_callback(); nil != err {
		fmt.Fprintf(os.Stderr, "TESTERROR: 16_scripting shutdown callback: %v\n", err)
		os.Exit(1)
	}
	if err := test_foreign_callback(); nil != err {
		fmt.Fprintf(os.Stderr, "TESTERROR: 16_scripting foreign callback: %v\n", err)
		os.Exit(1)
	}
	if err := test_parallel_vms(); nil != err {
		fmt.Fprintf(os.Stderr, "TESTERROR: 16_scripting parallel VMs: %v\n", err)
		os.Exit(1)
	}
	if err := test_server(); nil != err {
		fmt.Fprintln(os.Stderr, "TESTERROR: scripting server:", err)
		os.Exit(1)
	}
	fmt.Println("16_scripting passed")
}

/* vim: set ts=4 sw=4 et smartindent: */
