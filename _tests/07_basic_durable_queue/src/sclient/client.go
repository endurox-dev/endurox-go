package main

/*
#include <signal.h>
*/
import "C"

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
// So this is very simple queue use
// We will enqueue one message and dequeue it.
func main() {

	ret := SUCCEED
	defer func() { os.Exit(ret) }()

        // Have some core dumps...
        C.signal(11, nil);

	//Have some loop for memory leak checks...
	for i := 0; i < 1000; i++ {

		var ac *atmi.ATMICtx
		var err atmi.ATMIError
		//Allocate context
		ac, err = atmi.NewATMICtx()
		if nil != err {
			fmt.Errorf("Failed to allocate cotnext!", err)
			ret = FAIL
			return
		}

		var qctl atmi.TPQCTL

		testMessage := "Hello World from queue"

		buf, err := ac.NewString(testMessage)

		if err != nil {
			fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			return
		}

		//Enqueue the string
		if err := ac.TpEnqueue("QSPACE1", "MYQ1", &qctl, buf, 0); nil != err {
			fmt.Printf("TpEnqueue() failed: ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			return
		}

		fmt.Printf("Enqueue OK\n")

		//Allocate new return buffer, to ensure that this is different one..!
		buf2, err := ac.NewString("")

		if err != nil {
			fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			return
		}

		if err := ac.TpDequeue("QSPACE1", "MYQ1", &qctl, buf2, 0); nil != err {
			fmt.Printf("TpDequeue() failed: ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			return
		}

		//Print the output buffer
		fmt.Printf("Dequeued message: [%s]\n", buf2.GetString())

		if buf2.GetString() != testMessage {
			fmt.Printf("ERROR ! Enqueued [%s] but dequeued [%s]",
				testMessage, buf2.GetString())
			ret = FAIL
			return
		}
		fmt.Printf("Message machged ok!\n")

		runtime.GC()
	}

	return
}
