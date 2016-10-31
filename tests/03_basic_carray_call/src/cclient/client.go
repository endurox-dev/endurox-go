package main

import (
	"atmi"
	"fmt"
	"os"
	"runtime"
)

const (
	SUCCEED = 0
	FAIL    = -1
)

//Binary main entry
func main() {

	ret := SUCCEED
	//Return to the caller (kind of destructor..)
	defer func() {
		os.Exit(ret)
	}()

	for i := 0; i < 10000; i++ {
		var ac *atmi.ATMICtx
		var err atmi.ATMIError
		bytes := []byte{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}

		//Allocate context
		ac, err = atmi.NewATMICtx()
		if nil != err {
			fmt.Errorf("Failed to allocate cotnext!", err)
			ret = FAIL
			return
		}

		buf, err := ac.NewCarray(bytes)

		if err != nil {
			fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			return
		}

		fmt.Printf("Sending: [%v]\n", buf.GetBytes())

		//Call the server
		if _, err := ac.TpCall("TESTSVC", buf, 0); nil != err {
			fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			return
		}

		//Print the output buffer
		fmt.Printf("Got response: [%v]\n", buf.GetBytes())

		runtime.GC()
	}
}
