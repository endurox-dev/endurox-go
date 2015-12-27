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

        bytes := []byte{9,8,7,6,5,4,3,2,1,0}
        buf, err := atmi.NewCarray(bytes)

        if err != nil {
                fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message());
                ret = FAIL
                goto out
        }

        fmt.Printf("Sending: [%v]\n", buf.GetBytes())

        //Call the server
        if _, err:=atmi.TpCall("TESTSVC", buf, 0); nil!=err {
                fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message());
                ret = FAIL
                goto out
        }
        
        //Print the output buffer
        fmt.Printf("Got response: [%v]\n", buf.GetBytes())
        
out:
        //Close the ATMI session
        atmi.TpTerm()
        os.Exit(ret)
}
