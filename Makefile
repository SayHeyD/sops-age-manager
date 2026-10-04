GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
CGO_ENABLED ?= $(if $(filter darwin,$(GOOS)),1,0)
EXT := $(if $(filter windows,$(GOOS)),.exe,)
OUTPUT ?= sam$(EXT)

.PHONY: test build build-release build-all docs clean dmg lint winres

winres:
	@echo "Generating Windows resources... 🔄"
	@GOOS="" GOARCH="" go run github.com/tc-hib/go-winres@v0.3.3 simply --icon assets/Logo.png --arch amd64,arm64 --product-name "SOPS Age Manager" --file-description "SOPS Age Manager" --manifest gui && echo "Finished generating Windows resources ✅" || { echo "Failed to generate Windows resources ❌"; exit 1; }

build:
	@if [ "$(GOOS)" = "windows" ] && [ ! -f "rsrc_windows_$(GOARCH).syso" ]; then $(MAKE) winres; fi
	@echo "Building application ($(GOOS)/$(GOARCH))... 🔄"
	@mkdir -p bin
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -tags main -o bin/$(OUTPUT) . && echo "Finished building ✅" || { echo "Build failed ❌"; exit 1; }

build-release:
	@$(MAKE) build OUTPUT=sam_$(GOOS)_$(GOARCH)$(EXT)

build-all:
	@echo "Building all release targets... 🔄"
	@$(MAKE) build-release GOOS=darwin GOARCH=arm64
	@$(MAKE) build-release GOOS=darwin GOARCH=amd64
	@$(MAKE) build-release GOOS=linux GOARCH=amd64
	@$(MAKE) build-release GOOS=linux GOARCH=arm64
	@$(MAKE) build-release GOOS=windows GOARCH=amd64
	@$(MAKE) build-release GOOS=windows GOARCH=arm64
	@echo "All release targets built successfully ✅"

docs:
	@echo "Building docs application... 🔄"
	@mkdir -p bin
	@go build -tags docs -o bin/sam . && echo "Finished building ✅" || { echo "Build failed ❌"; exit 1; }
	@echo "\nGenerating docs... 🔄"
	@./bin/sam && echo "Finished docs generation ✅" || { echo "Docs generation failed ❌"; exit 1; }

test:
	@echo "Running test... 🔄"
	@go test ./... -v -parallel 4 -cover && echo "Tests finished successfully ✅" || { echo "Tests finished with errors ❌"; exit 1; }

lint:
	@echo "Running linter... 🔄"
	@golangci-lint run && echo "Lint passed ✅" || { echo "Lint failed ❌"; exit 1; }

clean:
	@echo "Cleaning build products... 🔄"
	@rm -rf ./bin ./build *.syso && echo "Cleaning done ✅" || { echo "Cleaning failed ❌"; exit 1; }

dmg: build
	@./macos_build.sh