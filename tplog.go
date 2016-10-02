package atmi

/*
#cgo LDFLAGS: -latmisrvinteg -latmi -lrt -lm -lubf -lnstd -ldl
#include <ndebug.h>
#include <xatmi.h>
#include <string.h>
#include <stdlib.h>
#include <ubf.h>

*/
import "C"
import (
	"fmt"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// Logging sub-system tplog*
////////////////////////////////////////////////////////////////////////////////

//Print the byte array buffer to Enduro/X logger (see tplogdump(3) manpage)
//@param lev     Logging level (see LOG_* constants)
//@param comment Title of the buffer dump
//@param ptr   Pointer to buffer for dump
//@param dumplen   Length of the bytes to dump
//@return 	atmiError (in case if invalid length we have for ptr and dumplen)
func TpLogDump(lev int, comment string, ptr []byte, dumplen int) ATMIError {

	c_comment := C.CString(comment)
	defer C.free(unsafe.Pointer(c_comment))
	l1 := len(ptr)

	/* Check the buffer sizes (both must be equal or larger then len) */
	if l1 < dumplen {
		return NewCustomAtmiError(TPEINVAL,
			fmt.Sprintf("ptr len is %d but must be >= %d (len param)",
				l1, dumplen))
	}

	c_ptr := C.malloc(C.size_t(l1))
	defer C.free(c_ptr)

	//Copy stuff to C memory (ptr1)
	for i := 0; i < l1; i++ {
		*(*C.char)(unsafe.Pointer(uintptr(c_ptr) + uintptr(i))) = C.char(ptr[i])
	}

	C.tplogdump(C.int(lev), c_comment, c_ptr, C.int(dumplen))

	return nil
}

//Function compares to byte array buffers and prints the differences to Enduro/X logger
//(see tplogdumpdiff(3) manpage)
//@param lev     Logging level (see LOG_* constants)
//@param comment Title of the buffer diff
//@param ptr1   Pointer to buffer1 for compare
//@param ptr2   Pointer to buffer2 for compare
//@param difflen   Length of the bytes to compare
//@return 	atmiError (in case if invalid length we have for ptr1/ptr2 and difflen)
func TpLogDumpDiff(lev int, comment string, ptr1 []byte, ptr2 []byte, difflen int) ATMIError {

	c_comment := C.CString(comment)
	defer C.free(unsafe.Pointer(c_comment))
	l1 := len(ptr1)
	l2 := len(ptr2)

	/* Check the buffer sizes (both must be equal or larger then len) */
	if l1 < difflen {
		return NewCustomAtmiError(TPEINVAL,
			fmt.Sprintf("ptr1 len is %d but must be >= %d (len param)",
				l1, difflen))
	}

	if l2 < difflen {
		return NewCustomAtmiError(TPEINVAL,
			fmt.Sprintf("ptr2 len is %d but must be >= %d (len param)",
				l2, difflen))
	}

	c_ptr1 := C.malloc(C.size_t(l1))
	defer C.free(c_ptr1)

	//Copy stuff to C memory (ptr1)
	for i := 0; i < l1; i++ {
		*(*C.char)(unsafe.Pointer(uintptr(c_ptr1) + uintptr(i))) = C.char(ptr1[i])
	}

	c_ptr2 := C.malloc(C.size_t(l2))
	defer C.free(c_ptr2)

	//Copy stuff to C memory (ptr1)
	for i := 0; i < l2; i++ {
		*(*C.char)(unsafe.Pointer(uintptr(c_ptr2) + uintptr(i))) = C.char(ptr2[i])
	}

	C.tplogdumpdiff(C.int(lev), c_comment, c_ptr1, c_ptr2, C.int(difflen))

	return nil
}

//Log the message to Enduro/X loggers (see tplog(3) manpage)
//@param lev	Logging level
//@param a	arguemnts for sprintf
//@param format Format string for loggers
func TpLog(lev int, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)

	c_msg := C.CString(msg)
	defer C.free(unsafe.Pointer(c_msg))

	C.tplog(C.int(lev), c_msg)
}

