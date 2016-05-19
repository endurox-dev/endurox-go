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
// So this is very simple queue use
// We will enqueue one message and dequeue it.
func main() {

	ret := SUCCEED
	defer func() { atmi.TpTerm(); os.Exit(ret) }()

	var qctl atmi.TPQCTL

	testMessage := "Hello World from queue"

	buf, err := atmi.NewString(testMessage)

	if err != nil {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
		return
	}

	//Enqueue the string
	if err := atmi.TpEnqueue("QSPACE1", "MYQ1", &qctl, buf, 0); nil != err {
		fmt.Printf("TpEnqueue() failed: ATMI Error %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
		return
	}

	fmt.Printf("Enqueue OK\n")

	//Allocate new return buffer, to ensure that this is different one..!
	buf2, err := atmi.NewString("")

	if err != nil {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
		return
	}

	if err := atmi.TpDequeue("QSPACE1", "MYQ1", &qctl, buf2, 0); nil != err {
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

	return
}
