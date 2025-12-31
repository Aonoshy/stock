package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server" json:"server"`
	Alltick  AlltickConfig  `yaml:"alltick" json:"alltick"`
	Database DatabaseConfig `yaml:"database" json:"database"`
}

type ServerConfig struct {
	Port string `yaml:"port" json:"port"`
	Host string `yaml:"host" json:"host"`
}

type AlltickConfig struct {
	APIKey  string `yaml:"api_key" json:"api_key"`
	BaseURL string `yaml:"base_url" json:"base_url"`
}

type DatabaseConfig struct {
	Type string `yaml:"type" json:"type"`
	URL  string `yaml:"url" json:"url"`
}

// NewConfig 创建配置实例，支持配置文件和环境变量
func NewConfig() *Config {
	config := &Config{
		Server: ServerConfig{
			Port: "8080",
			Host: "0.0.0.0",
		},
		Alltick: AlltickConfig{
			BaseURL: "https://quote.alltick.co/quote-stock-b-api",
		},
		Database: DatabaseConfig{
			Type: "memory",
		},
	}

	// 尝试从配置文件加载
	if err := loadConfigFile(config); err != nil {
		fmt.Printf("配置文件加载失败，使用默认配置: %v\n", err)
	}

	// 环境变量覆盖配置文件设置
	loadFromEnv(config)

	return config
}

// loadConfigFile 从YAML文件加载配置
func loadConfigFile(config *Config) error {
	// 尝试多个可能的配置文件路径
	configPaths := []string{
		"config.yaml",
		"config.yml",
		"configs/config.yaml",
		"configs/config.yml",
		"/etc/stock-api/config.yaml",
		filepath.Join(os.Getenv("HOME"), ".stock-api", "config.yaml"),
	}

	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				continue
			}

			if err := yaml.Unmarshal(data, config); err != nil {
				return fmt.Errorf("解析配置文件 %s 失败: %v", path, err)
			}

			fmt.Printf("成功加载配置文件: %s\n", path)
			return nil
		}
	}

	return fmt.Errorf("未找到配置文件")
}

// loadFromEnv 从环境变量加载配置（优先级最高）
func loadFromEnv(config *Config) {
	if port := os.Getenv("PORT"); port != "" {
		config.Server.Port = port
	}
	if host := os.Getenv("HOST"); host != "" {
		config.Server.Host = host
	}
	if apiKey := os.Getenv("ALLTICK_API_KEY"); apiKey != "" {
		config.Alltick.APIKey = apiKey
	}
	if baseURL := os.Getenv("ALLTICK_BASE_URL"); baseURL != "" {
		config.Alltick.BaseURL = baseURL
	}
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		config.Database.URL = dbURL
		config.Database.Type = "external"
	}
}

// GetAlltickAPIKey 获取Alltick API密钥（向后兼容）
func (c *Config) GetAlltickAPIKey() string {
	return c.Alltick.APIKey
}

// GetPort 获取服务端口（向后兼容）
func (c *Config) GetPort() string {
	return c.Server.Port
}