package atmi

/*
#cgo pkg-config: atmisrvinteg

#include <xatmi.h>
#include <string.h>
#include <stdlib.h>
#include <ubf.h>
#include <oubf.h>
#include <oatmi.h>
#include <sys_unix.h>


//Get the UBF Error code
static int WrapBerror(TPCONTEXT_T *p_ctx) {
	return OBerror(p_ctx);
}

//Cast the data type
static UBFH *GetU(char *data) {
	return (UBFH *)data;
}


//Get Char Ptr to void pointer
static char *GetCharPtr(void *ptr) {
	return (char *)ptr;
}

#define ATMI_MSG_MAX_SIZE	65536

//Get the value with buffer allocation
static char * c_Bget_str (UBFH * p_ub, BFLDID bfldid, BFLDOCC occ,
					BFLDLEN *len, int *err_code)
{
	char *ret = malloc(ATMI_MSG_MAX_SIZE);

	*len = ATMI_MSG_MAX_SIZE;
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

*/
import "C"
import "fmt"

import "unsafe"
import "runtime"

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
	BERFU2    = 16
	BERFU3    = 17
	BERFU4    = 18
	BERFU5    = 19
	BERFU6    = 20
	BERFU7    = 21
	BERFU8    = 22
	BMAXVAL   = 22 /* max error */
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
	c_ptr *C.char
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

//Get the field len
//@param fldid	Field ID
//@param occ 	Field occurance
//@return 	FIeld len, UBF error
func (u *TypedUBF) BLen(bfldid int, occ int) (int, UBFError) {
	ret := C.OBlen(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), C.BFLDID(bfldid), C.BFLDOCC(occ))
	if FAIL == ret {
		return FAIL, u.Buf.Ctx.NewUBFError()
	}
	return int(ret), nil
}

//Delete the field from buffer
//@param fldid	Field ID
//@param occ 	Field occurance
//@return 	UBF error
func (u *TypedUBF) BDel(bfldid int, occ int) UBFError {
	ret := C.OBdel(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), C.BFLDID(bfldid), C.BFLDOCC(occ))
	if FAIL == ret {
		return u.Buf.Ctx.NewUBFError()
	}
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

	if ret := C.OBproj(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), (*C.BFLDID)(unsafe.Pointer(c_val))); ret != SUCCEED {
		return u.Buf.Ctx.NewUBFError()
	}

	return nil
}

//Make a project copy of the fields (leave only those in array)
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

	if ret := C.OBprojcpy(&ac.c_ctx, C.GetU(dest.Buf.C_ptr), C.GetU(src.Buf.C_ptr),
		(*C.BFLDID)(unsafe.Pointer(c_val))); ret != SUCCEED {
		return ac.NewUBFError()
	}

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

	return C.GoString(ret), nil
}

//Check for field presence in buffer
//@param fldid	Field ID
//@param occ 	Field occurance
//@return 	true/false present/not present
func (u *TypedUBF) BPres(bfldid int, occ int) bool {
	ret := C.OBpres(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr),
		C.BFLDID(bfldid), C.BFLDOCC(occ))
	if 1 == ret {
		return true
	}
	return false
}

//Copy buffer
//@param dest Destination UBF buffer
//@param src		Source UBF buffer
//@return UBF error
func (ac *ATMICtx) BCpy(dest *TypedUBF, src *TypedUBF) UBFError {
	if ret := C.OBcpy(&ac.c_ctx, C.GetU(dest.Buf.C_ptr),
		C.GetU(src.Buf.C_ptr)); SUCCEED != ret {
		return ac.NewUBFError()
	}
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

	if ret := C.OBnext(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr),
		&fldid, &occ, nil, nil); 1 != ret {
		return FAIL, FAIL, u.Buf.Ctx.NewUBFError()
	}

	return int(fldid), int(occ), nil
}

