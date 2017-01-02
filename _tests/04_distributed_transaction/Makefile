export GOPATH=$(shell pwd)
export CGO_CFLAGS=-DOCI8_ENABLE_XA
all:
	$(MAKE) -C src

clean:
	rm -rf pkg bin



.PHONY: clean all

