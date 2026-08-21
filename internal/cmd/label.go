package cmd

import (
	"fmt"
	"strconv"

	"github.com/go-tapd/cli/internal/app"
	"github.com/go-tapd/tapd"
	"github.com/spf13/cobra"
)

func newLabelCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "label",
		Short: "Work with TAPD labels",
	}

	cmd.AddCommand(
		newLabelCreateCmd(rt),
		newLabelListCmd(rt),
		newLabelCountCmd(rt),
		newLabelUpdateCmd(rt),
	)
	return cmd
}

func newLabelCreateCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		name        string
		color       string
		creator     string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a label",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.CreateLabelRequest{
				WorkspaceID: new(workspaceID),
				Name:        new(name),
			}
			if color != "" {
				request.Color = new(tapd.LabelColor(color))
			}
			if creator != "" {
				request.Creator = new(creator)
			}

			label, _, err := client.LabelService.CreateLabel(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, labelTableHeaders(), labelRows([]*tapd.Label{label}), label)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().StringVar(&name, "name", "", "label name")
	cmd.Flags().StringVar(&color, "color", "", "label color: 1, 2, 3, or 4")
	cmd.Flags().StringVar(&creator, "creator", "", "label creator")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newLabelListCmd(rt *app.Runtime) *cobra.Command {
	var flags labelQueryFlags

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List labels",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request, err := newLabelsRequest(flags)
			if err != nil {
				return err
			}

			labels, _, err := client.LabelService.GetLabels(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, labelTableHeaders(), labelRows(labels), labels)
		},
	}

	addLabelQueryFlags(cmd, &flags, true)
	return cmd
}

func newLabelCountCmd(rt *app.Runtime) *cobra.Command {
	var flags labelQueryFlags

	cmd := &cobra.Command{
		Use:   "count",
		Short: "Count labels",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request, err := newLabelCountRequest(flags)
			if err != nil {
				return err
			}

			count, _, err := client.LabelService.GetLabelsCount(cmd.Context(), request)
			if err != nil {
				return err
			}

			data := map[string]any{"resource": "label", "workspace_id": flags.workspaceID, "count": count}
			rows := [][]string{{"label", strconv.Itoa(flags.workspaceID), strconv.Itoa(count)}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Resource", "WorkspaceID", "Count"}, rows, data)
		},
	}

	addLabelQueryFlags(cmd, &flags, false)
	return cmd
}

func newLabelUpdateCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		color       string
		modifier    string
	)

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a label",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseIntArg("label id", args[0])
			if err != nil {
				return err
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.UpdateLabelRequest{
				ID:          new(id),
				WorkspaceID: new(workspaceID),
			}
			if color != "" {
				request.Color = new(tapd.LabelColor(color))
			}
			if modifier != "" {
				request.Modifier = new(modifier)
			}
			if request.Color == nil && request.Modifier == nil {
				return fmt.Errorf("provide at least one of --color or --modifier")
			}

			label, _, err := client.LabelService.UpdateLabel(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, labelTableHeaders(), labelRows([]*tapd.Label{label}), label)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().StringVar(&color, "color", "", "label color: 1, 2, 3, or 4")
	cmd.Flags().StringVar(&modifier, "modifier", "", "label modifier")
	return cmd
}

type labelQueryFlags struct {
	workspaceID int
	limit       int
	page        int
	ids         string
	name        string
	creator     string
	created     string
}

func addLabelQueryFlags(cmd *cobra.Command, flags *labelQueryFlags, withPaging bool) {
	newWorkspaceFlag(cmd, &flags.workspaceID)
	if withPaging {
		newListFlags(cmd, &flags.limit, &flags.page)
	}
	cmd.Flags().StringVar(&flags.ids, "ids", "", "comma separated label IDs")
	cmd.Flags().StringVar(&flags.name, "name", "", "filter by label name")
	cmd.Flags().StringVar(&flags.creator, "creator", "", "filter by creator")
	cmd.Flags().StringVar(&flags.created, "created", "", "filter by created time expression")
}

func newLabelsRequest(flags labelQueryFlags) (*tapd.GetLabelsRequest, error) {
	request := &tapd.GetLabelsRequest{
		WorkspaceID: new(flags.workspaceID),
		Limit:       new(flags.limit),
		Page:        new(flags.page),
	}
	if flags.ids != "" {
		ids, err := strictIntMulti("label IDs", flags.ids)
		if err != nil {
			return nil, err
		}
		request.ID = ids
	}
	if flags.name != "" {
		request.Name = new(flags.name)
	}
	if flags.creator != "" {
		request.Creator = new(flags.creator)
	}
	if flags.created != "" {
		request.Created = new(flags.created)
	}
	return request, nil
}

func newLabelCountRequest(flags labelQueryFlags) (*tapd.GetLabelCountRequest, error) {
	request := &tapd.GetLabelCountRequest{
		WorkspaceID: new(flags.workspaceID),
	}
	if flags.ids != "" {
		ids, err := strictIntMulti("label IDs", flags.ids)
		if err != nil {
			return nil, err
		}
		request.ID = ids
	}
	if flags.name != "" {
		request.Name = new(flags.name)
	}
	if flags.creator != "" {
		request.Creator = new(flags.creator)
	}
	if flags.created != "" {
		request.Created = new(flags.created)
	}
	return request, nil
}

func labelTableHeaders() []string {
	return []string{"ID", "Name", "Color", "Creator", "Modifier", "Modified"}
}

func labelRows(labels []*tapd.Label) [][]string {
	rows := make([][]string, 0, len(labels))
	for _, item := range labels {
		rows = append(rows, []string{
			item.ID,
			item.Name,
			string(item.Color),
			item.Creator,
			item.Modifier,
			item.Modified,
		})
	}
	return rows
}