//Initialize/re-initialize UBF buffer
//@param u UBF buffer
//@param ulen	lenght of the buffer
//@return UBF error
func (ac *ATMICtx) BInit(u *TypedUBF, ulen int) UBFError {
	if ret := C.OBinit(&ac.c_ctx, C.GetU(u.Buf.C_ptr),
		C.BFLDLEN(ulen)); SUCCEED != ret {
		return ac.NewUBFError()
	}
	return nil
}

//Allocate the UBF buffer
//@param size	Buffer size in bytes
//@return UBF Handler, ATMI Error
func (ac *ATMICtx) UBFAlloc(size int64) (TypedUBF, ATMIError) {
	var err ATMIError
	var buf TypedUBF
	buf.Buf, err = ac.TpAlloc("UBF", "", size)
	return buf, err
}

//Get the UBF Handler
func (ac *ATMICtx) CastToUBF(abuf *ATMIBuf) (TypedUBF, ATMIError) {
	var buf TypedUBF

	//TODO: Check the buffer type!
	buf.Buf = abuf

	return buf, nil
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
		if ret := C.OBget(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), C.BFLDID(bfldid),
			C.BFLDOCC(occ), C.GetCharPtr(unsafe.Pointer(&c_val)), nil); ret != SUCCEED {
			return nil, u.Buf.Ctx.NewUBFError()
		}
		return int16(c_val), nil
	case BFLD_LONG:
		var c_val C.short
		if ret := C.OBget(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), C.BFLDID(bfldid),
			C.BFLDOCC(occ), C.GetCharPtr(unsafe.Pointer(&c_val)), nil); ret != SUCCEED {
			return nil, u.Buf.Ctx.NewUBFError()
		}
		return int64(c_val), nil
	case BFLD_CHAR: /* This is single byte... */
		var c_val C.char
		if ret := C.OBget(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), C.BFLDID(bfldid),
			C.BFLDOCC(occ), C.GetCharPtr(unsafe.Pointer(&c_val)), nil); ret != SUCCEED {
			return nil, u.Buf.Ctx.NewUBFError()
		}
		return byte(c_val), nil
	case BFLD_FLOAT:
		var c_val C.float
		if ret := C.OBget(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), C.BFLDID(bfldid),
			C.BFLDOCC(occ), C.GetCharPtr(unsafe.Pointer(&c_val)), nil); ret != SUCCEED {
			return nil, u.Buf.Ctx.NewUBFError()
		}
		return float32(c_val), nil
	case BFLD_DOUBLE:
		var c_val C.double
		if ret := C.OBget(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), C.BFLDID(bfldid),
			C.BFLDOCC(occ), C.GetCharPtr(unsafe.Pointer(&c_val)), nil); ret != SUCCEED {
			return nil, u.Buf.Ctx.NewUBFError()
		}
		return float64(c_val), nil
	case BFLD_STRING:
		var c_len C.BFLDLEN
		c_val := C.malloc(ATMI_MSG_MAX_SIZE)
		c_len = ATMI_MSG_MAX_SIZE

		if nil == c_val {
			return nil, NewCustomUBFError(BEUNIX, "Cannot alloc memory")
		}

		defer C.free(c_val)

		if ret := C.OBget(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), C.BFLDID(bfldid),
			C.BFLDOCC(occ), (*C.char)(c_val), &c_len); ret != SUCCEED {
			return nil, u.Buf.Ctx.NewUBFError()
		}

		return C.GoString((*C.char)(c_val)), nil

	case BFLD_CARRAY:
		var c_len C.BFLDLEN
		c_val := C.malloc(ATMI_MSG_MAX_SIZE)
		c_len = ATMI_MSG_MAX_SIZE

		if nil == c_val {
			return nil, NewCustomUBFError(BEUNIX, "Cannot alloc memory")
		}

		defer C.free(c_val)

		if ret := C.OBget(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), C.BFLDID(bfldid),
			C.BFLDOCC(occ), (*C.char)(c_val), &c_len); ret != SUCCEED {
			return nil, u.Buf.Ctx.NewUBFError()
		}

		g_val := make([]byte, c_len)

		for i := 0; i < int(c_len); i++ {
			g_val[i] = byte(*(*C.char)(unsafe.Pointer(uintptr(c_val) + uintptr(i))))
		}

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
	if ret := C.OCBget(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), C.BFLDID(bfldid),
		C.BFLDOCC(occ), C.GetCharPtr(unsafe.Pointer(&c_val)), nil, BFLD_SHORT); ret != SUCCEED {
		return 0, u.Buf.Ctx.NewUBFError()
	}
	return int16(c_val), nil
}

