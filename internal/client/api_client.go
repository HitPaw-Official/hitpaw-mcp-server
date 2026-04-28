package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

// APIClient 现有 HTTP API 的客户端
type APIClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewAPIClient 创建 API 客户端
func NewAPIClient(baseURL, apiKey string) *APIClient {
	return &APIClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 300 * time.Second,
		},
	}
}

// ==================== 请求/响应结构体 ====================

// APIResponse 通用 API 响应
type APIResponse struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

// PhotoEnhancerRequest 图片增强请求
type PhotoEnhancerRequest struct {
	ModelName string `json:"model_name"`
	ImgURL    string `json:"img_url"`
	Extension string `json:"extension"`
	Pid       int    `json:"pid,omitempty"`
	Upscale   int    `json:"upscale,omitempty"`
	Exif      bool   `json:"exif,omitempty"`
	DPI       int64  `json:"DPI,omitempty"`
}

// PhotoEnhancerResponse 图片增强响应
type PhotoEnhancerResponse struct {
	JobID        string `json:"job_id"`
	ConsumeCoins int64  `json:"consume_coins"`
}

// VideoEnhancerRequest 视频增强请求
type VideoEnhancerRequest struct {
	VideoURL           string `json:"video_url"`
	Extension          string `json:"extension"`
	ModelID            int    `json:"model_id,omitempty"`
	ModelName          string `json:"model_name"`
	Resolution         []int  `json:"resolution"`
	OriginalResolution []int  `json:"original_resolution,omitempty"`
}

// VideoEnhancerResponse 视频增强响应
type VideoEnhancerResponse struct {
	JobID        string `json:"job_id"`
	ConsumeCoins int64  `json:"consume_coins"`
}

// ImageSegmentationRequest 图像分割（抠图）请求
type ImageSegmentationRequest struct {
	FileURL    string `json:"file_url"`
	ReturnType string `json:"return_type,omitempty"`
}

// ImageSegmentationResponse 图像分割响应（与 PhotoEnhancerResponse 同结构）
type ImageSegmentationResponse struct {
	JobID        string `json:"job_id"`
	ConsumeCoins int64  `json:"consume_coins"`
}

// TaskStatusRequest 任务状态请求
type TaskStatusRequest struct {
	JobID string `json:"job_id"`
}

// TaskStatusResponse 任务状态响应
type TaskStatusResponse struct {
	JobID       string `json:"job_id"`
	Status      string `json:"status"`
	ResURL      string `json:"res_url"`
	OriginalURL string `json:"original_url"`
	MaskURL     string `json:"mask_url"` // 图像分割任务专属
}

// OSSTransferRequest OSS 转存请求
type OSSTransferRequest struct {
	URL      string `json:"url"`
	Filename string `json:"filename,omitempty"`
}

// OSSTransferResponse OSS 转存响应
type OSSTransferResponse struct {
	URL       string `json:"url"`
	ObjectKey string `json:"object_key"`
	Size      int64  `json:"size"`
}

// OSSBatchTransferRequest OSS 批量转存请求
type OSSBatchTransferRequest struct {
	URLs []string `json:"urls"`
}

// OSSBatchTransferResponse OSS 批量转存响应
type OSSBatchTransferResponse struct {
	Total   int                    `json:"total"`
	Success int                    `json:"success"`
	Failed  int                    `json:"failed"`
	Items   []OSSBatchTransferItem `json:"items"`
}

// OSSBatchTransferItem 批量转存单项结果
type OSSBatchTransferItem struct {
	SourceURL string `json:"source_url"`
	URL       string `json:"url,omitempty"`
	ObjectKey string `json:"object_key,omitempty"`
	Size      int64  `json:"size,omitempty"`
	Error     string `json:"error,omitempty"`
}

// OSSUploadResponse OSS 上传响应
type OSSUploadResponse struct {
	URL       string `json:"url"`
	ObjectKey string `json:"object_key"`
	Size      int64  `json:"size"`
}

// ==================== API 调用方法 ====================

