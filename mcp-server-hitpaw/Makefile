.PHONY: build build-all clean test run install

BINARY_NAME=hitpaw-mcp-server
VERSION=1.0.0
BUILD_DIR=build
LDFLAGS=-ldflags "-s -w -X main.serverVersion=$(VERSION)"

# 本地编译
build:
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/mcp-server/

# 跨平台编译（用于 npm 发布）
build-all: clean
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/mcp-server/
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/mcp-server/
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/mcp-server/
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/mcp-server/
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/mcp-server/
	@echo "✅ 跨平台编译完成，产物在 $(BUILD_DIR)/ 目录"

clean:
	rm -rf $(BUILD_DIR)

test:
	go test ./...

run:
	go run ./cmd/mcp-server/

install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@echo "✅ 已安装到 /usr/local/bin/$(BINARY_NAME)"
