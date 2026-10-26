#!/bin/bash

#
# @(#) DLM commands, asynchronous replies and listing
#

case "$#:${1:-}" in
    0:) echo "SKIPPED: DLM requires -dlm and Enduro/X EE"; exit 0 ;;
    1:-dlm) ;;
    *) echo "Usage: $0 [-dlm]" >&2; exit 2 ;;
esac

cd "$(dirname "$0")" || exit 1
export TEST_DIR=$(pwd)
if [[ ! -x "$TEST_DIR/bin/client" || ! -x "$TEST_DIR/bin/godlmsv" ]]; then
    echo "TESTERROR: build tests/17_dlm before running it" >&2
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
mkdir -p "$TEST_DIR/edb_log" "$TEST_DIR/edb_snapshot" || exit 1
if ! xadmin start -y; then
    echo "TESTERROR: failed to start the 17_dlm test domain" >&2
    exit 1
fi

ready=0
for attempt in $(seq 1 30); do
    if xadmin psc 2>/dev/null | grep -q '@DLGOTEST'; then
        ready=1
        break
    fi
    sleep 1
done
if [[ $ready -ne 1 ]]; then
    echo "TESTERROR: DLM leader service not available"
    exit 1
fi
sleep 5
"$TEST_DIR/bin/client"
ret=$?

if [[ $ret -eq 0 ]]; then
    echo "Test OK"
else
    echo "TESTERROR: 17_dlm cases failed ($ret)"
fi
exit "$ret"
