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
        somebytes := []byte{0,1,2,3,4,5,6,7,8,9,10}

	//Get Typed String Handler
	c, _ := atmi.CastToCarray(&svc.Data)

	//Print the buffer to stdout
        fmt.Printf("Incoming request: [%v]\n", c.GetBytes())

        //Resize buffer, to have some more space
        if err :=c.TpRealloc(128); err!=nil {
		fmt.Printf("Got error: %d:[%s]\n", err.Code(), err.Message())
                ret = FAIL
                goto out
        }

	//Send some bytes back...
        c.SetBytes(somebytes);

out:
        //Return to the caller
        if SUCCEED==ret {
        	atmi.TpReturn(atmi.TPSUCCESS, 0, &c, 0)
        } else {
        	atmi.TpReturn(atmi.TPFAIL, 0, &c, 0)
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

