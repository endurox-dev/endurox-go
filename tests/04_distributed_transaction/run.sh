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
client > test.out 2>&1
ret=$?

# Test the logfile for content

OUT=`grep -E 'Got customer id: [0-9]+' test.out`

if [[ "X$OUT" == "X" ]]; then
        echo "TESTERROR: Content not found"
        exit 1
fi

if [[ "$ret" == "0" ]]; then
        echo "Test OK"
else
        echo "Ret $ret -> FAIL"
fi

# shutdown the app server
xadmin stop -c -y

exit $ret
