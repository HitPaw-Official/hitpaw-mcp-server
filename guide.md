# HitPaw MCP Server - Complete Development and Deployment Guide

## 1. Core Concept: The Relationship Between MCP and Your Existing API

```
┌─────────────────────────────────────────────────────────────────┐
│                     Parts You Don't Need to Change              │
│                                                                  │
│   Your existing HTTP API server (api-service)                    │
│   Running at: https://api-base.hitpaw.com                        │
│   ├── POST /api/photo-enhancer     ← Photo Enhancement          │
│   ├── POST /api/video-enhancer     ← Video Enhancement          │
│   ├── POST /api/task-status        ← Task Status                │
│   ├── POST /api/oss/upload         ← OSS Upload                 │
│   ├── POST /api/oss/transfer       ← OSS Transfer               │
│   └── ...                                                        │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
                              ▲
                              │ HTTP Requests (X-API-KEY Auth)
                              │
┌─────────────────────────────────────────────────────────────────┐
│                     Parts You Need to Develop                    │
│                                                                  │
│   MCP Server (This project mcp-server-hitpaw)                    │
│   ├── A CLI command-line program (not a web service)             │
│   ├── Runs on the [User's Local Machine] (not your server)       │
│   ├── Communicates with Claude via stdin/stdout (JSON-RPC 2.0)   │
│   ├── Upon receiving tool call requests from Claude:             │
│   │   → Converts to HTTP requests, sent to api-base.hitpaw.com   │
│   │   → Returns the result to Claude                             │
│   └── Users install it via npm install                           │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
                              ▲
                              │ stdin/stdout (JSON-RPC 2.0)
                              │
┌─────────────────────────────────────────────────────────────────┐
│   Claude Desktop / Cursor IDE / Other MCP Clients                │
│   Running on the User's Local Machine                            │
└─────────────────────────────────────────────────────────────────┘
```

### Frequently Asked Questions

**Q: Do I need to redeploy my API service?**  
A: No. The MCP Server acts merely as an HTTP client calling your existing APIs. Your API service remains exactly as it is.

**Q: Where is the MCP Server deployed?**  
A: It runs on the user's local machine. The user installs it via `npm install`, and Claude Desktop starts it automatically.

**Q: Does it need to run as an independent web service?**  
A: No. It's a CLI program that communicates with Claude via stdin/stdout, and does not listen on any ports.

---

## 2. Development Workflow

### 2.1 Environmental Preparation

```bash
# Ensure Go 1.24+ is installed
go version
# go version go1.24.x linux/amd64

# Ensure Node.js 16+ is installed (for publishing the npm package)
node --version
npm --version
```

### 2.2 Project Initialization

```bash
cd mcp-server-hitpaw

# Initialize module dependencies
go mod tidy

# Build
make build
# Output artifact: build/hitpaw-mcp-server
```

### 2.3 Local Testing

```bash
# Set environment variables
export HITPAW_API_KEY=your_test_api_key
export HITPAW_API_BASE_URL=https://api-base.hitpaw.com

# Start (will block waiting for stdin input)
./build/hitpaw-mcp-server
```

#### Manually Send JSON-RPC Tests

```bash
# Test initialization
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}' | HITPAW_API_KEY=test ./build/hitpaw-mcp-server

# Test tool listing
echo '{"jsonrpc":"2.0","id":2,"method":"tools/list"}' | HITPAW_API_KEY=test ./build/hitpaw-mcp-server

# Test tool calling
echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"list_photo_models","arguments":{}}}' | HITPAW_API_KEY=test ./build/hitpaw-mcp-server
```

### 2.4 Debugging in Claude Desktop

Edit the configuration file to point to your local binary:

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`  
**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "hitpaw": {
      "command": "/your_absolute_path/mcp-server-hitpaw/build/hitpaw-mcp-server",
      "env": {
        "HITPAW_API_KEY": "your_api_key",
        "HITPAW_API_BASE_URL": "https://api-base.hitpaw.com"
      }
    }
  }
}
```

Restart Claude Desktop → You should see the hitpaw tools listed in the tools menu.

---

## 3. Build and Publish

### 3.1 Cross-Platform Compilation

```bash
bash scripts/build.sh 1.0.0

# Output artifacts:
# build/hitpaw-mcp-server-darwin-arm64    (macOS Apple Silicon)
# build/hitpaw-mcp-server-darwin-amd64    (macOS Intel)
# build/hitpaw-mcp-server-linux-amd64     (Linux x64)
# build/hitpaw-mcp-server-linux-arm64     (Linux ARM)
# build/hitpaw-mcp-server-windows-amd64.exe (Windows)
```

### 3.2 Publish to GitHub Releases

