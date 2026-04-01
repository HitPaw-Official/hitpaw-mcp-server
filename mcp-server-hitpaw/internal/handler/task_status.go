package handler

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hitpaw/mcp-server-hitpaw/internal/protocol"
)

// taskStatusArgs 任务状态查询参数
type taskStatusArgs struct {
	JobID string `json:"job_id"`
}

// HandleTaskStatus 处理任务状态查询
func (h *Handlers) HandleTaskStatus(arguments json.RawMessage) *protocol.CallToolResult {
	var args taskStatusArgs
	if err := json.Unmarshal(arguments, &args); err != nil {
		return protocol.ErrorResult(fmt.Sprintf("参数解析失败: %v", err))
	}

	if args.JobID == "" {
		return protocol.ErrorResult("job_id 为必填参数")
	}

	log.Printf("[TaskStatus] job_id=%s", args.JobID)

	resp, err := h.api.TaskStatus(args.JobID)
	if err != nil {
		log.Printf("[TaskStatus] API error: %v", err)
		return protocol.ErrorResult(fmt.Sprintf("查询任务状态失败: %v", err))
	}

	// 根据状态返回不同信息
	statusIcon := getStatusIcon(resp.Status)
	var result string

	switch resp.Status {
	case "COMPLETED":
		result = fmt.Sprintf(`%s 任务处理完成！

📋 任务详情：
- 任务ID: %s
- 状态: %s (已完成)
- 结果文件: %s
- 原始文件: %s

✅ 您可以通过结果文件URL下载处理后的文件。`, statusIcon, resp.JobID, resp.Status, resp.ResURL, resp.OriginalURL)

	case "ERROR", "ERROR_INTERRUPTION", "TIMEOUT":
		result = fmt.Sprintf(`%s 任务处理失败

📋 任务详情：
- 任务ID: %s
- 状态: %s

❌ 任务处理失败，可能原因：输入文件格式不支持、分辨率超限、服务繁忙等。
建议检查输入参数后重新提交。`, statusIcon, resp.JobID, resp.Status)

	default:
		result = fmt.Sprintf(`%s 任务处理中...

📋 任务详情：
- 任务ID: %s
- 状态: %s (%s)
- 原始文件: %s

⏳ 任务仍在处理中，请稍后再次查询。`, statusIcon, resp.JobID, resp.Status, getStatusDesc(resp.Status), resp.OriginalURL)
	}

	return protocol.SuccessResult(result)
}

// getStatusIcon 获取状态图标
func getStatusIcon(status string) string {
	switch status {
	case "COMPLETED":
		return "✅"
	case "ERROR", "ERROR_INTERRUPTION", "TIMEOUT", "REJECT", "NSFW":
		return "❌"
	case "CONVERTING":
		return "🔄"
	case "PENDING", "WAITING":
		return "⏳"
	default:
		return "📋"
	}
}

// getStatusDesc 获取状态描述
func getStatusDesc(status string) string {
	switch status {
	case "WAITING":
		return "等待中"
	case "PENDING":
		return "已提交，等待启动"
	case "INPUT_PREPARING":
		return "正在下载输入文件"
	case "CONVERTING":
		return "正在处理中"
	case "CONVERT_COMPLETED":
		return "处理完成，正在上传"
	case "OUTPUT_SAVING":
		return "正在保存结果"
	case "OUTPUT_SAVED":
		return "结果已保存"
	case "COMPLETED":
		return "已完成"
	case "ERROR":
		return "处理失败"
	case "ERROR_INTERRUPTION":
		return "处理中断"
	case "TIMEOUT":
		return "处理超时"
	case "CANCEL":
		return "已取消"
	case "REJECT":
		return "已拒绝"
	case "NSFW":
		return "内容不合规"
	default:
		return status
	}
}
