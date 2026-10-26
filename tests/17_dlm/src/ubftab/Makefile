EXFIELDTBL ?= $(shell pkg-config --variable=prefix atmisrvinteg)/share/endurox/ubftab/Exfields

all: Exfields.go

Exfields.go: $(EXFIELDTBL)
	mkfldhdr -m1 -pubftab "$(EXFIELDTBL)"

.PHONY: all
