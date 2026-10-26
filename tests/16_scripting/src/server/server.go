/**
 * @brief Buffer comparison, export and transport lengths
 *
 * @file server.go
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
	"fmt"
	atmi "github.com/endurox-dev/endurox-go/v9"
	"os"
)

// VM owned by the server lifecycle.
var vm *atmi.TpScrVm = nil

// Process the GOSCRIPT request through the native binding.
// @param ac server context
// @param svc service request
func GOSCRIPT(ac *atmi.ATMICtx, svc *atmi.TPSVCINFO) {
	reply, err := vm.TpScrExec("service", nil, &svc.Data, 37)
	if nil != err {
		ac.TpLogError("TESTERROR: server scripting: %v", err)
		ac.TpReturn(atmi.TPFAIL, 0, &svc.Data, 0)
		return
	}
	ac.TpReturn(atmi.TPSUCCESS, 0, reply, 0)
}

// Initialize the functional test service.
// @param ac server context
// @return SUCCEED or FAIL
func Init(ac *atmi.ATMICtx) int {
	var err atmi.ATMIError
	vm, err = ac.TpScrInit(nil, 0)
	if nil != err {
		fmt.Println("TESTERROR:", err)
		return atmi.FAIL
	}
	err = vm.TpScrRegCB("host", func(call *atmi.ATMICtx,
		vm *atmi.TpScrVm, name string, param atmi.TypedBuffer,
		input atmi.TypedBuffer, user interface{}, flags int64) (atmi.TypedBuffer, atmi.ATMIError) {
		data, err := call.CastToCarray(input.GetBuf())
		if nil != err {
			return nil, err
		}
		if 37 != flags || "request" != string(data.GetBytes()) {
			return nil, atmi.NewCustomATMIError(atmi.TPEINVAL,
				"server callback received wrong input or flags")
		}
		return call.NewCarray([]byte("script-server-reply"))
	}, nil, 0)
	if nil != err {
		fmt.Println("TESTERROR:", err)
		return atmi.FAIL
	}
	err = vm.TpScrLoadStr("service", `def main(name, param, incoming, flags):
    return host(param, incoming, flags)
`, atmi.NDRX_TPSCR_FLAT)
	if nil != err {
		fmt.Println("TESTERROR:", err)
		return atmi.FAIL
	}
	if err := ac.TpAdvertise("GOSCRIPT", "GOSCRIPT", GOSCRIPT); nil != err {
		fmt.Println("TESTERROR:", err)
		return atmi.FAIL
	}
	return atmi.SUCCEED
}

// Release service resources at shutdown.
// @param ac server context
func Uninit(ac *atmi.ATMICtx) {
	if nil != vm {
		if err := vm.TpScrUninit(0); nil != err {
			fmt.Println("TESTERROR:", err)
		}
	}
}

// Run the native XATMI server dispatcher.
func main() {
	atmi.RuntimeInit()
	ac, err := atmi.NewATMICtx()
	if nil != err {
		fmt.Fprintln(os.Stderr, "TESTERROR:", err)
		os.Exit(1)
	}
	if err := ac.TpRun(Init, Uninit); nil != err {
		fmt.Fprintln(os.Stderr, "TESTERROR:", err)
		os.Exit(1)
	}
}
