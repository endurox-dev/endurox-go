package atmi

/*
** VIEW buffer support - dynamic access
**
** @file typed_view.go
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

 #include <string.h>
 #include <stdlib.h>
 #include <ndebug.h>
 #include <oatmi.h>
 #include <ubf.h>
 #include <oubf.h>
 #include <oatmi.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

//View flags
const (
	BVACCESS_NOTNULL = 0x00000001 //View access mode (return non null values only)
	VIEW_NAME_LEN    = 33         //View name max len
	VIEW_CNAME_LEN   = 256        // View filed max len
)

///////////////////////////////////////////////////////////////////////////////////
// Buffer def, typedefs
///////////////////////////////////////////////////////////////////////////////////

//UBF Buffer
type TypedVIEW struct {
	view string //Cached view name
	Buf  *ATMIBuf
}

//State for view iterator
type BVNextState struct {
	state C.Bvnext_state_t
}

//Return The ATMI buffer to caller
func (u *TypedVIEW) GetBuf() *ATMIBuf {
	return u.Buf
}

///////////////////////////////////////////////////////////////////////////////////
// VIEW API
///////////////////////////////////////////////////////////////////////////////////

//Allocate the UBF buffer
//@param size	Buffer size in bytes
//@return UBF Handler, ATMI Error
func (ac *ATMICtx) VIEWAlloc(view string, size int64) (TypedVIEW, ATMIError) {
	var err ATMIError
	var buf TypedVIEW
	buf.Buf, err = ac.TpAlloc("VIEW", view, size)
	buf.view = view
	return buf, err
}

//Get the UBF Handler
func (ac *ATMICtx) CastToVIEW(abuf *ATMIBuf) (*TypedVIEW, ATMIError) {
	var buf TypedVIEW
	var itype string
	var subtype string

	// Check the buffer type & get view name
	if _, errA := ac.TpTypes(abuf, &itype, &subtype); nil != errA {
		return nil, errA
	}

	if (itype != "VIEW") && (itype != "VIEW32") {
		return nil, NewCustomATMIError(TPEINVAL, fmt.Sprintf("Invalid buffer type,"+
			" expected VIEW, got [%s]", itype))
	}

    ac.TpLogInfo("Got View: %s", subtype);

    buf.view = subtype
	buf.Buf = abuf

	return &buf, nil
}

//Return int16 value from buffer
//@param cname 	C field name for view
//@param occ	Occurrance
//@param flags	BVACCESS_NOTNULL or 0
//@return int16 val,	 UBF error
func (u *TypedVIEW) BVGetInt16(cname string, occ int, flags int64) (int16, UBFError) {
	var c_val C.short

	//char *view, char *cname
	c_cname := C.CString(cname)
	defer C.free(unsafe.Pointer(c_cname))

	//Get the view name
	c_view := C.CString(u.view)
	defer C.free(unsafe.Pointer(c_view))

	if ret := C.OCBvget(&u.Buf.Ctx.c_ctx, (*C.char)(unsafe.Pointer(u.Buf.C_ptr)),
		c_view, c_cname, C.BFLDOCC(occ),
		(*C.char)(unsafe.Pointer(unsafe.Pointer(&c_val))), nil, BFLD_SHORT, C.long(flags)); ret != SUCCEED {
		return 0, u.Buf.Ctx.NewUBFError()
	}
	return int16(c_val), nil
}

//Return int64 value from buffer
//@param bfldid 	Field ID
//@param occ	Occurrance
//@return int64 val,	 UBF error
func (u *TypedVIEW) BVGetInt64(cname string, occ int, flags int64) (int64, UBFError) {
	var c_val C.long

	//Get the view name
	c_view := C.CString(u.view)
	defer C.free(unsafe.Pointer(c_view))

	//Field name
	c_cname := C.CString(cname)
	defer C.free(unsafe.Pointer(c_cname))

	if ret := C.OCBvget(&u.Buf.Ctx.c_ctx, (*C.char)(unsafe.Pointer(u.Buf.C_ptr)),
		c_view, c_cname,
		C.BFLDOCC(occ), (*C.char)(unsafe.Pointer(unsafe.Pointer(&c_val))),
		nil, BFLD_LONG, C.long(flags)); ret != SUCCEED {
		return 0, u.Buf.Ctx.NewUBFError()
	}
	return int64(c_val), nil
}

//Return int (basicaly C long (int64) casted to) value from buffer
//@param bfldid 	Field ID
//@param occ	Occurrance
//@return int64 val,	 UBF error
func (u *TypedVIEW) BVGetInt(cname string, occ int, flags int64) (int, UBFError) {
	var c_val C.long

	//Get the view name
	c_view := C.CString(u.view)
	defer C.free(unsafe.Pointer(c_view))

	//Field name
	c_cname := C.CString(cname)
	defer C.free(unsafe.Pointer(c_cname))

	if ret := C.OCBvget(&u.Buf.Ctx.c_ctx, (*C.char)(unsafe.Pointer(u.Buf.C_ptr)),
		c_view, c_cname,
		C.BFLDOCC(occ), (*C.char)(unsafe.Pointer(unsafe.Pointer(&c_val))),
		nil, BFLD_INT, C.long(flags)); ret != SUCCEED {
		return 0, u.Buf.Ctx.NewUBFError()
	}
	return int(c_val), nil
}

//Return byte (c char) value from buffer
//@param bfldid 	Field ID
//@param occ	Occurrance
//@return byte val, UBF error
func (u *TypedVIEW) BVGetByte(cname string, occ int, flags int64) (byte, UBFError) {
	var c_val C.char

	//Get the view name
	c_view := C.CString(u.view)
	defer C.free(unsafe.Pointer(c_view))

	//Field name
	c_cname := C.CString(cname)
	defer C.free(unsafe.Pointer(c_cname))

	if ret := C.OCBvget(&u.Buf.Ctx.c_ctx, (*C.char)(unsafe.Pointer(u.Buf.C_ptr)),
		c_view, c_cname,
		C.BFLDOCC(occ), (*C.char)(unsafe.Pointer(unsafe.Pointer(&c_val))),
		nil, BFLD_CHAR, C.long(flags)); ret != SUCCEED {
		return 0, u.Buf.Ctx.NewUBFError()
	}
	return byte(c_val), nil
}

//Get float value from UBF buffer, see CBget(3)
//@param bfldid 	Field ID
//@param occ	Occurrance
//@return float, UBF error
func (u *TypedVIEW) BVGetFloat32(cname string, occ int, flags int64) (float32, UBFError) {
	var c_val C.float

	//Get the view name
	c_view := C.CString(u.view)
	defer C.free(unsafe.Pointer(c_view))

	//Field name
	c_cname := C.CString(cname)
	defer C.free(unsafe.Pointer(c_cname))

	if ret := C.OCBvget(&u.Buf.Ctx.c_ctx, (*C.char)(unsafe.Pointer(u.Buf.C_ptr)),
		c_view, c_cname,
		C.BFLDOCC(occ), (*C.char)(unsafe.Pointer(unsafe.Pointer(&c_val))),
		nil, BFLD_FLOAT, C.long(flags)); ret != SUCCEED {
		return 0, u.Buf.Ctx.NewUBFError()
	}
	return float32(c_val), nil
}

//Get double value
//@param bfldid 	Field ID
//@param occ	Occurrance
//@return double, UBF error
func (u *TypedVIEW) BVGetFloat64(cname string, occ int, flags int64) (float64, UBFError) {
	var c_val C.double

	//Get the view name
	c_view := C.CString(u.view)
	defer C.free(unsafe.Pointer(c_view))

	//Field name
	c_cname := C.CString(cname)
	defer C.free(unsafe.Pointer(c_cname))

	if ret := C.OCBvget(&u.Buf.Ctx.c_ctx, (*C.char)(unsafe.Pointer(u.Buf.C_ptr)),
		c_view, c_cname,
		C.BFLDOCC(occ), (*C.char)(unsafe.Pointer(unsafe.Pointer(&c_val))),
		nil, BFLD_DOUBLE, C.long(flags)); ret != SUCCEED {
		return 0, u.Buf.Ctx.NewUBFError()
	}
	return float64(c_val), nil
}

//Get string value
//@param bfldid 	Field ID
//@param occ	Occurrance
//@return string val, UBF error
func (u *TypedVIEW) BVGetString(cname string, occ int, flags int64) (string, UBFError) {
	var c_len C.BFLDLEN
	c_val := C.malloc(ATMI_MSG_MAX_SIZE)
	c_len = ATMI_MSG_MAX_SIZE

	if nil == c_val {
		return "", NewCustomUBFError(BEUNIX, "Cannot alloc memory")
	}

	defer C.free(c_val)

	//Get the view name
	c_view := C.CString(u.view)
	defer C.free(unsafe.Pointer(c_view))

	//Field name
	c_cname := C.CString(cname)
	defer C.free(unsafe.Pointer(c_cname))

	if ret := C.OCBvget(&u.Buf.Ctx.c_ctx, (*C.char)(unsafe.Pointer(u.Buf.C_ptr)),
		c_view, c_cname,
		C.BFLDOCC(occ), (*C.char)(c_val), &c_len, BFLD_STRING, C.long(flags)); ret != SUCCEED {
		return "", u.Buf.Ctx.NewUBFError()
	}

	return C.GoString((*C.char)(c_val)), nil
}

//Get string value
//@param bfldid 	Field ID
//@param occ	Occurrance
//@return string val, UBF error
func (u *TypedVIEW) BVGetByteArr(cname string, occ int, flags int64) ([]byte, UBFError) {
	var c_len C.BFLDLEN
	c_val := C.malloc(ATMI_MSG_MAX_SIZE)
	c_len = ATMI_MSG_MAX_SIZE

	if nil == c_val {
		return nil, NewCustomUBFError(BEUNIX, "Cannot alloc memory")
	}

	defer C.free(c_val)

	//Get the view name
	c_view := C.CString(u.view)
	defer C.free(unsafe.Pointer(c_view))

	//Field name
	c_cname := C.CString(cname)
	defer C.free(unsafe.Pointer(c_cname))

	if ret := C.OCBvget(&u.Buf.Ctx.c_ctx, (*C.char)(unsafe.Pointer(u.Buf.C_ptr)),
		c_view, c_cname,
		C.BFLDOCC(occ), (*C.char)(c_val),
		&c_len, BFLD_CARRAY, C.long(flags)); ret != SUCCEED {
		return nil, u.Buf.Ctx.NewUBFError()
	}

	g_val := make([]byte, c_len)

	for i := 0; i < int(c_len); i++ {
		g_val[i] = byte(*(*C.char)(unsafe.Pointer(uintptr(c_val) + uintptr(i))))
	}

	return g_val, nil
}

//Change field in buffer
//@param	bfldid	Field ID
//@param ival Input value
//@return UBF Error
func (u *TypedVIEW) BVChg(cname string, occ int, ival interface{}) UBFError {
	return u.BVChgCombined(cname, occ, ival)
}

//Set the field value. Combined supports change (chg) or add mode
//@param	bfldid	Field ID
//@param occ	Field Occurrance
//@param ival Input value
//@param	 do_add Adding mode true = add, false = change
//@return UBF Error
func (u *TypedVIEW) BVChgCombined(cname string, occ int, ival interface{}) UBFError {

	//Get the view name
	c_view := C.CString(u.view)
	defer C.free(unsafe.Pointer(c_view))

	//Field name
	c_cname := C.CString(cname)
	defer C.free(unsafe.Pointer(c_cname))

	switch ival.(type) {
	case int,
		int8,
		int16,
		int32,
		int64,
		uint,
		uint8,
		uint16,
		uint32,
		uint64:
		/* Cast the value to integer... */
		var val int64
		switch ival.(type) {
		case int:
			val = int64(ival.(int))
		case int8:
			val = int64(ival.(int8))
		case int16:
			val = int64(ival.(int16))
		case int32:
			val = int64(ival.(int32))
		case int64:
			val = int64(ival.(int64))
		case uint:
			val = int64(ival.(uint))
		case uint8:
			val = int64(ival.(uint8))
		case uint16:
			val = int64(ival.(uint16))
		case uint32:
			val = int64(ival.(uint32))
		case uint64:
			val = int64(ival.(uint64))
		}
		c_val := C.long(val)

		if ret := C.OCBvchg(&u.Buf.Ctx.c_ctx, (*C.char)(unsafe.Pointer(u.Buf.C_ptr)),
			c_view, c_cname, C.BFLDOCC(occ), (*C.char)(unsafe.Pointer(unsafe.Pointer(&c_val))), 0, BFLD_LONG); ret != SUCCEED {
			return u.Buf.Ctx.NewUBFError()
		}

	case float32:
		fval := ival.(float32)
		c_val := C.float(fval)
		if ret := C.OCBvchg(&u.Buf.Ctx.c_ctx, (*C.char)(unsafe.Pointer(u.Buf.C_ptr)),
			c_view, c_cname, C.BFLDOCC(occ),
			(*C.char)(unsafe.Pointer(unsafe.Pointer(&c_val))), 0, BFLD_FLOAT); ret != SUCCEED {
			return u.Buf.Ctx.NewUBFError()
		}
	case float64:
		dval := ival.(float64)
		c_val := C.double(dval)

		if ret := C.OCBvchg(&u.Buf.Ctx.c_ctx, (*C.char)(unsafe.Pointer(u.Buf.C_ptr)),
			c_view, c_cname, C.BFLDOCC(occ),
			(*C.char)(unsafe.Pointer(unsafe.Pointer(&c_val))), 0, BFLD_DOUBLE); ret != SUCCEED {
			return u.Buf.Ctx.NewUBFError()
		}

	case string:
		str := ival.(string)
		c_val := C.CString(str)
		defer C.free(unsafe.Pointer(c_val))

		if ret := C.OCBvchg(&u.Buf.Ctx.c_ctx, (*C.char)(unsafe.Pointer(u.Buf.C_ptr)),
			c_view, c_cname,
			C.BFLDOCC(occ), c_val, 0, BFLD_STRING); ret != SUCCEED {
			return u.Buf.Ctx.NewUBFError()
		}

	case []byte:
		arr := ival.([]byte)
		c_len := C.BFLDLEN(len(arr))
		c_arr := (*C.char)(unsafe.Pointer(&arr[0]))

		if ret := C.OCBvchg(&u.Buf.Ctx.c_ctx, (*C.char)(unsafe.Pointer(u.Buf.C_ptr)),
			c_view, c_cname,
			C.BFLDOCC(occ), c_arr, c_len, BFLD_CARRAY); ret != SUCCEED {
			return u.Buf.Ctx.NewUBFError()
		}
	default:
		/* TODO: Possibly we could take stuff from println to get string val... */
		return NewCustomUBFError(BEINVAL, "Cannot determine field type")
	}

	return nil
}

