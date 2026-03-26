package cmd

import (
	"strconv"

	"github.com/go-tapd/cli/internal/app"
	"github.com/go-tapd/tapd"
	"github.com/spf13/cobra"
)

func newWorkspaceCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workspace",
		Short: "Work with TAPD workspaces",
	}

	cmd.AddCommand(
		newWorkspaceViewCmd(rt),
		newWorkspaceUsersCmd(rt),
	)

	return cmd
}

func newWorkspaceViewCmd(rt *app.Runtime) *cobra.Command {
	var workspaceID int

	cmd := &cobra.Command{
		Use:   "view",
		Short: "Show workspace details",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			workspace, _, err := client.WorkspaceService.GetWorkspaceInfo(cmd.Context(), &tapd.GetWorkspaceInfoRequest{
				WorkspaceID: tapd.Ptr(workspaceID),
			})
			if err != nil {
				return err
			}

			rows := [][]string{{
				workspace.ID,
				workspace.Name,
				workspace.PrettyName,
				workspace.Status,
				workspace.Creator,
			}}

			return writeOutput(cmd, rt.OutputFormat, []string{"ID", "Name", "PrettyName", "Status", "Creator"}, rows, workspace)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	return cmd
}

func newWorkspaceUsersCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		fields      string
		user        string
	)

	cmd := &cobra.Command{
		Use:   "users",
		Short: "List workspace users",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetUsersRequest{
				WorkspaceID: tapd.Ptr(workspaceID),
				Fields:      fieldsMulti(fields),
			}
			if user != "" {
				request.User = tapd.NewMulti(user)
			}

			users, _, err := client.WorkspaceService.GetUsers(cmd.Context(), request)
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(users))
			for _, item := range users {
				rows = append(rows, []string{
					item.User,
					item.Name,
					strconv.Itoa(len(item.RoleID)),
					item.Status,
					item.RealJoinTime,
				})
			}

			return writeOutput(cmd, rt.OutputFormat, []string{"User", "Name", "Roles", "Status", "Joined"}, rows, users)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().StringVar(&fields, "fields", "", "comma separated fields")
	cmd.Flags().StringVar(&user, "user", "", "filter by TAPD user")
	return cmd
}
