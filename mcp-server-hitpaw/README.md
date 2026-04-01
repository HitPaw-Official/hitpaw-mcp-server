# HitPaw MCP Server

让 Claude 直接调用 HitPaw AI 图片/视频增强服务的 MCP（Model Context Protocol）服务器。

## 功能

- 🖼️ **图片增强** - 16种AI模型，支持人像美颜、通用超分、降噪、生成式修复
- 🎬 **视频增强** - 8种AI模型，支持人像修复、通用修复、超高清修复
- 📋 **任务查询** - 实时查询处理状态和结果
- 📁 **文件管理** - URL转存到OSS，批量转存
- 📖 **模型列表** - 查看所有可用模型和使用建议

## 快速开始

### 方式一：npx 直接运行（推荐）

```bash
export HITPAW_API_KEY=your_api_key_here
# 可选：
# export HITPAW_API_BASE_URL=https://api-base.hitpaw.com
npx @hitpaw/mcp-server
```

### 方式二：npm 全局安装

```bash
npm install -g @hitpaw/mcp-server
hitpaw-mcp-server
```

### 方式三：从源码编译

```bash
git clone https://github.com/hitpaw/mcp-server-hitpaw.git
cd mcp-server-hitpaw
make build
./build/hitpaw-mcp-server
```

## 配置 Claude Desktop

编辑配置文件：

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

从源码编译的可直接指向二进制：

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

配置完成后重启 Claude Desktop。

## 配置 Cursor IDE

编辑 `~/.cursor/mcp.json`：

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

## 环境变量

| 变量名 | 必填 | 默认值 | 说明 |
|--------|------|--------|------|
| `HITPAW_API_KEY` | ✅ | - | HitPaw API 密钥 |
| `HITPAW_API_BASE_URL` | ❌ | `https://api-base.hitpaw.com` | API 服务地址 |

## 使用示例

在 Claude 中对话：

```
用户: 帮我把这张照片变清晰 https://example.com/photo.jpg

Claude: 我来帮您增强这张图片。先查看一下可用的模型...
[调用 list_photo_models]
这张图片看起来是人像照片，我推荐使用 face_2x 模型进行2倍清晰化。
[调用 photo_enhance(model_name="face_2x", img_url="https://example.com/photo.jpg")]
任务已创建，任务ID为 xxx-xxx。让我查询一下处理状态...
[调用 task_status(job_id="xxx-xxx")]
处理完成！这是增强后的图片链接: https://...
```

## 可用工具

| 工具名 | 说明 |
|--------|------|
| `photo_enhance` | 图片增强/超分辨率 |
| `video_enhance` | 视频增强/超分辨率 |
| `task_status` | 查询任务状态 |
| `oss_transfer` | URL文件转存到OSS |
| `oss_batch_transfer` | 批量URL转存 |
| `list_photo_models` | 列出图片增强模型 |
| `list_video_models` | 列出视频增强模型 |

## 图片增强模型

### 人像模型
| 模型 | 放大 | 说明 |
|------|------|------|
| `face_2x` | 2x | 人像清晰柔和（美颜效果） |
| `face_4x` | 4x | 人像清晰柔和（美颜效果） |
| `face_v2_2x` | 2x | 人像自然真实（保留纹理） |
| `face_v2_4x` | 4x | 人像自然真实（保留纹理） |

### 通用模型
| 模型 | 放大 | 说明 |
|------|------|------|
| `general_2x` | 2x | 通用高清增强 |
| `general_4x` | 4x | 通用高清增强 |
| `high_fidelity_2x` | 2x | 高保真（保留原图风格） |
| `high_fidelity_4x` | 4x | 高保真（保留原图风格） |

### 降噪模型
| 模型 | 放大 | 说明 |
|------|------|------|
| `sharpen_denoise_1x` | 1x | 锐化降噪 |
| `detail_denoise_1x` | 1x | 极致保真降噪 |

### 生成式模型（SD）
| 模型 | 放大 | 说明 |
|------|------|------|
| `generative_portrait_1x/2x/4x` | 1x/2x/4x | 生成式人像修复 |
| `generative_1x/2x/4x` | 1x/2x/4x | 生成式通用修复 |

## 视频增强模型

| 模型 | 说明 |
|------|------|
| `face_soft_2x` | 人脸柔和增强2x |
| `portrait_restore_1x/2x` | 人像修复 |
| `general_restore_1x/2x/4x` | 通用修复 |
| `ultrahd_restore_2x` | 超高清修复（SD） |
| `generative_1x` | 生成式增强（SD） |

## 开发

```bash
make build          # 编译
make build-all      # 跨平台编译
make test           # 测试
HITPAW_API_KEY=xxx make run  # 本地运行
```

## 发布

```bash
# 1. 跨平台编译
bash scripts/build.sh 1.0.0

# 2. 创建 GitHub Release
gh release create v1.0.0 build/*

# 3. 发布 npm 包
cd npm && npm publish --access public
```

## License

MIT