//Get the number of field occurrances in buffer
//@param bfldid	Field ID
//@return count (or -1 on error), UBF error
func (u *TypedVIEW) BVOccur(cname string) (int, int, int, int64, int, UBFError) {

	//Get the view name
	c_view := C.CString(u.view)
	defer C.free(unsafe.Pointer(c_view))

	//Field name
	c_cname := C.CString(cname)
	defer C.free(unsafe.Pointer(c_cname))

	var c_maxocc C.BFLDOCC
	var c_realocc C.BFLDOCC
	var c_fldtype C.int
	var c_dim_size C.long

	c_ret := C.OBvoccur(&u.Buf.Ctx.c_ctx,
		(*C.char)(unsafe.Pointer(u.Buf.C_ptr)), c_view, c_cname,
		&c_maxocc, &c_realocc, &c_dim_size, &c_fldtype)

	if FAIL == c_ret {
		return 0, 0, 0, 0, 0, u.Buf.Ctx.NewUBFError()
	}

	return int(c_ret), int(c_maxocc), int(c_realocc), int64(c_dim_size), int(c_fldtype), nil
}

//Get the total buffer size
//@return bufer size, UBF error
func (ac *ATMICtx) BVSizeof(view string) (int64, UBFError) {

	c_view := C.CString(view)
	defer C.free(unsafe.Pointer(c_view))

	c_ret := C.OBvsizeof(&ac.c_ctx, c_view)

	if FAIL == c_ret {
		return FAIL, ac.NewUBFError()
	}

	return int64(c_ret), nil
}

