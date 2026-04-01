package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hitpaw/mcp-server-hitpaw/internal/client"
	"github.com/hitpaw/mcp-server-hitpaw/internal/protocol"
)

// ossTransferArgs OSS 转存参数
type ossTransferArgs struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
}

// HandleOSSTransfer 处理 OSS 转存
func (h *Handlers) HandleOSSTransfer(arguments json.RawMessage) *protocol.CallToolResult {
	var args ossTransferArgs
	if err := json.Unmarshal(arguments, &args); err != nil {
		return protocol.ErrorResult(fmt.Sprintf("参数解析失败: %v", err))
	}

	if args.URL == "" {
		return protocol.ErrorResult("url 为必填参数")
	}

	log.Printf("[OSSTransfer] url=%s", args.URL)

	resp, err := h.api.OSSTransfer(&client.OSSTransferRequest{
		URL:      args.URL,
		Filename: args.Filename,
	})
	if err != nil {
		return protocol.ErrorResult(fmt.Sprintf("文件转存失败: %v", err))
	}

	result := fmt.Sprintf(`文件转存成功！

📋 转存详情：
- 访问URL: %s
- 对象Key: %s
- 文件大小: %s`, resp.URL, resp.ObjectKey, formatFileSize(resp.Size))

	return protocol.SuccessResult(result)
}

// ossBatchTransferArgs OSS 批量转存参数
type ossBatchTransferArgs struct {
	URLs string `json:"urls"` // 逗号分隔的URL列表
}

// HandleOSSBatchTransfer 处理 OSS 批量转存
func (h *Handlers) HandleOSSBatchTransfer(arguments json.RawMessage) *protocol.CallToolResult {
	var args ossBatchTransferArgs
	if err := json.Unmarshal(arguments, &args); err != nil {
		return protocol.ErrorResult(fmt.Sprintf("参数解析失败: %v", err))
	}

	if args.URLs == "" {
		return protocol.ErrorResult("urls 为必填参数，多个URL用逗号分隔")
	}

	// 解析URL列表
	urls := strings.Split(args.URLs, ",")
	cleanURLs := make([]string, 0, len(urls))
	for _, u := range urls {
		u = strings.TrimSpace(u)
		if u != "" {
			cleanURLs = append(cleanURLs, u)
		}
	}

	if len(cleanURLs) == 0 {
		return protocol.ErrorResult("至少需要提供一个URL")
	}
	if len(cleanURLs) > 20 {
		return protocol.ErrorResult("批量转存最多支持20个URL")
	}

	log.Printf("[OSSBatchTransfer] urls count=%d", len(cleanURLs))

	resp, err := h.api.OSSBatchTransfer(&client.OSSBatchTransferRequest{
		URLs: cleanURLs,
	})
	if err != nil {
		return protocol.ErrorResult(fmt.Sprintf("批量转存失败: %v", err))
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("批量转存完成！\n\n📋 总计: %d | 成功: %d | 失败: %d\n\n", resp.Total, resp.Success, resp.Failed))

	for i, item := range resp.Items {
		if item.Error != "" {
			sb.WriteString(fmt.Sprintf("[%d] ❌ %s → 失败: %s\n", i+1, item.SourceURL, item.Error))
		} else {
			sb.WriteString(fmt.Sprintf("[%d] ✅ %s → %s (%s)\n", i+1, item.SourceURL, item.URL, formatFileSize(item.Size)))
		}
	}

	return protocol.SuccessResult(sb.String())
}

// ossUploadArgs OSS 上传参数
type ossUploadArgs struct {
	FileData string `json:"file_data"` // Base64 编码的文件内容
	Filename string `json:"filename"`  // 文件名
}

// HandleOSSUpload 处理本地文件上传到 OSS
func (h *Handlers) HandleOSSUpload(arguments json.RawMessage) *protocol.CallToolResult {
	var args ossUploadArgs
	if err := json.Unmarshal(arguments, &args); err != nil {
		return protocol.ErrorResult(fmt.Sprintf("参数解析失败: %v", err))
	}

	if args.FileData == "" {
		return protocol.ErrorResult("file_data 为必填参数，请提供Base64编码的文件内容")
	}
	if args.Filename == "" {
		return protocol.ErrorResult("filename 为必填参数，请提供文件名（含扩展名）")
	}

	// 解码 Base64 数据
	fileData, err := base64.StdEncoding.DecodeString(args.FileData)
	if err != nil {
		return protocol.ErrorResult(fmt.Sprintf("Base64解码失败: %v。请确保 file_data 是有效的Base64编码字符串", err))
	}

	log.Printf("[OSSUpload] filename=%s, size=%d bytes", args.Filename, len(fileData))

	// 调用 API 上传
	resp, err := h.api.OSSUploadFile(fileData, args.Filename)
	if err != nil {
		log.Printf("[OSSUpload] API error: %v", err)
		return protocol.ErrorResult(fmt.Sprintf("文件上传失败: %v", err))
	}

	result := fmt.Sprintf(`文件上传成功！

📋 上传详情：
- 访问URL: %s
- 对象Key: %s
- 文件大小: %s

💡 该URL可直接用于 photo_enhance 的 img_url 或 video_enhance 的 video_url 参数。`, resp.URL, resp.ObjectKey, formatFileSize(resp.Size))

	return protocol.SuccessResult(result)
}

// formatFileSize 格式化文件大小
func formatFileSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	}
	if size < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(size)/1024)
	}
	if size < 1024*1024*1024 {
		return fmt.Sprintf("%.1f MB", float64(size)/1024/1024)
	}
	return fmt.Sprintf("%.2f GB", float64(size)/1024/1024/1024)
}
