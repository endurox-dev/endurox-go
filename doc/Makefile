## 
## Enduro/X Go lang module documentation
##
## @file Makefile
## 
## -----------------------------------------------------------------------------
## Enduro/X Middleware Platform for Distributed Transaction Processing
## Copyright (C) 2015, ATR Baltic, SIA. All Rights Reserved.
## This software is released under one of the following licenses:
## GPL or ATR Baltic's license for commercial use.
## -----------------------------------------------------------------------------
## GPL license:
## 
## This program is free software; you can redistribute it and/or modify it under
## the terms of the GNU General Public License as published by the Free Software
## Foundation; either version 2 of the License, or (at your option) any later
## version.
##
## This program is distributed in the hope that it will be useful, but WITHOUT ANY
## WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
## PARTICULAR PURPOSE. See the GNU General Public License for more details.
##
## You should have received a copy of the GNU General Public License along with
## this program; if not, write to the Free Software Foundation, Inc., 59 Temple
## Place, Suite 330, Boston, MA 02111-1307 USA
##
## -----------------------------------------------------------------------------
## A commercial use license is available from ATR Baltic, SIA
## contact@atrbaltic.com
## -----------------------------------------------------------------------------
##
SOURCEDIR:=.
DOCOUT:=./out/
IMG:=./images/

SOURCES := $(shell find -name '*.adoc' -printf "%f\n")
HTML := $(addprefix $(DOCOUT), $(patsubst %.adoc, %.html, $(SOURCES)))


SOURCES_GRPAHS :=$(shell find -name '*.dia' -printf "%f\n")
GRPAHS  := $(addprefix $(IMG), $(patsubst %.dia, %.png, $(SOURCES_GRPAHS)))
#GRPAHS  := $(patsubst %.dia, %.png, $(SOURCES_GRPAHS))

all: $(GRPAHS) $(HTML)

$(info $$HTML is [${HTML}])

$(info $$GRPAHS is [${GRPAHS}])

#
# Build the charts (have to put the images locally so that we see them in 
#
$(GRPAHS): $(SOURCES_GRPAHS)
	dia --size=600x --filter=png -e ./images/$(patsubst %.dia,%,$<).png -t png $<
	cp ./images/$(patsubst %.dia,%,$<).png ./out/images

#
# Build the document
#
$(HTML): $(SOURCES)
	a2x -D ./out -f xhtml $<
	
clean:
	rm -f $(HTML) ./out/docbook-xsl.css $(addprefix ./out/images/, $(patsubst %.dia, %.png, $(SOURCES_GRPAHS)))

.PHONY: clean



