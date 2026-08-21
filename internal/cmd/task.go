package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

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
		newTaskCreateCmd(rt),
		newTaskViewCmd(rt),
		newTaskListCmd(rt),
		newTaskCountCmd(rt),
		newTaskUpdateCmd(rt),
		newTaskBatchUpdateCmd(rt),
		newTaskChangesCmd(rt),
		newTaskFieldsCmd(rt),
		newTaskRemovedCmd(rt),
	)
	return cmd
}

func newTaskCreateCmd(rt *app.Runtime) *cobra.Command {
	var flags taskMutationFlags

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a task",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.CreateTaskRequest{
				WorkspaceID: new(flags.workspaceID),
				Name:        new(flags.name),
			}
			if err := applyTaskCreateFlags(request, flags); err != nil {
				return err
			}

			task, _, err := client.TaskService.CreateTask(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, taskTableHeaders(), taskRows([]*tapd.Task{task}), task)
		},
	}

	newWorkspaceFlag(cmd, &flags.workspaceID)
	cmd.Flags().StringVar(&flags.name, "name", "", "task name")
	_ = cmd.MarkFlagRequired("name")
	addTaskMutationFlags(cmd, &flags, false)
	return cmd
}

func newTaskViewCmd(rt *app.Runtime) *cobra.Command {
	var workspaceID int

	cmd := &cobra.Command{
		Use:   "view <id>",
		Short: "Show task details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseInt64Arg("task id", args[0])
			if err != nil {
				return err
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			tasks, _, err := client.TaskService.GetTasks(cmd.Context(), &tapd.GetTasksRequest{
				WorkspaceID: new(workspaceID),
				ID:          tapd.NewMulti(id),
				Limit:       new(1),
				Page:        new(1),
			})
			if err != nil {
				return err
			}
			if len(tasks) == 0 {
				return fmt.Errorf("task %d not found", id)
			}

			return writeOutput(cmd, rt.OutputFormat, taskTableHeaders(), taskRows(tasks[:1]), tasks[0])
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	return cmd
}

func newTaskListCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		limit       int
		page        int
		fields      string
		ids         string
		name        string
		creator     string
		owner       string
		status      string
		storyIDs    string
		iterationID int64
		label       string
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
				WorkspaceID: new(workspaceID),
				Limit:       new(limit),
				Page:        new(page),
				Fields:      fieldsMulti(fields),
				ID:          int64Multi(ids),
			}
			if name != "" {
				request.Name = new(name)
			}
			if creator != "" {
				request.Creator = new(creator)
			}
			if owner != "" {
				request.Owner = new(owner)
			}
			if status != "" {
				request.Status = taskStatusEnum(status)
			}
			if storyIDs != "" {
				ids, err := strictInt64Multi("story IDs", storyIDs)
				if err != nil {
					return err
				}
				request.StoryID = ids
			}
			if iterationID > 0 {
				request.IterationID = tapd.NewEnum(iterationID)
			}
			if label != "" {
				request.Label = stringEnum(label)
			}

			items, _, err := client.TaskService.GetTasks(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, taskTableHeaders(), taskRows(items), items)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	newListFlags(cmd, &limit, &page)
	cmd.Flags().StringVar(&fields, "fields", "", "comma separated fields")
	cmd.Flags().StringVar(&ids, "ids", "", "comma separated task IDs")
	cmd.Flags().StringVar(&name, "name", "", "filter by task name")
	cmd.Flags().StringVar(&creator, "creator", "", "filter by creator")
	cmd.Flags().StringVar(&owner, "owner", "", "filter by owner")
	cmd.Flags().StringVar(&status, "status", "", "filter by status: open, progressing, or done")
	cmd.Flags().StringVar(&storyIDs, "story-ids", "", "comma separated related story IDs")
	cmd.Flags().Int64Var(&iterationID, "iteration-id", 0, "filter by iteration ID")
	cmd.Flags().StringVar(&label, "label", "", "filter by labels, comma separated")
	return cmd
}

func newTaskCountCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		ids         string
		name        string
		creator     string
		owner       string
		status      string
		storyIDs    string
		iterationID int64
		label       string
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
				WorkspaceID: new(workspaceID),
				ID:          int64Multi(ids),
			}
			if name != "" {
				request.Name = new(name)
			}
			if creator != "" {
				request.Creator = new(creator)
			}
			if owner != "" {
				request.Owner = new(owner)
			}
			if status != "" {
				request.Status = taskStatusEnum(status)
			}
			if storyIDs != "" {
				ids, err := strictInt64Multi("story IDs", storyIDs)
				if err != nil {
					return err
				}
				request.StoryID = ids
			}
			if iterationID > 0 {
				request.IterationID = tapd.NewEnum(iterationID)
			}
			if label != "" {
				request.Label = stringEnum(label)
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
	cmd.Flags().StringVar(&name, "name", "", "filter by task name")
	cmd.Flags().StringVar(&creator, "creator", "", "filter by creator")
	cmd.Flags().StringVar(&owner, "owner", "", "filter by owner")
	cmd.Flags().StringVar(&status, "status", "", "filter by status: open, progressing, or done")
	cmd.Flags().StringVar(&storyIDs, "story-ids", "", "comma separated related story IDs")
	cmd.Flags().Int64Var(&iterationID, "iteration-id", 0, "filter by iteration ID")
	cmd.Flags().StringVar(&label, "label", "", "filter by labels, comma separated")
	return cmd
}

func newTaskUpdateCmd(rt *app.Runtime) *cobra.Command {
	var flags taskMutationFlags

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseInt64Arg("task id", args[0])
			if err != nil {
				return err
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.UpdateTaskRequest{
				ID:          new(id),
				WorkspaceID: new(flags.workspaceID),
			}
			if err := applyTaskUpdateFlags(request, flags); err != nil {
				return err
			}

			task, _, err := client.TaskService.UpdateTask(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, taskTableHeaders(), taskRows([]*tapd.Task{task}), task)
		},
	}

	newWorkspaceFlag(cmd, &flags.workspaceID)
	addTaskMutationFlags(cmd, &flags, true)
	return cmd
}

func newTaskBatchUpdateCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		file        string
	)

	cmd := &cobra.Command{
		Use:   "batch-update",
		Short: "Batch update tasks from a JSON file",
		RunE: func(cmd *cobra.Command, _ []string) error {
			data, err := os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("read batch update file: %w", err)
			}

			request, err := decodeTaskBatchUpdate(workspaceID, data)
			if err != nil {
				return err
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			result, _, err := client.TaskService.BatchUpdateTasks(cmd.Context(), request)
			if err != nil {
				return err
			}

			rows := [][]string{{result.Msg}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Message"}, rows, result)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().StringVar(&file, "file", "", "JSON file containing an array of task update objects")
	_ = cmd.MarkFlagRequired("file")
	return cmd
}

func newTaskChangesCmd(rt *app.Runtime) *cobra.Command {
	var flags taskChangesFlags

	cmd := &cobra.Command{
		Use:   "changes",
		Short: "List task changes",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request, err := newTaskChangesRequest(flags)
			if err != nil {
				return err
			}

			changes, _, err := client.TaskService.GetTaskChanges(cmd.Context(), request)
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(changes))
			for _, item := range changes {
				rows = append(rows, []string{
					item.ID,
					item.TaskID,
					item.ChangeType,
					item.ChangeSummary,
					item.Creator,
					item.Created,
					strconv.Itoa(len(item.FieldChanges)),
				})
			}

			return writeOutput(
				cmd,
				rt.OutputFormat,
				[]string{"ID", "TaskID", "Type", "Summary", "Creator", "Created", "FieldChanges"},
				rows,
				changes,
			)
		},
	}

	addTaskChangesFlags(cmd, &flags, true)
	cmd.AddCommand(newTaskChangesCountCmd(rt))
	return cmd
}

