# Env needed for oci build:
export PKG_CONFIG_PATH=$(shell pwd)/oci8

$(info $$PKG_CONFIG_PATH is [${PKG_CONFIG_PATH}])

all: oci8 submakes

#
# Get the Oracle Driver patched for XA
#
oci8: 
	git clone https://github.com/endurox-dev/go-oci8 oci8
	cp oci8.pc oci8
	
#
# Run the submakes
#
submakes:
	$(MAKE) -C ubftab
	$(MAKE) -C client
	$(MAKE) -C server
