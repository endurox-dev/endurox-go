/**
 * @brief Unified Buffer Format (UBF) - Key value protocol buffer support
 *
 * @file typed_ubf.go
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

#include <errno.h>
#include <xatmi.h>
#include <string.h>
#include <stdlib.h>
#include <ubf.h>
#include <ndebug.h>
#include <odebug.h>
#include <oubf.h>
#include <oatmi.h>
#include <sys_unix.h>
#include <ndrstandard.h>


//Get the UBF Error code
static int WrapBerror(TPCONTEXT_T *p_ctx) {
	return OBerror(p_ctx);
}

//Get the value with buffer allocation
static char * c_Bget_str (UBFH * p_ub, BFLDID bfldid, BFLDOCC occ,
					BFLDLEN *len, int *err_code)
{
	char *ret = malloc(NDRX_MSGSIZEMAX);

	*len = NDRX_MSGSIZEMAX;
	*err_code = 0;

	if (NULL==ret)
	{
		*err_code = 1; //memory error
	}
	else
	{
		if (0!=Bget (p_ub, bfldid, occ, ret, len))
		{
			*err_code = 2; //Buffer
			free(ret);
			ret = NULL;
		}
	}

	return ret;
}

//Get integer size
static int c_sizeof_BFLDID(void)
{
	return sizeof(BFLDID);
}

//Go proxy function for expression evaluator
extern long go_expr_callback_proxy(char *buf, char *funcname);

//Proxy function for expression callback
static long c_expr_callback_proxy(UBFH *p_ub, char *funcname)
{
	//Call the service entry
	return go_expr_callback_proxy((char *)p_ub, funcname);
}

static int c_proxy_Bboolsetcbf(TPCONTEXT_T *p_ctx, char *funcname)
{
	return OBboolsetcbf(p_ctx, funcname, c_expr_callback_proxy);
}

//Will run the Btreefree in temp context
static void go_Btreefree(char *ptr)
{

    // Allocate new context + set it...
    TPCONTEXT_T c = tpnewctxt(0, 1);
    Btreefree(ptr);
    tpfreectxt(c);
}

//Reset location infos
static void reset_loc_info(Bfld_loc_info_t *loc)
{
	memset((void *)&loc, 0, sizeof(Bfld_loc_info_t));
}

typedef struct bfprintcb_data bfprintcb_data_t;
struct bfprintcb_data
{
	TPCONTEXT_T *p_ctx;
	char *buf;
	long cur_offset;
	long size;
};

//Callback for writting data to
//@param buffer output buffer to write
//@param datalen printed data including EOS
//@param dataptr1 custom data pointer, in this case bfprintcb_data_t
//@return EXSUCCEED/EXFAIL
static int bfprintcb_writef(char **buffer, long datalen, void *dataptr1,
	int *do_write, FILE *outf, BFLDID fid)
{
	int ret = EXSUCCEED;

	bfprintcb_data_t *data = (bfprintcb_data_t *)dataptr1;

	//-1 for skipping the EOS
	if (datalen > (data->size - data->cur_offset)-1)
	{
		OUBF_LOG(data->p_ctx, log_error, "Output buffer full: free %ld, new data: %ld",
			(long)((data->size - data->cur_offset)-1), datalen);
		EXFAIL_OUT(ret);
	}

	//now copy off the data

	if (data->cur_offset>0)
	{
		data->cur_offset--;
	}

	memcpy(data->buf + data->cur_offset, *buffer, datalen);

	data->cur_offset+=datalen;

out:

	return ret;
}

//Print UBF buffer to allocated string
//@param p_ctx ATMI Context
//@param p_ub buffer to print
//@return ptr to C allocate string with print data terminated with EOS
//	or NULL in case of error.
static char * BPrintStrC(TPCONTEXT_T *p_ctx, UBFH * p_ub)
{
	bfprintcb_data_t data;

	memset(&data, 0, sizeof(data));

	//Allocate the buffer, so lets allocate some heavy buffer for the string..
	//Also note that while we do not use the buffer, the good os actually won't
	//allocate any real resources. Thus I guess no problem here
	//just a virtual memory...

	data.size = Bsizeof(p_ub) * MAXTIDENT;
	data.buf = malloc(data.size);
	data.p_ctx = p_ctx;

	if (NULL==data.buf)
	{
		int err = errno;
		OUBF_LOG(p_ctx, log_error, "Failed to allocate %ld bytes for print buffer: %s",
				data.size, strerror(err));

		goto out;
	}

	if (EXSUCCEED!=OBfprintcb(p_ctx, p_ub, bfprintcb_writef, (char *)&data))
	{
		int err = errno;
		OUBF_LOG(p_ctx, log_error, "Failed to print to buffer / callback fail!");
		free(data.buf);
		data.buf = NULL;
		goto out;
	}

out:

	return data.buf;
}

*/
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

//Field types
const (
	BFLD_MIN    = 0
	BFLD_SHORT  = 0
	BFLD_LONG   = 1
	BFLD_CHAR   = 2
	BFLD_FLOAT  = 3
	BFLD_DOUBLE = 4
	BFLD_STRING = 5
	BFLD_CARRAY = 6
	BFLD_MAX    = 6

	BFLD_INT = 7 /* used for views only */
)

//Error codes
const (
	BMINVAL   = 0 /* min error */
	BERFU0    = 1
	BALIGNERR = 2
	BNOTFLD   = 3
	BNOSPACE  = 4
	BNOTPRES  = 5
	BBADFLD   = 6
	BTYPERR   = 7
	BEUNIX    = 8
	BBADNAME  = 9
	BMALLOC   = 10
	BSYNTAX   = 11
	BFTOPEN   = 12
	BFTSYNTAX = 13
	BEINVAL   = 14
	BERFU1    = 15
	BBADTBL   = 16
	BBADVIEW  = 17
	BVFSYNTAX = 18
	BVFOPEN   = 19
	BBADACM   = 20
	BNOCNAME  = 21
	BEBADOP   = 22

	BMAXVAL = 22 /* max error */
)

const (
	BBADFLDID   = 0
	BFIRSTFLDID = 0
)

///////////////////////////////////////////////////////////////////////////////////
// Buffer def, typedefs
///////////////////////////////////////////////////////////////////////////////////

//UBF Buffer
type TypedUBF struct {
	Buf *ATMIBuf
}

//Return The ATMI buffer to caller
func (u *TypedUBF) GetBuf() *ATMIBuf {
	return u.Buf
}

//Compiled Expression Tree
type ExprTree struct {
	//All object which have finalizer we need nop func to defer so that
	//during the c call the GC does not collect the object...
	gcoff int
	c_ptr *C.char
}

//Field location infos
type BFldLocInfo struct {
	loc C.Bfld_loc_info_t
}

///////////////////////////////////////////////////////////////////////////////////
// Error Handlers
///////////////////////////////////////////////////////////////////////////////////

