include vars.mk
include contrib/makefiles/pkg/base/base.mk
include contrib/makefiles/pkg/string/string.mk
include contrib/makefiles/pkg/color/color.mk
include contrib/makefiles/pkg/functions/functions.mk
include contrib/makefiles/target/git/git.mk
include contrib/makefiles/target/buildenv/buildenv.mk
include contrib/makefiles/target/go/go.mk
include contrib/makefiles/target/dare/files.mk
include contrib/makefiles/target/dare/encrypt.mk
include contrib/makefiles/target/dare/decrypt.mk
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
.SILENT: install
.PHONY: install
install: build
	- $(call print_running_target)
	- $(eval name=$(PROJECT_NAME))
	- $(eval command=$(RM) ~/.terraform.d/plugins/$(name) )
	- $(eval command=&& $(MKDIR) ~/.terraform.d/plugins )
	- $(eval command=&& $(CP) bin/$(name) ~/.terraform.d/plugins/$(name))
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell cmd="${command}"
	- $(call print_completed_target)


.PHONY: dare
.SILENT:dare
dare: demo-files
	- $(call print_running_target)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) dare-single-decrypt
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) dare-multi-encrypt
	- $(call print_completed_target)


