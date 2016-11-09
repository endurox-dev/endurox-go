package main

/*
#include <signal.h>
*/
import "C"

import (
	"atmi"
	"bytes"
	"fmt"
	"sync"
	"runtime"
	"ubftab"
	"os"
)

const (
	SUCCEED = 0
	FAIL    = -1
)

//Struct to marshal
type TestStruct struct {
	CharTest      byte   `ubf:"T_CHAR_FLD"`
	CharArrayTest []byte `ubf:"T_CHAR_FLD"`

	ShortTest      int16   `ubf:"T_SHORT_FLD"`
	ShortArrayTest []int16 `ubf:"T_SHORT_FLD"`

	LongTest      int64   `ubf:"T_LONG_FLD"`
	LongArrayTest []int64 `ubf:"T_LONG_FLD""`

	Float32Test      float32   `ubf:"T_FLOAT_FLD"`
	Float32ArrayTest []float32 `ubf:"T_FLOAT_FLD"`

	Float64Test      float64   `ubf:"T_DOUBLE_FLD"`
	Float64ArrayTest []float64 `ubf:"T_DOUBLE_FLD"`

	StringTest      string   `ubf:"T_STRING_FLD"`
	StringArrayTest []string `ubf:"T_STRING_FLD"`

	//By default this goes as array
	CarrayTest [][]byte `ubf:"T_CARRAY_FLD"`
}

var M_ret chan int
var M_wg sync.WaitGroup

