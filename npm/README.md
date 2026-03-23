# @hitpaw/mcp-server

An MCP (Model Context Protocol) server that brings HitPaw's powerful AI image and video enhancement capabilities directly to Claude Desktop and Cursor IDE.

## Features

- 🖼️ **Photo Enhancement**: 16 AI models including portrait beautification, general super-resolution, and generative restoration.
- 🎬 **Video Enhancement**: 8 AI models for portrait, general, and ultra-HD video restoration.
- 📋 **Task Management**: Check processing status and results in real-time.
- 📁 **OSS Integration**: Transfer URLs straight to OSS.

## Installation & Usage

You can use this MCP server directly with `npx` (no installation required) or install it globally.

### Prerequisites

You need a **HitPaw API Key**. Register at the [HitPaw Developer Console](https://api-base.hitpaw.com) to obtain one.

### For Claude Desktop

Edit your Claude Desktop configuration file:
- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

Add the following configuration:

```json
{
  "mcpServers": {
    "hitpaw": {
      "command": "npx",
      "args": ["-y", "@hitpaw/mcp-server"],
      "env": {
        "HITPAW_API_KEY": "your_api_key_here"
      }
    }
  }
}
```

Then restart Claude Desktop.

### Global Installation

```bash
npm install -g @hitpaw/mcp-server
```

When installed globally, you can configure it like this:

```json
{
  "mcpServers": {
    "hitpaw": {
      "command": "hitpaw-mcp-server",
      "env": {
        "HITPAW_API_KEY": "your_api_key_here"
      }
    }
  }
}
```

## Available Tool Requests (Inside Claude/Cursor)

Once configured, your AI assistant can invoke the following tools natively on your behalf:
- `photo_enhance`
- `video_enhance`
- `task_status`
- `oss_transfer`
- `oss_batch_transfer`
- `list_photo_models`
- `list_video_models`

For instance, you can simply ask: *"Can you enhance this portrait image for me? [image_url]"*

## Links
- [HitPaw GitHub Repository](https://github.com/HitPaw-Official/hitpaw-mcp-server)

## License

MIT
