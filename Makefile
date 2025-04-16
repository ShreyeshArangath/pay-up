# Variables
BINARY_NAME=payup
BUILD_DIR=./cmd/$(BINARY_NAME)
INSTALL_PATH=/usr/local/bin

# Build the executable
build:
	go build -o $(BINARY_NAME) $(BUILD_DIR)

# Install the executable to /usr/local/bin
install: build
	mv $(BINARY_NAME) $(INSTALL_PATH)

# Clean up the binary
clean:
	rm -f $(BINARY_NAME)

# Run the program
run:
	go run $(BUILD_DIR)/main.go
