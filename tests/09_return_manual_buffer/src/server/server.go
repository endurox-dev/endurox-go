package main

import (
	"atmi"
	"fmt"
	"os"
	"ubftab"
)

const (
	SUCCEED = 0
	FAIL    = -1
)

//TEST1 service
//Read the value of the T_STRING_FLD, alloc new buffer and set T_STRING2_FLD
//Forward to TEST2 service
func TEST1(ac *atmi.ATMICtx, svc *atmi.TPSVCINFO) {

	ret := SUCCEED

	//Get UBF Handler
	ub, _ := ac.CastToUBF(&svc.Data)

	//Return to the caller
	defer func() {

		ac.TpLogCloseReqFile()
		if SUCCEED == ret {
			ac.TpForward("TEST2", ub, 0)
		} else {
			ac.TpReturn(atmi.TPFAIL, 0, ub, 0)
		}
	}()

	//Set some field
	f, errB := ub.BGetString(ubftab.T_STRING_FLD, 0)

	if  errB != nil {
		ac.TpLogError("Bget() Got error: %s", errB.Error())
		ret = FAIL
		return
	}

	//Alloc new buffer
	ub, errA := ac.NewUBF(1024)

	if errA != nil {
		ac.TpLogError("ATMI Error: %s", errA.Error())
		ret=FAIL
		return
	}

	//Set one field for call
	if errB = ub.BChg(ubftab.T_STRING_2_FLD, 0, f); nil != errB {
		fmt.Printf("UBF Error: %s", errB.Error())
		ret = FAIL
		return
	}
}

//TEST2 service
//Read the value of the T_STRING2_FLD, alloc new buffer and set T_STRING3_FLD
//And return
func TEST2(ac *atmi.ATMICtx, svc *atmi.TPSVCINFO) {

	ret := SUCCEED

	//Get UBF Handler
	ub, _ := ac.CastToUBF(&svc.Data)

	//Return to the caller
	defer func() {

		ac.TpLogCloseReqFile()
		if SUCCEED == ret {
			ac.TpReturn(atmi.TPSUCCESS, 0, ub, 0)
		} else {
			ac.TpReturn(atmi.TPFAIL, 0, ub, 0)
		}
	}()

	//Set some field
	f, errB := ub.BGetString(ubftab.T_STRING_2_FLD, 0)

	if  errB != nil {
		ac.TpLogError("Bget() Got error: %s", errB.Error())
		ret = FAIL
		return
	}

	//Alloc new buffer
	ub, errA := ac.NewUBF(1024)

	if errA != nil {
		ac.TpLogError("ATMI Error: %s", errA.Error())
		ret=FAIL
		return
	}

	//Set one field for call
	if errB = ub.BChg(ubftab.T_STRING_3_FLD, 0, f); nil != errB {
		fmt.Printf("UBF Error: %s", errB.Error())
		ret = FAIL
		return
	}
}


//Server init
func Init(ac *atmi.ATMICtx) int {

	//Advertize TEST1
	if err := ac.TpAdvertise("TEST1", "TEST1", TEST1); err != nil {
		ac.TpLogError("TpAdvertise fail: %s", err.Error())
		return atmi.FAIL
	}

	//Advertize TEST2
	if err := ac.TpAdvertise("TEST2", "TEST2", TEST2); err != nil {
		ac.TpLogError("TpAdvertise fail: %s", err.Error())
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
	}
}
