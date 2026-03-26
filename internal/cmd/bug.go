package cmd

import (
	"github.com/go-tapd/cli/internal/app"
	"github.com/go-tapd/tapd"
	"github.com/spf13/cobra"
)

func newBugCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bug",
		Short: "Work with TAPD bugs",
	}

	cmd.AddCommand(newBugListCmd(rt))
	return cmd
}

func newBugListCmd(rt *app.Runtime) *cobra.Command {
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
		Short: "List bugs",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetBugsRequest{
				WorkspaceID: tapd.Ptr(workspaceID),
				Limit:       tapd.Ptr(limit),
				Page:        tapd.Ptr(page),
				Fields:      fieldsMulti(fields),
				ID:          int64Multi(ids),
			}
			if creator != "" {
				request.Reporter = tapd.NewMulti(creator)
			}
			if owner != "" {
				request.CurrentOwner = tapd.Ptr(owner)
			}

			items, _, err := client.BugService.GetBugs(cmd.Context(), request)
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(items))
			for _, item := range items {
				rows = append(rows, []string{
					item.ID,
					item.Title,
					item.Status,
					item.CurrentOwner,
					item.Reporter,
					item.Modified,
				})
			}

			return writeOutput(cmd, rt.OutputFormat, []string{"ID", "Title", "Status", "Owner", "Reporter", "Modified"}, rows, items)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	newListFlags(cmd, &limit, &page)
	cmd.Flags().StringVar(&fields, "fields", "", "comma separated fields")
	cmd.Flags().StringVar(&ids, "ids", "", "comma separated bug IDs")
	cmd.Flags().StringVar(&creator, "creator", "", "filter by creator")
	cmd.Flags().StringVar(&owner, "owner", "", "filter by current owner")
	return cmd
}
