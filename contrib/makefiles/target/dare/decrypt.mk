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
	- $(eval command=$(command) --input=${ARTIFACTS_ROOT}/no-extension.enc)
	- $(eval command=$(command) --output=${ARTIFACTS_ROOT}/decrypted)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell cmd="${command}"
ifneq ($(DELAY),)
	- sleep $(DELAY)
endif
	- $(call print_running_target, comparing '${ARTIFACTS_ROOT}/no-extension' and ${ARTIFACTS_ROOT}/decrypted${ARTIFACTS_ROOT}/no-extension)
	- $(eval command=md5sum ${ARTIFACTS_ROOT}/no-extension ${ARTIFACTS_ROOT}/decrypted${ARTIFACTS_ROOT}/no-extension)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell cmd="${command}"
ifneq ($(DELAY),)
	- sleep $(DELAY)
endif
	- $(call print_completed_target)

.SILENT: dare-multi-decrypt
.PHONY: dare-multi-decrypt
dare-multi-decrypt: dare-multi-encrypt
	- $(call print_running_target)
ifneq ($(FILE_COUNT),)
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
	- $(eval command=$(command) --input=${ARTIFACTS_ROOT}/encrypted)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell cmd="${command}"
ifneq ($(DELAY),)
	- sleep $(DELAY)
endif
	- $(call print_running_target, camparing hashes)
	- $(eval command=seq ${FILE_COUNT} | xargs -I {})
	- $(eval command=$(command) md5sum ${ARTIFACTS_ROOT}/${FILE_SIZE}/{}.${FILE_EXTENSION})
	- $(eval command=$(command) ${ARTIFACTS_ROOT}/encrypted${ARTIFACTS_ROOT}/${FILE_SIZE}/{}.${FILE_EXTENSION})
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell cmd="${command}"
ifneq ($(DELAY),)
	- sleep $(DELAY)
endif
endif
	- $(call print_completed_target)