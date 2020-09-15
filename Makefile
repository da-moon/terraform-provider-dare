include vars.mk
include contrib/makefiles/pkg/base/base.mk
include contrib/makefiles/pkg/string/string.mk
include contrib/makefiles/pkg/color/color.mk
include contrib/makefiles/pkg/functions/functions.mk
include contrib/makefiles/target/git/git.mk
include contrib/makefiles/target/buildenv/buildenv.mk
include contrib/makefiles/target/go/go.mk
SHELL := /bin/bash
THIS_FILE := $(firstword $(MAKEFILE_LIST))
SELF_DIR := $(dir $(THIS_FILE))

.SILENT: build
.PHONY: build
build:
	- $(call print_running_target)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) go-build
	- $(call print_completed_target)
.SILENT: clean
.PHONY: clean
clean:
	- $(call print_running_target)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) go-clean
	- $(call print_completed_target)

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
.SILENT: dare-single
.PHONY: dare-single
dare-single:
	- $(call print_running_target)
	- $(call print_running_target, encrypting single file '${ARTIFACTS_ROOT}/no-extension')
	- $(eval command=bin/dare encrypt)
ifneq (${LOG_LEVEL}, )
	- $(eval command=$(command) --log-level=$(LOG_LEVEL))
endif
ifneq (${ENCRYPTION_KEY}, )
	- $(eval command=$(command) --master-key=$(ENCRYPTION_KEY))
endif
	- $(eval command=$(command) --input=${ARTIFACTS_ROOT}/no-extension)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell cmd="${command}"
ifneq ($(DELAY),)
	- sleep $(DELAY)
endif
	- $(call print_completed_target)
.SILENT: dare-multi
.PHONY: dare-multi
dare-multi: 
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

.PHONY: dare
.SILENT:dare
dare: demo-files
	- $(call print_running_target)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) dare-single
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) dare-multi
	- $(call print_completed_target)
