module github.com/endurox-dev/endurox-go/tests/12_tpimportexp/src/client

go 1.24

require (
	atmi v1.0.0
	ubftab v1.0.0
)

replace atmi v1.0.0 => ../../../../

replace ubftab v1.0.0 => ../ubftab