//Return int64 value from buffer
//@param bfldid 	Field ID
//@param occ	Occurrance
//@return int64 val,	 UBF error
func (u *TypedUBF) BGetInt64(bfldid int, occ int) (int64, UBFError) {
	var c_val C.long
	if ret := C.OCBget(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), C.BFLDID(bfldid),
		C.BFLDOCC(occ), C.GetCharPtr(unsafe.Pointer(&c_val)), nil, BFLD_LONG); ret != SUCCEED {
		return 0, u.Buf.Ctx.NewUBFError()
	}
	return int64(c_val), nil
}

//Return byte (c char) value from buffer
//@param bfldid 	Field ID
//@param occ	Occurrance
//@return byte val, UBF error
func (u *TypedUBF) BGetByte(bfldid int, occ int) (byte, UBFError) {
	var c_val C.char
	if ret := C.OCBget(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), C.BFLDID(bfldid),
		C.BFLDOCC(occ), C.GetCharPtr(unsafe.Pointer(&c_val)), nil, BFLD_CHAR); ret != SUCCEED {
		return 0, u.Buf.Ctx.NewUBFError()
	}
	return byte(c_val), nil
}

//Get float value
//@param bfldid 	Field ID
//@param occ	Occurrance
//@return float, UBF error
func (u *TypedUBF) BGetFloat32(bfldid int, occ int) (float32, UBFError) {
	var c_val C.float
	if ret := C.OCBget(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), C.BFLDID(bfldid),
		C.BFLDOCC(occ), C.GetCharPtr(unsafe.Pointer(&c_val)), nil, BFLD_FLOAT); ret != SUCCEED {
		return 0, u.Buf.Ctx.NewUBFError()
	}
	return float32(c_val), nil
}

//Get double value
//@param bfldid 	Field ID
//@param occ	Occurrance
//@return double, UBF error
func (u *TypedUBF) BGetFloat64(bfldid int, occ int) (float64, UBFError) {
	var c_val C.double
	if ret := C.OCBget(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), C.BFLDID(bfldid),
		C.BFLDOCC(occ), C.GetCharPtr(unsafe.Pointer(&c_val)), nil, BFLD_DOUBLE); ret != SUCCEED {
		return 0, u.Buf.Ctx.NewUBFError()
	}
	return float64(c_val), nil
}

//Get string value
//@param bfldid 	Field ID
//@param occ	Occurrance
//@return string val, UBF error
func (u *TypedUBF) BGetString(bfldid int, occ int) (string, UBFError) {
	var c_len C.BFLDLEN
	c_val := C.malloc(ATMI_MSG_MAX_SIZE)
	c_len = ATMI_MSG_MAX_SIZE

	if nil == c_val {
		return "", NewCustomUBFError(BEUNIX, "Cannot alloc memory")
	}

	defer C.free(c_val)

	if ret := C.OCBget(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), C.BFLDID(bfldid),
		C.BFLDOCC(occ), (*C.char)(c_val), &c_len, BFLD_STRING); ret != SUCCEED {
		return "", u.Buf.Ctx.NewUBFError()
	}

	return C.GoString((*C.char)(c_val)), nil
}

