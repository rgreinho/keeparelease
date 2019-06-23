# Project name.
PROJECT_NAME = keeparelease

# Makefile parameters.
TAG ?= $(shell git describe)

# General.
SHELL = /bin/bash
TOPDIR = $(shell git rev-parse --show-toplevel)

# Project specifics.
BUILD_DIR = dist
PLATFORMS = linux darwin
OS = $(word 1, $@)
GOOS = $(shell uname -s | tr A-Z a-z)
GOARCH = amd64

default: build

.PHONY: help
help: # Display help
	@awk -F ':|##' \
		'/^[^\t].+?:.*?##/ {
			printf "\033[36m%-30s\033[0m %s\n", $$1, $$NF \
		}' $(MAKEFILE_LIST) | sort

.PHONY: build
build: ## Build the project for the current platform
	mkdir -p $(BUILD_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BUILD_DIR)/$(PROJECT_NAME)-$(TAG)-$(GOOS)-$(GOARCH)

.PHONY: test
test: ## Run the unit tests
	go test .../keeparelease

.PHONY: clean
clean: clean-code ## Clean everything (!DESTRUCTIVE!)

.PHONY: clean-code
clean-code: ## Remove unwanted files in this project (!DESTRUCTIVE!)
	@cd $(TOPDIR) && git clean -ffdx && git reset --hard

dist: $(PLATFORMS) ## Package the project for all available platforms

.PHONY: publish
publish: dist ## Create GitHub release
	keeparelease -a dist/*

.PHONY: setup
setup: ## Setup the full environment (default)
	go mod tidy

$(PLATFORMS): # Build the project for all available platforms
	mkdir -p $(BUILD_DIR)
	GOOS=$(OS) GOARCH=$(GOARCH) go build -o $(BUILD_DIR)/$(PROJECT_NAME)-$(TAG)-$(OS)-$(GOARCH)
.PHONY: $(PLATFORMS)
