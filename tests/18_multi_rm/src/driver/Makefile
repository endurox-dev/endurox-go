OS := $(shell uname -s)
ifeq ($(OS),Darwin)
LIBEXT := dylib
SHARED := -dynamiclib
else
LIBEXT := so
SHARED := -shared
endif
all:
	mkdir -p ../../bin
	$(CC) $(CFLAGS) -fPIC $(SHARED) mock.c -o ../../bin/libgormtest.$(LIBEXT) $(shell pkg-config --cflags --libs atmisrvinteg)

.PHONY: all
