package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	MinIO    MinIOConfig    `yaml:"minio"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

type DatabaseConfig struct {
	Driver       string `yaml:"driver"`
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	Name         string `yaml:"name"`
	Charset      string `yaml:"charset"`
	ParseTime    bool   `yaml:"parse_time"`
	Loc          string `yaml:"loc"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
}

type MinIOConfig struct {
	Endpoint       string `yaml:"endpoint"`
	AccessKey      string `yaml:"access_key"`
	SecretKey      string `yaml:"secret_key"`
	Bucket         string `yaml:"bucket"`
	UseSSL         bool   `yaml:"use_ssl"`
	PresignMinutes int    `yaml:"presign_minutes"`
	ThumbnailSize  int    `yaml:"thumbnail_size"`
	MaxUploadMB    int64  `yaml:"max_upload_mb"`
}

func (s ServerConfig) Address() string {
	return fmt.Sprintf(":%d", s.Port)
}

func LoadConfig(configPath string) (*AppConfig, error) {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	cfg := &AppConfig{}
	if err := yaml.Unmarshal(content, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
