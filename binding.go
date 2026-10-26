/**
 * @brief Shared cgo argument validation
 *
 * @file binding.go
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
*/
import "C"
import (
	"fmt"
	"reflect"
	"strings"
)

// checkContext rejects a nil or explicitly released context before entering C.
// @return Nil for a live context, otherwise TPEINVAL
func (ac *ATMICtx) checkContext() ATMIError {
	if nil == ac || nil == ac.c_ctx {
		return NewCustomATMIError(TPEINVAL, "ATMI context is nil or released")
	}
	return nil
}

// copyCText copies a NUL-free string into a fixed native field, including NUL.
// @param dst Native fixed character array including space for NUL
// @param value NUL-free input string
// @param field Field name included in validation errors
// @return Nil after copying, otherwise TPEINVAL
func copyCText(dst []C.char, value, field string) ATMIError {
	if len(value) >= len(dst) || strings.IndexByte(value, 0) >= 0 {
		return NewCustomATMIError(TPEINVAL,
			fmt.Sprintf("%s must contain at most %d bytes and no NUL", field,
				len(dst)-1))
	}
	clear(dst)
	for i := range len(value) {
		dst[i] = C.char(value[i])
	}
	return nil
}

// nativeBuffer resolves typed buffers, including a nil typed pointer.
// @param tb Typed buffer interface, including nil typed pointers
// @return Underlying ATMIBuf, or nil
func nativeBuffer(tb TypedBuffer) *ATMIBuf {
	if nil == tb {
		return nil
	}
	value := reflect.ValueOf(tb)
	if value.Kind() == reflect.Pointer && value.IsNil() {
		return nil
	}
	return tb.GetBuf()
}

/* vim: set ts=4 sw=4 et smartindent: */
