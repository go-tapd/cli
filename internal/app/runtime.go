package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/go-tapd/tapd"
)

const (
	DefaultBaseURL = "https://api.tapd.cn/"
	UserAgent      = "tapd-cli"
)

type Runtime struct {
	Stdout          io.Writer
	Stderr          io.Writer
	ConfigPath      string
	BaseURLOverride string
	OutputFormat    string
}

func NewRuntime() *Runtime {
	return &Runtime{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

type AuthSource string

const (
	AuthSourceEnv    AuthSource = "env"
	AuthSourceConfig AuthSource = "config"
)

type ResolvedAuth struct {
	Source  AuthSource
	BaseURL string
	Auth    AuthConfig
}

func (r *Runtime) ResolveAuth() (*ResolvedAuth, error) {
	if token := os.Getenv("TAPD_ACCESS_TOKEN"); token != "" {
		baseURL := valueOrDefault(os.Getenv("TAPD_BASE_URL"), r.BaseURLOverride, DefaultBaseURL)
		return &ResolvedAuth{
			Source:  AuthSourceEnv,
			BaseURL: baseURL,
			Auth: AuthConfig{
				Method:      AuthMethodPAT,
				AccessToken: token,
			},
		}, nil
	}

	clientID := os.Getenv("TAPD_CLIENT_ID")
	clientSecret := os.Getenv("TAPD_CLIENT_SECRET")
	if clientID != "" || clientSecret != "" {
		if clientID == "" || clientSecret == "" {
			return nil, errors.New("TAPD_CLIENT_ID and TAPD_CLIENT_SECRET must be set together")
		}

		baseURL := valueOrDefault(os.Getenv("TAPD_BASE_URL"), r.BaseURLOverride, DefaultBaseURL)
		return &ResolvedAuth{
			Source:  AuthSourceEnv,
			BaseURL: baseURL,
			Auth: AuthConfig{
				Method:       AuthMethodBasic,
				ClientID:     clientID,
				ClientSecret: clientSecret,
			},
		}, nil
	}

	cfg, err := LoadConfig(r.ConfigPath)
	if err != nil {
		return nil, err
	}
	if !cfg.Auth.Valid() {
		return nil, fmt.Errorf("no TAPD credentials configured, run `tapd login` first")
	}

	return &ResolvedAuth{
		Source:  AuthSourceConfig,
		BaseURL: valueOrDefault(r.BaseURLOverride, cfg.BaseURL, DefaultBaseURL),
		Auth:    cfg.Auth,
	}, nil
}

func (r *Runtime) NewClient() (*tapd.Client, *ResolvedAuth, error) {
	resolved, err := r.ResolveAuth()
	if err != nil {
		return nil, nil, err
	}

	opts := []tapd.ClientOption{
		tapd.WithBaseURL(resolved.BaseURL),
		tapd.WithUserAgent(UserAgent),
	}

	switch resolved.Auth.Method {
	case AuthMethodPAT:
		client, err := tapd.NewPATClient(resolved.Auth.AccessToken, opts...)
		return client, resolved, err
	case AuthMethodBasic:
		client, err := tapd.NewBasicAuthClient(resolved.Auth.ClientID, resolved.Auth.ClientSecret, opts...)
		return client, resolved, err
	default:
		return nil, nil, fmt.Errorf("unsupported auth method %q", resolved.Auth.Method)
	}
}

func ValidateCredentials(ctx context.Context, client *tapd.Client) error {
	_, _, err := client.UserService.GetRoles(ctx, &tapd.GetRolesRequest{})
	if err != nil {
		return fmt.Errorf("validate credentials: %w", err)
	}

	return nil
}

func valueOrDefault(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}

	return ""
}