//ATMI Error type
type ubfError struct {
	code    int
	message string
}

//ATMI error interface
type UBFError interface {
	Error() string
	Code() int
	Message() string
}

//Generate UBF error, read the codes
func (ac *ATMICtx) NewUBFError() UBFError {
	var err ubfError
	err.code = int(C.WrapBerror(&ac.c_ctx))
	err.message = C.GoString(C.OBstrerror(&ac.c_ctx, C.WrapBerror(&ac.c_ctx)))
	return err
}

//Build a custom error
//@param err		Error buffer to build
//@param code	Error code to setup
//@param msg		Error message
func NewCustomUBFError(code int, msg string) UBFError {
	var err ubfError
	err.code = code
	err.message = msg
	return err
}

//Standard error interface
func (e ubfError) Error() string {
	return fmt.Sprintf("%d: %s", e.code, e.message)
}

//code getter
func (e ubfError) Code() int {
	return e.code
}

//message getter
func (e ubfError) Message() string {
	return e.message
}

///////////////////////////////////////////////////////////////////////////////////
// Globals
///////////////////////////////////////////////////////////////////////////////////
type UBFExprFunc func(buf *TypedUBF, funcname string) int64

var exprfuncmap map[string]UBFExprFunc //callback mapping for UBF expression functions to go

///////////////////////////////////////////////////////////////////////////////////
// UBF API
///////////////////////////////////////////////////////////////////////////////////

//Do nothing, to trick the GC
func (expr *ExprTree) nop() int {
	expr.gcoff++
	return expr.gcoff
}

//Get the field len
//@param fldid	Field ID
//@param occ 	Field occurance
//@return 	FIeld len, UBF error
func (u *TypedUBF) BLen(bfldid int, occ int) (int, UBFError) {

	ret := C.OBlen(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), C.BFLDID(bfldid), C.BFLDOCC(occ))
	if FAIL == ret {
		return FAIL, u.Buf.Ctx.NewUBFError()
	}

	u.Buf.nop()

	return int(ret), nil
}

//Delete the field from buffer
//@param fldid	Field ID
//@param occ 	Field occurance
//@return 	UBF error
func (u *TypedUBF) BDel(bfldid int, occ int) UBFError {

	ret := C.OBdel(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), C.BFLDID(bfldid), C.BFLDOCC(occ))
	if FAIL == ret {
		return u.Buf.Ctx.NewUBFError()
	}

	u.Buf.nop()

	return nil
}

//Make a project copy of the fields (leave only those in array)
//@return UBF error
func (u *TypedUBF) BProj(fldlist []int) UBFError {

	c_fldidsize := int(C.c_sizeof_BFLDID())

	c_val := C.malloc(C.size_t(c_fldidsize * (len(fldlist) + 1)))

	if nil == c_val {
		return NewCustomUBFError(BEUNIX, "Cannot alloc memory")
	}

	defer C.free(c_val)
	var i int

	for i = 0; i < len(fldlist); i++ {
		*(*C.BFLDID)(unsafe.Pointer(uintptr(c_val) + uintptr(i*c_fldidsize))) =
			C.BFLDID(fldlist[i])
	}

	//Set last field to BBADFLDID
	*(*C.BFLDID)(unsafe.Pointer(uintptr(c_val) + uintptr(i*c_fldidsize))) = C.BFLDID(BBADFLDID)

	if ret := C.OBproj(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)),
		(*C.BFLDID)(unsafe.Pointer(c_val))); ret != SUCCEED {
		return u.Buf.Ctx.NewUBFError()
	}

	u.Buf.nop()

	return nil
}

//Make a project copy of the fields (leave only those in array)
//The terminator in the fildlist array are not required. The list shall
//contain fields only to copy. See *Bprojcpy(3)* for more details.
//NOTE! The dest buffer is erased before copying new data to
//@param dest is destination buffer
//@param src is source buffer to copy from
//@param fldlist list of fields to copy
//@return UBF error
func (ac *ATMICtx) BProjCpy(dest *TypedUBF, src *TypedUBF, fldlist []int) UBFError {

	c_fldidsize := int(C.c_sizeof_BFLDID())
	c_val := C.malloc(C.size_t(c_fldidsize * (len(fldlist) + 1)))

	if nil == c_val {
		return NewCustomUBFError(BEUNIX, "Cannot alloc memory")
	}

	defer C.free(c_val)
	var i int

	for i = 0; i < len(fldlist); i++ {
		*(*C.BFLDID)(unsafe.Pointer(uintptr(c_val) + uintptr(i*c_fldidsize))) =
			C.BFLDID(fldlist[i])
	}
	//Set last field to BBADFLDID
	*(*C.BFLDID)(unsafe.Pointer(uintptr(c_val) + uintptr(i*c_fldidsize))) = C.BFLDID(BBADFLDID)

	if ret := C.OBprojcpy(&ac.c_ctx, (*C.UBFH)(unsafe.Pointer(dest.Buf.C_ptr)),
		(*C.UBFH)(unsafe.Pointer(src.Buf.C_ptr)),
		(*C.BFLDID)(unsafe.Pointer(c_val))); ret != SUCCEED {
		return ac.NewUBFError()
	}

	dest.Buf.nop()
	src.Buf.nop()
	ac.nop()

	return nil
}

//Return field ID
//@param fldnm	Field name
//@return Field ID, UBF error
func (ac *ATMICtx) BFldId(fldnm string) (int, UBFError) {

	c_fldnm := C.CString(fldnm)

	defer C.free(unsafe.Pointer(c_fldnm))

	ret := C.OBfldid(&ac.c_ctx, c_fldnm)

	if FAIL == ret {
		return BBADFLDID, ac.NewUBFError()
	}

	ac.nop()

	return int(ret), nil
}

//Get field name
//@param bfldid Field ID
//@return Field name (or "" if error), UBF error
func (ac *ATMICtx) BFname(bfldid int) (string, UBFError) {

	ret := C.OBfname(&ac.c_ctx, C.BFLDID(bfldid))

	if nil == ret {
		return "", ac.NewUBFError()
	}

	ac.nop()

	return C.GoString(ret), nil
}

//Check for field presence in buffer
//@param fldid	Field ID
//@param occ 	Field occurance
//@return 	true/false present/not present
func (u *TypedUBF) BPres(bfldid int, occ int) bool {

	ret := C.OBpres(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)),
		C.BFLDID(bfldid), C.BFLDOCC(occ))
	if 1 == ret {
		return true
	}

	u.Buf.nop()

	return false
}

//Copy buffer
//@param dest Destination UBF buffer
//@param src		Source UBF buffer
//@return UBF error
func (ac *ATMICtx) BCpy(dest *TypedUBF, src *TypedUBF) UBFError {

	if ret := C.OBcpy(&ac.c_ctx, (*C.UBFH)(unsafe.Pointer(dest.Buf.C_ptr)),
		(*C.UBFH)(unsafe.Pointer(src.Buf.C_ptr))); SUCCEED != ret {
		return ac.NewUBFError()
	}

	dest.Buf.nop()
	src.Buf.nop()
	ac.nop()

	return nil
}

