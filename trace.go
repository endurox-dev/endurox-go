/**
 * @brief Trace context, application spans and context tags
 *
 * @file trace.go
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
#include <stdlib.h>
*/
import "C"
import (
	"runtime"
	"strings"
	"unsafe"
)

// TpTraceContextSet imports W3C traceparent; an empty string clears it.
// flags must be zero. The selected Enduro/X context owns the copied value.
// @param traceparent W3C traceparent string; an empty string clears the
//
//	propagated context
//
// @param flags Reserved; must be 0
// @return Nil on success, otherwise ATMIError
func (ac *ATMICtx) TpTraceContextSet(traceparent string, flags int64) ATMIError {
	defer runtime.KeepAlive(ac)
	if err := ac.checkContext(); nil != err {
		return err
	}
	if strings.IndexByte(traceparent, 0) >= 0 {
		return NewCustomATMIError(TPEINVAL, "traceparent contains NUL")
	}
	value := C.CString(traceparent)
	defer C.free(unsafe.Pointer(value))
	if C.Otptracecontextset(&ac.c_ctx, value, C.long(flags)) != SUCCEED {
		return ac.NewATMIError()
	}
	return nil
}

// TpTraceStateSet sets W3C tracestate on an active span; empty clears it.
// flags must be zero. The selected Enduro/X context owns the copied value.
// @param state W3C tracestate string; an empty string clears it
// @param flags Reserved; must be 0
// @return Nil on success, otherwise ATMIError
func (ac *ATMICtx) TpTraceStateSet(state string, flags int64) ATMIError {
	defer runtime.KeepAlive(ac)
	if err := ac.checkContext(); nil != err {
		return err
	}
	if strings.IndexByte(state, 0) >= 0 {
		return NewCustomATMIError(TPEINVAL, "state contains NUL")
	}
	value := C.CString(state)
	defer C.free(unsafe.Pointer(value))
	if C.Otptracestateset(&ac.c_ctx, value, C.long(flags)) != SUCCEED {
		return ac.NewATMIError()
	}
	return nil
}

// TpTraceContextGet returns the current W3C traceparent; flags must be zero.
// An empty string means no propagated value is present.
// @param flags Reserved; must be 0
// @return Current W3C traceparent and nil error, or empty string and
//
//	ATMIError; empty with nil error means no context
func (ac *ATMICtx) TpTraceContextGet(flags int64) (string, ATMIError) {
	defer runtime.KeepAlive(ac)
	if err := ac.checkContext(); nil != err {
		return "", err
	}
	var value [C.NDRX_TRACEPARENT_LEN + 1]C.char
	length := C.long(len(value))
	if C.Otptracecontextget(&ac.c_ctx, &value[0], &length, C.long(flags)) != SUCCEED {
		return "", ac.NewATMIError()
	}
	return C.GoString(&value[0]), nil
}

// TpTraceStateGet returns the current W3C tracestate; flags must be zero.
// An empty string means no propagated value is present.
// @param flags Reserved; must be 0
// @return Current W3C tracestate and nil error, or empty string and ATMIError;
//
//	empty with nil error means no state
func (ac *ATMICtx) TpTraceStateGet(flags int64) (string, ATMIError) {
	defer runtime.KeepAlive(ac)
	if err := ac.checkContext(); nil != err {
		return "", err
	}
	var value [C.NDRX_TRACE_STATE_MAX + 1]C.char
	length := C.long(len(value))
	if C.Otptracestateget(&ac.c_ctx, &value[0], &length, C.long(flags)) != SUCCEED {
		return "", ac.NewATMIError()
	}
	return C.GoString(&value[0]), nil
}

// TpTraceSpan identifies a span started by TpTraceSpanBegin. It is valid only
// in the ATMICtx that started it, until TpTraceSpanEnd.
type TpTraceSpan uint64

// TpTraceSpanBegin starts a span with optional UBF attributes and returns its
// handle. TPTRACE_F_ACTIVE spans must be ended in reverse nesting order.
// @param name NUL-free span name
// @param attributes Optional UBF attributes, or nil
// @param flags TPTRACE_F_ACTIVE to activate the span, TPTRACE_F_NEW_ROOT for a
//
//	new trace, or 0
//
// @return Span handle and nil error, or 0 and ATMIError
func (ac *ATMICtx) TpTraceSpanBegin(name string, attributes *TypedUBF,
	flags int64) (TpTraceSpan, ATMIError) {
	defer runtime.KeepAlive(ac)
	defer runtime.KeepAlive(attributes)
	if err := ac.checkContext(); nil != err {
		return 0, err
	}
	if strings.IndexByte(name, 0) >= 0 {
		return 0, NewCustomATMIError(TPEINVAL, "span name contains NUL")
	}
	var attrs *C.UBFH
	if nil != attributes {
		if nil == attributes.Buf {
			return 0, NewCustomATMIError(TPEINVAL, "attributes buffer is nil")
		}
		attrs = (*C.UBFH)(unsafe.Pointer(attributes.Buf.C_ptr))
	}
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	var handle C.uintptr_t
	if C.Otptracespanbegin(&ac.c_ctx, cName, attrs, C.long(flags), &handle) != SUCCEED {
		return 0, ac.NewATMIError()
	}
	return TpTraceSpan(handle), nil
}

