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

func newStoryCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "story",
		Short: "Work with TAPD stories",
	}

	cmd.AddCommand(
		newStoryCreateCmd(rt),
		newStoryViewCmd(rt),
		newStoryListCmd(rt),
		newStoryCountCmd(rt),
		newStoryUpdateCmd(rt),
		newStoryBatchUpdateCmd(rt),
		newStoryCategoriesCmd(rt),
		newStoryChangesCmd(rt),
		newStoryFieldsCmd(rt),
		newStoryFieldLabelsCmd(rt),
		newStoryTemplatesCmd(rt),
		newStoryTemplateFieldsCmd(rt),
		newStoryRemovedCmd(rt),
		newStoryRelatedBugsCmd(rt),
		newStoryRelatedTestCasesCmd(rt),
		newStoryByViewCmd(rt),
		newStoryConvertIDsCmd(rt),
	)
	return cmd
}

func newStoryListCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		limit       int
		page        int
		fields      string
		ids         string
		creator     string
		owner       string
		status      string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List stories",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetStoriesRequest{
				WorkspaceID: new(workspaceID),
				Limit:       new(limit),
				Page:        new(page),
				Fields:      fieldsMulti(fields),
				ID:          int64Multi(ids),
			}
			if creator != "" {
				request.Creator = new(creator)
			}
			if owner != "" {
				request.Owner = new(owner)
			}
			if status != "" {
				request.VStatus = new(status)
				request.WithVStatus = new("1")
			}

			stories, _, err := client.StoryService.GetStories(cmd.Context(), request)
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(stories))
			for _, item := range stories {
				rows = append(rows, []string{
					item.ID,
					item.Name,
					string(item.Status),
					item.Owner,
					item.Creator,
					item.Modified,
				})
			}

			return writeOutput(cmd, rt.OutputFormat, []string{"ID", "Name", "Status", "Owner", "Creator", "Modified"}, rows, stories)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	newListFlags(cmd, &limit, &page)
	cmd.Flags().StringVar(&fields, "fields", "", "comma separated fields")
	cmd.Flags().StringVar(&ids, "ids", "", "comma separated story IDs")
	cmd.Flags().StringVar(&creator, "creator", "", "filter by creator")
	cmd.Flags().StringVar(&owner, "owner", "", "filter by owner")
	cmd.Flags().StringVar(&status, "status", "", "filter by status or localized status text")

	return cmd
}

func newStoryCreateCmd(rt *app.Runtime) *cobra.Command {
	var flags storyMutationFlags

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a story",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.CreateStoryRequest{
				WorkspaceID: new(flags.workspaceID),
				Name:        new(flags.name),
			}
			applyStoryCreateFlags(request, flags)

			story, _, err := client.StoryService.CreateStory(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, storyTableHeaders(), storyRows([]*tapd.Story{story}), story)
		},
	}

	newWorkspaceFlag(cmd, &flags.workspaceID)
	cmd.Flags().StringVar(&flags.name, "name", "", "story name")
	_ = cmd.MarkFlagRequired("name")
	addStoryMutationFlags(cmd, &flags, false)
	return cmd
}

