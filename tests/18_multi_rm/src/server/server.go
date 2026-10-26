/**
 * @brief Multiple resource managers and manual branch enlistment
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

// Enlist the manual RM and return both native connection identities.
// @param ac server context
// @param svc incoming request
func GORM(ac *atmi.ATMICtx, svc *atmi.TPSVCINFO) {
	data, err := ac.CastToUBF(&svc.Data)
	if nil != err {
		ac.TpReturn(atmi.TPFAIL, 0, &svc.Data, 0)
		return
	}
	if err := ac.TpRMStart("B", 0); nil != err {
		ac.TpReturn(atmi.TPFAIL, 0, data, 0)
		return
	}
	a, err := ac.TpGetConnRM(1)
	if nil != err {
		ac.TpReturn(atmi.TPFAIL, 0, data, 0)
		return
	}
	b, err := ac.TpGetConnRM(2)
	if nil != err || nil == a || nil == b || a == b {
		ac.TpReturn(atmi.TPFAIL, 0, data, 0)
		return
	}
	if err := data.BChg(atmi.TMTXTRYCNT, 0, int64(uintptr(a))); nil != err {
		ac.TpReturn(atmi.TPFAIL, 0, data, 0)
		return
	}
	if err := data.BChg(atmi.TMTXTRYMAXCNT, 0, int64(uintptr(b))); nil != err {
		ac.TpReturn(atmi.TPFAIL, 0, data, 0)
		return
	}
	if err := ac.TpRMEnd("B", 0); nil != err {
		ac.TpReturn(atmi.TPFAIL, 0, data, 0)
		return
	}
	ac.TpReturn(atmi.TPSUCCESS, 0, data, 0)
}

// @param ac server context
// @return startup status
func Init(ac *atmi.ATMICtx) int {
	if err := ac.TpOpen2(atmi.TPTXOPENALL); nil != err {
		fmt.Println(err)
		return atmi.FAIL
	}
	if err := ac.TpAdvertise("GORM", "GORM", GORM); nil != err {
		return atmi.FAIL
	}
	return atmi.SUCCEED
}

// @param ac server context
func Uninit(ac *atmi.ATMICtx) { ac.TpClose() }

// Run the XATMI server.
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