func newTaskChangesCountCmd(rt *app.Runtime) *cobra.Command {
	var flags taskChangesFlags

	cmd := &cobra.Command{
		Use:   "count",
		Short: "Count task changes",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request, err := newTaskChangesCountRequest(flags)
			if err != nil {
				return err
			}

			count, _, err := client.TaskService.GetTaskChangesCount(cmd.Context(), request)
			if err != nil {
				return err
			}

			data := map[string]any{"resource": "task_change", "workspace_id": flags.workspaceID, "count": count}
			rows := [][]string{{"task_change", strconv.Itoa(flags.workspaceID), strconv.Itoa(count)}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Resource", "WorkspaceID", "Count"}, rows, data)
		},
	}

	addTaskChangesFlags(cmd, &flags, false)
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
				WorkspaceID: new(workspaceID),
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

func newTaskRemovedCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		limit       int
		page        int
		ids         string
		creator     string
		created     string
		deleted     string
		archived    bool
	)

	cmd := &cobra.Command{
		Use:   "removed",
		Short: "List removed tasks",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetRemovedTasksRequest{
				WorkspaceID: new(workspaceID),
				Limit:       new(limit),
				Page:        new(page),
			}
			if ids != "" {
				request.ID, err = strictIntMulti("task IDs", ids)
				if err != nil {
					return err
				}
			}
			if creator != "" {
				request.Creator = new(creator)
			}
			if created != "" {
				request.Created = new(created)
			}
			if deleted != "" {
				request.Deleted = new(deleted)
			}
			if archived {
				request.IsArchived = new(1)
			}

			tasks, _, err := client.TaskService.GetRemovedTasks(cmd.Context(), request)
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(tasks))
			for _, item := range tasks {
				rows = append(rows, []string{item.ID, item.Name, item.Creator, item.OperationUser, item.IsArchived, item.Deleted})
			}

			return writeOutput(cmd, rt.OutputFormat, []string{"ID", "Name", "Creator", "OperationUser", "Archived", "Deleted"}, rows, tasks)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	newListFlags(cmd, &limit, &page)
	cmd.Flags().StringVar(&ids, "ids", "", "comma separated task IDs")
	cmd.Flags().StringVar(&creator, "creator", "", "filter by creator")
	cmd.Flags().StringVar(&created, "created", "", "filter by created time expression")
	cmd.Flags().StringVar(&deleted, "deleted", "", "filter by deleted time expression")
	cmd.Flags().BoolVar(&archived, "archived", false, "only return archived tasks")
	return cmd
}

type taskMutationFlags struct {
	workspaceID        int
	name               string
	description        string
	currentUser        string
	creator            string
	status             string
	label              string
	owner              string
	cc                 string
	begin              string
	due                string
	storyIDs           string
	storyID            int64
	iterationID        int64
	priority           string
	priorityLabel      string
	progress           int
	completed          string
	effortCompleted    string
	exceed             float64
	remain             float64
	effort             string
	autoCompleteEffort bool
	customFields       []string
}

