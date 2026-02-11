package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config 应用配置结构体
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	JWT      JWTConfig      `json:"jwt"`
	Redis    RedisConfig    `json:"redis"`
	Upload   UploadConfig   `json:"upload"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int `json:"port"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret string `json:"secret"`
	Expire int    `json:"expire"` // 过期时间（小时）
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// UploadConfig 上传配置
type UploadConfig struct {
	Path         string   `json:"path"`
	MaxSize      int64    `json:"max_size"`
	AllowedTypes []string `json:"allowed_types"`
}

// LoadConfig 加载配置
func LoadConfig() (*Config, error) {
	// 默认配置
	config := &Config{
		Server: ServerConfig{
			Port: 8080,
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "password",
			Name:     "nookverse",
		},
		JWT: JWTConfig{
			Secret: "your-jwt-secret-key-here-change-in-production",
			Expire: 24,
		},
		Redis: RedisConfig{
			Host:     "localhost",
			Port:     6379,
			Password: "",
			DB:       0,
		},
		Upload: UploadConfig{
			Path:    "./uploads",
			MaxSize: 10485760, // 10MB
			AllowedTypes: []string{
				"image/jpeg",
				"image/png",
				"image/gif",
				"video/mp4",
			},
		},
	}

	// 尝试从环境变量加载配置文件路径
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.json"
	}

	// 如果配置文件存在，则加载它
	if _, err := os.Stat(configPath); err == nil {
		file, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := json.Unmarshal(file, config); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	// 从环境变量覆盖配置
	if port := os.Getenv("SERVER_PORT"); port != "" {
		// 这里可以解析端口号
	}

	return config, nil
}