//Iterate over the buffer
//NOTE: This is not multiple context safe. It stores iteration state internally
//@param first	TRUE start iteration, FALSE continue iteration
//@return Field ID, Field Occurrance, UBF Error
func (u *TypedUBF) BNext(first bool) (int, int, UBFError) {

	var fldid C.BFLDID
	var occ C.BFLDOCC

	if first {
		fldid = BFIRSTFLDID
	} else {
		//Get next saved in internal ptr in library
		fldid = FAIL
	}

	if ret := C.OBnext(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)),
		&fldid, &occ, nil, nil); 1 != ret {
		return FAIL, FAIL, u.Buf.Ctx.NewUBFError()
	}

	u.Buf.nop()
	return int(fldid), int(occ), nil
}

//Initialize/re-initialize UBF buffer
//@param u UBF buffer
//@param ulen	lenght of the buffer
//@return UBF error
func (ac *ATMICtx) BInit(u *TypedUBF, ulen int64) UBFError {

	if ret := C.OBinit(&ac.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)),
		C.BFLDLEN(ulen)); SUCCEED != ret {
		return ac.NewUBFError()
	}

	ac.nop()
	u.Buf.nop()

	return nil
}

//Allocate the UBF buffer
//@param size	Buffer size in bytes
//@return UBF Handler, ATMI Error
func (ac *ATMICtx) UBFAlloc(size int64) (TypedUBF, ATMIError) {

	var err ATMIError
	var buf TypedUBF
	buf.Buf, err = ac.TpAlloc("UBF", "", size)

	ac.nop()

	return buf, err
}

//Get the UBF Handler
func (ac *ATMICtx) CastToUBF(abuf *ATMIBuf) (*TypedUBF, ATMIError) {
	var buf TypedUBF

	//TODO: Check the buffer type!
	buf.Buf = abuf

	abuf.nop()
	ac.nop()

	return &buf, nil
}

//Get the field form buffer. This returns the interface to underlaying type
//@param bfldid 	Field ID
//@param occ	Occurrance
//@return interface to value,	 UBF error
func (u *TypedUBF) BGet(bfldid int, occ int) (interface{}, UBFError) {

	/* Determinte the type of the buffer */
	switch u.Buf.Ctx.BFldType(bfldid) {
	case BFLD_SHORT:
		var c_val C.short
		if ret := C.OBget(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), C.BFLDID(bfldid),
			C.BFLDOCC(occ), (*C.char)(unsafe.Pointer(unsafe.Pointer(&c_val))), nil); ret != SUCCEED {
			return nil, u.Buf.Ctx.NewUBFError()
		}
		u.Buf.nop()
		return int16(c_val), nil
	case BFLD_LONG:
		var c_val C.long
		if ret := C.OBget(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), C.BFLDID(bfldid),
			C.BFLDOCC(occ), (*C.char)(unsafe.Pointer(unsafe.Pointer(&c_val))), nil); ret != SUCCEED {
			return nil, u.Buf.Ctx.NewUBFError()
		}
		u.Buf.nop()
		return int64(c_val), nil
	case BFLD_CHAR: /* This is single byte... */
		var c_val C.char
		if ret := C.OBget(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), C.BFLDID(bfldid),
			C.BFLDOCC(occ), (*C.char)(unsafe.Pointer(unsafe.Pointer(&c_val))), nil); ret != SUCCEED {
			return nil, u.Buf.Ctx.NewUBFError()
		}
		u.Buf.nop()
		return byte(c_val), nil
	case BFLD_FLOAT:
		var c_val C.float
		if ret := C.OBget(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), C.BFLDID(bfldid),
			C.BFLDOCC(occ), (*C.char)(unsafe.Pointer(unsafe.Pointer(&c_val))), nil); ret != SUCCEED {
			return nil, u.Buf.Ctx.NewUBFError()
		}
		u.Buf.nop()
		return float32(c_val), nil
	case BFLD_DOUBLE:
		var c_val C.double
		if ret := C.OBget(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), C.BFLDID(bfldid),
			C.BFLDOCC(occ), (*C.char)(unsafe.Pointer(unsafe.Pointer(&c_val))), nil); ret != SUCCEED {
			return nil, u.Buf.Ctx.NewUBFError()
		}
		return float64(c_val), nil
	case BFLD_STRING:
		var c_len C.BFLDLEN
		c_val := C.malloc(C.size_t(ATMIMsgSizeMax()))
		c_len = C.BFLDLEN(ATMIMsgSizeMax())

		if nil == c_val {
			return nil, NewCustomUBFError(BEUNIX, "Cannot alloc memory")
		}

		defer C.free(c_val)

		if ret := C.OBget(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), C.BFLDID(bfldid),
			C.BFLDOCC(occ), (*C.char)(c_val), &c_len); ret != SUCCEED {
			return nil, u.Buf.Ctx.NewUBFError()
		}

		u.Buf.nop()
		return C.GoString((*C.char)(c_val)), nil

	case BFLD_CARRAY:
		var c_len C.BFLDLEN
		c_val := C.malloc(C.size_t(ATMIMsgSizeMax()))
		c_len = C.BFLDLEN(ATMIMsgSizeMax())

		if nil == c_val {
			return nil, NewCustomUBFError(BEUNIX, "Cannot alloc memory")
		}

		defer C.free(c_val)

		if ret := C.OBget(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), C.BFLDID(bfldid),
			C.BFLDOCC(occ), (*C.char)(c_val), &c_len); ret != SUCCEED {
			return nil, u.Buf.Ctx.NewUBFError()
		}

		g_val := make([]byte, c_len)

		for i := 0; i < int(c_len); i++ {
			g_val[i] = byte(*(*C.char)(unsafe.Pointer(uintptr(c_val) + uintptr(i))))
		}

		u.Buf.nop()
		return g_val, nil

	}
	/* Default case... */
	return nil, NewCustomUBFError(BEINVAL, "Invalid field")
}

//Return int16 value from buffer
//@param bfldid 	Field ID
//@param occ	Occurrance
//@return int16 val,	 UBF error
func (u *TypedUBF) BGetInt16(bfldid int, occ int) (int16, UBFError) {

	var c_val C.short
	if ret := C.OCBget(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), C.BFLDID(bfldid),
		C.BFLDOCC(occ), (*C.char)(unsafe.Pointer(unsafe.Pointer(&c_val))), nil, BFLD_SHORT); ret != SUCCEED {
		return 0, u.Buf.Ctx.NewUBFError()
	}

	u.Buf.nop()
	return int16(c_val), nil
}

