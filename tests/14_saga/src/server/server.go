/**
 * @brief Saga participant and completion services
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
	"sync"
)

// Completed callbacks retained until the client reads each record.
type actionRecord struct {
	action   byte
	payload  []byte
	sequence int64
}

var records = make(map[string]actionRecord)
var recordLock sync.Mutex
var sequence int64 = 0

// Register one participant in the service branch, then release its local section.
// @param ac server context
// @param svc request containing the participant key and optional leave mode
func GOREGISTER(ac *atmi.ATMICtx, svc *atmi.TPSVCINFO) {
	request, err := ac.CastToUBF(&svc.Data)
	if nil != err {
		ac.TpReturn(atmi.TPFAIL, 0, &svc.Data, 0)
		return
	}
	key, uerr := request.BGetString(atmi.EX_DATA_STR, 0)
	if nil != uerr {
		ac.TpReturn(atmi.TPFAIL, 0, request, 0)
		return
	}
	mode, uerr := request.BGetString(atmi.EX_DATA_STR, 1)
	if nil != uerr {
		ac.TpReturn(atmi.TPFAIL, 0, request, 0)
		return
	}
	var ctl atmi.TpSagaCtl
	if err := ctl.SetTimeout(30); nil != err {
		ac.TpReturn(atmi.TPFAIL, 0, request, 0)
		return
	}
	if err := ctl.SetSagaName("GO"); nil != err {
		ac.TpLogError("TESTERROR: native control: %v", err)
		ac.TpReturn(atmi.TPFAIL, 0, request, 0)
		return
	}
	if err := ctl.SetSagaKey(key); nil != err {
		ac.TpLogError("TESTERROR: native control: %v", err)
		ac.TpReturn(atmi.TPFAIL, 0, request, 0)
		return
	}
	if err := ctl.SetPayload([]byte("server")); nil != err {
		ac.TpLogError("TESTERROR: native control: %v", err)
		ac.TpReturn(atmi.TPFAIL, 0, request, 0)
		return
	}
	if err := ctl.SetContentType("text/plain"); nil != err {
		ac.TpLogError("TESTERROR: native control: %v", err)
		ac.TpReturn(atmi.TPFAIL, 0, request, 0)
		return
	}
	if err := ctl.SetCommitURL("xatmi://GOCOMMIT"); nil != err {
		ac.TpLogError("TESTERROR: native control: %v", err)
		ac.TpReturn(atmi.TPFAIL, 0, request, 0)
		return
	}
	if err := ctl.SetRollbackURL("xatmi://GOROLLBACK"); nil != err {
		ac.TpLogError("TESTERROR: native control: %v", err)
		ac.TpReturn(atmi.TPFAIL, 0, request, 0)
		return
	}
	if err := ac.TpSagaEnlist(&ctl, 0); nil != err {
		ac.TpLogError("TESTERROR: server registration: %v", err)
		ac.TpReturn(atmi.TPFAIL, 0, request, 0)
		return
	}
	duplicate := ctl
	if err := ac.TpSagaEnlist(&duplicate, 0); nil == err || atmi.TPEPROTO != err.Code() {
		ac.TpLogError("TESTERROR: duplicate registration: %v", err)
		ac.TpReturn(atmi.TPFAIL, 0, request, 0)
		return
	}
	if err := ac.TpSagaEnd(&ctl, 0); nil != err {
		ac.TpReturn(atmi.TPFAIL, 0, request, 0)
		return
	}
	if "leave" == mode {
		if err := ac.TpSagaLeave(&ctl, 0); nil != err {
			ac.TpReturn(atmi.TPFAIL, 0, request, 0)
			return
		}
	}
	if uerr := request.BChg(atmi.EX_DATA_STR, 0, ctl.GetBranchXID()); nil != uerr {
		ac.TpReturn(atmi.TPFAIL, 0, request, 0)
		return
	}
	xid := ctl.GetXID()
	global, err := ac.TpXIDToString(&xid)
	if nil != err {
		ac.TpReturn(atmi.TPFAIL, 0, request, 0)
		return
	}
	if uerr := request.BChg(atmi.EX_DATA_STR, 1, global); nil != uerr {
		ac.TpReturn(atmi.TPFAIL, 0, request, 0)
		return
	}
	if err := ctl.SetPayload(nil); nil != err {
		ac.TpReturn(atmi.TPFAIL, 0, request, 0)
		return
	}
	ac.TpReturn(atmi.TPSUCCESS, 0, request, 0)
}

// Record the action sent by tmsrv using the returned branch identifier.
// @param ac server context
// @param svc native Saga callback UBF
func GOACTION(ac *atmi.ATMICtx, svc *atmi.TPSVCINFO) {
	data, err := ac.CastToUBF(&svc.Data)
	if nil != err {
		ac.TpReturn(atmi.TPFAIL, 0, &svc.Data, 0)
		return
	}
	branch, uerr := data.BGetString(atmi.TMSAGA_BRXID, 0)
	if nil != uerr {
		ac.TpReturn(atmi.TPFAIL, 0, data, 0)
		return
	}
	action, uerr := data.BGetByte(atmi.TMSAGA_ACTION, 0)
	if nil != uerr {
		ac.TpReturn(atmi.TPFAIL, 0, data, 0)
		return
	}
	payload, uerr := data.BGetByteArr(atmi.TMSAGA_DATA, 0)
	if nil != uerr {
		ac.TpReturn(atmi.TPFAIL, 0, data, 0)
		return
	}
	recordLock.Lock()
	if _, exists := records[branch]; !exists {
		sequence++
		records[branch] = actionRecord{action: action, payload: payload, sequence: sequence}
	}
	recordLock.Unlock()
	ac.TpReturn(atmi.TPSUCCESS, 0, data, 0)
}

// Return and remove one callback record so load cases do not grow server state.
// @param ac server context
// @param svc query UBF containing a branch identifier
func GOSTATE(ac *atmi.ATMICtx, svc *atmi.TPSVCINFO) {
	data, err := ac.CastToUBF(&svc.Data)
	if nil != err {
		ac.TpReturn(atmi.TPFAIL, 0, &svc.Data, 0)
		return
	}
	branch, uerr := data.BGetString(atmi.EX_DATA_STR, 0)
	if nil != uerr {
		ac.TpReturn(atmi.TPFAIL, 0, data, 0)
		return
	}
	recordLock.Lock()
	record, exists := records[branch]
	delete(records, branch)
	recordLock.Unlock()
	count := int64(0)
	if exists {
		count = 1
	}
	if uerr := data.BChg(atmi.EX_NERROR_CODE, 0, count); nil != uerr {
		ac.TpReturn(atmi.TPFAIL, 0, data, 0)
		return
	}
	if uerr := data.BChg(atmi.TMSAGA_ACTION, 0, record.action); nil != uerr {
		ac.TpReturn(atmi.TPFAIL, 0, data, 0)
		return
	}
	if uerr := data.BChg(atmi.EX_DATA_STR, 1, string(record.payload)); nil != uerr {
		ac.TpReturn(atmi.TPFAIL, 0, data, 0)
		return
	}
	if uerr := data.BChg(atmi.TMTXTRYCNT, 0, record.sequence); nil != uerr {
		ac.TpReturn(atmi.TPFAIL, 0, data, 0)
		return
	}
	ac.TpReturn(atmi.TPSUCCESS, 0, data, 0)
}

// Initialize XA and advertise work, completion and inspection services.
// @param ac server context
// @return SUCCEED or FAIL
func Init(ac *atmi.ATMICtx) int {
	if err := ac.TpOpen(); nil != err {
		fmt.Println("TESTERROR:", err)
		return atmi.FAIL
	}
	if err := ac.TpAdvertise("GOREGISTER", "GOREGISTER", GOREGISTER); nil != err {
		return atmi.FAIL
	}
	if err := ac.TpAdvertise("GOCOMMIT", "GOACTION", GOACTION); nil != err {
		return atmi.FAIL
	}
	if err := ac.TpAdvertise("GOCOMMIT2", "GOACTION", GOACTION); nil != err {
		return atmi.FAIL
	}
	if err := ac.TpAdvertise("GOROLLBACK", "GOACTION", GOACTION); nil != err {
		return atmi.FAIL
	}
	if err := ac.TpAdvertise("GOSTATE", "GOSTATE", GOSTATE); nil != err {
		return atmi.FAIL
	}
	return atmi.SUCCEED
}

// Close the server's XA resource manager.
// @param ac server context
func Uninit(ac *atmi.ATMICtx) {
	if err := ac.TpClose(); nil != err {
		fmt.Println("TESTERROR:", err)
	}
}

// Run the native XATMI service dispatcher.
func main() {
	atmi.RuntimeInit()
	ac, err := atmi.NewATMICtx()
	if nil != err {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := ac.TpRun(Init, Uninit); nil != err {
		fmt.Fprintln(os.Stderr, "TESTERROR:", err)
		os.Exit(1)
	}
}
