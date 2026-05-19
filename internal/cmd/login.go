package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-tapd/cli/internal/app"
	"github.com/go-tapd/tapd"
	"github.com/spf13/cobra"
)

func newLoginCmd(rt *app.Runtime) *cobra.Command {
	var (
		authMethod   string
		token        string
		clientID     string
		clientSecret string
		workspaceID  int
	)

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Log in to TAPD",
		RunE: func(cmd *cobra.Command, _ []string) error {
			method, auth, err := resolveLoginAuth(authMethod, token, clientID, clientSecret)
			if err != nil {
				return err
			}

			clientRT := *rt
			clientRT.BaseURLOverride = strings.TrimSpace(clientRT.BaseURLOverride)

			client, _, err := newClientForLogin(&clientRT, auth)
			if err != nil {
				return err
			}
			if workspaceID > 0 {
				if err := validateLoginWithWorkspace(cmd.Context(), client, workspaceID); err != nil {
					return err
				}
			}

			cfg := &app.Config{
				BaseURL: loginBaseURL(clientRT.BaseURLOverride),
				Auth: app.AuthConfig{
					Method: method,
				},
			}
			switch method {
			case app.AuthMethodPAT:
				cfg.Auth.AccessToken = auth.AccessToken
			case app.AuthMethodBasic:
				cfg.Auth.ClientID = auth.ClientID
				cfg.Auth.ClientSecret = auth.ClientSecret
			}

			if err := app.SaveConfig(rt.ConfigPath, cfg); err != nil {
				return err
			}

			if workspaceID > 0 {
				fmt.Fprintf( //nolint:errcheck
					cmd.OutOrStdout(),
					"Logged in with %s and validated against workspace %d. Config saved to %s\n",
					method,
					workspaceID,
					rt.ConfigPath,
				)
				return nil
			}

			fmt.Fprintf( //nolint:errcheck
				cmd.OutOrStdout(),
				"Logged in with %s. Config saved to %s\n",
				method,
				rt.ConfigPath,
			)
			return nil
		},
	}

	cmd.Flags().StringVar(&authMethod, "auth-method", "", "auth method: pat|basic")
	cmd.Flags().StringVar(&token, "token", "", "personal access token")
	cmd.Flags().StringVar(&clientID, "client-id", "", "basic auth client ID")
	cmd.Flags().StringVar(&clientSecret, "client-secret", "", "basic auth client secret")
	cmd.Flags().IntVar(&workspaceID, "workspace-id", 0, "validate login against a workspace ID")

	return cmd
}

func resolveLoginAuth(authMethod, token, clientID, clientSecret string) (app.AuthMethod, app.AuthConfig, error) {
	method := strings.TrimSpace(authMethod)

	switch {
	case token != "":
		if clientID != "" || clientSecret != "" {
			return "", app.AuthConfig{}, fmt.Errorf("token auth cannot be mixed with basic auth flags")
		}
		return app.AuthMethodPAT, app.AuthConfig{
			Method:      app.AuthMethodPAT,
			AccessToken: token,
		}, nil
	case clientID != "" || clientSecret != "":
		if clientID == "" || clientSecret == "" {
			return "", app.AuthConfig{}, fmt.Errorf("client-id and client-secret must be set together")
		}
		return app.AuthMethodBasic, app.AuthConfig{
			Method:       app.AuthMethodBasic,
			ClientID:     clientID,
			ClientSecret: clientSecret,
		}, nil
	}

	switch method {
	case "", string(app.AuthMethodPAT):
		if method == "" {
			method = string(app.AuthMethodPAT)
		}
	case string(app.AuthMethodBasic):
	default:
		return "", app.AuthConfig{}, fmt.Errorf("unsupported auth method %q", authMethod)
	}

	if method == string(app.AuthMethodPAT) {
		value, err := app.PromptSecret("TAPD personal access token")
		if err != nil {
			return "", app.AuthConfig{}, err
		}
		if value == "" {
			return "", app.AuthConfig{}, fmt.Errorf("token cannot be empty")
		}
		return app.AuthMethodPAT, app.AuthConfig{
			Method:      app.AuthMethodPAT,
			AccessToken: value,
		}, nil
	}

	id, err := app.Prompt("TAPD client ID")
	if err != nil {
		return "", app.AuthConfig{}, err
	}
	secret, err := app.PromptSecret("TAPD client secret")
	if err != nil {
		return "", app.AuthConfig{}, err
	}
	if id == "" || secret == "" {
		return "", app.AuthConfig{}, fmt.Errorf("client-id and client-secret cannot be empty")
	}

	return app.AuthMethodBasic, app.AuthConfig{
		Method:       app.AuthMethodBasic,
		ClientID:     id,
		ClientSecret: secret,
	}, nil
}

func newClientForLogin(rt *app.Runtime, auth app.AuthConfig) (*tapd.Client, *app.ResolvedAuth, error) {
	resolved := &app.ResolvedAuth{
		Source:  app.AuthSourceConfig,
		BaseURL: loginBaseURL(rt.BaseURLOverride),
		Auth:    auth,
	}

	opts := []tapd.ClientOption{
		tapd.WithBaseURL(resolved.BaseURL),
		tapd.WithUserAgent(app.UserAgent),
	}

	switch auth.Method {
	case app.AuthMethodPAT:
		client, err := tapd.NewPATClient(auth.AccessToken, opts...)
		return client, resolved, err
	case app.AuthMethodBasic:
		client, err := tapd.NewBasicAuthClient(auth.ClientID, auth.ClientSecret, opts...)
		return client, resolved, err
	default:
		return nil, nil, fmt.Errorf("unsupported auth method %q", auth.Method)
	}
}

func loginBaseURL(override string) string {
	if override != "" {
		return override
	}

	return app.DefaultBaseURL
}

func validateLoginWithWorkspace(ctx context.Context, client *tapd.Client, workspaceID int) error {
	_, _, err := client.UserService.GetRoles(ctx, &tapd.GetRolesRequest{
		WorkspaceID: tapd.Ptr(workspaceID),
	})
	if err != nil {
		return fmt.Errorf("validate credentials with workspace %d: %w", workspaceID, err)
	}

	return nil
}
