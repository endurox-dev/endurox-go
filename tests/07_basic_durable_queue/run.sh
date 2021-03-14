#!/bin/bash

#
# @(#) Run the test case
#


#
# cleanup any leftovers from previous cases
#
rm -rf var/qspace1/prepared 2>/dev/null
rm -rf var/qspace1/committed 2>/dev/null
rm -rf var/qspace1/active 2>/dev/null
rm -rf var/tm2 2>/dev/null
mkdir var/tm2

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

#
# Move process to RM2/NULL switch
#
export NDRX_XA_RES_ID=2
export NDRX_XA_DRIVERLIB=libndrxxanulls.${NDRX_LIBEXT}

# should print some hello world
etclient

RET=$?


echo "Exit $RET"

# shutdown the app server
xadmin stop -c -y

exit $RET
