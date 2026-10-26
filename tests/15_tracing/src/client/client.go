/**
 * @brief Trace context, spans and context tags
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
	"strings"
	"time"

	atmi "github.com/endurox-dev/endurox-go/v9"
	"github.com/endurox-dev/endurox-go/v9/tests/testutil"
)

// Check trace propagation, context isolation, spans and context tags.
// @return Error if a trace value, span lifecycle or tag operation is invalid
func test_tracing() error {
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
	other, err := atmi.NewATMICtx()
	if nil != err {
		return err
	}
	defer other.FreeATMICtx()

	parent := "00-0123456789abcdef0123456789abcdef-0123456789abcdef-01"
	state := "vendor=one,1tenant@system=two"
	started = time.Now()
	for i = 0; i < 1 || time.Since(started) < duration; i++ {
		if err := ac.TpTraceContextSet("bad\x00parent", 0); nil == err || atmi.TPEINVAL != err.Code() {
			return fmt.Errorf("embedded NUL traceparent: %v", err)
		}
		if err := ac.TpTraceContextSet(parent, 0); nil != err {
			return err
		}
		got, err := ac.TpTraceContextGet(0)
		if nil != err || parent != got {
			return fmt.Errorf("traceparent round trip: %q, %v", got, err)
		}
		root, err := ac.TpTraceSpanBegin("go-root", nil, atmi.TPTRACE_F_ACTIVE)
		if nil != err {
			return err
		}
		activeParent, err := ac.TpTraceContextGet(0)
		if nil != err {
			return err
		}
		if err := ac.TpTraceStateSet(state, 0); nil != err {
			return err
		}
		got, err = ac.TpTraceStateGet(0)
		if nil != err || state != got {
			return fmt.Errorf("tracestate round trip: %q, %v", got, err)
		}
		request, err := ac.NewString("request")
		if nil != err {
			return err
		}
		if _, err := ac.TpCall("GOTRACE", request, 0); nil != err {
			return err
		}
		propagated := strings.Split(request.GetString(), "\n")
		if 2 != len(propagated) || len(parent) != len(propagated[0]) ||
			parent[:36] != propagated[0][:36] || state != propagated[1] {
			return fmt.Errorf("server trace propagation: %q", propagated)
		}
		ac.TpFree(request.GetBuf())
		got, err = other.TpTraceContextGet(0)
		if nil != err || "" != got {
			return fmt.Errorf("traceparent leaked into another context: %q, %v",
				got, err)
		}
		got, err = other.TpTraceStateGet(0)
		if nil != err || "" != got {
			return fmt.Errorf("tracestate leaked into another context: %q, %v",
				got, err)
		}

		handle, err := ac.TpTraceSpanBegin("go-tracing", nil,
			atmi.TPTRACE_F_ACTIVE)
		if nil != err {
			return err
		}
		if err := ac.TpTraceEvent("go-event", nil, 0); nil != err {
			return err
		}
		got, err = ac.TpTraceContextGet(0)
		if nil != err || len(parent) != len(got) {
			return fmt.Errorf("active traceparent: %q, %v", got, err)
		}
		if activeParent[:36] != got[:36] || activeParent[36:52] == got[36:52] {
			return fmt.Errorf("span did not retain trace ID and replace span ID")
		}
		if err := ac.TpTraceSpanEnd(handle, atmi.TPTRACE_STATUS_OK, nil, 0); nil != err {
			return err
		}
		got, err = ac.TpTraceContextGet(0)
		if nil != err || activeParent != got {
			return fmt.Errorf("parent was not restored: %q, %v", got, err)
		}

		if err := ac.TpContextTagSet("case", "tracing",
			atmi.TPCTXTAG_THREAD|atmi.TPCTXTAG_TRACE); nil != err {
			return err
		}
		if err := ac.TpContextTagSet("request", "one",
			atmi.TPCTXTAG_NEXT|atmi.TPCTXTAG_LOG); nil != err {
			return err
		}
		if err := ac.TpContextTagClear("case", atmi.TPCTXTAG_THREAD); nil != err {
			return err
		}
		if err := ac.TpContextTagClearAll(atmi.TPCTXTAG_NEXT); nil != err {
			return err
		}
		if err := ac.TpTraceStateSet("", 0); nil != err {
			return err
		}
		if err := ac.TpTraceSpanEnd(root, atmi.TPTRACE_STATUS_OK, nil, 0); nil != err {
			return err
		}
		if err := ac.TpTraceContextSet("", 0); nil != err {
			return err
		}
		got, err = ac.TpTraceStateGet(0)
		if nil != err || "" != got {
			return fmt.Errorf("tracestate was not cleared: %q, %v", got, err)
		}
		got, err = ac.TpTraceContextGet(0)
		if nil != err || "" != got {
			return fmt.Errorf("traceparent was not cleared: %q, %v", got, err)
		}
	}
	fmt.Printf("test_tracing loop 1: %d iterations in %.3fs\n", i, time.Since(started).Seconds())
	return nil
}

// Run the trace context, spans and context tags integration case.
func main() {
	if 1 != len(os.Args) {
		fmt.Fprintln(os.Stderr, "Usage: client")
		os.Exit(2)
	}
	atmi.RuntimeInit()
	fmt.Println("*** Trace context, spans and context tags (duration from NDRX_TEST_DURATATION; default 10s)")
	if err := test_tracing(); nil != err {
		fmt.Fprintf(os.Stderr, "TESTERROR: 15_tracing: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("15_tracing passed")
}

/* vim: set ts=4 sw=4 et smartindent: */
