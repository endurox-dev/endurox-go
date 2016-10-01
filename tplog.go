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

//Print the byte array buffer to Enduro/X logger
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

//Log the message to Enduro/X loggers
//@param lev	Logging level
//@param a	arguemnts for sprintf
//@param format Format string for loggers
func TpLog(lev int, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)

	c_msg := C.CString(msg)
	defer C.free(unsafe.Pointer(c_msg))

	C.tplog(C.int(lev), c_msg)
}
