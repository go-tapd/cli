package cmd

import (
	"github.com/go-tapd/cli/internal/app"
	"github.com/go-tapd/tapd"
	"github.com/spf13/cobra"
)

func newMeasureCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "measure",
		Short: "Work with TAPD measures",
	}

	cmd.AddCommand(newMeasureLifeTimesCmd(rt))
	return cmd
}

func newMeasureLifeTimesCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		entityID    int
		entityType  string
		created     string
		limit       int
		page        int
		fields      string
	)

	cmd := &cobra.Command{
		Use:   "life-times",
		Short: "List entity status life times",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.LifeTimesRequest{
				WorkspaceID: tapd.Ptr(workspaceID),
				EntityID:    tapd.Ptr(entityID),
				EntityType:  tapd.Ptr(tapd.EntityType(entityType)),
				Limit:       tapd.Ptr(limit),
				Page:        tapd.Ptr(page),
				Fields:      fieldsMulti(fields),
			}
			if created != "" {
				request.Created = tapd.Ptr(created)
			}

			items, _, err := client.MeasureService.LifeTimes(cmd.Context(), request)
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(items))
			for _, item := range items {
				rows = append(rows, []string{
					item.ID,
					string(item.EntityType),
					item.EntityID,
					item.Status,
					item.Owner,
					item.TimeCost,
					item.Created,
				})
			}

			return writeOutput(cmd, rt.OutputFormat, []string{"ID", "EntityType", "EntityID", "Status", "Owner", "TimeCost", "Created"}, rows, items)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	newListFlags(cmd, &limit, &page)
	cmd.Flags().IntVar(&entityID, "entity-id", 0, "entity ID")
	cmd.Flags().StringVar(&entityType, "entity-type", "", "entity type: story, task, or bug")
	cmd.Flags().StringVar(&created, "created", "", "filter by created time expression")
	cmd.Flags().StringVar(&fields, "fields", "", "comma separated fields")
	_ = cmd.MarkFlagRequired("entity-id")
	_ = cmd.MarkFlagRequired("entity-type")
	return cmd
}