//Return request logging file (if there is one currenlty in use)
// (see tploggetreqfile(3) manpage)
//@return Status (request logger open or not), full path to request file
func TpLogGetReqFile() (bool, string) {

	var status bool
	var reqfile string

	c_reqfile := C.malloc(C.PATH_MAX)
	c_reqfile_ptr := (*C.char)(unsafe.Pointer(c_reqfile))
	defer C.free(c_reqfile)

	if SUCCEED != C.tploggetreqfile(c_reqfile_ptr, C.PATH_MAX) {
		status = false
	} else {
		status = true
		reqfile = C.GoString(c_reqfile_ptr)
	}

	return status, reqfile
}

//Configure Enduro/X logger (see tplogconfig(3) manpage)
//@param logger is bitwise 'ored' (see LOG_FACILITY_*)
//@param lev is optional (if not set: -1), log level to be assigned to facilites
//@param debug_string optional Enduro/X debug string (see ndrxdebug.conf(5) manpage)
//@param new_file optional (if not set - empty string) logging output file, overrides debug_string file tag
//@return NSTDError - standard library error
func TpLogConfig(logger int, lev int, debug_string string, module string, new_file string) NSTDError {

	var err NSTDError
	c_debug_string := C.CString(debug_string)
	defer C.free(unsafe.Pointer(c_debug_string))

	c_module := C.CString(module)
	defer C.free(unsafe.Pointer(c_module))

	c_new_file := C.CString(new_file)
	defer C.free(unsafe.Pointer(c_new_file))

	if SUCCEED != C.tplogconfig(C.int(logger), C.int(lev), c_debug_string,
		c_module, c_new_file) {
		err = NewNstdError()
	}

	return err
}

//Close request logger (see tplogclosereqfile(3) manpage)
func TpLogCloseReqFile() {
	C.tplogclosereqfile()
}

//Close request logger (see tplogclosethread(3) manpage)
func TpLogCloseThread() {
	C.tplogclosethread()
}

//Set request logging file, direct version (see tplogsetreqfile_direct(3) manpage)
//Which does operate with thread local storage
//If fails to open request logging file, it will
//automatically fall-back to stderr.
//@param filename	Set file name to perform logging to
func TpLogSetReqFile_Direct(filename string) {
	c_filename := C.CString(filename)
	defer C.free(unsafe.Pointer(c_filename))

	C.tplogsetreqfile_direct(c_filename)
}

//Set request file to log to (see tplogsetreqfile(3) manpage)
//@param data	pointer to  XATMI buffer (must be UBF, others will cause error), optional
//@param filename	field name to set (this goes to UBF buffer too, if set), optional
//@param filesvc	XATMI service name to call for requesting the new request file name, optional
//@return	ATMI error
func TpLogSetReqFile(data TypedBuffer, filename string, filesvc string) ATMIError {
	var err ATMIError

	c_filename := C.CString(filename)
	defer C.free(unsafe.Pointer(c_filename))

	c_filesvc := C.CString(filesvc)
	defer C.free(unsafe.Pointer(c_filesvc))

	buf := data.GetBuf()

	if SUCCEED != C.tplogsetreqfile(&buf.C_ptr, c_filename, c_filesvc) {
		err = NewAtmiError()
	}

	return err
}

//Get the request file name from UBF buffer (see tploggetbufreqfile(3) manpage)
//@param data	XATMI buffer (must be UBF)
//@return file name, ATMI error
func TpLogGetBufReqFile(data TypedBuffer) (string, ATMIError) {
	var err ATMIError
	var reqfile string

	c_reqfile := C.malloc(C.PATH_MAX)
	c_reqfile_ptr := (*C.char)(unsafe.Pointer(c_reqfile))
	defer C.free(c_reqfile)

	buf := data.GetBuf()

	if SUCCEED != C.tploggetbufreqfile(buf.C_ptr, c_reqfile_ptr, C.PATH_MAX) {
		err = NewAtmiError()
	} else {
		reqfile = C.GoString(c_reqfile_ptr)
	}

	return reqfile, err
}

//Delete request file from UBF buffer (see tplogdelbufreqfile(3) manpage)
//@param data XATMI buffer, must be UBF type
//@return ATMI error
func TpLogDelBufReqFile(data TypedBuffer) ATMIError {
	var err ATMIError

	buf := data.GetBuf()

	if SUCCEED != C.tplogdelbufreqfile(buf.C_ptr) {
		err = NewAtmiError()
	}

	return err
}
