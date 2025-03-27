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
rm /tmp/01_CLIENT.log 2>/dev/null

# should print some hello world
client > /tmp/01_CLIENT.log

RET=$?

# Test the logfile for content
OUT=`grep 'Hello World from Enduro/X service' /tmp/01_CLIENT.log`

if [[ "X$OUT" == "X" ]]; then
        echo "TESTERROR: Content not found"
        RET=1
fi

echo "Runtime check (wait 60)"
# Check the runtime.. for intr/SIGURG (shall not pass to C)
# the signal ignoring shall do
runtimeck

if [ $? -ne 0 ]; then
    echo "Runtime check failed (problems with atmi.RuntimeInit() or paly with GODEBUG=\"asyncpreemptoff=...\")"
    RET=1
fi

if [[ "X$RET" == "X0" ]]; then
    echo "Test OK"

else
    echo "Test FAIL"
fi

# shutdown the app server
xadmin stop -c -y

exit $RET
