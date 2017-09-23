package main

import (
	"atmi"
	"fmt"
	//"runtime"
	//http "net/http"
	//_ "net/http/pprof"
	"os"
	"strconv"
	"unsafe"
	"C"
)

const (
	SUCCEED = 0
	FAIL    = -1
)

var M_ret int
var M_ac *atmi.ATMICtx

func assertEqual(a interface{}, b interface{}, message string) {
	aa := fmt.Sprintf("%v", a)
	bb := fmt.Sprintf("%v", b)

	if aa == bb {
		return
	}
	msg2 := fmt.Sprintf("%v != %v", a, b)
	M_ac.TpLogError("%s: %s", message, msg2)
	M_ret = FAIL
}

const test_count = 100
//Binary main entry
func main() {

	M_ret = SUCCEED

	//Return to the caller
	defer func() {
		os.Exit(M_ret)
	}()

	//Test the service call with view setup
	for i := 0; i <test_count && SUCCEED == M_ret; i++ {

		ac, err := atmi.NewATMICtx()
		M_ac = ac

		if nil != err {
			fmt.Errorf("Failed to allocate cotnext!", err)
			M_ret = FAIL
			return
		}

		buf, err := ac.NewVIEW("MYVIEW1", 0)
		//		buf, err := ac.NewUBF(1024)

		if err != nil {
			ac.TpLogError("ATMI Error %s", err.Error())
			M_ret = FAIL
			return
		}

		s := strconv.Itoa(i)

		//Set one field for call
		if errB := buf.BVChg("tshort1", 0, s); nil != errB {
			ac.TpLogError("VIEW Error: %s", errB.Error())
			M_ret = FAIL
			return
		}

		if errB := buf.BVChg("tint2", 1, 123456789); nil != errB {
			ac.TpLogError("VIEW Error: %s", errB.Error())
			M_ret = FAIL
			return
		}

		if errB := buf.BVChg("tchar2", 4, 'C'); nil != errB {
			ac.TpLogError("VIEW Error: %s", errB.Error())
			M_ret = FAIL
			return
		}

		if errB := buf.BVChg("tfloat2", 0, 0.11); nil != errB {
			ac.TpLogError("VIEW Error: %s", errB.Error())
			M_ret = FAIL
			return
		}

		if errB := buf.BVChg("tdouble2", 0, 110.099); nil != errB {
			ac.TpLogError("VIEW Error: %s", errB.Error())
			M_ret = FAIL
			return
		}

		var errB1 atmi.UBFError

		if errB1 = buf.BVChg("tdouble2", 1, 110.099); nil == errB1 {
			ac.TpLogError("MUST HAVWE ERROR tdouble occ=1 does not exists, but SUCCEED!")
			M_ret = FAIL
			return
		}

		if errB1.Code() != atmi.BEINVAL {
			ac.TpLogError("Expeced error code %d but got %d", atmi.BEINVAL, errB1.Code())
			M_ret = FAIL
			return
		}

		if errB := buf.BVChg("tstring0", 2, "HELLO ENDURO"); nil != errB {
			ac.TpLogError("VIEW Error: %s", errB.Error())
			M_ret = FAIL
			return
		}

		b := []byte{0, 1, 2, 3, 4, 5}

		if errB := buf.BVChg("tcarray2", 0, b); nil != errB {
			ac.TpLogError("VIEW Error: %s", errB.Error())
			M_ret = FAIL
			return
		}

		//Call the server
		if _, err := ac.TpCall("TEST1", buf, 0); nil != err {
			ac.TpLogError("ATMI Error: %s", err.Error())
			M_ret = FAIL
			return
		}

		//Test the result buffer, type should be MYVIEW2

		ttshort1, errB := buf.BVGetInt16("ttshort1", 0, 0)
		assertEqual(ttshort1, 2233, "ttshort1 !!!")
		assertEqual(errB, nil, "ttshort1 -> errB")

		ttstring1, errB := buf.BVGetString("ttstring1", 0, 0)
		assertEqual(ttstring1, "HELLO ENDURO", "ttstring1 !!!")
		assertEqual(errB, nil, "ttstring1 -> errB")

		var itype, subtype string

		// Check the buffer type & get view name
		if _, errA := ac.TpTypes(buf.Buf, &itype, &subtype); nil != errA {
			ac.TpLogError("Failed to get return buffer type: %s", errA.Message())
			M_ret = FAIL
			return
		}

		assertEqual(itype, "VIEW", "itype -> return buffer type")
		assertEqual(subtype, "MYVIEW2", "subtype -> return buffer type")

		ac.TpTerm()
		ac.FreeATMICtx()

		//runtime.GC()
	}

	for i:=0; i < test_count && SUCCEED==M_ret; i++ {

		//Test view 2 json and vice versa...
		ac, err := atmi.NewATMICtx()
		M_ac = ac

		if nil != err {
			fmt.Errorf("Failed to allocate cotnext!", err)
			M_ret = FAIL
			return
		}

		buf, err := ac.NewVIEW("MYVIEW2", 0)

		if err != nil {
			ac.TpLogError("ATMI Error %s", err.Message())
			M_ret = FAIL
			return
		}

		//The buffer shall be initialized to default values...

		//Conver with zeros
		strj, err := buf.TpVIEWToJSON(0)

		if nil != err {
			ac.TpLogError("Failed to convert view to JSON, ATMI Error %s",
				err.Error())
			M_ret = FAIL
			return
		}

		ac.TpLogDebug("Got JSON (with nulls): [%s]", strj)

		assertEqual(strj, "{\"MYVIEW2\":{\"ttshort1\":2000,\"ttlong1\":5,"+
			"\"ttchar1\":\"\",\"ttfloat1\":1.100000,\"ttdouble1\":0,"+
			"\"ttstring1\":\"HELLO WORLDBBB\",\"ttcarray1\":"+
			"\"LQAAAAAAAAAAAA==\"}}", "VIEW2JSON with NULLs")

		if errB := buf.BVChg("ttstring1", 0, "TEST JSON"); nil != errB {
			ac.TpLogError("VIEW Error: %s", errB.Error())
			M_ret = FAIL
			return
		}

		strj, err = buf.TpVIEWToJSON(atmi.BVACCESS_NOTNULL)

		if nil != err {
			ac.TpLogError("Failed to convert view to JSON, ATMI Error %s",
				err.Error())
			M_ret = FAIL
			return
		}

		ac.TpLogDebug("Got JSON (not null): [%s]", strj)
		assertEqual(strj, "{\"MYVIEW2\":{\"ttstring1\":\"TEST JSON\"}}",
			"VIEW2JSON with NULLs")

		//Restore VIEW from this json
		v, err:=ac.TpJSONToVIEW(strj)
		if nil != err {
			ac.TpLogError("Failed to convert JSON to VIEW, ATMI Error %s",
				err.Error())
			M_ret = FAIL
			return
		}
		ttstring1, errB:=v.BVGetString("ttstring1", 0, 0)
		assertEqual(errB, nil, "ttstring1 -> errB")
		assertEqual(ttstring1, "TEST JSON", "ttstring1 !!!")

		ac.TpTerm()
		ac.FreeATMICtx()
	}

	//Test of BVNext
	for i:=0; i<test_count && SUCCEED==M_ret; i++ {
		//Test view 2 json and vice versa...
		ac, err := atmi.NewATMICtx()
		M_ac = ac

		if nil != err {
			fmt.Errorf("Failed to allocate cotnext!", err)
			M_ret = FAIL
			return
		}

		buf, err := ac.NewVIEW("MYVIEW2", 0)

		if err != nil {
			ac.TpLogError("ATMI Error %s", err.Message())
			M_ret = FAIL
			return
		}

		var state atmi.BVNextState

		ret, cname, fldtyp, maxocc, dim_size, errU :=buf.BVNext(&state, true)

		assertEqual(1, ret, "BVNext 1: have next")
		assertEqual(cname, "ttshort1", "BVNext 1: Expect: ttshort1")
		assertEqual(fldtyp, atmi.BFLD_SHORT, "BVNext 1: Expect: BFLD_SHORT")
		assertEqual(maxocc, 1, "BVNext 1: Expect: maxocc")
		var s1 C.short
		assertEqual(dim_size, unsafe.Sizeof(s1), "BVNext 1: sizeof short")
		assertEqual(errU, nil, "BVNext 1: Expect: error nil")

		ret, cname, fldtyp, maxocc, dim_size, errU =buf.BVNext(&state, false)
		assertEqual(1, ret, "BVNext 1: have next")
		assertEqual(cname, "ttlong1", "BVNext 1: Expect: ttlong1")
		assertEqual(fldtyp, atmi.BFLD_LONG, "BVNext 1: Expect: BFLD_LONG")
		assertEqual(maxocc, 1, "BVNext 1: Expect: maxocc")
		var l1 C.long
		assertEqual(dim_size, unsafe.Sizeof(l1), "BVNext 1: sizeof long")
		assertEqual(errU, nil, "BVNext 1: Expect: error nil")

		ret, cname, fldtyp, maxocc, dim_size, errU =buf.BVNext(&state, false)
		assertEqual(1, ret, "BVNext 1: have next")
		assertEqual(cname, "ttchar1", "BVNext 1: Expect: ttchar1")
		assertEqual(fldtyp, atmi.BFLD_CHAR, "BVNext 1: Expect: BFLD_CHAR")
		assertEqual(maxocc, 1, "BVNext 1: Expect: maxocc")
		var c1 C.char
		assertEqual(dim_size, unsafe.Sizeof(c1), "BVNext 1: sizeof char")
		assertEqual(errU, nil, "BVNext 1: Expect: error nil")

		ret, cname, fldtyp, maxocc, dim_size, errU =buf.BVNext(&state, false)
		assertEqual(1, ret, "BVNext 1: have next")
		assertEqual(cname, "ttfloat1", "BVNext 1: Expect: ttfloat1")
		assertEqual(fldtyp, atmi.BFLD_FLOAT, "BVNext 1: Expect: BFLD_FLOAT")
		assertEqual(maxocc, 1, "BVNext 1: Expect: maxocc")
		var f1 C.float
		assertEqual(dim_size, unsafe.Sizeof(f1), "BVNext 1: sizeof float")
		assertEqual(errU, nil, "BVNext 1: Expect: error nil")

		ret, cname, fldtyp, maxocc, dim_size, errU =buf.BVNext(&state, false)
		assertEqual(1, ret, "BVNext 1: have next")
		assertEqual(cname, "ttdouble1", "BVNext 1: Expect: ttdouble1")
		assertEqual(fldtyp, atmi.BFLD_DOUBLE, "BVNext 1: Expect: BFLD_DOUBLE")
		assertEqual(maxocc, 1, "BVNext 1: Expect: maxocc")
		var d1 C.double
		assertEqual(dim_size, unsafe.Sizeof(d1), "BVNext 1: sizeof double")
		assertEqual(errU, nil, "BVNext 1: Expect: error nil")


		ret, cname, fldtyp, maxocc, dim_size, errU =buf.BVNext(&state, false)
		assertEqual(1, ret, "BVNext 1: have next")
		assertEqual(cname, "ttstring1", "BVNext 1: Expect: ttstring1")
		assertEqual(fldtyp, atmi.BFLD_STRING, "BVNext 1: Expect: BFLD_STRING")
		assertEqual(maxocc, 1, "BVNext 1: Expect: maxocc")
		assertEqual(dim_size, 15, "BVNext 1: sizeof string")
		assertEqual(errU, nil, "BVNext 1: Expect: error nil")

		ret, cname, fldtyp, maxocc, dim_size, errU =buf.BVNext(&state, false)
		assertEqual(1, ret, "BVNext 1: have next")
		assertEqual(cname, "ttcarray1", "BVNext 1: Expect: ttcarray1")
		assertEqual(fldtyp, atmi.BFLD_CARRAY, "BVNext 1: Expect: BFLD_CARRAY")
		assertEqual(maxocc, 1, "BVNext 1: Expect: maxocc")
		assertEqual(dim_size, 10, "BVNext 1: sizeof carray")
		assertEqual(errU, nil, "BVNext 1: Expect: error nil")

		ret, cname, fldtyp, maxocc, dim_size, errU =buf.BVNext(&state, false)
		assertEqual(0, ret, "BVNext 1: have next (EOF)")
		assertEqual(errU, nil, "BVNext 1: Expect: error nil on EOF")

		ac.TpTerm()
		ac.FreeATMICtx()
	}
/*
	//Test of BVSizeof
	for i:=0; i<test_count && SUCCEED==M_ret; i++ {
		//Test view 2 json and vice versa...
		ac, err := atmi.NewATMICtx()
		M_ac = ac

		if nil != err {
			fmt.Errorf("Failed to allocate cotnext!", err)
			M_ret = FAIL
			return
		}

		buf, err := ac.NewVIEW("MYVIEW1", 0)

		if err != nil {
			ac.TpLogError("ATMI Error %s", err.Message())
			M_ret = FAIL
			return
		}

		//Call TpTypes, must equal...

		//Call BVSizeof

		// Must compare with // func (ac *ATMICtx) BVSizeof(view string) (int64, UBFError) {

		ac.TpTerm()
		ac.FreeATMICtx()
	}
*/
	//TODO: Test of:
	// func (u *TypedVIEW) BVSetOccur(cname string, occ int) UBFError {
	// func (u *TypedVIEW) BVOccur(cname string) (int, int, int, int64, int, UBFError) {

}
