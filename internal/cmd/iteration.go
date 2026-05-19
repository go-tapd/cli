package cmd

import (
	"fmt"
	"strconv"

	"github.com/go-tapd/cli/internal/app"
	"github.com/go-tapd/tapd"
	"github.com/spf13/cobra"
)

func newIterationCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "iteration",
		Short: "Work with TAPD iterations",
	}

	cmd.AddCommand(
		newIterationCreateCmd(rt),
		newIterationViewCmd(rt),
		newIterationListCmd(rt),
		newIterationCountCmd(rt),
		newIterationUpdateCmd(rt),
		newIterationChangesCmd(rt),
		newIterationFieldsCmd(rt),
		newIterationWorkitemTypesCmd(rt),
		newIterationTemplatesCmd(rt),
		newIterationTemplateFieldsCmd(rt),
		newIterationLockCmd(rt),
		newIterationUnlockCmd(rt),
	)
	return cmd
}

func newIterationCreateCmd(rt *app.Runtime) *cobra.Command {
	var flags iterationMutationFlags

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an iteration",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.CreateIterationRequest{
				WorkspaceID: tapd.Ptr(flags.workspaceID),
				Name:        tapd.Ptr(flags.name),
				Description: tapd.Ptr(flags.description),
				StartDate:   tapd.Ptr(flags.startDate),
				EndDate:     tapd.Ptr(flags.endDate),
				Creator:     tapd.Ptr(flags.creator),
			}
			applyIterationCreateFlags(request, flags)

			iteration, _, err := client.IterationService.CreateIteration(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, iterationTableHeaders(), iterationRows([]*tapd.Iteration{iteration}), iteration)
		},
	}

	newWorkspaceFlag(cmd, &flags.workspaceID)
	addIterationMutationFlags(cmd, &flags, false)
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("description")
	_ = cmd.MarkFlagRequired("start-date")
	_ = cmd.MarkFlagRequired("end-date")
	_ = cmd.MarkFlagRequired("creator")
	return cmd
}

func newIterationViewCmd(rt *app.Runtime) *cobra.Command {
	var workspaceID int

	cmd := &cobra.Command{
		Use:   "view <id>",
		Short: "Show iteration details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseInt64Arg("iteration id", args[0])
			if err != nil {
				return err
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			iterations, _, err := client.IterationService.GetIterations(cmd.Context(), &tapd.GetIterationsRequest{
				WorkspaceID: tapd.Ptr(workspaceID),
				ID:          tapd.NewMulti(id),
				Limit:       tapd.Ptr(1),
				Page:        tapd.Ptr(1),
			})
			if err != nil {
				return err
			}
			if len(iterations) == 0 {
				return fmt.Errorf("iteration %d not found", id)
			}

			return writeOutput(cmd, rt.OutputFormat, iterationTableHeaders(), iterationRows(iterations[:1]), iterations[0])
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	return cmd
}

func newIterationListCmd(rt *app.Runtime) *cobra.Command {
	var flags iterationQueryFlags

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List iterations",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request, err := newIterationsRequest(flags)
			if err != nil {
				return err
			}

			iterations, _, err := client.IterationService.GetIterations(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, iterationTableHeaders(), iterationRows(iterations), iterations)
		},
	}

	addIterationQueryFlags(cmd, &flags, true)
	return cmd
}

func newIterationCountCmd(rt *app.Runtime) *cobra.Command {
	var flags iterationQueryFlags

	cmd := &cobra.Command{
		Use:   "count",
		Short: "Count iterations",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request, err := newIterationsCountRequest(flags)
			if err != nil {
				return err
			}

			count, _, err := client.IterationService.GetIterationsCount(cmd.Context(), request)
			if err != nil {
				return err
			}

			data := map[string]any{"resource": "iteration", "workspace_id": flags.workspaceID, "count": count}
			rows := [][]string{{"iteration", strconv.Itoa(flags.workspaceID), strconv.Itoa(count)}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Resource", "WorkspaceID", "Count"}, rows, data)
		},
	}

	addIterationQueryFlags(cmd, &flags, false)
	return cmd
}

