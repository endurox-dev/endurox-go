module github.com/endurox-dev/endurox-go/tests/14_saga/src/client

go 1.26.0

require (
	github.com/endurox-dev/endurox-go/v9 v9.0.0
	ubftab v1.0.0
)

replace github.com/endurox-dev/endurox-go/v9 => ../../../..

replace ubftab v1.0.0 => ../ubftab
