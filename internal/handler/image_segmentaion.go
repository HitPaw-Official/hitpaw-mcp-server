package handler

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hitpaw/mcp-server-hitpaw/internal/client"
	"github.com/hitpaw/mcp-server-hitpaw/internal/protocol"
)

// imageSegmentationArgs 图像分割参数
type imageSegmentationArgs struct {
	FileURL    string `json:"file_url"`
	ReturnType string `json:"return_type"`
}

// HandleImageSegmentation 处理图像分割（抠图）请求
func (h *Handlers) HandleImageSegmentation(arguments json.RawMessage) *protocol.CallToolResult {
	var args imageSegmentationArgs
	if err := json.Unmarshal(arguments, &args); err != nil {
		return protocol.ErrorResult(fmt.Sprintf("参数解析失败: %v", err))
	}

	// 参数校验
	if args.FileURL == "" {
		return protocol.ErrorResult("file_url 为必填参数，请提供图片的URL地址")
	}

	// 校验 return_type
	if args.ReturnType != "" {
		switch args.ReturnType {
		case "matted", "mask", "all":
			// ok
		default:
			return protocol.ErrorResult(fmt.Sprintf("return_type 取值错误: %s，仅支持 matted / mask / all", args.ReturnType))
		}
	}

	log.Printf("[ImageSegmentation] file_url=%s, return_type=%s", args.FileURL, args.ReturnType)

	// 调用 API
	resp, err := h.api.ImageSegmentation(&client.ImageSegmentationRequest{
		FileURL:    args.FileURL,
		ReturnType: args.ReturnType,
	})
	if err != nil {
		log.Printf("[ImageSegmentation] API error: %v", err)
		return protocol.ErrorResult(fmt.Sprintf("图像分割任务创建失败: %v", err))
	}

	returnTypeDesc := args.ReturnType
	if returnTypeDesc == "" {
		returnTypeDesc = "matted (默认透明背景图)"
	}

	result := fmt.Sprintf(`图像分割（抠图）任务已创建成功！

📋 任务详情：
- 任务ID: %s
- 返回类型: %s
- 消耗积分: %d

⏳ 任务正在后台处理中，请使用 task_status 工具查询处理结果。
完成后将返回：
- res_url: 抠图结果（透明背景图）
- mask_url: 主体掩码图（return_type=mask/all 时有值）

示例：task_status(job_id="%s")`, resp.JobID, returnTypeDesc, resp.ConsumeCoins, resp.JobID)

	return protocol.SuccessResult(result)
}
