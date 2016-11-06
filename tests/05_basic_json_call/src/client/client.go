package main

import (
	"atmi"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
)

const (
	SUCCEED = 0
	FAIL    = -1
)

type Message struct {
	T_CHAR_FLD   byte
	T_SHORT_FLD  int16
	T_LONG_FLD   int64
	T_FLOAT_FLD  float32
	T_DOUBLE_FLD float64
	T_STRING_FLD string
	T_CARRAY_FLD []byte
}

//Binary main entry
func main() {

	ret := SUCCEED

	defer func() {
		os.Exit(ret)
	}()
	//Have some loop for memory leak checks...
	for i := 0; i < 10000; i++ {
		var ac *atmi.ATMICtx
		var err atmi.ATMIError
		//Allocate context
		ac, err = atmi.NewATMICtx()
		if nil != err {
			fmt.Errorf("Failed to allocate cotnext!", err)
			ret = FAIL
			return
		}

		m := Message{65, 100, 1294706395881547000, 66.77, 11111111222.77, "Hello Wolrd", []byte{0, 1, 2, 3}}

		b, _ := json.Marshal(m)

		bb := []byte{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}

		fmt.Printf("Got JSON [%s]\n", string(b))

		buf, err := ac.NewJSON(b)

		if err != nil {
			fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			return
		}

		//Call the server
		if _, err := ac.TpCall("TESTSVC", buf, 0); nil != err {
			fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
			ret = FAIL
			return
		}

		json.Unmarshal(buf.GetJSON(), &m)

		if m.T_STRING_FLD != "Hello World from Enduro/X service" {
			fmt.Printf("Invalid message recieved: [%s]\n", m.T_STRING_FLD)
			ret = FAIL
			goto out
		}

		if 0 != bytes.Compare(bb, m.T_CARRAY_FLD) {
			fmt.Printf("Invalid c array received...")
			ret = FAIL
			return
		}

		fmt.Println(m)
		//Close the ATMI session
		runtime.GC()
	}

out:

	os.Exit(ret)
}
