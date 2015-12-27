package main

import (
	"atmi"
	"fmt"
)

const (
	SUCCEED = 0
	FAIL    = -1
)

//TESTSVC service
func TESTSVC(svc *atmi.TPSVCINFO) {

	ret := SUCCEED

	//Get Typed String Handler
	s, _ := atmi.CastToString(&svc.Data)

	//Print the buffer to stdout
        fmt.Printf("Incoming request: [%s]\n", s.GetString())

        //Resize buffer, to have some more space
        if err :=s.TpRealloc(128); err!=nil {
		fmt.Printf("Got error: %d:[%s]\n", err.Code(), err.Message())
                ret = FAIL
                goto out
        }

	//Send string back
        s.SetString("Hello From TESTSVC. This string is bit longer than receved in req");

out:
        //Return to the caller
        if SUCCEED==ret {
        	atmi.TpReturn(atmi.TPSUCCESS, 0, &s, 0)
        } else {
        	atmi.TpReturn(atmi.TPFAIL, 0, &s, 0)
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

