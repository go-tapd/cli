package cmd

import (
	"errors"
	"strconv"

	"github.com/go-tapd/cli/internal/app"
	"github.com/go-tapd/tapd"
	"github.com/spf13/cobra"
)

func newReleaseCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "release",
		Short: "Work with TAPD releases",
	}

	cmd.AddCommand(
		newReleaseCreateCmd(rt),
		newReleaseViewCmd(rt),
		newReleaseListCmd(rt),
		newReleaseCountCmd(rt),
		newReleaseUpdateCmd(rt),
	)
	return cmd
}

func newReleaseCreateCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		name        string
		description string
		startDate   string
		endDate     string
		creator     string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a release",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.CreateReleaseRequest{
				WorkspaceID: tapd.Ptr(workspaceID),
				Name:        tapd.Ptr(name),
				StartDate:   tapd.Ptr(startDate),
				EndDate:     tapd.Ptr(endDate),
			}
			if description != "" {
				request.Description = tapd.Ptr(description)
			}
			if creator != "" {
				request.Creator = tapd.Ptr(creator)
			}

			release, _, err := client.ReleaseService.CreateRelease(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, releaseHeaders(), releaseRows([]*tapd.Release{release}), release)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().StringVar(&name, "name", "", "release name")
	cmd.Flags().StringVar(&description, "description", "", "release description")
	cmd.Flags().StringVar(&startDate, "start-date", "", "start date")
	cmd.Flags().StringVar(&endDate, "end-date", "", "end date")
	cmd.Flags().StringVar(&creator, "creator", "", "release creator")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("start-date")
	_ = cmd.MarkFlagRequired("end-date")
	return cmd
}

func newReleaseViewCmd(rt *app.Runtime) *cobra.Command {
	var workspaceID int

	cmd := &cobra.Command{
		Use:   "view <id>",
		Short: "Show release details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseInt64Arg("release id", args[0])
			if err != nil {
				return err
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			releases, _, err := client.ReleaseService.GetReleases(cmd.Context(), &tapd.GetReleasesRequest{
				WorkspaceID: tapd.Ptr(workspaceID),
				ID:          tapd.NewMulti(id),
				Limit:       tapd.Ptr(1),
				Page:        tapd.Ptr(1),
			})
			if err != nil {
				return err
			}
			if len(releases) == 0 {
				return errors.New("release not found")
			}

			return writeOutput(cmd, rt.OutputFormat, releaseHeaders(), releaseRows(releases[:1]), releases[0])
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	return cmd
}

func newReleaseListCmd(rt *app.Runtime) *cobra.Command {
	var flags releaseQueryFlags

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List releases",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request, err := newGetReleasesRequest(flags)
			if err != nil {
				return err
			}

			releases, _, err := client.ReleaseService.GetReleases(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, releaseHeaders(), releaseRows(releases), releases)
		},
	}

	addReleaseQueryFlags(cmd, &flags, true)
	return cmd
}

func newReleaseCountCmd(rt *app.Runtime) *cobra.Command {
	var flags releaseQueryFlags

	cmd := &cobra.Command{
		Use:   "count",
		Short: "Count releases",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request, err := newGetReleasesCountRequest(flags)
			if err != nil {
				return err
			}

			count, _, err := client.ReleaseService.GetReleasesCount(cmd.Context(), request)
			if err != nil {
				return err
			}

			data := map[string]any{"resource": "release", "workspace_id": flags.workspaceID, "count": count}
			rows := [][]string{{"release", strconv.Itoa(flags.workspaceID), strconv.Itoa(count)}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Resource", "WorkspaceID", "Count"}, rows, data)
		},
	}

	addReleaseQueryFlags(cmd, &flags, false)
	return cmd
}

func newReleaseUpdateCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		name        string
		description string
		startDate   string
		endDate     string
		status      string
	)

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a release",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseInt64Arg("release id", args[0])
			if err != nil {
				return err
			}
			if name == "" && description == "" && startDate == "" && endDate == "" && status == "" {
				return errors.New("at least one release field flag is required")
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.UpdateReleaseRequest{
				WorkspaceID: tapd.Ptr(workspaceID),
				ID:          tapd.Ptr(id),
			}
			if name != "" {
				request.Name = tapd.Ptr(name)
			}
			if description != "" {
				request.Description = tapd.Ptr(description)
			}
			if startDate != "" {
				request.StartDate = tapd.Ptr(startDate)
			}
			if endDate != "" {
				request.EndDate = tapd.Ptr(endDate)
			}
			if status != "" {
				request.Status = tapd.Ptr(status)
			}

			release, _, err := client.ReleaseService.UpdateRelease(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, releaseHeaders(), releaseRows([]*tapd.Release{release}), release)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().StringVar(&name, "name", "", "release name")
	cmd.Flags().StringVar(&description, "description", "", "release description")
	cmd.Flags().StringVar(&startDate, "start-date", "", "start date")
	cmd.Flags().StringVar(&endDate, "end-date", "", "end date")
	cmd.Flags().StringVar(&status, "status", "", "release status")
	return cmd
}