//Return int64 value from buffer
//@param bfldid 	Field ID
//@param occ	Occurrance
//@return int64 val,	 UBF error
func (u *TypedUBF) BGetInt64(bfldid int, occ int) (int64, UBFError) {

	var c_val C.long
	if ret := C.OCBget(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), C.BFLDID(bfldid),
		C.BFLDOCC(occ), (*C.char)(unsafe.Pointer(unsafe.Pointer(&c_val))), nil, BFLD_LONG); ret != SUCCEED {
		return 0, u.Buf.Ctx.NewUBFError()
	}

	u.Buf.nop()
	return int64(c_val), nil
}

//Return int (basicaly C long (int64) casted to) value from buffer
//@param bfldid 	Field ID
//@param occ	Occurrance
//@return int64 val,	 UBF error
func (u *TypedUBF) BGetInt(bfldid int, occ int) (int, UBFError) {

	var c_val C.long
	if ret := C.OCBget(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), C.BFLDID(bfldid),
		C.BFLDOCC(occ), (*C.char)(unsafe.Pointer(unsafe.Pointer(&c_val))), nil, BFLD_LONG); ret != SUCCEED {
		return 0, u.Buf.Ctx.NewUBFError()
	}

	u.Buf.nop()
	return int(c_val), nil
}

//Return byte (c char) value from buffer
//@param bfldid 	Field ID
//@param occ	Occurrance
//@return byte val, UBF error
func (u *TypedUBF) BGetByte(bfldid int, occ int) (byte, UBFError) {

	var c_val C.char
	if ret := C.OCBget(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), C.BFLDID(bfldid),
		C.BFLDOCC(occ), (*C.char)(unsafe.Pointer(unsafe.Pointer(&c_val))), nil, BFLD_CHAR); ret != SUCCEED {
		return 0, u.Buf.Ctx.NewUBFError()
	}

	u.Buf.nop()
	return byte(c_val), nil
}

//Get float value from UBF buffer, see CBget(3)
//@param bfldid 	Field ID
//@param occ	Occurrance
//@return float, UBF error
func (u *TypedUBF) BGetFloat32(bfldid int, occ int) (float32, UBFError) {

	var c_val C.float
	if ret := C.OCBget(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), C.BFLDID(bfldid),
		C.BFLDOCC(occ), (*C.char)(unsafe.Pointer(unsafe.Pointer(&c_val))), nil, BFLD_FLOAT); ret != SUCCEED {
		return 0, u.Buf.Ctx.NewUBFError()
	}

	u.Buf.nop()
	return float32(c_val), nil
}

//Get double value
//@param bfldid 	Field ID
//@param occ	Occurrance
//@return double, UBF error
func (u *TypedUBF) BGetFloat64(bfldid int, occ int) (float64, UBFError) {

	var c_val C.double
	if ret := C.OCBget(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), C.BFLDID(bfldid),
		C.BFLDOCC(occ), (*C.char)(unsafe.Pointer(unsafe.Pointer(&c_val))), nil, BFLD_DOUBLE); ret != SUCCEED {
		return 0, u.Buf.Ctx.NewUBFError()
	}

	u.Buf.nop()
	return float64(c_val), nil
}

//Get string value
//@param bfldid 	Field ID
//@param occ	Occurrance
//@return string val, UBF error
func (u *TypedUBF) BGetString(bfldid int, occ int) (string, UBFError) {

	var c_len C.BFLDLEN
	c_val := C.malloc(C.size_t(ATMIMsgSizeMax()))
	c_len = C.BFLDLEN(ATMIMsgSizeMax())

	if nil == c_val {
		return "", NewCustomUBFError(BEUNIX, "Cannot alloc memory")
	}

	defer C.free(c_val)

	if ret := C.OCBget(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), C.BFLDID(bfldid),
		C.BFLDOCC(occ), (*C.char)(c_val), &c_len, BFLD_STRING); ret != SUCCEED {
		return "", u.Buf.Ctx.NewUBFError()
	}

	u.Buf.nop()
	return C.GoString((*C.char)(c_val)), nil
}

//Get string value
//@param bfldid 	Field ID
//@param occ	Occurrance
//@return string val, UBF error
func (u *TypedUBF) BGetByteArr(bfldid int, occ int) ([]byte, UBFError) {

	var c_len C.BFLDLEN
	c_val := C.malloc(C.size_t(ATMIMsgSizeMax()))
	c_len = C.BFLDLEN(ATMIMsgSizeMax())

	if nil == c_val {
		return nil, NewCustomUBFError(BEUNIX, "Cannot alloc memory")
	}

	defer C.free(c_val)

	if ret := C.OCBget(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), C.BFLDID(bfldid),
		C.BFLDOCC(occ), (*C.char)(c_val), &c_len, BFLD_CARRAY); ret != SUCCEED {
		return nil, u.Buf.Ctx.NewUBFError()
	}

	g_val := make([]byte, c_len)

	for i := 0; i < int(c_len); i++ {
		g_val[i] = byte(*(*C.char)(unsafe.Pointer(uintptr(c_val) + uintptr(i))))
	}

	u.Buf.nop()
	return g_val, nil
}

//Change field in buffer
//@param	bfldid	Field ID
//@param ival Input value
//@return UBF Error
func (u *TypedUBF) BChg(bfldid int, occ int, ival interface{}) UBFError {
	return u.BChgCombined(bfldid, occ, ival, false)
}

