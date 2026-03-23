package configs

import (
	"fmt"
	"os"
	"strings"
)

// Config MCP Server 配置
type Config struct {
	// API 服务地址（你现有 HTTP API 的地址）
	APIBaseURL string
	// API Key（用户的认证密钥）
	APIKey string
}

// LoadConfig 从环境变量加载配置
func LoadConfig() (*Config, error) {
	cfg := &Config{
		APIBaseURL: getEnv("HITPAW_API_BASE_URL", "https://api-base.hitpaw.com"),
		APIKey:     getEnv("HITPAW_API_KEY", ""),
	}

	// API Key 是必须的
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("HITPAW_API_KEY 环境变量未设置。请设置后重试，例如：\n  export HITPAW_API_KEY=your_api_key_here")
	}

	// 去掉末尾的 /
	cfg.APIBaseURL = strings.TrimRight(cfg.APIBaseURL, "/")

	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
