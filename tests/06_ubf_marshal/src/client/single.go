package main

/*
#include <signal.h>
*/
import "C"

import (
	"atmi"
	"bytes"
	"fmt"
	"runtime"

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

//Binary main entry
func async_main_one() {

	ret := SUCCEED

	//Close the ATMI session at exit.
	defer func() {
		M_wg.Done()
		M_ret <- ret
	}()

	//Have some loop for memory leak checks...
	for i := 0; i < 100; i++ {
		var ac *atmi.ATMICtx
		var err atmi.ATMIError
		//Allocate context
		ac, err = atmi.NewATMICtx()
		if nil != err {
			fmt.Errorf("Failed to allocate cotnext!", err)
			ret = FAIL
			return
		}
		buf, err := ac.NewUBF(1024)
		buf2, err := ac.NewUBF(1024)

		var s TestStruct

		////////////////////////////////////////////////////////////////////////
		// Set the test data
		////////////////////////////////////////////////////////////////////////
		if err != nil {
			fmt.Printf("async_main_one: ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			return
		}

		if err := loadbufferdata(buf); nil != err {
			fmt.Printf("async_main_one: Unexpected Failed to set UBF data: %d:[%s]\n",
				err.Code(), err.Message())
			ret = FAIL
			return
		}

		//Unmarshal single, non existing instance, no errors
		if err := buf.UnmarshalSingle(&s, 5); nil != err {
			fmt.Printf("async_main_one: Unexpected error at UnmarshalSingle: %d:[%s]\n",
				err.Code(), err.Message())
			ret = FAIL
			return
		}

		////////////////////////////////////////////////////////////////////////
		// Do the Unmarshal & test the values
		// UBF -> Struct
		////////////////////////////////////////////////////////////////////////
		buf.UnmarshalSingle(&s, 1)

		fmt.Printf("Got the struct [%v]\n", s)

		//CHAR
		if s.CharTest != 66 {
			fmt.Printf("s.CharTest invalid value!\n")
			ret = FAIL
			return
		}

		if len(s.CharArrayTest) != 1 || s.CharArrayTest[0] != 66 {
			fmt.Printf("s.CharArrayTest invalid value!\n")
			ret = FAIL
			return
		}

		//SHORT
		if s.ShortTest != 32001 {
			fmt.Printf("s.ShortTest invalid value!\n")
			ret = FAIL
			return
		}

		if len(s.ShortArrayTest) != 1 || s.ShortArrayTest[0] != 32001 {
			fmt.Printf("s.ShortArrayTest invalid value!\n")
			ret = FAIL
			return
		}

		//LONG
		if s.LongTest != 199102 {
			fmt.Printf("s.LongTest invalid value!\n")
			ret = FAIL
			return
		}

		if len(s.LongArrayTest) != 1 || s.LongArrayTest[0] != 199102 {
			fmt.Printf("s.LongArrayTest invalid value!\n")
			ret = FAIL
			return
		}

		//FLOAT
		if s.Float32Test-9.22 > 0.0001 {
			fmt.Printf("s.Float32Test invalid value!\n")
			ret = FAIL
			return
		}

		if len(s.Float32ArrayTest) != 1 || s.Float32ArrayTest[0]-9.22 > 0.0001 {
			fmt.Printf("s.Float32ArrayTest invalid value!\n")
			ret = FAIL
			return
		}

		//FLOAT64
		if s.Float64Test-19910.999 > 0.0001 {
			fmt.Printf("s.Float64Test invalid value!\n")
			ret = FAIL
			return
		}

		if len(s.Float64ArrayTest) != 1 || s.Float64ArrayTest[0]-19910.999 > 0.0001 {
			fmt.Printf("s.Float64ArrayTest invalid value!\n")
			ret = FAIL
			return
		}

		//STRING
		if s.StringTest != "HELLO STRING 2" {
			fmt.Printf("s.StringTest invalid value!\n")
			ret = FAIL
			return
		}

		if len(s.StringArrayTest) != 1 || s.StringArrayTest[0] != "HELLO STRING 2" {
			fmt.Printf("s.StringArrayTest invalid value!\n")
			ret = FAIL
			return
		}

		//CARRAY
		if len(s.CarrayTest) != 1 || 0 != bytes.Compare(s.CarrayTest[0], M_b2) {
			fmt.Printf("s.CarrayTest invalid value!\n")
			ret = FAIL
			return
		}

		fmt.Println("Unmarshal tests ok...\n")

		//Marshal of non existing array element shall cause error!

		if err := loadbufferdata(buf); nil != err {
			fmt.Printf("async_main_one: Unexpected Failed to set UBF data: %d:[%s]\n",
				err.Code(), err.Message())
			ret = FAIL
			return
		}

		if err := buf.Unmarshal(&s); err != nil {
			fmt.Printf("async_main_one: Failed to Unmarshal: %d:[%s]\n",
				err.Code(), err.Message())
			ret = FAIL
			return
		}

		if err := buf2.MarshalSingle(&s, 2); err == nil {
			fmt.Printf("Expected error at occurrence 2, but got OK!")
			ret = FAIL
			return
		}

		////////////////////////////////////////////////////////////////////////
		// DO the marshal tests...
		// Test the stuff with boolean expressions...
		////////////////////////////////////////////////////////////////////////
		//Reset the buffer & generate it again...
		if err := buf2.MarshalSingle(&s, 1); err != nil {
			fmt.Printf("Failed to marshal [%s]\n", err.Error())
			ret = FAIL
			return
		}
		buf2.BPrint()

		//CHAR
		if res, err := buf2.BQBoolEv("T_CHAR_FLD[0]=='B'"); !res || nil != err {
			if nil != err {
				fmt.Printf("char: Expression failed: %s\n", err.Error())
				ret = FAIL
				return
			} else {
				fmt.Printf("char: Expression is false\n")
				ret = FAIL
				return
			}
		}

		//SHORT
		if res, err := buf2.BQBoolEv("T_SHORT_FLD[0]==32001"); !res || nil != err {
			if nil != err {
				fmt.Printf("short: Expression failed: %s\n", err.Error())
				ret = FAIL
				return
			} else {
				fmt.Printf("short: Expression is false\n")
				ret = FAIL
				return
			}
		}

		//LONG
		if res, err := buf2.BQBoolEv("T_LONG_FLD[0]==199102"); !res || nil != err {
			if nil != err {
				fmt.Printf("long: Expression failed: %s\n", err.Error())
				ret = FAIL
				return
			} else {
				fmt.Printf("long: Expression is false\n")
				ret = FAIL
				return
			}
		}

		//FLOAT
		if res, err := buf2.BQBoolEv("T_FLOAT_FLD[0]==9.22"); !res || nil != err {
			if nil != err {
				fmt.Printf("float: Expression failed: %s\n", err.Error())
				ret = FAIL
				return
			} else {
				fmt.Printf("float: Expression is false\n")
				ret = FAIL
				return
			}
		}

		//DOUBLE
		if res, err := buf2.BQBoolEv("T_DOUBLE_FLD[0]==19910.999"); !res || nil != err {
			if nil != err {
				fmt.Printf("double: Expression failed: %s\n", err.Error())
				ret = FAIL
				return
			} else {
				fmt.Printf("double: Expression is false\n")
				ret = FAIL
				return
			}
		}

		//STRING
		if res, err := buf2.BQBoolEv("T_STRING_FLD[0]=='HELLO STRING 2'"); !res || nil != err {
			if nil != err {
				fmt.Printf("string: Expression failed: %s\n", err.Error())
				ret = FAIL
				return
			} else {
				fmt.Printf("string: Expression is false\n")
				ret = FAIL
				return
			}
		}

		//Bool eval does not work on carray...
		tb1, _ := buf2.BGetByteArr(ubftab.T_CARRAY_FLD, 0)

		if 0 != bytes.Compare(tb1, M_b2) {
			fmt.Printf("carray: invalid array!\n")
			ret = FAIL
			return
		}

		fmt.Println("MarshalSingle tests ok...\n")

		runtime.GC()
	}
}
