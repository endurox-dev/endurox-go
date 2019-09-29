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
#valgrind --tool=memcheck --leak-check=yes client
client

ret=$?

if [[ $ret -eq 0 ]]; then
	echo "Test OK"
else
	echo "Test failed"
fi

# shutdown the app server
xadmin stop -c -y

exit $ret
