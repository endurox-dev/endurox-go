export FIELDTBLS=test.fd
export FLDTBLDIR=$(shell pwd)

all: test.fd.go

test.fd.go: test.fd
	mkfldhdr -m1 -pubftab

.PHONY: all
