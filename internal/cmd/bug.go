package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/go-tapd/cli/internal/app"
	"github.com/go-tapd/tapd"
	"github.com/spf13/cobra"
)

func newBugCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bug",
		Short: "Work with TAPD bugs",
	}

	cmd.AddCommand(
		newBugCreateCmd(rt),
		newBugViewCmd(rt),
		newBugListCmd(rt),
		newBugCountCmd(rt),
		newBugUpdateCmd(rt),
		newBugCopyCmd(rt),
		newBugBatchUpdateCmd(rt),
		newBugChangesCmd(rt),
		newBugFieldsCmd(rt),
		newBugCustomFieldSettingsCmd(rt),
		newBugFieldLabelsCmd(rt),
		newBugTemplatesCmd(rt),
		newBugTemplateFieldsCmd(rt),
		newBugRemovedCmd(rt),
		newBugRelatedStoriesCmd(rt),
		newBugLinksCmd(rt),
		newBugLinkCmd(rt),
		newBugUnlinkCmd(rt),
		newBugUpdateSystemOptionsCmd(rt),
		newBugByViewCmd(rt),
		newBugConvertIDsCmd(rt),
	)
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
				WorkspaceID: new(workspaceID),
				Limit:       new(limit),
				Page:        new(page),
				Fields:      fieldsMulti(fields),
				ID:          int64Multi(ids),
			}
			if creator != "" {
				request.Reporter = tapd.NewMulti(creator)
			}
			if owner != "" {
				request.CurrentOwner = new(owner)
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

func newBugCreateCmd(rt *app.Runtime) *cobra.Command {
	var flags bugMutationFlags

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a bug",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.CreateBugRequest{
				WorkspaceID: new(flags.workspaceID),
				Title:       new(flags.title),
			}
			applyBugCreateFlags(request, flags)

			bug, _, err := client.BugService.CreateBug(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, bugTableHeaders(), bugRows([]*tapd.Bug{bug}), bug)
		},
	}

	newWorkspaceFlag(cmd, &flags.workspaceID)
	cmd.Flags().StringVar(&flags.title, "title", "", "bug title")
	_ = cmd.MarkFlagRequired("title")
	addBugMutationFlags(cmd, &flags, false)
	return cmd
}

func newBugViewCmd(rt *app.Runtime) *cobra.Command {
	var workspaceID int

	cmd := &cobra.Command{
		Use:   "view <id>",
		Short: "Show bug details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseInt64Arg("bug id", args[0])
			if err != nil {
				return err
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			bugs, _, err := client.BugService.GetBugs(cmd.Context(), &tapd.GetBugsRequest{
				WorkspaceID: new(workspaceID),
				ID:          tapd.NewMulti(id),
				Limit:       new(1),
				Page:        new(1),
			})
			if err != nil {
				return err
			}
			if len(bugs) == 0 {
				return fmt.Errorf("bug %d not found", id)
			}

			return writeOutput(cmd, rt.OutputFormat, bugTableHeaders(), bugRows(bugs[:1]), bugs[0])
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	return cmd
}

func newBugCountCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		ids         string
		creator     string
		owner       string
	)

	cmd := &cobra.Command{
		Use:   "count",
		Short: "Count bugs",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetBugsCountRequest{
				WorkspaceID: new(workspaceID),
				ID:          int64Multi(ids),
			}
			if creator != "" {
				request.Reporter = tapd.NewMulti(creator)
			}
			if owner != "" {
				request.CurrentOwner = new(owner)
			}

			count, _, err := client.BugService.GetBugsCount(cmd.Context(), request)
			if err != nil {
				return err
			}

			data := map[string]any{
				"resource":     "bug",
				"workspace_id": workspaceID,
				"count":        count,
			}
			rows := [][]string{{"bug", strconv.Itoa(workspaceID), strconv.Itoa(count)}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Resource", "WorkspaceID", "Count"}, rows, data)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().StringVar(&ids, "ids", "", "comma separated bug IDs")
	cmd.Flags().StringVar(&creator, "creator", "", "filter by creator")
	cmd.Flags().StringVar(&owner, "owner", "", "filter by current owner")
	return cmd
}

func newBugCopyCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID    int
		sourceBugID    int64
		dstWorkspaceID int
		syncFields     string
	)

	cmd := &cobra.Command{
		Use:   "copy",
		Short: "Copy a bug to another workspace",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			bug, _, err := client.BugService.CopyBug(cmd.Context(), &tapd.CopyBugRequest{
				WorkspaceID:    new(workspaceID),
				SourceBugID:    new(sourceBugID),
				DstWorkspaceID: new(dstWorkspaceID),
				SyncFields:     fieldsMulti(syncFields),
			})
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, bugTableHeaders(), bugRows([]*tapd.Bug{bug}), bug)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().Int64Var(&sourceBugID, "source-bug-id", 0, "source bug ID")
	cmd.Flags().IntVar(&dstWorkspaceID, "dst-workspace-id", 0, "destination workspace ID")
	cmd.Flags().StringVar(&syncFields, "sync-fields", "", "comma separated fields to copy")
	_ = cmd.MarkFlagRequired("source-bug-id")
	_ = cmd.MarkFlagRequired("dst-workspace-id")
	return cmd
}

func newBugFieldsCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		allOptions  bool
	)

	cmd := &cobra.Command{
		Use:   "fields",
		Short: "List bug fields",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetBugFieldsInfoRequest{
				WorkspaceID: new(workspaceID),
			}
			if allOptions {
				request.AllOptions = new(1)
			}

			fields, _, err := client.BugService.GetBugFieldsInfo(cmd.Context(), request)
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
	cmd.Flags().BoolVar(&allOptions, "all-options", false, "include closed field options")
	return cmd
}

func newBugCustomFieldSettingsCmd(rt *app.Runtime) *cobra.Command {
	var workspaceID int

	cmd := &cobra.Command{
		Use:   "custom-field-settings",
		Short: "List bug custom field settings",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			settings, _, err := client.BugService.GetBugCustomFieldsSettings(
				cmd.Context(),
				&tapd.GetBugCustomFieldsSettingsRequest{WorkspaceID: new(workspaceID)},
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

func newBugUpdateCmd(rt *app.Runtime) *cobra.Command {
	var flags bugMutationFlags

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a bug",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseInt64Arg("bug id", args[0])
			if err != nil {
				return err
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.UpdateBugRequest{
				ID:          new(id),
				WorkspaceID: new(flags.workspaceID),
			}
			applyBugUpdateFlags(request, flags)

			bug, _, err := client.BugService.UpdateBug(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, bugTableHeaders(), bugRows([]*tapd.Bug{bug}), bug)
		},
	}

	newWorkspaceFlag(cmd, &flags.workspaceID)
	addBugMutationFlags(cmd, &flags, true)
	return cmd
}

func newBugUpdateSystemOptionsCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		field       string
		values      string
	)

	cmd := &cobra.Command{
		Use:   "update-system-options",
		Short: "Update bug system select field options",
		RunE: func(cmd *cobra.Command, _ []string) error {
			items := splitCSV(values)
			if len(items) == 0 {
				return errors.New("values cannot be empty")
			}

			options := make([]*tapd.BugSystemSelectFieldOption, 0, len(items))
			for _, item := range items {
				options = append(options, &tapd.BugSystemSelectFieldOption{Value: new(item)})
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			ok, _, err := client.BugService.UpdateBugSystemSelectFieldOptions(
				cmd.Context(),
				&tapd.UpdateBugSystemSelectFieldOptionsRequest{
					WorkspaceID: new(workspaceID),
					Field:       new(field),
					Options:     options,
				},
			)
			if err != nil {
				return err
			}

			rows := [][]string{{field, strconv.Itoa(len(options)), strconv.FormatBool(ok)}}
			return writeOutput(
				cmd,
				rt.OutputFormat,
				[]string{"Field", "Options", "Success"},
				rows,
				map[string]any{"field": field, "options": options, "success": ok},
			)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().StringVar(&field, "field", "bugtype", "system select field name")
	cmd.Flags().StringVar(&values, "values", "", "comma separated option values")
	_ = cmd.MarkFlagRequired("values")
	return cmd
}

func newBugBatchUpdateCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		file        string
	)

	cmd := &cobra.Command{
		Use:   "batch-update",
		Short: "Batch update bugs from a JSON file",
		RunE: func(cmd *cobra.Command, _ []string) error {
			data, err := os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("read batch update file: %w", err)
			}

			request, err := decodeBugBatchUpdate(workspaceID, data)
			if err != nil {
				return err
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			result, _, err := client.BugService.BatchUpdateBugs(cmd.Context(), request)
			if err != nil {
				return err
			}

			rows := [][]string{{result.Msg}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Message"}, rows, result)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().StringVar(&file, "file", "", "JSON file containing an array of bug update objects")
	_ = cmd.MarkFlagRequired("file")
	return cmd
}

func newBugChangesCmd(rt *app.Runtime) *cobra.Command {
	var flags bugChangesFlags

	cmd := &cobra.Command{
		Use:   "changes",
		Short: "List bug changes",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request, err := newBugChangesRequest(flags)
			if err != nil {
				return err
			}

			changes, _, err := client.BugService.GetBugChanges(cmd.Context(), request)
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(changes))
			for _, item := range changes {
				memo := ""
				if item.Memo != nil {
					memo = *item.Memo
				}
				rows = append(rows, []string{
					item.ID,
					item.BugID,
					item.Field,
					item.OldValue,
					item.NewValue,
					item.Author,
					item.Created,
					memo,
				})
			}

			return writeOutput(
				cmd,
				rt.OutputFormat,
				[]string{"ID", "BugID", "Field", "OldValue", "NewValue", "Author", "Created", "Memo"},
				rows,
				changes,
			)
		},
	}

	addBugChangesFlags(cmd, &flags, true)
	cmd.AddCommand(newBugChangesCountCmd(rt))
	return cmd
}

func newBugChangesCountCmd(rt *app.Runtime) *cobra.Command {
	var flags bugChangesFlags

	cmd := &cobra.Command{
		Use:   "count",
		Short: "Count bug changes",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request, err := newBugChangesCountRequest(flags)
			if err != nil {
				return err
			}

			count, _, err := client.BugService.GetBugChangesCount(cmd.Context(), request)
			if err != nil {
				return err
			}

			data := map[string]any{"resource": "bug_change", "workspace_id": flags.workspaceID, "count": count}
			rows := [][]string{{"bug_change", strconv.Itoa(flags.workspaceID), strconv.Itoa(count)}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Resource", "WorkspaceID", "Count"}, rows, data)
		},
	}

	addBugChangesFlags(cmd, &flags, false)
	return cmd
}

func newBugFieldLabelsCmd(rt *app.Runtime) *cobra.Command {
	var workspaceID int

	cmd := &cobra.Command{
		Use:   "field-labels",
		Short: "List bug field labels",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			labels, _, err := client.BugService.GetBugFieldsLabel(cmd.Context(), &tapd.GetBugFieldsLabelRequest{
				WorkspaceID: new(workspaceID),
			})
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(labels))
			for _, item := range labels {
				rows = append(rows, []string{item.EN, item.CN})
			}

			return writeOutput(cmd, rt.OutputFormat, []string{"Field", "Label"}, rows, labels)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	return cmd
}

func newBugTemplatesCmd(rt *app.Runtime) *cobra.Command {
	var workspaceID int

	cmd := &cobra.Command{
		Use:   "templates",
		Short: "List bug templates",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			templates, _, err := client.BugService.GetBugTemplates(cmd.Context(), &tapd.GetBugTemplatesRequest{
				WorkspaceID: new(workspaceID),
			})
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(templates))
			for _, item := range templates {
				rows = append(rows, []string{item.ID, item.Name, item.Default, item.Creator, item.EditorType})
			}

			return writeOutput(cmd, rt.OutputFormat, []string{"ID", "Name", "Default", "Creator", "EditorType"}, rows, templates)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	return cmd
}

func newBugTemplateFieldsCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID      int
		templateID       int64
		usePriorityLabel bool
	)

	cmd := &cobra.Command{
		Use:   "template-fields",
		Short: "List bug template fields",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetBugTemplateFieldsRequest{
				WorkspaceID: new(workspaceID),
				TemplateID:  new(templateID),
			}
			if usePriorityLabel {
				request.UsePriorityLabel = new(1)
			}

			fields, _, err := client.BugService.GetBugTemplateFields(cmd.Context(), request)
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
	cmd.Flags().Int64Var(&templateID, "template-id", 0, "bug template ID")
	cmd.Flags().BoolVar(&usePriorityLabel, "use-priority-label", false, "replace priority field with priority_label")
	_ = cmd.MarkFlagRequired("template-id")
	return cmd
}

func newBugRemovedCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		limit       int
		page        int
		ids         string
		creator     string
		created     string
		modified    string
		includeAll  bool
	)

	cmd := &cobra.Command{
		Use:   "removed",
		Short: "List removed bugs",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetRemovedBugsRequest{
				WorkspaceID: new(workspaceID),
				Limit:       new(limit),
				Page:        new(page),
			}
			if ids != "" {
				request.ID, err = strictInt64Multi("bug IDs", ids)
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
			if modified != "" {
				request.Modified = new(modified)
			}
			if includeAll {
				request.IncludeAll = new(1)
			}

			bugs, _, err := client.BugService.GetRemovedBugs(cmd.Context(), request)
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(bugs))
			for _, item := range bugs {
				rows = append(rows, []string{item.ID, item.Title, item.Reporter, item.OperationUser, item.Type, item.Modified})
			}

			return writeOutput(cmd, rt.OutputFormat, []string{"ID", "Title", "Reporter", "OperationUser", "Type", "Modified"}, rows, bugs)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	newListFlags(cmd, &limit, &page)
	cmd.Flags().StringVar(&ids, "ids", "", "comma separated bug IDs")
	cmd.Flags().StringVar(&creator, "creator", "", "filter by creator")
	cmd.Flags().StringVar(&created, "created", "", "filter by created time expression")
	cmd.Flags().StringVar(&modified, "modified", "", "filter by removed time expression")
	cmd.Flags().BoolVar(&includeAll, "include-all", false, "include moved, merged, and deleted bugs")
	return cmd
}

func newBugLinksCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		bugID       int64
	)

	cmd := &cobra.Command{
		Use:   "links",
		Short: "List bug link relations",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			relations, _, err := client.BugService.GetBugLinkBugs(cmd.Context(), &tapd.GetBugLinkBugsRequest{
				WorkspaceID: new(workspaceID),
				BugID:       new(bugID),
			})
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(relations))
			for _, item := range relations {
				rows = append(rows, []string{
					item.LinkID,
					item.ID,
					item.Type,
					item.ActAs,
					item.WorkspaceID,
					strconv.Itoa(item.LinkedWorkspaceID),
				})
			}

			return writeOutput(
				cmd,
				rt.OutputFormat,
				[]string{"LinkID", "BugID", "Type", "ActAs", "WorkspaceID", "LinkedWorkspaceID"},
				rows,
				relations,
			)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().Int64Var(&bugID, "bug-id", 0, "bug ID")
	_ = cmd.MarkFlagRequired("bug-id")
	return cmd
}

func newBugRelatedStoriesCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		bugIDs      string
	)

	cmd := &cobra.Command{
		Use:   "related-stories",
		Short: "List stories related to bugs",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ids, err := strictInt64Multi("bug IDs", bugIDs)
			if err != nil {
				return err
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			items, _, err := client.BugService.GetBugRelatedStories(cmd.Context(), &tapd.GetBugRelatedStoriesRequest{
				WorkspaceID: new(workspaceID),
				BugID:       ids,
			})
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(items))
			for _, item := range items {
				rows = append(rows, []string{item.WorkspaceID, item.BugID, item.StoryID})
			}

			return writeOutput(cmd, rt.OutputFormat, []string{"WorkspaceID", "BugID", "StoryID"}, rows, items)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().StringVar(&bugIDs, "bug-ids", "", "comma separated bug IDs")
	_ = cmd.MarkFlagRequired("bug-ids")
	return cmd
}

func newBugLinkCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID  int
		bugID        int64
		relateBugIDs string
	)

	cmd := &cobra.Command{
		Use:   "link",
		Short: "Link bugs together",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ids, err := strictInt64Multi("related bug IDs", relateBugIDs)
			if err != nil {
				return err
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			ok, _, err := client.BugService.LinkBugs(cmd.Context(), &tapd.LinkBugsRequest{
				WorkspaceID: new(workspaceID),
				BugID:       new(bugID),
				RelateBugs:  ids,
			})
			if err != nil {
				return err
			}

			rows := [][]string{{strconv.FormatBool(ok)}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Success"}, rows, map[string]bool{"success": ok})
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().Int64Var(&bugID, "bug-id", 0, "source bug ID")
	cmd.Flags().StringVar(&relateBugIDs, "relate-bug-ids", "", "comma separated related bug IDs")
	_ = cmd.MarkFlagRequired("bug-id")
	_ = cmd.MarkFlagRequired("relate-bug-ids")
	return cmd
}

func newBugUnlinkCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		bugID       int64
		linkIDs     string
	)

	cmd := &cobra.Command{
		Use:   "unlink",
		Short: "Remove bug links",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ids, err := strictInt64Multi("link IDs", linkIDs)
			if err != nil {
				return err
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			ok, _, err := client.BugService.DeleteLinkBugs(cmd.Context(), &tapd.DeleteLinkBugsRequest{
				WorkspaceID: new(workspaceID),
				BugID:       new(bugID),
				LinkIDs:     ids,
			})
			if err != nil {
				return err
			}

			rows := [][]string{{strconv.FormatBool(ok)}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Success"}, rows, map[string]bool{"success": ok})
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().Int64Var(&bugID, "bug-id", 0, "bug ID")
	cmd.Flags().StringVar(&linkIDs, "link-ids", "", "comma separated bug link relation IDs")
	_ = cmd.MarkFlagRequired("bug-id")
	_ = cmd.MarkFlagRequired("link-ids")
	return cmd
}

func newBugByViewCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		limit       int
		page        int
		fields      string
		viewConfID  int64
		currentUser string
		creator     string
		owner       string
		status      string
	)

	cmd := &cobra.Command{
		Use:   "by-view",
		Short: "List bugs by view configuration",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetBugsByViewConfIDRequest{
				ViewConfID: new(viewConfID),
				GetBugsRequest: tapd.GetBugsRequest{
					WorkspaceID: new(workspaceID),
					Limit:       new(limit),
					Page:        new(page),
					Fields:      fieldsMulti(fields),
				},
			}
			if currentUser != "" {
				request.CurrentUser = new(currentUser)
			}
			if creator != "" {
				request.Reporter = tapd.NewMulti(creator)
			}
			if owner != "" {
				request.CurrentOwner = new(owner)
			}
			if status != "" {
				request.Status = tapd.NewEnum(splitCSV(status)...)
			}

			bugs, _, err := client.BugService.GetBugsByViewConfID(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, bugTableHeaders(), bugRows(bugs), bugs)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	newListFlags(cmd, &limit, &page)
	cmd.Flags().StringVar(&fields, "fields", "", "comma separated fields")
	cmd.Flags().Int64Var(&viewConfID, "view-conf-id", 0, "view configuration ID")
	cmd.Flags().StringVar(&currentUser, "current-user", "", "current user for personal views")
	cmd.Flags().StringVar(&creator, "creator", "", "filter by creator")
	cmd.Flags().StringVar(&owner, "owner", "", "filter by current owner")
	cmd.Flags().StringVar(&status, "status", "", "filter by status")
	_ = cmd.MarkFlagRequired("view-conf-id")
	return cmd
}

func newBugConvertIDsCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		ids         string
	)

	cmd := &cobra.Command{
		Use:   "convert-ids",
		Short: "Convert bug IDs to a list query token",
		RunE: func(cmd *cobra.Command, _ []string) error {
			bugIDs, err := strictInt64Multi("bug IDs", ids)
			if err != nil {
				return err
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			result, _, err := client.BugService.GetConvertBugIDsToQueryToken(
				cmd.Context(),
				&tapd.GetConvertBugIDsToQueryTokenRequest{
					WorkspaceID: new(workspaceID),
					BugIDs:      bugIDs,
				},
			)
			if err != nil {
				return err
			}

			rows := [][]string{{result.QueryToken, result.Href}}
			return writeOutput(cmd, rt.OutputFormat, []string{"QueryToken", "Href"}, rows, result)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().StringVar(&ids, "ids", "", "comma separated bug IDs")
	_ = cmd.MarkFlagRequired("ids")
	return cmd
}

type bugMutationFlags struct {
	workspaceID   int
	title         string
	description   string
	priority      string
	priorityLabel string
	severity      string
	module        string
	currentOwner  string
	cc            string
	reporter      string
	participator  string
	te            string
	de            string
	auditer       string
	confirmer     string
	fixer         string
	closer        string
	status        string
	vStatus       string
	versionReport string
	versionTest   string
	versionFix    string
	versionClose  string
	iterationID   int64
	releaseID     int64
	source        string
	bugType       string
	label         string
	deadline      string
	feature       string
	frequency     string
	originPhase   string
	sourcePhase   string
	resolution    string
	estimate      int
}

func addBugMutationFlags(cmd *cobra.Command, flags *bugMutationFlags, update bool) {
	if update {
		cmd.Flags().StringVar(&flags.title, "title", "", "bug title")
		cmd.Flags().StringVar(&flags.status, "status", "", "bug status")
		cmd.Flags().StringVar(&flags.vStatus, "v-status", "", "localized bug status")
		cmd.Flags().StringVar(&flags.participator, "participator", "", "bug participators, comma separated")
		cmd.Flags().StringVar(&flags.auditer, "auditer", "", "bug auditer")
		cmd.Flags().StringVar(&flags.confirmer, "confirmer", "", "bug confirmer")
		cmd.Flags().StringVar(&flags.fixer, "fixer", "", "bug fixer")
		cmd.Flags().StringVar(&flags.closer, "closer", "", "bug closer")
		cmd.Flags().StringVar(&flags.versionTest, "version-test", "", "test version")
		cmd.Flags().StringVar(&flags.versionFix, "version-fix", "", "fixed version")
		cmd.Flags().StringVar(&flags.versionClose, "version-close", "", "closed version")
		cmd.Flags().StringVar(&flags.feature, "feature", "", "bug feature")
		cmd.Flags().StringVar(&flags.frequency, "frequency", "", "bug recurrence frequency")
		cmd.Flags().StringVar(&flags.originPhase, "origin-phase", "", "bug origin phase")
		cmd.Flags().StringVar(&flags.sourcePhase, "source-phase", "", "bug source phase")
		cmd.Flags().StringVar(&flags.resolution, "resolution", "", "bug resolution")
		cmd.Flags().IntVar(&flags.estimate, "estimate", 0, "estimated fix time")
	}
	cmd.Flags().StringVar(&flags.description, "description", "", "bug description")
	cmd.Flags().StringVar(&flags.priority, "priority", "", "bug priority")
	cmd.Flags().StringVar(&flags.priorityLabel, "priority-label", "", "bug priority label")
	cmd.Flags().StringVar(&flags.severity, "severity", "", "bug severity: fatal, serious, normal, prompt, or advice")
	cmd.Flags().StringVar(&flags.module, "module", "", "bug module")
	cmd.Flags().StringVar(&flags.currentOwner, "owner", "", "bug current owner")
	cmd.Flags().StringVar(&flags.cc, "cc", "", "bug CC users")
	cmd.Flags().StringVar(&flags.reporter, "reporter", "", "bug reporter")
	cmd.Flags().StringVar(&flags.te, "te", "", "bug tester")
	cmd.Flags().StringVar(&flags.de, "de", "", "bug developer")
	cmd.Flags().StringVar(&flags.versionReport, "version-report", "", "reported version")
	cmd.Flags().Int64Var(&flags.iterationID, "iteration-id", 0, "iteration ID")
	cmd.Flags().Int64Var(&flags.releaseID, "release-id", 0, "release ID")
	cmd.Flags().StringVar(&flags.source, "source", "", "bug source")
	cmd.Flags().StringVar(&flags.bugType, "type", "", "bug type")
	cmd.Flags().StringVar(&flags.label, "label", "", "bug labels")
	cmd.Flags().StringVar(&flags.deadline, "deadline", "", "bug deadline")
}

func applyBugCreateFlags(request *tapd.CreateBugRequest, flags bugMutationFlags) {
	if flags.description != "" {
		request.Description = new(flags.description)
	}
	if flags.priority != "" {
		request.Priority = new(flags.priority)
	}
	if flags.priorityLabel != "" {
		request.PriorityLabel = new(tapd.PriorityLabel(flags.priorityLabel))
	}
	if flags.severity != "" {
		request.Severity = new(tapd.BugSeverity(flags.severity))
	}
	if flags.module != "" {
		request.Module = new(flags.module)
	}
	if flags.currentOwner != "" {
		request.CurrentOwner = new(flags.currentOwner)
	}
	if flags.cc != "" {
		request.CC = new(flags.cc)
	}
	if flags.reporter != "" {
		request.Reporter = new(flags.reporter)
	}
	if flags.te != "" {
		request.TE = new(flags.te)
	}
	if flags.de != "" {
		request.DE = new(flags.de)
	}
	if flags.versionReport != "" {
		request.VersionReport = new(flags.versionReport)
	}
	if flags.iterationID > 0 {
		request.IterationID = new(flags.iterationID)
	}
	if flags.releaseID > 0 {
		request.ReleaseID = new(flags.releaseID)
	}
	if flags.source != "" {
		request.Source = new(flags.source)
	}
	if flags.bugType != "" {
		request.BugType = new(flags.bugType)
	}
	if flags.label != "" {
		request.Label = new(flags.label)
	}
	if flags.deadline != "" {
		request.Deadline = new(flags.deadline)
	}
}

func applyBugUpdateFlags(request *tapd.UpdateBugRequest, flags bugMutationFlags) { //nolint:gocyclo
	if flags.title != "" {
		request.Title = new(flags.title)
	}
	if flags.description != "" {
		request.Description = new(flags.description)
	}
	if flags.priority != "" {
		request.Priority = new(flags.priority)
	}
	if flags.priorityLabel != "" {
		request.PriorityLabel = new(tapd.PriorityLabel(flags.priorityLabel))
	}
	if flags.severity != "" {
		request.Severity = tapd.NewEnum(tapd.BugSeverity(flags.severity))
	}
	if flags.status != "" {
		request.Status = stringEnum(flags.status)
	}
	if flags.vStatus != "" {
		request.VStatus = new(flags.vStatus)
	}
	if flags.label != "" {
		request.Label = stringEnum(flags.label)
	}
	if flags.iterationID > 0 {
		request.IterationID = tapd.NewEnum(strconv.FormatInt(flags.iterationID, 10))
	}
	if flags.module != "" {
		request.Module = stringEnum(flags.module)
	}
	if flags.releaseID > 0 {
		request.ReleaseID = new(int(flags.releaseID))
	}
	if flags.versionReport != "" {
		request.VersionReport = stringEnum(flags.versionReport)
	}
	if flags.versionTest != "" {
		request.VersionTest = new(flags.versionTest)
	}
	if flags.versionFix != "" {
		request.VersionFix = new(flags.versionFix)
	}
	if flags.versionClose != "" {
		request.VersionClose = new(flags.versionClose)
	}
	if flags.feature != "" {
		request.Feature = new(flags.feature)
	}
	if flags.currentOwner != "" {
		request.CurrentOwner = new(flags.currentOwner)
	}
	if flags.cc != "" {
		request.CC = new(flags.cc)
	}
	if flags.reporter != "" {
		request.Reporter = tapd.NewMulti(splitCSV(flags.reporter)...)
	}
	if flags.participator != "" {
		request.Participator = tapd.NewMulti(splitCSV(flags.participator)...)
	}
	if flags.te != "" {
		request.TE = new(flags.te)
	}
	if flags.de != "" {
		request.DE = new(flags.de)
	}
	if flags.auditer != "" {
		request.Auditer = new(flags.auditer)
	}
	if flags.confirmer != "" {
		request.Confirmer = new(flags.confirmer)
	}
	if flags.fixer != "" {
		request.Fixer = new(flags.fixer)
	}
	if flags.closer != "" {
		request.Closer = new(flags.closer)
	}
	if flags.deadline != "" {
		request.Deadline = new(flags.deadline)
	}
	if flags.source != "" {
		request.Source = stringEnum(flags.source)
	}
	if flags.bugType != "" {
		request.BugType = new(flags.bugType)
	}
	if flags.frequency != "" {
		request.Frequency = stringEnum(flags.frequency)
	}
	if flags.originPhase != "" {
		request.OriginPhase = new(flags.originPhase)
	}
	if flags.sourcePhase != "" {
		request.SourcePhase = new(flags.sourcePhase)
	}
	if flags.resolution != "" {
		request.Resolution = stringEnum(flags.resolution)
	}
	if flags.estimate > 0 {
		request.Estimate = new(flags.estimate)
	}
}

func decodeBugBatchUpdate(workspaceID int, data []byte) (*tapd.BatchUpdateBugsRequest, error) {
	var request tapd.BatchUpdateBugsRequest
	if err := json.Unmarshal(data, &request); err == nil && len(request.Workitems) > 0 {
		request.ProjectID = new(workspaceID)
		for _, item := range request.Workitems {
			if item != nil && item.WorkspaceID == nil {
				item.WorkspaceID = new(workspaceID)
			}
		}
		return &request, nil
	}

	var items []*tapd.UpdateBugRequest
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("decode batch update file as object or array: %w", err)
	}
	if len(items) == 0 {
		return nil, errors.New("batch update file must contain at least one bug update")
	}
	for _, item := range items {
		if item != nil && item.WorkspaceID == nil {
			item.WorkspaceID = new(workspaceID)
		}
	}
	return &tapd.BatchUpdateBugsRequest{
		ProjectID: new(workspaceID),
		Workitems: items,
	}, nil
}

type bugChangesFlags struct {
	workspaceID   int
	limit         int
	page          int
	fields        string
	ids           string
	bugIDs        string
	author        string
	field         string
	oldValue      string
	newValue      string
	memo          string
	created       string
	includeAddBug bool
}

func addBugChangesFlags(cmd *cobra.Command, flags *bugChangesFlags, withPaging bool) {
	newWorkspaceFlag(cmd, &flags.workspaceID)
	if withPaging {
		newListFlags(cmd, &flags.limit, &flags.page)
		cmd.Flags().StringVar(&flags.fields, "fields", "", "comma separated fields")
		cmd.Flags().BoolVar(&flags.includeAddBug, "include-add-bug", false, "include bug creation records")
	}
	cmd.Flags().StringVar(&flags.ids, "ids", "", "comma separated change IDs")
	cmd.Flags().StringVar(&flags.bugIDs, "bug-ids", "", "comma separated bug IDs")
	cmd.Flags().StringVar(&flags.author, "author", "", "filter by author")
	cmd.Flags().StringVar(&flags.field, "field", "", "filter by changed field")
	cmd.Flags().StringVar(&flags.oldValue, "old-value", "", "filter by old value")
	cmd.Flags().StringVar(&flags.newValue, "new-value", "", "filter by new value")
	cmd.Flags().StringVar(&flags.memo, "memo", "", "filter by memo")
	cmd.Flags().StringVar(&flags.created, "created", "", "filter by created time expression")
}

func newBugChangesRequest(flags bugChangesFlags) (*tapd.GetBugChangesRequest, error) {
	request := &tapd.GetBugChangesRequest{
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
	if flags.bugIDs != "" {
		ids, err := strictInt64Multi("bug IDs", flags.bugIDs)
		if err != nil {
			return nil, err
		}
		request.BugID = ids
	}
	if flags.author != "" {
		request.Author = new(flags.author)
	}
	if flags.field != "" {
		request.Field = new(flags.field)
	}
	if flags.oldValue != "" {
		request.OldValue = new(flags.oldValue)
	}
	if flags.newValue != "" {
		request.NewValue = new(flags.newValue)
	}
	if flags.memo != "" {
		request.Memo = new(flags.memo)
	}
	if flags.created != "" {
		request.Created = new(flags.created)
	}
	if flags.includeAddBug {
		request.IncludeAddBug = new(1)
	}
	return request, nil
}

func newBugChangesCountRequest(flags bugChangesFlags) (*tapd.GetBugChangesCountRequest, error) {
	request := &tapd.GetBugChangesCountRequest{
		WorkspaceID: new(flags.workspaceID),
	}
	if flags.ids != "" {
		ids, err := strictInt64Multi("change IDs", flags.ids)
		if err != nil {
			return nil, err
		}
		request.ID = ids
	}
	if flags.bugIDs != "" {
		ids, err := strictInt64Multi("bug IDs", flags.bugIDs)
		if err != nil {
			return nil, err
		}
		request.BugID = ids
	}
	if flags.author != "" {
		request.Author = new(flags.author)
	}
	if flags.field != "" {
		request.Field = new(flags.field)
	}
	if flags.oldValue != "" {
		request.OldValue = new(flags.oldValue)
	}
	if flags.newValue != "" {
		request.NewValue = new(flags.newValue)
	}
	if flags.memo != "" {
		request.Memo = new(flags.memo)
	}
	if flags.created != "" {
		request.Created = new(flags.created)
	}
	return request, nil
}

func stringEnum(csv string) *tapd.Enum[string] {
	items := splitCSV(csv)
	if len(items) == 0 {
		return nil
	}
	return tapd.NewEnum(items...)
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func bugTableHeaders() []string {
	return []string{"ID", "Title", "Status", "Owner", "Reporter", "Modified"}
}

func bugRows(bugs []*tapd.Bug) [][]string {
	rows := make([][]string, 0, len(bugs))
	for _, item := range bugs {
		rows = append(rows, []string{
			item.ID,
			item.Title,
			item.Status,
			item.CurrentOwner,
			item.Reporter,
			item.Modified,
		})
	}
	return rows
}