type releaseQueryFlags struct {
	workspaceID int
	limit       int
	page        int
	fields      string
	ids         string
	name        string
	description string
	startDate   string
	endDate     string
	creator     string
	created     string
	modified    string
	status      string
}

func addReleaseQueryFlags(cmd *cobra.Command, flags *releaseQueryFlags, withPaging bool) {
	newWorkspaceFlag(cmd, &flags.workspaceID)
	if withPaging {
		newListFlags(cmd, &flags.limit, &flags.page)
		cmd.Flags().StringVar(&flags.fields, "fields", "", "comma separated fields")
	}
	cmd.Flags().StringVar(&flags.ids, "ids", "", "comma separated release IDs")
	cmd.Flags().StringVar(&flags.name, "name", "", "filter by release name")
	cmd.Flags().StringVar(&flags.description, "description", "", "filter by description")
	cmd.Flags().StringVar(&flags.startDate, "start-date", "", "filter by start date")
	cmd.Flags().StringVar(&flags.endDate, "end-date", "", "filter by end date")
	cmd.Flags().StringVar(&flags.creator, "creator", "", "filter by creator")
	cmd.Flags().StringVar(&flags.created, "created", "", "filter by created time expression")
	cmd.Flags().StringVar(&flags.modified, "modified", "", "filter by modified time expression")
	cmd.Flags().StringVar(&flags.status, "status", "", "filter by status")
}

func newGetReleasesRequest(flags releaseQueryFlags) (*tapd.GetReleasesRequest, error) {
	request := &tapd.GetReleasesRequest{
		WorkspaceID: tapd.Ptr(flags.workspaceID),
		Limit:       tapd.Ptr(flags.limit),
		Page:        tapd.Ptr(flags.page),
		Fields:      fieldsMulti(flags.fields),
	}
	if err := applyReleaseFilters(request, flags); err != nil {
		return nil, err
	}
	return request, nil
}

func newGetReleasesCountRequest(flags releaseQueryFlags) (*tapd.GetReleasesCountRequest, error) {
	request := &tapd.GetReleasesCountRequest{WorkspaceID: tapd.Ptr(flags.workspaceID)}
	if flags.ids != "" {
		ids, err := strictInt64Multi("release IDs", flags.ids)
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
	if flags.creator != "" {
		request.Creator = tapd.Ptr(flags.creator)
	}
	if flags.created != "" {
		request.Created = tapd.Ptr(flags.created)
	}
	if flags.modified != "" {
		request.Modified = tapd.Ptr(flags.modified)
	}
	if flags.status != "" {
		request.Status = tapd.Ptr(flags.status)
	}
	return request, nil
}

func applyReleaseFilters(request *tapd.GetReleasesRequest, flags releaseQueryFlags) error {
	if flags.ids != "" {
		ids, err := strictInt64Multi("release IDs", flags.ids)
		if err != nil {
			return err
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
	if flags.creator != "" {
		request.Creator = tapd.Ptr(flags.creator)
	}
	if flags.created != "" {
		request.Created = tapd.Ptr(flags.created)
	}
	if flags.modified != "" {
		request.Modified = tapd.Ptr(flags.modified)
	}
	if flags.status != "" {
		request.Status = tapd.Ptr(flags.status)
	}
	return nil
}

func releaseHeaders() []string {
	return []string{"ID", "Name", "Status", "StartDate", "EndDate", "Creator", "Modified"}
}

func releaseRows(releases []*tapd.Release) [][]string {
	rows := make([][]string, 0, len(releases))
	for _, item := range releases {
		if item == nil {
			continue
		}
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

func newLaunchFormCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "launch-form",
		Short: "Work with TAPD launch forms",
	}

	cmd.AddCommand(
		newLaunchFormCreateCmd(rt),
		newLaunchFormListCmd(rt),
		newLaunchFormCountCmd(rt),
		newLaunchFormFieldsCmd(rt),
		newLaunchFormTemplatesCmd(rt),
		newLaunchFormLogsCmd(rt),
	)
	return cmd
}

func newLaunchFormCreateCmd(rt *app.Runtime) *cobra.Command {
	var flags launchFormCreateFlags

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a launch form",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.CreateLaunchFormRequest{
				WorkspaceID: tapd.Ptr(flags.workspaceID),
				Creator:     tapd.Ptr(flags.creator),
				TemplateID:  tapd.Ptr(flags.templateID),
			}
			applyLaunchFormCreateFlags(request, flags)

			form, _, err := client.ReleaseService.CreateLaunchForm(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, launchFormHeaders(), launchFormRows([]*tapd.LaunchForm{form}), form)
		},
	}

	newWorkspaceFlag(cmd, &flags.workspaceID)
	cmd.Flags().StringVar(&flags.creator, "creator", "", "launch form creator")
	cmd.Flags().StringVar(&flags.templateID, "template-id", "", "launch form template ID")
	cmd.Flags().StringVar(&flags.title, "title", "", "launch form title")
	cmd.Flags().StringVar(&flags.versionType, "version-type", "", "version type")
	cmd.Flags().StringVar(&flags.baseline, "baseline", "", "baseline")
	cmd.Flags().StringVar(&flags.releaseModel, "release-model", "", "release model")
	cmd.Flags().StringVar(&flags.roadmapVersion, "roadmap-version", "", "roadmap version")
	cmd.Flags().StringVar(&flags.releaseType, "release-type", "", "release type")
	cmd.Flags().StringVar(&flags.signedBy, "signed-by", "", "signer")
	cmd.Flags().StringVar(&flags.archivedBy, "archived-by", "", "archive confirmer")
	cmd.Flags().StringVar(&flags.cc, "cc", "", "CC users")
	_ = cmd.MarkFlagRequired("creator")
	_ = cmd.MarkFlagRequired("template-id")
	return cmd
}

