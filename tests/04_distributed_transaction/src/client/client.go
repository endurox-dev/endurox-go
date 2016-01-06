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

//Binary main entry
func main() {

	ret := SUCCEED

        // Allocate some buffer
	buf, err := atmi.NewUBF(1024)

	if err != nil {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
		goto out
        }

	//Open The XA sub-sysitem
	if err := atmi.TpOpen(); err != nil {
		fmt.Printf("ATMI Error: [%s]\n", err.Message())
		ret = FAIL
		goto out
	}

	//Begin transaction, timeout 60 sec
	if err := atmi.TpBegin(60, 0); err != nil {
		fmt.Printf("ATMI Error: [%s]\n", err.Message())
		ret = FAIL
		goto out
	}

	//Set customer name field
	if err := buf.BChg(ubftab.T_CUSTOMER_NAME, 0, "John"); nil != err {
		fmt.Printf("UBF Error %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
		goto out
	}

	//Set city field
	if err := buf.BChg(ubftab.T_CITY, 0, "Riga"); nil != err {
		fmt.Printf("UBF Error %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
		goto out
	}

	//Call the server
        //Will use TRANSUSPEND as we run on the same RMID
        //On one RMID there can be only one resource client active
        //Or otherwise we could use dynamic registration
	if _, err := atmi.TpCall("MKCUST", buf, atmi.TPTRANSUSPEND); nil != err {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
		goto out
	}

	//Print the output buffer
	buf.BPrint()
	
        //Print the customer id
	if cust_id, err:=buf.BGetInt64(ubftab.T_CUSTOMER_ID, 0); nil!=err {
		fmt.Printf("UBF Error %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
		goto out
	} else {
		fmt.Printf("Got customer id: %d\n", cust_id)
	}
	
	//Commit the transaction
	if err := atmi.TpCommit(0); err != nil {
		fmt.Printf("Got error: [%s]\n", err.Message())
		ret = FAIL
		goto out
	}


out:

	//Abort transaction, if we failed.
	if SUCCEED != ret {
		if err := atmi.TpAbort(0); err != nil {
			fmt.Printf("Got error: [%s]\n", err.Message())
			ret = FAIL
		}	
	}
	
	//Close the XA sub-system
	atmi.TpClose()
	
	//Close the ATMI session
	atmi.TpTerm()
	os.Exit(ret)
}
