package cmd

import (
	"strconv"
	"strings"

	"github.com/go-tapd/cli/internal/app"
	"github.com/go-tapd/tapd"
	"github.com/spf13/cobra"
)

func newTestCaseCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test-case",
		Short: "Work with TAPD test cases",
	}

	cmd.AddCommand(
		newTestCaseCreateCmd(rt),
		newTestCaseListCmd(rt),
		newTestCaseCountCmd(rt),
		newTestCaseUpdateCmd(rt),
		newTestCaseCategoriesCmd(rt),
		newTestCaseFieldsCmd(rt),
		newTestCaseResultsCmd(rt),
	)
	return cmd
}

func newTestPlanCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test-plan",
		Short: "Work with TAPD test plans",
	}

	cmd.AddCommand(
		newTestPlanCreateCmd(rt),
		newTestPlanListCmd(rt),
		newTestPlanCountCmd(rt),
		newTestPlanUpdateCmd(rt),
		newTestPlanProgressCmd(rt),
		newTestPlanResultCmd(rt),
		newTestPlanRelatedBugsCmd(rt),
		newTestPlanRelatedStoriesCmd(rt),
	)
	return cmd
}

func newTestCaseCreateCmd(rt *app.Runtime) *cobra.Command {
	var flags testCaseMutationFlags

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a test case",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.CreateTestCaseRequest{
				WorkspaceID: tapd.Ptr(flags.workspaceID),
				Name:        tapd.Ptr(flags.name),
			}
			applyTestCaseMutationFlags(request, flags)

			testCase, _, err := client.TestService.CreateTestCase(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, testCaseTableHeaders(), testCaseRows([]*tapd.TestCase{testCase}), testCase)
		},
	}

	newWorkspaceFlag(cmd, &flags.workspaceID)
	addTestCaseMutationFlags(cmd, &flags, false)
	cmd.Flags().StringVar(&flags.name, "name", "", "test case name")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newTestCaseListCmd(rt *app.Runtime) *cobra.Command {
	var flags testCaseQueryFlags

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List test cases",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request, err := newTestCasesRequest(flags)
			if err != nil {
				return err
			}

			testCases, _, err := client.TestService.GetTestCases(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, testCaseTableHeaders(), testCaseRows(testCases), testCases)
		},
	}

	addTestCaseQueryFlags(cmd, &flags, true)
	return cmd
}

func newTestCaseCountCmd(rt *app.Runtime) *cobra.Command {
	var flags testCaseQueryFlags

	cmd := &cobra.Command{
		Use:   "count",
		Short: "Count test cases",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request, err := newTestCasesCountRequest(flags)
			if err != nil {
				return err
			}

			count, _, err := client.TestService.GetTestCasesCount(cmd.Context(), request)
			if err != nil {
				return err
			}

			data := map[string]any{"resource": "test_case", "workspace_id": flags.workspaceID, "count": count}
			rows := [][]string{{"test_case", strconv.Itoa(flags.workspaceID), strconv.Itoa(count)}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Resource", "WorkspaceID", "Count"}, rows, data)
		},
	}

	addTestCaseQueryFlags(cmd, &flags, false)
	return cmd
}

func newTestCaseUpdateCmd(rt *app.Runtime) *cobra.Command {
	var flags testCaseMutationFlags

	cmd := &cobra.Command{
		Use:   "update <test-case-id>",
		Short: "Update a test case",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseInt64Arg("test case ID", args[0])
			if err != nil {
				return err
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.UpdateTestCaseRequest{
				ID:          tapd.Ptr(id),
				WorkspaceID: tapd.Ptr(flags.workspaceID),
			}
			applyTestCaseMutationFlags(request, flags)

			testCase, _, err := client.TestService.UpdateTestCase(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, testCaseTableHeaders(), testCaseRows([]*tapd.TestCase{testCase}), testCase)
		},
	}

	newWorkspaceFlag(cmd, &flags.workspaceID)
	addTestCaseMutationFlags(cmd, &flags, true)
	return cmd
}

func newTestCaseCategoriesCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		limit       int
		page        int
		fields      string
		ids         string
		name        string
		description string
		parentID    int64
		creator     string
		modifier    string
		created     string
		modified    string
	)

	cmd := &cobra.Command{
		Use:   "categories",
		Short: "List test case categories",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetTestCaseCategoriesRequest{
				WorkspaceID: tapd.Ptr(workspaceID),
				Limit:       tapd.Ptr(limit),
				Page:        tapd.Ptr(page),
				Fields:      fieldsMulti(fields),
			}
			if ids != "" {
				request.ID, err = strictInt64Multi("test case category IDs", ids)
				if err != nil {
					return err
				}
			}
			if name != "" {
				request.Name = tapd.Ptr(name)
			}
			if description != "" {
				request.Description = tapd.Ptr(description)
			}
			if parentID > 0 {
				request.ParentID = tapd.Ptr(parentID)
			}
			if creator != "" {
				request.Creator = tapd.Ptr(creator)
			}
			if modifier != "" {
				request.Modifier = tapd.Ptr(modifier)
			}
			if created != "" {
				request.Created = tapd.Ptr(created)
			}
			if modified != "" {
				request.Modified = tapd.Ptr(modified)
			}

			categories, _, err := client.TestService.GetTestCaseCategories(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, testCaseCategoryTableHeaders(), testCaseCategoryRows(categories), categories)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	newListFlags(cmd, &limit, &page)
	cmd.Flags().StringVar(&fields, "fields", "", "comma separated fields")
	cmd.Flags().StringVar(&ids, "ids", "", "comma separated test case category IDs")
	cmd.Flags().StringVar(&name, "name", "", "filter by category name")
	cmd.Flags().StringVar(&description, "description", "", "filter by category description")
	cmd.Flags().Int64Var(&parentID, "parent-id", 0, "filter by parent category ID")
	cmd.Flags().StringVar(&creator, "creator", "", "filter by creator")
	cmd.Flags().StringVar(&modifier, "modifier", "", "filter by modifier")
	cmd.Flags().StringVar(&created, "created", "", "filter by created time expression")
	cmd.Flags().StringVar(&modified, "modified", "", "filter by modified time expression")
	return cmd
}

func newTestCaseFieldsCmd(rt *app.Runtime) *cobra.Command {
	var workspaceID int

	cmd := &cobra.Command{
		Use:   "fields",
		Short: "List test case fields",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			fields, _, err := client.TestService.GetTestCaseFieldsInfo(cmd.Context(), &tapd.GetTestCaseFieldsInfoRequest{
				WorkspaceID: tapd.Ptr(workspaceID),
			})
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, testFieldsTableHeaders(), testCaseFieldsRows(fields), fields)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	return cmd
}

func newTestCaseResultsCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		testPlanID  int64
		testCaseID  int64
	)

	cmd := &cobra.Command{
		Use:   "results",
		Short: "List test case execution results",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			results, _, err := client.TestService.GetTestCaseResults(cmd.Context(), &tapd.GetTestCaseResultsRequest{
				WorkspaceID: tapd.Ptr(workspaceID),
				TestPlanID:  tapd.Ptr(testPlanID),
				TestCaseID:  tapd.Ptr(testCaseID),
			})
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, testCaseResultTableHeaders(), testCaseResultRows(results), results)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().Int64Var(&testPlanID, "test-plan-id", 0, "test plan ID")
	cmd.Flags().Int64Var(&testCaseID, "test-case-id", 0, "test case ID")
	_ = cmd.MarkFlagRequired("test-plan-id")
	_ = cmd.MarkFlagRequired("test-case-id")
	return cmd
}

func newTestPlanCreateCmd(rt *app.Runtime) *cobra.Command {
	var flags testPlanMutationFlags

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a test plan",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.CreateTestPlanRequest{
				WorkspaceID: tapd.Ptr(flags.workspaceID),
				Name:        tapd.Ptr(flags.name),
			}
			applyTestPlanCreateFlags(request, flags)

			plan, _, err := client.TestService.CreateTestPlan(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, testPlanTableHeaders(), testPlanRows([]*tapd.TestPlan{plan}), plan)
		},
	}

	newWorkspaceFlag(cmd, &flags.workspaceID)
	addTestPlanMutationFlags(cmd, &flags, false)
	cmd.Flags().StringVar(&flags.name, "name", "", "test plan name")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newTestPlanListCmd(rt *app.Runtime) *cobra.Command {
	var flags testPlanQueryFlags

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List test plans",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request, err := newTestPlansRequest(flags)
			if err != nil {
				return err
			}

			plans, _, err := client.TestService.GetTestPlans(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, testPlanTableHeaders(), testPlanRows(plans), plans)
		},
	}

	addTestPlanQueryFlags(cmd, &flags, true)
	return cmd
}

