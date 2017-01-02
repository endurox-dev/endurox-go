package main

// #include <signal.h>
import "C"

import (
	"atmi"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"sync"
	"ubftab"
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

var M_ret chan int
var M_wg sync.WaitGroup

//Binary main entry
func async_main() {

	ret := SUCCEED

	defer func() {
		M_wg.Done()
		M_ret <- ret
	}()
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
			return
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

}

//Do some tests with UBF convert to/from JSON
func test_buffer_convert() int {

	var ac *atmi.ATMICtx
	var err atmi.ATMIError
	//Allocate context
	ac, err = atmi.NewATMICtx()
	if nil != err {
		fmt.Errorf("Failed to allocate contex [%s]t!\n", err)
		return atmi.FAIL

	}
	//Create UBF buffer
	u, err := ac.NewUBF(1024)
	if err != nil {
		fmt.Errorf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		return atmi.FAIL
	}

	//Set some fields to buffer
	u.BChg(ubftab.T_CHAR_FLD, 0, "A")
	u.BChg(ubftab.T_DOUBLE_2_FLD, 0, 12.9999)
	u.BChg(ubftab.T_STRING_FLD, 0, "Hello World")
	u.BChg(ubftab.T_STRING_FLD, 2, "Yes...")

	u.TpLogPrintUBF(atmi.LOG_WARN, "Original UBF")

	//Conver to JSON...
	str, err := u.TpUBFToJSON()

	if err != nil {
		fmt.Printf("TpUBFToJSON() fail - ATMI Error %d:[%s]\n",
			err.Code(), err.Message())
		return atmi.FAIL
	}

	fmt.Printf("Got json string from UBF: [%s]\n", str)

	//Convert the JSON back to buffer...
	u.BProj([]int{0}) //Reset the the buffer (0 - bad field id...)

	//Check that reset was ok
	if u.BPres(ubftab.T_CHAR_FLD, 0) {
		fmt.Printf("ubftab.T_CHAR_FLD must not exist!\n")
		return atmi.FAIL
	}

	u.TpLogPrintUBF(atmi.LOG_WARN, "Cleared UBF")

	//Convert build UBF from json
	err = u.TpJSONToUBF(str)
	if err != nil {
		fmt.Printf("TpUBFToJSON() fail - ATMI Error %d:[%s]\n",
			err.Code(), err.Message())
		return atmi.FAIL
	}

	//Check the buffer for values...
	if res, err := u.BQBoolEv("T_CHAR_FLD=='A' && T_STRING_FLD[2]=='Yes...'"); !res || nil != err {
		if nil != err {
			fmt.Printf("long: Expression failed: %s\n", err.Error())
			return atmi.FAIL
		} else {
			fmt.Printf("long: Expression is false\n")
			return atmi.FAIL
		}
	}

	u.TpLogPrintUBF(atmi.LOG_WARN, "Restored UBF")

	return atmi.SUCCEED
}

func main() {

	// you can also add these one at
	// a time if you need to
	M_ret = make(chan int, 10)
	M_wg.Add(10)
	// Have some core dumps...
	C.signal(11, nil)

	for i := 0; i < 10; i++ {
		go async_main()
	}

	M_wg.Wait()

	i := 0
	for ret := range M_ret {
		fmt.Println(ret)
		if ret == FAIL {
			os.Exit(FAIL)
		}
		i++
		//For some reason the for loop does not terminate by it self..
		if i >= 10 {
			break
		}
	}

	ret2 := test_buffer_convert()

	if atmi.SUCCEED != ret2 {
		os.Exit(FAIL)
	}

	os.Exit(SUCCEED)
}
