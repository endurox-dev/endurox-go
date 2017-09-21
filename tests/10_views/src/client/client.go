package main

import (
	"atmi"
	"fmt"
	//"log"
	//http "net/http"
	//_ "net/http/pprof"
	"os"
	//"strconv"
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

		buf, err := ac.NewVIEW("MYVIEW1", 0)
//		buf, err := ac.NewUBF(1024)

		if err != nil {
			ac.TpLogError("ATMI Error %s", err.Error())
			ret = FAIL
			goto out
		}

/*
		s := strconv.Itoa(i)

		//Set one field for call
		if errB := buf.BVChg("tshort1", 0, s); nil != errB {
			ac.TpLogError("VIEW Error: %s", errB.Error())
			ret = FAIL
			goto out
		}

		if errB := buf.BVChg("tint2", 1, 123456789); nil != errB {
			ac.TpLogError("VIEW Error: %s", errB.Error())
			ret = FAIL
			goto out
		}

		if errB := buf.BVChg("tchar2", 4, 'A'); nil != errB {
			ac.TpLogError("VIEW Error: %s", errB.Error())
			ret = FAIL
			goto out
		}

		if errB := buf.BVChg("tfloat2", 0, 0.11); nil != errB {
			ac.TpLogError("VIEW Error: %s", errB.Error())
			ret = FAIL
			goto out
		}

		if errB := buf.BVChg("tdouble2", 0, 110.099); nil != errB {
			ac.TpLogError("VIEW Error: %s", errB.Error())
			ret = FAIL
			goto out
		}

		var errB1 atmi.UBFError

		if errB1 = buf.BVChg("tdouble2", 1, 110.099); nil == errB1 {
			ac.TpLogError("MUST HAVWE ERROR tdouble occ=1 does not exists, but SUCCEED!")
			ret = FAIL
			goto out
		}

		if errB1.Code() != atmi.BEINVAL {
			ac.TpLogError("Expeced error code %d but got %d", atmi.BEINVAL, errB1.Code())
			ret = FAIL
			goto out
		}

		if errB := buf.BVChg("tstring0", 2, "HELLO ENDURO"); nil != errB {
			ac.TpLogError("VIEW Error: %s", errB.Error())
			ret = FAIL
			goto out
		}

		b := []byte{0, 1, 2, 3, 4, 5}

		if errB := buf.BVChg("tcarray2", 0, b); nil != errB {
			ac.TpLogError("VIEW Error: %s", errB.Error())
			ret = FAIL
			goto out
		}
*/

		//Call the server
		if _, err := ac.TpCall("TEST1", buf, 0); nil != err {
			ac.TpLogError("ATMI Error: %s", err.Error())
			ret = FAIL
			goto out
		}

		ac.TpTerm()
		ac.FreeATMICtx()
	}

out:
	//Close the ATMI session

	os.Exit(ret)
}
