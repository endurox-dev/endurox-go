#!/bin/bash

#
# @(#) Saga and transaction argument validation
#

if [[ $# -ne 0 ]]; then
    echo "Usage: $0" >&2
    exit 2
fi

cd "$(dirname "$0")" || exit 1
export TEST_DIR=$(pwd)
if [[ ! -x "$TEST_DIR/bin/client" || ! -x "$TEST_DIR/bin/gosagasv" ]]; then
    echo "TESTERROR: build tests/14_saga before running it" >&2
    exit 1
fi

pushd conf || exit 1
. ./setndrx || exit 1
popd || exit 1

mkdir -p "$TEST_DIR/tm-log" || exit 1

# Stop the case's domain on success, failure, or interruption.
trap '
    ret=$?
    trap - EXIT
    xadmin stop -c -y
    stop_ret=$?
    if [[ $ret -eq 0 && $stop_ret -ne 0 ]]; then
        ret=$stop_ret
    fi
    exit "$ret"
' EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
if ! xadmin start -y; then
    echo "TESTERROR: failed to start the 14_saga test domain" >&2
    exit 1
fi

"$TEST_DIR/bin/client"
ret=$?

if [[ $ret -eq 0 ]]; then
    echo "Test OK"
else
    echo "TESTERROR: 14_saga cases failed ($ret)"
fi
exit "$ret"