func addTaskMutationFlags(cmd *cobra.Command, flags *taskMutationFlags, update bool) {
	if update {
		cmd.Flags().StringVar(&flags.name, "name", "", "task name")
		cmd.Flags().StringVar(&flags.currentUser, "current-user", "", "current operator")
		cmd.Flags().Int64Var(&flags.storyID, "story-id", 0, "related story ID")
		cmd.Flags().BoolVar(&flags.autoCompleteEffort, "auto-complete-effort", false, "auto fill completed effort when completing task")
	} else {
		cmd.Flags().StringVar(&flags.storyIDs, "story-ids", "", "comma separated related story IDs")
		cmd.Flags().StringVar(&flags.completed, "completed", "", "completed time")
		cmd.Flags().StringVar(&flags.effortCompleted, "effort-completed", "", "completed effort")
		cmd.Flags().Float64Var(&flags.exceed, "exceed", 0, "exceeded effort")
		cmd.Flags().Float64Var(&flags.remain, "remain", 0, "remaining effort")
	}
	cmd.Flags().StringVar(&flags.description, "description", "", "task description")
	cmd.Flags().StringVar(&flags.creator, "creator", "", "task creator")
	cmd.Flags().StringVar(&flags.status, "status", "", "task status: open, progressing, or done")
	cmd.Flags().StringVar(&flags.label, "label", "", "task labels, comma separated")
	cmd.Flags().StringVar(&flags.owner, "owner", "", "task owner")
	cmd.Flags().StringVar(&flags.cc, "cc", "", "task CC users")
	cmd.Flags().StringVar(&flags.begin, "begin", "", "expected begin date")
	cmd.Flags().StringVar(&flags.due, "due", "", "expected due date")
	cmd.Flags().Int64Var(&flags.iterationID, "iteration-id", 0, "iteration ID")
	cmd.Flags().StringVar(&flags.priority, "priority", "", "task priority")
	cmd.Flags().StringVar(&flags.priorityLabel, "priority-label", "", "task priority label")
	cmd.Flags().IntVar(&flags.progress, "progress", 0, "task progress")
	cmd.Flags().StringVar(&flags.effort, "effort", "", "estimated effort")
	cmd.Flags().StringArrayVar(&flags.customFields, "field", nil, "set custom field as key=value; repeatable, for example custom_field_one=开发阶段")
}

func applyTaskCreateFlags(request *tapd.CreateTaskRequest, flags taskMutationFlags) error {
	if flags.description != "" {
		request.Description = new(flags.description)
	}
	if flags.creator != "" {
		request.Creator = new(flags.creator)
	}
	if flags.status != "" {
		request.Status = taskStatusEnum(flags.status)
	}
	if flags.label != "" {
		request.Label = stringEnum(flags.label)
	}
	if flags.owner != "" {
		request.Owner = new(flags.owner)
	}
	if flags.cc != "" {
		request.CC = new(flags.cc)
	}
	if flags.begin != "" {
		request.Begin = new(flags.begin)
	}
	if flags.due != "" {
		request.Due = new(flags.due)
	}
	if flags.storyIDs != "" {
		ids, err := strictInt64Multi("story IDs", flags.storyIDs)
		if err != nil {
			return err
		}
		request.StoryID = ids
	}
	if flags.iterationID > 0 {
		request.IterationID = tapd.NewEnum(flags.iterationID)
	}
	if flags.priority != "" {
		request.Priority = new(flags.priority)
	}
	if flags.priorityLabel != "" {
		request.PriorityLabel = new(tapd.PriorityLabel(flags.priorityLabel))
	}
	if flags.progress > 0 {
		request.Progress = new(flags.progress)
	}
	if flags.completed != "" {
		request.Completed = new(flags.completed)
	}
	if flags.effortCompleted != "" {
		request.EffortCompleted = new(flags.effortCompleted)
	}
	if flags.exceed > 0 {
		request.Exceed = new(flags.exceed)
	}
	if flags.remain > 0 {
		request.Remain = new(flags.remain)
	}
	if flags.effort != "" {
		request.Effort = new(flags.effort)
	}
	if err := applyTaskCustomFields(request, flags.customFields); err != nil {
		return err
	}
	return nil
}

func applyTaskUpdateFlags(request *tapd.UpdateTaskRequest, flags taskMutationFlags) error {
	if flags.currentUser != "" {
		request.CurrentUser = new(flags.currentUser)
	}
	if flags.name != "" {
		request.Name = new(flags.name)
	}
	if flags.description != "" {
		request.Description = new(flags.description)
	}
	if flags.creator != "" {
		request.Creator = new(flags.creator)
	}
	if flags.status != "" {
		status := tapd.TaskStatus(flags.status)
		request.Status = new(status)
	}
	if flags.label != "" {
		request.Label = stringEnum(flags.label)
	}
	if flags.owner != "" {
		request.Owner = new(flags.owner)
	}
	if flags.cc != "" {
		request.CC = new(flags.cc)
	}
	if flags.begin != "" {
		request.Begin = new(flags.begin)
	}
	if flags.due != "" {
		request.Due = new(flags.due)
	}
	if flags.storyID > 0 {
		request.StoryID = new(flags.storyID)
	}
	if flags.iterationID > 0 {
		request.IterationID = new(flags.iterationID)
	}
	if flags.priority != "" {
		request.Priority = new(flags.priority)
	}
	if flags.priorityLabel != "" {
		request.PriorityLabel = new(tapd.PriorityLabel(flags.priorityLabel))
	}
	if flags.progress > 0 {
		request.Progress = new(flags.progress)
	}
	if flags.effort != "" {
		request.Effort = new(flags.effort)
	}
	if flags.autoCompleteEffort {
		request.AutoCompleteEffort = new(1)
	}
	if err := applyTaskCustomFields(request, flags.customFields); err != nil {
		return err
	}
	return nil
}

