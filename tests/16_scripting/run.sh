#!/bin/bash

#
# @(#) Scripting lifecycle, callbacks and buffer ownership
#

case "$#:${1:-}" in
    0:)
        echo "SKIPPED: scripting requires separately licensed endurox-python; use -scripting"
        exit 0
        ;;
    1:-scripting) ;;
    *)
        echo "Usage: $0 [-scripting]" >&2
        exit 2
        ;;
esac

cd "$(dirname "$0")" || exit 1
export TEST_DIR=$(pwd)
if [[ ! -x "$TEST_DIR/bin/client" || ! -x "$TEST_DIR/bin/goscriptsv" ]]; then
    echo "TESTERROR: build tests/16_scripting before running it" >&2
    exit 1
fi

pushd conf || exit 1
. ./setndrx || exit 1
popd || exit 1

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
    echo "TESTERROR: failed to start the 16_scripting test domain" >&2
    exit 1
fi

"$TEST_DIR/bin/client" -scripting
ret=$?

if [[ $ret -eq 0 ]]; then
    echo "Test OK"
else
    echo "TESTERROR: 16_scripting cases failed ($ret)"
fi
exit "$ret"
