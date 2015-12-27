#!/bin/bash

#
# @(#) Integration tests
#

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

echo "*** SUMMARY $M_tests tests executed. $M_ok passes, $M_fail failures"
