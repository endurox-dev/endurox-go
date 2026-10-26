EXFIELDTBL ?= $(shell pkg-config --variable=prefix atmisrvinteg)/share/endurox/ubftab/Exfields

all: test.fd.go Exfields.go

test.fd.go: test.fd
	mkfldhdr -m1 -pubftab test.fd

Exfields.go: $(EXFIELDTBL)
	mkfldhdr -m1 -pubftab "$(EXFIELDTBL)"

.PHONY: all