//Fast add of filed to buffer (assuming buffer not changed and adding the same
//type of field. NOTE ! Types must be matched with UBF field type
//@param bfldid field id to add
//@param ival value to add
//@param loc location data (last saved or new data) - initialized by first flag
//@param first set to true, if 'loc' is not inialised
//@return UBF error or nil
func (u *TypedUBF) BAddFast(bfldid int, ival interface{}, loc *BFldLocInfo, first bool) UBFError {

	if nil == loc {
		return NewCustomUBFError(BEINVAL, "loc cannot be nil!")
	}

	if first {
		C.reset_loc_info(&loc.loc)
	}

	fldtyp := u.GetBuf().Ctx.BFldType(bfldid)

	switch ival.(type) {
	case int8,
		int16,
		uint8,
		uint16:

		//Validate type code
		var val int16

		if BFLD_SHORT != fldtyp {
			return NewCustomUBFError(BEINVAL, fmt.Sprintf("expected BFLD_SHORT got: %d",
				fldtyp))
		}

		switch ival.(type) {
		case int8:
			val = int16(ival.(int8))
		case int16:
			val = int16(ival.(int16))
		case uint8:
			val = int16(ival.(uint8))
		case uint16:
			val = int16(ival.(uint16))
		}

		c_val := C.short(val)

		if ret := C.OBaddfast(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), C.BFLDID(bfldid),
			(*C.char)(unsafe.Pointer(unsafe.Pointer(&c_val))), 0, &loc.loc); ret != SUCCEED {
			return u.Buf.Ctx.NewUBFError()
		}

		break

	case int32,
		int,
		uint,
		int64,
		uint32,
		uint64:
		/* Cast the value to integer... */
		var val int64

		if BFLD_LONG != fldtyp {
			return NewCustomUBFError(BEINVAL, fmt.Sprintf("expected BFLD_LONG got: %d",
				fldtyp))
		}

		switch ival.(type) {
		case int:
			val = int64(ival.(int))
		case int32:
			val = int64(ival.(int32))
		case int64:
			val = int64(ival.(int64))
		case uint:
			val = int64(ival.(uint))
		case uint32:
			val = int64(ival.(uint32))
		case uint64:
			val = int64(ival.(uint64))
		}
		c_val := C.long(val)

		if ret := C.OBaddfast(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), C.BFLDID(bfldid),
			(*C.char)(unsafe.Pointer(unsafe.Pointer(&c_val))), 0, &loc.loc); ret != SUCCEED {
			return u.Buf.Ctx.NewUBFError()
		}

	case float32:

		if BFLD_FLOAT != fldtyp {
			return NewCustomUBFError(BEINVAL, fmt.Sprintf("expected BFLD_FLOAT got: %d",
				fldtyp))
		}

		fval := ival.(float32)
		c_val := C.float(fval)

		if ret := C.OBaddfast(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), C.BFLDID(bfldid),
			(*C.char)(unsafe.Pointer(unsafe.Pointer(&c_val))), 0, &loc.loc); ret != SUCCEED {
			return u.Buf.Ctx.NewUBFError()
		}
	case float64:

		if BFLD_DOUBLE != fldtyp {
			return NewCustomUBFError(BEINVAL, fmt.Sprintf("expected BFLD_DOUBLE got: %d",
				fldtyp))
		}

		dval := ival.(float64)
		c_val := C.double(dval)

		if ret := C.OBaddfast(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), C.BFLDID(bfldid),
			(*C.char)(unsafe.Pointer(unsafe.Pointer(&c_val))), 0, &loc.loc); ret != SUCCEED {
			return u.Buf.Ctx.NewUBFError()
		}

	case string:

		if BFLD_STRING != fldtyp && BFLD_CHAR != fldtyp {
			return NewCustomUBFError(BEINVAL, fmt.Sprintf("expected BFLD_STRING or "+
				"BFLD_CHAR but got: %d",
				fldtyp))
		}

		str := ival.(string)
		c_val := C.CString(str)
		defer C.free(unsafe.Pointer(c_val))

		if ret := C.OBaddfast(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), C.BFLDID(bfldid),
			c_val, 0, &loc.loc); ret != SUCCEED {
			return u.Buf.Ctx.NewUBFError()
		}

	case []byte:

		if BFLD_CARRAY != fldtyp {
			return NewCustomUBFError(BEINVAL, fmt.Sprintf("expected BFLD_CARRAY or "+
				"BFLD_CHAR but got: %d",
				fldtyp))
		}

		arr := ival.([]byte)
		c_len := C.BFLDLEN(len(arr))
		c_arr := (*C.char)(unsafe.Pointer(&arr[0]))

		if ret := C.OBaddfast(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), C.BFLDID(bfldid),
			c_arr, c_len, &loc.loc); ret != SUCCEED {
			return u.Buf.Ctx.NewUBFError()
		}
	default:
		/* TODO: Possibly we could take stuff from println to get string val... */
		return NewCustomUBFError(BEINVAL, "Cannot determine field type")
	}

	u.Buf.nop()
	return nil
}

//Add field to buffer
//@param	bfldid	Field ID
//@param ival Input value
//@return UBF Error
func (u *TypedUBF) BAdd(bfldid int, ival interface{}) UBFError {

	return u.BChgCombined(bfldid, 0, ival, true)
}

//Set the field value. Combined supports change (chg) or add mode
//@param	bfldid	Field ID
//@param occ	Field Occurrance
//@param ival Input value
//@param	 do_add Adding mode true = add, false = change
//@return UBF Error
func (u *TypedUBF) BChgCombined(bfldid int, occ int, ival interface{}, do_add bool) UBFError {

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

		if do_add {
			if ret := C.OCBadd(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), C.BFLDID(bfldid),
				(*C.char)(unsafe.Pointer(unsafe.Pointer(&c_val))), 0, BFLD_LONG); ret != SUCCEED {
				return u.Buf.Ctx.NewUBFError()
			}
		} else {
			if ret := C.OCBchg(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), C.BFLDID(bfldid),
				C.BFLDOCC(occ), (*C.char)(unsafe.Pointer(unsafe.Pointer(&c_val))), 0, BFLD_LONG); ret != SUCCEED {
				return u.Buf.Ctx.NewUBFError()
			}
		}
	case float32:
		fval := ival.(float32)
		c_val := C.float(fval)
		if do_add {
			if ret := C.OCBadd(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), C.BFLDID(bfldid),
				(*C.char)(unsafe.Pointer(unsafe.Pointer(&c_val))), 0, BFLD_FLOAT); ret != SUCCEED {
				return u.Buf.Ctx.NewUBFError()
			}
		} else {
			if ret := C.OCBchg(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), C.BFLDID(bfldid),
				C.BFLDOCC(occ), (*C.char)(unsafe.Pointer(unsafe.Pointer(&c_val))), 0, BFLD_FLOAT); ret != SUCCEED {
				return u.Buf.Ctx.NewUBFError()
			}
		}
	case float64:
		dval := ival.(float64)
		c_val := C.double(dval)
		if do_add {
			if ret := C.OCBadd(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), C.BFLDID(bfldid),
				(*C.char)(unsafe.Pointer(unsafe.Pointer(&c_val))), 0, BFLD_DOUBLE); ret != SUCCEED {
				return u.Buf.Ctx.NewUBFError()
			}
		} else {
			if ret := C.OCBchg(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), C.BFLDID(bfldid),
				C.BFLDOCC(occ), (*C.char)(unsafe.Pointer(unsafe.Pointer(&c_val))), 0, BFLD_DOUBLE); ret != SUCCEED {
				return u.Buf.Ctx.NewUBFError()
			}
		}
	case string:
		str := ival.(string)
		c_val := C.CString(str)
		defer C.free(unsafe.Pointer(c_val))
		if do_add {
			if ret := C.OCBadd(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), C.BFLDID(bfldid),
				c_val, 0, BFLD_STRING); ret != SUCCEED {
				return u.Buf.Ctx.NewUBFError()
			}
		} else {
			if ret := C.OCBchg(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), C.BFLDID(bfldid),
				C.BFLDOCC(occ), c_val, 0, BFLD_STRING); ret != SUCCEED {
				return u.Buf.Ctx.NewUBFError()
			}
		}
	case []byte:
		arr := ival.([]byte)
		c_len := C.BFLDLEN(len(arr))
		var c_arr *C.char

		if c_len > 0 {
			c_arr = (*C.char)(unsafe.Pointer(&arr[0]))
		} else {
			dumdata := [...]byte{0x0}
			// set some pointer..., not used really as len is 0, but we need some ptr
			c_arr = (*C.char)(unsafe.Pointer(&dumdata[0]))
		}

		if do_add {
			if ret := C.OCBadd(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), C.BFLDID(bfldid),
				c_arr, c_len, BFLD_CARRAY); ret != SUCCEED {
				return u.Buf.Ctx.NewUBFError()
			}
		} else {
			if ret := C.OCBchg(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), C.BFLDID(bfldid),
				C.BFLDOCC(occ), c_arr, c_len, BFLD_CARRAY); ret != SUCCEED {
				return u.Buf.Ctx.NewUBFError()
			}
		}
		/*
				- Currently not supported!
			case fmt.Stringer:
				str := ival.(fmt.Stringer).String()
				c_val := C.CString(str)
				defer C.free(unsafe.Pointer(c_val))
				if ret := C.CBchg((*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr), C.BFLDID(bfldid),
					C.BFLDOCC(occ), c_val, 0, BFLD_STRING); ret != SUCCEED {
					return NewUBFError()
				}
		*/
	default:
		/* TODO: Possibly we could take stuff from println to get string val... */
		return NewCustomUBFError(BEINVAL, "Cannot determine field type")
	}

	u.Buf.nop()
	return nil
}