func newLaunchFormListCmd(rt *app.Runtime) *cobra.Command {
	var flags launchFormQueryFlags

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List launch forms",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := newGetLaunchFormsRequest(flags)
			forms, _, err := client.ReleaseService.GetLaunchForms(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, launchFormHeaders(), launchFormRows(forms), forms)
		},
	}

	addLaunchFormQueryFlags(cmd, &flags, true)
	return cmd
}

func newLaunchFormCountCmd(rt *app.Runtime) *cobra.Command {
	var flags launchFormQueryFlags

	cmd := &cobra.Command{
		Use:   "count",
		Short: "Count launch forms",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			count, _, err := client.ReleaseService.GetLaunchFormsCount(cmd.Context(), newGetLaunchFormsCountRequest(flags))
			if err != nil {
				return err
			}

			data := map[string]any{"resource": "launch_form", "workspace_id": flags.workspaceID, "count": count}
			rows := [][]string{{"launch_form", strconv.Itoa(flags.workspaceID), strconv.Itoa(count)}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Resource", "WorkspaceID", "Count"}, rows, data)
		},
	}

	addLaunchFormQueryFlags(cmd, &flags, false)
	return cmd
}

func newLaunchFormFieldsCmd(rt *app.Runtime) *cobra.Command {
	var workspaceID int

	cmd := &cobra.Command{
		Use:   "fields",
		Short: "List launch form custom field settings",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			settings, _, err := client.ReleaseService.GetLaunchFormCustomFieldsSettings(
				cmd.Context(),
				&tapd.GetLaunchFormCustomFieldsSettingsRequest{WorkspaceID: tapd.Ptr(workspaceID)},
			)
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(settings))
			for _, item := range settings {
				if item == nil {
					continue
				}
				rows = append(rows, []string{item.ID, item.CustomField, item.Name, item.Type, item.Enabled, stringValue(item.Sort)})
			}
			return writeOutput(cmd, rt.OutputFormat, []string{"ID", "CustomField", "Name", "Type", "Enabled", "Sort"}, rows, settings)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	return cmd
}

func newLaunchFormTemplatesCmd(rt *app.Runtime) *cobra.Command {
	var workspaceID int

	cmd := &cobra.Command{
		Use:   "templates",
		Short: "List launch form templates",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			templates, _, err := client.ReleaseService.GetLaunchFormTemplates(
				cmd.Context(),
				&tapd.GetLaunchFormTemplatesRequest{WorkspaceID: tapd.Ptr(workspaceID)},
			)
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(templates))
			for _, item := range templates {
				if item == nil {
					continue
				}
				rows = append(rows, []string{item.ID, item.Name})
			}
			return writeOutput(cmd, rt.OutputFormat, []string{"ID", "Name"}, rows, templates)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	return cmd
}

func newLaunchFormLogsCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		formID      int64
	)

	cmd := &cobra.Command{
		Use:   "logs",
		Short: "List launch form activity logs",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			logs, _, err := client.ReleaseService.GetLaunchFormActivityLogs(
				cmd.Context(),
				&tapd.GetLaunchFormActivityLogsRequest{
					WorkspaceID: tapd.Ptr(workspaceID),
					FormID:      tapd.Ptr(formID),
				},
			)
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(logs))
			for _, item := range logs {
				if item == nil {
					continue
				}
				rows = append(rows, []string{
					item.ID,
					item.FormID,
					item.Field,
					item.Operation,
					item.CreatedBy,
					item.Created,
				})
			}
			return writeOutput(cmd, rt.OutputFormat, []string{"ID", "FormID", "Field", "Operation", "CreatedBy", "Created"}, rows, logs)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().Int64Var(&formID, "form-id", 0, "launch form ID")
	_ = cmd.MarkFlagRequired("form-id")
	return cmd
}

type launchFormCreateFlags struct {
	workspaceID    int
	creator        string
	templateID     string
	title          string
	versionType    string
	baseline       string
	releaseModel   string
	roadmapVersion string
	releaseType    string
	signedBy       string
	archivedBy     string
	cc             string
}

func applyLaunchFormCreateFlags(request *tapd.CreateLaunchFormRequest, flags launchFormCreateFlags) {
	if flags.title != "" {
		request.Title = tapd.Ptr(flags.title)
	}
	if flags.versionType != "" {
		request.VersionType = tapd.Ptr(flags.versionType)
	}
	if flags.baseline != "" {
		request.Baseline = tapd.Ptr(flags.baseline)
	}
	if flags.releaseModel != "" {
		request.ReleaseModel = tapd.Ptr(flags.releaseModel)
	}
	if flags.roadmapVersion != "" {
		request.RoadmapVersion = tapd.Ptr(flags.roadmapVersion)
	}
	if flags.releaseType != "" {
		request.ReleaseType = tapd.Ptr(flags.releaseType)
	}
	if flags.signedBy != "" {
		request.SignedBy = tapd.Ptr(flags.signedBy)
	}
	if flags.archivedBy != "" {
		request.ArchivedBy = tapd.Ptr(flags.archivedBy)
	}
	if flags.cc != "" {
		request.CC = tapd.Ptr(flags.cc)
	}
}

type launchFormQueryFlags struct {
	workspaceID    int
	limit          int
	page           int
	fields         string
	id             int64
	creator        string
	created        string
	title          string
	status         string
	versionType    string
	baseline       string
	releaseModel   string
	roadmapVersion string
	releaseType    string
	changeType     string
	signedBy       string
	archivedBy     string
	cc             string
	changeNotifier string
}

func addLaunchFormQueryFlags(cmd *cobra.Command, flags *launchFormQueryFlags, withPaging bool) {
	newWorkspaceFlag(cmd, &flags.workspaceID)
	if withPaging {
		newListFlags(cmd, &flags.limit, &flags.page)
		cmd.Flags().StringVar(&flags.fields, "fields", "", "comma separated fields")
	}
	cmd.Flags().Int64Var(&flags.id, "id", 0, "launch form ID")
	cmd.Flags().StringVar(&flags.creator, "creator", "", "filter by creator")
	cmd.Flags().StringVar(&flags.created, "created", "", "filter by created time expression")
	cmd.Flags().StringVar(&flags.title, "title", "", "filter by title")
	cmd.Flags().StringVar(&flags.status, "status", "", "filter by status")
	cmd.Flags().StringVar(&flags.versionType, "version-type", "", "filter by version type")
	cmd.Flags().StringVar(&flags.baseline, "baseline", "", "filter by baseline")
	cmd.Flags().StringVar(&flags.releaseModel, "release-model", "", "filter by release model")
	cmd.Flags().StringVar(&flags.roadmapVersion, "roadmap-version", "", "filter by roadmap version")
	cmd.Flags().StringVar(&flags.releaseType, "release-type", "", "filter by release type")
	cmd.Flags().StringVar(&flags.changeType, "change-type", "", "filter by change type")
	cmd.Flags().StringVar(&flags.signedBy, "signed-by", "", "filter by signer")
	cmd.Flags().StringVar(&flags.archivedBy, "archived-by", "", "filter by archive confirmer")
	cmd.Flags().StringVar(&flags.cc, "cc", "", "filter by CC users")
	cmd.Flags().StringVar(&flags.changeNotifier, "change-notifier", "", "filter by change notifier")
}

