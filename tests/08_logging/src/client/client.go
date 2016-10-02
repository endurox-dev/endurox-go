package main

import (
	"atmi"
	"fmt"
	"os"
	"runtime"
	"time"
	"ubftab"
)

const (
	SUCCEED = 0
	FAIL    = -1
)

//Test the threading
func testGo1() {
	runtime.LockOSThread()

	atmi.TpLogConfig(atmi.LOG_FACILITY_TP_THREAD,
		-1, "file=/tmp/08_th1.log ndrx=5 ubf=5 tp=5", "TH1", "")

	time.Sleep(1000 * time.Millisecond)

	atmi.TpLog(atmi.LOG_ALWAYS, "Hello from TH1")

	time.Sleep(1000 * time.Millisecond)

}

func testGo2() {
	runtime.LockOSThread()

	atmi.TpLogConfig(atmi.LOG_FACILITY_TP_THREAD,
		-1, "file=/tmp/08_th2.log ndrx=5 ubf=5 tp=5", "TH2", "")

	time.Sleep(1000 * time.Millisecond)

	atmi.TpLog(atmi.LOG_ALWAYS, "Hello from TH2")

	time.Sleep(1000 * time.Millisecond)

}

//Binary main entry
func main() {

	runtime.LockOSThread()

	ret := SUCCEED

	buf, err := atmi.NewUBF(16024)

	//Return to the caller (kind of destructor..)
	defer func() {
		atmi.TpTerm()
		os.Exit(ret)
	}()

	if err != nil {
		atmi.TpLog(atmi.LOG_ERROR, "ATMI Error %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
		return
	}

	ptr1 := []byte{0, 1, 2, 3, 4, 5, 6, 8, 9}
	ptr2 := []byte{2, 3, 4, 5, 6, 7, 8, 9, 10}

	atmi.TpLogConfig(atmi.LOG_FACILITY_NDRX|atmi.LOG_FACILITY_UBF|atmi.LOG_FACILITY_TP,
		-1, "file=/tmp/08_client_process.log ndrx=5 ubf=5 tp=5", "SUPERUSER", "")

	atmi.TpLog(atmi.LOG_DEBUG, "Hello From GO this is %d and %d!", 5, 6)

	atmi.TpLogDump(atmi.LOG_WARN, "Doing buffer HEX dump", ptr1, len(ptr1))

	atmi.TpLogDumpDiff(atmi.LOG_DEBUG, "Test HEX dump diff", ptr1, ptr2, len(ptr1))

	//Open request logging

	atmi.TpLogSetReqFile_Direct("/tmp/08_single-request-1.log")

	atmi.TpLog(atmi.LOG_ERROR, "Hello from request")

	atmi.TpLogCloseReqFile()

	atmi.TpLogConfig(atmi.LOG_FACILITY_TP_THREAD,
		-1, "file=/tmp/08_mainth.log ndrx=5 ubf=5 tp=5", "MTH", "")

	atmi.TpLog(atmi.LOG_ERROR, "Hello from main th")
	go testGo1()

	go testGo2()

	time.Sleep(4000 * time.Millisecond)

	//Set one field for call

	for i := 0; i < 100; i++ {

		atmi.TpLogDelBufReqFile(buf.GetBuf())
		if err := atmi.TpLogSetReqFile(buf.GetBuf(), "", "GETLOGFILE"); nil != err {
			atmi.TpLog(atmi.LOG_ERROR, "ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			return
		}

		atmi.TpLog(atmi.LOG_ERROR, "New string [%s]", fmt.Sprintf("HELLO FROM CLIENT %d", i))
		if err := buf.BChg(ubftab.T_CARRAY_FLD, i, fmt.Sprintf("HELLO FROM CLIENT %d abc", i)); nil != err {
			atmi.TpLog(atmi.LOG_ERROR, "ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			return
		}

		tmp, _ := buf.BGet(ubftab.T_CARRAY_FLD, i)

		atmi.TpLog(atmi.LOG_ERROR, "New string from buffer [%s]", tmp)

		buf.TpLogPrintUBF(atmi.LOG_ERROR, "Buffer before call")

		//Call the server
		if _, err := atmi.TpCall("TESTSVC", buf, 0); nil != err {
			atmi.TpLog(atmi.LOG_ERROR, "ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			return
		}

		buf.TpLogPrintUBF(atmi.LOG_ERROR, "Buffer after call")

		atmi.TpLogCloseReqFile()
	}

}