func newTestPlanCountCmd(rt *app.Runtime) *cobra.Command {
	var flags testPlanQueryFlags

	cmd := &cobra.Command{
		Use:   "count",
		Short: "Count test plans",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request, err := newTestPlansCountRequest(flags)
			if err != nil {
				return err
			}

			count, _, err := client.TestService.GetTestPlansCount(cmd.Context(), request)
			if err != nil {
				return err
			}

			data := map[string]any{"resource": "test_plan", "workspace_id": flags.workspaceID, "count": count}
			rows := [][]string{{"test_plan", strconv.Itoa(flags.workspaceID), strconv.Itoa(count)}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Resource", "WorkspaceID", "Count"}, rows, data)
		},
	}

	addTestPlanQueryFlags(cmd, &flags, false)
	return cmd
}

func newTestPlanUpdateCmd(rt *app.Runtime) *cobra.Command {
	var flags testPlanMutationFlags

	cmd := &cobra.Command{
		Use:   "update <test-plan-id>",
		Short: "Update a test plan",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseInt64Arg("test plan ID", args[0])
			if err != nil {
				return err
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.UpdateTestPlanRequest{
				ID:          tapd.Ptr(id),
				WorkspaceID: tapd.Ptr(flags.workspaceID),
			}
			applyTestPlanUpdateFlags(request, flags)

			plan, _, err := client.TestService.UpdateTestPlan(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, testPlanTableHeaders(), testPlanRows([]*tapd.TestPlan{plan}), plan)
		},
	}

	newWorkspaceFlag(cmd, &flags.workspaceID)
	addTestPlanMutationFlags(cmd, &flags, true)
	return cmd
}

func newTestPlanProgressCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		id          int64
	)

	cmd := &cobra.Command{
		Use:   "progress",
		Short: "Show test plan progress",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			progress, _, err := client.TestService.GetTestPlanProgress(cmd.Context(), &tapd.GetTestPlanProgressRequest{
				WorkspaceID: tapd.Ptr(workspaceID),
				ID:          tapd.Ptr(id),
			})
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, testPlanProgressTableHeaders(), testPlanProgressRows(progress), progress)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().Int64Var(&id, "id", 0, "test plan ID")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func newTestPlanResultCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID   int
		id            int64
		lastExecutor  string
		includeRepeat bool
	)

	cmd := &cobra.Command{
		Use:   "result",
		Short: "List test plan result details",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetTestPlanResultRequest{
				WorkspaceID: tapd.Ptr(workspaceID),
				ID:          tapd.Ptr(id),
			}
			if lastExecutor != "" {
				request.LastExecutor = tapd.Ptr(lastExecutor)
			}
			if includeRepeat {
				request.IncludeRepeat = tapd.Ptr(1)
			}

			testCases, _, err := client.TestService.GetTestPlanResult(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, testCaseTableHeaders(), testCaseRows(testCases), testCases)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().Int64Var(&id, "id", 0, "test plan ID")
	cmd.Flags().StringVar(&lastExecutor, "last-executor", "", "filter by last executor")
	cmd.Flags().BoolVar(&includeRepeat, "include-repeat", false, "include repeated test case results")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func newTestPlanRelatedBugsCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		id          int64
		limit       int
		page        int
		fields      string
	)

	cmd := &cobra.Command{
		Use:   "related-bugs",
		Short: "List bugs related to a test plan",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			items, _, err := client.TestService.GetTestPlanRelatedBugs(cmd.Context(), &tapd.GetTestPlanRelatedBugsRequest{
				WorkspaceID: tapd.Ptr(workspaceID),
				ID:          tapd.Ptr(id),
				Limit:       tapd.Ptr(limit),
				Page:        tapd.Ptr(page),
				Fields:      fieldsMulti(fields),
			})
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, testPlanRelatedBugTableHeaders(), testPlanRelatedBugRows(items), items)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().Int64Var(&id, "id", 0, "test plan ID")
	newListFlags(cmd, &limit, &page)
	cmd.Flags().StringVar(&fields, "fields", "", "comma separated fields")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func newTestPlanRelatedStoriesCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		testPlanID  int64
	)

	cmd := &cobra.Command{
		Use:   "related-stories",
		Short: "List stories related to a test plan",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			storyIDs, _, err := client.TestService.GetTestPlanRelatedStories(cmd.Context(), &tapd.GetTestPlanRelatedStoriesRequest{
				WorkspaceID: tapd.Ptr(workspaceID),
				TestPlanID:  tapd.Ptr(testPlanID),
			})
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(storyIDs))
			for _, storyID := range storyIDs {
				rows = append(rows, []string{strconv.Itoa(workspaceID), strconv.FormatInt(testPlanID, 10), storyID})
			}
			data := map[string]any{"workspace_id": workspaceID, "test_plan_id": testPlanID, "story_ids": storyIDs}
			return writeOutput(cmd, rt.OutputFormat, []string{"WorkspaceID", "TestPlanID", "StoryID"}, rows, data)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().Int64Var(&testPlanID, "test-plan-id", 0, "test plan ID")
	_ = cmd.MarkFlagRequired("test-plan-id")
	return cmd
}

