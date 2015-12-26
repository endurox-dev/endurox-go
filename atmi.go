package atmi

/*
#cgo LDFLAGS: -latmi -lrt -lm -lubf -lnstd -ldl

#include <xatmi.h>
#include <string.h>
#include <stdlib.h>

// Wrapper for TPNIT
static int go_tpinit(void) {
	return tpinit(NULL);
}

static int go_tperrno(void) {
	return tperrno;
}


static void free_string(char* s) { free(s); }
static char * malloc_string(int size) { return malloc(size); }

*/
import "C"
import "unsafe"
import "fmt"
import "runtime"

//TODO: Think about runtime.SetFinalizer - might be usable for ATMI buffer free
//      and for UBF expression dealloc

/*
 * SUCCEED/FAIL flags
 */
const (
	SUCCEED = 0
	FAIL    = -1
)

/*
 * List of ATMI Error codes
 */
const (
	TPMINVAL      = 0
	TPEABORT      = 1
	TPEBADDESC    = 2
	TPEBLOCK      = 3
	TPEINVAL      = 4
	TPELIMIT      = 5
	TPENOENT      = 6
	TPEOS         = 7
	TPEPERM       = 8
	TPEPROTO      = 9
	TPESVCERR     = 10
	TPESVCFAIL    = 11
	TPESYSTEM     = 12
	TPETIME       = 13
	TPETRAN       = 14
	TPGOTSIG      = 15
	TPERMERR      = 16
	TPEITYPE      = 17
	TPEOTYPE      = 18
	TPERELEASE    = 19
	TPEHAZARD     = 20
	TPEHEURISTIC  = 21
	TPEEVENT      = 22
	TPEMATCH      = 23
	TPEDIAGNOSTIC = 24
	TPEMIB        = 25
	TPINITFAIL    = 30
	TPMAXVAL      = 31
)

/*
 * flag bits for C language xatmi routines
 */
const (
	TPNOBLOCK     = 0x00000001
	TPSIGRSTRT    = 0x00000002
	TPNOREPLY     = 0x00000004
	TPNOTRAN      = 0x00000008
	TPTRAN        = 0x00000010
	TPNOTIME      = 0x00000020
	TPGETANY      = 0x00000080
	TPNOCHANGE    = 0x00000100
	TPCONV        = 0x00000400
	TPSENDONLY    = 0x00000800
	TPRECVONLY    = 0x00001000
	TPTRANSUSPEND = 0x00040000 /* Suspend current transaction */
)

/*
 * values for rval in tpreturn
 */
const (
	TPFAIL    = 0x0001
	TPSUCCESS = 0x0002
)

/*
 * Max message size (int bytes)
 */
const (
	ATMI_MSG_MAX_SIZE = 65536
)

/*
 * Transaction ID type
 */
type TPTRANID struct {
	c_tptranid C.TPTRANID
}

/*
 * Server context data (used for server's main thread
 * switching taks to worker thread)
 */
type TPSRVCTXDATA struct {
	c_ptr *C.char
}

/*
 * Event controll struct
 */
type TPEVCTL struct {
	flags int64
	name1 string
	name2 string
}

///////////////////////////////////////////////////////////////////////////////////
// ATMI Buffers section
///////////////////////////////////////////////////////////////////////////////////

//ATMI buffer
type ATMIBuf struct {
	C_ptr *C.char
	//We will need some API for length & buffer setting
	//Probably we need a wrapper for lenght function
	C_len C.long
}

//Base interface for typed buffer
type TypedBuffer interface {
	GetBuf() *ATMIBuf
}

//Have inteface to base ATMI buffer
func (u *ATMIBuf) GetBuf() *ATMIBuf {
	return u
}

///////////////////////////////////////////////////////////////////////////////////
// Error Handlers
///////////////////////////////////////////////////////////////////////////////////

//ATMI Error type
type atmiError struct {
	code    int
	message string
}

//ATMI error interface
type ATMIError interface {
	Error() string
	Code() int
	Message() string
}

//Generate ATMI error, read the codes
func NewAtmiError() ATMIError {
	var err atmiError
	err.code = int(C.go_tperrno())
	err.message = C.GoString(C.tpstrerror(C.go_tperrno()))
	return err
}

//Build a custom error
//@param err		Error buffer to build
//@param code	Error code to setup
//@param msg		Error message
func NewCustomAtmiError(code int, msg string) ATMIError {
	var err atmiError
	err.code = code
	err.message = msg
	return err
}

//Standard error interface
func (e atmiError) Error() string {
	return fmt.Sprintf("%d: %s", e.code, e.message)
}

//code getter
func (e atmiError) Code() int {
	return e.code
}

//message getter
func (e atmiError) Message() string {
	return e.message
}

///////////////////////////////////////////////////////////////////////////////////
// API Section
///////////////////////////////////////////////////////////////////////////////////

//TODO, maybe we need to use error deligates, so that user can override the error handling object?