//Get string value
//@param bfldid 	Field ID
//@param occ	Occurrance
//@return string val, UBF error
func (u *TypedUBF) BGetByteArr(bfldid int, occ int) ([]byte, UBFError) {
	var c_len C.BFLDLEN
	c_val := C.malloc(ATMI_MSG_MAX_SIZE)
	c_len = ATMI_MSG_MAX_SIZE

	if nil == c_val {
		return nil, NewCustomUBFError(BEUNIX, "Cannot alloc memory")
	}

	defer C.free(c_val)

	if ret := C.OCBget(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), C.BFLDID(bfldid),
		C.BFLDOCC(occ), (*C.char)(c_val), &c_len, BFLD_CARRAY); ret != SUCCEED {
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
func (u *TypedUBF) BChg(bfldid int, occ int, ival interface{}) UBFError {
	return u.BChgCombined(bfldid, occ, ival, false)
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
			if ret := C.OCBadd(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), C.BFLDID(bfldid),
				C.GetCharPtr(unsafe.Pointer(&c_val)), 0, BFLD_LONG); ret != SUCCEED {
				return u.Buf.Ctx.NewUBFError()
			}
		} else {
			if ret := C.OCBchg(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), C.BFLDID(bfldid),
				C.BFLDOCC(occ), C.GetCharPtr(unsafe.Pointer(&c_val)), 0, BFLD_LONG); ret != SUCCEED {
				return u.Buf.Ctx.NewUBFError()
			}
		}
	case float32:
		fval := ival.(float32)
		c_val := C.float(fval)
		if do_add {
			if ret := C.OCBadd(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), C.BFLDID(bfldid),
				C.GetCharPtr(unsafe.Pointer(&c_val)), 0, BFLD_FLOAT); ret != SUCCEED {
				return u.Buf.Ctx.NewUBFError()
			}
		} else {
			if ret := C.OCBchg(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), C.BFLDID(bfldid),
				C.BFLDOCC(occ), C.GetCharPtr(unsafe.Pointer(&c_val)), 0, BFLD_FLOAT); ret != SUCCEED {
				return u.Buf.Ctx.NewUBFError()
			}
		}
	case float64:
		dval := ival.(float64)
		c_val := C.double(dval)
		if do_add {
			if ret := C.OCBadd(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), C.BFLDID(bfldid),
				C.GetCharPtr(unsafe.Pointer(&c_val)), 0, BFLD_DOUBLE); ret != SUCCEED {
				return u.Buf.Ctx.NewUBFError()
			}
		} else {
			if ret := C.OCBchg(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), C.BFLDID(bfldid),
				C.BFLDOCC(occ), C.GetCharPtr(unsafe.Pointer(&c_val)), 0, BFLD_DOUBLE); ret != SUCCEED {
				return u.Buf.Ctx.NewUBFError()
			}
		}
	case string:
		str := ival.(string)
		c_val := C.CString(str)
		defer C.free(unsafe.Pointer(c_val))
		if do_add {
			if ret := C.OCBadd(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), C.BFLDID(bfldid),
				c_val, 0, BFLD_STRING); ret != SUCCEED {
				return u.Buf.Ctx.NewUBFError()
			}
		} else {
			if ret := C.OCBchg(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), C.BFLDID(bfldid),
				C.BFLDOCC(occ), c_val, 0, BFLD_STRING); ret != SUCCEED {
				return u.Buf.Ctx.NewUBFError()
			}
		}
	case []byte:
		arr := ival.([]byte)
		c_len := C.BFLDLEN(len(arr))
		c_arr := (*C.char)(unsafe.Pointer(&arr[0]))

		if do_add {
			if ret := C.OCBadd(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), C.BFLDID(bfldid),
				c_arr, c_len, BFLD_CARRAY); ret != SUCCEED {
				return u.Buf.Ctx.NewUBFError()
			}
		} else {
			if ret := C.OCBchg(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), C.BFLDID(bfldid),
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
				if ret := C.CBchg(C.GetU(u.Buf.C_ptr), C.BFLDID(bfldid),
					C.BFLDOCC(occ), c_val, 0, BFLD_STRING); ret != SUCCEED {
					return NewUBFError()
				}
		*/
	default:
		/* TODO: Possibly we could take stuff from println to get string val... */
		return NewCustomUBFError(BEINVAL, "Cannot determine field type")
	}

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

	return &tree, nil
}

//Free the expression buffer
func (ac *ATMICtx) BTreeFree(tree *ExprTree) {
	//Unset the finalizer
	C.OBtreefree(&ac.c_ctx, tree.c_ptr)
	tree.c_ptr = nil
}

//Internal version (uses temp context)
func btreeFree(tree *ExprTree) {
	if nil!=tree.c_ptr {
		C.go_Btreefree(tree.c_ptr)
		tree.c_ptr = nil
	}
}

//Test the expresion tree to current UBF buffer
//@param tree	compiled expression tree
//@return true (buffer matched expression) or false (buffer not matched expression)
func (u *TypedUBF) BBoolEv(tree *ExprTree) bool {
	c_ret := C.OBboolev(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), tree.c_ptr)

	if 1 == c_ret {
		return true
	}

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
	c_ret := C.OBfloatev(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), tree.c_ptr)

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

	return int(c_ret), nil
}

//Get the number of field occurrances in buffer
//@param bfldid	Field ID
//@return count (or -1 on error), UBF error
func (u *TypedUBF) BOccur(bfldid int) (int, UBFError) {
	c_ret := C.OBoccur(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), C.BFLDID(bfldid))

	if FAIL == c_ret {
		return FAIL, u.Buf.Ctx.NewUBFError()
	}

	return int(c_ret), nil
}

//Get the number of bytes used in UBF buffer
//@return number of byptes used, UBF error
func (u *TypedUBF) BUsed() (int64, UBFError) {
	c_ret := C.OBused(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr))

	if FAIL == c_ret {
		return FAIL, u.Buf.Ctx.NewUBFError()
	}

	return int64(c_ret), nil
}

//Get the number of free bytes of UBF buffer
//@return buffer free bytes, UBF error
func (u *TypedUBF) BUnused() (int64, UBFError) {
	c_ret := C.OBunused(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr))

	if FAIL == c_ret {
		return FAIL, u.Buf.Ctx.NewUBFError()
	}

	return int64(c_ret), nil
}

//Get the total buffer size
//@return bufer size, UBF error
func (u *TypedUBF) BSizeof() (int64, UBFError) {
	c_ret := C.OBsizeof(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr))

	if FAIL == c_ret {
		return FAIL, u.Buf.Ctx.NewUBFError()
	}

	return int64(c_ret), nil
}

//Return the field type
//@param bfldid field id
//@return field type
func (ac *ATMICtx) BFldType(bfldid int) int {
	c_ret := C.OBfldtype(&ac.c_ctx, C.BFLDID(bfldid))
	return int(c_ret)
}

//Return field number
//@param bfldid field id
//@return field number
func (ac *ATMICtx) BFldNo(bfldid int) int {
	c_ret := C.OBfldno(&ac.c_ctx, C.BFLDID(bfldid))
	return int(c_ret)
}

//Delete field (all occurrances) from buffer
//@param bfldid field ID
//@return UBF error
func (u *TypedUBF) BDelAll(bfldid int) UBFError {
	if ret := C.OBdelall(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), C.BFLDID(bfldid)); SUCCEED != ret {
		return u.Buf.Ctx.NewUBFError()
	}
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

	if ret := C.OBdelete(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), (*C.BFLDID)(unsafe.Pointer(c_val))); ret != SUCCEED {
		return u.Buf.Ctx.NewUBFError()
	}

	return nil
}

//Return field name in string
//@param bfldid field ID
//@return field type, UBF error
func (u *TypedUBF) BType(bfldid int) (string, UBFError) {
	ret := C.OBtype(&u.Buf.Ctx.c_ctx, C.BFLDID(bfldid))

	if nil == ret {
		return "", u.Buf.Ctx.NewUBFError()
	}

	return C.GoString(ret), nil
}

