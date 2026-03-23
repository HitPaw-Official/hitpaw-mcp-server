package protocol

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// ToolHandler 工具处理函数类型
type ToolHandler func(arguments json.RawMessage) *CallToolResult

// MCPServer MCP 服务器
type MCPServer struct {
	serverInfo   ServerInfo
	capabilities ServerCapabilities
	instructions string

	tools    []Tool
	handlers map[string]ToolHandler

	reader *bufio.Reader
	writer io.Writer
	mu     sync.Mutex
}

// NewMCPServer 创建 MCP 服务器
func NewMCPServer(name, version, instructions string) *MCPServer {
	return &MCPServer{
		serverInfo: ServerInfo{
			Name:    name,
			Version: version,
		},
		capabilities: ServerCapabilities{
			Tools: &ToolCapability{ListChanged: false},
		},
		instructions: instructions,
		tools:        make([]Tool, 0),
		handlers:     make(map[string]ToolHandler),
		reader:       bufio.NewReader(os.Stdin),
		writer:       os.Stdout,
	}
}

// RegisterTool 注册工具
func (s *MCPServer) RegisterTool(tool Tool, handler ToolHandler) {
	s.tools = append(s.tools, tool)
	s.handlers[tool.Name] = handler
}

// Run 启动 MCP 服务器主循环
func (s *MCPServer) Run() error {
	log.SetOutput(os.Stderr) // 日志输出到 stderr，stdout 给 JSON-RPC 用
	log.Println("[MCP] Server starting...")

	scanner := bufio.NewScanner(s.reader)
	// 增大缓冲区，防止大消息截断
	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		// 复制一份数据避免 scanner 覆盖
		msg := make([]byte, len(line))
		copy(msg, line)
		go s.handleMessage(msg)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	log.Println("[MCP] Server stopped.")
	return nil
}

// handleMessage 处理单条消息
func (s *MCPServer) handleMessage(data []byte) {
	var req JSONRPCRequest
	if err := json.Unmarshal(data, &req); err != nil {
		log.Printf("[MCP] Failed to parse message: %v", err)
		return
	}

	log.Printf("[MCP] Received method: %s", req.Method)

	// 如果没有 id，说明是通知，不需要回复
	if req.ID == nil || string(req.ID) == "null" {
		s.handleNotification(req.Method, req.Params)
		return
	}

	// 有 id，需要回复
	var result interface{}
	var rpcErr *JSONRPCError

	switch req.Method {
	case "initialize":
		result = s.handleInitialize()
	case "tools/list":
		result = s.handleToolsList()
	case "tools/call":
		result, rpcErr = s.handleToolsCall(req.Params)
	case "ping":
		result = map[string]interface{}{}
	default:
		rpcErr = &JSONRPCError{
			Code:    -32601,
			Message: fmt.Sprintf("method not found: %s", req.Method),
		}
	}

	s.sendResponse(req.ID, result, rpcErr)
}

// handleNotification 处理通知消息
func (s *MCPServer) handleNotification(method string, params json.RawMessage) {
	switch method {
	case "notifications/initialized":
		log.Println("[MCP] Client initialized notification received")
	case "notifications/cancelled":
		log.Println("[MCP] Request cancelled notification received")
	default:
		log.Printf("[MCP] Unknown notification: %s", method)
	}
}

// handleInitialize 处理初始化请求
func (s *MCPServer) handleInitialize() *InitializeResult {
	return &InitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities:    s.capabilities,
		ServerInfo:      s.serverInfo,
		Instructions:    s.instructions,
	}
}

// handleToolsList 处理工具列表请求
func (s *MCPServer) handleToolsList() *ToolsListResult {
	return &ToolsListResult{
		Tools: s.tools,
	}
}

// handleToolsCall 处理工具调用请求
func (s *MCPServer) handleToolsCall(params json.RawMessage) (*CallToolResult, *JSONRPCError) {
	var callParams CallToolParams
	if err := json.Unmarshal(params, &callParams); err != nil {
		return nil, &JSONRPCError{
			Code:    -32602,
			Message: fmt.Sprintf("invalid params: %v", err),
		}
	}

	handler, ok := s.handlers[callParams.Name]
	if !ok {
		return nil, &JSONRPCError{
			Code:    -32602,
			Message: fmt.Sprintf("unknown tool: %s", callParams.Name),
		}
	}

	result := handler(callParams.Arguments)
	return result, nil
}

// sendResponse 发送响应
func (s *MCPServer) sendResponse(id json.RawMessage, result interface{}, rpcErr *JSONRPCError) {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
	}

	if rpcErr != nil {
		resp.Error = rpcErr
	} else {
		resp.Result = result
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("[MCP] Failed to marshal response: %v", err)
		return
	}

	// MCP 使用 newline-delimited JSON
	data = append(data, '\n')
	if _, err := s.writer.Write(data); err != nil {
		log.Printf("[MCP] Failed to write response: %v", err)
	}
}
