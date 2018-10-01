##
## @brief Enduro/X Go lang module documentation
##
## @file Makefile
##
## -----------------------------------------------------------------------------
## Enduro/X Middleware Platform for Distributed Transaction Processing
## Copyright (C) 2009-2016, ATR Baltic, Ltd. All Rights Reserved.
## Copyright (C) 2017-2018, Mavimax, Ltd. All Rights Reserved.
## This software is released under one of the following licenses:
## AGPL or Mavimax's license for commercial use.
## -----------------------------------------------------------------------------
## AGPL license:
## 
## This program is free software; you can redistribute it and/or modify it under
## the terms of the GNU Affero General Public License, version 3 as published
## by the Free Software Foundation;
##
## This program is distributed in the hope that it will be useful, but WITHOUT ANY
## WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
## PARTICULAR PURPOSE. See the GNU Affero General Public License, version 3
## for more details.
##
## You should have received a copy of the GNU Affero General Public License along 
## with this program; if not, write to the Free Software Foundation, Inc., 
## 59 Temple Place, Suite 330, Boston, MA 02111-1307 USA
##
## -----------------------------------------------------------------------------
## A commercial use license is available from Mavimax, Ltd
## contact@mavimax.com
## -----------------------------------------------------------------------------
##
SOURCEDIR:=.
DOCOUT:=./out/
IMG:=./images/


DOCSOURCES := $(shell ls -1 ../*.go ./docupdate.pl);
SOURCES := $(shell find -name '*.adoc' -printf "%f\n")
HTML := $(addprefix $(DOCOUT), $(patsubst %.adoc, %.html, $(SOURCES)))


SOURCES_GRPAHS :=$(shell find -name '*.dia' -printf "%f\n")
GRPAHS  := $(addprefix $(IMG), $(patsubst %.dia, %.png, $(SOURCES_GRPAHS)))
#GRPAHS  := $(patsubst %.dia, %.png, $(SOURCES_GRPAHS))

all: endurox-go-book.adoc $(GRPAHS) $(HTML) 

$(info $$HTML is [${HTML}])

$(info $$GRPAHS is [${GRPAHS}])

$(info $$DOCSOURCES is [${DOCSOURCES}])

#
# Generate the book from javadoc comments
#
endurox-go-book.adoc: $(DOCSOURCES)
	./docupdate.pl ../*

#
# Build the charts (have to put the images locally so that we see them in 
#
$(GRPAHS): $(SOURCES_GRPAHS)
	dia --size=1200x --filter=png -e ./images/$(patsubst %.dia,%,$<).png -t png $<
	cp ./images/$(patsubst %.dia,%,$<).png ./out/images

#
# Build the document
#
$(HTML): $(SOURCES)
	a2x -D ./out -f xhtml $<
	
clean:
	rm -f $(HTML) ./out/docbook-xsl.css $(addprefix ./out/images/, $(patsubst %.dia, %.png, $(SOURCES_GRPAHS)))

.PHONY: clean



# vim: set ts=4 sw=4 et smartindent:
