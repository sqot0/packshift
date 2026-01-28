package config

import (
	"encoding/json"
	"errors"
	"os"
	"path"
)

type FTPConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	SSL      bool   `json:"ssl"`
}

type Config struct {
	FTPConfig    FTPConfig         `json:"ftpConfig"`
	PathMappings map[string]string `json:"pathMappings"`
	path         string
}

func Init(projectPath string, ftp *FTPConfig, pathMap map[string]string) error {
	cfg := Config{*ftp, pathMap, projectPath}
	err := write(cfg)
	if err != nil {
		return err
	}
	return nil
}

func Load(projectPath string) (*Config, error) {
	cfgFile := path.Join(projectPath, "packshift.json")
	data, err := os.ReadFile(cfgFile)
	if err != nil {
		return nil, errors.New("initialize project before using other commands")
	}

	var cfg Config
	cfg.path = projectPath
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func write(cfg Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	cfgFile := path.Join(cfg.path, "packshift.json")
	err = os.WriteFile(cfgFile, data, 0o644)
	if err != nil {
		return err
	}
	return nil
}