type testCaseMutationFlags struct {
	workspaceID  int
	name         string
	steps        string
	categoryID   int64
	status       string
	precondition string
	expectation  string
	testType     string
	priority     string
	creator      string
}

func addTestCaseMutationFlags(cmd *cobra.Command, flags *testCaseMutationFlags, update bool) {
	if update {
		cmd.Flags().StringVar(&flags.name, "name", "", "test case name")
	}
	cmd.Flags().StringVar(&flags.steps, "steps", "", "test case steps")
	cmd.Flags().Int64Var(&flags.categoryID, "category-id", 0, "test case category ID")
	cmd.Flags().StringVar(&flags.status, "status", "", "test case status: normal, updating, or abandon")
	cmd.Flags().StringVar(&flags.precondition, "precondition", "", "test case precondition")
	cmd.Flags().StringVar(&flags.expectation, "expectation", "", "test case expectation")
	cmd.Flags().StringVar(&flags.testType, "type", "", "test case type")
	cmd.Flags().StringVar(&flags.priority, "priority", "", "test case priority")
	cmd.Flags().StringVar(&flags.creator, "creator", "", "test case creator")
}

func applyTestCaseMutationFlags(request *tapd.CreateTestCaseRequest, flags testCaseMutationFlags) {
	if flags.name != "" {
		request.Name = tapd.Ptr(flags.name)
	}
	if flags.steps != "" {
		request.Steps = tapd.Ptr(flags.steps)
	}
	if flags.categoryID > 0 {
		request.CategoryID = tapd.Ptr(flags.categoryID)
	}
	if flags.status != "" {
		request.Status = tapd.Ptr(tapd.TestCaseStatus(flags.status))
	}
	if flags.precondition != "" {
		request.Precondition = tapd.Ptr(flags.precondition)
	}
	if flags.expectation != "" {
		request.Expectation = tapd.Ptr(flags.expectation)
	}
	if flags.testType != "" {
		request.Type = tapd.Ptr(flags.testType)
	}
	if flags.priority != "" {
		request.Priority = tapd.Ptr(flags.priority)
	}
	if flags.creator != "" {
		request.Creator = tapd.Ptr(flags.creator)
	}
}

type testCaseQueryFlags struct {
	workspaceID        int
	limit              int
	page               int
	fields             string
	ids                string
	name               string
	steps              string
	categoryID         int64
	status             string
	precondition       string
	expectation        string
	testType           string
	priority           string
	creator            string
	modifier           string
	created            string
	modified           string
	isAutomated        string
	automationType     string
	automationPlatform string
	isServing          string
	testPlanID         int64
}

func addTestCaseQueryFlags(cmd *cobra.Command, flags *testCaseQueryFlags, withPaging bool) {
	newWorkspaceFlag(cmd, &flags.workspaceID)
	if withPaging {
		newListFlags(cmd, &flags.limit, &flags.page)
		cmd.Flags().StringVar(&flags.fields, "fields", "", "comma separated fields")
	}
	cmd.Flags().StringVar(&flags.ids, "ids", "", "comma separated test case IDs")
	cmd.Flags().StringVar(&flags.name, "name", "", "filter by test case name")
	cmd.Flags().StringVar(&flags.steps, "steps", "", "filter by test case steps")
	cmd.Flags().Int64Var(&flags.categoryID, "category-id", 0, "filter by category ID")
	cmd.Flags().StringVar(&flags.status, "status", "", "filter by status")
	cmd.Flags().StringVar(&flags.precondition, "precondition", "", "filter by precondition")
	cmd.Flags().StringVar(&flags.expectation, "expectation", "", "filter by expectation")
	cmd.Flags().StringVar(&flags.testType, "type", "", "filter by test case type")
	cmd.Flags().StringVar(&flags.priority, "priority", "", "filter by priority")
	cmd.Flags().StringVar(&flags.creator, "creator", "", "filter by creator")
	cmd.Flags().StringVar(&flags.modifier, "modifier", "", "filter by modifier")
	cmd.Flags().StringVar(&flags.created, "created", "", "filter by created time expression")
	cmd.Flags().StringVar(&flags.modified, "modified", "", "filter by modified time expression")
	cmd.Flags().StringVar(&flags.isAutomated, "is-automated", "", "filter by automation flag")
	cmd.Flags().StringVar(&flags.automationType, "automation-type", "", "filter by automation type")
	cmd.Flags().StringVar(&flags.automationPlatform, "automation-platform", "", "filter by automation platform")
	cmd.Flags().StringVar(&flags.isServing, "is-serving", "", "filter by serving flag")
	if !withPaging {
		cmd.Flags().Int64Var(&flags.testPlanID, "test-plan-id", 0, "filter by test plan ID")
	}
}

