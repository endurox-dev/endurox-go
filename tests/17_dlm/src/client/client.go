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
	"bytes"
	"fmt"
	atmi "github.com/endurox-dev/endurox-go/v9"
	"github.com/endurox-dev/endurox-go/v9/tests/testutil"
	"os"
	"time"
	"ubftab"
)

// Exercise DLM via a Go server, native async replies and listing.
// @return failure details
func test_dlm() error {
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
	data, err := ac.NewUBF(8192)
	if nil != err {
		return err
	}
	defer ac.TpFree(data.GetBuf())
	effective, attempt, block, err := ac.TpDLMToutGet(5000)
	if nil != err || 5000 != effective || 0 >= attempt || 0 >= block {
		return fmt.Errorf("DLM timeouts: %d/%d/%d, %v", effective, attempt, block, err)
	}
	count, err := ac.TpDLMAttemptsGet()
	if nil != err || 0 > count {
		return fmt.Errorf("DLM attempt limit: %d, %v", count, err)
	}
	payload := []byte{0, 1, 255, 9}
	started = time.Now()
	for i = 0; i < 1 || time.Since(started) < duration; i++ {
		key := fmt.Sprintf("GO17_%d", i)
		if err := ac.BInit(data, 8192); nil != err {
			return err
		}
		if err := data.BChg(ubftab.EX_DLM_OP, 0, atmi.NDRX_DLM_OP_DLMCMD); nil != err {
			return err
		}
		if err := data.BChg(ubftab.EX_DLM_CMD, 0, atmi.NDRX_DLM_CMD_SET); nil != err {
			return err
		}
		if err := data.BChg(ubftab.EX_DLM_KEY, 0, key); nil != err {
			return err
		}
		if err := data.BChg(ubftab.EX_DLM_FLAGS, 0, int64(atmi.NDRX_DLM_FL_FLATBLOB)); nil != err {
			return err
		}
		if err := data.BChg(ubftab.EX_DLM_DATA_BLOB, 0, payload); nil != err {
			return err
		}
		if _, err := ac.TpCall("GODLM", data, 0); nil != err {
			return err
		}
		rval, uerr := data.BGetInt64(ubftab.EX_DLM_RVAL1, 0)
		if nil != uerr || 0 != rval {
			return fmt.Errorf("SET result %d: %v", rval, uerr)
		}
		if err := ac.BInit(data, 8192); nil != err {
			return err
		}
		if err := data.BChg(ubftab.EX_DLM_OP, 0, atmi.NDRX_DLM_OP_DLMCMD); nil != err {
			return err
		}
		if err := data.BChg(ubftab.EX_DLM_CMD, 0, atmi.NDRX_DLM_CMD_GET); nil != err {
			return err
		}
		if err := data.BChg(ubftab.EX_DLM_KEY, 0, key); nil != err {
			return err
		}
		if err := data.BChg(ubftab.EX_DLM_FLAGS, 0, int64(atmi.NDRX_DLM_FL_FLATBLOB)); nil != err {
			return err
		}
		var ctl atmi.TpDlmCtl
		if err := ctl.SetDLMSpace("GOTEST"); nil != err {
			return err
		}
		if err := ac.TpDLMMkCall(data, &ctl); nil != err {
			return err
		}
		cd, err := ac.TpDLMACall(ctl.GetSvcName(), data, 0)
		if nil != err {
			return err
		}
		if _, err := ac.TpGetRply(&cd, data, 0); nil != err {
			return err
		}
		value, uerr := data.BGetByteArr(ubftab.EX_DLM_DATA_BLOB, 0)
		if nil != uerr || !bytes.Equal(value, payload) {
			return fmt.Errorf("GET payload %v: %v", value, uerr)
		}
		if err := ac.BInit(data, 8192); nil != err {
			return err
		}
		if err := data.BChg(ubftab.EX_DLM_OP, 0, atmi.NDRX_DLM_OP_DLMCMD); nil != err {
			return err
		}
		if err := data.BChg(ubftab.EX_DLM_CMD, 0, atmi.NDRX_DLM_CMD_DEL); nil != err {
			return err
		}
		if err := data.BChg(ubftab.EX_DLM_FLAGS, 0, int64(atmi.NDRX_DLM_FL_FLATBLOB)); nil != err {
			return err
		}
		if err := data.BChg(ubftab.EX_DLM_KEY, 0, key); nil != err {
			return err
		}
		ctl = atmi.TpDlmCtl{}
		if err := ctl.SetDLMSpace("GOTEST"); nil != err {
			return err
		}
		if err := ac.TpDLMMkCall(data, &ctl); nil != err {
			return err
		}
		if _, err := ac.TpDLMCall(ctl.GetSvcName(), data, 0); nil != err {
			return err
		}
		rval, uerr = data.BGetInt64(ubftab.EX_DLM_RVAL1, 0)
		if nil != uerr || 0 != rval {
			return fmt.Errorf("DEL result %d: %v", rval, uerr)
		}
	}
	fmt.Printf("test_dlm loop 1: %d iterations in %.3fs\n", i, time.Since(started).Seconds())
	if err := ac.BInit(data, 8192); nil != err {
		return err
	}
	if err := data.BChg(ubftab.EX_DLM_OP, 0, atmi.NDRX_DLM_OP_LSSTORE); nil != err {
		return err
	}
	if err := data.BChg(ubftab.EX_DLM_KEY, 0, "GO17_.*"); nil != err {
		return err
	}
	if err := data.BChg(ubftab.EX_DLM_OPFLAGS, 0, int64(atmi.NDRX_DLM_FL_REGEX|atmi.NDRX_DLM_FL_NODATA)); nil != err {
		return err
	}
	var ctl atmi.TpDlmCtl
	if err := ctl.SetDLMSpace("GOTEST"); nil != err {
		return err
	}
	if err := ac.TpDLMMkCall(data, &ctl); nil != err {
		return err
	}
	cd, err := ac.TpDLMConnect(ctl.GetSvcName(), data, atmi.TPRECVONLY)
	if nil != err {
		return err
	}
	var event int64
	err = ac.TpRecv(cd, data, 0, &event)
	if nil == err || atmi.TPEEVENT != err.Code() || atmi.TPEV_SVCSUCC != event {
		ac.TpDiscon(cd)
		return fmt.Errorf("listing terminal result %d: %v", event, err)
	}
	return nil
}

// Run the DLM integration case.
func main() {
	atmi.RuntimeInit()
	if err := test_dlm(); nil != err {
		fmt.Fprintln(os.Stderr, "TESTERROR:", err)
		os.Exit(1)
	}
	fmt.Println("17_dlm passed")
}
