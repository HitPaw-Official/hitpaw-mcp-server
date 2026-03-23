package handler

import (
	"github.com/hitpaw/mcp-server-hitpaw/internal/client"
	"github.com/hitpaw/mcp-server-hitpaw/internal/protocol"
)

// RegisterAllTools 注册所有工具到 MCP Server
func RegisterAllTools(server *protocol.MCPServer, apiClient *client.APIClient) {
	h := &Handlers{api: apiClient}

	server.RegisterTool(photoEnhanceTool(), h.HandlePhotoEnhance)
	server.RegisterTool(videoEnhanceTool(), h.HandleVideoEnhance)
	server.RegisterTool(taskStatusTool(), h.HandleTaskStatus)
	server.RegisterTool(ossTransferTool(), h.HandleOSSTransfer)
	server.RegisterTool(ossBatchTransferTool(), h.HandleOSSBatchTransfer)
	server.RegisterTool(listPhotoModelsTool(), h.HandleListPhotoModels)
	server.RegisterTool(listVideoModelsTool(), h.HandleListVideoModels)
}

// Handlers 工具处理器集合
type Handlers struct {
	api *client.APIClient
}

// ==================== 工具定义 ====================

func photoEnhanceTool() protocol.Tool {
	minUpscale := float64(1)
	maxUpscale := float64(4)
	return protocol.Tool{
		Name: "photo_enhance",
		Description: `图片增强/超分辨率工具 - 提升图片清晰度和分辨率。返回 job_id 用于查询任务状态。

━━━ 模型选择决策树（必须严格遵循） ━━━

第一步：判断图片内容
├── 含人脸 → 转第二步
└── 无人脸（风景/建筑/物品/动物）→ 转第三步

第二步：人脸图片选模型
├── 极低质量/大头照/侧脸/老照片/表情包 → generative_portrait 系列（SD生成式补全）
│   ├── 不放大 → generative_portrait_1x
│   ├── 2倍放大 → generative_portrait_2x
│   └── 4倍放大 → generative_portrait_4x
├── 用户要求"保留纹理/自然/不磨皮/保留皱纹胡须" → face_v2 系列（自然纹理保真）
│   ├── 2倍放大或未指定 → face_v2_2x
│   └── 4倍放大 → face_v2_4x
└── 其他人脸需求（默认/美颜/柔和/清晰）→ face 系列（柔和美颜）
    ├── 2倍放大或未指定 → face_2x
    └── 4倍放大 → face_4x

第三步：通用图片选模型
├── 仅降噪不放大
│   ├── "锐化/手机抓拍" → sharpen_denoise_1x（锐化降噪，生成新纹理）
│   └── "保真/保留细节/夜景/复杂噪声" → detail_denoise_1x（保真降噪）
├── 极低质量图片需生成补全
│   ├── 要求"快速/效率优先"
│   │   ├── 2倍 → generative_enhance_fast_2x
│   │   └── 4倍 → generative_enhance_fast_4x
│   ├── 不放大 → generative_1x
│   ├── 2倍放大 → generative_2x
│   └── 4倍放大 → generative_4x
├── "高保真/保留风格/摄影作品/艺术照"
│   ├── 2倍放大 → high_fidelity_2x
│   └── 4倍放大 → high_fidelity_4x
└── 常规清晰化/画质提升（默认）
    ├── 2倍放大或未指定 → general_2x
    └── 4倍放大 → general_4x

━━━ 各模型详细说明 ━━━

【人像清晰柔和系列 - 柔和美颜效果】
- face_2x: 人脸+背景双模型，2倍放大。适合模糊自拍/街拍/合影/证件照。避免过度锐化导致假脸。用户提"清晰/修复/美颜/自拍/老照片/人像高清"且2倍或未指定倍数时使用。
- face_4x: 同上但4倍放大。仅当用户明确要求"4倍放大/大幅放大"时使用。

【人像自然真实系列 - 保留面部真实纹理】
- face_v2_2x: 人脸v2+背景双模型，2倍放大。保留皱纹/胡须/装饰等特征。用户提"保留纹理/自然质感/不磨皮"且2倍或未指定时使用。
- face_v2_4x: 同上但4倍放大。仅当用户明确要求4倍放大+保留纹理时使用。

【通用高清增强系列 - 常规清晰化】
- general_2x: 通用超分辨率，2倍放大。适合风景/建筑/物品/动物。禁止用于含人脸图片。
- general_4x: 同上但4倍放大。

【高保真系列 - 保留原图风格】
- high_fidelity_2x: 高保真超分，2倍放大。消除伪影/噪声，保留原图纹理风格。适合摄影作品/艺术照/商品图/海报。
- high_fidelity_4x: 同上但4倍放大。

【降噪系列 - 不放大仅优化】
- sharpen_denoise_1x: 锐化降噪1x。去除常规噪声+锐化+生成新纹理。适合手机抓拍/运动模糊。禁止用于需要放大的场景。
- detail_denoise_1x: 极致保真降噪1x。去除复杂噪声+保留原图纹理。适合夜景/高感/会议抓拍。禁止用于需要锐化的场景。

【生成式人像系列 - SD模型，极低质量人像】
- generative_portrait_1x: Diffusion人像增强，不放大。补全人脸五官，生成自然纹理。适合极低质量大头照/侧脸/老照片/表情包。
- generative_portrait_2x: 同上+2倍放大。
- generative_portrait_4x: 同上+4倍放大。

【生成式通用系列 - SD模型，极低质量通用图】
- generative_1x: Diffusion通用增强，不放大。锐利化生成，适合极模糊/压缩图。
- generative_2x: 同上+2倍放大。
- generative_4x: 同上+4倍放大。

【生成式快速系列 - SD模型，效率优先】
- generative_enhance_fast_2x: SD快速生成2倍放大。效率优先，质量略低于generative_2x。仅当用户提"快速/效率优先"时使用。
- generative_enhance_fast_4x: SD快速生成4倍放大。效率优先，质量略低于generative_4x。仅当用户提"快速/效率优先"时使用。

━━━ 输入输出限制 ━━━
- 最大输入: 67MP | 最大输出: 432MP(非SD) / 34MP(SD)
- 输入格式: .jpg .jpeg .png .webp .bmp .tif .tiff .heic .heif
- 输出格式: .jpg .jpeg .png .webp .bmp .tif .tiff`,
		InputSchema: protocol.InputSchema{
			Type: "object",
			Properties: map[string]protocol.PropertySchema{
				"model_name": {
					Type: "string",
					Description: `增强模型名称（必填）。请根据上方决策树严格选择。
人像: face_2x, face_4x, face_v2_2x, face_v2_4x, generative_portrait_1x/2x/4x
通用: general_2x, general_4x, high_fidelity_2x/4x, generative_1x/2x/4x, generative_enhance_fast_2x/4x
降噪: sharpen_denoise_1x, detail_denoise_1x`,
					Enum: []string{
						"face_2x", "face_4x",
						"face_v2_2x", "face_v2_4x",
						"general_2x", "general_4x",
						"high_fidelity_2x", "high_fidelity_4x",
						"sharpen_denoise_1x", "detail_denoise_1x",
						"generative_portrait_1x", "generative_portrait_2x", "generative_portrait_4x",
						"generative_1x", "generative_2x", "generative_4x",
						"generative_enhance_fast_2x", "generative_enhance_fast_4x",
					},
				},
				"img_url": {
					Type:        "string",
					Description: "待增强图片的URL地址（必填）。支持: .jpg .jpeg .png .webp .bmp .tif .tiff .heic .heif",
				},
				"extension": {
					Type:        "string",
					Description: "输出图片格式，默认 .jpg",
					Enum:        []string{".jpg", ".jpeg", ".png", ".webp", ".bmp", ".tif", ".tiff"},
					Default:     ".jpg",
				},
				"upscale": {
					Type:        "integer",
					Description: "放大倍数（1-4）。每个模型有固定倍数（名称后缀1x/2x/4x），通常不需要手动设置",
					Minimum:     &minUpscale,
					Maximum:     &maxUpscale,
				},
			},
			Required: []string{"model_name", "img_url"},
		},
	}
}

