#build/exec env setting
DOCKER_ENV:=false
DELAY:=
# golang specific vars
GO_PKG:=github.com/da-moon/terraform-provider-dare
GO_IMAGE=golang:buster
MOD=on
# variables for demo
LOG_LEVEL=TRACE
FILE_SIZE:=5
FILE_EXTENSION:=demo
ARTIFACTS_ROOT:=/tmp/artifacts
FILE_COUNT:=3
ENCRYPTION_KEY:=092fd9b84801ad6ba6bf9b1119087c4b6cc075c62e551bac8c21be8023f935a9
ENCRYPTION_KEY_FILE:=${ARTIFACTS_ROOT}/.secret
