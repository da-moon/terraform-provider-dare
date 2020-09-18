.SILENT: dare-single-decrypt
.PHONY: dare-single-decrypt
dare-single-decrypt: dare-single-encrypt
	- $(call print_running_target)
	- $(eval name=$(@:%-decrypt=%))
	- $(eval NONCE=$(shell jq -c '.[].random_nonce'  /tmp/$(name).json))
	- $(call print_running_target, decrypting single file '${ARTIFACTS_ROOT}/no-extension' with encryption key '$(ENCRYPTION_KEY)' and nonce '$(NONCE)')
	- $(eval command=bin/dare decrypt)
ifneq (${LOG_LEVEL}, )
	- $(eval command=$(command) --log-level=$(LOG_LEVEL))
endif
ifneq (${ENCRYPTION_KEY}, )
	- $(eval command=$(command) --master-key=$(ENCRYPTION_KEY))
endif
	- $(eval command=$(command) --nonce=$(NONCE))
	- $(eval command=$(command) --input=${ARTIFACTS_ROOT}/no-extension)
	- $(eval command=$(command))
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell cmd="${command}"
ifneq ($(DELAY),)
	- sleep $(DELAY)
endif
	- $(call print_completed_target)