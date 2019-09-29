/**
 * @brief JSON IPC Buffer support
 *
 * @file typed_json.go
 */
/* -----------------------------------------------------------------------------
 * Enduro/X Middleware Platform for Distributed Transaction Processing
 * Copyright (C) 2009-2016, ATR Baltic, Ltd. All Rights Reserved.
 * Copyright (C) 2017-2019, Mavimax, Ltd. All Rights Reserved.
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

#include <xatmi.h>
#include <string.h>
#include <stdlib.h>
#include <ubf.h>

*/
import "C"
import "unsafe"

//UBF Buffer
type TypedJSON struct {
	Buf *ATMIBuf
}

//Return The ATMI buffer to caller
func (u *TypedJSON) GetBuf() *ATMIBuf {
	return u.Buf
}

//Allocate new JSON buffer
//@param s - source string
func (ac *ATMICtx) NewJSON(b []byte) (*TypedJSON, ATMIError) {
	var buf TypedJSON

	c_val := C.CString(string(b))
	defer C.free(unsafe.Pointer(c_val))

	size := int64(C.strlen(c_val) + 1) /* 1 for EOS. */

	if ptr, err := ac.TpAlloc("JSON", "", size); nil != err {
		return nil, err
	} else {
		buf.Buf = ptr
		C.strcpy(buf.Buf.C_ptr, c_val)

		buf.Buf.TpSetCtxt(ac)

		return &buf, nil
	}
}

//Get the JSON Handler from ATMI Buffer
func (ac *ATMICtx) CastToJSON(abuf *ATMIBuf) (*TypedJSON, ATMIError) {
	var buf TypedJSON

	buf.Buf = abuf

	return &buf, nil
}

//Get the string value out from buffer
//@return JSON value
func (j *TypedJSON) GetJSONText() string {
	return C.GoString(j.Buf.C_ptr)
}

//Get JSON bytes..
func (j *TypedJSON) GetJSON() []byte {
	return []byte(C.GoString(j.Buf.C_ptr))
}

//Set JSON bytes
func (j *TypedJSON) SetJSON(b []byte) ATMIError {
	return j.SetJSONText(string(b))
}

//Set the string to the buffer
//@param str 	JSON value
func (j *TypedJSON) SetJSONText(gs string) ATMIError {
	c_val := C.CString(gs)
	defer C.free(unsafe.Pointer(c_val))

	new_size := int64(C.strlen(c_val) + 1) /* 1 for EOS. */

	if cur_size, err := j.Buf.Ctx.TpTypes(j.Buf, nil, nil); nil != err {
		return err
	} else {
		if cur_size >= new_size {
			C.strcpy(j.Buf.C_ptr, c_val)
		} else if err := j.Buf.TpRealloc(new_size); nil != err {
			return err
		} else {
			C.strcpy(j.Buf.C_ptr, c_val)
		}
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////
// Wrappers for memory management
///////////////////////////////////////////////////////////////////////////////////

func (u *TypedJSON) TpRealloc(size int64) ATMIError {
	return u.Buf.TpRealloc(size)
}
/* vim: set ts=4 sw=4 et smartindent: */