func newGetLaunchFormsRequest(flags launchFormQueryFlags) *tapd.GetLaunchFormsRequest {
	request := &tapd.GetLaunchFormsRequest{
		WorkspaceID: tapd.Ptr(flags.workspaceID),
		Limit:       tapd.Ptr(flags.limit),
		Page:        tapd.Ptr(flags.page),
		Fields:      fieldsMulti(flags.fields),
	}
	applyLaunchFormFilters(request, flags)
	return request
}

func newGetLaunchFormsCountRequest(flags launchFormQueryFlags) *tapd.GetLaunchFormsCountRequest {
	request := &tapd.GetLaunchFormsCountRequest{WorkspaceID: tapd.Ptr(flags.workspaceID)}
	if flags.id > 0 {
		request.ID = tapd.Ptr(flags.id)
	}
	if flags.creator != "" {
		request.Creator = tapd.Ptr(flags.creator)
	}
	if flags.created != "" {
		request.Created = tapd.Ptr(flags.created)
	}
	if flags.title != "" {
		request.Title = tapd.Ptr(flags.title)
	}
	if flags.status != "" {
		request.Status = tapd.Ptr(flags.status)
	}
	if flags.versionType != "" {
		request.VersionType = tapd.Ptr(flags.versionType)
	}
	if flags.baseline != "" {
		request.Baseline = tapd.Ptr(flags.baseline)
	}
	if flags.releaseModel != "" {
		request.ReleaseModel = tapd.Ptr(flags.releaseModel)
	}
	if flags.roadmapVersion != "" {
		request.RoadmapVersion = tapd.Ptr(flags.roadmapVersion)
	}
	if flags.releaseType != "" {
		request.ReleaseType = tapd.Ptr(flags.releaseType)
	}
	if flags.changeType != "" {
		request.ChangeType = tapd.Ptr(flags.changeType)
	}
	if flags.signedBy != "" {
		request.SignedBy = tapd.Ptr(flags.signedBy)
	}
	if flags.archivedBy != "" {
		request.ArchivedBy = tapd.Ptr(flags.archivedBy)
	}
	if flags.cc != "" {
		request.CC = tapd.Ptr(flags.cc)
	}
	if flags.changeNotifier != "" {
		request.ChangeNotifier = tapd.Ptr(flags.changeNotifier)
	}
	return request
}

func applyLaunchFormFilters(request *tapd.GetLaunchFormsRequest, flags launchFormQueryFlags) {
	if flags.id > 0 {
		request.ID = tapd.Ptr(flags.id)
	}
	if flags.creator != "" {
		request.Creator = tapd.Ptr(flags.creator)
	}
	if flags.created != "" {
		request.Created = tapd.Ptr(flags.created)
	}
	if flags.title != "" {
		request.Title = tapd.Ptr(flags.title)
	}
	if flags.status != "" {
		request.Status = tapd.Ptr(flags.status)
	}
	if flags.versionType != "" {
		request.VersionType = tapd.Ptr(flags.versionType)
	}
	if flags.baseline != "" {
		request.Baseline = tapd.Ptr(flags.baseline)
	}
	if flags.releaseModel != "" {
		request.ReleaseModel = tapd.Ptr(flags.releaseModel)
	}
	if flags.roadmapVersion != "" {
		request.RoadmapVersion = tapd.Ptr(flags.roadmapVersion)
	}
	if flags.releaseType != "" {
		request.ReleaseType = tapd.Ptr(flags.releaseType)
	}
	if flags.changeType != "" {
		request.ChangeType = tapd.Ptr(flags.changeType)
	}
	if flags.signedBy != "" {
		request.SignedBy = tapd.Ptr(flags.signedBy)
	}
	if flags.archivedBy != "" {
		request.ArchivedBy = tapd.Ptr(flags.archivedBy)
	}
	if flags.cc != "" {
		request.CC = tapd.Ptr(flags.cc)
	}
	if flags.changeNotifier != "" {
		request.ChangeNotifier = tapd.Ptr(flags.changeNotifier)
	}
}

func launchFormHeaders() []string {
	return []string{"ID", "Name", "Title", "Status", "Creator", "ReleaseID", "Created"}
}

func launchFormRows(forms []*tapd.LaunchForm) [][]string {
	rows := make([][]string, 0, len(forms))
	for _, item := range forms {
		if item == nil {
			continue
		}
		rows = append(rows, []string{
			item.ID,
			item.Name,
			stringValue(item.Title),
			item.Status,
			item.Creator,
			stringValue(item.ReleaseID),
			item.Created,
		})
	}
	return rows
}
