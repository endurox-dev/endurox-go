package main

import (
	"atmi"
	"bytes"
	"fmt"
	"os"
	"ubftab"
)

const (
	SUCCEED = 0
	FAIL    = -1
)

//Struct to marshal
type TestStruct struct {
	CharTest      byte   `ubf:"T_CHAR_FLD"`
	CharArrayTest []byte `ubf:"T_CHAR_FLD"`

	ShortTest      int16   `ubf:"T_SHORT_FLD"`
	ShortArrayTest []int16 `ubf:"T_SHORT_FLD"`

	LongTest      int64   `ubf:"T_LONG_FLD"`
	LongArrayTest []int64 `ubf:"T_LONG_FLD""`

	Float32Test      float32   `ubf:"T_FLOAT_FLD"`
	Float32ArrayTest []float32 `ubf:"T_FLOAT_FLD"`

	Float64Test      float64   `ubf:"T_DOUBLE_FLD"`
	Float64ArrayTest []float64 `ubf:"T_DOUBLE_FLD"`

	StringTest      string   `ubf:"T_STRING_FLD"`
	StringArrayTest []string `ubf:"T_STRING_FLD"`

	//By default this goes as array
	CarrayTest [][]byte `ubf:"T_CARRAY_FLD"`
}

//Binary main entry
func main() {

	ret := SUCCEED
	var s TestStruct

	buf, err := atmi.NewUBF(1024)

	b1 := []byte{0, 1, 2, 3}

	b2 := []byte{4, 3, 2, 1, 0}

	////////////////////////////////////////////////////////////////////////
	// Set the test data
	////////////////////////////////////////////////////////////////////////
	if err != nil {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
		goto out
	}

	if err := buf.BChg(ubftab.T_CHAR_FLD, 0, 65); nil != err {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
		goto out
	}

	if err := buf.BChg(ubftab.T_CHAR_FLD, 1, 66); nil != err {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
		goto out
	}

	if err := buf.BChg(ubftab.T_SHORT_FLD, 0, 32000); nil != err {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
		goto out
	}

	if err := buf.BChg(ubftab.T_SHORT_FLD, 1, 32001); nil != err {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
		goto out
	}

	if err := buf.BChg(ubftab.T_LONG_FLD, 0, 9999999101); nil != err {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
		goto out
	}

	if err := buf.BChg(ubftab.T_LONG_FLD, 1, 9999999102); nil != err {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
		goto out
	}

	if err := buf.BChg(ubftab.T_FLOAT_FLD, 0, 9999999101.11); nil != err {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
		goto out
	}

	if err := buf.BChg(ubftab.T_FLOAT_FLD, 1, 9999999102.22); nil != err {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
		goto out
	}

	if err := buf.BChg(ubftab.T_DOUBLE_FLD, 0, 999999910122.888); nil != err {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
		goto out
	}

	if err := buf.BChg(ubftab.T_DOUBLE_FLD, 1, 999999910122.999); nil != err {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
		goto out
	}

	if err := buf.BChg(ubftab.T_STRING_FLD, 0, "HELLO STRING 1"); nil != err {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
		goto out
	}

	if err := buf.BChg(ubftab.T_STRING_FLD, 1, "HELLO STRING 2"); nil != err {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
		goto out
	}

	if err := buf.BChg(ubftab.T_CARRAY_FLD, 0, b1); nil != err {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
		goto out
	}

	if err := buf.BChg(ubftab.T_CARRAY_FLD, 1, b2); nil != err {
		fmt.Printf("ATMI Error %d:[%s]\n", err.Code(), err.Message())
		ret = FAIL
		goto out
	}
	////////////////////////////////////////////////////////////////////////
	// Do the Unmarshal & test the values
	////////////////////////////////////////////////////////////////////////
	buf.Unmarshal(&s)

	fmt.Printf("Got the struct [%v]\n", s)

	//CHAR
	if s.CharTest != 65 {
		fmt.Printf("s.CharTest invalid value!\n")
		ret = FAIL
	}

	if len(s.CharArrayTest) != 2 || s.CharArrayTest[0] != 65 ||
		s.CharArrayTest[1] != 66 {
		fmt.Printf("s.CharArrayTest invalid value!\n")
		ret = FAIL
	}

	//SHORT
	if s.ShortTest != 32000 {
		fmt.Printf("s.ShortTest invalid value!\n")
		ret = FAIL
	}

	if len(s.ShortArrayTest) != 2 || s.ShortArrayTest[0] != 32000 ||
		s.ShortArrayTest[1] != 32001 {
		fmt.Printf("s.ShortArrayTest invalid value!\n")
		ret = FAIL
	}

	//LONG
	if s.LongTest != 9999999101 {
		fmt.Printf("s.LongTest invalid value!\n")
		ret = FAIL
	}

	if len(s.LongArrayTest) != 2 || s.LongArrayTest[0] != 9999999101 ||
		s.LongArrayTest[1] != 9999999102 {
		fmt.Printf("s.LongArrayTest invalid value!\n")
		ret = FAIL
	}

	//FLOAT
	if s.Float32Test-9999999101.11 > 0.0001 {
		fmt.Printf("s.Float32Test invalid value!\n")
		ret = FAIL
	}

	if len(s.Float32ArrayTest) != 2 || s.Float32ArrayTest[0]-9999999101.11 > 0.0001 ||
		s.Float32ArrayTest[1]-9999999101.22 > 0.0001 {
		fmt.Printf("s.Float32ArrayTest invalid value!\n")
		ret = FAIL
	}

	//FLOAT64
	if s.Float64Test-999999910122.888 > 0.0001 {
		fmt.Printf("s.Float64Test invalid value!\n")
		ret = FAIL
	}

	if len(s.Float64ArrayTest) != 2 || s.Float64ArrayTest[0]-999999910122.888 > 0.0001 ||
		s.Float64ArrayTest[1]-999999910122.999 > 0.0001 {
		fmt.Printf("s.Float64ArrayTest invalid value!\n")
		ret = FAIL
	}

	//STRING
	if s.StringTest != "HELLO STRING 1" {
		fmt.Printf("s.StringTest invalid value!\n")
		ret = FAIL
	}

	if len(s.StringArrayTest) != 2 || s.StringArrayTest[0] != "HELLO STRING 1" ||
		s.StringArrayTest[1] != "HELLO STRING 2" {
		fmt.Printf("s.StringArrayTest invalid value!\n")
		ret = FAIL
	}

	//CARRAY
	if len(s.CarrayTest) != 2 || 0 != bytes.Compare(s.CarrayTest[0], b1) ||
		0 != bytes.Compare(s.CarrayTest[1], b2) {
		fmt.Printf("s.CarrayTest invalid value!\n")
		ret = FAIL
	}

	fmt.Println("Unmarshal tests ok...")

out:
	//Close the ATMI session
	atmi.TpTerm()
	os.Exit(ret)
}
