.PHONY: build

NAME	:= yddns
VERSION	?= dev
PACKAGE	:= github.com/drieschel/$(NAME)

GOOS ?=
ifeq ($(GOOS),)
	GOOS = linux
	ifeq ($(OS),Windows_NT)
		GOOS = windows
	endif
endif

GOARCH ?=
ifeq ($(GOARCH),)
	GOARCH = $(shell echo $(PROCESSOR_ARCHITECTURE) | tr '[:upper:]' '[:lower:]')
endif

ifeq ($(GOARCH),x86)
	GOARCH = 386
endif

build:
	GOARCH=$(GOARCH) GOOS=$(GOOS) go build -ldflags "-X $(PACKAGE)/cmd.version=$(VERSION)"