func newStoryViewCmd(rt *app.Runtime) *cobra.Command {
	var workspaceID int

	cmd := &cobra.Command{
		Use:   "view <id>",
		Short: "Show story details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseInt64Arg("story id", args[0])
			if err != nil {
				return err
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			stories, _, err := client.StoryService.GetStories(cmd.Context(), &tapd.GetStoriesRequest{
				WorkspaceID: new(workspaceID),
				ID:          tapd.NewMulti(id),
				Limit:       new(1),
				Page:        new(1),
			})
			if err != nil {
				return err
			}
			if len(stories) == 0 {
				return fmt.Errorf("story %d not found", id)
			}

			return writeOutput(cmd, rt.OutputFormat, storyTableHeaders(), storyRows(stories[:1]), stories[0])
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	return cmd
}

func newStoryCountCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		ids         string
		creator     string
		owner       string
		status      string
	)

	cmd := &cobra.Command{
		Use:   "count",
		Short: "Count stories",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetStoriesCountRequest{
				WorkspaceID: new(workspaceID),
				ID:          int64Multi(ids),
			}
			if creator != "" {
				request.Creator = new(creator)
			}
			if owner != "" {
				request.Owner = new(owner)
			}
			if status != "" {
				request.VStatus = new(status)
				request.WithVStatus = new("1")
			}

			count, _, err := client.StoryService.GetStoriesCount(cmd.Context(), request)
			if err != nil {
				return err
			}

			data := map[string]any{
				"resource":     "story",
				"workspace_id": workspaceID,
				"count":        count,
			}
			rows := [][]string{{"story", strconv.Itoa(workspaceID), strconv.Itoa(count)}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Resource", "WorkspaceID", "Count"}, rows, data)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().StringVar(&ids, "ids", "", "comma separated story IDs")
	cmd.Flags().StringVar(&creator, "creator", "", "filter by creator")
	cmd.Flags().StringVar(&owner, "owner", "", "filter by owner")
	cmd.Flags().StringVar(&status, "status", "", "filter by status or localized status text")
	return cmd
}

func newStoryFieldsCmd(rt *app.Runtime) *cobra.Command {
	var workspaceID int

	cmd := &cobra.Command{
		Use:   "fields",
		Short: "List story fields",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			fields, _, err := client.StoryService.GetStoryFieldsInfo(cmd.Context(), &tapd.GetStoryFieldsInfoRequest{
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
					strconv.Itoa(len(field.UserOptions)),
				})
			}

			return writeOutput(
				cmd,
				rt.OutputFormat,
				[]string{"Name", "Label", "Type", "Options", "ColorOptions", "PureOptions", "UserOptions"},
				rows,
				fields,
			)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	return cmd
}

func newStoryUpdateCmd(rt *app.Runtime) *cobra.Command {
	var flags storyMutationFlags

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a story",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseInt64Arg("story id", args[0])
			if err != nil {
				return err
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.UpdateStoryRequest{
				ID:          new(id),
				WorkspaceID: new(flags.workspaceID),
			}
			applyStoryUpdateFlags(request, flags)

			story, _, err := client.StoryService.UpdateStory(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, storyTableHeaders(), storyRows([]*tapd.Story{story}), story)
		},
	}

	newWorkspaceFlag(cmd, &flags.workspaceID)
	addStoryMutationFlags(cmd, &flags, true)
	return cmd
}

func newStoryBatchUpdateCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		file        string
	)

	cmd := &cobra.Command{
		Use:   "batch-update",
		Short: "Batch update stories from a JSON file",
		RunE: func(cmd *cobra.Command, _ []string) error {
			data, err := os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("read batch update file: %w", err)
			}

			request, err := decodeStoryBatchUpdate(workspaceID, data)
			if err != nil {
				return err
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			result, _, err := client.StoryService.BatchUpdateStories(cmd.Context(), request)
			if err != nil {
				return err
			}

			rows := [][]string{{result.Msg}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Message"}, rows, result)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().StringVar(&file, "file", "", "JSON file containing an array of story update objects")
	_ = cmd.MarkFlagRequired("file")
	return cmd
}

func newStoryCategoriesCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		limit       int
		page        int
		fields      string
		ids         string
		name        string
		parentID    int
	)

	cmd := &cobra.Command{
		Use:   "categories",
		Short: "List story categories",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetStoryCategoriesRequest{
				WorkspaceID: new(workspaceID),
				Limit:       new(limit),
				Page:        new(page),
				Fields:      fieldsMulti(fields),
				ID:          int64Multi(ids),
			}
			if name != "" {
				request.Name = new(name)
			}
			if parentID > 0 {
				request.ParentID = new(parentID)
			}

			categories, _, err := client.StoryService.GetStoryCategories(cmd.Context(), request)
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(categories))
			for _, item := range categories {
				rows = append(rows, []string{
					item.ID,
					item.Name,
					item.ParentID,
					item.Creator,
					item.Modified,
				})
			}

			return writeOutput(cmd, rt.OutputFormat, []string{"ID", "Name", "ParentID", "Creator", "Modified"}, rows, categories)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	newListFlags(cmd, &limit, &page)
	cmd.Flags().StringVar(&fields, "fields", "", "comma separated fields")
	cmd.Flags().StringVar(&ids, "ids", "", "comma separated category IDs")
	cmd.Flags().StringVar(&name, "name", "", "filter by category name")
	cmd.Flags().IntVar(&parentID, "parent-id", 0, "filter by parent category ID")
	cmd.AddCommand(newStoryCategoriesCountCmd(rt))
	return cmd
}

func newStoryCategoriesCountCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		ids         string
		name        string
		parentID    int
	)

	cmd := &cobra.Command{
		Use:   "count",
		Short: "Count story categories",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetStoryCategoriesCountRequest{
				WorkspaceID: new(workspaceID),
				ID:          int64Multi(ids),
			}
			if name != "" {
				request.Name = new(name)
			}
			if parentID > 0 {
				request.ParentID = new(parentID)
			}

			count, _, err := client.StoryService.GetStoryCategoriesCount(cmd.Context(), request)
			if err != nil {
				return err
			}

			data := map[string]any{"resource": "story_category", "workspace_id": workspaceID, "count": count}
			rows := [][]string{{"story_category", strconv.Itoa(workspaceID), strconv.Itoa(count)}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Resource", "WorkspaceID", "Count"}, rows, data)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().StringVar(&ids, "ids", "", "comma separated category IDs")
	cmd.Flags().StringVar(&name, "name", "", "filter by category name")
	cmd.Flags().IntVar(&parentID, "parent-id", 0, "filter by parent category ID")
	return cmd
}

func newStoryChangesCmd(rt *app.Runtime) *cobra.Command {
	var flags storyChangesFlags

	cmd := &cobra.Command{
		Use:   "changes",
		Short: "List story changes",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request, err := newStoryChangesRequest(flags)
			if err != nil {
				return err
			}

			changes, _, err := client.StoryService.GetStoryChanges(cmd.Context(), request)
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(changes))
			for _, item := range changes {
				rows = append(rows, []string{
					item.ID,
					item.StoryID,
					string(item.ChangeType),
					item.ChangeSummary,
					item.Creator,
					item.Created,
					strconv.Itoa(len(item.FieldChanges)),
				})
			}

			return writeOutput(
				cmd,
				rt.OutputFormat,
				[]string{"ID", "StoryID", "Type", "Summary", "Creator", "Created", "FieldChanges"},
				rows,
				changes,
			)
		},
	}

	addStoryChangesFlags(cmd, &flags, true)
	cmd.AddCommand(newStoryChangesCountCmd(rt))
	return cmd
}

func newStoryChangesCountCmd(rt *app.Runtime) *cobra.Command {
	var flags storyChangesFlags

	cmd := &cobra.Command{
		Use:   "count",
		Short: "Count story changes",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request, err := newStoryChangesCountRequest(flags)
			if err != nil {
				return err
			}

			count, _, err := client.StoryService.GetStoryChangesCount(cmd.Context(), request)
			if err != nil {
				return err
			}

			data := map[string]any{"resource": "story_change", "workspace_id": flags.workspaceID, "count": count}
			rows := [][]string{{"story_change", strconv.Itoa(flags.workspaceID), strconv.Itoa(count)}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Resource", "WorkspaceID", "Count"}, rows, data)
		},
	}

	addStoryChangesFlags(cmd, &flags, false)
	return cmd
}

