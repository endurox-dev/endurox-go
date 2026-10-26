/**
 * @brief Buffer comparison, export and transport lengths
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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"

	atmi "github.com/endurox-dev/endurox-go/v9"
	"github.com/endurox-dev/endurox-go/v9/tests/testutil"
)

// Check comparison ordering and dynamic export ownership.
// @return Error if any comparison or round trip fails
func test_buffers() error {
	duration, durationErr := testutil.Duration()
	if nil != durationErr {
		return durationErr
	}
	var started time.Time
	var i int64
	ctx, err := atmi.NewATMICtx()
	if nil != err {
		return err
	}
	defer ctx.FreeATMICtx()
	left, err := ctx.NewCarray([]byte("a"))
	if nil != err {
		return err
	}
	defer ctx.TpFree(left.GetBuf())
	right, err := ctx.NewCarray([]byte("b"))
	if nil != err {
		return err
	}
	defer ctx.TpFree(right.GetBuf())
	roundtrip, err := ctx.NewCarray([]byte("initial"))
	if nil != err {
		return err
	}
	defer ctx.TpFree(roundtrip.GetBuf())
	payload := bytes.Repeat([]byte{0, 1, 255, 127}, 32768)
	large, err := ctx.NewCarray(payload)
	if nil != err {
		return err
	}
	defer ctx.TpFree(large.GetBuf())
	started = time.Now()
	for i = 0; i < 1 || time.Since(started) < duration; i++ {
		result, err := ctx.TpBufCmp(left, right, 0)
		if nil != err || result != -1 {
			return fmt.Errorf("iteration %d: negative ordering = %d, %v", i, result, err)
		}
		result, err = ctx.TpBufCmp(right, left, 0)
		if nil != err || result != 1 {
			return fmt.Errorf("positive ordering = %d, %v", result, err)
		}
		text, err := ctx.TpExport(large, 0)
		if nil != err {
			return err
		}
		if err := ctx.TpImport(text, roundtrip, 0); nil != err {
			return err
		}
		if !bytes.Equal(roundtrip.GetBytes(), payload) {
			return fmt.Errorf("large export/import changed bytes")
		}
	}
	fmt.Printf("test_buffers loop 1: %d iterations in %.3fs\n", i, time.Since(started).Seconds())
	started = time.Now()
	for i = 0; i < 1 || time.Since(started) < duration; i++ {
		if err := large.GetBuf().SetLen(4); nil != err {
			return err
		}
		if !bytes.Equal(large.GetBytes(), payload[:4]) {
			return fmt.Errorf("SetLen did not select the shorter view")
		}
		if err := large.GetBuf().SetLen(-1); nil == err {
			return fmt.Errorf("negative view length accepted")
		}
	}
	fmt.Printf("test_buffers loop 2: %d iterations in %.3fs\n", i, time.Since(started).Seconds())

	return nil
}

// Check synchronous and asynchronous service calls through buffer export.
// @return Error if a server reply or its typed-buffer replacement is invalid
func test_server() error {
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
		request, err := ac.NewCarray([]byte("request"))
		if nil != err {
			return err
		}
		if _, err := ac.TpCall("GOBUFFER", request, 0); nil != err {
			return err
		}
		reply, err := ac.CastToString(request.GetBuf())
		if nil != err {
			return err
		}
		var exported struct {
			Data string `json:"data"`
		}
		if err := json.Unmarshal([]byte(reply.GetString()), &exported); nil != err {
			return err
		}
		decoded, decodeErr := base64.StdEncoding.DecodeString(exported.Data)
		if nil != decodeErr || "request" != string(decoded) {
			return fmt.Errorf("synchronous export reply: %q, %v", decoded, decodeErr)
		}
		ac.TpFree(request.GetBuf())

		request, err = ac.NewCarray([]byte("request"))
		if nil != err {
			return err
		}
		cd, err := ac.TpACall("GOBUFFER", request, 0)
		if nil != err {
			return err
		}
		secondCd, err := ac.TpACall("GOBUFFER", request, 0)
		if nil != err {
			return err
		}
		requestedCd := secondCd

		ac.TpFree(request.GetBuf())
		reply, err = ac.NewString("")
		if nil != err {
			return err
		}
		if _, err := ac.TpGetRply(&secondCd, reply, 0); nil != err {
			return err
		}
		if requestedCd != secondCd {
			return fmt.Errorf("collected wrong call descriptor")
		}
		if _, err := ac.TpGetRply(&cd, reply, 0); nil != err {
			return err
		}

		if err := json.Unmarshal([]byte(reply.GetString()), &exported); nil != err {
			return err
		}
		decoded, decodeErr = base64.StdEncoding.DecodeString(exported.Data)
		if nil != decodeErr || "request" != string(decoded) {
			return fmt.Errorf("asynchronous export reply: %q, %v", decoded, decodeErr)
		}
		ac.TpFree(reply.GetBuf())
	}
	fmt.Printf("test_server loop 3: %d iterations in %.3fs\n", i, time.Since(started).Seconds())
	return nil
}

// Run the buffer comparison, export and transport lengths integration case.
func main() {
	if 1 != len(os.Args) {
		fmt.Fprintln(os.Stderr, "Usage: client")
		os.Exit(2)
	}
	atmi.RuntimeInit()
	fmt.Println("*** Buffer comparison, export and transport lengths (duration from NDRX_TEST_DURATATION; default 10s)")
	if err := test_buffers(); nil != err {
		fmt.Fprintf(os.Stderr, "TESTERROR: 13_buffer_utils: %v\n", err)
		os.Exit(1)
	}
	if err := test_server(); nil != err {
		fmt.Fprintln(os.Stderr, "TESTERROR: buffer server:", err)
		os.Exit(1)
	}
	fmt.Println("13_buffer_utils passed")
}

/* vim: set ts=4 sw=4 et smartindent: */
