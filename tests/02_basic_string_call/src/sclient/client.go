package main

import (
	"atmi"
	"fmt"
        "os"
)

const (
        SUCCEED = 0
        FAIL = -1
)

//Binary main entry
func main() {

        ret:=SUCCEED

        buf, err := atmi.NewString("Hello World")

        if err != nil {
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
        fmt.Printf("Got response: [%s]\n", buf.GetString())
        
out:
        //Close the ATMI session
        atmi.TpTerm()
        os.Exit(ret)
}
