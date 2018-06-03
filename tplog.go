package atmi

/*
** TPLOG - text logging and debuging API provided by Enduro/X
**
** @file tplog.go
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
#include <ndebug.h>
#include <ondebug.h>
#include <xatmi.h>
#include <oatmi.h>
#include <string.h>
#include <stdlib.h>
#include <ubf.h>
#include <oubf.h>
#include <userlog.h>

*/
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

const (
	DETAIL_MODE = "detailed" //Slower, not for production !
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
func (ac *ATMICtx) TpLogDump(lev int, comment string, ptr []byte, dumplen int) ATMIError {

	//NOTE: Checking level here - assume C call in cheaper than
	//String formatting and memory copy in Go!
	if lev <= int(C.debug_get_tp_level()) {
		c_comment := C.CString(comment)
		defer C.free(unsafe.Pointer(c_comment))
		l1 := len(ptr)

		/* Check the buffer sizes (both must be equal or larger then len) */
		if l1 < dumplen {
			return NewCustomATMIError(TPEINVAL,
				fmt.Sprintf("ptr len is %d but must be >= %d (len param)",
					l1, dumplen))
		}

		c_ptr := C.malloc(C.size_t(l1))
		defer C.free(c_ptr)

		//Copy stuff to C memory (ptr1)
		for i := 0; i < l1; i++ {
			*(*C.char)(unsafe.Pointer(uintptr(c_ptr) + uintptr(i))) = C.char(ptr[i])
		}

		C.Otplogdump(&ac.c_ctx, C.int(lev), c_comment, c_ptr, C.int(dumplen))
	}

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
func (ac *ATMICtx) TpLogDumpDiff(lev int, comment string, ptr1 []byte, ptr2 []byte, difflen int) ATMIError {

	if lev <= int(C.debug_get_tp_level()) {
		c_comment := C.CString(comment)
		defer C.free(unsafe.Pointer(c_comment))
		l1 := len(ptr1)
		l2 := len(ptr2)

		/* Check the buffer sizes (both must be equal or larger then len) */
		if l1 < difflen {
			return NewCustomATMIError(TPEINVAL,
				fmt.Sprintf("ptr1 len is %d but must be >= %d (len param)",
					l1, difflen))
		}

		if l2 < difflen {
			return NewCustomATMIError(TPEINVAL,
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

		C.Otplogdumpdiff(&ac.c_ctx, C.int(lev), c_comment, c_ptr1, c_ptr2, C.int(difflen))
	}
	return nil
}

//Log the message to Enduro/X loggers (see tplog(3) manpage)
//This version does not check the debug level (kind of internal one)
//@param lev	Logging level
//@param a	arguemnts for sprintf
//@param format Format string for loggers
func (ac *ATMICtx) tpLog(lev int, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)

	c_msg := C.CString(msg)
	defer C.free(unsafe.Pointer(c_msg))

	if ac.TpLogGetIflags() != DETAIL_MODE {
		C.Otplog(&ac.c_ctx, C.int(lev), c_msg)
	} else {
		//Get the stack and give file name and line
		_, file, line, _ := runtime.Caller(2)

		c_file := C.CString(file)
		defer C.free(unsafe.Pointer(c_file))

		C.Otplogex(&ac.c_ctx, C.int(lev), c_file, C.long(line), c_msg)
	}
}

//Log the message to Enduro/X loggers (see tplog(3) manpage)
//@param lev	Logging level
//@param a	arguemnts for sprintf
//@param format Format string for loggers
func (ac *ATMICtx) TpLog(lev int, format string, a ...interface{}) {

	if lev <= int(C.debug_get_tp_level()) {
		msg := fmt.Sprintf(format, a...)

		c_msg := C.CString(msg)
		defer C.free(unsafe.Pointer(c_msg))

		C.Otplog(&ac.c_ctx, C.int(lev), c_msg)
	}
}

//Log the message to Enduro/X loggers (see tplog(3) manpage), internal ndrx only
//@param lev	Logging level
//@param a	arguemnts for sprintf
//@param format Format string for loggers
func (ac *ATMICtx) ndrxLog(lev int, format string, a ...interface{}) {
	if lev <= int(C.debug_get_ndrx_level()) {

		msg := fmt.Sprintf(format, a...)

		c_msg := C.CString(msg)
		defer C.free(unsafe.Pointer(c_msg))

		C.Ondrxlog(&ac.c_ctx, C.int(lev), c_msg)
	}
}

//Log the message to Enduro/X loggers (see tplog(3) manpage), internal ubf only
//@param lev	Logging level
//@param a	arguemnts for sprintf
//@param format Format string for loggers
func (ac *ATMICtx) ubfLog(lev int, format string, a ...interface{}) {
	if lev <= int(C.debug_get_ubf_level()) {
		msg := fmt.Sprintf(format, a...)

		c_msg := C.CString(msg)
		defer C.free(unsafe.Pointer(c_msg))

		C.Oubflog(&ac.c_ctx, C.int(lev), c_msg)
	}
}

//Log the message to Enduro/X loggers (see tplog(3) manpage)
//Debug level wrapper
//@param a	arguemnts for sprintf
//@param format Format string for loggers
func (ac *ATMICtx) TpLogDebug(format string, a ...interface{}) {
	if LOG_DEBUG <= int(C.debug_get_tp_level()) {
		ac.tpLog(LOG_DEBUG, format, a...)
	}
}

//Log the message to Enduro/X loggers (see tplog(3) manpage)
//Info level wrapper
//@param a	arguemnts for sprintf
//@param format Format string for loggers
func (ac *ATMICtx) TpLogInfo(format string, a ...interface{}) {
	if LOG_INFO <= int(C.debug_get_tp_level()) {
		ac.tpLog(LOG_INFO, format, a...)
	}
}

//Log the message to Enduro/X loggers (see tplog(3) manpage)
//Warning level wrapper
//@param a	arguemnts for sprintf
//@param format Format string for loggers
func (ac *ATMICtx) TpLogWarn(format string, a ...interface{}) {
	if LOG_WARN <= int(C.debug_get_tp_level()) {
		ac.tpLog(LOG_WARN, format, a...)
	}
}

//Log the message to Enduro/X loggers (see tplog(3) manpage)
//Error level wrapper
//@param a	arguemnts for sprintf
//@param format Format string for loggers
func (ac *ATMICtx) TpLogError(format string, a ...interface{}) {
	if LOG_WARN <= int(C.debug_get_tp_level()) {
		ac.tpLog(LOG_ERROR, format, a...)
	}
}

//Log the message to Enduro/X loggers (see tplog(3) manpage)
//Fatal/Always level wrapper
//@param a	arguemnts for sprintf
//@param format Format string for loggers
func (ac *ATMICtx) TpLogAlways(format string, a ...interface{}) {
	if LOG_ALWAYS <= int(C.debug_get_tp_level()) {
		ac.tpLog(LOG_ALWAYS, format, a...)
	}
}

//Log the message to Enduro/X loggers (see tplog(3) manpage)
//Fatal/Always level wrapper
//@param a	arguemnts for sprintf
//@param format Format string for loggers
func (ac *ATMICtx) TpLogFatal(format string, a ...interface{}) {
	if LOG_ALWAYS <= int(C.debug_get_tp_level()) {
		ac.tpLog(LOG_ALWAYS, format, a...)
	}
}

//Return request logging file (if there is one currenlty in use)
// (see tploggetreqfile(3) manpage)
//@return Status (request logger open or not), full path to request file
func (ac *ATMICtx) TpLogGetReqFile() (bool, string) {

	var status bool
	var reqfile string

	c_reqfile := C.malloc(C.PATH_MAX)
	c_reqfile_ptr := (*C.char)(unsafe.Pointer(c_reqfile))
	defer C.free(c_reqfile)

	if SUCCEED != C.Otploggetreqfile(&ac.c_ctx, c_reqfile_ptr, C.PATH_MAX) {
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
func (ac *ATMICtx) TpLogConfig(logger int, lev int, debug_string string, module string, new_file string) NSTDError {

	var err NSTDError
	c_debug_string := C.CString(debug_string)
	defer C.free(unsafe.Pointer(c_debug_string))

	c_module := C.CString(module)
	defer C.free(unsafe.Pointer(c_module))

	c_new_file := C.CString(new_file)
	defer C.free(unsafe.Pointer(c_new_file))

	if SUCCEED != C.Otplogconfig(&ac.c_ctx, C.int(logger), C.int(lev), c_debug_string,
		c_module, c_new_file) {
		err = ac.NewNstdError()
	}

	return err
}

//Close request logger (see tplogclosereqfile(3) manpage)
func (ac *ATMICtx) TpLogCloseReqFile() {
	C.Otplogclosereqfile(&ac.c_ctx)
}

//Close request logger (see tplogclosethread(3) manpage)
func (ac *ATMICtx) TpLogCloseThread() {
	C.Otplogclosethread(&ac.c_ctx)
}

//Set request logging file, direct version (see tplogsetreqfile_direct(3) manpage)
//Which does operate with thread local storage
//If fails to open request logging file, it will
//automatically fall-back to stderr.
//@param filename	Set file name to perform logging to
func (ac *ATMICtx) TpLogSetReqFileDirect(filename string) {
	c_filename := C.CString(filename)
	defer C.free(unsafe.Pointer(c_filename))

	C.Otplogsetreqfile_direct(&ac.c_ctx, c_filename)
}

//Set request file to log to (see tplogsetreqfile(3) manpage)
//@param data	pointer to  XATMI buffer (must be UBF, others will cause error), optional
//@param filename	field name to set (this goes to UBF buffer too, if set), optional
//@param filesvc	XATMI service name to call for requesting the new request file name, optional
//@return	ATMI error
func (ac *ATMICtx) TpLogSetReqFile(data TypedBuffer, filename string, filesvc string) ATMIError {
	var err ATMIError

	c_filename := C.CString(filename)
	defer C.free(unsafe.Pointer(c_filename))

	c_filesvc := C.CString(filesvc)
	defer C.free(unsafe.Pointer(c_filesvc))

	buf := data.GetBuf()

	if SUCCEED != C.Otplogsetreqfile(&ac.c_ctx, &buf.C_ptr, c_filename, c_filesvc) {
		err = ac.NewATMIError()
	}

	return err
}

//Get the request file name from UBF buffer (see tploggetbufreqfile(3) manpage)
//@param data	XATMI buffer (must be UBF)
//@return file name, ATMI error
func (ac *ATMICtx) TpLogGetBufReqFile(data TypedBuffer) (string, ATMIError) {
	var err ATMIError
	var reqfile string

	c_reqfile := C.malloc(C.PATH_MAX)
	c_reqfile_ptr := (*C.char)(unsafe.Pointer(c_reqfile))
	defer C.free(c_reqfile)

	buf := data.GetBuf()

	if SUCCEED != C.Otploggetbufreqfile(&ac.c_ctx, buf.C_ptr, c_reqfile_ptr, C.PATH_MAX) {
		err = ac.NewATMIError()
	} else {
		reqfile = C.GoString(c_reqfile_ptr)
	}

	return reqfile, err
}

//Delete request file from UBF buffer (see tplogdelbufreqfile(3) manpage)
//@param data XATMI buffer, must be UBF type
//@return ATMI error
func (ac *ATMICtx) TpLogDelBufReqFile(data TypedBuffer) ATMIError {
	var err ATMIError

	buf := data.GetBuf()

	if SUCCEED != C.Otplogdelbufreqfile(&ac.c_ctx, buf.C_ptr) {
		err = ac.NewATMIError()
	}

	return err
}

//Return integration flags
//Well we will run it in cached mode...
func (ac *ATMICtx) TpLogGetIflags() string {

	return C.GoString(C.tploggetiflags())
}

//Do the user logging. This prints the message to ULOG. Suitable for system wide
//critical message notifications
//@param format	format string
//@param a list of data fields for format string
func (ac *ATMICtx) UserLog(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)

	c_msg := C.CString(msg)
	defer C.free(unsafe.Pointer(c_msg))

	C.userlog_const(c_msg)
}
