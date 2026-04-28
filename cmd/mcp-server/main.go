package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hitpaw/mcp-server-hitpaw/configs"
	"github.com/hitpaw/mcp-server-hitpaw/internal/client"
	"github.com/hitpaw/mcp-server-hitpaw/internal/handler"
	"github.com/hitpaw/mcp-server-hitpaw/internal/protocol"
)

const (
	serverName    = "hitpaw-mcp-server"
	serverVersion = "1.1.1"
	instructions  = `HitPaw AI 图片/视频增强 MCP 服务。

可用工具：
1. photo_enhance - 图片增强（20个模型：超分辨率2x/4x、降噪、人像优化、生成式修复2x/4x/6x/8x）
2. video_enhance - 视频增强（8个模型：人像修复、通用修复、超高清修复）
3. image_segmentation - 图像分割/抠图（自动分离主体与背景，生成透明背景图）
4. task_status - 查询任务处理状态和结果（支持 photo/video/segmentation 全部任务类型）
5. oss_transfer - 文件URL转存到OSS
6. oss_batch_transfer - 批量文件URL转存
7. oss_upload - 上传本地文件到OSS（Base64编码），获取URL后可用于增强或抠图
8. list_photo_models - 查看图片增强全部20个模型详情
9. list_video_models - 查看视频增强全部8个模型详情

标准使用流程：
1. 用户描述需求 → 根据工具描述自动选择最佳工具/模型
2. 如果用户提供本地文件 → 先用 oss_upload 上传获取URL
3. 调用 photo_enhance / video_enhance / image_segmentation 创建任务 → 获得 job_id
4. 调用 task_status 查询结果 → 获得结果URL（抠图任务额外返回 mask_url）

需求路由速查：
- 抠图/去背景/透明背景/分离主体 → image_segmentation
- 放大/高清/清晰/降噪 → photo_enhance
- 视频增强/视频高清 → video_enhance

图片模型快速选择（详情见 photo_enhance 工具描述）：
- 含人脸 + 美颜 → face_2x/face_4x
- 含人脸 + 保留纹理 → face_v2_2x/face_v2_4x
- 含人脸 + 极低质量 → generative_portrait_1x/2x/4x
- 风景/建筑 + 常规清晰 → general_2x/general_4x
- 摄影/艺术 + 保留风格 → high_fidelity_2x/high_fidelity_4x
- 仅降噪 + 锐化 → sharpen_denoise_1x
- 仅降噪 + 保真 → detail_denoise_1x
- 极低质量 + 生成补全 → generative_1x/2x/4x
- 6倍/8倍放大 → 仅 generative_6x / generative_8x 可用
- 极低质量 + 快速 → generative_enhance_fast_2x/4x`
)

func main() {
	// 加载配置
	cfg, err := configs.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "配置加载失败: %v\n", err)
		os.Exit(1)
	}

	log.SetOutput(os.Stderr)
	log.Printf("[Main] HitPaw MCP Server v%s 启动中...", serverVersion)
	log.Printf("[Main] API Base URL: %s", cfg.APIBaseURL)

	// 创建 API 客户端
	apiClient := client.NewAPIClient(cfg.APIBaseURL, cfg.APIKey)

	// 创建 MCP Server
	server := protocol.NewMCPServer(serverName, serverVersion, instructions)

	// 注册所有工具
	handler.RegisterAllTools(server, apiClient)

	log.Println("[Main] 所有工具已注册（9个工具，20个图片模型，8个视频模型，1个抠图工具），等待客户端连接...")

	// 启动 MCP Server
	if err := server.Run(); err != nil {
		log.Fatalf("[Main] Server error: %v", err)
	}
}
