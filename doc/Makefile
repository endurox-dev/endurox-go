
SOURCES := $(shell find $(SOURCEDIR) -name '*.adoc')
HTML = $(patsubst %.adoc, %.html, $(SOURCES))

all: $(HTML)

$(info $$HTML is [${HTML}])

$(HTML): $(SOURCES)
	a2x -D ./compiled -f xhtml $<
	
clean:
	rm -f $(HTML) docbook-xsl.css

.PHONY: all clean
