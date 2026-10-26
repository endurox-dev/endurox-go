/**
 * @brief Distributed lock manager bindings
 *
 * @file dlm.go
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

#include <oatmi.h>
#include <ubf.h>
#include <stdlib.h>
*/
import "C"
import (
	"runtime"
	"strings"
	"unsafe"
)

// TpDlmCtl holds one native DLM request control.
// The zero value is ready to use. Copies of an initialized value share the native record; use a new
// value for independent state. Access a record serially.
type TpDlmCtl struct {
	c_ctl *C.TPDLMCTL
}

// native returns the single backing C record, allocating it on first use.
// @return Native record shared by this wrapper and its copies
func (ctl *TpDlmCtl) native() *C.TPDLMCTL {
	if nil == ctl.c_ctl {
		ctl.c_ctl = new(C.TPDLMCTL)
	}
	return ctl.c_ctl
}

// GetFlags reads DLM call flags from the native record.
// @return Current DLM call flags
func (ctl *TpDlmCtl) GetFlags() int64 {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return 0
	}
	return int64(ctl.native().flags)
}

// SetFlags writes DLM call flags to the native record.
// @param value New DLM call flags
// @return Nil on success, otherwise TPEINVAL without changing the field
func (ctl *TpDlmCtl) SetFlags(value int64) ATMIError {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return NewCustomATMIError(TPEINVAL, "TpDlmCtl is nil")
	}
	if int64(C.long(value)) != value {
		return NewCustomATMIError(TPEINVAL, "Flags exceeds native long")
	}
	ctl.native().flags = C.long(value)
	return nil
}

// GetDLMSpace reads DLM space name from the native record.
// @return Current DLM space name
func (ctl *TpDlmCtl) GetDLMSpace() string {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return ""
	}
	return C.GoString(&ctl.native().dlmspace[0])
}

// SetDLMSpace writes DLM space name to the native record.
// @param value New DLM space name
// @return Nil on success, otherwise TPEINVAL without changing the field
func (ctl *TpDlmCtl) SetDLMSpace(value string) ATMIError {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return NewCustomATMIError(TPEINVAL, "TpDlmCtl is nil")
	}
	return copyCText(ctl.native().dlmspace[:], value, "DLMSpace")
}

// GetSvcName reads resolved DLM service name from the native record.
// @return Current resolved DLM service name
func (ctl *TpDlmCtl) GetSvcName() string {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return ""
	}
	return C.GoString(&ctl.native().svcname[0])
}

// SetSvcName writes resolved DLM service name to the native record.
// @param value New resolved DLM service name
// @return Nil on success, otherwise TPEINVAL without changing the field
func (ctl *TpDlmCtl) SetSvcName(value string) ATMIError {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return NewCustomATMIError(TPEINVAL, "TpDlmCtl is nil")
	}
	return copyCText(ctl.native().svcname[:], value, "SvcName")
}

// GetWaitTime reads lock wait time in milliseconds from the native record.
// @return Current lock wait time in milliseconds
func (ctl *TpDlmCtl) GetWaitTime() int64 {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return 0
	}
	return int64(ctl.native().wait_time)
}

// SetWaitTime writes lock wait time in milliseconds to the native record.
// @param value New lock wait time in milliseconds
// @return Nil on success, otherwise TPEINVAL without changing the field
func (ctl *TpDlmCtl) SetWaitTime(value int64) ATMIError {
	defer runtime.KeepAlive(ctl)
	if nil == ctl {
		return NewCustomATMIError(TPEINVAL, "TpDlmCtl is nil")
	}
	if int64(C.long(value)) != value {
		return NewCustomATMIError(TPEINVAL, "WaitTime exceeds native long")
	}
	ctl.native().wait_time = C.long(value)
	return nil
}

// TpDLMMkCall prepares ub and ctl for dispatch. It may replace ub's allocation.
// @param ub Request UBF; its allocation and length are updated if native
//
//	preparation replaces it
//
// @param ctl DLM control configured with SetDLMSpace and SetWaitTime (milliseconds);
//
//	GetSvcName returns the selected service
//
// @return Nil on success, otherwise ATMIError
func (ac *ATMICtx) TpDLMMkCall(ub *TypedUBF, ctl *TpDlmCtl) ATMIError {
	defer runtime.KeepAlive(ac)
	defer runtime.KeepAlive(ub)
	defer runtime.KeepAlive(ctl)
	if err := ac.checkContext(); nil != err {
		return err
	}
	if nil == ub || nil == ub.Buf || nil == ctl {
		return NewCustomATMIError(TPEINVAL, "DLM buffer and control are required")
	}
	ptr := (*C.UBFH)(unsafe.Pointer(ub.Buf.C_ptr))
	ret := C.Otpdlmmkcall(&ac.c_ctx, &ptr, ctl.native())
	var err ATMIError
	if ret != SUCCEED {
		err = ac.NewATMIError()
	}
	ub.Buf.C_ptr = (*C.char)(unsafe.Pointer(ptr))
	if nil != ptr {
		ub.Buf.C_len = C.long(C.Bsizeof(ptr))
	}
	return err
}

