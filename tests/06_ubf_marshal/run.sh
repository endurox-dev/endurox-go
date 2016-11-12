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
client

ret=$?

echo "Test OK"

# shutdown the app server
xadmin stop -c -y

exit $?
