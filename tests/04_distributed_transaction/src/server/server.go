package main

import (
	"atmi"
	"database/sql"
	"fmt"
	_ "oci8"
	"ubftab"
)

const (
	SUCCEED = 0
	FAIL    = -1
)

//Connection to DB
var M_db *sql.DB

//MKCUST service
func MKCUST(svc *atmi.TPSVCINFO) {

	ret := SUCCEED

        var id, city, cust_name string
        var uerr        atmi.UBFError

	//Get UBF Handler
	ub, _ := atmi.CastToUBF(&svc.Data)

	//Print the buffer to stdout
	fmt.Println("Incoming request:")
	ub.BPrint()

	//Ensure that we have space for answer...
	if err := ub.TpRealloc(64); err != nil {
		fmt.Printf("Got error: %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
                goto out
	}

	// Get the next customer id
	if rows, err := M_db.Query("select nvl(max(CUSTOMER_ID),0)+1  from customers"); err != nil {
		fmt.Println("Query error: ", err)
	} else {
	//	defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(&id); err != nil {
				ret = FAIL
                                goto out
			}
			fmt.Printf("Got new ID: [%d]", id)
		}
	}

	//Return ID back to caller.
	if err := ub.BChg(ubftab.T_CUSTOMER_ID, 0, id); err != nil {
		fmt.Printf("Failed to set T_CUSTOMER_ID: %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
                goto out
	}

	//Read the other fields
	cust_name, uerr = ub.BGetString(ubftab.T_CUSTOMER_NAME, 0)
	if uerr != nil {
		fmt.Printf("Failed to get T_CUSTOMER_NAME: %d:[%s]\n", uerr.Code(), uerr.Message())
		ret = FAIL
                goto out
	}

	city, uerr = ub.BGetString(ubftab.T_CITY, 0)
	if uerr != nil {
		fmt.Printf("Failed to get T_CITY: %d:[%s]\n", uerr.Code(), uerr.Message())
		ret = FAIL
                goto out
	}

	//Now insert the record
	if _, err := M_db.Exec("INSERT INTO customers (customer_id, customer_name, city) " +
		" VALUES (" + id + ", '" + cust_name + "', '" + city + "')"); err != nil {
		fmt.Printf("Failed to create customer: %s\n", err)
		ret = FAIL
                goto out
	}

out:
	//Return to the caller
	if SUCCEED == ret {
		atmi.TpReturn(atmi.TPSUCCESS, 0, &ub, 0)
	} else {
		atmi.TpReturn(atmi.TPFAIL, 0, &ub, 0)
	}
	return
}

//Server init
func Init() int {

	if err := atmi.TpOpen(); err != nil {
		fmt.Println(err)
		return atmi.FAIL
	}

	//Advertize MKCUST
	if err := atmi.TpAdvertise("MKCUST", "MKCUST", MKCUST); err != nil {
		fmt.Println(err)
		return atmi.FAIL
	}

	//Connect to XA driver (empty conn string...) & get the SQL handler.
	if db, err := sql.Open("oci8", "dummy_user:dummy_pass@localhost:1111/SID?enable_xa=YES"); err != nil {
//	if db, err := sql.Open("oci8", "endurotest:endurotest1@localhost:1521/ROCKY"); err != nil {
		fmt.Printf("Failed to get SQL handler: %s\n", err)
		return atmi.FAIL
	} else {
		M_db = db
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
