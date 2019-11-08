/**
 * @brief Plain text IPC buffer support
 *
 * @file typed_string.go
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
type TypedString struct {
	Buf *ATMIBuf
}

//Return The ATMI buffer to caller
func (u *TypedString) GetBuf() *ATMIBuf {
	return u.Buf
}

//Allocate new string buffer
//@param s - source string
func (ac *ATMICtx) NewString(gs string) (*TypedString, ATMIError) {
	var buf TypedString

	c_val := C.CString(gs)
	defer C.free(unsafe.Pointer(c_val))

	size := int64(C.strlen(c_val) + 1) /* 1 for EOS. */

	if ptr, err := ac.TpAlloc("STRING", "", size); nil != err {
		return nil, err
	} else {
		buf.Buf = ptr
		C.strcpy(buf.Buf.C_ptr, c_val)
		buf.Buf.TpSetCtxt(ac)

		return &buf, nil
	}
}

//Get the String Handler from ATMI Buffer
func (ac *ATMICtx) CastToString(abuf *ATMIBuf) (*TypedString, ATMIError) {
	var buf TypedString

	buf.Buf = abuf

	return &buf, nil
}

//Get the string value out from buffer
//@return String value
func (s *TypedString) GetString() string {
	ret := C.GoString(s.Buf.C_ptr)

	s.Buf.nop()
	return ret

}

//Set the string to the buffer
//@param str 	String value
func (s *TypedString) SetString(gs string) ATMIError {

	c_val := C.CString(gs)
	defer C.free(unsafe.Pointer(c_val))

	new_size := int64(C.strlen(c_val) + 1) /* 1 for EOS. */

	if cur_size, err := s.Buf.Ctx.TpTypes(s.Buf, nil, nil); nil != err {
		return err
	} else {
		if cur_size >= new_size {
			C.strcpy(s.Buf.C_ptr, c_val)
		} else if err := s.Buf.TpRealloc(new_size); nil != err {
			return err
		} else {
			C.strcpy(s.Buf.C_ptr, c_val)
		}
	}

	s.Buf.nop()

	return nil
}

///////////////////////////////////////////////////////////////////////////////////
// Wrappers for memory management
///////////////////////////////////////////////////////////////////////////////////

func (u *TypedString) TpRealloc(size int64) ATMIError {
	return u.Buf.TpRealloc(size)
}

/* vim: set ts=4 sw=4 et smartindent: */
