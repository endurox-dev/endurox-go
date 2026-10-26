/**
 * @brief DLM commands, asynchronous replies and listing
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
	"fmt"
	atmi "github.com/endurox-dev/endurox-go/v9"
	"os"
)

// Forward a request through DLM preparation and retry handling.
// @param ac server context
// @param svc incoming UBF request
func GODLM(ac *atmi.ATMICtx, svc *atmi.TPSVCINFO) {
	data, err := ac.CastToUBF(&svc.Data)
	if nil != err {
		ac.TpReturn(atmi.TPFAIL, 0, &svc.Data, 0)
		return
	}
	var ctl atmi.TpDlmCtl
	if err := ctl.SetDLMSpace("GOTEST"); nil != err {
		ac.TpLogError("TESTERROR: native control: %v", err)
		ac.TpReturn(atmi.TPFAIL, 0, data, 0)
		return
	}
	if err := ac.TpDLMMkCall(data, &ctl); nil != err {
		ac.TpLogError("TESTERROR: DLM preparation: %v", err)
		ac.TpReturn(atmi.TPFAIL, 0, data, 0)
		return
	}
	if _, err := ac.TpDLMCall(ctl.GetSvcName(), data, 0); nil != err {
		ac.TpLogError("TESTERROR: DLM call: %v", err)
		ac.TpReturn(atmi.TPFAIL, 0, data, 0)
		return
	}
	ac.TpReturn(atmi.TPSUCCESS, 0, data, 0)
}

// @param ac server context
// @return startup status
func Init(ac *atmi.ATMICtx) int {
	if err := ac.TpAdvertise("GODLM", "GODLM", GODLM); nil != err {
		return atmi.FAIL
	}
	return atmi.SUCCEED
}

// @param ac server context
func Uninit(ac *atmi.ATMICtx) {}

// Run the XATMI service dispatcher.
func main() {
	atmi.RuntimeInit()
	ac, err := atmi.NewATMICtx()
	if nil != err {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := ac.TpRun(Init, Uninit); nil != err {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
