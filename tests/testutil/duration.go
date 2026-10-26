/**
 * @brief Shared integration test duration
 *
 * @file duration.go
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
package testutil

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Duration reads the stress-scenario duration; an unset or empty value selects ten seconds.
// @return Positive duration, or an error for an invalid NDRX_TEST_DURATATION
func Duration() (time.Duration, error) {
	value := os.Getenv("NDRX_TEST_DURATATION")
	if value == "" {
		value = "10"
	}
	for _, digit := range value {
		if digit < '0' || digit > '9' {
			return 0, fmt.Errorf("NDRX_TEST_DURATATION must be a positive integer in seconds")
		}
	}
	seconds, err := strconv.ParseInt(value, 10, 32)
	if err != nil || seconds <= 0 {
		return 0, fmt.Errorf("NDRX_TEST_DURATATION must be between 1 and 2147483647 seconds")
	}
	return time.Duration(seconds) * time.Second, nil
}

/* vim: set ts=4 sw=4 et smartindent: */
