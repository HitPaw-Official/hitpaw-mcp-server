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
	serverVersion = "1.0.0"
	instructions  = `HitPaw AI 图片/视频增强 MCP 服务。

可用工具：
1. photo_enhance - 图片增强（18个模型：超分辨率、降噪、人像优化、生成式修复）
2. video_enhance - 视频增强（8个模型：人像修复、通用修复、超高清修复）
3. task_status - 查询任务处理状态和结果
4. oss_transfer - 文件URL转存到OSS
5. oss_batch_transfer - 批量文件URL转存
6. oss_upload - 上传本地文件到OSS（Base64编码），获取URL后可用于增强
7. list_photo_models - 查看图片增强全部18个模型详情
8. list_video_models - 查看视频增强全部8个模型详情

标准使用流程：
1. 用户描述需求 → 根据 photo_enhance 工具描述中的决策树自动选择最佳模型
2. 如果用户提供本地文件 → 先用 oss_upload 上传获取URL
3. 调用 photo_enhance 或 video_enhance 创建任务 → 获得 job_id
4. 调用 task_status 查询结果 → 获得结果URL

图片模型快速选择（详情见 photo_enhance 工具描述）：
- 含人脸 + 美颜 → face_2x/face_4x
- 含人脸 + 保留纹理 → face_v2_2x/face_v2_4x
- 含人脸 + 极低质量 → generative_portrait_1x/2x/4x
- 风景/建筑 + 常规清晰 → general_2x/general_4x
- 摄影/艺术 + 保留风格 → high_fidelity_2x/high_fidelity_4x
- 仅降噪 + 锐化 → sharpen_denoise_1x
- 仅降噪 + 保真 → detail_denoise_1x
- 极低质量 + 生成补全 → generative_1x/2x/4x
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

	log.Println("[Main] 所有工具已注册（8个工具，18个图片模型，8个视频模型），等待客户端连接...")

	// 启动 MCP Server
	if err := server.Run(); err != nil {
		log.Fatalf("[Main] Server error: %v", err)
	}
}
