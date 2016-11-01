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

	var ac *atmi.ATMICtx
	var err atmi.ATMIError

	ret := SUCCEED
	//Return to the caller (kind of destructor..)
	defer func() {
		os.Exit(ret)
	}()

	var cust_id_first int64

	//Allocate context
	ac, err = atmi.NewATMICtx()
	if nil != err {
		fmt.Errorf("Failed to allocate cotnext!", err)
		ret = FAIL
		return
	}

	// Allocate some buffer
	buf, err := ac.NewUBF(1024)

	if err != nil {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
		goto out
	}

	//Open The XA sub-sysitem
	if err := ac.TpOpen(); err != nil {
		fmt.Printf("ATMI Error: [%s]\n", err.Message())
		ret = FAIL
		goto out
	}

	//Begin transaction, timeout 60 sec
	if err := ac.TpBegin(60, 0); err != nil {
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

	////////////////////////////////////////////////////////////////////////
	//First call!
	////////////////////////////////////////////////////////////////////////
	//Call the server
	//Will use TRANSUSPEND as we run on the same RMID
	//On one RMID there can be only one resource client active
	//Or otherwise we could use dynamic registration
	if _, err := ac.TpCall("MKCUST", buf, atmi.TPTRANSUSPEND); nil != err {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
		goto out
	}
	////////////////////////////////////////////////////////////////////////

	//Print the output buffer
	buf.BPrint()

	//Get the first call customer id.
	//After txn abort, must be the same at next call.

	cust_id_first, _ = buf.BGetInt64(ubftab.T_CUSTOMER_ID, 0)

	//Do the abort & call again, the ID must be the same if tran works.

	//Abort the transaction
	if err := ac.TpAbort(0); err != nil {
		fmt.Printf("Got error: [%s]\n", err.Message())
		ret = FAIL
		goto out
	}

	//Begin transaction, timeout 60 sec
	if err := ac.TpBegin(60, 0); err != nil {
		fmt.Printf("ATMI Error: [%s]\n", err.Message())
		ret = FAIL
		goto out
	}

	////////////////////////////////////////////////////////////////////////
	//Second call!
	////////////////////////////////////////////////////////////////////////
	//Call the server
	//Will use TRANSUSPEND as we run on the same RMID
	//On one RMID there can be only one resource client active
	//Or otherwise we could use dynamic registration
	if _, err := ac.TpCall("MKCUST", buf, atmi.TPTRANSUSPEND); nil != err {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
		goto out
	}

	//Print the customer id
	if cust_id, err := buf.BGetInt64(ubftab.T_CUSTOMER_ID, 0); nil != err {
		fmt.Printf("UBF Error %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
		goto out
	} else {
		fmt.Printf("Got customer id: %d\n", cust_id)

		if cust_id_first != cust_id {
			fmt.Printf("XA transaction fail, first call id: %d, second %d\n",
				cust_id_first, cust_id)
			ret = FAIL
			goto out
		}
	}

	//Commit the transaction
	if err := ac.TpCommit(0); err != nil {
		fmt.Printf("Got error: [%s]\n", err.Message())
		ret = FAIL
		goto out
	}

out:

	//Abort transaction, if we failed.
	if SUCCEED != ret {
		if err := ac.TpAbort(0); err != nil {
			fmt.Printf("Got error: [%s]\n", err.Message())
		}
	}

	//Close the XA sub-system
	ac.TpClose()

	//Close the ATMI session
	ac.TpTerm()

	os.Exit(ret)
}