func newTestCasesRequest(flags testCaseQueryFlags) (*tapd.GetTestCasesRequest, error) {
	request := &tapd.GetTestCasesRequest{
		WorkspaceID: tapd.Ptr(flags.workspaceID),
		Limit:       tapd.Ptr(flags.limit),
		Page:        tapd.Ptr(flags.page),
		Fields:      fieldsMulti(flags.fields),
	}
	if flags.ids != "" {
		ids, err := strictInt64Multi("test case IDs", flags.ids)
		if err != nil {
			return nil, err
		}
		request.ID = ids
	}
	applyTestCaseQueryFilters(request, flags)
	return request, nil
}

func newTestCasesCountRequest(flags testCaseQueryFlags) (*tapd.GetTestCasesCountRequest, error) {
	request := &tapd.GetTestCasesCountRequest{
		WorkspaceID: tapd.Ptr(flags.workspaceID),
	}
	if flags.ids != "" {
		ids, err := strictInt64Multi("test case IDs", flags.ids)
		if err != nil {
			return nil, err
		}
		request.ID = ids
	}
	applyTestCaseCountFilters(request, flags)
	return request, nil
}

func applyTestCaseQueryFilters(request *tapd.GetTestCasesRequest, flags testCaseQueryFlags) {
	if flags.name != "" {
		request.Name = tapd.Ptr(flags.name)
	}
	if flags.steps != "" {
		request.Steps = tapd.Ptr(flags.steps)
	}
	if flags.categoryID > 0 {
		request.CategoryID = tapd.Ptr(flags.categoryID)
	}
	if flags.status != "" {
		request.Status = tapd.NewEnum(tapd.TestCaseStatus(flags.status))
	}
	if flags.precondition != "" {
		request.Precondition = tapd.Ptr(flags.precondition)
	}
	if flags.expectation != "" {
		request.Expectation = tapd.Ptr(flags.expectation)
	}
	if flags.testType != "" {
		request.Type = tapd.Ptr(flags.testType)
	}
	if flags.priority != "" {
		request.Priority = tapd.Ptr(flags.priority)
	}
	if flags.creator != "" {
		request.Creator = tapd.Ptr(flags.creator)
	}
	if flags.modifier != "" {
		request.Modifier = tapd.Ptr(flags.modifier)
	}
	if flags.created != "" {
		request.Created = tapd.Ptr(flags.created)
	}
	if flags.modified != "" {
		request.Modified = tapd.Ptr(flags.modified)
	}
	if flags.isAutomated != "" {
		request.IsAutomated = tapd.Ptr(flags.isAutomated)
	}
	if flags.automationType != "" {
		request.AutomationType = tapd.Ptr(flags.automationType)
	}
	if flags.automationPlatform != "" {
		request.AutomationPlatform = tapd.Ptr(flags.automationPlatform)
	}
	if flags.isServing != "" {
		request.IsServing = tapd.Ptr(flags.isServing)
	}
}

func applyTestCaseCountFilters(request *tapd.GetTestCasesCountRequest, flags testCaseQueryFlags) {
	if flags.name != "" {
		request.Name = tapd.Ptr(flags.name)
	}
	if flags.steps != "" {
		request.Steps = tapd.Ptr(flags.steps)
	}
	if flags.categoryID > 0 {
		request.CategoryID = tapd.Ptr(flags.categoryID)
	}
	if flags.status != "" {
		request.Status = tapd.NewEnum(tapd.TestCaseStatus(flags.status))
	}
	if flags.precondition != "" {
		request.Precondition = tapd.Ptr(flags.precondition)
	}
	if flags.expectation != "" {
		request.Expectation = tapd.Ptr(flags.expectation)
	}
	if flags.testType != "" {
		request.Type = tapd.Ptr(flags.testType)
	}
	if flags.priority != "" {
		request.Priority = tapd.Ptr(flags.priority)
	}
	if flags.creator != "" {
		request.Creator = tapd.Ptr(flags.creator)
	}
	if flags.modifier != "" {
		request.Modifier = tapd.Ptr(flags.modifier)
	}
	if flags.created != "" {
		request.Created = tapd.Ptr(flags.created)
	}
	if flags.modified != "" {
		request.Modified = tapd.Ptr(flags.modified)
	}
	if flags.isAutomated != "" {
		request.IsAutomated = tapd.Ptr(flags.isAutomated)
	}
	if flags.automationType != "" {
		request.AutomationType = tapd.Ptr(flags.automationType)
	}
	if flags.automationPlatform != "" {
		request.AutomationPlatform = tapd.Ptr(flags.automationPlatform)
	}
	if flags.isServing != "" {
		request.IsServing = tapd.Ptr(flags.isServing)
	}
	if flags.testPlanID > 0 {
		request.TestPlanID = tapd.Ptr(flags.testPlanID)
	}
}

