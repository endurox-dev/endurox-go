package main

/*
#include <signal.h>
*/
import "C"

import (
	"atmi"
	"fmt"
	"os"
	"runtime"
//	"time"
	"ubftab"
)

const (
	SUCCEED = 0
	FAIL    = -1
)

//Test the threading
func testGo1() {
	//Allocate context
	ac, err := atmi.NewATMICtx()
	if nil != err {
		fmt.Errorf("Failed to allocate cotnext!", err)
		return
	}
	ac.TpLogConfig(atmi.LOG_FACILITY_TP_THREAD,
		-1, "file=/tmp/08_th1.log ndrx=5 ubf=5 tp=5", "TH1", "")

	//time.Sleep(1000 * time.Millisecond)

	ac.TpLog(atmi.LOG_ALWAYS, "Hello from TH1")

	//time.Sleep(1000 * time.Millisecond)

}

func testGo2() {

	//Allocate context
	ac, err := atmi.NewATMICtx()
	if nil != err {
		fmt.Errorf("Failed to allocate cotnext!", err)
		return
	}

	ac.TpLogConfig(atmi.LOG_FACILITY_TP_THREAD,
		-1, "file=/tmp/08_th2.log ndrx=5 ubf=5 tp=5", "TH2", "")

	//time.Sleep(1000 * time.Millisecond)

	ac.TpLog(atmi.LOG_ALWAYS, "Hello from TH2")

	//time.Sleep(1000 * time.Millisecond)

}

//Binary main entry
func main() {

	ret := SUCCEED
	//Return to the caller (kind of destructor..)
	defer func() {
		os.Exit(ret)
	}()

	// Have some core dumps...
	C.signal(11, nil)

	for i := 0; i < 100; i++ {
		var ac *atmi.ATMICtx
		var err atmi.ATMIError
		//bytes := []byte{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}

		//Allocate context
		ac, err = atmi.NewATMICtx()
		if nil != err {
			fmt.Errorf("Failed to allocate cotnext!", err)
			ret = FAIL
			return
		}

		buf, err := ac.NewUBF(16024)

		if err != nil {
			ac.TpLog(atmi.LOG_ERROR, "ATMI Error %d:[%s]\n",
				err.Code(), err.Message())
			ret = FAIL
			return
		}

		ptr1 := []byte{0, 1, 2, 3, 4, 5, 6, 8, 9}
		ptr2 := []byte{2, 3, 4, 5, 6, 7, 8, 9, 10}

		ac.TpLogConfig(atmi.LOG_FACILITY_NDRX|atmi.LOG_FACILITY_UBF|atmi.LOG_FACILITY_TP,
			-1, "file=/tmp/08_client_process.log ndrx=5 ubf=5 tp=5", "SUPERUSER", "")

		ac.TpLog(atmi.LOG_DEBUG, "Hello From GO this is %d and %d!", 5, 6)

		ac.TpLogDump(atmi.LOG_WARN, "Doing buffer HEX dump", ptr1, len(ptr1))

		ac.TpLogDumpDiff(atmi.LOG_DEBUG, "Test HEX dump diff", ptr1, ptr2, len(ptr1))

		//Open request logging

		ac.TpLogSetReqFile_Direct("/tmp/08_single-request-1.log")

		ac.TpLog(atmi.LOG_ERROR, "Hello from request")

		ac.TpLogCloseReqFile()

		ac.TpLogConfig(atmi.LOG_FACILITY_TP_THREAD,
			-1, "file=/tmp/08_mainth.log ndrx=5 ubf=5 tp=5", "MTH", "")

		ac.TpLog(atmi.LOG_ERROR, "Hello from main th")
		go testGo1()

		go testGo2()

		//time.Sleep(4000 * time.Millisecond)

		//Set one field for call

		for i := 0; i < 10; i++ {

			ac.TpLogDelBufReqFile(buf.GetBuf())
			if err := ac.TpLogSetReqFile(buf.GetBuf(), "", "GETLOGFILE"); nil != err {
				ac.TpLog(atmi.LOG_ERROR, "ATMI Error %d:[%s]\n", err.Code(), err.Message())
				ret = FAIL
				return
			}

			ac.TpLog(atmi.LOG_ERROR, "New string [%s]", fmt.Sprintf("HELLO FROM CLIENT %d", i))
			if err := buf.BChg(ubftab.T_CARRAY_FLD, i, fmt.Sprintf("HELLO FROM CLIENT %d abc", i)); nil != err {
				ac.TpLog(atmi.LOG_ERROR, "ATMI Error %d:[%s]\n", err.Code(), err.Message())
				ret = FAIL
				return
			}

			tmp, _ := buf.BGet(ubftab.T_CARRAY_FLD, i)

			ac.TpLog(atmi.LOG_ERROR, "New string from buffer [%s]", tmp)

			buf.TpLogPrintUBF(atmi.LOG_ERROR, "Buffer before call")

			//Call the server
			if _, err := ac.TpCall("TESTSVC", buf, 0); nil != err {
				ac.TpLog(atmi.LOG_ERROR, "ATMI Error %d:[%s]\n", err.Code(), err.Message())
				ret = FAIL
				return
			}

			buf.TpLogPrintUBF(atmi.LOG_ERROR, "Buffer after call")

			ac.TpLogCloseReqFile()
		}

		runtime.GC()
	}
}
