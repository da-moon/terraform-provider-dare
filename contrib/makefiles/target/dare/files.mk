
.SILENT :demo-files
.PHONY :demo-files
demo-files: 
	- $(call print_running_target)
	- $(eval command=$(RM) ${ARTIFACTS_ROOT})
	- $(eval command=${command} && $(MKDIR) ${ARTIFACTS_ROOT} && )
ifneq ($(FILE_COUNT),)
	- $(eval command=$(MKDIR) ${ARTIFACTS_ROOT}/${FILE_SIZE} && ) 
	- $(eval command=${command}seq ${FILE_COUNT} | xargs -I {} ) 
endif
ifneq ($(shell which dd), )
	- $(eval command=${command}dd if=/dev/urandom  bs=1048576 count=${FILE_SIZE} of=${ARTIFACTS_ROOT}/${FILE_SIZE})
else
	- $(eval command=${command}bin$(PSEP)dare dd --size=${FILE_SIZE}MB --path=${ARTIFACTS_ROOT}/${FILE_SIZE})
endif
ifneq ($(FILE_COUNT),)
	- $(eval command=${command}/{})
endif
	- $(eval command=${command}.${FILE_EXTENSION})
ifneq ($(shell which dd), )
	- $(eval command=${command} && dd if=/dev/urandom  bs=1048576 count=${FILE_SIZE} of=${ARTIFACTS_ROOT}/no-extension)
else
	- $(eval command=${command} && bin$(PSEP)dare dd --size=${FILE_SIZE}MB --path=${ARTIFACTS_ROOT}/no-extension)
endif
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell cmd="${command}"
ifneq ($(DELAY),)
	- sleep $(DELAY)
endif
	- $(call print_completed_target)