type testPlanMutationFlags struct {
	workspaceID int
	name        string
	description string
	creator     string
	modifier    string
	owner       string
	startDate   string
	endDate     string
	iterationID int64
	version     string
	testType    string
	status      string
	templateID  int64
}

func addTestPlanMutationFlags(cmd *cobra.Command, flags *testPlanMutationFlags, update bool) {
	if update {
		cmd.Flags().StringVar(&flags.name, "name", "", "test plan name")
	}
	cmd.Flags().StringVar(&flags.description, "description", "", "test plan description")
	if !update {
		cmd.Flags().StringVar(&flags.creator, "creator", "", "test plan creator")
	}
	cmd.Flags().StringVar(&flags.modifier, "modifier", "", "test plan modifier")
	cmd.Flags().StringVar(&flags.owner, "owner", "", "test plan owner")
	cmd.Flags().StringVar(&flags.startDate, "start-date", "", "planned start date")
	cmd.Flags().StringVar(&flags.endDate, "end-date", "", "planned end date")
	if !update {
		cmd.Flags().Int64Var(&flags.iterationID, "iteration-id", 0, "related iteration ID")
	}
	cmd.Flags().StringVar(&flags.version, "version", "", "test plan version")
	cmd.Flags().StringVar(&flags.testType, "type", "", "test plan type")
	cmd.Flags().StringVar(&flags.status, "status", "", "test plan status")
	if update {
		cmd.Flags().Int64Var(&flags.templateID, "template-id", 0, "test plan template ID")
	}
}

func applyTestPlanCreateFlags(request *tapd.CreateTestPlanRequest, flags testPlanMutationFlags) {
	if flags.description != "" {
		request.Description = tapd.Ptr(flags.description)
	}
	if flags.creator != "" {
		request.Creator = tapd.Ptr(flags.creator)
	}
	if flags.modifier != "" {
		request.Modifier = tapd.Ptr(flags.modifier)
	}
	if flags.owner != "" {
		request.Owner = tapd.Ptr(flags.owner)
	}
	if flags.startDate != "" {
		request.StartDate = tapd.Ptr(flags.startDate)
	}
	if flags.endDate != "" {
		request.EndDate = tapd.Ptr(flags.endDate)
	}
	if flags.iterationID > 0 {
		request.IterationID = tapd.Ptr(flags.iterationID)
	}
	if flags.version != "" {
		request.Version = tapd.Ptr(flags.version)
	}
	if flags.testType != "" {
		request.Type = tapd.Ptr(flags.testType)
	}
	if flags.status != "" {
		request.Status = tapd.Ptr(flags.status)
	}
}

func applyTestPlanUpdateFlags(request *tapd.UpdateTestPlanRequest, flags testPlanMutationFlags) {
	if flags.name != "" {
		request.Name = tapd.Ptr(flags.name)
	}
	if flags.description != "" {
		request.Description = tapd.Ptr(flags.description)
	}
	if flags.modifier != "" {
		request.Modifier = tapd.Ptr(flags.modifier)
	}
	if flags.owner != "" {
		request.Owner = tapd.Ptr(flags.owner)
	}
	if flags.startDate != "" {
		request.StartDate = tapd.Ptr(flags.startDate)
	}
	if flags.endDate != "" {
		request.EndDate = tapd.Ptr(flags.endDate)
	}
	if flags.version != "" {
		request.Version = tapd.Ptr(flags.version)
	}
	if flags.testType != "" {
		request.Type = tapd.Ptr(flags.testType)
	}
	if flags.status != "" {
		request.Status = tapd.Ptr(flags.status)
	}
	if flags.templateID > 0 {
		request.TemplateID = tapd.Ptr(flags.templateID)
	}
}

type testPlanQueryFlags struct {
	workspaceID int
	limit       int
	page        int
	fields      string
	ids         string
	name        string
	description string
	creator     string
	modifier    string
	owner       string
	startDate   string
	endDate     string
	iterationID int64
	version     string
	testType    string
	status      string
}

