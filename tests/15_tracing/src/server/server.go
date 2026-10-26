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

// Process the GOTRACE request through the native binding.
// @param ac server context
// @param svc service request
func GOTRACE(ac *atmi.ATMICtx, svc *atmi.TPSVCINFO) {
	parent, err := ac.TpTraceContextGet(0)
	if nil != err {
		ac.TpReturn(atmi.TPFAIL, 0, &svc.Data, 0)
		return
	}
	state, err := ac.TpTraceStateGet(0)
	if nil != err {
		ac.TpReturn(atmi.TPFAIL, 0, &svc.Data, 0)
		return
	}
	reply, err := ac.NewString(parent + "\n" + state)
	if nil != err {
		ac.TpReturn(atmi.TPFAIL, 0, &svc.Data, 0)
		return
	}
	ac.TpReturn(atmi.TPSUCCESS, 0, reply, 0)
}

// Initialize the functional test service.
// @param ac server context
// @return SUCCEED or FAIL
func Init(ac *atmi.ATMICtx) int {
	if err := ac.TpAdvertise("GOTRACE", "GOTRACE", GOTRACE); nil != err {
		fmt.Println("TESTERROR:", err)
		return atmi.FAIL
	}
	return atmi.SUCCEED
}

// Release service resources at shutdown.
// @param ac server context
func Uninit(ac *atmi.ATMICtx) {
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
