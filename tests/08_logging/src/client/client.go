package main

import (
	"atmi"
	"os"
)

const (
	SUCCEED = 0
	FAIL    = -1
)

//Binary main entry
func main() {

	ret := SUCCEED

	/*
		buf, err := atmi.NewUBF(1024)

		if err != nil {
			fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			goto out
		}
	*/

	ptr1 := []byte{0, 1, 2, 3, 4, 5, 6, 8, 9}
	ptr2 := []byte{2, 3, 4, 5, 6, 7, 8, 9, 10}

	atmi.TpLog(atmi.LOG_DEBUG, "Hello From GO this is %d and %d!", 5, 6)

	atmi.TpLogDump(atmi.LOG_WARN, "Doing buffer HEX dump", ptr1, len(ptr1))

	atmi.TpLogDumpDiff(atmi.LOG_DEBUG, "Test HEX dump diff", ptr1, ptr2, len(ptr1))

	/*
		//Set one field for call
		if err := buf.BChg(ubftab.T_CARRAY_FLD, 0, "HELLO FROM CLIENT"); nil != err {
			fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			goto out
		}

		//Call the server
		if _, err := atmi.TpCall("TESTSVC", buf, 0); nil != err {
			fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			goto out
		}

		//Print the output buffer
		buf.BPrint()
	*/
	//out:
	//Close the ATMI session
	atmi.TpTerm()
	os.Exit(ret)
}