//Binary main entry
func async_main() {

        C.signal(11, nil);

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
		b1 := []byte{0, 1, 2, 3}

		b2 := []byte{4, 3, 2, 1, 0}

		////////////////////////////////////////////////////////////////////////
		// Set the test data
		////////////////////////////////////////////////////////////////////////
		if err != nil {
			fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			return
		}

		if err := buf.BChg(ubftab.T_CHAR_FLD, 0, 65); nil != err {
			fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			return
		}

		if err := buf.BChg(ubftab.T_CHAR_FLD, 1, 66); nil != err {
			fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			return
		}

		if err := buf.BChg(ubftab.T_SHORT_FLD, 0, 32000); nil != err {
			fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			return
		}

		if err := buf.BChg(ubftab.T_SHORT_FLD, 1, 32001); nil != err {
			fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			return
		}

		if err := buf.BChg(ubftab.T_LONG_FLD, 0, 9999999101); nil != err {
			fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			return
		}

		if err := buf.BChg(ubftab.T_LONG_FLD, 1, 9999999102); nil != err {
			fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			return
		}

		if err := buf.BChg(ubftab.T_FLOAT_FLD, 0, 9.11); nil != err {
			fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			return
		}

		if err := buf.BChg(ubftab.T_FLOAT_FLD, 1, 9.22); nil != err {
			fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			return
		}

		if err := buf.BChg(ubftab.T_DOUBLE_FLD, 0, 999999910.888); nil != err {
			fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			return
		}

		if err := buf.BChg(ubftab.T_DOUBLE_FLD, 1, 999999910.999); nil != err {
			fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			return
		}

		if err := buf.BChg(ubftab.T_STRING_FLD, 0, "HELLO STRING 1"); nil != err {
			fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			return
		}

		if err := buf.BChg(ubftab.T_STRING_FLD, 1, "HELLO STRING 2"); nil != err {
			fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			return
		}

		if err := buf.BChg(ubftab.T_CARRAY_FLD, 0, b1); nil != err {
			fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			return
		}

		if err := buf.BChg(ubftab.T_CARRAY_FLD, 1, b2); nil != err {
			fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			return
		}
		////////////////////////////////////////////////////////////////////////
		// Do the Unmarshal & test the values
		////////////////////////////////////////////////////////////////////////
		buf.Unmarshal(&s)

		fmt.Printf("Got the struct [%v]\n", s)

		//CHAR
		if s.CharTest != 65 {
			fmt.Printf("s.CharTest invalid value!\n")
			ret = FAIL
			return
		}

		if len(s.CharArrayTest) != 2 || s.CharArrayTest[0] != 65 ||
			s.CharArrayTest[1] != 66 {
			fmt.Printf("s.CharArrayTest invalid value!\n")
			ret = FAIL
			return
		}

		//SHORT
		if s.ShortTest != 32000 {
			fmt.Printf("s.ShortTest invalid value!\n")
			ret = FAIL
			return
		}

		if len(s.ShortArrayTest) != 2 || s.ShortArrayTest[0] != 32000 ||
			s.ShortArrayTest[1] != 32001 {
			fmt.Printf("s.ShortArrayTest invalid value!\n")
			ret = FAIL
			return
		}

		//LONG
		if s.LongTest != 9999999101 {
			fmt.Printf("s.LongTest invalid value!\n")
			ret = FAIL
			return
		}

		if len(s.LongArrayTest) != 2 || s.LongArrayTest[0] != 9999999101 ||
			s.LongArrayTest[1] != 9999999102 {
			fmt.Printf("s.LongArrayTest invalid value!\n")
			ret = FAIL
			return
		}

		//FLOAT
		if s.Float32Test-9.11 > 0.0001 {
			fmt.Printf("s.Float32Test invalid value!\n")
			ret = FAIL
			return
		}

		if len(s.Float32ArrayTest) != 2 || s.Float32ArrayTest[0]-9.11 > 0.0001 ||
			s.Float32ArrayTest[1]-9.22 > 0.0001 {
			fmt.Printf("s.Float32ArrayTest invalid value!\n")
			ret = FAIL
			return
		}

		//FLOAT64
		if s.Float64Test-999999910.888 > 0.0001 {
			fmt.Printf("s.Float64Test invalid value!\n")
			ret = FAIL
			return
		}

		if len(s.Float64ArrayTest) != 2 || s.Float64ArrayTest[0]-999999910.888 > 0.0001 ||
			s.Float64ArrayTest[1]-999999910.999 > 0.0001 {
			fmt.Printf("s.Float64ArrayTest invalid value!\n")
			ret = FAIL
			return
		}

		//STRING
		if s.StringTest != "HELLO STRING 1" {
			fmt.Printf("s.StringTest invalid value!\n")
			ret = FAIL
			return
		}

		if len(s.StringArrayTest) != 2 || s.StringArrayTest[0] != "HELLO STRING 1" ||
			s.StringArrayTest[1] != "HELLO STRING 2" {
			fmt.Printf("s.StringArrayTest invalid value!\n")
			ret = FAIL
			return
		}

		//CARRAY
		if len(s.CarrayTest) != 2 || 0 != bytes.Compare(s.CarrayTest[0], b1) ||
			0 != bytes.Compare(s.CarrayTest[1], b2) {
			fmt.Printf("s.CarrayTest invalid value!\n")
			ret = FAIL
			return
		}

		fmt.Println("Unmarshal tests ok...\n")

		////////////////////////////////////////////////////////////////////////
		// DO the marshal tests...
		// Test the stuff with boolean expressions...
		////////////////////////////////////////////////////////////////////////
		//Reset the buffer & generate it again...
		if err := buf2.Marshal(&s); err != nil {
			fmt.Printf("Failed to marshal [%s]\n", err.Error())
			ret = FAIL
			return
		}
		buf2.BPrint()

		//CHAR
		if res, err := buf2.BQBoolEv("T_CHAR_FLD=='A' && T_CHAR_FLD[1]=='B'"); !res || nil != err {
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
		if res, err := buf2.BQBoolEv("T_SHORT_FLD==32000 && T_SHORT_FLD[1]==32001"); !res || nil != err {
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
		if res, err := buf2.BQBoolEv("T_LONG_FLD==9999999101 && T_LONG_FLD[1]==9999999102"); !res || nil != err {
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
		if res, err := buf2.BQBoolEv("T_FLOAT_FLD==9.11 && T_FLOAT_FLD[1]==9.22"); !res || nil != err {
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
		if res, err := buf2.BQBoolEv("T_DOUBLE_FLD==999999910.888 && T_DOUBLE_FLD[1]==999999910.999"); !res || nil != err {
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
		if res, err := buf2.BQBoolEv("T_STRING_FLD=='HELLO STRING 1' && T_STRING_FLD[1]=='HELLO STRING 2'"); !res || nil != err {
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
		tb2, _ := buf2.BGetByteArr(ubftab.T_CARRAY_FLD, 1)

		if 0 != bytes.Compare(tb1, b1) ||
			0 != bytes.Compare(tb2, b2) {
			fmt.Printf("carray: invalid array!\n")
			ret = FAIL
			return
		}

		fmt.Println("Marshal tests ok...\n")

		runtime.GC()
	}

}

func main() {

        // you can also add these one at
        // a time if you need to
        M_ret = make(chan int, 10)
        M_wg.Add(10)
        // Have some core dumps...
        C.signal(11, nil);

        for i := 0; i < 10; i++ {
                go async_main()
        }

        M_wg.Wait()

	i:=0
        for ret := range M_ret {
                fmt.Println(ret)
		i++
                if ret == FAIL {
                        os.Exit(-1)
                }
		//For some reason the for loop does not terminate by it self..
		if i>= 10 {
			break
		}
        }

        os.Exit(0)
}
