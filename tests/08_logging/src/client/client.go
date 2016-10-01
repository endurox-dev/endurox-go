package main

import (
	"atmi"
	"os"
	"runtime"
	"time"
)

const (
	SUCCEED = 0
	FAIL    = -1
)

//Test the threading
func testGo1() {
	runtime.LockOSThread()

	atmi.TpLogConfig(atmi.LOG_FACILITY_TP_THREAD,
		-1, "file=./th1.out ndrx=5 ubf=5 tp=5", "TH1", "")

	time.Sleep(1000 * time.Millisecond)

	atmi.TpLog(atmi.LOG_ALWAYS, "Hello from TH1")

	time.Sleep(1000 * time.Millisecond)

}

func testGo2() {
	runtime.LockOSThread()

	atmi.TpLogConfig(atmi.LOG_FACILITY_TP_THREAD,
		-1, "file=./th2.out ndrx=5 ubf=5 tp=5", "TH2", "")

	time.Sleep(1000 * time.Millisecond)

	atmi.TpLog(atmi.LOG_ALWAYS, "Hello from TH2")

	time.Sleep(1000 * time.Millisecond)

}

//Binary main entry
func main() {

	runtime.LockOSThread()

	ret := SUCCEED

	/*
		buf, err := atmi.NewUBF(1024)

		if err != nil {
			fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			goto out
		}
	*/

	ptr1 := []byte{0, 1, 2, 3, 4, 5, 6, 8, 9}
	ptr2 := []byte{2, 3, 4, 5, 6, 7, 8, 9, 10}

	atmi.TpLogConfig(atmi.LOG_FACILITY_NDRX|atmi.LOG_FACILITY_UBF|atmi.LOG_FACILITY_TP,
		-1, "file=./test.out ndrx=5 ubf=5 tp=5", "SUPERUSER", "")

	atmi.TpLog(atmi.LOG_DEBUG, "Hello From GO this is %d and %d!", 5, 6)

	atmi.TpLogDump(atmi.LOG_WARN, "Doing buffer HEX dump", ptr1, len(ptr1))

	atmi.TpLogDumpDiff(atmi.LOG_DEBUG, "Test HEX dump diff", ptr1, ptr2, len(ptr1))

	//Open request logging

	atmi.TpLogSetReqFile_Direct("./request-log-1")

	atmi.TpLog(atmi.LOG_ERROR, "Hello from request")

	atmi.TpLogCloseReqFile()

	atmi.TpLogConfig(atmi.LOG_FACILITY_TP_THREAD,
		-1, "file=./mainth.out ndrx=5 ubf=5 tp=5", "MTH", "")

	atmi.TpLog(atmi.LOG_ERROR, "Hello from main th")
	go testGo1()

	go testGo2()

	time.Sleep(4000 * time.Millisecond)
	/*
		//Set one field for call
		if err := buf.BChg(ubftab.T_CARRAY_FLD, 0, "HELLO FROM CLIENT"); nil != err {
			fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			goto out
		}

		//Call the server
		if _, err := atmi.TpCall("TESTSVC", buf, 0); nil != err {
			fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			goto out
		}

		//Print the output buffer
		buf.BPrint()
	*/
	//out:
	//Close the ATMI session
	atmi.TpTerm()
	os.Exit(ret)
}
