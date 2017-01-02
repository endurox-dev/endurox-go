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

# Run the client
client

RET=$?

if [ "X$RET" != "X0" ]; then
        echo "Invalid exit code $RET"
        exit $RET
fi

# Test the logfile for content
OUT=`grep '00 01 02 03 04 05 06 08 09' /tmp/08_client_process.log`
if [[ "X$OUT" == "X" ]]; then
        echo "TESTERROR: [00 01 02 03 04 05 06 08 09] not found in /tmp/08_client_process.log"
        exit 1
fi

# Test the logfile for content
OUT=`grep '02 03 04 05 06 07 08 09 0a' /tmp/08_client_process.log`
if [[ "X$OUT" == "X" ]]; then
        echo "TESTERROR: [02 03 04 05 06 07 08 09 0a] not found in /tmp/08_client_process.log"
        exit 1
fi

# test Th1
OUT=`grep 'Hello from TH1' /tmp/08_th1.log`
if [[ "X$OUT" == "X" ]]; then
        echo "TESTERROR: [Hello from TH1] not found in /tmp/08_th1.log"
        exit 1
fi

# test Th2
OUT=`grep 'Hello from TH2' /tmp/08_th2.log`
if [[ "X$OUT" == "X" ]]; then
        echo "TESTERROR: [Hello from TH2] not found in /tmp/08_th2.log"
        exit 1
fi

# Test reqeust logging
OUT=`grep 'HELLO FROM CLIENT' /tmp/08_request95.log`
if [[ "X$OUT" == "X" ]]; then
        echo "TESTERROR: [HELLO FROM CLIENT] not found in /tmp/08_request95.log"
        exit 1
fi

echo "Test OK"

# shutdown the app server
xadmin stop -c -y

exit 0