func (u *TypedVIEW) BVSizeof() (int64, UBFError) {

	//Get the view name
	c_view := C.CString(u.view)
	defer C.free(unsafe.Pointer(c_view))

	c_ret := C.OBvsizeof(&u.Buf.Ctx.c_ctx, c_view)

	if FAIL == c_ret {
		return FAIL, u.Buf.Ctx.NewUBFError()
	}

	return int64(c_ret), nil
}

//Delete field (all occurrances) from buffer
//@param bfldid field ID
//@return UBF error
func (u *TypedVIEW) BVSetOccur(cname string, occ int) UBFError {

	//Get the view name
	c_view := C.CString(u.view)
	defer C.free(unsafe.Pointer(c_view))

	//Field name
	c_cname := C.CString(cname)
	defer C.free(unsafe.Pointer(c_cname))

	if ret := C.OBvsetoccur(&u.Buf.Ctx.c_ctx, (*C.char)(unsafe.Pointer(u.Buf.C_ptr)),
		c_view, c_cname, C.BFLDOCC(occ)); SUCCEED != ret {
		return u.Buf.Ctx.NewUBFError()
	}
	return nil
}

//Allocate the new UBF buffer
//NOTE: realloc or other ATMI ops you can do with TypedVIEW.Buf
//@param size - buffer size
//@return Typed UBF, ATMI error
func (ac *ATMICtx) NewVIEW(view string, size int64) (*TypedVIEW, ATMIError) {

	var buf TypedVIEW
	buf.view = view

	if ptr, err := ac.TpAlloc("VIEW", view, size); nil != err {
		return nil, err
	} else {
		buf.Buf = ptr
		buf.Buf.Ctx = ac
		return &buf, nil
	}
}

