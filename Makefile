.PHONY: build clean all test

NAME := yddns
VERSION	?= dev
PACKAGE	:= github.com/drieschel/$(NAME)
BUILD_DIR := dist/$(VERSION)

GO = go
LDFLAGS := "-X $(PACKAGE)/cmd.version=$(VERSION)"

OS ?=
GOOS ?=
ifeq ($(GOOS),)
	ifeq ($(OS),Windows_NT)
		GOOS := windows
	else
		GOOS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
	endif
endif

GOARCH ?=
ifeq ($(GOARCH),)
	GOARCH := $(shell uname -m)
	ifeq ($(GOARCH),x86_64)
		GOARCH := amd64
	else ifeq ($(GOARCH),aarch64)
		GOARCH := arm64
	else ifeq ($(GOARCH),i386)
    	GOARCH := 386
	endif
endif


OUTPUT_NAME := $(GOOS)-$(GOARCH)
ifeq ($(GOOS),windows)
	OUTPUT_NAME := $(OUTPUT_NAME).exe
endif

build:
	mkdir -p "$(BUILD_DIR)/cache"
	cp -ur config.toml.example templates $(BUILD_DIR)
	touch $(BUILD_DIR)/config.toml
	GOARCH=$(GOARCH) GOOS=$(GOOS) $(GO) build -o $(BUILD_DIR)/bin/$(OUTPUT_NAME) -ldflags $(LDFLAGS)

clean:
	rm -rf $(BUILD_DIR)

all: clean build

test:
	GOOS="" GOARCH="" $(GO) test ./... -ldflags $(LDFLAGS)