// TpDLMCall dispatches a DLM request with retries and updates tb with the reply.
// Inspect per-command EX_DLM_RVAL1 even when the ATMI call succeeds.
// @param svc DLM service name returned in TpDlmCtl.GetSvcName()
// @param tb Request buffer, updated with the reply pointer and length even on
//
//	failure
//
// @param flags XATMI call flags passed to the native DLM retry API
// @return 0 and nil error on successful transport, otherwise -1 and ATMIError;
//
//	inspect EX_DLM_RVAL1 for individual command results
func (ac *ATMICtx) TpDLMCall(svc string, tb TypedBuffer, flags int64) (int, ATMIError) {
	defer runtime.KeepAlive(ac)
	defer runtime.KeepAlive(tb)
	if err := ac.checkContext(); nil != err {
		return FAIL, err
	}
	buf := nativeBuffer(tb)
	if nil == buf || strings.IndexByte(svc, 0) >= 0 {
		return FAIL, NewCustomATMIError(TPEINVAL, "DLM buffer is nil or service contains NUL")
	}
	cSvc := C.CString(svc)
	defer C.free(unsafe.Pointer(cSvc))
	ret := C.Otpdlmcall(&ac.c_ctx, cSvc, buf.C_ptr, buf.C_len,
		&buf.C_ptr, &buf.C_len, C.long(flags))
	var err ATMIError
	if ret != SUCCEED {
		err = ac.NewATMIError()
	}
	if view, ok := tb.(*TypedVIEW); ok && nil != buf.C_ptr {
		var kind string
		_, typeErr := ac.TpTypes(buf, &kind, &view.view)
		if nil == err {
			err = typeErr
		}
	}
	return int(ret), err
}

// TpDLMACall dispatches an asynchronous DLM request with native retry handling.
// svc is resolved by TpDLMMkCall; the result is the call/conversation descriptor.
// @param svc DLM service name returned in TpDlmCtl.GetSvcName()
// @param tb Request typed buffer, retained by the caller
// @param flags XATMI asynchronous call flags
// @return Asynchronous call descriptor and nil error, or -1 and ATMIError
func (ac *ATMICtx) TpDLMACall(svc string, tb TypedBuffer, flags int64) (int, ATMIError) {
	defer runtime.KeepAlive(ac)
	defer runtime.KeepAlive(tb)
	if err := ac.checkContext(); nil != err {
		return FAIL, err
	}
	buf := nativeBuffer(tb)
	if nil == buf || strings.IndexByte(svc, 0) >= 0 {
		return FAIL, NewCustomATMIError(TPEINVAL, "DLM buffer is nil or service contains NUL")
	}
	cSvc := C.CString(svc)
	defer C.free(unsafe.Pointer(cSvc))
	ret := C.Otpdlmacall(&ac.c_ctx, cSvc, buf.C_ptr, buf.C_len, C.long(flags))
	if ret == FAIL {
		return int(ret), ac.NewATMIError()
	}
	return int(ret), nil
}

// TpDLMConnect opens a DLM conversation with native retry handling.
// svc is resolved by TpDLMMkCall; the result is the call/conversation descriptor.
// @param svc DLM service name returned in TpDlmCtl.GetSvcName()
// @param tb Initial conversation buffer, retained by the caller
// @param flags XATMI conversation flags
// @return Conversation descriptor and nil error, or -1 and ATMIError
func (ac *ATMICtx) TpDLMConnect(svc string, tb TypedBuffer, flags int64) (int, ATMIError) {
	defer runtime.KeepAlive(ac)
	defer runtime.KeepAlive(tb)
	if err := ac.checkContext(); nil != err {
		return FAIL, err
	}
	buf := nativeBuffer(tb)
	if nil == buf || strings.IndexByte(svc, 0) >= 0 {
		return FAIL, NewCustomATMIError(TPEINVAL, "DLM buffer is nil or service contains NUL")
	}
	cSvc := C.CString(svc)
	defer C.free(unsafe.Pointer(cSvc))
	ret := C.Otpdlmconnect(&ac.c_ctx, cSvc, buf.C_ptr, buf.C_len, C.long(flags))
	if ret == FAIL {
		return int(ret), ac.NewATMIError()
	}
	return int(ret), nil
}

// TpDLMToutGet returns effective and per-attempt milliseconds, and the
// per-attempt ATMI timeout in seconds. Nonpositive waitTime uses ATMI settings.
// @param waitTime Requested total wait in milliseconds; nonpositive values use
//
//	ATMI timeout settings
//
// @return Effective milliseconds, per-attempt milliseconds, per-attempt ATMI
//
//	seconds, and ATMIError if calculation fails
func (ac *ATMICtx) TpDLMToutGet(waitTime int64) (int64, int64, int64, ATMIError) {
	defer runtime.KeepAlive(ac)
	if err := ac.checkContext(); nil != err {
		return 0, 0, 0, err
	}
	var effective, attempt, block C.long
	if C.Otpdlmtoutget(&ac.c_ctx, C.long(waitTime), &effective, &attempt,
		&block) != SUCCEED {
		return 0, 0, 0, ac.NewATMIError()
	}
	return int64(effective), int64(attempt), int64(block), nil
}

// TpDLMAttemptsGet returns the configured attempt limit, not a call counter.
// @return Configured retry attempt limit and nil error, or 0 and ATMIError;
//
//	this is not the count consumed by a call
func (ac *ATMICtx) TpDLMAttemptsGet() (int, ATMIError) {
	defer runtime.KeepAlive(ac)
	if err := ac.checkContext(); nil != err {
		return 0, err
	}
	count := C.Otpdlmattemptsget(&ac.c_ctx)
	if count == FAIL {
		return 0, ac.NewATMIError()
	}
	return int(count), nil
}

/* vim: set ts=4 sw=4 et smartindent: */
