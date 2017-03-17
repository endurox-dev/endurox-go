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
rm /tmp/09_CLIENT.log 2>/dev/null

# should print some hello world
client 

RET=$?

# shutdown the app server
xadmin stop -c -y

exit $?
