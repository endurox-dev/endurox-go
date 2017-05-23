package atmi
/*
** Typed C-Array (binary array) IPC buffer support
**
** @file typed_carray.go
**
** -----------------------------------------------------------------------------
** Enduro/X Middleware Platform for Distributed Transaction Processing
** Copyright (C) 2015, Mavimax, Ltd. All Rights Reserved.
** This software is released under one of the following licenses:
** GPL or Mavimax's license for commercial use.
** -----------------------------------------------------------------------------
** GPL license:
**
** This program is free software; you can redistribute it and/or modify it under
** the terms of the GNU General Public License as published by the Free Software
** Foundation; either version 2 of the License, or (at your option) any later
** version.
**
** This program is distributed in the hope that it will be useful, but WITHOUT ANY
** WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
** PARTICULAR PURPOSE. See the GNU General Public License for more details.
**
** You should have received a copy of the GNU General Public License along with
** this program; if not, write to the Free Software Foundation, Inc., 59 Temple
** Place, Suite 330, Boston, MA 02111-1307 USA
**
** -----------------------------------------------------------------------------
** A commercial use license is available from Mavimax, Ltd
** contact@mavimax.com
** -----------------------------------------------------------------------------
*/

/*
#cgo pkg-config: atmisrvinteg

#include <xatmi.h>
#include <string.h>
#include <stdlib.h>
#include <ubf.h>

void * c_get_void_ptr(char * ptr)
{
	return (void *)ptr;
}

*/
import "C"
import "unsafe"

//UBF Buffer
type TypedCarray struct {
	Buf *ATMIBuf
}

//Return The ATMI buffer to caller
func (u *TypedCarray) GetBuf() *ATMIBuf {
	return u.Buf
}

func cpyGo2C(c *C.char, b []byte) {
	for i := 0; i < len(b); i++ {
		*(*C.char)(unsafe.Pointer(uintptr(C.c_get_void_ptr(c)) + uintptr(i))) = C.char(b[i])
	}
}

//Allocate new string buffer
//@param s - source string
func (ac *ATMICtx) NewCarray(b []byte) (*TypedCarray, ATMIError) {
	var buf TypedCarray

	if ptr, err := ac.TpAlloc("CARRAY", "", int64(len(b))); nil != err {
		return nil, err
	} else {
		buf.Buf = ptr

		/* Copy off the bytes to C buf */
		cpyGo2C(buf.Buf.C_ptr, b)
		buf.Buf.C_len = C.long(len(b))
		buf.Buf.TpSetCtxt(ac)

		return &buf, nil
	}
}

//Get the String Handler
func (ac *ATMICtx) CastToCarray(abuf *ATMIBuf) (*TypedCarray, ATMIError) {
	var buf TypedCarray
	buf.Buf = abuf
	return &buf, nil
}

//Get the string value out from buffer
//@return String value
func (s *TypedCarray) GetBytes() []byte {
	b := make([]byte, s.Buf.C_len)

	for i := 0; i < len(b); i++ {
		b[i] = byte(*(*C.char)(unsafe.Pointer(uintptr(C.c_get_void_ptr(s.Buf.C_ptr)) + uintptr(i))))
	}
	return b

}

//@param str 	String value
func (s *TypedCarray) SetBytes(b []byte) ATMIError {

	new_size := int64(len(b))

	if cur_size, err := s.Buf.Ctx.TpTypes(s.Buf, nil, nil); nil != err {
		return err
	} else {
		if cur_size >= new_size {
			cpyGo2C(s.Buf.C_ptr, b)
			s.Buf.C_len = C.long(new_size)
		} else if err := s.Buf.TpRealloc(new_size); nil != err {
			return err
		} else {
			cpyGo2C(s.Buf.C_ptr, b)
			s.Buf.C_len = C.long(new_size)
		}
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////
// Wrappers for memory management
///////////////////////////////////////////////////////////////////////////////////

func (u *TypedCarray) TpRealloc(size int64) ATMIError {
	return u.Buf.TpRealloc(size)
}