//Compile boolean expression
//TODO: might want auto finalizer with Btreefree!
//@param	expr Expression string
//@return Expression tree (ptr or nil on error), UBF error
func (ac *ATMICtx) BBoolCo(expr string) (*ExprTree, UBFError) {

	c_str := C.CString(expr)

	defer C.free(unsafe.Pointer(c_str))

	c_ptr := C.OBboolco(&ac.c_ctx, c_str)

	if nil == c_ptr {
		return nil, ac.NewUBFError()
	}

	var tree ExprTree

	tree.c_ptr = c_ptr

	//Free up the data once GCed
	//Well we might have issue here, the ATMI Context might be already
	//Deallocated, thus we need to have temp context free op.
	runtime.SetFinalizer(&tree, btreeFree)

	ac.nop() //keep context until the end of the func, and only then allow gc
	return &tree, nil
}

//Free the expression buffer
func (ac *ATMICtx) BTreeFree(tree *ExprTree) {

	//Unset the finalizer
	C.OBtreefree(&ac.c_ctx, tree.c_ptr)
	tree.c_ptr = nil

	ac.nop() //keep context until the end of the func, and only then allow gc
	tree.nop()
}

//Internal version (uses temp context)
func btreeFree(tree *ExprTree) {

	if nil != tree.c_ptr {
		C.go_Btreefree(tree.c_ptr)
		tree.c_ptr = nil
	}

	tree.nop()
}

//Test the expresion tree to current UBF buffer
//@param tree	compiled expression tree
//@return true (buffer matched expression) or false (buffer not matched expression)
func (u *TypedUBF) BBoolEv(tree *ExprTree) bool {

	c_ret := C.OBboolev(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), tree.c_ptr)

	if 1 == c_ret {
		return true
	}

	tree.nop()
	u.Buf.nop()
	return false
}

//Quick eval of the expression (compiles & frees the handler automatically)
//@param expr Expression tree
//@return result: true or false, UBF error
func (u *TypedUBF) BQBoolEv(expr string) (bool, UBFError) {

	h_exp, err := u.Buf.Ctx.BBoolCo(expr)

	if err == nil {
		defer u.Buf.Ctx.BTreeFree(h_exp)
	} else {
		return false, err
	}

	ret := u.BBoolEv(h_exp)

	return ret, nil
}

//Evalute expression value in float64 format
//@param tree	compiled expression tree
//@return expression value
func (u *TypedUBF) BFloatEv(tree *ExprTree) float64 {

	c_ret := C.OBfloatev(&u.Buf.Ctx.c_ctx,
		(*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), tree.c_ptr)

	tree.nop()
	u.Buf.nop()

	return float64(c_ret)
}

//Generate Field ID
//@param fldtype Field type (see BFLD_SHORT cost list)
//@param bfldid field number
//@return field id or 0 if error, UBF error
func (ac *ATMICtx) BMkFldId(fldtype int, bfldid int) (int, UBFError) {

	c_ret := C.OBmkfldid(&ac.c_ctx, C.int(fldtype), C.BFLDID(bfldid))
	if BBADFLDID == c_ret {
		return BBADFLDID, ac.NewUBFError()
	}

	ac.nop() //keep context until the end of the func, and only then allow gc
	return int(c_ret), nil
}

//Get the number of field occurrances in buffer
//@param bfldid	Field ID
//@return count (or -1 on error), UBF error
func (u *TypedUBF) BOccur(bfldid int) (int, UBFError) {

	c_ret := C.OBoccur(&u.Buf.Ctx.c_ctx,
		(*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), C.BFLDID(bfldid))

	if FAIL == c_ret {
		return FAIL, u.Buf.Ctx.NewUBFError()
	}

	u.Buf.nop()
	return int(c_ret), nil
}

//Get the number of bytes used in UBF buffer
//@return number of byptes used, UBF error
func (u *TypedUBF) BUsed() (int64, UBFError) {

	c_ret := C.OBused(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)))

	if FAIL == c_ret {
		return FAIL, u.Buf.Ctx.NewUBFError()
	}

	u.Buf.nop()
	return int64(c_ret), nil
}

//Get the number of free bytes of UBF buffer
//@return buffer free bytes, UBF error
func (u *TypedUBF) BUnused() (int64, UBFError) {

	c_ret := C.OBunused(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)))

	if FAIL == c_ret {
		return FAIL, u.Buf.Ctx.NewUBFError()
	}

	u.Buf.nop()
	return int64(c_ret), nil
}

//Get the total buffer size
//@return bufer size, UBF error
func (u *TypedUBF) BSizeof() (int64, UBFError) {

	c_ret := C.OBsizeof(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)))

	if FAIL == c_ret {
		return FAIL, u.Buf.Ctx.NewUBFError()
	}

	u.Buf.nop()
	return int64(c_ret), nil
}

//Return the field type
//@param bfldid field id
//@return field type
func (ac *ATMICtx) BFldType(bfldid int) int {

	c_ret := C.OBfldtype(&ac.c_ctx, C.BFLDID(bfldid))

	ac.nop()
	return int(c_ret)
}

//Return field number
//@param bfldid field id
//@return field number
func (ac *ATMICtx) BFldNo(bfldid int) int {

	c_ret := C.OBfldno(&ac.c_ctx, C.BFLDID(bfldid))

	ac.nop()
	return int(c_ret)
}

//Delete field (all occurrances) from buffer
//@param bfldid field ID
//@return UBF error
func (u *TypedUBF) BDelAll(bfldid int) UBFError {

	if ret := C.OBdelall(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)),
		C.BFLDID(bfldid)); SUCCEED != ret {
		return u.Buf.Ctx.NewUBFError()
	}

	u.Buf.nop()
	return nil
}

