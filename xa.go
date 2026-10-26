/**
 * @brief Multiple XA resource manager bindings
 *
 * @file xa.go
 */
/* -----------------------------------------------------------------------------
 * Enduro/X Middleware Platform for Distributed Transaction Processing
 * Copyright (C) 2009-2016, ATR Baltic, Ltd. All Rights Reserved.
 * Copyright (C) 2017-2026, Mavimax, Ltd. All Rights Reserved.
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

#include <oatmi.h>
#include <stdlib.h>
*/
import "C"
import (
	"runtime"
	"strings"
	"unsafe"
)

// TpOpen2 opens automatic resource managers, or all with TPTXOPENALL.
// flags selects the resource managers; errors are reported through ATMIError.
// @param flags 0 for automatic resource managers, or TPTXOPENALL to include
//
//	manually managed resource managers
//
// @return Nil on success, otherwise ATMIError
func (ac *ATMICtx) TpOpen2(flags int64) ATMIError {
	defer runtime.KeepAlive(ac)
	if err := ac.checkContext(); nil != err {
		return err
	}
	if C.Otpopen2(&ac.c_ctx, C.long(flags)) != SUCCEED {
		return ac.NewATMIError()
	}
	return nil
}

// TpRMOpen opens the named resource manager.
// name is the configured RM name; flags follows the native API.
// @param name Configured resource-manager name
// @param flags Reserved; must be 0
// @return Nil on success, otherwise ATMIError
func (ac *ATMICtx) TpRMOpen(name string, flags int64) ATMIError {
	defer runtime.KeepAlive(ac)
	if err := ac.checkContext(); nil != err {
		return err
	}
	if strings.IndexByte(name, 0) >= 0 {
		return NewCustomATMIError(TPEINVAL, "RM name contains NUL")
	}
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	if C.Otprmopen(&ac.c_ctx, cName, C.long(flags)) != SUCCEED {
		return ac.NewATMIError()
	}
	return nil
}

// TpRMClose closes the named resource manager.
// name is the configured RM name; flags follows the native API.
// @param name Configured resource-manager name
// @param flags Reserved; must be 0
// @return Nil on success, otherwise ATMIError
func (ac *ATMICtx) TpRMClose(name string, flags int64) ATMIError {
	defer runtime.KeepAlive(ac)
	if err := ac.checkContext(); nil != err {
		return err
	}
	if strings.IndexByte(name, 0) >= 0 {
		return NewCustomATMIError(TPEINVAL, "RM name contains NUL")
	}
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	if C.Otprmclose(&ac.c_ctx, cName, C.long(flags)) != SUCCEED {
		return ac.NewATMIError()
	}
	return nil
}

// TpRMStart associates a manual resource manager with the current transaction.
// name is the configured RM name; flags follows the native API.
// @param name Configured resource-manager name
// @param flags Reserved; must be 0
// @return Nil on success, otherwise ATMIError
func (ac *ATMICtx) TpRMStart(name string, flags int64) ATMIError {
	defer runtime.KeepAlive(ac)
	if err := ac.checkContext(); nil != err {
		return err
	}
	if strings.IndexByte(name, 0) >= 0 {
		return NewCustomATMIError(TPEINVAL, "RM name contains NUL")
	}
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	if C.Otprmstart(&ac.c_ctx, cName, C.long(flags)) != SUCCEED {
		return ac.NewATMIError()
	}
	return nil
}

// TpRMEnd ends the manual resource manager association.
// name is the configured RM name; flags follows the native API.
// @param name Configured resource-manager name
// @param flags Reserved; must be 0
// @return Nil on success, otherwise ATMIError
func (ac *ATMICtx) TpRMEnd(name string, flags int64) ATMIError {
	defer runtime.KeepAlive(ac)
	if err := ac.checkContext(); nil != err {
		return err
	}
	if strings.IndexByte(name, 0) >= 0 {
		return NewCustomATMIError(TPEINVAL, "RM name contains NUL")
	}
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	if C.Otprmend(&ac.c_ctx, cName, C.long(flags)) != SUCCEED {
		return ac.NewATMIError()
	}
	return nil
}

// TpGetConnRM returns the RM-owned native connection for rmid.
// The pointer is borrowed and must not be freed; closing the RM invalidates it.
// @param rmid Configured numeric resource-manager ID
// @return Borrowed native connection pointer and nil error, or nil and
//
//	ATMIError; the RM owns the connection
func (ac *ATMICtx) TpGetConnRM(rmid int) (unsafe.Pointer, ATMIError) {
	defer runtime.KeepAlive(ac)
	if err := ac.checkContext(); nil != err {
		return nil, err
	}
	if int(C.int(rmid)) != rmid {
		return nil, NewCustomATMIError(TPEINVAL, "RMID exceeds native int")
	}
	ptr := C.Otpgetconnrm(&ac.c_ctx, C.int(rmid))
	if nil == ptr {
		return nil, ac.NewATMIError()
	}
	return ptr, nil
}

/* vim: set ts=4 sw=4 et smartindent: */