//Allocate buffer
//Accepts the standard ATMI values
//We should add error handling here
//@param	 b_type 		Buffer type
//@param	 b_subtype 	Buffer sub-type
//@param	 size		Buffer size request
//@return 			ATMI Buffer, atmiError
func TpAlloc(b_type string, b_subtype string, size int64) (*ATMIBuf, ATMIError) {
	var buf ATMIBuf
	var err ATMIError

	c_type := C.CString(b_type)
	c_subtype := C.CString(b_subtype)

	size_l := C.long(size)

	buf.C_ptr = C.tpalloc(c_type, c_subtype, size_l)

	//Check the error
	if nil == buf.C_ptr {
		err = NewAtmiError()
	}

	C.free(unsafe.Pointer(c_type))
	C.free(unsafe.Pointer(c_subtype))

	runtime.SetFinalizer(&buf, TpFree)

	return &buf, err
}

//Reallocate the buffer
//@param buf		ATMI buffer
//@return 		ATMI Error
func (buf *ATMIBuf) TpRealloc(size int64) ATMIError {
	var err ATMIError

	buf.C_ptr = C.tprealloc(buf.C_ptr, C.long(size))

	if nil == buf.C_ptr {
		err = NewAtmiError()
	}

	return err
}

//Initialize client
//@return		ATMI Error
func TpInit() ATMIError {
	var err ATMIError

	if SUCCEED != C.go_tpinit() {
		err = NewAtmiError()
	}

	return err
}

// Do the service call, assume using the same buffer
// for return value.
// This works for self describing buffers. Otherwise we need a buffer size in
// ATMIBuf.
// @param svc	service name
// @param buf	ATMI buffer
// @param flags 	Flags to be used
// @return atmiError
func TpCall(svc string, tb TypedBuffer, flags int64) (int, ATMIError) {
	var err ATMIError
	c_svc := C.CString(svc)

	buf := tb.GetBuf()

	ret := C.tpcall(c_svc, buf.C_ptr, buf.C_len, &buf.C_ptr, &buf.C_len, C.long(flags))

	if SUCCEED != ret {
		err = NewAtmiError()
	}

	C.free(unsafe.Pointer(c_svc))

	return int(ret), err
}

//TP Async call
//@param svc		Service Name to call
//@param buf		ATMI buffer
//@param flags	Flags to be used for call (see flags section)
//@return		Call Descriptor (cd), ATMI Error
func TpACall(svc string, tb TypedBuffer, flags int64) (int, ATMIError) {
	var err ATMIError
	c_svc := C.CString(svc)

	buf := tb.GetBuf()

	ret := C.tpacall(c_svc, buf.C_ptr, buf.C_len, C.long(flags))

	if FAIL == ret {
		err = NewAtmiError()
	}

	C.free(unsafe.Pointer(c_svc))

	return int(ret), err
}

//Get async call reply
//@param cd	call
//@param buf	ATMI buffer
//@param flags call flags
func TpGetRply(cd *int, tb TypedBuffer, flags int64) (int, ATMIError) {
	var err ATMIError
	var c_cd C.int

	buf := tb.GetBuf()

	ret := C.tpgetrply(&c_cd, &buf.C_ptr, &buf.C_len, C.long(flags))

	if SUCCEED != ret {
		err = NewAtmiError()
	} else {
		*cd = int(c_cd)
	}

	return int(ret), err
}

//Cancel async call
//@param cd		Call descriptor
//@return ATMI error
func TpCancel(cd int) ATMIError {
	var err ATMIError

	ret := C.tpcancel(C.int(cd))

	if SUCCEED != ret {
		err = NewAtmiError()
	}

	return err
}

//Connect to service in conversational mode
//@param svc		Service name
//@param data	ATMI buffers
//@param flags	Flags
//@return		call descriptor (cd), ATMI error
func TpConnect(svc string, tb TypedBuffer, flags int64) (int, ATMIError) {
	var err ATMIError
	c_svc := C.CString(svc)

	data := tb.GetBuf()

	ret := C.tpconnect(c_svc, data.C_ptr, data.C_len, C.long(flags))

	if FAIL == ret {
		err = NewAtmiError()
	}

	C.free(unsafe.Pointer(c_svc))

	return int(ret), err
}

//Disconnect from conversation
//@param cd		Call Descriptor
//@return ATMI Error
func TpDiscon(cd int) ATMIError {
	var err ATMIError

	ret := C.tpdiscon(C.int(cd))

	if SUCCEED != ret {
		err = NewAtmiError()
	}

	return err
}

//Receive data from conversation
//@param cd			call descriptor
//@param	 data		ATMI buffer
//@param revent		Return Event
//@return			ATMI Error
func TpRecv(cd int, tb TypedBuffer, flags int64, revent *int64) ATMIError {
	var err ATMIError

	c_revent := C.long(*revent)

	data := tb.GetBuf()

	ret := C.tprecv(C.int(cd), &data.C_ptr, &data.C_len, C.long(flags), &c_revent)

	if FAIL == ret {
		err = NewAtmiError()
	}

	*revent = int64(c_revent)

	return err
}