```bash
# Install GitHub CLI
brew install gh  # or apt install gh

# Initialize repository and push
git init
git add .
git commit -m "feat: HitPaw MCP Server v1.0.0"
git remote add origin https://github.com/hitpaw/mcp-server-hitpaw.git
git push -u origin main

# Create Release and upload all binaries
gh release create v1.0.0 build/* --title "v1.0.0" --notes "Initial release"
```

### 3.3 Publish to npm

```bash
# Login to npm (npm account required)
npm login

# If using the @hitpaw organization name, you must create the organization on npmjs.com first
# Publish
cd npm
npm publish --access public
```

---

## 4. User Installation & Usage Flow

### Complete User Journey

```
Step 1: Obtain API Key
  └── Register on HitPaw website → Developer Console → Create API Key

Step 2: Configure Claude Desktop (One-time)
  └── Edit claude_desktop_config.json → Restart

Step 3: Direct Usage within Claude Chat
  └── "Please make this photo clearer: https://..."
```

### User Configuration Methods

#### Method A: npx (Recommended, no install needed)

```json
{
  "mcpServers": {
    "hitpaw": {
      "command": "npx",
      "args": ["-y", "@hitpaw/mcp-server"],
      "env": {
        "HITPAW_API_KEY": "user_api_key"
      }
    }
  }
}
```

#### Method B: Global npm install

```bash
npm install -g @hitpaw/mcp-server
```

```json
{
  "mcpServers": {
    "hitpaw": {
      "command": "hitpaw-mcp-server",
      "env": {
        "HITPAW_API_KEY": "user_api_key"
      }
    }
  }
}
```

#### Method C: Download Binary Directly

Download the binary for the corresponding platform from GitHub Releases → place it in your `PATH`.

---

## 5. Promoting Your MCP Service

### 5.1 Submit to the Official MCP Directory

1. Fork https://github.com/modelcontextprotocol/servers
2. Add your server's details to the README
3. Submit a Pull Request

### 5.2 Organic npm Search

After publishing to npm, users can find it by searching for keywords like "mcp server image enhance".

### 5.3 Product Documentation Integration

Add an "MCP Integration Guide" section to the HitPaw official developer documentation.

### 5.4 Communities

- Anthropic Discord
- GitHub Discussions
- Tech Blogs

---

## 6. Important Notes for Your API

### 6.1 Authentication Method

The MCP Server appends `X-API-KEY` to the request header. Ensure your API middleware accurately reads it:

```go
// Extracted from headers in your existing code
ApiKey := c.GetHeader("X-API-KEY")
```

If the header name is different, you'll need to modify it inside `internal/client/api_client.go`:

```go
httpReq.Header.Set("X-API-KEY", c.apiKey)
```

### 6.2 Response Format

The MCP Server expects your API to return:

```json
{
  "code": 0,
  "msg": "success",
  "data": { ... }
}
```

If the format is different, modify the `APIResponse` struct in `api_client.go`.

### 6.3 No CORS Required

The MCP Server is a local CLI program, not a browser, so CORS configuration is not necessary.

---

## 7. Version Updates

```bash
# 1. Update the version number (in three places)
#    - cmd/mcp-server/main.go → serverVersion
#    - npm/package.json → version
#    - npm/install.js → VERSION

# 2. Build + Publish
bash scripts/build.sh 1.1.0
git add . && git commit -m "feat: v1.1.0"
git tag v1.1.0 && git push && git push --tags
gh release create v1.1.0 build/*
cd npm && npm publish --access public
```

---

## 8. Project File Structure

```
mcp-server-hitpaw/
├── go.mod                              # Go module definition (Go 1.24)
├── Makefile                            # Compilation commands
├── .gitignore                          # Git ignore rules
├── README.md                           # Project Description
├── GUIDE.md                            # This Guide
├── cmd/
│   └── mcp-server/
│       └── main.go                     # Program entry point
├── configs/
│   └── config.go                       # Environment variable config
├── internal/
│   ├── protocol/
│   │   ├── types.go                    # MCP protocol types (JSON-RPC 2.0)
│   │   └── server.go                   # MCP Server core (stdin/stdout)
│   ├── client/
│   │   └── api_client.go              # HTTP client (calls api-base.hitpaw.com)
│   └── handler/
│       ├── tools.go                    # Registering 7 tools
│       ├── photo_enhancer.go           # Photo enhancer handling
│       ├── video_enhancer.go           # Video enhancer handling
│       ├── task_status.go              # Task status query
│       ├── oss.go                      # OSS file operations
│       └── models.go                   # Model list query
├── npm/
│   ├── package.json                    # npm package configuration
│   ├── index.js                        # npm entry (starts Go binary)
│   └── install.js                      # install script (downloads binary)
└── scripts/
    └── build.sh                        # Cross-platform build script
```
