package main

import (
	"atmi"
	"fmt"
	"os"
)

const (
	SUCCEED = 0
	FAIL    = -1
)

var M_counter int = 0
var M_ret int
var M_ac *atmi.ATMICtx

func assertEqual(a interface{}, b interface{}, message string) {
        aa:= fmt.Sprintf("%v", a)
        bb:= fmt.Sprintf("%v", b)

        if aa == bb {
                return
        }
        msg2:= fmt.Sprintf("%v != %v", a, b)
        M_ac.TpLogError("%s: %s", message, msg2)
        M_ret = FAIL
}


//TEST1 service
//Read the value of the T_STRING_FLD, alloc new buffer and set T_STRING2_FLD
//Forward to TEST2 service
func TEST1(ac *atmi.ATMICtx, svc *atmi.TPSVCINFO) {

	M_ret = SUCCEED
	M_ac = ac

	//Get UBF Handler
	v, err := ac.CastToVIEW(&svc.Data)
	if err != nil {
		ac.TpLogError("Failed to cast to view: %s", err.Error())
		M_ret = FAIL
		return //drop the message...
	}

	//Return to the caller
	defer func() {
		M_counter++
		if SUCCEED == M_ret {
			ac.TpReturn(atmi.TPSUCCESS, 0, v, 0)
		} else {
			ac.TpReturn(atmi.TPFAIL, 0, v, 0)
		}
	}()

	////////////////////////////////////////////////////////////////////////
	//Test the values received
	////////////////////////////////////////////////////////////////////////
	tshort1, errV := v.BVGetInt16("tshort1", 0, 0)
	assertEqual(tshort1, M_counter, "tshort1")
	assertEqual(errV, nil, "tshort1 -> errV")

	tint2, errV := v.BVGetInt("tint2", 1, 0)
	assertEqual(tint2, 123456789, "tint2")
	assertEqual(errV, nil, "tint2 -> errV")

	tchar2, errV := v.BVGetString("tchar2", 4, 0)
	assertEqual(tchar2, "C", "tchar2")
	assertEqual(errV, nil, "tchar2 -> errV")

	tfloat2, errV := v.BVGetFloat32("tfloat2", 0, 0)
	assertEqual(tfloat2, 0.11, "tfloat2")
	assertEqual(errV, nil, "tfloat2 -> errV")

	tdouble2, errV := v.BVGetFloat64("tdouble2", 0, 0)
	assertEqual(tdouble2, 110.099, "tdouble2")
	assertEqual(errV, nil, "tdouble2 -> errV")

	tstring0, errV := v.BVGetString("tstring0", 2, 0)
	assertEqual(tstring0, "HELLO ENDURO", "tstring0")
	assertEqual(errV, nil, "tstring0 -> errV")

	b := []byte{0, 1, 2, 3, 4, 5}

	tcarray2, errV := v.BVGetByteArr("tcarray2", 0, 0)
	for i:=0; i<5; i++ {
		assertEqual(tcarray2[i], b[i], "tcarray2")
	}
	//assertEqual(tcarray2, b, "tcarray2")
	assertEqual(errV, nil, "tcarray2 -> errV")
	////////////////////////////////////////////////////////////////////////
	//Test BVACCESS_NOTNULL functionality...
	////////////////////////////////////////////////////////////////////////

	tshort2, errV := v.BVGetInt16("tshort2", 1, 0)
	assertEqual(tshort2, 2001, "tshort2")
	assertEqual(errV, nil, "tshort2 -> errV")

	tshort2_2, errV := v.BVGetInt16("tshort2", 1, atmi.BVACCESS_NOTNULL)
	assertEqual(tshort2_2, 0, "tshort2")
	assertEqual(errV.Code(), atmi.BNOTPRES, "tshort2_2 -> must not be present"+
			" with atmi.BVACCESS_NOTNULL")

	_, errV = v.BVGetInt16("tshortX", 1, atmi.BVACCESS_NOTNULL)
	assertEqual(errV.Code(), atmi.BNOCNAME, "tshortX")


	v, err = ac.NewVIEW("MYVIEW2", 0);
	if err != nil {
		ac.TpLogError("Failed to cast to view: %s", err.Error())
		M_ret = FAIL
		return
	}

	if errB := v.BVChg("ttshort1", 0, 2233); nil != errB {
		ac.TpLogError("VIEW Error: %s", errB.Error())
		M_ret = FAIL
		return
	}

	if errB := v.BVChg("ttstring1", 0, "HELLO ENDURO"); nil != errB {
		ac.TpLogError("VIEW Error: %s", errB.Error())
		M_ret = FAIL
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
		fmt.Errorf("Failed to allocate cotnext: %s!", err.Message())
		os.Exit(atmi.FAIL)
	} else {
		//Run as server
		ac.TpRun(Init, Uninit)
	}
}