func videoEnhanceTool() protocol.Tool {
	return protocol.Tool{
		Name: "video_enhance",
		Description: `视频增强/超分辨率工具 - 提升视频清晰度和分辨率。返回 job_id 用于查询任务状态。

━━━ 模型选择指引 ━━━

人脸/人像视频：
- face_soft_2x: 人脸柔和增强2x
- portrait_restore_1x: 人像修复，不放大
- portrait_restore_2x: 人像修复，2倍放大

通用视频：
- general_restore_1x: 通用修复，不放大
- general_restore_2x: 通用修复，2倍放大
- general_restore_4x: 通用修复，4倍放大

高级模型（SD）：
- ultrahd_restore_2x: 超高清修复2x
- generative_1x: 生成式增强

━━━ 限制 ━━━
- 最大输入: 36MP | 时长: 0.5秒~60分钟
- 输入: .mp4 .mov .avi .mkv .webm .flv .ts 等
- 输出: .mp4 .mov .mkv .m4v .avi .gif
- 积分 = 分辨率 × 帧率 × 时长`,
		InputSchema: protocol.InputSchema{
			Type: "object",
			Properties: map[string]protocol.PropertySchema{
				"video_url": {
					Type:        "string",
					Description: "待增强视频的URL地址（必填）",
				},
				"model_name": {
					Type:        "string",
					Description: "视频增强模型名称（必填）",
					Enum: []string{
						"face_soft_2x", "portrait_restore_2x", "portrait_restore_1x",
						"general_restore_1x", "general_restore_2x", "general_restore_4x",
						"ultrahd_restore_2x", "generative_1x",
					},
				},
				"resolution": {
					Type:        "string",
					Description: "输出分辨率（必填），格式 '宽x高'，例如 '1920x1080'",
				},
				"extension": {
					Type:        "string",
					Description: "输出视频格式，默认 .mp4",
					Enum:        []string{".mp4", ".mov", ".mkv", ".m4v", ".avi", ".gif"},
					Default:     ".mp4",
				},
			},
			Required: []string{"video_url", "model_name", "resolution"},
		},
	}
}