//Converts string JSON buffer passed in 'buffer' to UBF buffer. This function will
//automatically allocate the free space in UBF to fit the JSON. The size will be
//determinated by string length. See tpjsontoubf(3) C call for more information.
//@param buffer	String buffer containing JSON message. The format must be one level
//JSON containing UBF_FIELD:Value. The value can be array, then it is loaded into
//occurrences.
//@return UBFError ('BEINVAL' if failed to convert, 'BMALLOC' if buffer resize failed)
func (ac *ATMICtx) TpJSONToVIEW(buffer string) (*TypedVIEW, ATMIError) {

	c_buffer := C.CString(buffer)
	defer C.free(unsafe.Pointer(c_buffer))

	c_view := C.malloc(VIEW_NAME_LEN + 1)
	c_view_ptr := (*C.char)(unsafe.Pointer(c_view))
	defer C.free(unsafe.Pointer(c_view))

	var ret *C.char

	if ret = C.Otpjsontoview(&ac.c_ctx, c_view_ptr, c_buffer); ret == nil {
		return nil, ac.NewATMIError()
	}

	var atmiBuf ATMIBuf
	atmiBuf.C_ptr = ret
	atmiBuf.Ctx = ac
	view := C.GoString(c_view_ptr)
	len, errA := ac.BVSizeof(view)

	if nil != errA {
		return nil, errA
	}

	var tv TypedVIEW

	atmiBuf.C_len = C.long(len)
	tv.Buf = &atmiBuf
	tv.view = view

	return &tv, nil
}

