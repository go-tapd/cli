package cmd

import (
	"github.com/go-tapd/cli/internal/app"
	"github.com/go-tapd/tapd"
	"github.com/spf13/cobra"
)

func newUserCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Work with TAPD users",
	}

	cmd.AddCommand(newUserRolesCmd(rt))
	return cmd
}

func newUserRolesCmd(rt *app.Runtime) *cobra.Command {
	var workspaceID int

	cmd := &cobra.Command{
		Use:   "roles",
		Short: "List user roles",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			roles, _, err := client.UserService.GetRoles(cmd.Context(), &tapd.GetRolesRequest{
				WorkspaceID: new(workspaceID),
			})
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(roles))
			for _, item := range roles {
				rows = append(rows, []string{item.ID, item.Name})
			}

			return writeOutput(cmd, rt.OutputFormat, []string{"ID", "Name"}, rows, roles)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	return cmd
}
