BINARY_NAME=rproxctl
BUILD_DIR=cmd/$(BINARY_NAME)
INSTALL_PATH=$(HOME)/go/bin

.PHONY: all build install clean

all: build install

build:
	@echo "Building $(BINARY_NAME)..."
	go build -mod=vendor -o $(BINARY_NAME) ./$(BUILD_DIR)

install:
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	mkdir -p $(INSTALL_PATH)
	cp $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Installed successfully. Make sure $(INSTALL_PATH) is in your PATH."

clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
