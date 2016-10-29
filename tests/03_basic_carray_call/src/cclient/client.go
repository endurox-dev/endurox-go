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

//Binary main entry
func main() {

	ret := SUCCEED
	var ac *atmi.ATMICtx
	var err atmi.ATMIError
	//Return to the caller (kind of destructor..)
	defer func() {
		if nil != ac {
			ac.TpTerm()
			ac.FreeATMICtx() // Kill the context
		}
		os.Exit(ret)
	}()

	bytes := []byte{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}
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

}
