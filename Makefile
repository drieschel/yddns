.PHONY: build

NAME := yddns
VERSION	?= dev
PACKAGE	:= github.com/drieschel/$(NAME)

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

build:
	GOARCH=$(GOARCH) GOOS=$(GOOS) go build -ldflags "-X $(PACKAGE)/cmd.version=$(VERSION)"