package cmd

import (
	"strconv"

	"github.com/go-tapd/cli/internal/app"
	"github.com/go-tapd/tapd"
	"github.com/spf13/cobra"
)

func newWorkflowCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow",
		Short: "Work with TAPD workflows",
	}

	cmd.AddCommand(newWorkflowLastStepsCmd(rt))
	return cmd
}

func newWorkflowLastStepsCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		system      string
		groupKey    string
	)

	cmd := &cobra.Command{
		Use:   "last-steps",
		Short: "List workflow final statuses",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetAllLastStepsRequest{
				WorkspaceID: new(workspaceID),
			}
			if system != "" {
				request.System = new(system)
			}
			if groupKey != "" {
				request.GroupKey = new(groupKey)
			}

			steps, _, err := client.WorkflowService.GetAllLastSteps(cmd.Context(), request)
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(steps))
			for _, item := range steps {
				rows = append(rows, []string{item.Key, strconv.Itoa(len(item.Status))})
			}

			return writeOutput(cmd, rt.OutputFormat, []string{"Key", "StatusCount"}, rows, steps)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().StringVar(&system, "system", "story", "workflow system name")
	cmd.Flags().StringVar(&groupKey, "group-key", "", "group by workflow_id or workitem_type_id")
	return cmd
}