//Receive data from conversation
//@param cd			call descriptor
//@param	 data		ATMI buffer
//@param revent		Return Event
//@return			ATMI Error
func TpSend(cd int, tb TypedBuffer, flags int64, revent *int64) ATMIError {
	var err ATMIError

	c_revent := C.long(*revent)

	data := tb.GetBuf()

	ret := C.tpsend(C.int(cd), data.C_ptr, data.C_len, C.long(flags), &c_revent)

	if SUCCEED != ret {
		err = NewAtmiError()
	}

	*revent = int64(c_revent)

	return err
}

//Free the ATMI buffer
//@param buf		ATMI buffer
func TpFree(buf *ATMIBuf) {
	C.tpfree(buf.C_ptr)
	buf.C_ptr = nil
}

//Commit global transaction
//@param	 flags		flags for abort operation
func TpCommit(flags int64) ATMIError {
	var err ATMIError

	ret := C.tpcommit(C.long(flags))

	if SUCCEED != ret {
		err = NewAtmiError()
	}

	return err
}

//Abort global transaction
//@param	 flags		flags for abort operation (must be 0)
//@return ATMI Error
func TpAbort(flags int64) ATMIError {
	var err ATMIError

	ret := C.tpabort(C.long(flags))

	if SUCCEED != ret {
		err = NewAtmiError()
	}

	return err
}

//Open XA Sub-system
//@return ATMI Error
func TpOpen() ATMIError {
	var err ATMIError

	ret := C.tpopen()

	if SUCCEED != ret {
		err = NewAtmiError()
	}

	return err
}

// Close XA Sub-system
//@return ATMI Error
func TpClose() ATMIError {
	var err ATMIError

	ret := C.tpclose()

	if SUCCEED != ret {
		err = NewAtmiError()
	}

	return err
}

//Check are we in globa transaction?
//@return 	0 - not in global Tx, 1 - in global Tx
func TpGetLev() int {

	ret := C.tpgetlev()

	return int(ret)
}

//Begin transaction
//@param timeout		Transaction Timeout
//@param flags		Transaction flags
//@return	ATMI Error
func TpBegin(timeout uint64, flags int64) ATMIError {

	var err ATMIError

	ret := C.tpbegin(C.ulong(timeout), C.long(flags))

	if SUCCEED != ret {
		err = NewAtmiError()
	}

	return err
}

//Suspend transaction
//@param tranid	Transaction Id reference
//@param flags	Flags for suspend (must be 0)
//@return 	ATMI Error
func TpSuspend(tranid *TPTRANID, flags int64) ATMIError {
	var err ATMIError

	ret := C.tpsuspend(&tranid.c_tptranid, C.long(flags))

	if SUCCEED != ret {
		err = NewAtmiError()
	}

	return err
}

//Resume transaction
//@param tranid	Transaction Id reference
//@param flags	Flags for tran resume (must be 0)
//@return 	ATMI Error
func TpResume(tranid *TPTRANID, flags int64) ATMIError {
	var err ATMIError

	ret := C.tpresume(&tranid.c_tptranid, C.long(flags))

	if SUCCEED != ret {
		err = NewAtmiError()
	}

	return err
}

//Get cluster node id
//@return		Node Id
func TpGetnodeId() int64 {
	ret := C.tpgetnodeid()
	return int64(ret)
}

//Post the event to subscribers
//@param eventname	Name of the event to post
//@param data		ATMI buffer
//@param flags		flags
//@return		Number Of events posted, ATMI error
func TpPost(eventname string, tb TypedBuffer, len int64, flags int64) (int, ATMIError) {
	var err ATMIError
	c_eventname := C.CString(eventname)

	data := tb.GetBuf()
	ret := C.tppost(c_eventname, data.C_ptr, data.C_len, C.long(flags))

	if FAIL == ret {
		err = NewAtmiError()
	}

	C.free(unsafe.Pointer(c_eventname))

	return int(ret), err
}

//Return ATMI buffer info
//@param ptr 	Pointer to ATMI buffer
//@param itype	ptr to string to return the buffer type  (can be nil)
//@param subtype ptr to string to return sub-type (can be nil)
func TpTypes(ptr *ATMIBuf, itype *string, subtype *string) (int64, ATMIError) {
	var err ATMIError

	/* we should allocat the fields there...  */

	var c_type *C.char
	var c_subtype *C.char

	c_type = C.malloc_string(16)
	c_subtype = C.malloc_string(16)

	ret := C.tptypes(ptr.C_ptr, c_type, c_subtype)

	if FAIL == ret {
		err = NewAtmiError()
	} else {
		if nil != itype && nil != c_type {
			*itype = C.GoString(c_type)
		}

		if nil != subtype && nil != c_subtype {
			*subtype = C.GoString(c_subtype)
		}
	}

	if nil != c_type {
		C.free_string(c_type)
	}

	if nil != c_subtype {
		C.free_string(c_subtype)
	}

	return int64(ret), err
}
