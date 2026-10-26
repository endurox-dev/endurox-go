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
	"github.com/endurox-dev/endurox-go/v9/tests/testutil"
	"os"
	"time"
	"ubftab"
)

// Exercise named opens and manual enlistment with a server. Both RMs use the
// null switch, which provides no native connection handle.
// @return failure details
func test_rms() error {
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
	if err := ac.TpOpen2(atmi.TPTXOPENALL); nil != err {
		return err
	}
	defer ac.TpClose()
	data, err := ac.NewUBF(4096)
	if nil != err {
		return err
	}
	defer ac.TpFree(data.GetBuf())
	if err := ac.TpRMClose("B", 0); nil != err {
		return err
	}
	if err := ac.TpRMStart("B", 0); nil == err || atmi.TPEPROTO != err.Code() {
		return fmt.Errorf("closed RM start result: %v", err)
	}
	if err := ac.TpRMOpen("B", 0); nil != err {
		return err
	}
	started = time.Now()
	for i = 0; i < 2 || i%2 != 0 || time.Since(started) < duration; i++ {
		if err := ac.TpBegin(60, 0); nil != err {
			return err
		}
		if err := ac.TpRMStart("B", 0); nil != err {
			return err
		}
		if _, err := ac.TpGetConnRM(1); nil == err || atmi.TPENOENT != err.Code() {
			return fmt.Errorf("RM A connection result: %v", err)
		}
		if _, err := ac.TpGetConnRM(2); nil == err || atmi.TPENOENT != err.Code() {
			return fmt.Errorf("RM B connection result: %v", err)
		}
		if err := data.BChg(ubftab.T_LONG_FLD, 0, int64(0)); nil != err {
			return err
		}
		if _, err := ac.TpCall("GORM", data, 0); nil != err {
			return err
		}
		enlisted, uerr := data.BGetInt64(ubftab.T_LONG_FLD, 0)
		if nil != uerr || 1 != enlisted {
			return fmt.Errorf("server RM B enlistment: %d, %v", enlisted, uerr)
		}
		if err := ac.TpRMEnd("B", 0); nil != err {
			return err
		}
		if 0 == i%2 {
			if err := ac.TpCommit(0); nil != err {
				return err
			}
		} else {
			if err := ac.TpAbort(0); nil != err {
				return err
			}
		}
	}
	fmt.Printf("test_rms loop 1: %d iterations in %.3fs\n", i, time.Since(started).Seconds())
	return nil
}

// Run the multiple-resource-manager integration case.
func main() {
	atmi.RuntimeInit()
	if err := test_rms(); nil != err {
		fmt.Fprintln(os.Stderr, "TESTERROR:", err)
		os.Exit(1)
	}
	fmt.Println("18_multi_rm passed")
}
