package cmd

import (
	"fmt"
	"strconv"

	"github.com/go-tapd/cli/internal/app"
	"github.com/go-tapd/tapd"
	"github.com/spf13/cobra"
)

func newWorkspaceCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workspace",
		Short: "Work with TAPD workspaces",
	}

	cmd.AddCommand(
		newWorkspaceViewCmd(rt),
		newWorkspaceUsersCmd(rt),
		newWorkspaceRolesCmd(rt),
		newWorkspaceSubWorkspacesCmd(rt),
		newWorkspaceCompanyWorkspacesCmd(rt),
		newWorkspaceParticipantWorkspacesCmd(rt),
		newWorkspaceAddMemberCmd(rt),
		newWorkspaceUpdateCmd(rt),
		newWorkspaceCustomFieldsCmd(rt),
		newWorkspaceDocumentsCmd(rt),
		newWorkspaceShortIDCmd(rt),
		newWorkspaceMemberActivityLogCmd(rt),
		newWorkspaceCalendarCmd(rt),
	)

	return cmd
}

func newWorkspaceViewCmd(rt *app.Runtime) *cobra.Command {
	var workspaceID int

	cmd := &cobra.Command{
		Use:   "view",
		Short: "Show workspace details",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			workspace, _, err := client.WorkspaceService.GetWorkspaceInfo(cmd.Context(), &tapd.GetWorkspaceInfoRequest{
				WorkspaceID: new(workspaceID),
			})
			if err != nil {
				return err
			}

			rows := [][]string{{
				workspace.ID,
				workspace.Name,
				workspace.PrettyName,
				workspace.Status,
				workspace.Creator,
			}}

			return writeOutput(cmd, rt.OutputFormat, []string{"ID", "Name", "PrettyName", "Status", "Creator"}, rows, workspace)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	return cmd
}

func newWorkspaceUsersCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		fields      string
		user        string
	)

	cmd := &cobra.Command{
		Use:   "users",
		Short: "List workspace users",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetUsersRequest{
				WorkspaceID: new(workspaceID),
				Fields:      fieldsMulti(fields),
			}
			if user != "" {
				request.User = tapd.NewMulti(user)
			}

			users, _, err := client.WorkspaceService.GetUsers(cmd.Context(), request)
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(users))
			for _, item := range users {
				rows = append(rows, []string{
					item.User,
					item.Name,
					strconv.Itoa(len(item.RoleID)),
					item.Status,
					item.RealJoinTime,
				})
			}

			return writeOutput(cmd, rt.OutputFormat, []string{"User", "Name", "Roles", "Status", "Joined"}, rows, users)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().StringVar(&fields, "fields", "", "comma separated fields")
	cmd.Flags().StringVar(&user, "user", "", "filter by TAPD user")
	return cmd
}

func newWorkspaceRolesCmd(rt *app.Runtime) *cobra.Command {
	var workspaceID int

	cmd := &cobra.Command{
		Use:   "roles",
		Short: "List workspace roles",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			roles, _, err := client.WorkspaceService.GetWorkspaceRoles(cmd.Context(), &tapd.GetWorkspaceRolesRequest{
				WorkspaceID: new(workspaceID),
			})
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(roles))
			for _, item := range roles {
				if item == nil {
					continue
				}
				rows = append(rows, []string{item.ID, item.Name})
			}

			return writeOutput(cmd, rt.OutputFormat, []string{"ID", "Name"}, rows, roles)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	return cmd
}

func newWorkspaceSubWorkspacesCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		templateID  int
	)

	cmd := &cobra.Command{
		Use:   "sub-workspaces",
		Short: "List sub-workspaces",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetSubWorkspacesRequest{WorkspaceID: new(workspaceID)}
			if templateID > 0 {
				request.TemplateID = new(templateID)
			}

			workspaces, _, err := client.WorkspaceService.GetSubWorkspaces(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, workspaceHeaders(), workspaceRows(workspaces), workspaces)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().IntVar(&templateID, "template-id", 0, "sub-workspace template ID")
	return cmd
}

func newWorkspaceCompanyWorkspacesCmd(rt *app.Runtime) *cobra.Command {
	var (
		companyID   int
		category    string
		withExtends bool
	)

	cmd := &cobra.Command{
		Use:   "company-workspaces",
		Short: "List company workspaces",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetCompanyWorkspacesRequest{CompanyID: new(companyID)}
			if category != "" {
				request.Category = new(category)
			}
			if withExtends {
				request.WithExtends = new(1)
			}

			workspaces, _, err := client.WorkspaceService.GetCompanyWorkspaces(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, workspaceHeaders(), workspaceRows(workspaces), workspaces)
		},
	}

	cmd.Flags().IntVar(&companyID, "company-id", 0, "company ID")
	cmd.Flags().StringVar(&category, "category", "", "workspace category: project or mini_project")
	cmd.Flags().BoolVar(&withExtends, "with-extends", false, "include workspace extension data")
	_ = cmd.MarkFlagRequired("company-id")
	return cmd
}

