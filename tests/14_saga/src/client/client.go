/**
 * @brief Saga and transaction argument validation
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
	"os"
	"time"

	atmi "github.com/endurox-dev/endurox-go/v9"
	"github.com/endurox-dev/endurox-go/v9/tests/testutil"
)

// Check Saga payload and transaction identifier validation.
// @return Error if a malformed request is accepted
func test_saga() error {
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
		var ctl atmi.TpSagaCtl
		if err := ctl.SetPayload([]byte{1, 2, 3}); nil != err {
			return err
		}
		if err := ctl.SetPayload(make([]byte, atmi.NDRX_SAGA_PAYLOAD_MAX+1)); nil == err || atmi.TPEINVAL != err.Code() {
			return fmt.Errorf("oversized Saga payload: %v", err)
		}
		if 3 != ctl.GetPayloadLen() || 1 != ctl.GetPayload()[0] {
			return fmt.Errorf("invalid setter changed native payload")
		}
		alias := ctl
		if err := alias.SetPayloadLen(6); nil != err {
			return err
		}
		payload := ctl.GetPayload()
		if 6 != len(payload) || 0 != payload[5] {
			return fmt.Errorf("native payload view did not grow safely")
		}
		if err := alias.SetPayload(nil); nil != err {
			return err
		}
		if 0 != ctl.GetPayloadLen() {
			return fmt.Errorf("shared control did not clear payload")
		}
		if _, err := ac.TpXIDFromString("bad\x00xid"); nil == err || atmi.TPEINVAL != err.Code() {
			return fmt.Errorf("embedded NUL XID: %v", err)
		}
		var status atmi.TpTxStatus
		if err := status.SetTranIDString("not-base64"); nil == err || atmi.TPEINVAL != err.Code() {
			return fmt.Errorf("invalid suspend token: %v", err)
		}
	}
	fmt.Printf("test_saga loop 1: %d iterations in %.3fs\n", i, time.Since(started).Seconds())
	return nil
}

// Exercise real transaction-manager registration and completion callbacks.
// @return Error if a transaction, recovery API or callback result differs
func test_transactions() error {
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
	if err := ac.TpOpen(); nil != err {
		return err
	}
	defer ac.TpClose()
	data, err := ac.NewUBF(8192)
	if nil != err {
		return err
	}
	defer ac.TpFree(data.GetBuf())
	started = time.Now()
	for i = 0; i < 3 || i%3 != 0 || time.Since(started) < duration; i++ {
		if err := ac.TpBegin(60, atmi.TPSAGALRA); nil != err {
			return err
		}
		var first atmi.TpSagaCtl
		if err := first.SetTimeout(40); nil != err {
			return err
		}
		if err := first.SetSagaName("GO"); nil != err {
			return err
		}
		if err := first.SetSagaKey(fmt.Sprintf("client-%d", i)); nil != err {
			return err
		}
		if err := first.SetPayload([]byte("client")); nil != err {
			return err
		}
		if err := first.SetContentType("text/plain"); nil != err {
			return err
		}
		if err := first.SetCommitURL("xatmi://GOCOMMIT"); nil != err {
			return err
		}
		if err := first.SetRollbackURL("xatmi://GOROLLBACK"); nil != err {
			return err
		}
		if err := ac.TpSagaEnlist(&first, 0); nil != err {
			return err
		}
		if err := ac.TpSagaEnd(&first, 0); nil != err {
			return err
		}
		mode := "normal"
		if 2 == i%3 {
			mode = "leave"
		}
		if err := data.BChg(atmi.EX_DATA_STR, 0, fmt.Sprintf("server-%d", i)); nil != err {
			return err
		}
		if err := data.BChg(atmi.EX_DATA_STR, 1, mode); nil != err {
			return err
		}
		if _, err := ac.TpCall("GOREGISTER", data, 0); nil != err {
			return err
		}
		second, uerr := data.BGetString(atmi.EX_DATA_STR, 0)
		if nil != uerr {
			return uerr
		}
		global, uerr := data.BGetString(atmi.EX_DATA_STR, 1)
		if nil != uerr {
			return uerr
		}
		xid, err := ac.TpXIDFromString(global)
		if nil != err {
			return err
		}
		encoded, err := ac.TpXIDToString(xid)
		if nil != err || global != encoded {
			return fmt.Errorf("global XID round trip: %v", err)
		}
		var status atmi.TpTxStatus
		if err := status.SetFlags(atmi.TPTXS_XID); nil != err {
			return err
		}
		if err := status.SetXID(*xid); nil != err {
			return err
		}
		if err := ac.TpTxStatus(&status, 0); nil != err {
			return err
		}
		if atmi.TPTXS_STATUS_ACTIVE != status.GetStatus() {
			return fmt.Errorf("unexpected active state %d", status.GetStatus())
		}
		token, err := ac.TpSuspendString(0)
		if nil != err {
			return err
		}
		status = atmi.TpTxStatus{}
		if err := status.SetFlags(atmi.TPTXS_TRANID); nil != err {
			return err
		}
		if err := status.SetTranIDString(token); nil != err {
			return err
		}
		if err := ac.TpTxStatus(&status, 0); nil != err {
			return err
		}
		nativeToken := status.GetTranID()
		if err := status.SetTranID(nativeToken); nil != err {
			return err
		}
		if err := ac.TpResume(&nativeToken, 0); nil != err {
			return err
		}
		status = atmi.TpTxStatus{}
		if err := status.SetFlags(atmi.TPTXS_BRXID); nil != err {
			return err
		}
		if err := status.SetBranchXID(first.GetBranchXID()); nil != err {
			return err
		}
		if err := status.SetTimeout(90); nil != err {
			return err
		}
		if err := ac.TpTxRenew(&status, 0); nil != err {
			return err
		}
		if 0 >= status.GetTimeoutLeft() || 30 < status.GetTimeoutLeft() || atmi.TPTXS_STATUS_ACTIVE != status.GetStatus() {
			return fmt.Errorf("renewal result: %+v", status)
		}
		if "leave" != mode {
			var lookup atmi.TpSagaCtl
			if err := lookup.SetBranchXID(second); nil != err {
				return err
			}
			if err := ac.TpSagaGet(&lookup, 0); nil != err {
				return err
			}
			if 30 != lookup.GetTimeout() || "GO" != lookup.GetSagaName() || "xatmi://GOCOMMIT" != lookup.GetCommitURL() {
				return fmt.Errorf("participant lookup: %+v", lookup)
			}
			if err := lookup.SetPayload([]byte("updated")); nil != err {
				return err
			}
			if err := ac.TpSagaUpdate(&lookup, 0); nil != err {
				return err
			}
			if err := lookup.SetCommitURL("xatmi://GOCOMMIT2"); nil != err {
				return err
			}
			if err := ac.TpSagaUpdate(&lookup, atmi.TPSAGA_URLSONLY); nil != err {
				return err
			}
			if err := ac.TpSagaGet(&lookup, 0); nil != err {
				return err
			}
			if "xatmi://GOCOMMIT2" != lookup.GetCommitURL() {
				return fmt.Errorf("URL update missing")
			}
			if err := lookup.SetPayload(nil); nil != err {
				return err
			}
		}
		if err := first.SetPayload(nil); nil != err {
			return err
		}
		expected := byte(atmi.NDRX_SAGA_ACTION_COMMIT)
		if 1 == i%3 {
			expected = byte(atmi.NDRX_SAGA_ACTION_ROLLBACK)
			if err := ac.TpAbort(0); nil != err {
				return err
			}
		} else {
			if err := ac.TpCommit(0); nil != err {
				return err
			}
		}
		var firstSequence int64
		for index, branch := range []string{first.GetBranchXID(), second} {
			if err := data.BChg(atmi.EX_DATA_STR, 0, branch); nil != err {
				return err
			}
			if _, err := ac.TpCall("GOSTATE", data, atmi.TPNOTRAN); nil != err {
				return err
			}
			count, err := data.BGetInt64(atmi.EX_NERROR_CODE, 0)
			if nil != err {
				return err
			}
			if 1 == index && "leave" == mode {
				if 0 != count {
					return fmt.Errorf("removed participant received callback")
				}
				continue
			}
			action, err := data.BGetByte(atmi.TMSAGA_ACTION, 0)
			if nil != err {
				return err
			}
			payload, err := data.BGetString(atmi.EX_DATA_STR, 1)
			if nil != err {
				return err
			}
			sequence, err := data.BGetInt64(atmi.TMTXTRYCNT, 0)
			if nil != err {
				return err
			}
			if 1 != count || expected != action {
				return fmt.Errorf("callback count/action %d/%c", count, action)
			}
			if 0 == index {
				firstSequence = sequence
				if "client" != payload {
					return fmt.Errorf("client payload %q", payload)
				}
			} else {
				if "updated" != payload {
					return fmt.Errorf("updated payload %q", payload)
				}
				if byte(atmi.NDRX_SAGA_ACTION_COMMIT) == expected && firstSequence >= sequence {
					return fmt.Errorf("commit order reversed")
				}
				if byte(atmi.NDRX_SAGA_ACTION_ROLLBACK) == expected && firstSequence <= sequence {
					return fmt.Errorf("rollback order not reversed")
				}
			}
		}
	}
	fmt.Printf("test_transactions loop 2: %d iterations in %.3fs\n", i, time.Since(started).Seconds())
	return nil
}

// Run the saga and transaction argument validation integration case.
func main() {
	if 1 != len(os.Args) {
		fmt.Fprintln(os.Stderr, "Usage: client")
		os.Exit(2)
	}
	atmi.RuntimeInit()
	fmt.Println("*** Saga and transaction argument validation (duration from NDRX_TEST_DURATATION; default 10s)")
	if err := test_saga(); nil != err {
		fmt.Fprintf(os.Stderr, "TESTERROR: 14_saga: %v\n", err)
		os.Exit(1)
	}
	if err := test_transactions(); nil != err {
		fmt.Fprintln(os.Stderr, "TESTERROR: Saga services:", err)
		os.Exit(1)
	}
	fmt.Println("14_saga passed")
}

/* vim: set ts=4 sw=4 et smartindent: */
