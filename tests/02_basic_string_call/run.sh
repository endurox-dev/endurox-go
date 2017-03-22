#!/bin/bash

#
# @(#) Run the test case
#

pushd .
cd conf
. setndrx
popd

# Start the enduro/x app server (which will boot the our server executable)

xadmin start -y

# should print some hello world
etclient | tee test.out

# Test the logfile for content

OUT=`grep 'Hello From TESTSVC. This string is bit longer than receved in req]' test.out`

if [[ "X$OUT" == "X" ]]; then
        echo "TESTERROR: Content not found"
        exit 1
fi

echo "Test OK"

# shutdown the app server
xadmin stop -c -y

exit 0