//Delete listed fields from UBF buffer
//@param fldlist list of fields (array)
//@return UBF error
func (u *TypedUBF) BDelete(fldlist []int) UBFError {

	c_fldidsize := int(C.c_sizeof_BFLDID())

	c_val := C.malloc(C.size_t(c_fldidsize * (len(fldlist) + 1)))

	if nil == c_val {
		return NewCustomUBFError(BEUNIX, "Cannot alloc memory")
	}

	defer C.free(c_val)
	var i int

	for i = 0; i < len(fldlist); i++ {
		*(*C.BFLDID)(unsafe.Pointer(uintptr(c_val) + uintptr(i*c_fldidsize))) =
			C.BFLDID(fldlist[i])
	}

	//Set last field to BBADFLDID
	*(*C.BFLDID)(unsafe.Pointer(uintptr(c_val) + uintptr(i*c_fldidsize))) = C.BFLDID(BBADFLDID)

	if ret := C.OBdelete(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)),
		(*C.BFLDID)(unsafe.Pointer(c_val))); ret != SUCCEED {
		return u.Buf.Ctx.NewUBFError()
	}

	u.Buf.nop()
	return nil
}

//Return type descriptor of the field - string format.
//possible values: short, long, char, float, double, string, carray
//@param bfldid field ID
//@return field type, UBF error
func (u *TypedUBF) BType(bfldid int) (string, UBFError) {

	ret := C.OBtype(&u.Buf.Ctx.c_ctx, C.BFLDID(bfldid))

	if nil == ret {
		return "", u.Buf.Ctx.NewUBFError()
	}

	u.Buf.nop()
	return C.GoString(ret), nil
}

//Update dest buffer with source buffer data
//@param dest 	dest buffer
//@param src		source buffer
//@return UBF error
func (ac *ATMICtx) BUpdate(dest *TypedUBF, src *TypedUBF) UBFError {

	if ret := C.OBupdate(&ac.c_ctx, (*C.UBFH)(unsafe.Pointer(dest.Buf.C_ptr)),
		(*C.UBFH)(unsafe.Pointer(src.Buf.C_ptr))); ret != SUCCEED {
		return ac.NewUBFError()
	}

	ac.nop()
	dest.Buf.nop()
	src.Buf.nop()

	return nil
}

//Contact the buffers
//@param dest 	dest buffer
//@param src		source buffer
//@return UBF error
func (ac *ATMICtx) BConcat(dest *TypedUBF, src *TypedUBF) UBFError {

	if ret := C.OBconcat(&ac.c_ctx, (*C.UBFH)(unsafe.Pointer(dest.Buf.C_ptr)),
		(*C.UBFH)(unsafe.Pointer(src.Buf.C_ptr))); ret != SUCCEED {
		return ac.NewUBFError()
	}

	ac.nop()
	dest.Buf.nop()
	src.Buf.nop()

	return nil
}

//Print the buffer to stdout
//@return UBF error
func (u *TypedUBF) BPrint() UBFError {

	if ret := C.OBprint(&u.Buf.Ctx.c_ctx,
		(*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr))); ret != SUCCEED {
		return u.Buf.Ctx.NewUBFError()
	}

	u.Buf.nop()

	return nil
}

//Print the buffer to stdout
//@return UBF error
func (u *TypedUBF) TpLogPrintUBF(lev int, title string) {

	c_title := C.CString(title)

	C.Otplogprintubf(&u.Buf.Ctx.c_ctx, C.int(lev), c_title,
		(*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)))

	C.free(unsafe.Pointer(c_title))
	u.Buf.nop()

	return
}

/*
Replaced with callback version with EX 6+
//Alternative for Bfprint. Will return the output in string variable
//So that caller can do anything it wants with the string output
//@return Printed buffer, UBF error
func (u *TypedUBF) BSprint() (string, UBFError) {

	c_val := C.calloc(C.size_t(ATMIMsgSizeMax()), 10)
	c_len := C.size_t(C.size_t(ATMIMsgSizeMax()) * 10)

	if nil == c_val {
		return "", NewCustomUBFError(BEUNIX, "Cannot alloc memory")
	}

	c_mode := C.CString("w")

	defer C.free(unsafe.Pointer(c_mode))
	defer C.free(c_val)

	f := C.fmemopen(c_val, c_len, c_mode)

	if nil == f {
		return "", NewCustomUBFError(BEUNIX, "Cannot open FILE handle")
	}

	defer C.fclose(f)

	if ret := C.OBfprint(&u.Buf.Ctx.c_ctx,
		(*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), f); ret != SUCCEED {
		return "", u.Buf.Ctx.NewUBFError()
	}

	return C.GoString((*C.char)(c_val)), nil
}
*/

//Print UBF buffer to string. The output string buffer at C side is composed
//as UBF buffer size of multiplied by MAXTIDENT (currently 30). The total size
//is used for purpuse so that Go developer can used extended buffer size in case
//if there is no free space (returned error BEUNIX)
//@returns BPrint format string or "" in case of error. Second argument is
//	UBF error set in case of error, else it is nil
func (u *TypedUBF) BSprint() (string, UBFError) {

	c_str := C.BPrintStrC(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)))

	if nil != c_str {

		str := C.GoString(c_str)
		C.free(unsafe.Pointer(c_str))

		return str, nil
	}

	u.Buf.nop()

	return "", NewCustomUBFError(BEUNIX, "Failed to print UBF buffer to string, "+
		"either insufficient memory or other error. See UBF logs.")
}

//Read the bufer content from string
//@param s String buffer representation
//@return UBF error
func (u *TypedUBF) BExtRead(s string) UBFError {

	c_val := C.CString(s)
	defer C.free(unsafe.Pointer(c_val))

	c_len := C.strlen(c_val)

	c_mode := C.CString("r")
	defer C.free(unsafe.Pointer(c_mode))

	f := C.fmemopen(unsafe.Pointer(c_val), c_len, c_mode)

	if nil == f {
		return NewCustomUBFError(BEUNIX, "Cannot open FILE handle")
	}

	defer C.fclose(f)

	if ret := C.OBextread(&u.Buf.Ctx.c_ctx,
		(*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), f); ret != SUCCEED {
		return u.Buf.Ctx.NewUBFError()
	}

	u.Buf.nop()
	return nil
}

//Print the expression tree
//@param tree 	Compiled expression tree
//@return printed expresion string, ubf error
func (ac *ATMICtx) BBoolPr(tree *ExprTree) (string, UBFError) {

	c_val := C.calloc(C.size_t(ATMIMsgSizeMax()), 10)
	c_len := C.size_t(C.size_t(ATMIMsgSizeMax()) * 10)

	if nil == c_val {
		return "", NewCustomUBFError(BEUNIX, "Cannot alloc memory")
	}

	c_mode := C.CString("w")

	defer C.free(unsafe.Pointer(c_mode))
	defer C.free(c_val)

	f := C.fmemopen(c_val, c_len, c_mode)

	if nil == f {
		return "", NewCustomUBFError(BEUNIX, "Cannot open FILE handle")
	}

	defer C.fclose(f)

	C.OBboolpr(&ac.c_ctx, tree.c_ptr, f)

	ac.nop()
	tree.nop()

	return C.GoString((*C.char)(c_val)), nil
}

