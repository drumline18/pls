BINARY := pls
BINARY_PATH := bin/$(BINARY)
INSTALL_DIR ?= $(HOME)/.local/bin

.PHONY: build test install uninstall print-config-path

build:
	go build -o $(BINARY_PATH) ./cmd/pls

test:
	go test ./...

install:
	mkdir -p $(INSTALL_DIR)
	GOBIN=$(INSTALL_DIR) go install ./cmd/pls
	@echo "Installed $(BINARY) to $(INSTALL_DIR)/$(BINARY)"
	@echo 'If that directory is not in PATH, add: export PATH="$HOME/.local/bin:$PATH"'

uninstall:
	rm -f $(INSTALL_DIR)/$(BINARY)

print-config-path:
	go run ./cmd/pls --print-config-path
