package main

/*
#include <signal.h>
*/
import "C"

import (
	"atmi"
	"bytes"
	"fmt"
	"os"
	"runtime"
	"sync"
	"ubftab"
)

const (
	SUCCEED = 0
	FAIL    = -1
)

var M_ret chan int
var M_wg sync.WaitGroup

//Binary main entry
func async_main() {

	b1 := []byte{0, 1, 2, 3}

	b2 := []byte{4, 3, 2, 1, 0}

	C.signal(11, nil)

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
			fmt.Errorf("ERROR ! Failed to allocate cotnext!", err)
			ret = FAIL
			return
		}

		buf, err := ac.NewUBF(1024)

		if nil != err {
			ac.TpLogError("ERROR ! Failed to allocate buffer 1 %s!", err.Error())
			ret = FAIL
			return
		}

		buf2, err := ac.NewUBF(1024)

		if nil != err {
			ac.TpLogError("ERROR ! Failed to allocate buffer 2 %s!", err.Error())
			ret = FAIL
			return
		}

		err = loadbufferdata(buf)

		if nil != err {
			ac.TpLogError("ERROR ! Failed to load UBF data! %s", err.Error())
			ret = FAIL
			return
		}

		//Convert to string

		jstr, err := ac.TpExport(buf, 0)

		if nil != err {
			ac.TpLogError("ERROR ! Failed to TpExport! %s", err.Error())
			ret = FAIL
			return
		}

		ac.TpLogDebug("Converted JSON: [%s]", jstr)

		err = ac.TpImport(jstr, buf2, 0)

		if nil != err {
			ac.TpLogError("ERROR ! Failed to TpImport! %s", err.Error())
			ret = FAIL
			return
		}

		//CHAR
		if res, err := buf2.BQBoolEv("T_CHAR_FLD=='A' && T_CHAR_FLD[1]=='B'"); !res || nil != err {
			if nil != err {
				ac.TpLogError("char: Expression failed: %s\n", err.Error())
				ret = FAIL
				return
			} else {
				ac.TpLogError("char: Expression is false")
				ret = FAIL
				return
			}
		}

		//SHORT
		if res, err := buf2.BQBoolEv("T_SHORT_FLD==32000 && T_SHORT_FLD[1]==32001"); !res || nil != err {
			if nil != err {
				ac.TpLogError("ERROR ! short: Expression failed: %s\n", err.Error())
				ret = FAIL
				return
			} else {
				ac.TpLogError("ERROR ! short: Expression is false")
				ret = FAIL
				return
			}
		}

		//LONG
		if res, err := buf2.BQBoolEv("T_LONG_FLD==9999999101 && T_LONG_FLD[1]==9999999102"); !res || nil != err {
			if nil != err {
				ac.TpLogError("ERROR ! long: Expression failed: %s", err.Error())
				ret = FAIL
				return
			} else {
				ac.TpLogError("ERROR ! long: Expression is false")
				ret = FAIL
				return
			}
		}

		//FLOAT
		if res, err := buf2.BQBoolEv("T_FLOAT_FLD==9.11 && T_FLOAT_FLD[1]==9.22"); !res || nil != err {
			if nil != err {
				ac.TpLogError("ERROR ! float: Expression failed: %s", err.Error())
				ret = FAIL
				return
			} else {
				ac.TpLogError("ERROR ! float: Expression is false")
				ret = FAIL
				return
			}
		}

		//DOUBLE
		if res, err := buf2.BQBoolEv("T_DOUBLE_FLD==999999910.888 && T_DOUBLE_FLD[1]==999999910.999"); !res || nil != err {
			if nil != err {
				ac.TpLogError("ERROR ! double: Expression failed: %s", err.Error())
				ret = FAIL
				return
			} else {
				ac.TpLogError("ERROR ! double: Expression is false")
				ret = FAIL
				return
			}
		}

		//STRING
		if res, err := buf2.BQBoolEv("T_STRING_FLD=='HELLO STRING 1' && T_STRING_FLD[1]=='HELLO STRING 2'"); !res || nil != err {
			if nil != err {
				ac.TpLogError("ERROR ! string: Expression failed: %s", err.Error())
				ret = FAIL
				return
			} else {
				ac.TpLogError("ERROR ! string: Expression is false")
				ret = FAIL
				return
			}
		}

		//Bool eval does not work on carray...
		tb1, _ := buf2.BGetByteArr(ubftab.T_CARRAY_FLD, 0)
		tb2, _ := buf2.BGetByteArr(ubftab.T_CARRAY_FLD, 1)

		if 0 != bytes.Compare(tb1, b1) ||
			0 != bytes.Compare(tb2, b2) {
			ac.TpLogError("ERROR ! carray: invalid array!\n")
			ret = FAIL
			return
		}

		fmt.Println("TpImport / TpExport tests OK\n")

		runtime.GC()

	}

}

func main() {

	//Set some test data...
	M_b1 = []byte{0, 1, 2, 3}

	M_b2 = []byte{4, 3, 2, 1, 0}

	// you can also add these one at
	// a time if you need to

	// Have some core dumps...
	C.signal(11, nil)

	M_ret = make(chan int, 20)
	M_wg.Add(10)
	for i := 0; i < 10; i++ {
		go async_main()
	}
	M_wg.Wait()

	i := 0
	for ret := range M_ret {
		fmt.Println(ret)
		i++
		if ret == FAIL {
			os.Exit(-1)
		}
		//For some reason the for loop does not terminate by it self..
		if i >= 10 {
			break
		}
	}

	os.Exit(0)
}