func addTestPlanQueryFlags(cmd *cobra.Command, flags *testPlanQueryFlags, withPaging bool) {
	newWorkspaceFlag(cmd, &flags.workspaceID)
	if withPaging {
		newListFlags(cmd, &flags.limit, &flags.page)
		cmd.Flags().StringVar(&flags.fields, "fields", "", "comma separated fields")
	}
	cmd.Flags().StringVar(&flags.ids, "ids", "", "comma separated test plan IDs")
	cmd.Flags().StringVar(&flags.name, "name", "", "filter by test plan name")
	cmd.Flags().StringVar(&flags.description, "description", "", "filter by description")
	cmd.Flags().StringVar(&flags.creator, "creator", "", "filter by creator")
	cmd.Flags().StringVar(&flags.modifier, "modifier", "", "filter by modifier")
	cmd.Flags().StringVar(&flags.owner, "owner", "", "filter by owner")
	cmd.Flags().StringVar(&flags.startDate, "start-date", "", "filter by start date")
	cmd.Flags().StringVar(&flags.endDate, "end-date", "", "filter by end date")
	cmd.Flags().Int64Var(&flags.iterationID, "iteration-id", 0, "filter by iteration ID")
	cmd.Flags().StringVar(&flags.version, "version", "", "filter by version")
	cmd.Flags().StringVar(&flags.testType, "type", "", "filter by test plan type")
	cmd.Flags().StringVar(&flags.status, "status", "", "filter by status")
}

func newTestPlansRequest(flags testPlanQueryFlags) (*tapd.GetTestPlansRequest, error) {
	request := &tapd.GetTestPlansRequest{
		WorkspaceID: tapd.Ptr(flags.workspaceID),
		Limit:       tapd.Ptr(flags.limit),
		Page:        tapd.Ptr(flags.page),
		Fields:      fieldsMulti(flags.fields),
	}
	if flags.ids != "" {
		ids, err := strictInt64Multi("test plan IDs", flags.ids)
		if err != nil {
			return nil, err
		}
		request.ID = ids
	}
	applyTestPlanQueryFilters(request, flags)
	return request, nil
}

func newTestPlansCountRequest(flags testPlanQueryFlags) (*tapd.GetTestPlansCountRequest, error) {
	request := &tapd.GetTestPlansCountRequest{
		WorkspaceID: tapd.Ptr(flags.workspaceID),
	}
	if flags.ids != "" {
		ids, err := strictInt64Multi("test plan IDs", flags.ids)
		if err != nil {
			return nil, err
		}
		request.ID = ids
	}
	applyTestPlanCountFilters(request, flags)
	return request, nil
}

func applyTestPlanQueryFilters(request *tapd.GetTestPlansRequest, flags testPlanQueryFlags) {
	if flags.name != "" {
		request.Name = tapd.Ptr(flags.name)
	}
	if flags.description != "" {
		request.Description = tapd.Ptr(flags.description)
	}
	if flags.creator != "" {
		request.Creator = tapd.Ptr(flags.creator)
	}
	if flags.modifier != "" {
		request.Modifier = tapd.Ptr(flags.modifier)
	}
	if flags.owner != "" {
		request.Owner = tapd.Ptr(flags.owner)
	}
	if flags.startDate != "" {
		request.StartDate = tapd.Ptr(flags.startDate)
	}
	if flags.endDate != "" {
		request.EndDate = tapd.Ptr(flags.endDate)
	}
	if flags.iterationID > 0 {
		request.IterationID = tapd.Ptr(flags.iterationID)
	}
	if flags.version != "" {
		request.Version = tapd.Ptr(flags.version)
	}
	if flags.testType != "" {
		request.Type = tapd.Ptr(flags.testType)
	}
	if flags.status != "" {
		request.Status = tapd.Ptr(flags.status)
	}
}

func applyTestPlanCountFilters(request *tapd.GetTestPlansCountRequest, flags testPlanQueryFlags) {
	if flags.name != "" {
		request.Name = tapd.Ptr(flags.name)
	}
	if flags.description != "" {
		request.Description = tapd.Ptr(flags.description)
	}
	if flags.creator != "" {
		request.Creator = tapd.Ptr(flags.creator)
	}
	if flags.modifier != "" {
		request.Modifier = tapd.Ptr(flags.modifier)
	}
	if flags.owner != "" {
		request.Owner = tapd.Ptr(flags.owner)
	}
	if flags.startDate != "" {
		request.StartDate = tapd.Ptr(flags.startDate)
	}
	if flags.endDate != "" {
		request.EndDate = tapd.Ptr(flags.endDate)
	}
	if flags.iterationID > 0 {
		request.IterationID = tapd.Ptr(flags.iterationID)
	}
	if flags.version != "" {
		request.Version = tapd.Ptr(flags.version)
	}
	if flags.testType != "" {
		request.Type = tapd.Ptr(flags.testType)
	}
	if flags.status != "" {
		request.Status = tapd.Ptr(flags.status)
	}
}

