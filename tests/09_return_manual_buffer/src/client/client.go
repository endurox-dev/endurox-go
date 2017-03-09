package main

import (
	"atmi"
	"fmt"
	//"log"
	//http "net/http"
	//_ "net/http/pprof"
	"os"
	"strconv"
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

	for i := 0; i < 100000; i++ {

		ac, err := atmi.NewATMICtx()

		if nil != err {
			fmt.Errorf("Failed to allocate cotnext!", err)
			os.Exit(atmi.FAIL)
		}

		buf, err := ac.NewUBF(1024)

		if err != nil {
			ac.TpLogError("ATMI Error %s", err.Error())
			ret = FAIL
			goto out
		}

		s := strconv.Itoa(i)

		//Set one field for call
		if errB := buf.BChg(ubftab.T_STRING_FLD, 0, s); nil != errB {
			ac.TpLogError("UBF Error: %s", errB.Error())
			ret = FAIL
			goto out
		}

		//Call the server
		if _, err := ac.TpCall("TEST1", buf, 0); nil != err {
			ac.TpLogError("ATMI Error: %s", err.Error())
			ret = FAIL
			goto out
		}

		res, errB := buf.BGetString(ubftab.T_STRING_3_FLD, 0)

		if nil!=errB {
			ac.TpLogError("UBF Error: %s", errB.Error())
			ret = FAIL
			goto out
		}

		if res!=s {
			ac.TpLogError("Sent %s, but got [%s]")
		}

		ac.TpTerm()
		ac.FreeATMICtx()
	}

out:
	//Close the ATMI session

	os.Exit(ret)
}
