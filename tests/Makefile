# Do recursive builds
all:
	$(MAKE) -C 01_basic_ubf_call
	$(MAKE) -C 02_basic_string_call
	$(MAKE) -C 03_basic_carray_call
	$(MAKE) -C 05_basic_json_call
	$(MAKE) -C 06_ubf_marshal
	$(MAKE) -C 07_basic_durable_queue
	$(MAKE) -C 08_logging
	$(MAKE) -C 09_return_manual_buffer
	$(MAKE) -C 10_views
	$(MAKE) -C 11_ubf_tests

clean:
	$(MAKE) -C 01_basic_ubf_call clean
	$(MAKE) -C 02_basic_string_call clean
	$(MAKE) -C 03_basic_carray_call clean
	$(MAKE) -C 05_basic_json_call clean
	$(MAKE) -C 06_ubf_marshal clean
	$(MAKE) -C 07_basic_durable_queue clean
	$(MAKE) -C 08_logging clean
	$(MAKE) -C 09_return_manual_buffer
	$(MAKE) -C 10_views
	$(MAKE) -C 11_ubf_tests

