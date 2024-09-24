# Define variables
BINARY_NAME = tracker
BUILD_DIR = bin
PKG_ROOT = pkgroot
PKG_FILE = $(BINARY_NAME).pkg
INSTALL_LOCATION = /usr/local/bin
VERSION = 1.0
IDENTIFIER = com.yourcompany.tracker

# Build for macOS
build:
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME) cmd/main.go

# Run the Go application
run:
	go run cmd/main.go

# Create directory structure for packaging
pkgdir:
	mkdir -p $(PKG_ROOT)$(INSTALL_LOCATION)
	cp $(BUILD_DIR)/$(BINARY_NAME) $(PKG_ROOT)$(INSTALL_LOCATION)/

# Package the application into a .pkg file
package: build pkgdir
	pkgbuild --root $(PKG_ROOT) --identifier $(IDENTIFIER) --version $(VERSION) --install-location $(INSTALL_LOCATION) $(PKG_FILE)

# Clean up temporary files
clean:
	rm -rf $(BUILD_DIR) $(PKG_ROOT) $(PKG_FILE)

# Install dependencies
deps:
	go mod tidy

# Run tests
test:
	go test ./... -v

# Default target
.DEFAULT_GOAL := build

.PHONY: build run pkgdir package clean deps test