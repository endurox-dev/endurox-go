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
	C.signal(11, nil)

	//Have some loop for memory leak checks...
	for i := 0; i < 100; i++ {

		var ac *atmi.ATMICtx
		var err atmi.ATMIError
		//Allocate context
		ac, err = atmi.NewATMICtx()
		if nil != err {
			fmt.Errorf("Failed to allocate cotnext!", err)
			ret = FAIL
			return
		}

		if err := ac.TpOpen(); nil != err {
			fmt.Printf("TpOpen() failed: ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			return
		}

		//////////////////////////////////////////////////////////////////////////
		// Test with out transaction
		//////////////////////////////////////////////////////////////////////////
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

		//////////////////////////////////////////////////////////////////////////
		// Transaction mode test...
		//////////////////////////////////////////////////////////////////////////

		for j := 0; j < 2; j++ {

			//Check transactional functionality
			if err := ac.TpBegin(60, 0); nil != err {
				fmt.Printf("TpBegin() failed: ATMI Error %d:[%s]\n", err.Code(), err.Message())
				ret = FAIL
				return
			}

			if ac.TpGetLev() != 1 {
				fmt.Printf("TpGetLev() failed: not in transaction: %d (1)\n", ac.TpGetLev())
				ret = FAIL
				return
			}

			for n := 0; n < 10; n++ {
				testMessage = fmt.Sprintf("Hello World from queue %d %d %d", n, j, i)

				buf, err = ac.NewString(testMessage)

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
			}

			fmt.Printf("Enqueue OK (TX) \n")

			var tranid atmi.TPTRANID

			//Suspend transaction
			if err := ac.TpSuspend(&tranid, 0); nil != err {
				fmt.Printf("TpSuspend() failed: ATMI Error %d:[%s]\n", err.Code(), err.Message())
				ret = FAIL
				return
			}

			if ac.TpGetLev() != 0 {
				fmt.Printf("TpGetLev() failed: must not be in transaction but is: %d (2)\n", ac.TpGetLev())
				ret = FAIL
				return
			}

			//Resume transaction
			if err := ac.TpResume(&tranid, 0); nil != err {
				fmt.Printf("TpResume() failed: ATMI Error %d:[%s]\n", err.Code(), err.Message())
				ret = FAIL
				return
			}

			if ac.TpGetLev() != 1 {
				fmt.Printf("TpGetLev() failed: must be in transaction but is not: %d (3)\n", ac.TpGetLev())
				ret = FAIL
				return
			}

			//Suspend with base64
			btid, err := ac.TpSuspendString(0)
			if nil != err {
				fmt.Printf("TpSuspendString() failed: ATMI Error %d:[%s]\n", err.Code(), err.Message())
				ret = FAIL
				return
			}

			if ac.TpGetLev() != 0 {
				fmt.Printf("TpGetLev() failed: must not be in transaction but is: %d (4)\n", ac.TpGetLev())
				ret = FAIL
				return
			}

			fmt.Printf("Got trans id [%s]\n", btid)

			//Resume transaction
			if err := ac.TpResumeString(btid, 0); nil != err {
				fmt.Printf("TpResumeString() failed: ATMI Error %d:[%s]\n", err.Code(), err.Message())
				ret = FAIL
				return
			}

			//Resume transaction
			err = ac.TpResumeString("aGVsbG8K", 0)

			if nil == err {
				fmt.Printf("TpResumeString() shall fail but didn't\n")
				ret = FAIL
				return
			}

			if atmi.TPEINVAL != err.Code() {
				fmt.Printf("TpResumeString() Invalid error code, expected %d got %d\n",
					atmi.TPEINVAL, err.Code())
				ret = FAIL
				return
			}

			if ac.TpGetLev() != 1 {
				fmt.Printf("TpGetLev() failed: must be in transaction but is not: %d (5)\n", ac.TpGetLev())
				ret = FAIL
				return
			}

			if j == 0 {

				if err := ac.TpAbort(0); nil != err {
					fmt.Printf("TpAbort() failed: ATMI Error %d:[%s]\n", err.Code(), err.Message())
					ret = FAIL
					return
				}

				err := ac.TpDequeue("QSPACE1", "MYQ1", &qctl, buf2, 0)

				if nil == err {
					fmt.Printf("TpDequeue() did not fail (msgs are aborted... %d %d)\n",
						i, j)
					ret = FAIL
					return
				}

				if atmi.TPEDIAGNOSTIC != err.Code() {
					fmt.Printf("TpDequeue() Invalid error code, expected %d got %d: %s\n",
						atmi.TPEDIAGNOSTIC, err.Code(), err.Message())
					ret = FAIL
					return
				}

				//Check that it is no msg..

				if atmi.QMENOMSG != qctl.Diagnostic {
					fmt.Printf("TpDequeue() expected diag %d got %d",
						atmi.QMENOMSG, qctl.Diagnostic)
					ret = FAIL
					return
				}
			} else {

				if err := ac.TpCommit(0); nil != err {
					fmt.Printf("TpCommit() (2) failed: ATMI Error %d:[%s]\n", err.Code(), err.Message())
					ret = FAIL
					return
				}

				for n := 0; n < 10; n++ {

					testMessage = fmt.Sprintf("Hello World from queue %d %d %d", n, j, i)

					if err := ac.TpDequeue("QSPACE1", "MYQ1", &qctl, buf2, 0); nil != err {
						fmt.Printf("TpDequeue() failed: ATMI Error %d:[%s]\n", err.Code(), err.Message())
						ret = FAIL
						return
					}

					//Print the output buffer
					fmt.Printf("Dequeued message: [%s]\n", buf2.GetString())

					if buf2.GetString() != testMessage {
						fmt.Printf("ERROR ! Enqueued [%s] but dequeued [%s]\n",
							testMessage, buf2.GetString())
						ret = FAIL
						return
					}
					fmt.Printf("Message machged ok (2)!\n")
				}

			}

		}

		if err := ac.TpClose(); nil != err {
			fmt.Printf("TpClose() failed: ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			return
		}

		runtime.GC()
	}

	return
}