// PhotoEnhance 调用图片增强接口
func (c *APIClient) PhotoEnhance(req *PhotoEnhancerRequest) (*PhotoEnhancerResponse, error) {
	var resp PhotoEnhancerResponse
	err := c.doJSON("POST", "/api/photo-enhancer", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// VideoEnhance 调用视频增强接口
func (c *APIClient) VideoEnhance(req *VideoEnhancerRequest) (*VideoEnhancerResponse, error) {
	var resp VideoEnhancerResponse
	err := c.doJSON("POST", "/api/video-enhancer", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// ImageSegmentation 调用图像分割（抠图）接口
func (c *APIClient) ImageSegmentation(req *ImageSegmentationRequest) (*ImageSegmentationResponse, error) {
	var resp ImageSegmentationResponse
	err := c.doJSON("POST", "/api/image-segmentation", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// TaskStatus 查询任务状态
func (c *APIClient) TaskStatus(jobID string) (*TaskStatusResponse, error) {
	var resp TaskStatusResponse
	err := c.doJSON("POST", "/api/task-status", &TaskStatusRequest{JobID: jobID}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// OSSTransfer 远程URL转存到OSS
func (c *APIClient) OSSTransfer(req *OSSTransferRequest) (*OSSTransferResponse, error) {
	var resp OSSTransferResponse
	err := c.doJSON("POST", "/api/oss/transfer", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// OSSBatchTransfer 批量远程URL转存
func (c *APIClient) OSSBatchTransfer(req *OSSBatchTransferRequest) (*OSSBatchTransferResponse, error) {
	var resp OSSBatchTransferResponse
	err := c.doJSON("POST", "/api/oss/batch-transfer", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// OSSUploadFile 上传本地文件（multipart）
func (c *APIClient) OSSUploadFile(fileData []byte, filename string) (*OSSUploadResponse, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("create form file failed: %w", err)
	}
	if _, err := part.Write(fileData); err != nil {
		return nil, fmt.Errorf("write file data failed: %w", err)
	}
	writer.Close()

	reqURL := c.baseURL + "/api/oss/upload"
	httpReq, err := http.NewRequest("POST", reqURL, body)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())
	httpReq.Header.Set("APIKEY", c.apiKey)

	log.Printf("[APIClient] POST /api/oss/upload filename=%s size=%d", filename, len(fileData))

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, _ := io.ReadAll(httpResp.Body)

	log.Printf("[APIClient] Response status=%d body=%s", httpResp.StatusCode, string(respBody))

	if httpResp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", httpResp.StatusCode, string(respBody))
	}

	var apiResp APIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("parse response failed: %w, body: %s", err, string(respBody))
	}

	if apiResp.Code != 200 {
		return nil, fmt.Errorf("API error: code=%d, msg=%s", apiResp.Code, apiResp.Msg)
	}

	var resp OSSUploadResponse
	if err := json.Unmarshal(apiResp.Data, &resp); err != nil {
		return nil, fmt.Errorf("parse data failed: %w", err)
	}
	return &resp, nil
}

// ==================== 内部方法 ====================

// doJSON 发送 JSON 请求并解析响应
func (c *APIClient) doJSON(method, path string, reqBody interface{}, result interface{}) error {
	reqURL := c.baseURL + path

	var bodyReader io.Reader
	if reqBody != nil {
		data, err := json.Marshal(reqBody)
		if err != nil {
			return fmt.Errorf("marshal request failed: %w", err)
		}
		bodyReader = bytes.NewReader(data)
		log.Printf("[APIClient] %s %s body=%s", method, path, string(data))
	}

	httpReq, err := http.NewRequest(method, reqURL, bodyReader)
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("APIKEY", c.apiKey)

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return fmt.Errorf("read response failed: %w", err)
	}

	log.Printf("[APIClient] Response status=%d body=%s", httpResp.StatusCode, string(respBody))

	if httpResp.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d: %s", httpResp.StatusCode, string(respBody))
	}

	// 解析通用响应
	var apiResp APIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return fmt.Errorf("parse response failed: %w, body: %s", err, string(respBody))
	}

	if apiResp.Code != 200 {
		return fmt.Errorf("API error: code=%d, msg=%s", apiResp.Code, apiResp.Msg)
	}

	// 解析 data 字段
	if result != nil && len(apiResp.Data) > 0 {
		if err := json.Unmarshal(apiResp.Data, result); err != nil {
			return fmt.Errorf("parse data failed: %w", err)
		}
	}

	return nil
}
