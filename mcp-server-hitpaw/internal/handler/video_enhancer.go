package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hitpaw/mcp-server-hitpaw/internal/client"
	"github.com/hitpaw/mcp-server-hitpaw/internal/protocol"
)

// videoEnhanceArgs 视频增强参数
type videoEnhanceArgs struct {
	VideoURL   string `json:"video_url"`
	ModelName  string `json:"model_name"`
	Resolution string `json:"resolution"` // "1920x1080" 格式
	Extension  string `json:"extension"`
}

// HandleVideoEnhance 处理视频增强请求
func (h *Handlers) HandleVideoEnhance(arguments json.RawMessage) *protocol.CallToolResult {
	var args videoEnhanceArgs
	if err := json.Unmarshal(arguments, &args); err != nil {
		return protocol.ErrorResult(fmt.Sprintf("参数解析失败: %v", err))
	}

	// 参数校验
	if args.VideoURL == "" {
		return protocol.ErrorResult("video_url 为必填参数，请提供视频的URL地址")
	}
	if args.ModelName == "" {
		return protocol.ErrorResult("model_name 为必填参数，请使用 list_video_models 查看可用模型")
	}
	if args.Resolution == "" {
		return protocol.ErrorResult("resolution 为必填参数，格式为 '宽x高'，例如 '1920x1080'")
	}

	// 解析分辨率
	resolution, err := parseResolution(args.Resolution)
	if err != nil {
		return protocol.ErrorResult(fmt.Sprintf("分辨率格式错误: %v。正确格式: '1920x1080'", err))
	}

	// 设置默认值
	if args.Extension == "" {
		args.Extension = ".mp4"
	}

	log.Printf("[VideoEnhance] model=%s, video_url=%s, resolution=%v", args.ModelName, args.VideoURL, resolution)

	// 调用 API
	resp, err := h.api.VideoEnhance(&client.VideoEnhancerRequest{
		VideoURL:   args.VideoURL,
		ModelName:  args.ModelName,
		Resolution: resolution,
		Extension:  args.Extension,
	})
	if err != nil {
		log.Printf("[VideoEnhance] API error: %v", err)
		return protocol.ErrorResult(fmt.Sprintf("视频增强任务创建失败: %v", err))
	}

	result := fmt.Sprintf(`视频增强任务已创建成功！

📋 任务详情：
- 任务ID: %s
- 使用模型: %s
- 目标分辨率: %dx%d
- 消耗积分: %d
- 输出格式: %s

⏳ 视频处理时间较长，请使用 task_status 工具查询处理结果。
示例：task_status(job_id="%s")`, resp.JobID, args.ModelName, resolution[0], resolution[1], resp.ConsumeCoins, args.Extension, resp.JobID)

	return protocol.SuccessResult(result)
}

// parseResolution 解析分辨率字符串 "1920x1080" → [1920, 1080]
func parseResolution(s string) ([]int, error) {
	s = strings.TrimSpace(s)
	// 支持 "x", "X", "*", "×" 分隔
	var parts []string
	for _, sep := range []string{"×", "x", "X", "*"} {
		if strings.Contains(s, sep) {
			parts = strings.SplitN(s, sep, 2)
			break
		}
	}

	if len(parts) != 2 {
		return nil, fmt.Errorf("无法解析 '%s'，请使用 '宽x高' 格式", s)
	}

	width, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil, fmt.Errorf("宽度无效: %s", parts[0])
	}

	height, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return nil, fmt.Errorf("高度无效: %s", parts[1])
	}

	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("分辨率必须为正整数")
	}

	return []int{width, height}, nil
}
