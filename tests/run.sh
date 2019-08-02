#!/bin/bash

#
# @(#) Integration tests
#

> ./test.out
# Have some terminal output...
tail -f test.out &

(

M_tests=0
M_ok=0
M_fail=0

run_test () {

        test=$1
        M_tests=$((M_tests + 1))
        echo "*** RUNNING [$test]"

        pushd .
        cd $test
        ./run.sh
        ret=$?
        popd
        
        echo "*** RESULT [$test] $ret"
        
        if [[ $ret -eq 0 ]]; then
                M_ok=$((M_ok + 1))
        else
                M_fail=$((M_fail + 1))
        fi
}

run_test "01_basic_ubf_call"
run_test "02_basic_string_call"
run_test "03_basic_carray_call"
run_test "05_basic_json_call"
run_test "06_ubf_marshal"
run_test "07_basic_durable_queue"
run_test "08_logging"
run_test "09_return_manual_buffer"
run_test "10_views"
run_test "11_ubf_tests"
run_test "12_tpimportexp"

echo "*** SUMMARY $M_tests tests executed. $M_ok passes, $M_fail failures"

xadmin killall tail

exit $M_fail

) > test.out 2>&1

