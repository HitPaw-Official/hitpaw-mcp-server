# HitPaw MCP Server

An MCP (Model Context Protocol) server that allows Claude to directly invoke the HitPaw AI photo/video enhancement services.

## Features

- 🖼️ **Photo Enhancement** - 16 AI models, supporting portrait beautification, general super-resolution, denoising, and generative restoration
- 🎬 **Video Enhancement** - 8 AI models, supporting portrait restoration, general restoration, and ultra-HD restoration
- 📋 **Task Query** - Real-time query of processing status and results
- 📁 **File Management** - Save URLs to OSS, batch transfer
- 📖 **Model List** - View all available models and usage suggestions

## Quick Start

### Method 1: Run directly with npx (Recommended)

```bash
export HITPAW_API_KEY=your_api_key_here
# Optional:
# export HITPAW_API_BASE_URL=https://api-base.hitpaw.com
npx @hitpaw/mcp-server
```

### Method 2: Global installation via npm

```bash
npm install -g @hitpaw/mcp-server
hitpaw-mcp-server
```

### Method 3: Build from source

```bash
git clone https://github.com/hitpaw/mcp-server-hitpaw.git
cd mcp-server-hitpaw
make build
./build/hitpaw-mcp-server
```

## Configure Claude Desktop

Edit the configuration file:

- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "hitpaw": {
      "command": "npx",
      "args": ["-y", "@hitpaw/mcp-server"],
      "env": {
        "HITPAW_API_KEY": "your_api_key_here",
        "HITPAW_API_BASE_URL": "https://api-base.hitpaw.com"
      }
    }
  }
}
```

If built from source, you can point directly to the binary:

```json
{
  "mcpServers": {
    "hitpaw": {
      "command": "/path/to/hitpaw-mcp-server",
      "env": {
        "HITPAW_API_KEY": "your_api_key_here",
        "HITPAW_API_BASE_URL": "https://api-base.hitpaw.com"
      }
    }
  }
}
```

Restart Claude Desktop after configuration.

## Configure Cursor IDE

Edit `~/.cursor/mcp.json`:

```json
{
  "mcpServers": {
    "hitpaw": {
      "command": "npx",
      "args": ["-y", "@hitpaw/mcp-server"],
      "env": {
        "HITPAW_API_KEY": "your_api_key_here",
        "HITPAW_API_BASE_URL": "https://api-base.hitpaw.com"
      }
    }
  }
}
```

## Environment Variables

| Variable | Required | Default | Description |
|--------|------|--------|------|
| `HITPAW_API_KEY` | ✅ | - | HitPaw API Key |
| `HITPAW_API_BASE_URL` | ❌ | `https://api-base.hitpaw.com` | API Service URL |

## Usage Examples

Chat in Claude:

```
User: Please enhance this photo https://example.com/photo.jpg

Claude: I'll help you enhance this image. Let me check the available models first...
[Calls list_photo_models]
This photo looks like a portrait, I recommend using the face_2x model to upscale it by 2x.
[Calls photo_enhance(model_name="face_2x", img_url="https://example.com/photo.jpg")]
Task created, task ID is xxx-xxx. Let me check the processing status...
[Calls task_status(job_id="xxx-xxx")]
Processing complete! Here is the link to the enhanced photo: https://...
```

## Available Tools

| Tool Name | Description |
|--------|------|
| `photo_enhance` | Photo enhancement/super resolution |
| `video_enhance` | Video enhancement/super resolution |
| `task_status` | Query task status |
| `oss_transfer` | Transfer URL file to OSS |
| `oss_batch_transfer` | Batch transfer URLs to OSS |
| `list_photo_models` | List photo enhancement models |
| `list_video_models` | List video enhancement models |

## Photo Enhancement Models

### Portrait Models
| Model | Upscale | Description |
|------|------|------|
| `face_2x` | 2x | Portrait clear and soft (beautification effect) |
| `face_4x` | 4x | Portrait clear and soft (beautification effect) |
| `face_v2_2x` | 2x | Portrait natural and realistic (preserves texture) |
| `face_v2_4x` | 4x | Portrait natural and realistic (preserves texture) |

### General Models
| Model | Upscale | Description |
|------|------|------|
| `general_2x` | 2x | General HD enhancement |
| `general_4x` | 4x | General HD enhancement |
| `high_fidelity_2x` | 2x | High fidelity (preserves original style) |
| `high_fidelity_4x` | 4x | High fidelity (preserves original style) |

### Denoising Models
| Model | Upscale | Description |
|------|------|------|
| `sharpen_denoise_1x` | 1x | Sharpen and denoise |
| `detail_denoise_1x` | 1x | Ultimate fidelity denoising |

### Generative Models (SD)
| Model | Upscale | Description |
|------|------|------|
| `generative_portrait_1x/2x/4x` | 1x/2x/4x | Generative portrait restoration |
| `generative_1x/2x/4x` | 1x/2x/4x | Generative general restoration |

## Video Enhancement Models

| Model | Description |
|------|------|
| `face_soft_2x` | Face soft enhancement 2x |
| `portrait_restore_1x/2x` | Portrait restoration |
| `general_restore_1x/2x/4x` | General restoration |
| `ultrahd_restore_2x` | Ultra-HD restoration (SD) |
| `generative_1x` | Generative enhancement (SD) |

## Development

```bash
make build          # Build
make build-all      # Cross-platform build
make test           # Test
HITPAW_API_KEY=xxx make run  # Run locally
```

## Release

```bash
# 1. Cross-platform build
bash scripts/build.sh 1.0.0

# 2. Create GitHub Release
gh release create v1.0.0 build/*

# 3. Publish npm package
cd npm && npm publish --access public
```

## License

MIT
