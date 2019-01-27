#
# Tweak the variables based on your project.
#

# Target binaries.
TARGET := depInstallScriptGenerator

# Project main package location (can be multiple ones).
CMD_DIR := .

# Project output directory.
OUTPUT_DIR := ./bin

#
# Define all targets. At least the following commands are required:
#

.PHONY: build test clean

build:
	  go build -i -v -o $(OUTPUT_DIR)/$(TARGET) $(CMD_DIR);

mod-reset-vendor:
	@$(shell [ -f go.mod ] && go mod vendor)

test:
	@go test ./...

clean:
	@rm -vrf ${OUTPUT_DIR}/*
