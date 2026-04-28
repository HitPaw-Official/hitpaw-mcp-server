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
	server.RegisterTool(imageSegmentationTool(), h.HandleImageSegmentation)
	server.RegisterTool(taskStatusTool(), h.HandleTaskStatus)
	server.RegisterTool(ossTransferTool(), h.HandleOSSTransfer)
	server.RegisterTool(ossBatchTransferTool(), h.HandleOSSBatchTransfer)
	server.RegisterTool(ossUploadTool(), h.HandleOSSUpload)
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
	maxUpscale := float64(8)
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
│   ├── 4倍放大 → generative_4x
│   ├── 6倍放大（极致放大需求）→ generative_6x
│   └── 8倍放大（极致放大需求）→ generative_8x
├── "高保真/保留风格/摄影作品/艺术照"
│   ├── 2倍放大 → high_fidelity_2x
│   └── 4倍放大 → high_fidelity_4x
└── 常规清晰化/画质提升（默认）
    ├── 2倍放大或未指定 → general_2x
    └── 4倍放大 → general_4x

⚠️ 重要：6x/8x 仅有 generative_6x / generative_8x 两个模型（生成式SD）。
   如用户需要 6倍/8倍 放大，只能使用这两个模型。
   其他模型系列（face/face_v2/general/high_fidelity）目前最高支持 4x。

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
- generative_6x: 同上+6倍放大。⚠️ 仅当用户明确需要6倍放大时使用，注意输入分辨率不能过大（输出 ≤ 34MP）。
- generative_8x: 同上+8倍放大。⚠️ 仅当用户明确需要8倍放大时使用，注意输入分辨率不能过大（输出 ≤ 34MP）。

【生成式快速系列 - SD模型，效率优先】
- generative_enhance_fast_2x: SD快速生成2倍放大。效率优先，质量略低于generative_2x。仅当用户提"快速/效率优先"时使用。
- generative_enhance_fast_4x: SD快速生成4倍放大。效率优先，质量略低于generative_4x。仅当用户提"快速/效率优先"时使用。

━━━ 输入输出限制 ━━━
- 最大输入: 67MP | 最大输出: 432MP(非SD) / 34MP(SD)
- 输入格式: .jpg .jpeg .png .webp .bmp .tif .tiff .heic .heif
- 输出格式: .jpg .jpeg .png .webp .bmp .tif .tiff
- 6x/8x 为 SD 模型，输出上限 34MP，请控制输入图片尺寸`,
		InputSchema: protocol.InputSchema{
			Type: "object",
			Properties: map[string]protocol.PropertySchema{
				"model_name": {
					Type: "string",
					Description: `增强模型名称（必填）。请根据上方决策树严格选择。
人像: face_2x, face_4x, face_v2_2x, face_v2_4x, generative_portrait_1x/2x/4x
通用: general_2x, general_4x, high_fidelity_2x/4x, generative_1x/2x/4x/6x/8x, generative_enhance_fast_2x/4x
降噪: sharpen_denoise_1x, detail_denoise_1x
注意: 6x/8x 仅 generative_6x、generative_8x 可用`,
					Enum: []string{
						"face_2x", "face_4x",
						"face_v2_2x", "face_v2_4x",
						"general_2x", "general_4x",
						"high_fidelity_2x", "high_fidelity_4x",
						"sharpen_denoise_1x", "detail_denoise_1x",
						"generative_portrait_1x", "generative_portrait_2x", "generative_portrait_4x",
						"generative_1x", "generative_2x", "generative_4x", "generative_6x", "generative_8x",
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
					Description: "放大倍数（1-8）。每个模型有固定倍数（名称后缀1x/2x/4x/6x/8x），通常不需要手动设置",
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

func imageSegmentationTool() protocol.Tool {
	return protocol.Tool{
		Name: "image_segmentation",
		Description: `AI图像分割工具（抠图）- 自动识别图像主体并与背景分离，生成透明背景图。返回 job_id 用于查询任务状态。

━━━ 使用场景 ━━━
✅ 用户提出 "抠图/去背景/分离主体/透明背景/换背景" 等需求
✅ 用户需要把图片里的主体（人像/物品/动物）从背景中分离出来
✅ 用户提到 "保留主体/背景透明/PNG透明图/白底/透明底" 等关键词
✅ 关键词: 抠图, 去背景, 背景透明, 分离主体, 透明底, PNG, 主体提取

❌ 不要用于以下场景：
- 需要放大/清晰化/修复图片 → 使用 photo_enhance
- 需要生成/编辑图片内容 → 使用对应的图片编辑工具
- 非图片格式的文件处理 → 使用其他文件处理工具

