package cmd

import (
	"fmt"

	"github.com/go-tapd/cli/internal/app"
	"github.com/spf13/cobra"
)

func newAuthCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage TAPD authentication",
	}

	cmd.AddCommand(
		newAuthStatusCmd(rt),
		newLogoutCmd(rt),
	)

	return cmd
}

func newAuthStatusCmd(rt *app.Runtime) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show authentication status",
		RunE: func(cmd *cobra.Command, _ []string) error {
			resolved, err := rt.ResolveAuth()
			if err != nil {
				return err
			}

			status := map[string]any{
				"source":   resolved.Source,
				"base_url": resolved.BaseURL,
				"auth": map[string]any{
					"method":  resolved.Auth.Method,
					"summary": resolved.Auth.RedactedSummary(),
				},
				"config_path": rt.ConfigPath,
			}

			if rt.OutputFormat == "json" {
				return app.WriteJSON(cmd.OutOrStdout(), status)
			}

			_, err = fmt.Fprintf(
				cmd.OutOrStdout(),
				"source: %s\nbase_url: %s\nauth: %s\nconfig_path: %s\n",
				resolved.Source,
				resolved.BaseURL,
				resolved.Auth.RedactedSummary(),
				rt.ConfigPath,
			)
			return err
		},
	}
}

func newLogoutCmd(rt *app.Runtime) *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Remove stored TAPD credentials",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := app.RemoveConfig(rt.ConfigPath); err != nil {
				return err
			}

			_, err := fmt.Fprintf(cmd.OutOrStdout(), "Removed %s\n", rt.ConfigPath)
			return err
		},
	}
}
