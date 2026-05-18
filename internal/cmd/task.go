package cmd

import (
	"strconv"

	"github.com/go-tapd/cli/internal/app"
	"github.com/go-tapd/tapd"
	"github.com/spf13/cobra"
)

func newTaskCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "Work with TAPD tasks",
	}

	cmd.AddCommand(
		newTaskListCmd(rt),
		newTaskCountCmd(rt),
		newTaskFieldsCmd(rt),
	)
	return cmd
}

func newTaskListCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		limit       int
		page        int
		fields      string
		ids         string
		creator     string
		owner       string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List tasks",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetTasksRequest{
				WorkspaceID: tapd.Ptr(workspaceID),
				Limit:       tapd.Ptr(limit),
				Page:        tapd.Ptr(page),
				Fields:      fieldsMulti(fields),
				ID:          int64Multi(ids),
			}
			if creator != "" {
				request.Creator = tapd.Ptr(creator)
			}
			if owner != "" {
				request.Owner = tapd.Ptr(owner)
			}

			items, _, err := client.TaskService.GetTasks(cmd.Context(), request)
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(items))
			for _, item := range items {
				rows = append(rows, []string{
					item.ID,
					item.Name,
					item.Status.String(),
					item.Owner,
					item.Creator,
					item.Progress,
				})
			}

			return writeOutput(cmd, rt.OutputFormat, []string{"ID", "Name", "Status", "Owner", "Creator", "Progress"}, rows, items)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	newListFlags(cmd, &limit, &page)
	cmd.Flags().StringVar(&fields, "fields", "", "comma separated fields")
	cmd.Flags().StringVar(&ids, "ids", "", "comma separated task IDs")
	cmd.Flags().StringVar(&creator, "creator", "", "filter by creator")
	cmd.Flags().StringVar(&owner, "owner", "", "filter by owner")
	return cmd
}

func newTaskCountCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		ids         string
		creator     string
		owner       string
	)

	cmd := &cobra.Command{
		Use:   "count",
		Short: "Count tasks",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetTasksCountRequest{
				WorkspaceID: tapd.Ptr(workspaceID),
				ID:          int64Multi(ids),
			}
			if creator != "" {
				request.Creator = tapd.Ptr(creator)
			}
			if owner != "" {
				request.Owner = tapd.Ptr(owner)
			}

			count, _, err := client.TaskService.GetTasksCount(cmd.Context(), request)
			if err != nil {
				return err
			}

			data := map[string]any{
				"resource":     "task",
				"workspace_id": workspaceID,
				"count":        count,
			}
			rows := [][]string{{"task", strconv.Itoa(workspaceID), strconv.Itoa(count)}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Resource", "WorkspaceID", "Count"}, rows, data)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().StringVar(&ids, "ids", "", "comma separated task IDs")
	cmd.Flags().StringVar(&creator, "creator", "", "filter by creator")
	cmd.Flags().StringVar(&owner, "owner", "", "filter by owner")
	return cmd
}

func newTaskFieldsCmd(rt *app.Runtime) *cobra.Command {
	var workspaceID int

	cmd := &cobra.Command{
		Use:   "fields",
		Short: "List task fields",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			fields, _, err := client.TaskService.GetTaskFieldsInfo(cmd.Context(), &tapd.GetTaskFieldsInfoRequest{
				WorkspaceID: tapd.Ptr(workspaceID),
			})
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(fields))
			for _, field := range fields {
				rows = append(rows, []string{
					field.Name,
					field.Label,
					string(field.HTMLType),
					strconv.Itoa(len(field.Options)),
					strconv.Itoa(len(field.ColorOptions)),
					strconv.Itoa(len(field.PureOptions)),
				})
			}

			return writeOutput(
				cmd,
				rt.OutputFormat,
				[]string{"Name", "Label", "Type", "Options", "ColorOptions", "PureOptions"},
				rows,
				fields,
			)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	return cmd
}
