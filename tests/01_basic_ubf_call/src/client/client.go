package main

import (
	"atmi"
	"fmt"
	//"log"
	//http "net/http"
	//_ "net/http/pprof"
	"os"
	"runtime"
	"ubftab"
)

const (
	SUCCEED = 0
	FAIL    = -1
)

//Binary main entry
func main() {

	ret := SUCCEED

	// Run profiler
	// go func() {
	//	log.Println(http.ListenAndServe("localhost:6060", nil))
	//}()

	for i := 0; i < 10000; i++ {

		ac, err := atmi.NewATMICtx()

		if nil != err {
			fmt.Errorf("Failed to allocate cotnext!", err)
			os.Exit(atmi.FAIL)
		}

		buf, err := ac.NewUBF(1024)

		if err != nil {
			fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			goto out
		}

		//Set one field for call
		if err := buf.BChg(ubftab.T_CARRAY_FLD, 0, "HELLO FROM CLIENT"); nil != err {
			fmt.Printf("UBF Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			goto out
		}
		//test that we can run zero lenght byte arrays
		var nulbuf []byte
		if err := buf.BChg(ubftab.T_CARRAY_2_FLD, 0, nulbuf); nil != err {
			fmt.Printf("UBF Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			goto out
		}

		//Get zero len array
		nulout, err := buf.BGetByteArr(ubftab.T_CARRAY_2_FLD, 0)

		if nil != err {
			fmt.Printf("Failed to get 0len array %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			goto out
		}

		if  len(nulout) != 0 {
			fmt.Printf("nulout not 0: %d\n", len(nulout))
			ret = FAIL
			goto out
		}

		//Call the server
		if _, err := ac.TpCall("TESTSVC", buf, 0); nil != err {
			fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			goto out
		}

		//Print the output buffer
		//buf.BPrint()
		buf.TpLogPrintUBF(atmi.LOG_DEBUG, "Got response")
		ac.TpTerm()
		ac.FreeATMICtx()
		runtime.GC()
	}

	//Run some bigger message tests
	fmt.Printf("Message size: %d\n", atmi.ATMIMsgSizeMax())
	if atmi.ATMIMsgSizeMax() > 68000 {
		for i := 0; i < 10000; i++ {
			testdata :=make([]byte, 1024*1024)
			ac, err := atmi.NewATMICtx()

			if nil != err {
				fmt.Errorf("Failed to allocate cotnext!", err)
				os.Exit(atmi.FAIL)
			}

			//2MB buffer
			buf, err := ac.NewUBF(1024*1024 + 1024)

			if err != nil {
				fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
				ret = FAIL
				goto out
			}

			for i := 0; i < len(testdata); i++ {
				testdata[i] = byte((i + 1) % 255)
			}

			//Set one field for call
			if err := buf.BChg(ubftab.T_CARRAY_FLD, 0, testdata); nil != err {
				fmt.Printf("! ATMI Error %d:[%s]\n", err.Code(), err.Message())
				ret = FAIL
				goto out
			}

			//Call the server
			if _, err := ac.TpCall("BIGMSG", buf, 0); nil != err {
				//RPI and others, ignore big fails when message size unavailable
				if err.Code() == atmi.TPEINVAL {
					continue
				}

				fmt.Printf("ATMI Error BIGMSG %d:[%s]\n", err.Code(), err.Message())

				ret = FAIL
				goto out
			}

			testdata, err = buf.BGetByteArr(ubftab.T_CARRAY_FLD, 0)

			if nil != err {
				fmt.Printf("! Failed to get rsp: ATMI Error %d:[%s]\n", err.Code(), err.Message())
				ret = FAIL
				goto out
			}

			for i := 0; i < len(testdata); i++ {
				if testdata[i] != byte((i+2)%255) {
					ac.TpLogError("TESTERROR: Error at index %d expected %d got: %d",
						i, (i+2)%255, testdata[i])
					ret = FAIL
					goto out
				}
			}

		}
	}

out:
	//Close the ATMI session

	os.Exit(ret)
}
