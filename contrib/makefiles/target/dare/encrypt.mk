.SILENT: dare-single-encrypt
.PHONY: dare-single-encrypt
dare-single-encrypt: dare-single-encrypt
	- $(call print_running_target)
	- $(eval name=$(@:%-encrypt=%))
	- $(call print_running_target, encrypting single file '${ARTIFACTS_ROOT}/no-extension')
	- $(eval command=bin/dare encrypt)
ifneq (${LOG_LEVEL}, )
	- $(eval command=$(command) --log-level=$(LOG_LEVEL))
endif
ifneq (${ENCRYPTION_KEY}, )
	- $(eval command=$(command) --master-key=$(ENCRYPTION_KEY))
endif
	- $(eval command=$(command) --input=${ARTIFACTS_ROOT}/no-extension)
	- $(eval command=$(command) > /tmp/$(name).json)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell cmd="${command}"
ifneq ($(DELAY),)
	- sleep $(DELAY)
endif
	- $(call print_completed_target)


.SILENT: dare-multi-encrypt
.PHONY: dare-multi-encrypt
dare-multi-encrypt: 
	- $(call print_running_target)
	- $(call print_running_target, encrypting directory '${ARTIFACTS_ROOT}/' with key file '${ENCRYPTION_KEY_FILE}' and regex '*.${FILE_EXTENSION}' for recursive-search and storing in ${ARTIFACTS_ROOT}/encrypted directory)
	- $(eval command=echo '$(ENCRYPTION_KEY)' > ${ENCRYPTION_KEY_FILE})
	- $(eval command=$(command) && bin/dare encrypt)
ifneq (${LOG_LEVEL}, )
	- $(eval command=$(command) --log-level=$(LOG_LEVEL))
endif
ifneq (${ENCRYPTION_KEY_FILE}, )
	- $(eval command=$(command) --master-key-file=$(ENCRYPTION_KEY_FILE))
endif
	- $(eval command=$(command) --input=${ARTIFACTS_ROOT}/)
	- $(eval command=$(command) --output=${ARTIFACTS_ROOT}/encrypted)
	- $(eval command=$(command) --regex='*.${FILE_EXTENSION}')
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell cmd="${command}"
ifneq ($(DELAY),)
	- sleep $(DELAY)
endif
	- $(call print_completed_target)
