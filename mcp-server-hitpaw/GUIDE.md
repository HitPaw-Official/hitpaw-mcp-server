# HitPaw MCP Server - 完整开发部署指南

## 一、核心概念：MCP 和你现有 API 的关系

```
┌─────────────────────────────────────────────────────────────────┐
│                     你不需要改动的部分                            │
│                                                                  │
│   你现有的 HTTP API 服务器（api-service）                         │
│   运行在: https://api-base.hitpaw.com                            │
│   ├── POST /api/photo-enhancer     ← 图片增强                   │
│   ├── POST /api/video-enhancer     ← 视频增强                   │
│   ├── POST /api/task-status        ← 任务状态                   │
│   ├── POST /api/oss/upload         ← OSS上传                    │
│   ├── POST /api/oss/transfer       ← OSS转存                    │
│   └── ...                                                        │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
                              ▲
                              │ HTTP 请求（通过 X-API-KEY 认证）
                              │
┌─────────────────────────────────────────────────────────────────┐
│                     你需要新开发的部分                            │
│                                                                  │
│   MCP Server（本项目 mcp-server-hitpaw）                         │
│   ├── 是一个 CLI 命令行程序（不是 web 服务）                      │
│   ├── 运行在【用户的本地电脑】上（不是你的服务器）                  │
│   ├── 通过 stdin/stdout 和 Claude 通信（JSON-RPC 2.0）           │
│   ├── 收到 Claude 的工具调用请求后                                │
│   │   → 转化为 HTTP 请求，发送给 api-base.hitpaw.com             │
│   │   → 把结果返回给 Claude                                     │
│   └── 用户通过 npm install 安装                                  │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
                              ▲
                              │ stdin/stdout（JSON-RPC 2.0）
                              │
┌─────────────────────────────────────────────────────────────────┐
│   Claude Desktop / Cursor IDE / 其他 MCP 客户端                  │
│   运行在用户的电脑上                                              │
└─────────────────────────────────────────────────────────────────┘
```

### 常见问题

**Q: 需要重新部署我的 API 服务吗？**  
A: 不需要。MCP Server 只是一个 HTTP 客户端，调用你已有的 API。你的 API 服务保持原样运行。

**Q: MCP Server 部署在哪里？**  
A: 运行在用户的本地电脑上。用户通过 npm install 安装，Claude Desktop 自动启动它。

**Q: 需要独立运行一个 web 服务吗？**  
A: 不需要。它是 CLI 程序，通过 stdin/stdout 和 Claude 通信，不监听任何端口。

---

## 二、开发流程

### 2.1 环境准备

```bash
# 确保安装了 Go 1.24+
go version
# go version go1.24.x linux/amd64

# 确保安装了 Node.js 16+（发布 npm 包用）
node --version
npm --version
```

### 2.2 初始化项目

```bash
cd mcp-server-hitpaw

# 初始化模块依赖
go mod tidy

# 编译
make build
# 产物: build/hitpaw-mcp-server
```

### 2.3 本地测试

```bash
# 设置环境变量
export HITPAW_API_KEY=your_test_api_key
export HITPAW_API_BASE_URL=https://api-base.hitpaw.com

# 启动（会阻塞等待 stdin 输入）
./build/hitpaw-mcp-server
```

#### 手动发送 JSON-RPC 测试

```bash
# 测试初始化
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}' | HITPAW_API_KEY=test ./build/hitpaw-mcp-server

# 测试获取工具列表
echo '{"jsonrpc":"2.0","id":2,"method":"tools/list"}' | HITPAW_API_KEY=test ./build/hitpaw-mcp-server

# 测试调用模型列表
echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"list_photo_models","arguments":{}}}' | HITPAW_API_KEY=test ./build/hitpaw-mcp-server
```

### 2.4 在 Claude Desktop 中调试

编辑配置指向本地二进制：

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`  
**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "hitpaw": {
      "command": "/你的绝对路径/mcp-server-hitpaw/build/hitpaw-mcp-server",
      "env": {
        "HITPAW_API_KEY": "your_api_key",
        "HITPAW_API_BASE_URL": "https://api-base.hitpaw.com"
      }
    }
  }
}
```

重启 Claude Desktop → 设置里应该能看到 hitpaw 工具列表。

---

## 三、编译和发布

### 3.1 跨平台编译

```bash
bash scripts/build.sh 1.0.0

# 产物:
# build/hitpaw-mcp-server-darwin-arm64    (macOS Apple Silicon)
# build/hitpaw-mcp-server-darwin-amd64    (macOS Intel)
# build/hitpaw-mcp-server-linux-amd64     (Linux x64)
# build/hitpaw-mcp-server-linux-arm64     (Linux ARM)
# build/hitpaw-mcp-server-windows-amd64.exe (Windows)
```

### 3.2 发布到 GitHub Releases

