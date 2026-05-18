package cmd

import (
	"strconv"

	"github.com/go-tapd/cli/internal/app"
	"github.com/go-tapd/tapd"
	"github.com/spf13/cobra"
)

func newStoryCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "story",
		Short: "Work with TAPD stories",
	}

	cmd.AddCommand(
		newStoryListCmd(rt),
		newStoryCountCmd(rt),
		newStoryFieldsCmd(rt),
	)
	return cmd
}

func newStoryListCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		limit       int
		page        int
		fields      string
		ids         string
		creator     string
		owner       string
		status      string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List stories",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetStoriesRequest{
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
			if status != "" {
				request.VStatus = tapd.Ptr(status)
				request.WithVStatus = tapd.Ptr("1")
			}

			stories, _, err := client.StoryService.GetStories(cmd.Context(), request)
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(stories))
			for _, item := range stories {
				rows = append(rows, []string{
					item.ID,
					item.Name,
					string(item.Status),
					item.Owner,
					item.Creator,
					item.Modified,
				})
			}

			return writeOutput(cmd, rt.OutputFormat, []string{"ID", "Name", "Status", "Owner", "Creator", "Modified"}, rows, stories)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	newListFlags(cmd, &limit, &page)
	cmd.Flags().StringVar(&fields, "fields", "", "comma separated fields")
	cmd.Flags().StringVar(&ids, "ids", "", "comma separated story IDs")
	cmd.Flags().StringVar(&creator, "creator", "", "filter by creator")
	cmd.Flags().StringVar(&owner, "owner", "", "filter by owner")
	cmd.Flags().StringVar(&status, "status", "", "filter by status or localized status text")

	return cmd
}

func newStoryCountCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		ids         string
		creator     string
		owner       string
		status      string
	)

	cmd := &cobra.Command{
		Use:   "count",
		Short: "Count stories",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetStoriesCountRequest{
				WorkspaceID: tapd.Ptr(workspaceID),
				ID:          int64Multi(ids),
			}
			if creator != "" {
				request.Creator = tapd.Ptr(creator)
			}
			if owner != "" {
				request.Owner = tapd.Ptr(owner)
			}
			if status != "" {
				request.VStatus = tapd.Ptr(status)
				request.WithVStatus = tapd.Ptr("1")
			}

			count, _, err := client.StoryService.GetStoriesCount(cmd.Context(), request)
			if err != nil {
				return err
			}

			data := map[string]any{
				"resource":     "story",
				"workspace_id": workspaceID,
				"count":        count,
			}
			rows := [][]string{{"story", strconv.Itoa(workspaceID), strconv.Itoa(count)}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Resource", "WorkspaceID", "Count"}, rows, data)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().StringVar(&ids, "ids", "", "comma separated story IDs")
	cmd.Flags().StringVar(&creator, "creator", "", "filter by creator")
	cmd.Flags().StringVar(&owner, "owner", "", "filter by owner")
	cmd.Flags().StringVar(&status, "status", "", "filter by status or localized status text")
	return cmd
}

func newStoryFieldsCmd(rt *app.Runtime) *cobra.Command {
	var workspaceID int

	cmd := &cobra.Command{
		Use:   "fields",
		Short: "List story fields",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			fields, _, err := client.StoryService.GetStoryFieldsInfo(cmd.Context(), &tapd.GetStoryFieldsInfoRequest{
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
					strconv.Itoa(len(field.UserOptions)),
				})
			}

			return writeOutput(
				cmd,
				rt.OutputFormat,
				[]string{"Name", "Label", "Type", "Options", "ColorOptions", "PureOptions", "UserOptions"},
				rows,
				fields,
			)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	return cmd
}
