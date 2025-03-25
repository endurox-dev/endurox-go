module github.com/endurox-dev/endurox-go/tests/04_distributed_transaction/src/server

require (
	atmi v1.0.0
	ubftab v1.0.0
	oci8 v1.0.0
)

replace atmi v1.0.0 => ../../../../
replace ubftab v1.0.0 => ../ubftab
replace oci8 v1.0.0 => ../oci8
