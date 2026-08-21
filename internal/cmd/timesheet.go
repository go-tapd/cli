package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/go-tapd/cli/internal/app"
	"github.com/go-tapd/tapd"
	"github.com/spf13/cobra"
)

func newTimesheetCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "timesheet",
		Short: "Work with TAPD timesheets",
	}

	cmd.AddCommand(
		newTimesheetCreateCmd(rt),
		newTimesheetListCmd(rt),
		newTimesheetCountCmd(rt),
		newTimesheetUpdateCmd(rt),
		newTimesheetDeleteCmd(rt),
	)
	return cmd
}

func newTimesheetCreateCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		entityType  string
		entityID    int64
		timespent   string
		timeremain  string
		spentdate   string
		owner       string
		memo        string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a timesheet",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.CreateTimesheetRequest{
				WorkspaceID: new(workspaceID),
				EntityType:  new(tapd.EntityType(entityType)),
				EntityID:    new(entityID),
				Timespent:   new(timespent),
				Owner:       new(owner),
			}
			if timeremain != "" {
				request.Timeremain = new(timeremain)
			}
			if spentdate != "" {
				request.Spentdate = new(spentdate)
			}
			if memo != "" {
				request.Memo = new(memo)
			}

			timesheet, _, err := client.TimesheetService.CreateTimesheet(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, timesheetHeaders(), timesheetRows([]*tapd.Timesheet{timesheet}), timesheet)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().StringVar(&entityType, "entity-type", "", "entity type: story, task, or bug")
	cmd.Flags().Int64Var(&entityID, "entity-id", 0, "entity ID")
	cmd.Flags().StringVar(&timespent, "timespent", "", "spent time")
	cmd.Flags().StringVar(&timeremain, "timeremain", "", "remaining time")
	cmd.Flags().StringVar(&spentdate, "spentdate", "", "spent date, for example 2026-06-01")
	cmd.Flags().StringVar(&owner, "owner", "", "timesheet owner")
	cmd.Flags().StringVar(&memo, "memo", "", "timesheet memo")
	_ = cmd.MarkFlagRequired("entity-type")
	_ = cmd.MarkFlagRequired("entity-id")
	_ = cmd.MarkFlagRequired("timespent")
	_ = cmd.MarkFlagRequired("owner")
	return cmd
}

func newTimesheetListCmd(rt *app.Runtime) *cobra.Command {
	var flags timesheetQueryFlags

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List timesheets",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request, err := newGetTimesheetsRequest(flags)
			if err != nil {
				return err
			}

			timesheets, _, err := client.TimesheetService.GetTimesheets(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, timesheetHeaders(), timesheetRows(timesheets), timesheets)
		},
	}

	addTimesheetQueryFlags(cmd, &flags, true)
	return cmd
}

func newTimesheetCountCmd(rt *app.Runtime) *cobra.Command {
	var flags timesheetQueryFlags

	cmd := &cobra.Command{
		Use:   "count",
		Short: "Count timesheets",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request, err := newGetTimesheetsCountRequest(flags)
			if err != nil {
				return err
			}

			count, _, err := client.TimesheetService.GetTimesheetsCount(cmd.Context(), request)
			if err != nil {
				return err
			}

			data := map[string]any{"resource": "timesheet", "workspace_id": flags.workspaceID, "count": count}
			rows := [][]string{{"timesheet", strconv.Itoa(flags.workspaceID), strconv.Itoa(count)}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Resource", "WorkspaceID", "Count"}, rows, data)
		},
	}

	addTimesheetQueryFlags(cmd, &flags, false)
	return cmd
}

func newTimesheetUpdateCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		timespent   string
		timeremain  string
		memo        string
	)

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a timesheet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseInt64Arg("timesheet id", args[0])
			if err != nil {
				return err
			}
			if timespent == "" && timeremain == "" && memo == "" {
				return errors.New("at least one of --timespent, --timeremain, or --memo is required")
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.UpdateTimesheetRequest{
				ID:          new(id),
				WorkspaceID: new(workspaceID),
			}
			if timespent != "" {
				request.Timespent = new(timespent)
			}
			if timeremain != "" {
				request.Timeremain = new(timeremain)
			}
			if memo != "" {
				request.Memo = new(memo)
			}

			timesheet, _, err := client.TimesheetService.UpdateTimesheet(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, timesheetHeaders(), timesheetRows([]*tapd.Timesheet{timesheet}), timesheet)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().StringVar(&timespent, "timespent", "", "spent time")
	cmd.Flags().StringVar(&timeremain, "timeremain", "", "remaining time")
	cmd.Flags().StringVar(&memo, "memo", "", "timesheet memo")
	return cmd
}

func newTimesheetDeleteCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		entityType  string
		entityID    int64
		costIDs     string
	)

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete timesheets",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ids, err := strictInt64Slice("cost IDs", costIDs)
			if err != nil {
				return err
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			result, _, err := client.TimesheetService.DeleteTimesheets(cmd.Context(), &tapd.DeleteTimesheetsRequest{
				WorkspaceID: new(workspaceID),
				EntityType:  new(tapd.EntityType(entityType)),
				EntityID:    new(entityID),
				CostIDs:     new(ids),
			})
			if err != nil {
				return err
			}

			rows := [][]string{{
				result.Msg,
				strconv.Itoa(len(result.Data.Success.CostIDs)),
				strconv.Itoa(len(result.Data.Failed)),
			}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Message", "Success", "Failed"}, rows, result)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().StringVar(&entityType, "entity-type", "", "entity type: story, task, or bug")
	cmd.Flags().Int64Var(&entityID, "entity-id", 0, "entity ID")
	cmd.Flags().StringVar(&costIDs, "cost-ids", "", "comma separated timesheet cost IDs")
	_ = cmd.MarkFlagRequired("entity-type")
	_ = cmd.MarkFlagRequired("entity-id")
	_ = cmd.MarkFlagRequired("cost-ids")
	return cmd
}

type timesheetQueryFlags struct {
	workspaceID        int
	limit              int
	page               int
	fields             string
	ids                string
	entityType         string
	entityID           int64
	timespent          string
	spentdate          string
	modified           string
	owner              string
	excludeParentStory bool
	created            string
	memo               string
	includeDeleted     bool
}

func addTimesheetQueryFlags(cmd *cobra.Command, flags *timesheetQueryFlags, withPaging bool) {
	newWorkspaceFlag(cmd, &flags.workspaceID)
	if withPaging {
		newListFlags(cmd, &flags.limit, &flags.page)
		cmd.Flags().StringVar(&flags.fields, "fields", "", "comma separated fields")
	}
	cmd.Flags().StringVar(&flags.ids, "ids", "", "comma separated timesheet IDs")
	cmd.Flags().StringVar(&flags.entityType, "entity-type", "", "entity type: story, task, or bug")
	cmd.Flags().Int64Var(&flags.entityID, "entity-id", 0, "entity ID")
	cmd.Flags().StringVar(&flags.timespent, "timespent", "", "filter by spent time")
	cmd.Flags().StringVar(&flags.spentdate, "spentdate", "", "filter by spent date expression")
	cmd.Flags().StringVar(&flags.modified, "modified", "", "filter by modified time expression")
	cmd.Flags().StringVar(&flags.owner, "owner", "", "filter by timesheet owner")
	cmd.Flags().BoolVar(&flags.excludeParentStory, "exclude-parent-story", false, "exclude parent story timesheets")
	cmd.Flags().StringVar(&flags.created, "created", "", "filter by created time expression")
	cmd.Flags().StringVar(&flags.memo, "memo", "", "filter by memo")
	cmd.Flags().BoolVar(&flags.includeDeleted, "include-deleted", false, "include deleted timesheets")
}

func newGetTimesheetsRequest(flags timesheetQueryFlags) (*tapd.GetTimesheetsRequest, error) {
	request := &tapd.GetTimesheetsRequest{
		WorkspaceID: new(flags.workspaceID),
		Limit:       new(flags.limit),
		Page:        new(flags.page),
		Fields:      fieldsMulti(flags.fields),
	}
	if err := applyTimesheetFilters(request, flags); err != nil {
		return nil, err
	}
	return request, nil
}

func newGetTimesheetsCountRequest(flags timesheetQueryFlags) (*tapd.GetTimesheetsCountRequest, error) {
	request := &tapd.GetTimesheetsCountRequest{WorkspaceID: new(flags.workspaceID)}
	if flags.ids != "" {
		ids, err := strictInt64Multi("timesheet IDs", flags.ids)
		if err != nil {
			return nil, err
		}
		request.ID = ids
	}
	if flags.entityType != "" {
		request.EntityType = new(tapd.EntityType(flags.entityType))
	}
	if flags.entityID > 0 {
		request.EntityID = new(flags.entityID)
	}
	if flags.timespent != "" {
		request.Timespent = new(flags.timespent)
	}
	if flags.spentdate != "" {
		request.Spentdate = new(flags.spentdate)
	}
	if flags.modified != "" {
		request.Modified = new(flags.modified)
	}
	if flags.owner != "" {
		request.Owner = new(flags.owner)
	}
	if flags.excludeParentStory {
		request.IncludeParentStoryTimesheet = new(0)
	}
	if flags.created != "" {
		request.Created = new(flags.created)
	}
	if flags.memo != "" {
		request.Memo = new(flags.memo)
	}
	if flags.includeDeleted {
		request.IsDelete = new(1)
	}
	return request, nil
}

func applyTimesheetFilters(request *tapd.GetTimesheetsRequest, flags timesheetQueryFlags) error {
	if flags.ids != "" {
		ids, err := strictInt64Multi("timesheet IDs", flags.ids)
		if err != nil {
			return err
		}
		request.ID = ids
	}
	if flags.entityType != "" {
		request.EntityType = new(tapd.EntityType(flags.entityType))
	}
	if flags.entityID > 0 {
		request.EntityID = new(flags.entityID)
	}
	if flags.timespent != "" {
		request.Timespent = new(flags.timespent)
	}
	if flags.spentdate != "" {
		request.Spentdate = new(flags.spentdate)
	}
	if flags.modified != "" {
		request.Modified = new(flags.modified)
	}
	if flags.owner != "" {
		request.Owner = new(flags.owner)
	}
	if flags.excludeParentStory {
		request.IncludeParentStoryTimesheet = new(0)
	}
	if flags.created != "" {
		request.Created = new(flags.created)
	}
	if flags.memo != "" {
		request.Memo = new(flags.memo)
	}
	if flags.includeDeleted {
		request.IsDelete = new(1)
	}
	return nil
}

func strictInt64Slice(name, csv string) ([]int64, error) {
	items := splitCSV(csv)
	if len(items) == 0 {
		return nil, fmt.Errorf("%s cannot be empty", name)
	}

	values := make([]int64, 0, len(items))
	for _, item := range items {
		v, err := parseInt64Arg(name, item)
		if err != nil {
			return nil, err
		}
		values = append(values, v)
	}
	return values, nil
}

func timesheetHeaders() []string {
	return []string{"ID", "EntityType", "EntityID", "Timespent", "Spentdate", "Owner", "Modified"}
}

func timesheetRows(timesheets []*tapd.Timesheet) [][]string {
	rows := make([][]string, 0, len(timesheets))
	for _, item := range timesheets {
		rows = append(rows, []string{
			item.ID,
			string(item.EntityType),
			item.EntityID,
			item.Timespent,
			item.Spentdate,
			item.Owner,
			item.Modified,
		})
	}
	return rows
}
