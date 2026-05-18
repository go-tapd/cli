package cmd

import (
	"github.com/go-tapd/cli/internal/app"
	"github.com/go-tapd/tapd"
	"github.com/spf13/cobra"
)

func newReportCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report",
		Short: "Work with TAPD reports",
	}

	cmd.AddCommand(newReportListCmd(rt))
	return cmd
}

func newReportListCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		limit       int
		page        int
		fields      string
		id          int
		title       string
		author      string
		created     string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List reports",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetReportsRequest{
				WorkspaceID: tapd.Ptr(workspaceID),
				Limit:       tapd.Ptr(limit),
				Page:        tapd.Ptr(page),
				Fields:      fieldsMulti(fields),
			}
			if id > 0 {
				request.ID = tapd.Ptr(id)
			}
			if title != "" {
				request.Title = tapd.Ptr(title)
			}
			if author != "" {
				request.Author = tapd.Ptr(author)
			}
			if created != "" {
				request.Created = tapd.Ptr(created)
			}

			reports, _, err := client.ReportService.GetReports(cmd.Context(), request)
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(reports))
			for _, item := range reports {
				rows = append(rows, []string{
					item.ID,
					item.Title,
					string(item.ReportType),
					string(item.Status),
					item.Author,
					item.Created,
				})
			}

			return writeOutput(cmd, rt.OutputFormat, []string{"ID", "Title", "Type", "Status", "Author", "Created"}, rows, reports)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	newListFlags(cmd, &limit, &page)
	cmd.Flags().StringVar(&fields, "fields", "", "comma separated fields")
	cmd.Flags().IntVar(&id, "id", 0, "report ID")
	cmd.Flags().StringVar(&title, "title", "", "filter by report title")
	cmd.Flags().StringVar(&author, "author", "", "filter by report author")
	cmd.Flags().StringVar(&created, "created", "", "filter by created time expression")
	return cmd
}
