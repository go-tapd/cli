package cmd

import (
	"github.com/go-tapd/cli/internal/app"
	"github.com/go-tapd/tapd"
	"github.com/spf13/cobra"
)

func newStoryCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "story",
		Short: "Work with TAPD stories",
	}

	cmd.AddCommand(newStoryListCmd(rt))
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