// TpTraceSpanEnd ends span with a TPTRACE_STATUS_* result and optional UBF
// attributes. The native engine releases the span; the handle must not be reused.
// @param span Live handle returned by TpTraceSpanBegin; cannot be reused
//
//	after completion
//
// @param status TPTRACE_STATUS_UNSET, TPTRACE_STATUS_OK or
//
//	TPTRACE_STATUS_ERROR
//
// @param attributes Optional UBF completion attributes, or nil
// @param flags Reserved; must be 0
// @return Nil on success, otherwise ATMIError
func (ac *ATMICtx) TpTraceSpanEnd(span TpTraceSpan, status int,
	attributes *TypedUBF, flags int64) ATMIError {
	defer runtime.KeepAlive(ac)
	defer runtime.KeepAlive(attributes)
	if err := ac.checkContext(); nil != err {
		return err
	}
	var attrs *C.UBFH
	if nil != attributes {
		if nil == attributes.Buf {
			return NewCustomATMIError(TPEINVAL, "attributes buffer is nil")
		}
		attrs = (*C.UBFH)(unsafe.Pointer(attributes.Buf.C_ptr))
	}
	if C.Otptracespanend(&ac.c_ctx, C.uintptr_t(span), C.int(status), attrs,
		C.long(flags)) != SUCCEED {
		return ac.NewATMIError()
	}
	return nil
}

// TpTraceEvent records name and optional UBF attributes on the current span.
// @param name NUL-free event name
// @param attributes Optional UBF attributes, or nil
// @param flags Reserved; must be 0
// @return Nil on success, otherwise ATMIError
func (ac *ATMICtx) TpTraceEvent(name string, attributes *TypedUBF,
	flags int64) ATMIError {
	defer runtime.KeepAlive(ac)
	defer runtime.KeepAlive(attributes)
	if err := ac.checkContext(); nil != err {
		return err
	}
	if strings.IndexByte(name, 0) >= 0 {
		return NewCustomATMIError(TPEINVAL, "event name contains NUL")
	}
	var attrs *C.UBFH
	if nil != attributes {
		if nil == attributes.Buf {
			return NewCustomATMIError(TPEINVAL, "attributes buffer is nil")
		}
		attrs = (*C.UBFH)(unsafe.Pointer(attributes.Buf.C_ptr))
	}
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	if C.Otptraceevent(&ac.c_ctx, cName, attrs, C.long(flags)) != SUCCEED {
		return ac.NewATMIError()
	}
	return nil
}

// TpContextTagSet sets a tag in the TPCTXTAG_* scope/destinations in flags.
// NEXT tags persist until explicitly cleared; key and value are copied by C.
// @param key NUL-free tag name
// @param value NUL-free tag value copied by the native API
// @param flags TPCTXTAG_PROCESS, TPCTXTAG_THREAD or TPCTXTAG_NEXT scope with
//
//	TPCTXTAG_LOG and/or TPCTXTAG_TRACE
//
// @return Nil on success, otherwise ATMIError
func (ac *ATMICtx) TpContextTagSet(key, value string, flags int64) ATMIError {
	defer runtime.KeepAlive(ac)
	if err := ac.checkContext(); nil != err {
		return err
	}
	if strings.IndexByte(key, 0) >= 0 || strings.IndexByte(value, 0) >= 0 {
		return NewCustomATMIError(TPEINVAL, "tag key or value contains NUL")
	}
	cKey, cValue := C.CString(key), C.CString(value)
	defer C.free(unsafe.Pointer(cKey))
	defer C.free(unsafe.Pointer(cValue))
	if C.Otpcontexttagset(&ac.c_ctx, cKey, cValue, C.long(flags)) != SUCCEED {
		return ac.NewATMIError()
	}
	return nil
}

// TpContextTagClear removes key from the scopes/destinations selected by flags.
// @param key Tag name to remove
// @param flags TPCTXTAG scope and signal selection
// @return Nil on success, otherwise ATMIError
func (ac *ATMICtx) TpContextTagClear(key string, flags int64) ATMIError {
	defer runtime.KeepAlive(ac)
	if err := ac.checkContext(); nil != err {
		return err
	}
	if strings.IndexByte(key, 0) >= 0 {
		return NewCustomATMIError(TPEINVAL, "tag key contains NUL")
	}
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	if C.Otpcontexttagclear(&ac.c_ctx, cKey, C.long(flags)) != SUCCEED {
		return ac.NewATMIError()
	}
	return nil
}

// TpContextTagClearAll clears the scopes/destinations selected by flags.
// @param flags TPCTXTAG scope and signal selection to clear
// @return Nil on success, otherwise ATMIError
func (ac *ATMICtx) TpContextTagClearAll(flags int64) ATMIError {
	defer runtime.KeepAlive(ac)
	if err := ac.checkContext(); nil != err {
		return err
	}
	if C.Otpcontexttagclearall(&ac.c_ctx, C.long(flags)) != SUCCEED {
		return ac.NewATMIError()
	}
	return nil
}

/* vim: set ts=4 sw=4 et smartindent: */
