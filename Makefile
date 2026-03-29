BINARY := pls
BINARY_PATH := bin/$(BINARY)
INSTALL_DIR ?= $(HOME)/.local/bin
GO ?= go

.PHONY: build test install uninstall print-config-path doctor go-version release-snapshot

build:
	$(GO) build -o $(BINARY_PATH) ./cmd/pls

test:
	$(GO) test ./...

install:
	mkdir -p $(INSTALL_DIR)
	GOBIN=$(INSTALL_DIR) $(GO) install ./cmd/pls
	@echo "Installed $(BINARY) to $(INSTALL_DIR)/$(BINARY)"
	@echo 'If that directory is not in PATH, add: export PATH="$$HOME/.local/bin:$$PATH"'

uninstall:
	rm -f $(INSTALL_DIR)/$(BINARY)

print-config-path:
	$(GO) run ./cmd/pls --print-config-path

doctor:
	$(GO) run ./cmd/pls doctor

go-version:
	$(GO) version

release-snapshot:
	goreleaser release --snapshot --clean --config .goreleaser.yaml --skip=publish,announce,sign
