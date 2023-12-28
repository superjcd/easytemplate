.DEFAULT_GOAL := all 

# =============================================================================
# Globals:
ROOT_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
OUTPUT_DIR  := $(ROOT_DIR)/build
PLATFORMS ?=   linux_amd64 windows_amd64
VERSION = 0.1.0
COMMAND = easytemplate

# =============================================================================

.PHONY: all 
all: tidy build 

     
.PHONY: tidy 
tidy:
	@go mod tidy
.PHONY: build
build:  $(foreach P,${PLATFORMS}, $(addprefix build., $(P)))

.PHONY: build.%
build.%:
	$(eval OS:= $(word 1,$(subst _, ,$*)))
	$(eval ARCH := $(word 2,$(subst _, ,$*)))  
	$(if $(findstring windows, $(OS)), $(eval EXE_SUFFIX:=.exe), $(eval EXE_SUFFIX:=''))
	@go env -w CGO_ENABLED=0  GOOS=$(OS) GOARCH=$(ARCH)
	@echo "====>Build binary for ${COMMAND}, with OS: $(OS), ARCH:$(ARCH)"
	@go build -o $(OUTPUT_DIR)/$(COMMAND)$(EXE_SUFFIX)  $(ROOT_DIR)/cmd



.PHONY: clean
clean:
	$(if $(findstring Windows, $(OS)), $(shell rmdir /Q /S ${OUTPUT_DIR}), $(shell rm -rf ${OUTPUT_DIR}))
	@echo "====>output directory  is removed sucessfully"
