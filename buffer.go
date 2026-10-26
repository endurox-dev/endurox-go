/**
 * @brief Buffer comparison and reply-queue integration
 *
 * @file buffer.go
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
*/
import "C"
import (
	"runtime"
)

// TpBufCmp compares typed values. -1, 0 and 1 are successful order results;
// native -2 is returned with an ATMI error. Nil denotes a NULL typed buffer.
// @param left Left typed buffer, or nil for a native NULL buffer
// @param right Right typed buffer, or nil for a native NULL buffer
// @param flags Native comparison flags; use 0
// @return -1, 0 or 1 for ordering and nil error; -2 and ATMIError on failure
func (ac *ATMICtx) TpBufCmp(left, right TypedBuffer, flags int64) (int, ATMIError) {
	defer runtime.KeepAlive(ac)
	defer runtime.KeepAlive(left)
	defer runtime.KeepAlive(right)
	if err := ac.checkContext(); nil != err {
		return -2, err
	}
	var lp, rp *C.char
	var ll, rl C.long
	if b := nativeBuffer(left); nil != b {
		lp, ll = b.C_ptr, b.C_len
	}
	if b := nativeBuffer(right); nil != b {
		rp, rl = b.C_ptr, b.C_len
	}
	result := C.Otpbufcmp(&ac.c_ctx, lp, ll, rp, rl, C.long(flags))
	if result == -2 {
		return int(result), ac.NewATMIError()
	}
	return int(result), nil
}

// TpExtGetReplyQFD returns the context's borrowed, pollable reply-queue FD.
// Do not close it. Unsupported polling backends return the native ATMI error.
// @return Borrowed pollable reply-queue descriptor and nil error, or -1 and
//
//	ATMIError; do not close the descriptor
func (ac *ATMICtx) TpExtGetReplyQFD() (int, ATMIError) {
	defer runtime.KeepAlive(ac)
	if err := ac.checkContext(); nil != err {
		return FAIL, err
	}
	fd := C.Otpext_getreplyqfd(&ac.c_ctx)
	if fd == FAIL {
		return FAIL, ac.NewATMIError()
	}
	return int(fd), nil
}

// Len returns the buffer's current transport length.
// @return Current transport length in bytes; 0 for a nil buffer
func (buf *ATMIBuf) Len() int64 {
	if nil == buf {
		return 0
	}
	return int64(buf.C_len)
}

// SetLen selects the data length within the existing allocation.
// It does not resize the buffer; zero is valid for self-describing buffers.
// @param length Nonnegative transport length within the existing allocation; 0
//
//	is valid for self-describing buffers
//
// @return Nil on success, otherwise ATMIError; the allocation is not resized
func (buf *ATMIBuf) SetLen(length int64) ATMIError {
	defer runtime.KeepAlive(buf)
	if nil == buf || length < 0 || int64(C.long(length)) != length {
		return NewCustomATMIError(TPEINVAL, "invalid buffer or length")
	}
	if nil == buf.C_ptr {
		if length != 0 {
			return NewCustomATMIError(TPEINVAL, "NULL buffer has nonzero length")
		}
		buf.C_len = 0
		return nil
	}
	if err := buf.Ctx.checkContext(); nil != err {
		return err
	}
	capacity, err := buf.Ctx.TpTypes(buf, nil, nil)
	if nil != err {
		return err
	}
	if length > capacity {
		return NewCustomATMIError(TPEINVAL, "length exceeds buffer allocation")
	}
	buf.C_len = C.long(length)
	return nil
}

/* vim: set ts=4 sw=4 et smartindent: */