func applyTaskCustomFields(request any, fields []string) error {
	for _, field := range fields {
		name, value, ok := strings.Cut(field, "=")
		name = strings.TrimSpace(name)
		value = strings.TrimSpace(value)
		if !ok || name == "" {
			return fmt.Errorf("invalid custom field %q, expected key=value", field)
		}
		if err := setTaskCustomField(request, name, value); err != nil {
			return err
		}
	}
	return nil
}

func setTaskCustomField(request any, name, value string) error {
	v := reflect.ValueOf(request)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return errors.New("task custom fields require a non-nil request pointer")
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return errors.New("task custom fields require a struct request")
	}

	t := v.Type()
	for i := range v.NumField() {
		fieldType := t.Field(i)
		jsonName, _, _ := strings.Cut(fieldType.Tag.Get("json"), ",")
		if jsonName != name {
			continue
		}
		if !strings.HasPrefix(jsonName, "custom_field_") {
			return unsupportedTaskCustomFieldError(name)
		}
		field := v.Field(i)
		if !field.CanSet() || field.Kind() != reflect.Pointer || field.Type().Elem().Kind() != reflect.String {
			return fmt.Errorf("unsupported task custom field %q, expected string custom field", name)
		}

		valueCopy := value
		field.Set(reflect.ValueOf(&valueCopy))
		return nil
	}

	return unsupportedTaskCustomFieldError(name)
}

func unsupportedTaskCustomFieldError(name string) error {
	return fmt.Errorf(
		"unsupported task custom field %q, expected custom_field_one through custom_field_eight or custom_field_9 through custom_field_50",
		name,
	)
}

func decodeTaskBatchUpdate(workspaceID int, data []byte) (*tapd.BatchUpdateTasksRequest, error) {
	var request tapd.BatchUpdateTasksRequest
	if err := json.Unmarshal(data, &request); err == nil && len(request.Workitems) > 0 {
		request.WorkspaceID = new(workspaceID)
		for _, item := range request.Workitems {
			if item != nil && item.WorkspaceID == nil {
				item.WorkspaceID = new(workspaceID)
			}
		}
		return &request, nil
	}

	var items []*tapd.UpdateTaskRequest
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("decode batch update file as object or array: %w", err)
	}
	if len(items) == 0 {
		return nil, errors.New("batch update file must contain at least one task update")
	}
	for _, item := range items {
		if item != nil && item.WorkspaceID == nil {
			item.WorkspaceID = new(workspaceID)
		}
	}
	return &tapd.BatchUpdateTasksRequest{
		WorkspaceID: new(workspaceID),
		Workitems:   items,
	}, nil
}

type taskChangesFlags struct {
	workspaceID      int
	limit            int
	page             int
	fields           string
	ids              string
	taskID           int64
	creator          string
	created          string
	changeSummary    string
	comment          string
	changes          string
	entityType       string
	needParseChanges bool
}

