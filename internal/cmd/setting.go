package cmd

import (
	"strconv"

	"github.com/go-tapd/cli/internal/app"
	"github.com/go-tapd/tapd"
	"github.com/spf13/cobra"
)

func newSettingCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setting",
		Short: "Work with TAPD settings",
	}

	cmd.AddCommand(newSettingWorkspaceCmd(rt))
	return cmd
}

func newSettingWorkspaceCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		settingType string
	)

	cmd := &cobra.Command{
		Use:   "workspace",
		Short: "Show workspace settings",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetWorkspaceSettingRequest{
				WorkspaceID: new(workspaceID),
			}
			if settingType != "" {
				request.Type = new(settingType)
			}

			settings, _, err := client.SettingService.GetWorkspaceSetting(cmd.Context(), request)
			if err != nil {
				return err
			}

			storyCategory := ""
			if settings.IsEnabledStoryCategory != nil {
				storyCategory = strconv.Itoa(*settings.IsEnabledStoryCategory)
			}
			metrology := ""
			if settings.WorkspaceMetrology != nil {
				metrology = *settings.WorkspaceMetrology
			}

			rows := [][]string{{storyCategory, metrology}}
			return writeOutput(cmd, rt.OutputFormat, []string{"StoryCategoryEnabled", "Metrology"}, rows, settings)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().StringVar(&settingType, "type", "", "setting name")
	return cmd
}
