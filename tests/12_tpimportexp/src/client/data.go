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

var M_b1 []byte
var M_b2 []byte

//
//Load  test data into buffer
//
func loadbufferdata(buf *atmi.TypedUBF) atmi.ATMIError {
	if err := buf.BChg(ubftab.T_CHAR_FLD, 0, 65); nil != err {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		return err
	}

	if err := buf.BChg(ubftab.T_CHAR_FLD, 1, 66); nil != err {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		return err
	}

	if err := buf.BChg(ubftab.T_SHORT_FLD, 0, 32000); nil != err {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		return err
	}

	if err := buf.BChg(ubftab.T_SHORT_FLD, 1, 32001); nil != err {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		return err
	}

	if err := buf.BChg(ubftab.T_LONG_FLD, 0, 199101); nil != err {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		return err
	}

	if err := buf.BChg(ubftab.T_LONG_FLD, 1, 199102); nil != err {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		return err
	}

	if err := buf.BChg(ubftab.T_FLOAT_FLD, 0, 9.11); nil != err {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		return err
	}

	if err := buf.BChg(ubftab.T_FLOAT_FLD, 1, 9.22); nil != err {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		return err
	}

	if err := buf.BChg(ubftab.T_DOUBLE_FLD, 0, 19910.888); nil != err {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		return err
	}

	if err := buf.BChg(ubftab.T_DOUBLE_FLD, 1, 19910.999); nil != err {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		return err
	}

	if err := buf.BChg(ubftab.T_STRING_FLD, 0, "HELLO STRING 1"); nil != err {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		return err
	}

	if err := buf.BChg(ubftab.T_STRING_FLD, 1, "HELLO STRING 2"); nil != err {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		return err
	}

	if err := buf.BChg(ubftab.T_CARRAY_FLD, 0, M_b1); nil != err {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		return err
	}

	if err := buf.BChg(ubftab.T_CARRAY_FLD, 1, M_b2); nil != err {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		return err
	}

	return nil
}
