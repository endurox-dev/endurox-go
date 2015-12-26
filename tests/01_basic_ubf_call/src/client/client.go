package main

import (
	"atmi"
	"fmt"
        "os"
        "ubftab"
)

const (
        SUCCEED = 0
        FAIL = -1
)

//Binary main entry
func main() {

        ret:=SUCCEED

        buf, err := atmi.NewUBF(1024)

        if err != nil {
                fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message());
                ret = FAIL
                goto out
        }

        //Set one field for call
        if err:=buf.BChg(ubftab.T_CARRAY_FLD, 0, "HELLO FROM CLIENT"); nil!=err {
                fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message());
                ret = FAIL
                goto out
        }

        //Call the server
        if _, err:=atmi.TpCall("TESTSVC", buf, 0); nil!=err {
                fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message());
                ret = FAIL
                goto out
        }
        
        //Print the output buffer
        buf.BPrint()
        
out:
        //Close the ATMI session
        atmi.TpTerm()
        os.Exit(ret)
}
