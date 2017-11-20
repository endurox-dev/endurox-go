#!/bin/bash

#
# @(#) Run the test case
#

pushd .
cd conf
. setndrx
popd

# Seems on freebsd we have an issue with stack sizes, as cgo limit is 2M
# Support #251
export NDRX_MSGSIZEMAX=100000

# Start the enduro/x app server (which will boot the our server executable)

xadmin start -y

# should print some hello world
etclient

RET=$?


echo "Exit $RET"

# shutdown the app server
xadmin stop -c -y

exit $RET