func addTaskChangesFlags(cmd *cobra.Command, flags *taskChangesFlags, withPaging bool) {
	newWorkspaceFlag(cmd, &flags.workspaceID)
	if withPaging {
		newListFlags(cmd, &flags.limit, &flags.page)
		cmd.Flags().StringVar(&flags.fields, "fields", "", "comma separated fields")
		cmd.Flags().BoolVar(&flags.needParseChanges, "need-parse-changes", true, "return parsed field changes")
	}
	cmd.Flags().StringVar(&flags.ids, "ids", "", "comma separated change IDs")
	cmd.Flags().Int64Var(&flags.taskID, "task-id", 0, "task ID")
	cmd.Flags().StringVar(&flags.creator, "creator", "", "filter by creator")
	cmd.Flags().StringVar(&flags.created, "created", "", "filter by created time expression")
	cmd.Flags().StringVar(&flags.changeSummary, "change-summary", "", "filter by change summary")
	cmd.Flags().StringVar(&flags.comment, "comment", "", "filter by comment")
	cmd.Flags().StringVar(&flags.changes, "changes", "", "filter by changes text")
	cmd.Flags().StringVar(&flags.entityType, "entity-type", "", "filter by entity type")
}

func newTaskChangesRequest(flags taskChangesFlags) (*tapd.GetTaskChangesRequest, error) {
	request := &tapd.GetTaskChangesRequest{
		WorkspaceID: new(flags.workspaceID),
		Limit:       new(flags.limit),
		Page:        new(flags.page),
		Fields:      fieldsMulti(flags.fields),
	}
	if flags.ids != "" {
		ids, err := strictInt64Multi("change IDs", flags.ids)
		if err != nil {
			return nil, err
		}
		request.ID = ids
	}
	if flags.taskID > 0 {
		request.TaskID = new(flags.taskID)
	}
	if flags.creator != "" {
		request.Creator = new(flags.creator)
	}
	if flags.created != "" {
		request.Created = new(flags.created)
	}
	if flags.changeSummary != "" {
		request.ChangeSummary = new(flags.changeSummary)
	}
	if flags.comment != "" {
		request.Comment = new(flags.comment)
	}
	if flags.changes != "" {
		request.Changes = new(flags.changes)
	}
	if flags.entityType != "" {
		request.EntityType = new(flags.entityType)
	}
	if !flags.needParseChanges {
		request.NeedParseChanges = new(0)
	}
	return request, nil
}

func newTaskChangesCountRequest(flags taskChangesFlags) (*tapd.GetTaskChangesCountRequest, error) {
	request := &tapd.GetTaskChangesCountRequest{
		WorkspaceID: new(flags.workspaceID),
	}
	if flags.ids != "" {
		ids, err := strictInt64Multi("change IDs", flags.ids)
		if err != nil {
			return nil, err
		}
		request.ID = ids
	}
	if flags.taskID > 0 {
		request.TaskID = new(flags.taskID)
	}
	if flags.creator != "" {
		request.Creator = new(flags.creator)
	}
	if flags.created != "" {
		request.Created = new(flags.created)
	}
	if flags.changeSummary != "" {
		request.ChangeSummary = new(flags.changeSummary)
	}
	if flags.comment != "" {
		request.Comment = new(flags.comment)
	}
	if flags.changes != "" {
		request.Changes = new(flags.changes)
	}
	if flags.entityType != "" {
		request.EntityType = new(flags.entityType)
	}
	return request, nil
}

func taskStatusEnum(csv string) *tapd.Enum[tapd.TaskStatus] {
	items := splitCSV(csv)
	if len(items) == 0 {
		return nil
	}
	statuses := make([]tapd.TaskStatus, 0, len(items))
	for _, item := range items {
		statuses = append(statuses, tapd.TaskStatus(item))
	}
	return tapd.NewEnum(statuses...)
}

func taskTableHeaders() []string {
	return []string{"ID", "Name", "Status", "Owner", "Creator", "Progress"}
}

func taskRows(tasks []*tapd.Task) [][]string {
	rows := make([][]string, 0, len(tasks))
	for _, item := range tasks {
		rows = append(rows, []string{
			item.ID,
			item.Name,
			item.Status.String(),
			item.Owner,
			item.Creator,
			item.Progress,
		})
	}
	return rows
}