//Update dest buffer with source buffer data
//@param dest 	dest buffer
//@param src		source buffer
//@return UBF error
func (ac *ATMICtx) BUpdate(dest *TypedUBF, src *TypedUBF) UBFError {

	if ret := C.OBupdate(&ac.c_ctx, C.GetU(dest.Buf.C_ptr), C.GetU(src.Buf.C_ptr)); ret != SUCCEED {
		return ac.NewUBFError()
	}
	return nil
}

//Contact the buffers
//@param dest 	dest buffer
//@param src		source buffer
//@return UBF error
func (ac *ATMICtx) BConcat(dest *TypedUBF, src *TypedUBF) UBFError {
	if ret := C.OBconcat(&ac.c_ctx, C.GetU(dest.Buf.C_ptr), C.GetU(src.Buf.C_ptr)); ret != SUCCEED {
		return ac.NewUBFError()
	}
	return nil
}

//Print the buffer to stdout
//@return UBF error
func (u *TypedUBF) BPrint() UBFError {
	if ret := C.OBprint(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr)); ret != SUCCEED {
		return u.Buf.Ctx.NewUBFError()
	}
	return nil
}

//Print the buffer to stdout
//@return UBF error
func (u *TypedUBF) TpLogPrintUBF(lev int, title string) {

	c_title := C.CString(title)
	defer C.free(unsafe.Pointer(c_title))

	C.Otplogprintubf(&u.Buf.Ctx.c_ctx, C.int(lev), c_title, C.GetU(u.Buf.C_ptr))

	return
}

//Alternative for Bfprint. Will return the output in string variable
//So that caller can do anything it wants with the string output
//@return Printed buffer, UBF error
func (u *TypedUBF) BSprint() (string, UBFError) {

	c_val := C.calloc(ATMI_MSG_MAX_SIZE, 10)
	c_len := C.size_t(ATMI_MSG_MAX_SIZE * 10)

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

	if ret := C.OBfprint(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), f); ret != SUCCEED {
		return "", u.Buf.Ctx.NewUBFError()
	}

	return C.GoString((*C.char)(c_val)), nil
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

	if ret := C.OBextread(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), f); ret != SUCCEED {
		return u.Buf.Ctx.NewUBFError()
	}

	return nil
}

//Print the expression tree
//@param tree 	Compiled expression tree
//@return printed expresion string, ubf error
func (ac *ATMICtx) BBoolPr(tree *ExprTree) (string, UBFError) {

	c_val := C.calloc(ATMI_MSG_MAX_SIZE, 10)
	c_len := C.size_t(ATMI_MSG_MAX_SIZE * 10)

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

	return C.GoString((*C.char)(c_val)), nil
}

//Serialize the UBF buffer
//@return serialized bytes, UBF error
func (u *TypedUBF) BWrite() ([]byte, UBFError) {

	c_val := C.calloc(ATMI_MSG_MAX_SIZE, 1)
	c_len := C.size_t(ATMI_MSG_MAX_SIZE)

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

	if SUCCEED != C.OBwrite(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), f) {
		return nil, u.Buf.Ctx.NewUBFError()
	}

	size := C.ftell(f)

	var array = make([]byte, int(size))

	for i := 0; i < int(size); i++ {
		array[i] = byte(*(*C.char)(unsafe.Pointer(uintptr(c_val) + uintptr(i))))
	}

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

	if SUCCEED != C.OBread(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr), f) {
		return u.Buf.Ctx.NewUBFError()
	}

	return nil
}

//Test C buffer for UBF format
//@return TRUE - buffer is UBF, FALSE - not UBF
func (u *TypedUBF) BIsUBF() bool {
	c_ret := C.OBisubf(&u.Buf.Ctx.c_ctx, C.GetU(u.Buf.C_ptr))
	if 1 == c_ret {
		return true
	}
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

///////////////////////////////////////////////////////////////////////////////////
// Wrappers for memory management
///////////////////////////////////////////////////////////////////////////////////

func (u *TypedUBF) TpRealloc(size int64) ATMIError {
	return u.Buf.TpRealloc(size)
}
