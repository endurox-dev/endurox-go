/*
** Test buffer printing routines
**
** @file bprint.go
** -----------------------------------------------------------------------------
** Enduro/X Middleware Platform for Distributed Transaction Processing
** Copyright (C) 2015, ATR Baltic, Ltd. All Rights Reserved.
** This software is released under one of the following licenses:
** GPL or ATR Baltic's license for commercial use.
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
** A commercial use license is available from ATR Baltic, Ltd
** contact@atrbaltic.com
** -----------------------------------------------------------------------------
 */
package main

/*
#include <signal.h>
*/
import "C"

import (
	"atmi"
	"fmt"
	"ubftab"
)

//Perform fast add tests
func test_BprintStr() error {

	for i := 0; i < 10000; i++ {

		ac, err := atmi.NewATMICtx()

		if nil != err {
			return fmt.Errorf("Failed to allocate cotnext!", err)
		}

		buf, err := ac.NewUBF(1024)

		if err != nil {
			return fmt.Errorf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		}

		var loc atmi.BFldLocInfo
		//////////////////////////////////////////////////////////////////////////
		//Short tests
		//////////////////////////////////////////////////////////////////////////
		var s int16
		s = 5

		if err := buf.BAddFast(ubftab.T_SHORT_FLD, s, &loc, true); nil != err {
			return fmt.Errorf("UBF Error %d:[%s]", err.Code(), err.Message())
		}

		s = 7
		if err := buf.BAddFast(ubftab.T_SHORT_FLD, s, &loc, false); nil != err {
			return fmt.Errorf("UBF Error %d:[%s]", err.Code(), err.Message())
		}

		s = 8
		if err := buf.BAddFast(ubftab.T_SHORT_FLD, s, &loc, false); nil != err {
			return fmt.Errorf("UBF Error %d:[%s]", err.Code(), err.Message())
		}

		//////////////////////////////////////////////////////////////////////////
		//Long tests
		//////////////////////////////////////////////////////////////////////////
		var l int64
		l = 5777

		if err := buf.BAddFast(ubftab.T_LONG_2_FLD, l, &loc, true); nil != err {
			return fmt.Errorf("UBF Error %d:[%s]", err.Code(), err.Message())
		}

		l = 7888
		if err := buf.BAddFast(ubftab.T_LONG_2_FLD, l, &loc, false); nil != err {
			return fmt.Errorf("UBF Error %d:[%s]", err.Code(), err.Message())
		}

		l = -8999
		if err := buf.BAddFast(ubftab.T_LONG_2_FLD, l, &loc, false); nil != err {
			return fmt.Errorf("UBF Error %d:[%s]", err.Code(), err.Message())
		}

		//////////////////////////////////////////////////////////////////////////
		//Char tests
		//////////////////////////////////////////////////////////////////////////
		var c string
		c = "C"

		if err := buf.BAddFast(ubftab.T_CHAR_FLD, c, &loc, true); nil != err {
			return fmt.Errorf("UBF Error %d:[%s]", err.Code(), err.Message())
		}

		c = "X"
		if err := buf.BAddFast(ubftab.T_CHAR_FLD, c, &loc, false); nil != err {
			return fmt.Errorf("UBF Error %d:[%s]", err.Code(), err.Message())
		}

		c = "a"
		if err := buf.BAddFast(ubftab.T_CHAR_FLD, c, &loc, false); nil != err {
			return fmt.Errorf("UBF Error %d:[%s]", err.Code(), err.Message())
		}

		//////////////////////////////////////////////////////////////////////////
		//float tests
		//////////////////////////////////////////////////////////////////////////
		var f float32
		f = 1.1

		if err := buf.BAddFast(ubftab.T_FLOAT_FLD, f, &loc, true); nil != err {
			return fmt.Errorf("UBF Error %d:[%s]", err.Code(), err.Message())
		}

		f = -2.2
		if err := buf.BAddFast(ubftab.T_FLOAT_FLD, f, &loc, false); nil != err {
			return fmt.Errorf("UBF Error %d:[%s]", err.Code(), err.Message())
		}

		f = 3.14
		if err := buf.BAddFast(ubftab.T_FLOAT_FLD, f, &loc, false); nil != err {
			return fmt.Errorf("UBF Error %d:[%s]", err.Code(), err.Message())
		}

		//////////////////////////////////////////////////////////////////////////
		//double tests
		//////////////////////////////////////////////////////////////////////////
		var d float64
		d = 111.1

		if err := buf.BAddFast(ubftab.T_DOUBLE_FLD, d, &loc, true); nil != err {
			return fmt.Errorf("UBF Error %d:[%s]", err.Code(), err.Message())
		}

		d = -222.2
		if err := buf.BAddFast(ubftab.T_DOUBLE_FLD, d, &loc, false); nil != err {
			return fmt.Errorf("UBF Error %d:[%s]", err.Code(), err.Message())
		}

		d = 333.14
		if err := buf.BAddFast(ubftab.T_DOUBLE_FLD, d, &loc, false); nil != err {
			return fmt.Errorf("UBF Error %d:[%s]", err.Code(), err.Message())
		}

		//////////////////////////////////////////////////////////////////////////
		//string tests
		//////////////////////////////////////////////////////////////////////////
		var r string
		r = "hello"

		if err := buf.BAddFast(ubftab.T_STRING_FLD, r, &loc, true); nil != err {
			return fmt.Errorf("UBF Error %d:[%s]", err.Code(), err.Message())
		}

		r = "world"
		if err := buf.BAddFast(ubftab.T_STRING_FLD, r, &loc, false); nil != err {
			return fmt.Errorf("UBF Error %d:[%s]", err.Code(), err.Message())
		}

		r = "from mars"
		if err := buf.BAddFast(ubftab.T_STRING_FLD, r, &loc, false); nil != err {
			return fmt.Errorf("UBF Error %d:[%s]", err.Code(), err.Message())
		}

		//////////////////////////////////////////////////////////////////////////
		//carray tests
		//////////////////////////////////////////////////////////////////////////

		a := []byte("this")

		if err := buf.BAddFast(ubftab.T_CARRAY_FLD, a, &loc, true); nil != err {
			return fmt.Errorf("UBF Error %d:[%s]", err.Code(), err.Message())
		}

		a = []byte("is")
		if err := buf.BAddFast(ubftab.T_CARRAY_FLD, a, &loc, false); nil != err {
			return fmt.Errorf("UBF Error %d:[%s]", err.Code(), err.Message())
		}

		a = []byte("carray")
		if err := buf.BAddFast(ubftab.T_CARRAY_FLD, a, &loc, false); nil != err {
			return fmt.Errorf("UBF Error %d:[%s]", err.Code(), err.Message())
		}

		//Test some error too, missing loc
		err = buf.BAddFast(ubftab.T_CARRAY_FLD, a, nil, false)

		if nil == err {
			return fmt.Errorf("Expected error bot ok! (1)")
		}

		if atmi.BEINVAL != err.Code() {
			return fmt.Errorf("Expected error code BEINVAL but got %d (1)", err.Code())
		}

		err = buf.BAddFast(ubftab.T_SHORT_FLD, a, &loc, false)

		if nil == err {
			return fmt.Errorf("Expected error bot ok! (2)")
		}

		if atmi.BEINVAL != err.Code() {
			return fmt.Errorf("Expected error code BEINVAL but got %d (2)", err.Code())
		}

		//Now validate the buffer by mega boolean expression

		buf.BPrint()

		//Now transfer the buffer to string

		str, err := buf.BSprint()

		if nil != err {
			return fmt.Errorf("Failed to print to str: %s\n", err.Error())
		}

		fmt.Printf("Got string: [%s]\n", str)

		//Reset buffer
		buf, err = ac.NewUBF(1024)

		if err != nil {
			return fmt.Errorf("Failed to realloc %d:[%s]\n", err.Code(), err.Message())
		}

		err = buf.BExtRead(str)

		if nil != err {
			return fmt.Errorf("Failed to extread str [%s]: %s\n", str, err.Error())
		}

		res, err := buf.BQBoolEv("T_SHORT_FLD[0]==5 && T_SHORT_FLD[1]==7 && T_SHORT_FLD[2]==8 && " +
			"T_LONG_2_FLD[0]==5777 && T_LONG_2_FLD[1]==7888 && T_LONG_2_FLD[2]==-8999 && " +
			"T_CHAR_FLD[0]=='C' && T_CHAR_FLD[1]=='X' && T_CHAR_FLD[2]=='a' && " +
			"T_FLOAT_FLD[0]==1.1 && T_FLOAT_FLD[1]==-2.2 && T_FLOAT_FLD[2]==3.14 && " +
			"T_DOUBLE_FLD[0]==111.1 && T_DOUBLE_FLD[1]==-222.2 && T_DOUBLE_FLD[2]==333.14 && " +
			"T_STRING_FLD[0]=='hello' && T_STRING_FLD[1]=='world' && T_STRING_FLD[2]=='from mars' && " +
			"T_CARRAY_FLD[0]=='this' && T_CARRAY_FLD[1]=='is' && T_CARRAY_FLD[2]=='carray'")
		if nil != err {
			return fmt.Errorf("BQBoolEv failed %d:[%s]", err.Code(), err.Message())
		}

		if !res {
			return fmt.Errorf("Expected expression to be true, but got false!")
		}

	}

	return nil
}
