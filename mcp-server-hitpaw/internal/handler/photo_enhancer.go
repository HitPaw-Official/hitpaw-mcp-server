package handler

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hitpaw/mcp-server-hitpaw/internal/client"
	"github.com/hitpaw/mcp-server-hitpaw/internal/protocol"
)

// photoEnhanceArgs 图片增强参数
type photoEnhanceArgs struct {
	ModelName string `json:"model_name"`
	ImgURL    string `json:"img_url"`
	Extension string `json:"extension"`
	Upscale   int    `json:"upscale"`
}

// HandlePhotoEnhance 处理图片增强请求
func (h *Handlers) HandlePhotoEnhance(arguments json.RawMessage) *protocol.CallToolResult {
	var args photoEnhanceArgs
	if err := json.Unmarshal(arguments, &args); err != nil {
		return protocol.ErrorResult(fmt.Sprintf("参数解析失败: %v", err))
	}

	// 参数校验
	if args.ModelName == "" {
		return protocol.ErrorResult("model_name 为必填参数，请使用 list_photo_models 查看可用模型")
	}
	if args.ImgURL == "" {
		return protocol.ErrorResult("img_url 为必填参数，请提供图片的URL地址")
	}

	// 设置默认值
	if args.Extension == "" {
		args.Extension = ".jpg"
	}

	log.Printf("[PhotoEnhance] model=%s, img_url=%s, extension=%s", args.ModelName, args.ImgURL, args.Extension)

	// 调用 API
	resp, err := h.api.PhotoEnhance(&client.PhotoEnhancerRequest{
		ModelName: args.ModelName,
		ImgURL:    args.ImgURL,
		Extension: args.Extension,
		Upscale:   args.Upscale,
	})
	if err != nil {
		log.Printf("[PhotoEnhance] API error: %v", err)
		return protocol.ErrorResult(fmt.Sprintf("图片增强任务创建失败: %v", err))
	}

	result := fmt.Sprintf(`图片增强任务已创建成功！

📋 任务详情：
- 任务ID: %s
- 使用模型: %s
- 消耗积分: %d
- 输出格式: %s

⏳ 任务正在后台处理中，请使用 task_status 工具查询处理结果。
示例：task_status(job_id="%s")`, resp.JobID, args.ModelName, resp.ConsumeCoins, args.Extension, resp.JobID)

	return protocol.SuccessResult(result)
}
