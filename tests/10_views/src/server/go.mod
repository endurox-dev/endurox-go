module github.com/endurox-dev/endurox-go/tests/10_views/src/server

go 1.24

require (
	atmi v1.0.0
	ubftab v1.0.0
)

replace atmi v1.0.0 => ../../../../

replace ubftab v1.0.0 => ../ubftab
