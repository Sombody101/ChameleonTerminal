# Check if Go is installed
ifeq (, $(shell which go))
	$(error Go is not installed!)
endif

PACKAGE := gecko

# Local to the root of the project, not the repository
output_path ?= ../build/gecko

# Get version using git describe
VERSION := $(shell git describe --tags --always --abbrev=0 --match='v[0-9]*.[0-9]*.[0-9]*' 2> /dev/null || echo "v0.0.0")

# Get commit hash using git rev-parse
COMMIT_HASH := $(shell git rev-parse --short HEAD)

# Get build timestamp
BUILD_TIMESTAMP := $(shell date '+%Y-%m-%dT%H:%M:%S')

# Define ldflags with version information
LDFLAGS := \
  -X $(PACKAGE)/version.VERSION=$(VERSION)+RELEASE \
  -X $(PACKAGE)/version.COMMIT_HASH=$(COMMIT_HASH) \
  -X $(PACKAGE)/version.BUILD_TIME=$(BUILD_TIMESTAMP)

build:
	go build -C ./src/ -ldflags \
		"-w -s $(LDFLAGS)" \
		-o $(output_path)
	@echo Built Gecko v$(VERSION) to $(output_path)

.PHONY: build