```bash
# 安装 GitHub CLI
brew install gh  # 或 apt install gh

# 初始化仓库并推送
git init
git add .
git commit -m "feat: HitPaw MCP Server v1.0.0"
git remote add origin https://github.com/hitpaw/mcp-server-hitpaw.git
git push -u origin main

# 创建 Release 并上传所有二进制
gh release create v1.0.0 build/* --title "v1.0.0" --notes "Initial release"
```

### 3.3 发布到 npm

```bash
# 登录 npm（需要 npm 账号）
npm login

# 如果使用 @hitpaw 组织名，需要先在 npmjs.com 创建组织
# 发布
cd npm
npm publish --access public
```

---

## 四、用户安装和使用流程

### 完整用户旅程

```
步骤1: 获取 API Key
  └── 在 HitPaw 官网注册 → 开发者后台 → 创建 API Key

步骤2: 配置 Claude Desktop（一次性）
  └── 编辑 claude_desktop_config.json → 重启

步骤3: 直接在 Claude 中对话使用
  └── "帮我把这张照片变清晰 https://..."
```

### 用户配置方式

#### 方式 A: npx（推荐，无需安装）

```json
{
  "mcpServers": {
    "hitpaw": {
      "command": "npx",
      "args": ["-y", "@hitpaw/mcp-server"],
      "env": {
        "HITPAW_API_KEY": "用户的api_key"
      }
    }
  }
}
```

#### 方式 B: npm 全局安装

```bash
npm install -g @hitpaw/mcp-server
```

```json
{
  "mcpServers": {
    "hitpaw": {
      "command": "hitpaw-mcp-server",
      "env": {
        "HITPAW_API_KEY": "用户的api_key"
      }
    }
  }
}
```

#### 方式 C: 直接下载二进制

从 GitHub Releases 下载对应平台文件 → 放到 PATH 中。

---

## 五、推广你的 MCP 服务

### 5.1 提交到 MCP 官方目录

1. Fork https://github.com/modelcontextprotocol/servers
2. 在 README 中添加你的 server 信息
3. 提交 Pull Request

### 5.2 npm 自然搜索

发布到 npm 后，用户搜 "mcp server image enhance" 等关键词可以找到。

### 5.3 产品文档集成

在 HitPaw 官网开发者文档中添加 "MCP 集成指南" 章节。

### 5.4 社区

- Anthropic Discord
- GitHub Discussions
- 技术博客

---

## 六、你的 API 需要注意的事项

### 6.1 认证方式

MCP Server 在请求头带 `X-API-KEY`。确保你的 API 中间件能识别：

```go
// 你现有代码中的 ext.ApiKey 就是从 Header 取的
ApiKey := c.GetHeader(ext.ApiKey)
```

如果 Header 名称不同，需要修改 `internal/client/api_client.go` 中的：

```go
httpReq.Header.Set("X-API-KEY", c.apiKey)
```

### 6.2 响应格式

MCP Server 期望你的 API 返回：

```json
{
  "code": 0,
  "msg": "success",
  "data": { ... }
}
```

如果格式不同，修改 `api_client.go` 的 `APIResponse` 结构体。

### 6.3 不需要 CORS

MCP Server 是本地 CLI 程序，不是浏览器，所以不需要 CORS 配置。

---

## 七、版本更新

```bash
# 1. 修改版本号（三处）
#    - cmd/mcp-server/main.go → serverVersion
#    - npm/package.json → version
#    - npm/install.js → VERSION

# 2. 编译 + 发布
bash scripts/build.sh 1.1.0
git add . && git commit -m "feat: v1.1.0"
git tag v1.1.0 && git push && git push --tags
gh release create v1.1.0 build/*
cd npm && npm publish --access public
```

---

## 八、项目文件清单

```
mcp-server-hitpaw/
├── go.mod                              # Go 模块定义（Go 1.24）
├── Makefile                            # 编译命令
├── .gitignore                          # Git 忽略规则
├── README.md                           # 项目说明
├── GUIDE.md                            # 本文档
├── cmd/
│   └── mcp-server/
│       └── main.go                     # 程序入口
├── configs/
│   └── config.go                       # 环境变量配置
├── internal/
│   ├── protocol/
│   │   ├── types.go                    # MCP 协议类型（JSON-RPC 2.0）
│   │   └── server.go                   # MCP Server 核心（stdin/stdout）
│   ├── client/
│   │   └── api_client.go              # HTTP 客户端（调用 api-base.hitpaw.com）
│   └── handler/
│       ├── tools.go                    # 7个工具注册
│       ├── photo_enhancer.go           # 图片增强处理
│       ├── video_enhancer.go           # 视频增强处理
│       ├── task_status.go              # 任务状态查询
│       ├── oss.go                      # OSS 文件操作
│       └── models.go                   # 模型列表查询
├── npm/
│   ├── package.json                    # npm 包配置
│   ├── index.js                        # npm 入口（启动 Go 二进制）
│   └── install.js                      # 安装脚本（下载二进制）
└── scripts/
    └── build.sh                        # 跨平台编译脚本
```
