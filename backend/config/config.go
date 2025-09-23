package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config 应用程序配置结构
type Config struct {
	Server     ServerConfig     `json:"server"`
	Volcengine VolcengineConfig `json:"volcengine"`
	Google     GoogleConfig     `json:"google"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port                int `json:"port"`
	ReadTimeoutSeconds  int `json:"readTimeoutSeconds"`
	WriteTimeoutSeconds int `json:"writeTimeoutSeconds"`
}

// VolcengineConfig 火山引擎翻译API配置
type VolcengineConfig struct {
	AccessKey    string `json:"accessKey"`
	SecretKey    string `json:"secretKey"`
	Region       string `json:"region"`
	Endpoint     string `json:"endpoint"`
	TranslateURL string `json:"translateURL"`
}

// GoogleConfig Google翻译API配置
type GoogleConfig struct {
	APIKey string `json:"apiKey"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:                8080,
			ReadTimeoutSeconds:  10,
			WriteTimeoutSeconds: 10,
		},
		Volcengine: VolcengineConfig{
			Region:       "cn-north-1",
			Endpoint:     "open.volcengineapi.com",
			TranslateURL: "https://translate.volcengineapi.com",
		},
		Google: GoogleConfig{},
	}
}

// Load 加载配置文件
func Load() (*Config, error) {
	// 首先使用默认配置
	cfg := DefaultConfig()

	// 尝试从配置文件加载
	configPath := getConfigPath()
	if _, err := os.Stat(configPath); err == nil {
		file, err := os.Open(configPath)
		if err != nil {
			return nil, fmt.Errorf("无法打开配置文件: %w", err)
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(cfg); err != nil {
			return nil, fmt.Errorf("解析配置文件失败: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("检查配置文件状态失败: %w", err)
	}

	// 从环境变量覆盖配置
	overrideFromEnv(cfg)

	return cfg, nil
}

// getConfigPath 获取配置文件路径
func getConfigPath() string {
	// 首先检查环境变量中是否指定了配置文件路径
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		return path
	}

	// 默认使用当前目录下的config.json
	execPath, err := os.Executable()
	if err != nil {
		return "config.json"
	}

	execDir := filepath.Dir(execPath)
	return filepath.Join(execDir, "config.json")
}

// overrideFromEnv 从环境变量覆盖配置
func overrideFromEnv(cfg *Config) {
	// 服务器配置
	if port := os.Getenv("SERVER_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &cfg.Server.Port)
	}

	// 火山引擎配置
	if key := os.Getenv("VOLCENGINE_ACCESS_KEY"); key != "" {
		cfg.Volcengine.AccessKey = key
	}
	if key := os.Getenv("VOLCENGINE_SECRET_KEY"); key != "" {
		cfg.Volcengine.SecretKey = key
	}
	if region := os.Getenv("VOLCENGINE_REGION"); region != "" {
		cfg.Volcengine.Region = region
	}
	if endpoint := os.Getenv("VOLCENGINE_ENDPOINT"); endpoint != "" {
		cfg.Volcengine.Endpoint = endpoint
	}
	if url := os.Getenv("VOLCENGINE_TRANSLATE_URL"); url != "" {
		cfg.Volcengine.TranslateURL = url
	}

	// Google配置
	if key := os.Getenv("GOOGLE_API_KEY"); key != "" {
		cfg.Google.APIKey = key
	}
}