package main

import (
	"atmi"
	"fmt"
	"os"
	"runtime"
	"ubftab"
)

const (
	SUCCEED = 0
	FAIL    = -1
)

var M_counter = 0

//Set the request file
func GETLOGFILE(ac *atmi.ATMICtx, svc *atmi.TPSVCINFO) {
	runtime.LockOSThread()
	ret := SUCCEED

	//Get UBF Handler
	ub, _ := ac.CastToUBF(&svc.Data)

	//Return to the caller
	defer func() {

		ac.TpLogCloseReqFile()
		if SUCCEED == ret {
			ac.TpReturn(atmi.TPSUCCESS, 0, &ub, 0)
		} else {
			ac.TpReturn(atmi.TPFAIL, 0, &ub, 0)
		}
	}()

	M_counter++
	ac.TpLog(atmi.LOG_DEBUG, "Current counter = %d", M_counter)

	ac.TpLogSetReqFile(&svc.Data, fmt.Sprintf("/tmp/08_request%d.log", M_counter), "")

	ac.TpLog(atmi.LOG_WARN, "Hello from GETLOGFILE!")

}

//TESTSVC service
func TESTSVC(ac *atmi.ATMICtx, svc *atmi.TPSVCINFO) {
	ret := SUCCEED

	//Get UBF Handler
	ub, _ := ac.CastToUBF(&svc.Data)

	ac.TpLogSetReqFile(&svc.Data, "", "")
	//Print the buffer to stdout

	ub.TpLogPrintUBF(atmi.LOG_ERROR, "Got call")

	//Resize buffer, to have some more space
	size, _ := ub.BSizeof()
	if err := ub.TpRealloc(size + 1024); err != nil {
		ac.TpLog(atmi.LOG_ERROR, "Got error: %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
		goto out
	}

	//Set some field
	if err := ub.BAdd(ubftab.T_STRING_FLD, "Hello World from Enduro/X service"); err != nil {
		ac.TpLog(atmi.LOG_ERROR, "Got error: %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
		goto out
	}
	//Set second occurance too of the T_STRING_FLD field
	if err := ub.BAdd(ubftab.T_STRING_FLD, "This is line2"); err != nil {
		ac.TpLog(atmi.LOG_ERROR, "Got error: %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
		goto out
	}

out:
	ac.TpLog(atmi.LOG_ERROR, "Returning... %d", ret)

	ac.TpLogCloseReqFile()

	ac.TpLog(atmi.LOG_DEBUG, "bank to main")

	//Return to the caller
	if SUCCEED == ret {
		ac.TpReturn(atmi.TPSUCCESS, 0, &ub, 0)
	} else {
		ac.TpReturn(atmi.TPFAIL, 0, &ub, 0)
	}
	return
}

//Server init
func Init(ac *atmi.ATMICtx) int {

	//Configure logger
	ac.TpLogConfig(atmi.LOG_FACILITY_NDRX|atmi.LOG_FACILITY_UBF|atmi.LOG_FACILITY_TP,
		-1, "file=/tmp/08_server.log ndrx=5 ubf=0 tp=5", "SRV", "")

	//Advertize TESTSVC
	if err := ac.TpAdvertise("TESTSVC", "TESTSVC", TESTSVC); err != nil {
		ac.TpLog(atmi.LOG_ERROR, fmt.Sprint(err))
		return atmi.FAIL
	}

	if err := ac.TpAdvertise("GETLOGFILE", "GETLOGFILE", GETLOGFILE); err != nil {
		ac.TpLog(atmi.LOG_ERROR, fmt.Sprint(err))
		return atmi.FAIL
	}

	return atmi.SUCCEED
}

//Server shutdown
func Uninit(ac *atmi.ATMICtx) {
	fmt.Println("Server shutting down...")
}

//Executable main entry point
func main() {
	//Have some context
	ac, err := atmi.NewATMICtx()

	if nil != err {
		fmt.Errorf("Failed to allocate cotnext!", err)
		os.Exit(atmi.FAIL)
	} else {
		//Run as server
		ac.TpRun(Init, Uninit)
		ac.FreeATMICtx()
	}
}