func newWorkspaceParticipantWorkspacesCmd(rt *app.Runtime) *cobra.Command {
	var (
		nick      string
		companyID int
	)

	cmd := &cobra.Command{
		Use:   "participant-workspaces",
		Short: "List workspaces a user participates in",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			workspaces, _, err := client.WorkspaceService.GetUserParticipantWorkspaces(
				cmd.Context(),
				&tapd.GetUserParticipantWorkspacesRequest{
					Nick:      new(nick),
					CompanyID: new(companyID),
				},
			)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, workspaceHeaders(), workspaceRows(workspaces), workspaces)
		},
	}

	cmd.Flags().StringVar(&nick, "nick", "", "TAPD user nickname")
	cmd.Flags().IntVar(&companyID, "company-id", 0, "company ID")
	_ = cmd.MarkFlagRequired("nick")
	_ = cmd.MarkFlagRequired("company-id")
	return cmd
}

func newWorkspaceAddMemberCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		nick        string
		companyID   int
		roleIDs     string
	)

	cmd := &cobra.Command{
		Use:   "add-member",
		Short: "Add a workspace member",
		RunE: func(cmd *cobra.Command, _ []string) error {
			roles, err := strictInt64Multi("role IDs", roleIDs)
			if err != nil {
				return err
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.AddWorkspaceMemberRequest{
				WorkspaceID: new(workspaceID),
				Nick:        new(nick),
				RoleIDs:     roles,
			}
			if companyID > 0 {
				request.CompanyID = new(companyID)
			}

			result, _, err := client.WorkspaceService.AddWorkspaceMember(cmd.Context(), request)
			if err != nil {
				return err
			}

			rows := [][]string{{strconv.FormatBool(result.Success)}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Success"}, rows, result)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().StringVar(&nick, "nick", "", "TAPD user nickname")
	cmd.Flags().IntVar(&companyID, "company-id", 0, "company ID")
	cmd.Flags().StringVar(&roleIDs, "role-ids", "", "comma separated role IDs")
	_ = cmd.MarkFlagRequired("nick")
	_ = cmd.MarkFlagRequired("role-ids")
	return cmd
}

func newWorkspaceUpdateCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		field       string
		value       string
	)

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update workspace information",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			result, _, err := client.WorkspaceService.UpdateWorkspaceInfo(cmd.Context(), &tapd.UpdateWorkspaceInfoRequest{
				WorkspaceID: new(workspaceID),
				Field:       new(field),
				Value:       new(value),
			})
			if err != nil {
				return err
			}

			rows := [][]string{{field, result}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Field", "Result"}, rows, map[string]string{"field": field, "result": result})
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().StringVar(&field, "field", "", "workspace field name")
	cmd.Flags().StringVar(&value, "value", "", "workspace field value")
	_ = cmd.MarkFlagRequired("field")
	_ = cmd.MarkFlagRequired("value")
	return cmd
}

func newWorkspaceCustomFieldsCmd(rt *app.Runtime) *cobra.Command {
	var workspaceID int

	cmd := &cobra.Command{
		Use:   "custom-fields",
		Short: "List workspace custom field settings",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			settings, _, err := client.WorkspaceService.GetWorkspaceCustomFieldsSettings(
				cmd.Context(),
				&tapd.GetWorkspaceCustomFieldsSettingsRequest{WorkspaceID: new(workspaceID)},
			)
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(settings))
			for _, item := range settings {
				if item == nil {
					continue
				}
				rows = append(rows, []string{
					item.ID,
					item.CustomField,
					item.Name,
					item.Type,
					item.Enabled,
					item.Freeze,
					stringValue(item.Sort),
				})
			}

			return writeOutput(
				cmd,
				rt.OutputFormat,
				[]string{"ID", "CustomField", "Name", "Type", "Enabled", "Freeze", "Sort"},
				rows,
				settings,
			)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	return cmd
}

func newWorkspaceDocumentsCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		limit       int
		page        int
		fields      string
	)

	cmd := &cobra.Command{
		Use:   "documents",
		Short: "List workspace documents",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			documents, _, err := client.WorkspaceService.GetWorkspaceDocuments(cmd.Context(), &tapd.GetWorkspaceDocumentsRequest{
				WorkspaceID: new(workspaceID),
				Limit:       new(limit),
				Page:        new(page),
				Fields:      fieldsMulti(fields),
			})
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(documents))
			for _, item := range documents {
				if item == nil {
					continue
				}
				rows = append(rows, []string{
					item.ID,
					item.Name,
					item.Type,
					item.FolderID,
					item.Creator,
					stringValue(item.Status),
					item.Modified,
				})
			}

			return writeOutput(cmd, rt.OutputFormat, []string{"ID", "Name", "Type", "FolderID", "Creator", "Status", "Modified"}, rows, documents)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	newListFlags(cmd, &limit, &page)
	cmd.Flags().StringVar(&fields, "fields", "", "comma separated fields")
	return cmd
}

func newWorkspaceShortIDCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "short-id",
		Short: "Work with workspace short IDs",
	}
	cmd.AddCommand(newWorkspaceShortIDConvertCmd(rt))
	return cmd
}

func newWorkspaceShortIDConvertCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		entityType  string
		shortIDs    string
		longIDs     string
	)

	cmd := &cobra.Command{
		Use:   "convert",
		Short: "Convert work item short and long IDs",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if shortIDs == "" && longIDs == "" {
				return fmt.Errorf("one of --short-ids or --long-ids is required")
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetWorkItemsLongIDByShortIDsRequest{
				WorkspaceID: new(workspaceID),
				EntityType:  new(tapd.EntityType(entityType)),
			}
			if shortIDs != "" {
				request.ShortIDs = new(shortIDs)
			}
			if longIDs != "" {
				request.LongIDs = new(longIDs)
			}

			result, _, err := client.WorkspaceService.GetWorkItemsLongIDByShortIDs(cmd.Context(), request)
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(result.ValidIDMap))
			for _, item := range result.ValidIDMap {
				if item == nil {
					continue
				}
				rows = append(rows, []string{item.ShortID, item.LongID, string(item.EntityType), item.WorkspaceID, item.CompanyID})
			}
			return writeOutput(cmd, rt.OutputFormat, []string{"ShortID", "LongID", "EntityType", "WorkspaceID", "CompanyID"}, rows, result)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().StringVar(&entityType, "entity-type", "", "entity type: story, task, or bug")
	cmd.Flags().StringVar(&shortIDs, "short-ids", "", "semicolon separated short IDs")
	cmd.Flags().StringVar(&longIDs, "long-ids", "", "semicolon separated long IDs")
	_ = cmd.MarkFlagRequired("entity-type")
	return cmd
}

func newWorkspaceMemberActivityLogCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID    int
		companyOnly    bool
		limit          int
		page           int
		startTime      string
		endTime        string
		operator       string
		operateType    string
		operatorObject string
		ip             string
	)

	cmd := &cobra.Command{
		Use:   "member-activity-log",
		Short: "List member activity logs",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetMemberActivityLogRequest{
				WorkspaceID: new(workspaceID),
				Limit:       new(limit),
				Page:        new(page),
			}
			if companyOnly {
				request.CompanyOnly = new(1)
			}
			if startTime != "" {
				request.StartTime = new(startTime)
			}
			if endTime != "" {
				request.EndTime = new(endTime)
			}
			if operator != "" {
				request.Operator = new(operator)
			}
			if operateType != "" {
				request.OperateType = new(tapd.OperateType(operateType))
			}
			if operatorObject != "" {
				request.OperatorObject = new(tapd.OperateObject(operatorObject))
			}
			if ip != "" {
				request.IP = new(ip)
			}

			result, _, err := client.WorkspaceService.GetMemberActivityLog(cmd.Context(), request)
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(result.Records))
			for _, item := range result.Records {
				if item == nil {
					continue
				}
				rows = append(rows, []string{
					item.ID,
					item.Action,
					item.Creator,
					string(item.OperateType),
					string(item.OperateObject),
					item.Title,
					item.Created,
				})
			}
			return writeOutput(cmd, rt.OutputFormat, []string{"ID", "Action", "Creator", "OperateType", "OperateObject", "Title", "Created"}, rows, result)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	newListFlags(cmd, &limit, &page)
	cmd.Flags().BoolVar(&companyOnly, "company-only", false, "query company-level activity logs")
	cmd.Flags().StringVar(&startTime, "start-time", "", "start time, format YYYY-MM-DD HH:MM")
	cmd.Flags().StringVar(&endTime, "end-time", "", "end time, format YYYY-MM-DD HH:MM")
	cmd.Flags().StringVar(&operator, "operator", "", "filter by operator nickname")
	cmd.Flags().StringVar(&operateType, "operate-type", "", "operate type: add, delete, download, or upload")
	cmd.Flags().StringVar(&operatorObject, "operate-object", "", "operate object, such as story, bug, task, wiki")
	cmd.Flags().StringVar(&ip, "ip", "", "filter by IP")
	return cmd
}

func newWorkspaceCalendarCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "calendar",
		Short: "Work with workspace work calendars",
	}
	cmd.AddCommand(
		newWorkspaceCalendarSetCustomCmd(rt),
		newWorkspaceCalendarEnableCmd(rt),
		newWorkspaceCalendarViewCustomCmd(rt),
		newWorkspaceCalendarSettingsCmd(rt),
	)
	return cmd
}

func newWorkspaceCalendarSetCustomCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		year        string
		weekdays    string
		holidays    string
		workdays    string
	)

	cmd := &cobra.Command{
		Use:   "set-custom",
		Short: "Set custom work calendar",
		RunE: func(cmd *cobra.Command, _ []string) error {
			weekdayItems, err := optionalIntSlice("weekdays", weekdays)
			if err != nil {
				return err
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.SetCustomWorkCalendarRequest{
				WorkspaceID: new(workspaceID),
				Year:        new(year),
			}
			if weekdays != "" {
				request.Weekdays = new(weekdayItems)
			}
			if holidays != "" {
				request.Holidays = new(splitCSV(holidays))
			}
			if workdays != "" {
				request.Workdays = new(splitCSV(workdays))
			}

			result, _, err := client.WorkspaceService.SetCustomWorkCalendar(cmd.Context(), request)
			if err != nil {
				return err
			}

			rows := [][]string{{strconv.FormatBool(result.Success)}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Success"}, rows, result)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().StringVar(&year, "year", "", "calendar year")
	cmd.Flags().StringVar(&weekdays, "weekdays", "", "comma separated work weekdays, 1-7")
	cmd.Flags().StringVar(&holidays, "holidays", "", "comma separated extra holidays")
	cmd.Flags().StringVar(&workdays, "workdays", "", "comma separated extra workdays")
	_ = cmd.MarkFlagRequired("year")
	return cmd
}

func newWorkspaceCalendarEnableCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID  int
		calendarType string
	)

	cmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable a work calendar type",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			result, _, err := client.WorkspaceService.EnableWorkCalendar(cmd.Context(), &tapd.EnableWorkCalendarRequest{
				WorkspaceID: new(workspaceID),
				Type:        new(tapd.WorkCalendarType(calendarType)),
			})
			if err != nil {
				return err
			}

			rows := [][]string{{calendarType, strconv.FormatBool(result.Success)}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Type", "Success"}, rows, result)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().StringVar(&calendarType, "type", "", "calendar type: system or custom")
	_ = cmd.MarkFlagRequired("type")
	return cmd
}

func newWorkspaceCalendarViewCustomCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		year        string
	)

	cmd := &cobra.Command{
		Use:   "view-custom",
		Short: "Show custom work calendar",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			calendar, _, err := client.WorkspaceService.GetCustomWorkCalendar(cmd.Context(), &tapd.GetCustomWorkCalendarRequest{
				WorkspaceID: new(workspaceID),
				Year:        new(year),
			})
			if err != nil {
				return err
			}

			rows := [][]string{{strconv.Itoa(len(calendar.Weekdays)), strconv.Itoa(len(calendar.Holidays)), strconv.Itoa(len(calendar.Workdays))}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Weekdays", "Holidays", "Workdays"}, rows, calendar)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().StringVar(&year, "year", "", "calendar year")
	_ = cmd.MarkFlagRequired("year")
	return cmd
}

func newWorkspaceCalendarSettingsCmd(rt *app.Runtime) *cobra.Command {
	var workspaceID int

	cmd := &cobra.Command{
		Use:   "settings",
		Short: "List work calendar settings",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			settings, _, err := client.WorkspaceService.GetWorkCalendarSettings(cmd.Context(), &tapd.GetWorkCalendarSettingsRequest{
				WorkspaceID: new(workspaceID),
			})
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(settings))
			for _, item := range settings {
				if item == nil {
					continue
				}
				rows = append(rows, []string{item.Name, string(item.Type), strconv.FormatBool(item.Enable)})
			}
			return writeOutput(cmd, rt.OutputFormat, []string{"Name", "Type", "Enable"}, rows, settings)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	return cmd
}

func workspaceHeaders() []string {
	return []string{"ID", "Name", "PrettyName", "Category", "Status", "Creator"}
}

func workspaceRows(workspaces []*tapd.Workspace) [][]string {
	rows := make([][]string, 0, len(workspaces))
	for _, item := range workspaces {
		if item == nil {
			continue
		}
		rows = append(rows, []string{item.ID, item.Name, item.PrettyName, item.Category, item.Status, item.Creator})
	}
	return rows
}

func optionalIntSlice(name, csv string) ([]int, error) {
	items := splitCSV(csv)
	if len(items) == 0 {
		return nil, nil
	}

	values := make([]int, 0, len(items))
	for _, item := range items {
		value, err := parseIntArg(name, item)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return values, nil
}