//Serialize the UBF buffer
//@return serialized bytes, UBF error
func (u *TypedUBF) BWrite() ([]byte, UBFError) {

	c_val := C.calloc(C.size_t(ATMIMsgSizeMax()), 1)
	c_len := C.size_t(C.size_t(ATMIMsgSizeMax()))

	if nil == c_val {
		return nil, NewCustomUBFError(BEUNIX, "Cannot alloc memory")
	}

	c_mode := C.CString("wb")

	defer C.free(unsafe.Pointer(c_mode))
	defer C.free(c_val)

	f := C.fmemopen(c_val, c_len, c_mode)

	if nil == f {
		return nil, NewCustomUBFError(BEUNIX, "Cannot open FILE handle")
	}

	defer C.fclose(f)

	if SUCCEED != C.OBwrite(&u.Buf.Ctx.c_ctx,
		(*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), f) {
		return nil, u.Buf.Ctx.NewUBFError()
	}

	size := C.ftell(f)

	var array = make([]byte, int(size))

	for i := 0; i < int(size); i++ {
		array[i] = byte(*(*C.char)(unsafe.Pointer(uintptr(c_val) + uintptr(i))))
	}

	u.Buf.nop()
	return array, nil
}

//Serialize the UBF buffer
//@return serialized bytes, UBF error
func (u *TypedUBF) BRead(dump []byte) UBFError {

	c_val := C.malloc(C.size_t(len(dump)))
	c_len := C.size_t(len(dump))

	if nil == c_val {
		return NewCustomUBFError(BEUNIX, "Cannot alloc memory")
	}

	//Copy stuff to C memory
	for i := 0; i < len(dump); i++ {
		*(*C.char)(unsafe.Pointer(uintptr(c_val) + uintptr(i))) = C.char(dump[i])
	}

	c_mode := C.CString("rb")

	defer C.free(unsafe.Pointer(c_mode))
	defer C.free(c_val)

	f := C.fmemopen(c_val, c_len, c_mode)

	if nil == f {
		return NewCustomUBFError(BEUNIX, "Cannot open FILE handle")
	}

	defer C.fclose(f)

	if SUCCEED != C.OBread(&u.Buf.Ctx.c_ctx,
		(*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)), f) {
		return u.Buf.Ctx.NewUBFError()
	}

	u.Buf.nop()
	return nil
}

//Test C buffer for UBF format
//@return TRUE - buffer is UBF, FALSE - not UBF
func (u *TypedUBF) BIsUBF() bool {

	c_ret := C.OBisubf(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)))
	if 1 == c_ret {
		return true
	}

	u.Buf.nop()
	return false
}

//export go_expr_callback_proxy
func go_expr_callback_proxy(buf *C.char, funcname *C.char) C.long {
	var u TypedUBF
	u.Buf.C_ptr = buf

	ret := exprfuncmap[C.GoString(funcname)](&u, C.GoString(funcname))

	//Map the call
	return C.long(ret)
}

//Set custom callback function for UBF buffer expression evaluator
//@param funcname Name of the function to be used in expression
//@param f callback to function
//@return UBF error
func (ac *ATMICtx) BBoolSetCBF(funcname string, f UBFExprFunc) UBFError {

	if nil == f || "" == funcname {
		return NewCustomUBFError(BEINVAL, "func nil or func name empty!")
	}

	c_funcname := C.CString(funcname)

	defer C.free(unsafe.Pointer(c_funcname))

	if SUCCEED != C.c_proxy_Bboolsetcbf(&ac.c_ctx, c_funcname) {
		return ac.NewUBFError()
	} else {
		exprfuncmap[funcname] = f
	}

	ac.nop()
	return nil
}

//Allocate the new UBF buffer
//NOTE: realloc or other ATMI ops you can do with TypedUBF.Buf
//@param size - buffer size
//@return Typed UBF, ATMI error
func (ac *ATMICtx) NewUBF(size int64) (*TypedUBF, ATMIError) {

	var buf TypedUBF

	if ptr, err := ac.TpAlloc("UBF", "", size); nil != err {
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
func (u *TypedUBF) TpJSONToUBF(buffer string) UBFError {

	size := int64(len(buffer))
	sizeof, _ := u.BSizeof()
	unused, _ := u.BUnused()
	alloc := size - unused

	c_buffer := C.CString(buffer)

	defer C.free(unsafe.Pointer(c_buffer))

	u.Buf.Ctx.ndrxLog(LOG_INFO, "Data size: %d, UBF sizeof: %d, "+
		"unused: %d, about to alloc (if >0) %d",
		size, sizeof, unused, alloc)

	if alloc > 0 {
		if err := u.TpRealloc(sizeof + alloc); nil != err {
			return NewCustomUBFError(BMALLOC, err.Message())
		}
	}

	if ret := C.Otpjsontoubf(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)),
		c_buffer); ret != 0 {
		return u.Buf.Ctx.NewUBFError()
	}

	u.Buf.nop()

	return nil
}

//Convert given UBF buffer to JSON block, see tpubftojson(3) C call
//Output string is automatically allocated
//@return JSON string (if converted ok), ATMIError in case of failure. More detailed
//infos in case of error is found in 'ubf' and 'ndrx' facility logs.
func (u *TypedUBF) TpUBFToJSON() (string, ATMIError) {

	used, _ := u.BUsed()

	ret_size := used * 10

	u.Buf.Ctx.ndrxLog(LOG_INFO, "TpUBFToJSON: used %d allocating %d", used, ret_size)

	c_buffer := C.malloc(C.size_t(ret_size))

	if nil == c_buffer {
		return "", NewCustomUBFError(BEUNIX, "Cannot alloc memory")
	}

	defer C.free(c_buffer)

	if ret := C.Otpubftojson(&u.Buf.Ctx.c_ctx, (*C.UBFH)(unsafe.Pointer(u.Buf.C_ptr)),
		(*C.char)(unsafe.Pointer(c_buffer)), C.int(ret_size)); ret != 0 {
		return "", u.Buf.Ctx.NewATMIError()
	}

	u.Buf.nop()

	return C.GoString((*C.char)(c_buffer)), nil

}

///////////////////////////////////////////////////////////////////////////////////
// Wrappers for memory management
///////////////////////////////////////////////////////////////////////////////////

func (u *TypedUBF) TpRealloc(size int64) ATMIError {

	return u.Buf.TpRealloc(size)
}

/* vim: set ts=4 sw=4 et smartindent: */
