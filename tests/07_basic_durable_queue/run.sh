#!/bin/bash

#
# @(#) Run the test case
#


#
# cleanup any leftovers from previous cases
#
rm -rf var/prepared 2>/dev/null
rm -rf var/committed 2>/dev/null
rm -rf var/active 2>/dev/null

pushd .
cd conf
. setndrx
popd

# Seems on freebsd we have an issue with stack sizes, as cgo limit is 2M
# Support #251
export NDRX_MSGSIZEMAX=100000

# Start the enduro/x app server (which will boot the our server executable)

xadmin down -y
xadmin start -y

# should print some hello world
etclient

RET=$?


echo "Exit $RET"

# shutdown the app server
xadmin stop -c -y

exit $RET