func testCaseTableHeaders() []string {
	return []string{"ID", "Name", "Status", "CategoryID", "Priority", "Type", "Creator", "Modified"}
}

func testCaseRows(testCases []*tapd.TestCase) [][]string {
	rows := make([][]string, 0, len(testCases))
	for _, item := range testCases {
		if item == nil {
			continue
		}
		rows = append(rows, []string{
			item.ID,
			item.Name,
			string(item.Status),
			item.CategoryID,
			item.Priority,
			item.Type,
			item.Creator,
			item.Modified,
		})
	}
	return rows
}

func testPlanTableHeaders() []string {
	return []string{"ID", "Name", "Status", "Owner", "Type", "Version", "StartDate", "EndDate", "Modified"}
}

func testPlanRows(plans []*tapd.TestPlan) [][]string {
	rows := make([][]string, 0, len(plans))
	for _, item := range plans {
		if item == nil {
			continue
		}
		rows = append(rows, []string{
			item.ID,
			item.Name,
			item.Status,
			item.Owner,
			item.Type,
			item.Version,
			stringValue(item.StartDate),
			stringValue(item.EndDate),
			item.Modified,
		})
	}
	return rows
}

func testCaseCategoryTableHeaders() []string {
	return []string{"ID", "Name", "ParentID", "Creator", "Modifier", "Modified"}
}

func testCaseCategoryRows(categories []*tapd.TestCaseCategory) [][]string {
	rows := make([][]string, 0, len(categories))
	for _, item := range categories {
		if item == nil {
			continue
		}
		rows = append(rows, []string{
			item.ID,
			item.Name,
			item.ParentID,
			stringValue(item.Creator),
			stringValue(item.Modifier),
			item.Modified,
		})
	}
	return rows
}

func testFieldsTableHeaders() []string {
	return []string{"Name", "Label", "Type", "Default", "Options"}
}

func testCaseFieldsRows(fields []*tapd.TestCaseFieldsInfo) [][]string {
	rows := make([][]string, 0, len(fields))
	for _, item := range fields {
		if item == nil {
			continue
		}
		rows = append(rows, []string{
			item.Name,
			item.Label,
			string(item.HTMLType),
			item.Default,
			strconv.Itoa(len(item.Options)),
		})
	}
	return rows
}

func testCaseResultTableHeaders() []string {
	return []string{"ID", "TestCaseID", "TestPlanID", "Status", "Executor", "ExecutedAt", "BugIDs"}
}

func testCaseResultRows(items []*tapd.TestCaseResultItem) [][]string {
	rows := make([][]string, 0, len(items))
	for _, item := range items {
		if item == nil || item.Result == nil {
			continue
		}
		result := item.Result
		rows = append(rows, []string{
			item.ID,
			result.TestCaseID,
			result.TestPlanID,
			string(result.ResultStatus),
			result.Executor,
			result.ExecutedAt,
			strings.Join(result.BugID, ","),
		})
	}
	return rows
}

func testPlanProgressTableHeaders() []string {
	return []string{"StoryCount", "TestCaseCount", "ExecutedRate", "Pass", "NoPass", "Block", "Unexecuted"}
}

func testPlanProgressRows(progress *tapd.TestPlanProgress) [][]string {
	if progress == nil {
		return nil
	}
	return [][]string{{
		strconv.Itoa(progress.StoryCount),
		strconv.Itoa(progress.TestCaseCount),
		progress.ExecutedRate,
		strconv.Itoa(progress.StatusCounter[tapd.TestCaseResultStatusPass]),
		strconv.Itoa(progress.StatusCounter[tapd.TestCaseResultStatusNoPass]),
		strconv.Itoa(progress.StatusCounter[tapd.TestCaseResultStatusBlock]),
		strconv.Itoa(progress.StatusCounter[tapd.TestCaseResultStatusUnexecuted]),
	}}
}

func testPlanRelatedBugTableHeaders() []string {
	return []string{"ID", "Name", "ResultCount"}
}

func testPlanRelatedBugRows(items []*tapd.TestPlanRelatedBug) [][]string {
	rows := make([][]string, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		rows = append(rows, []string{
			item.ID,
			item.Name,
			strconv.Itoa(len(item.TestCaseResultRelatedBugs)),
		})
	}
	return rows
}
