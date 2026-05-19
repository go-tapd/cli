package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	ConfigDirName  = ".tapd"
	ConfigFileName = "config.json"
)

type AuthMethod string

const (
	AuthMethodPAT   AuthMethod = "pat"
	AuthMethodBasic AuthMethod = "basic"
)

type AuthConfig struct {
	Method       AuthMethod `json:"method,omitempty"`
	AccessToken  string     `json:"access_token,omitempty"`
	ClientID     string     `json:"client_id,omitempty"`
	ClientSecret string     `json:"client_secret,omitempty"`
}

type Config struct {
	BaseURL string     `json:"base_url,omitempty"`
	Auth    AuthConfig `json:"auth,omitzero"`
}

func DefaultConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get user home dir: %w", err)
	}

	return filepath.Join(homeDir, ConfigDirName, ConfigFileName), nil
}

func ResolveConfigPath(custom string) (string, error) {
	if custom != "" {
		return custom, nil
	}

	return DefaultConfigPath()
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return &Config{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	cfg := new(Config)
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}

	return cfg, nil
}

func SaveConfig(path string, cfg *Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("encode config: %w", err)
	}
	data = append(data, '\n')

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}

func RemoveConfig(path string) error {
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("remove config: %w", err)
	}

	return nil
}

func (c AuthConfig) Valid() bool {
	switch c.Method {
	case AuthMethodPAT:
		return c.AccessToken != ""
	case AuthMethodBasic:
		return c.ClientID != "" && c.ClientSecret != ""
	default:
		return false
	}
}

func (c AuthConfig) RedactedSummary() string {
	switch c.Method {
	case AuthMethodPAT:
		return fmt.Sprintf("pat %s", redact(c.AccessToken))
	case AuthMethodBasic:
		return fmt.Sprintf("basic client_id=%s client_secret=%s", redact(c.ClientID), redact(c.ClientSecret))
	default:
		return "not configured"
	}
}

func redact(s string) string {
	if s == "" {
		return "<empty>"
	}
	if len(s) <= 8 {
		return "********"
	}

	return s[:4] + "..." + s[len(s)-4:]
}