func newIterationUpdateCmd(rt *app.Runtime) *cobra.Command {
	var flags iterationMutationFlags

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an iteration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseInt64Arg("iteration id", args[0])
			if err != nil {
				return err
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.UpdateIterationRequest{
				ID:          tapd.Ptr(id),
				WorkspaceID: tapd.Ptr(flags.workspaceID),
				CurrentUser: tapd.Ptr(flags.currentUser),
			}
			applyIterationUpdateFlags(request, flags)

			iteration, _, err := client.IterationService.UpdateIteration(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, iterationTableHeaders(), iterationRows([]*tapd.Iteration{iteration}), iteration)
		},
	}

	newWorkspaceFlag(cmd, &flags.workspaceID)
	addIterationMutationFlags(cmd, &flags, true)
	_ = cmd.MarkFlagRequired("current-user")
	return cmd
}

func newIterationChangesCmd(rt *app.Runtime) *cobra.Command {
	var flags iterationChangesFlags

	cmd := &cobra.Command{
		Use:   "changes",
		Short: "List iteration changes",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request, err := newIterationChangesRequest(flags)
			if err != nil {
				return err
			}

			changes, _, err := client.IterationService.GetIterationChanges(cmd.Context(), request)
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(changes))
			for _, item := range changes {
				oldValue := ""
				if item.OldValue != nil {
					oldValue = *item.OldValue
				}
				newValue := ""
				if item.NewValue != nil {
					newValue = *item.NewValue
				}
				rows = append(rows, []string{
					item.ID,
					item.IterationID,
					item.Field,
					oldValue,
					newValue,
					item.Author,
					item.Created,
				})
			}

			return writeOutput(
				cmd,
				rt.OutputFormat,
				[]string{"ID", "IterationID", "Field", "OldValue", "NewValue", "Author", "Created"},
				rows,
				changes,
			)
		},
	}

	newWorkspaceFlag(cmd, &flags.workspaceID)
	newListFlags(cmd, &flags.limit, &flags.page)
	cmd.Flags().Int64Var(&flags.iterationID, "iteration-id", 0, "iteration ID")
	cmd.Flags().StringVar(&flags.ids, "ids", "", "comma separated change IDs")
	cmd.Flags().StringVar(&flags.author, "author", "", "filter by author")
	cmd.Flags().StringVar(&flags.field, "field", "", "filter by changed field")
	cmd.Flags().StringVar(&flags.oldValue, "old-value", "", "filter by old value")
	cmd.Flags().StringVar(&flags.newValue, "new-value", "", "filter by new value")
	cmd.Flags().StringVar(&flags.created, "created", "", "filter by created time expression")
	cmd.Flags().StringVar(&flags.fields, "fields", "", "comma separated fields")
	_ = cmd.MarkFlagRequired("iteration-id")
	return cmd
}

func newIterationFieldsCmd(rt *app.Runtime) *cobra.Command {
	var workspaceID int

	cmd := &cobra.Command{
		Use:   "fields",
		Short: "List iteration custom field settings",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			fields, _, err := client.IterationService.GetIterationCustomFieldsSettings(
				cmd.Context(),
				&tapd.GetIterationCustomFieldsSettingsRequest{WorkspaceID: tapd.Ptr(workspaceID)},
			)
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(fields))
			for _, item := range fields {
				rows = append(rows, []string{item.ID, item.CustomField, item.Name, item.Type, item.Enabled})
			}

			return writeOutput(cmd, rt.OutputFormat, []string{"ID", "Field", "Name", "Type", "Enabled"}, rows, fields)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	return cmd
}

func newIterationWorkitemTypesCmd(rt *app.Runtime) *cobra.Command {
	var workspaceID int

	cmd := &cobra.Command{
		Use:   "workitem-types",
		Short: "List iteration workitem types",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			items, _, err := client.IterationService.GetWorkitemTypes(
				cmd.Context(),
				&tapd.GetWorkitemTypesRequest{WorkspaceID: tapd.Ptr(workspaceID)},
			)
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(items))
			for _, item := range items {
				rows = append(rows, []string{item.ID, item.Name, item.EntityType, item.Creator, item.Modified})
			}

			return writeOutput(cmd, rt.OutputFormat, []string{"ID", "Name", "EntityType", "Creator", "Modified"}, rows, items)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	return cmd
}

func newIterationTemplatesCmd(rt *app.Runtime) *cobra.Command {
	var workspaceID int

	cmd := &cobra.Command{
		Use:   "templates",
		Short: "List iteration templates",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			items, _, err := client.IterationService.GetTemplateList(
				cmd.Context(),
				&tapd.GetTemplateListRequest{WorkspaceID: tapd.Ptr(workspaceID)},
			)
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(items))
			for _, item := range items {
				rows = append(rows, []string{item.ID, item.Name, item.Type, item.Creator, item.Modified})
			}

			return writeOutput(cmd, rt.OutputFormat, []string{"ID", "Name", "Type", "Creator", "Modified"}, rows, items)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	return cmd
}

func newIterationTemplateFieldsCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID    int
		templateID     int64
		workitemTypeID int64
	)

	cmd := &cobra.Command{
		Use:   "template-fields",
		Short: "List iteration template fields",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if (templateID == 0 && workitemTypeID == 0) || (templateID > 0 && workitemTypeID > 0) {
				return fmt.Errorf("provide exactly one of --template-id or --workitem-type-id")
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			var fields []*tapd.IterationTemplateField
			if templateID > 0 {
				fields, _, err = client.IterationService.GetIterationTemplateFields(
					cmd.Context(),
					&tapd.GetIterationTemplateFieldsRequest{
						WorkspaceID: tapd.Ptr(workspaceID),
						TemplateID:  tapd.Ptr(templateID),
					},
				)
			} else {
				fields, _, err = client.IterationService.GetIterationDefaultTemplateFields(
					cmd.Context(),
					&tapd.GetIterationDefaultTemplateFieldsRequest{
						WorkspaceID:    tapd.Ptr(workspaceID),
						WorkitemTypeID: tapd.Ptr(workitemTypeID),
					},
				)
			}
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(fields))
			for _, item := range fields {
				rows = append(rows, []string{item.ID, item.Field, item.Value, item.Required, item.Sort})
			}

			return writeOutput(cmd, rt.OutputFormat, []string{"ID", "Field", "Value", "Required", "Sort"}, rows, fields)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().Int64Var(&templateID, "template-id", 0, "iteration template ID")
	cmd.Flags().Int64Var(&workitemTypeID, "workitem-type-id", 0, "iteration workitem type ID")
	return cmd
}

func newIterationLockCmd(rt *app.Runtime) *cobra.Command {
	return newIterationLockToggleCmd(rt, true)
}

func newIterationUnlockCmd(rt *app.Runtime) *cobra.Command {
	return newIterationLockToggleCmd(rt, false)
}