func newStoryFieldLabelsCmd(rt *app.Runtime) *cobra.Command {
	var workspaceID int

	cmd := &cobra.Command{
		Use:   "field-labels",
		Short: "List story field labels",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			labels, _, err := client.StoryService.GetStoryFieldsLabel(cmd.Context(), &tapd.GetStoryFieldsLabelRequest{
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

func newStoryTemplatesCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID    int
		workitemTypeID int
	)

	cmd := &cobra.Command{
		Use:   "templates",
		Short: "List story templates",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetStoryTemplatesRequest{
				WorkspaceID: new(workspaceID),
			}
			if workitemTypeID > 0 {
				request.WorkitemTypeID = new(workitemTypeID)
			}

			templates, _, err := client.StoryService.GetStoryTemplates(cmd.Context(), request)
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
	cmd.Flags().IntVar(&workitemTypeID, "workitem-type-id", 0, "filter by story workitem type ID")
	return cmd
}

func newStoryTemplateFieldsCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		templateID  int64
	)

	cmd := &cobra.Command{
		Use:   "template-fields",
		Short: "List story template fields",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			fields, _, err := client.StoryService.GetStoryTemplateFields(cmd.Context(), &tapd.GetStoryTemplateFieldsRequest{
				WorkspaceID: new(workspaceID),
				TemplateID:  new(templateID),
			})
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
	cmd.Flags().Int64Var(&templateID, "template-id", 0, "story template ID")
	_ = cmd.MarkFlagRequired("template-id")
	return cmd
}

func newStoryRemovedCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		limit       int
		page        int
		ids         string
		creator     string
		archived    bool
	)

	cmd := &cobra.Command{
		Use:   "removed",
		Short: "List removed stories",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetRemovedStoriesRequest{
				WorkspaceID: new(workspaceID),
				Limit:       new(limit),
				Page:        new(page),
			}
			if ids != "" {
				request.ID, err = strictIntMulti("story IDs", ids)
				if err != nil {
					return err
				}
			}
			if creator != "" {
				request.Creator = new(creator)
			}
			if archived {
				request.IsArchived = new(1)
			}

			stories, _, err := client.StoryService.GetRemovedStories(cmd.Context(), request)
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(stories))
			for _, item := range stories {
				rows = append(rows, []string{item.ID, item.Name, item.Creator, item.OperationUser, item.IsArchived, item.Deleted})
			}

			return writeOutput(cmd, rt.OutputFormat, []string{"ID", "Name", "Creator", "OperationUser", "Archived", "Deleted"}, rows, stories)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	newListFlags(cmd, &limit, &page)
	cmd.Flags().StringVar(&ids, "ids", "", "comma separated story IDs")
	cmd.Flags().StringVar(&creator, "creator", "", "filter by creator")
	cmd.Flags().BoolVar(&archived, "archived", false, "only return archived stories")
	return cmd
}

func newStoryRelatedBugsCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		storyIDs    string
	)

	cmd := &cobra.Command{
		Use:   "related-bugs",
		Short: "List bugs related to stories",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ids, err := strictInt64Multi("story IDs", storyIDs)
			if err != nil {
				return err
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			items, _, err := client.StoryService.GetStoryRelatedBugs(cmd.Context(), &tapd.GetStoryRelatedBugsRequest{
				WorkspaceID: new(workspaceID),
				StoryID:     ids,
			})
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(items))
			for _, item := range items {
				rows = append(rows, []string{strconv.Itoa(item.WorkspaceID), item.StoryID, item.BugID})
			}

			return writeOutput(cmd, rt.OutputFormat, []string{"WorkspaceID", "StoryID", "BugID"}, rows, items)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().StringVar(&storyIDs, "story-ids", "", "comma separated story IDs")
	_ = cmd.MarkFlagRequired("story-ids")
	return cmd
}

func newStoryRelatedTestCasesCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID     int
		storyID         int64
		includeTestPlan bool
	)

	cmd := &cobra.Command{
		Use:   "related-test-cases",
		Short: "List test cases related to a story",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			include := 0
			if includeTestPlan {
				include = 1
			}
			items, _, err := client.StoryService.GetStoryTestCaseRelation(cmd.Context(), &tapd.GetStoryTestCaseRelationRequest{
				WorkspaceID:     new(workspaceID),
				StoryID:         new(storyID),
				IncludeTestPlan: new(include),
			})
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(items))
			for _, item := range items {
				rows = append(rows, []string{item.ID, item.StoryID, item.TcaseID, item.TestPlanID, item.Creator, item.Created})
			}

			return writeOutput(cmd, rt.OutputFormat, []string{"ID", "StoryID", "TestCaseID", "TestPlanID", "Creator", "Created"}, rows, items)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().Int64Var(&storyID, "story-id", 0, "story ID")
	cmd.Flags().BoolVar(&includeTestPlan, "include-test-plan", true, "include test plan relations")
	_ = cmd.MarkFlagRequired("story-id")
	return cmd
}

func newStoryByViewCmd(rt *app.Runtime) *cobra.Command {
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
		Short: "List stories by view configuration",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetStoriesByViewConfIDRequest{
				ViewConfID: new(viewConfID),
				GetStoriesRequest: tapd.GetStoriesRequest{
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
				request.Creator = new(creator)
			}
			if owner != "" {
				request.Owner = new(owner)
			}
			if status != "" {
				request.VStatus = new(status)
				request.WithVStatus = new("1")
			}

			stories, _, err := client.StoryService.GetStoriesByViewConfID(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, storyTableHeaders(), storyRows(stories), stories)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	newListFlags(cmd, &limit, &page)
	cmd.Flags().StringVar(&fields, "fields", "", "comma separated fields")
	cmd.Flags().Int64Var(&viewConfID, "view-conf-id", 0, "view configuration ID")
	cmd.Flags().StringVar(&currentUser, "current-user", "", "current user for personal views")
	cmd.Flags().StringVar(&creator, "creator", "", "filter by creator")
	cmd.Flags().StringVar(&owner, "owner", "", "filter by owner")
	cmd.Flags().StringVar(&status, "status", "", "filter by status or localized status text")
	_ = cmd.MarkFlagRequired("view-conf-id")
	return cmd
}

func newStoryConvertIDsCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		ids         string
	)

	cmd := &cobra.Command{
		Use:   "convert-ids",
		Short: "Convert story IDs to a list query token",
		RunE: func(cmd *cobra.Command, _ []string) error {
			storyIDs, err := strictInt64Multi("story IDs", ids)
			if err != nil {
				return err
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			result, _, err := client.StoryService.GetConvertStoryIDsToQueryToken(
				cmd.Context(),
				&tapd.GetConvertStoryIDsToQueryTokenRequest{
					WorkspaceID: new(workspaceID),
					StoryIDs:    storyIDs,
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
	cmd.Flags().StringVar(&ids, "ids", "", "comma separated story IDs")
	_ = cmd.MarkFlagRequired("ids")
	return cmd
}

type storyMutationFlags struct {
	workspaceID     int
	name            string
	description     string
	owner           string
	creator         string
	currentUser     string
	status          string
	vStatus         string
	priority        string
	priorityLabel   string
	categoryID      int64
	workitemTypeID  int
	iterationID     string
	parentID        int64
	releaseID       int
	begin           string
	due             string
	label           string
	module          string
	version         string
	source          string
	storyType       string
	developer       string
	testFocus       string
	effort          string
	effortCompleted string
	size            int
	businessValue   int
	autoCloseTask   bool
}

func addStoryMutationFlags(cmd *cobra.Command, flags *storyMutationFlags, update bool) {
	if update {
		cmd.Flags().StringVar(&flags.name, "name", "", "story name")
		cmd.Flags().StringVar(&flags.currentUser, "current-user", "", "current operator")
		cmd.Flags().StringVar(&flags.status, "status", "", "story status")
		cmd.Flags().StringVar(&flags.vStatus, "v-status", "", "localized story status")
		cmd.Flags().BoolVar(&flags.autoCloseTask, "auto-close-task", false, "auto close related tasks")
	}
	cmd.Flags().StringVar(&flags.description, "description", "", "story description")
	cmd.Flags().StringVar(&flags.owner, "owner", "", "story owner")
	if !update {
		cmd.Flags().StringVar(&flags.creator, "creator", "", "story creator")
	}
	cmd.Flags().StringVar(&flags.priority, "priority", "", "story priority")
	cmd.Flags().StringVar(&flags.priorityLabel, "priority-label", "", "story priority label")
	cmd.Flags().IntVar(&flags.businessValue, "business-value", 0, "business value")
	cmd.Flags().IntVar(&flags.size, "size", 0, "story size")
	cmd.Flags().Int64Var(&flags.categoryID, "category-id", 0, "story category ID")
	if !update {
		cmd.Flags().IntVar(&flags.workitemTypeID, "workitem-type-id", 0, "story workitem type ID")
	}
	cmd.Flags().StringVar(&flags.iterationID, "iteration-id", "", "iteration ID")
	if !update {
		cmd.Flags().Int64Var(&flags.parentID, "parent-id", 0, "parent story ID")
	}
	cmd.Flags().IntVar(&flags.releaseID, "release-id", 0, "release ID")
	cmd.Flags().StringVar(&flags.begin, "begin", "", "expected begin date")
	cmd.Flags().StringVar(&flags.due, "due", "", "expected due date")
	cmd.Flags().StringVar(&flags.label, "label", "", "story labels")
	cmd.Flags().StringVar(&flags.module, "module", "", "story module")
	cmd.Flags().StringVar(&flags.version, "version", "", "story version")
	cmd.Flags().StringVar(&flags.source, "source", "", "story source")
	cmd.Flags().StringVar(&flags.storyType, "type", "", "story type")
	cmd.Flags().StringVar(&flags.developer, "developer", "", "story developer")
	cmd.Flags().StringVar(&flags.testFocus, "test-focus", "", "story test focus")
	cmd.Flags().StringVar(&flags.effort, "effort", "", "estimated effort")
	cmd.Flags().StringVar(&flags.effortCompleted, "effort-completed", "", "completed effort")
}

func applyStoryCreateFlags(request *tapd.CreateStoryRequest, flags storyMutationFlags) {
	if flags.description != "" {
		request.Description = new(flags.description)
	}
	if flags.owner != "" {
		request.Owner = new(flags.owner)
	}
	if flags.creator != "" {
		request.Creator = new(flags.creator)
	}
	if flags.priority != "" {
		request.Priority = new(flags.priority)
	}
	if flags.priorityLabel != "" {
		request.PriorityLabel = new(tapd.PriorityLabel(flags.priorityLabel))
	}
	if flags.businessValue > 0 {
		request.BusinessValue = new(flags.businessValue)
	}
	if flags.size > 0 {
		request.Size = new(flags.size)
	}
	if flags.categoryID > 0 {
		request.CategoryID = new(int(flags.categoryID))
	}
	if flags.workitemTypeID > 0 {
		request.WorkitemTypeID = new(flags.workitemTypeID)
	}
	if flags.iterationID != "" {
		request.IterationID = new(flags.iterationID)
	}
	if flags.parentID > 0 {
		request.ParentID = new(int(flags.parentID))
	}
	if flags.releaseID > 0 {
		request.ReleaseID = new(flags.releaseID)
	}
	if flags.begin != "" {
		request.Begin = new(flags.begin)
	}
	if flags.due != "" {
		request.Due = new(flags.due)
	}
	if flags.label != "" {
		request.Label = new(flags.label)
	}
	if flags.module != "" {
		request.Module = new(flags.module)
	}
	if flags.version != "" {
		request.Version = new(flags.version)
	}
	if flags.source != "" {
		request.Source = new(flags.source)
	}
	if flags.storyType != "" {
		request.Type = new(flags.storyType)
	}
	if flags.developer != "" {
		request.Developer = new(flags.developer)
	}
	if flags.testFocus != "" {
		request.TestFocus = new(flags.testFocus)
	}
	if flags.effort != "" {
		request.Effort = new(flags.effort)
	}
	if flags.effortCompleted != "" {
		request.EffortCompleted = new(flags.effortCompleted)
	}
}

func applyStoryUpdateFlags(request *tapd.UpdateStoryRequest, flags storyMutationFlags) {
	if flags.name != "" {
		request.Name = new(flags.name)
	}
	if flags.description != "" {
		request.Description = new(flags.description)
	}
	if flags.owner != "" {
		request.Owner = new(flags.owner)
	}
	if flags.currentUser != "" {
		request.CurrentUser = new(flags.currentUser)
	}
	if flags.status != "" {
		request.Status = new(flags.status)
	}
	if flags.vStatus != "" {
		request.VStatus = new(flags.vStatus)
	}
	if flags.priority != "" {
		request.Priority = new(flags.priority)
	}
	if flags.priorityLabel != "" {
		request.PriorityLabel = new(tapd.PriorityLabel(flags.priorityLabel))
	}
	if flags.businessValue > 0 {
		request.BusinessValue = new(flags.businessValue)
	}
	if flags.size > 0 {
		request.Size = new(flags.size)
	}
	if flags.categoryID > 0 {
		request.CategoryID = new(flags.categoryID)
	}
	if flags.iterationID != "" {
		request.IterationID = new(flags.iterationID)
	}
	if flags.releaseID > 0 {
		request.ReleaseID = new(flags.releaseID)
	}
	if flags.begin != "" {
		request.Begin = new(flags.begin)
	}
	if flags.due != "" {
		request.Due = new(flags.due)
	}
	if flags.label != "" {
		request.Label = new(flags.label)
	}
	if flags.module != "" {
		request.Module = new(flags.module)
	}
	if flags.version != "" {
		request.Version = new(flags.version)
	}
	if flags.source != "" {
		request.Source = new(flags.source)
	}
	if flags.storyType != "" {
		request.Type = new(flags.storyType)
	}
	if flags.developer != "" {
		request.Developer = new(flags.developer)
	}
	if flags.testFocus != "" {
		request.TestFocus = new(flags.testFocus)
	}
	if flags.effort != "" {
		request.Effort = new(flags.effort)
	}
	if flags.effortCompleted != "" {
		request.EffortCompleted = new(flags.effortCompleted)
	}
	if flags.autoCloseTask {
		request.IsAutoCloseTask = new(1)
	}
}

func decodeStoryBatchUpdate(workspaceID int, data []byte) (*tapd.BatchUpdateStoriesRequest, error) {
	var request tapd.BatchUpdateStoriesRequest
	if err := json.Unmarshal(data, &request); err == nil && len(request.Workitems) > 0 {
		request.WorkspaceID = new(workspaceID)
		for _, item := range request.Workitems {
			if item != nil && item.WorkspaceID == nil {
				item.WorkspaceID = new(workspaceID)
			}
		}
		return &request, nil
	}

	var items []*tapd.UpdateStoryRequest
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("decode batch update file as object or array: %w", err)
	}
	if len(items) == 0 {
		return nil, errors.New("batch update file must contain at least one story update")
	}
	for _, item := range items {
		if item != nil && item.WorkspaceID == nil {
			item.WorkspaceID = new(workspaceID)
		}
	}
	return &tapd.BatchUpdateStoriesRequest{
		WorkspaceID: new(workspaceID),
		Workitems:   items,
	}, nil
}

type storyChangesFlags struct {
	workspaceID   int
	limit         int
	page          int
	fields        string
	ids           string
	storyIDs      string
	creator       string
	changeType    string
	changeSummary string
	comment       string
	entityType    string
	changeField   string
}

func addStoryChangesFlags(cmd *cobra.Command, flags *storyChangesFlags, withPaging bool) {
	newWorkspaceFlag(cmd, &flags.workspaceID)
	if withPaging {
		newListFlags(cmd, &flags.limit, &flags.page)
		cmd.Flags().StringVar(&flags.fields, "fields", "", "comma separated fields")
	}
	cmd.Flags().StringVar(&flags.ids, "ids", "", "comma separated change IDs")
	cmd.Flags().StringVar(&flags.storyIDs, "story-ids", "", "comma separated story IDs")
	cmd.Flags().StringVar(&flags.creator, "creator", "", "filter by creator")
	cmd.Flags().StringVar(&flags.changeType, "change-type", "", "filter by change type")
	cmd.Flags().StringVar(&flags.changeSummary, "change-summary", "", "filter by change summary")
	cmd.Flags().StringVar(&flags.comment, "comment", "", "filter by comment")
	cmd.Flags().StringVar(&flags.entityType, "entity-type", "", "filter by entity type")
	cmd.Flags().StringVar(&flags.changeField, "change-field", "", "filter by changed field")
}

func newStoryChangesRequest(flags storyChangesFlags) (*tapd.GetStoryChangesRequest, error) {
	request := &tapd.GetStoryChangesRequest{
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
	if flags.storyIDs != "" {
		ids, err := strictInt64Multi("story IDs", flags.storyIDs)
		if err != nil {
			return nil, err
		}
		request.StoryID = ids
	}
	if flags.creator != "" {
		request.Creator = new(flags.creator)
	}
	if flags.changeType != "" {
		request.ChangeType = new(tapd.StoreChangeType(flags.changeType))
	}
	if flags.changeSummary != "" {
		request.ChangeSummary = new(flags.changeSummary)
	}
	if flags.comment != "" {
		request.Comment = new(flags.comment)
	}
	if flags.entityType != "" {
		request.EntityType = new(flags.entityType)
	}
	if flags.changeField != "" {
		request.ChangeField = new(flags.changeField)
	}
	return request, nil
}

func newStoryChangesCountRequest(flags storyChangesFlags) (*tapd.GetStoryChangesCountRequest, error) {
	request := &tapd.GetStoryChangesCountRequest{
		WorkspaceID: new(flags.workspaceID),
	}
	if flags.ids != "" {
		ids, err := strictInt64Multi("change IDs", flags.ids)
		if err != nil {
			return nil, err
		}
		request.ID = ids
	}
	if flags.storyIDs != "" {
		ids, err := strictInt64Multi("story IDs", flags.storyIDs)
		if err != nil {
			return nil, err
		}
		request.StoryID = ids
	}
	if flags.creator != "" {
		request.Creator = new(flags.creator)
	}
	if flags.changeType != "" {
		request.ChangeType = new(tapd.StoreChangeType(flags.changeType))
	}
	if flags.changeSummary != "" {
		request.ChangeSummary = new(flags.changeSummary)
	}
	if flags.comment != "" {
		request.Comment = new(flags.comment)
	}
	if flags.entityType != "" {
		request.EntityType = new(flags.entityType)
	}
	if flags.changeField != "" {
		request.ChangeField = new(flags.changeField)
	}
	return request, nil
}

func storyTableHeaders() []string {
	return []string{"ID", "Name", "Status", "Owner", "Creator", "Modified"}
}

func storyRows(stories []*tapd.Story) [][]string {
	rows := make([][]string, 0, len(stories))
	for _, item := range stories {
		rows = append(rows, []string{
			item.ID,
			item.Name,
			string(item.Status),
			item.Owner,
			item.Creator,
			item.Modified,
		})
	}
	return rows
}
