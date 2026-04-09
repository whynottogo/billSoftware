package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
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