func newIterationLockToggleCmd(rt *app.Runtime, lock bool) *cobra.Command {
	var (
		workspaceID int
		iterationID int64
		lockTypes   string
	)

	use := "lock"
	short := "Lock an iteration"
	if !lock {
		use = "unlock"
		short = "Unlock an iteration"
	}

	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, _ []string) error {
			types, err := strictStringMulti("lock types", lockTypes)
			if err != nil {
				return err
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			var result string
			if lock {
				result, _, err = client.IterationService.LockIteration(cmd.Context(), &tapd.LockIterationRequest{
					WorkspaceID: tapd.Ptr(workspaceID),
					IterationID: tapd.Ptr(iterationID),
					LockTypes:   types,
				})
			} else {
				result, _, err = client.IterationService.UnlockIteration(cmd.Context(), &tapd.UnlockIterationRequest{
					WorkspaceID: tapd.Ptr(workspaceID),
					IterationID: tapd.Ptr(iterationID),
					LockTypes:   types,
				})
			}
			if err != nil {
				return err
			}

			rows := [][]string{{result}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Result"}, rows, map[string]string{"result": result})
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().Int64Var(&iterationID, "iteration-id", 0, "iteration ID")
	cmd.Flags().StringVar(&lockTypes, "lock-types", "", "comma separated lock types")
	_ = cmd.MarkFlagRequired("iteration-id")
	_ = cmd.MarkFlagRequired("lock-types")
	return cmd
}

type iterationMutationFlags struct {
	workspaceID    int
	name           string
	description    string
	startDate      string
	endDate        string
	creator        string
	currentUser    string
	status         string
	label          string
	workitemTypeID int
	planAppID      int
}

func addIterationMutationFlags(cmd *cobra.Command, flags *iterationMutationFlags, update bool) {
	cmd.Flags().StringVar(&flags.name, "name", "", "iteration name")
	cmd.Flags().StringVar(&flags.description, "description", "", "iteration description")
	cmd.Flags().StringVar(&flags.startDate, "start-date", "", "iteration start date")
	cmd.Flags().StringVar(&flags.endDate, "end-date", "", "iteration end date")
	cmd.Flags().StringVar(&flags.creator, "creator", "", "iteration creator")
	cmd.Flags().StringVar(&flags.status, "status", "", "iteration status")
	cmd.Flags().StringVar(&flags.label, "label", "", "iteration labels, comma separated")
	if update {
		cmd.Flags().StringVar(&flags.currentUser, "current-user", "", "current operator")
	} else {
		cmd.Flags().IntVar(&flags.workitemTypeID, "workitem-type-id", 0, "iteration workitem type ID")
		cmd.Flags().IntVar(&flags.planAppID, "plan-app-id", 0, "plan app ID")
	}
}

func applyIterationCreateFlags(request *tapd.CreateIterationRequest, flags iterationMutationFlags) {
	if flags.workitemTypeID > 0 {
		request.WorkitemTypeID = tapd.Ptr(flags.workitemTypeID)
	}
	if flags.planAppID > 0 {
		request.PlanAppID = tapd.Ptr(flags.planAppID)
	}
	if flags.status != "" {
		request.Status = tapd.Ptr(flags.status)
	}
	if flags.label != "" {
		request.Label = stringEnum(flags.label)
	}
}

func applyIterationUpdateFlags(request *tapd.UpdateIterationRequest, flags iterationMutationFlags) {
	if flags.name != "" {
		request.Name = tapd.Ptr(flags.name)
	}
	if flags.description != "" {
		request.Description = tapd.Ptr(flags.description)
	}
	if flags.startDate != "" {
		request.StartDate = tapd.Ptr(flags.startDate)
	}
	if flags.endDate != "" {
		request.EndDate = tapd.Ptr(flags.endDate)
	}
	if flags.creator != "" {
		request.Creator = tapd.Ptr(flags.creator)
	}
	if flags.status != "" {
		request.Status = tapd.Ptr(flags.status)
	}
	if flags.label != "" {
		request.Label = stringEnum(flags.label)
	}
}

type iterationQueryFlags struct {
	workspaceID    int
	limit          int
	page           int
	fields         string
	ids            string
	name           string
	description    string
	startDate      string
	endDate        string
	workitemTypeID int
	planAppID      int
	status         string
	creator        string
	created        string
	modified       string
	completed      string
}

func addIterationQueryFlags(cmd *cobra.Command, flags *iterationQueryFlags, withPaging bool) {
	newWorkspaceFlag(cmd, &flags.workspaceID)
	if withPaging {
		newListFlags(cmd, &flags.limit, &flags.page)
		cmd.Flags().StringVar(&flags.fields, "fields", "", "comma separated fields")
	}
	cmd.Flags().StringVar(&flags.ids, "ids", "", "comma separated iteration IDs")
	cmd.Flags().StringVar(&flags.name, "name", "", "filter by name")
	cmd.Flags().StringVar(&flags.description, "description", "", "filter by description")
	cmd.Flags().StringVar(&flags.startDate, "start-date", "", "filter by start date expression")
	cmd.Flags().StringVar(&flags.endDate, "end-date", "", "filter by end date expression")
	cmd.Flags().IntVar(&flags.workitemTypeID, "workitem-type-id", 0, "filter by workitem type ID")
	cmd.Flags().IntVar(&flags.planAppID, "plan-app-id", 0, "filter by plan app ID")
	cmd.Flags().StringVar(&flags.status, "status", "", "filter by status")
	cmd.Flags().StringVar(&flags.creator, "creator", "", "filter by creator")
	cmd.Flags().StringVar(&flags.created, "created", "", "filter by created time expression")
	cmd.Flags().StringVar(&flags.modified, "modified", "", "filter by modified time expression")
	cmd.Flags().StringVar(&flags.completed, "completed", "", "filter by completed time expression")
}

func newIterationsRequest(flags iterationQueryFlags) (*tapd.GetIterationsRequest, error) {
	request := &tapd.GetIterationsRequest{
		WorkspaceID: tapd.Ptr(flags.workspaceID),
		Limit:       tapd.Ptr(flags.limit),
		Page:        tapd.Ptr(flags.page),
		Fields:      fieldsMulti(flags.fields),
	}
	if flags.ids != "" {
		ids, err := strictInt64Multi("iteration IDs", flags.ids)
		if err != nil {
			return nil, err
		}
		request.ID = ids
	}
	if flags.name != "" {
		request.Name = tapd.Ptr(flags.name)
	}
	if flags.description != "" {
		request.Description = tapd.Ptr(flags.description)
	}
	if flags.startDate != "" {
		request.StartDate = tapd.Ptr(flags.startDate)
	}
	if flags.endDate != "" {
		request.EndDate = tapd.Ptr(flags.endDate)
	}
	if flags.workitemTypeID > 0 {
		request.WorkitemTypeID = tapd.Ptr(flags.workitemTypeID)
	}
	if flags.planAppID > 0 {
		request.PlanAppID = tapd.Ptr(flags.planAppID)
	}
	if flags.status != "" {
		request.Status = tapd.Ptr(flags.status)
	}
	if flags.creator != "" {
		request.Creator = tapd.Ptr(flags.creator)
	}
	if flags.created != "" {
		request.Created = tapd.Ptr(flags.created)
	}
	if flags.modified != "" {
		request.Modified = tapd.Ptr(flags.modified)
	}
	if flags.completed != "" {
		request.Completed = tapd.Ptr(flags.completed)
	}
	return request, nil
}

func newIterationsCountRequest(flags iterationQueryFlags) (*tapd.GetIterationsCountRequest, error) {
	request := &tapd.GetIterationsCountRequest{
		WorkspaceID: tapd.Ptr(flags.workspaceID),
	}
	if flags.ids != "" {
		ids, err := strictInt64Multi("iteration IDs", flags.ids)
		if err != nil {
			return nil, err
		}
		request.ID = ids
	}
	if flags.name != "" {
		request.Name = tapd.Ptr(flags.name)
	}
	if flags.description != "" {
		request.Description = tapd.Ptr(flags.description)
	}
	if flags.startDate != "" {
		request.StartDate = tapd.Ptr(flags.startDate)
	}
	if flags.endDate != "" {
		request.EndDate = tapd.Ptr(flags.endDate)
	}
	if flags.workitemTypeID > 0 {
		request.WorkitemTypeID = tapd.Ptr(flags.workitemTypeID)
	}
	if flags.planAppID > 0 {
		request.PlanAppID = tapd.Ptr(flags.planAppID)
	}
	if flags.status != "" {
		request.Status = tapd.Ptr(flags.status)
	}
	if flags.creator != "" {
		request.Creator = tapd.Ptr(flags.creator)
	}
	if flags.created != "" {
		request.Created = tapd.Ptr(flags.created)
	}
	if flags.modified != "" {
		request.Modified = tapd.Ptr(flags.modified)
	}
	if flags.completed != "" {
		request.Completed = tapd.Ptr(flags.completed)
	}
	return request, nil
}

type iterationChangesFlags struct {
	workspaceID int
	iterationID int64
	limit       int
	page        int
	fields      string
	ids         string
	author      string
	field       string
	oldValue    string
	newValue    string
	created     string
}

func newIterationChangesRequest(flags iterationChangesFlags) (*tapd.GetIterationChangesRequest, error) {
	request := &tapd.GetIterationChangesRequest{
		WorkspaceID: tapd.Ptr(flags.workspaceID),
		IterationID: tapd.Ptr(flags.iterationID),
		Limit:       tapd.Ptr(flags.limit),
		Page:        tapd.Ptr(flags.page),
		Fields:      fieldsMulti(flags.fields),
	}
	if flags.ids != "" {
		ids, err := strictInt64Multi("change IDs", flags.ids)
		if err != nil {
			return nil, err
		}
		request.ID = ids
	}
	if flags.author != "" {
		request.Author = tapd.Ptr(flags.author)
	}
	if flags.field != "" {
		request.Field = tapd.Ptr(flags.field)
	}
	if flags.oldValue != "" {
		request.OldValue = tapd.Ptr(flags.oldValue)
	}
	if flags.newValue != "" {
		request.NewValue = tapd.Ptr(flags.newValue)
	}
	if flags.created != "" {
		request.Created = tapd.Ptr(flags.created)
	}
	return request, nil
}

func iterationTableHeaders() []string {
	return []string{"ID", "Name", "Status", "StartDate", "EndDate", "Creator", "Modified"}
}

func iterationRows(iterations []*tapd.Iteration) [][]string {
	rows := make([][]string, 0, len(iterations))
	for _, item := range iterations {
		rows = append(rows, []string{
			item.ID,
			item.Name,
			item.Status,
			item.StartDate,
			item.EndDate,
			item.Creator,
			item.Modified,
		})
	}
	return rows
}