func taskStatusTool() protocol.Tool {
	return protocol.Tool{
		Name: "task_status",
		Description: `查询图片/视频增强任务的处理状态。
状态: WAITING/PENDING(等待) → INPUT_PREPARING(下载) → CONVERTING(处理中) → COMPLETED(完成,含结果URL)
异常: ERROR/ERROR_INTERRUPTION/TIMEOUT(失败) | NSFW(不合规) | CANCEL(取消)`,
		InputSchema: protocol.InputSchema{
			Type: "object",
			Properties: map[string]protocol.PropertySchema{
				"job_id": {
					Type:        "string",
					Description: "任务ID（必填），由 photo_enhance 或 video_enhance 返回",
				},
			},
			Required: []string{"job_id"},
		},
	}
}

func ossTransferTool() protocol.Tool {
	return protocol.Tool{
		Name: "oss_transfer",
		Description: "将远程URL文件转存到OSS存储，返回稳定的访问URL。适用于将外部图片/视频URL转存为持久链接。",
		InputSchema: protocol.InputSchema{
			Type: "object",
			Properties: map[string]protocol.PropertySchema{
				"url": {
					Type:        "string",
					Description: "要转存的远程文件URL（必填）",
				},
				"filename": {
					Type:        "string",
					Description: "可选的文件名",
				},
			},
			Required: []string{"url"},
		},
	}
}

func ossBatchTransferTool() protocol.Tool {
	return protocol.Tool{
		Name: "oss_batch_transfer",
		Description: "批量将远程URL文件转存到OSS存储，最多20个URL。",
		InputSchema: protocol.InputSchema{
			Type: "object",
			Properties: map[string]protocol.PropertySchema{
				"urls": {
					Type:        "string",
					Description: "要转存的远程文件URL列表（必填），用逗号分隔",
				},
			},
			Required: []string{"urls"},
		},
	}
}

func listPhotoModelsTool() protocol.Tool {
	return protocol.Tool{
		Name: "list_photo_models",
		Description: "列出所有18个图片增强模型的详细信息，包括适用场景、触发条件、放大倍数。帮助选择最合适的模型。",
		InputSchema: protocol.InputSchema{
			Type:       "object",
			Properties: map[string]protocol.PropertySchema{},
		},
	}
}

func listVideoModelsTool() protocol.Tool {
	return protocol.Tool{
		Name: "list_video_models",
		Description: "列出所有8个视频增强模型的详细信息，包括适用场景、支持分辨率。",
		InputSchema: protocol.InputSchema{
			Type:       "object",
			Properties: map[string]protocol.PropertySchema{},
		},
	}
}
