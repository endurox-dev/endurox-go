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

	ac, err = atmi.NewATMICtx()

	if nil != err {
		fmt.Errorf("Failed to allocate cotnext!", err)
		ret = FAIL
		return
	}

	buf, err := ac.NewString("Hello World")

	if err != nil {
		ac.TpLogError("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
		return
	}

	//Call the server
	if _, err := ac.TpCall("TESTSVC", buf, 0); nil != err {
		ac.TpLogError("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
		return
	}

	//Print the output buffer
	fmt.Printf("Got response: [%s]\n", buf.GetString())
}
