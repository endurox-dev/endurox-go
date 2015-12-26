package main

import (
	"atmi"
	"fmt"
	"ubftab"
)

const (
	SUCCEED = 0
	FAIL    = -1
)

//TESTSVC service
func TESTSVC(svc *atmi.TPSVCINFO) {

	ret := SUCCEED

	//Get UBF Handler
	ub, _ := atmi.CastToUBF(&svc.Data)

	//Print the buffer to stdout
        fmt.Println("Incoming request:")
	ub.BPrint()

        //Resize buffer, to have some more space
        if err :=ub.TpRealloc(1024); err!=nil {
		fmt.Printf("Got error: %d:[%s]\n", err.Code(), err.Message())
                ret = FAIL
                goto out
        }

	//Set some file
	if err := ub.BChg(ubftab.T_STRING_FLD, 0, "Hello World from Enduro/X service"); err != nil {
		fmt.Printf("Got error: %d:[%s]\n", err.Code(), err.Message())
                ret = FAIL
                goto out
	}
	if err:=ub.BChg(ubftab.T_STRING_FLD, 1, "This is line2"); err!=nil {
		fmt.Printf("Got error: %d:[%s]\n", err.Code(), err.Message())
                ret = FAIL
                goto out
        }

out:
        //Return to the caller
        if SUCCEED==ret {
        	atmi.TpReturn(atmi.TPSUCCESS, 0, &ub, 0)
        } else {
        	atmi.TpReturn(atmi.TPFAIL, 0, &ub, 0)
        }
	return
}

//Server init
func Init() int {

	//Advertize TESTSVC
	if err := atmi.TpAdvertise("TESTSVC", "TESTSVC", TESTSVC); err != nil {
		fmt.Println(err)
		return atmi.FAIL
	}

	return atmi.SUCCEED
}

//Server shutdown
func Uninit() {
	fmt.Println("Server shutting down...")
}

//Executable main entry point
func main() {

	//Run as server
	atmi.TpRun(Init, Uninit)
}