━━━ 能力说明 ━━━
- 自动识别图像主体，分离主体与背景，生成透明背景图
- 支持返回多种结果格式：透明背景图、mask掩码图、全部结果
- 支持公网图片URL输入（如需上传本地文件，请先用 oss_upload 获取URL）

━━━ return_type 参数说明 ━━━
- matted (默认): 仅返回透明背景的抠图结果
- mask: 仅返回主体的 mask 掩码图（黑白二值图）
- all: 同时返回透明背景图、mask图和缩略图

━━━ 处理流程 ━━━
1. 调用本工具创建任务 → 获得 job_id
2. 异步处理，需通过 task_status 工具轮询查询结果
3. 完成后返回 res_url（透明背景图）、mask_url（掩码图，return_type=mask/all 时有值）

━━━ 限制 ━━━
- 异步处理，不支持实时同步返回
- 单次任务扣除积分（具体扣费以接口返回为准）
- 输出URL有效期 24 小时`,
		InputSchema: protocol.InputSchema{
			Type: "object",
			Properties: map[string]protocol.PropertySchema{
				"file_url": {
					Type:        "string",
					Description: "待抠图的图片URL地址（必填）。如果是本地文件，请先使用 oss_upload 工具上传获取URL",
				},
				"return_type": {
					Type:        "string",
					Description: "返回结果类型。matted=仅透明背景图（默认），mask=仅mask掩码图，all=透明背景图+mask图+缩略图",
					Enum:        []string{"matted", "mask", "all"},
					Default:     "matted",
				},
			},
			Required: []string{"file_url"},
		},
	}
}

func taskStatusTool() protocol.Tool {
	return protocol.Tool{
		Name: "task_status",
		Description: `查询图片/视频/抠图任务的处理状态。
状态: WAITING/PENDING(等待) → INPUT_PREPARING(下载) → CONVERTING(处理中) → COMPLETED(完成,含结果URL)
异常: ERROR/ERROR_INTERRUPTION/TIMEOUT(失败) | NSFW(不合规) | CANCEL(取消)

返回字段说明：
- res_url: 结果文件URL（增强图片/视频URL，或抠图后的透明背景图）
- mask_url: 抠图任务专属，主体mask掩码图URL（仅图像分割任务返回，其他任务为空）
- original_url: 原始输入文件URL`,
		InputSchema: protocol.InputSchema{
			Type: "object",
			Properties: map[string]protocol.PropertySchema{
				"job_id": {
					Type:        "string",
					Description: "任务ID（必填），由 photo_enhance / video_enhance / image_segmentation 返回",
				},
			},
			Required: []string{"job_id"},
		},
	}
}

func ossTransferTool() protocol.Tool {
	return protocol.Tool{
		Name:        "oss_transfer",
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
		Name:        "oss_batch_transfer",
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

func ossUploadTool() protocol.Tool {
	return protocol.Tool{
		Name: "oss_upload",
		Description: `上传本地文件到OSS存储，返回可访问的URL。适用于将本地图片/视频文件上传到OSS，获取URL后可用于 photo_enhance / video_enhance / image_segmentation。

使用场景：
- 用户提供了本地文件（base64编码），需要先上传到OSS获取URL
- 上传后返回的URL可直接作为 photo_enhance 的 img_url、video_enhance 的 video_url、image_segmentation 的 file_url 使用`,
		InputSchema: protocol.InputSchema{
			Type: "object",
			Properties: map[string]protocol.PropertySchema{
				"file_data": {
					Type:        "string",
					Description: "文件内容的Base64编码字符串（必填）",
				},
				"filename": {
					Type:        "string",
					Description: "文件名（必填），包含扩展名，例如 'photo.jpg'、'video.mp4'",
				},
			},
			Required: []string{"file_data", "filename"},
		},
	}
}

func listPhotoModelsTool() protocol.Tool {
	return protocol.Tool{
		Name:        "list_photo_models",
		Description: "列出所有20个图片增强模型的详细信息，包括适用场景、触发条件、放大倍数。帮助选择最合适的模型。",
		InputSchema: protocol.InputSchema{
			Type:       "object",
			Properties: map[string]protocol.PropertySchema{},
		},
	}
}

func listVideoModelsTool() protocol.Tool {
	return protocol.Tool{
		Name:        "list_video_models",
		Description: "列出所有8个视频增强模型的详细信息，包括适用场景、支持分辨率。",
		InputSchema: protocol.InputSchema{
			Type:       "object",
			Properties: map[string]protocol.PropertySchema{},
		},
	}
}