//Convert given UBF buffer to JSON block, see tpubftojson(3) C call
//Output string is automatically allocated
//@return JSON string (if converted ok), ATMIError in case of failure. More detailed
//infos in case of error is found in 'ubf' and 'ndrx' facility logs.
func (u *TypedVIEW) TpVIEWToJSON(flags int64) (string, ATMIError) {

	//Get the view name
	c_view := C.CString(u.view)
	defer C.free(unsafe.Pointer(c_view))

	used, _ := u.BVSizeof()
	ret_size := used * 10

	u.Buf.Ctx.ndrxLog(LOG_INFO, "TpVIEWToJSON: sizeof %d allocating %d",
		used, ret_size)

	c_buffer := C.malloc(C.size_t(ret_size))

	if nil == c_buffer {
		return "", NewCustomUBFError(BEUNIX, "Cannot alloc memory")
	}

	defer C.free(c_buffer)

	if ret := C.Otpviewtojson(&u.Buf.Ctx.c_ctx, (*C.char)(unsafe.Pointer(u.Buf.C_ptr)),
		c_view, (*C.char)(unsafe.Pointer(c_buffer)), C.int(ret_size), C.long(flags)); ret != 0 {
		return "", u.Buf.Ctx.NewUBFError()
	}

	return C.GoString((*C.char)(c_buffer)), nil

}

//Iterate over the view structure
func (u *TypedVIEW) BVNext(state *BVNextState, start bool) (int, string, int, int, int64, UBFError) {

	var c_view *C.char = nil

	if start {
		c_view = C.CString(u.view)
		defer C.free(unsafe.Pointer(c_view))
	}

	c_cname := C.malloc(VIEW_CNAME_LEN + 1)
	c_cname_ptr := (*C.char)(unsafe.Pointer(c_cname))
	defer C.free(unsafe.Pointer(c_cname))

	var c_fldtype C.int
	var c_maxocc C.BFLDOCC
	var c_dim_size C.long

	if ret := C.OBvnext(&u.Buf.Ctx.c_ctx, (*C.struct_Bvnext_state)(unsafe.Pointer(&state.state)),
		c_view, c_cname_ptr, &c_fldtype, &c_maxocc, &c_dim_size); ret >= 0 {
		return int(ret), C.GoString(c_cname_ptr), int(c_fldtype),
			int(c_maxocc), int64(c_dim_size), nil
	}

	//We have a failure
	return -1, "", 0, 0, 0, u.Buf.Ctx.NewUBFError()

}

///////////////////////////////////////////////////////////////////////////////////
// Wrappers for memory management
///////////////////////////////////////////////////////////////////////////////////

func (v *TypedVIEW) TpRealloc(size int64) ATMIError {
	return v.Buf.TpRealloc(size)
}
