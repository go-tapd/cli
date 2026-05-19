package cmd

import (
	"errors"
	"strconv"

	"github.com/go-tapd/cli/internal/app"
	"github.com/go-tapd/tapd"
	"github.com/spf13/cobra"
)

func newWikiCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wiki",
		Short: "Work with TAPD wikis",
	}

	cmd.AddCommand(
		newWikiCreateCmd(rt),
		newWikiListCmd(rt),
		newWikiCountCmd(rt),
		newWikiUpdateCmd(rt),
		newWikiDrawioCmd(rt),
		newWikiFollowersCmd(rt),
		newWikiPermissionsCmd(rt),
		newWikiTagsCmd(rt),
		newWikiAttachmentsCmd(rt),
	)
	return cmd
}

func newWikiCreateCmd(rt *app.Runtime) *cobra.Command {
	var flags wikiMutationFlags

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a wiki",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.CreateWikiRequest{
				WorkspaceID: tapd.Ptr(flags.workspaceID),
				Name:        tapd.Ptr(flags.name),
				Creator:     tapd.Ptr(flags.creator),
			}
			if flags.markdownDescription != "" {
				request.MarkdownDescription = tapd.Ptr(flags.markdownDescription)
			}
			if flags.description != "" {
				request.Description = tapd.Ptr(flags.description)
			}
			if flags.note != "" {
				request.Note = tapd.Ptr(flags.note)
			}
			if flags.parentWikiID != "" {
				request.ParentWikiID = tapd.Ptr(flags.parentWikiID)
			}

			wiki, _, err := client.WikiService.CreateWiki(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, wikiHeaders(), wikiRows([]*tapd.Wiki{wiki}), wiki)
		},
	}

	newWorkspaceFlag(cmd, &flags.workspaceID)
	cmd.Flags().StringVar(&flags.name, "name", "", "wiki title")
	cmd.Flags().StringVar(&flags.creator, "creator", "", "wiki creator")
	addWikiMutationFlags(cmd, &flags, false)
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("creator")
	return cmd
}

func newWikiListCmd(rt *app.Runtime) *cobra.Command {
	var flags wikiQueryFlags

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List wikis",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			wikis, _, err := client.WikiService.GetWikis(cmd.Context(), newGetWikisRequest(flags))
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, wikiHeaders(), wikiRows(wikis), wikis)
		},
	}

	addWikiQueryFlags(cmd, &flags, true)
	return cmd
}

func newWikiCountCmd(rt *app.Runtime) *cobra.Command {
	var flags wikiQueryFlags

	cmd := &cobra.Command{
		Use:   "count",
		Short: "Count wikis",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			count, _, err := client.WikiService.GetWikisCount(cmd.Context(), newGetWikisCountRequest(flags))
			if err != nil {
				return err
			}

			data := map[string]any{"resource": "wiki", "workspace_id": flags.workspaceID, "count": count}
			rows := [][]string{{"wiki", strconv.Itoa(flags.workspaceID), strconv.Itoa(count)}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Resource", "WorkspaceID", "Count"}, rows, data)
		},
	}

	addWikiQueryFlags(cmd, &flags, false)
	return cmd
}

func newWikiUpdateCmd(rt *app.Runtime) *cobra.Command {
	var flags wikiMutationFlags

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a wiki",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseInt64Arg("wiki id", args[0])
			if err != nil {
				return err
			}
			if !hasWikiUpdate(flags) {
				return errors.New("at least one wiki field flag is required")
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.UpdateWikiRequest{
				WorkspaceID: tapd.Ptr(flags.workspaceID),
				ID:          tapd.Ptr(id),
			}
			if flags.name != "" {
				request.Name = tapd.Ptr(flags.name)
			}
			if flags.markdownDescription != "" {
				request.MarkdownDescription = tapd.Ptr(flags.markdownDescription)
			}
			if flags.description != "" {
				request.Description = tapd.Ptr(flags.description)
			}
			if flags.note != "" {
				request.Note = tapd.Ptr(flags.note)
			}
			if flags.parentWikiID != "" {
				request.ParentWikiID = tapd.Ptr(flags.parentWikiID)
			}

			wiki, _, err := client.WikiService.UpdateWiki(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, wikiHeaders(), wikiRows([]*tapd.Wiki{wiki}), wiki)
		},
	}

	newWorkspaceFlag(cmd, &flags.workspaceID)
	cmd.Flags().StringVar(&flags.name, "name", "", "wiki title")
	addWikiMutationFlags(cmd, &flags, true)
	return cmd
}

func newWikiDrawioCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		id          int64
		token       string
	)

	cmd := &cobra.Command{
		Use:   "drawio",
		Short: "Show wiki drawio data",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetWikiDrawioDataRequest{
				WorkspaceID: tapd.Ptr(workspaceID),
				ID:          tapd.Ptr(id),
			}
			if token != "" {
				request.Token = tapd.Ptr(token)
			}

			data, _, err := client.WikiService.GetWikiDrawioData(cmd.Context(), request)
			if err != nil {
				return err
			}

			rows := [][]string{{data.ID, strconv.Itoa(len(data.Values))}}
			return writeOutput(cmd, rt.OutputFormat, []string{"ID", "ValuesBytes"}, rows, data)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().Int64Var(&id, "id", 0, "drawio data ID")
	cmd.Flags().StringVar(&token, "token", "", "drawio token")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func newWikiFollowersCmd(rt *app.Runtime) *cobra.Command {
	var flags wikiFollowerFlags

	cmd := &cobra.Command{
		Use:   "followers",
		Short: "List wiki followers",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			followers, _, err := client.WikiService.GetWikiFollowers(cmd.Context(), newGetWikiFollowersRequest(flags))
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, []string{"ID", "WikiID", "User", "Created"}, wikiFollowerRows(followers), followers)
		},
	}

	addWikiFollowerFlags(cmd, &flags, true)
	cmd.AddCommand(newWikiFollowersCountCmd(rt))
	return cmd
}

func newWikiFollowersCountCmd(rt *app.Runtime) *cobra.Command {
	var flags wikiFollowerFlags

	cmd := &cobra.Command{
		Use:   "count",
		Short: "Count wiki followers",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			count, _, err := client.WikiService.GetWikiFollowersCount(cmd.Context(), newGetWikiFollowersCountRequest(flags))
			if err != nil {
				return err
			}

			data := map[string]any{"resource": "wiki_follower", "workspace_id": flags.workspaceID, "count": count}
			rows := [][]string{{"wiki_follower", strconv.Itoa(flags.workspaceID), strconv.Itoa(count)}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Resource", "WorkspaceID", "Count"}, rows, data)
		},
	}

	addWikiFollowerFlags(cmd, &flags, false)
	return cmd
}

func newWikiPermissionsCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		wikiID      int64
		targetType  string
		targetID    string
		limit       int
		page        int
		fields      string
	)

	cmd := &cobra.Command{
		Use:   "permissions",
		Short: "List wiki entity permissions",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetWikiEntityPermissionsRequest{
				WorkspaceID: tapd.Ptr(workspaceID),
				WikiID:      tapd.Ptr(wikiID),
				Limit:       tapd.Ptr(limit),
				Page:        tapd.Ptr(page),
				Fields:      fieldsMulti(fields),
			}
			if targetType != "" {
				request.TargetType = tapd.Ptr(targetType)
			}
			if targetID != "" {
				request.TargetID = tapd.Ptr(targetID)
			}

			permissions, _, err := client.WikiService.GetWikiEntityPermissions(cmd.Context(), request)
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(permissions))
			for _, item := range permissions {
				if item == nil {
					continue
				}
				rows = append(rows, []string{item.ID, item.WikiID, item.TargetType, item.TargetID, item.EntryType})
			}
			return writeOutput(cmd, rt.OutputFormat, []string{"ID", "WikiID", "TargetType", "TargetID", "EntryType"}, rows, permissions)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	newListFlags(cmd, &limit, &page)
	cmd.Flags().Int64Var(&wikiID, "wiki-id", 0, "wiki ID")
	cmd.Flags().StringVar(&targetType, "target-type", "", "filter by target type")
	cmd.Flags().StringVar(&targetID, "target-id", "", "filter by target ID")
	cmd.Flags().StringVar(&fields, "fields", "", "comma separated fields")
	_ = cmd.MarkFlagRequired("wiki-id")
	return cmd
}

func newWikiTagsCmd(rt *app.Runtime) *cobra.Command {
	var flags wikiTagFlags

	cmd := &cobra.Command{
		Use:   "tags",
		Short: "List wiki tags",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			tags, _, err := client.WikiService.GetWikiTags(cmd.Context(), newGetWikiTagsRequest(flags))
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, []string{"WikiID", "Tag", "Creator", "Created"}, wikiTagRows(tags), tags)
		},
	}

	addWikiTagFlags(cmd, &flags, true)
	cmd.AddCommand(newWikiTagsCountCmd(rt))
	return cmd
}

func newWikiTagsCountCmd(rt *app.Runtime) *cobra.Command {
	var flags wikiTagFlags

	cmd := &cobra.Command{
		Use:   "count",
		Short: "Count wiki tags",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			count, _, err := client.WikiService.GetWikiTagsCount(cmd.Context(), newGetWikiTagsCountRequest(flags))
			if err != nil {
				return err
			}

			data := map[string]any{"resource": "wiki_tag", "workspace_id": flags.workspaceID, "count": count}
			rows := [][]string{{"wiki_tag", strconv.Itoa(flags.workspaceID), strconv.Itoa(count)}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Resource", "WorkspaceID", "Count"}, rows, data)
		},
	}

	addWikiTagFlags(cmd, &flags, false)
	return cmd
}

func newWikiAttachmentsCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attachments",
		Short: "Work with wiki attachments",
	}
	cmd.AddCommand(newWikiAttachmentsCountCmd(rt))
	return cmd
}

func newWikiAttachmentsCountCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		id          int64
		filename    string
		size        int
		owner       string
		created     string
		modified    string
		wikiID      int64
	)

	cmd := &cobra.Command{
		Use:   "count",
		Short: "Count wiki attachments",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetWikiAttachmentsCountRequest{WorkspaceID: tapd.Ptr(workspaceID)}
			if id > 0 {
				request.ID = tapd.Ptr(id)
			}
			if filename != "" {
				request.Filename = tapd.Ptr(filename)
			}
			if size > 0 {
				request.Size = tapd.Ptr(size)
			}
			if owner != "" {
				request.Owner = tapd.Ptr(owner)
			}
			if created != "" {
				request.Created = tapd.Ptr(created)
			}
			if modified != "" {
				request.Modified = tapd.Ptr(modified)
			}
			if wikiID > 0 {
				request.WikiID = tapd.Ptr(wikiID)
			}

			count, _, err := client.WikiService.GetWikiAttachmentsCount(cmd.Context(), request)
			if err != nil {
				return err
			}

			data := map[string]any{"resource": "wiki_attachment", "workspace_id": workspaceID, "count": count}
			rows := [][]string{{"wiki_attachment", strconv.Itoa(workspaceID), strconv.Itoa(count)}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Resource", "WorkspaceID", "Count"}, rows, data)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().Int64Var(&id, "id", 0, "attachment ID")
	cmd.Flags().StringVar(&filename, "filename", "", "filter by filename")
	cmd.Flags().IntVar(&size, "size", 0, "filter by size")
	cmd.Flags().StringVar(&owner, "owner", "", "filter by owner")
	cmd.Flags().StringVar(&created, "created", "", "filter by created time expression")
	cmd.Flags().StringVar(&modified, "modified", "", "filter by modified time expression")
	cmd.Flags().Int64Var(&wikiID, "wiki-id", 0, "wiki ID")
	return cmd
}

type wikiMutationFlags struct {
	workspaceID         int
	name                string
	markdownDescription string
	description         string
	creator             string
	note                string
	parentWikiID        string
}

func addWikiMutationFlags(cmd *cobra.Command, flags *wikiMutationFlags, _ bool) {
	cmd.Flags().StringVar(&flags.markdownDescription, "markdown-description", "", "wiki markdown content")
	cmd.Flags().StringVar(&flags.description, "description", "", "wiki rich text content")
	cmd.Flags().StringVar(&flags.note, "note", "", "wiki note")
	cmd.Flags().StringVar(&flags.parentWikiID, "parent-wiki-id", "", "parent wiki ID")
}

func hasWikiUpdate(flags wikiMutationFlags) bool {
	return flags.name != "" ||
		flags.markdownDescription != "" ||
		flags.description != "" ||
		flags.note != "" ||
		flags.parentWikiID != ""
}

type wikiQueryFlags struct {
	workspaceID int
	limit       int
	page        int
	fields      string
	id          int64
	name        string
	modifier    string
	creator     string
	note        string
	viewCount   string
	created     string
	modified    string
}

func addWikiQueryFlags(cmd *cobra.Command, flags *wikiQueryFlags, withPaging bool) {
	newWorkspaceFlag(cmd, &flags.workspaceID)
	if withPaging {
		newListFlags(cmd, &flags.limit, &flags.page)
		cmd.Flags().StringVar(&flags.fields, "fields", "", "comma separated fields")
	}
	cmd.Flags().Int64Var(&flags.id, "id", 0, "wiki ID")
	cmd.Flags().StringVar(&flags.name, "name", "", "filter by title")
	cmd.Flags().StringVar(&flags.modifier, "modifier", "", "filter by modifier")
	cmd.Flags().StringVar(&flags.creator, "creator", "", "filter by creator")
	cmd.Flags().StringVar(&flags.note, "note", "", "filter by note")
	cmd.Flags().StringVar(&flags.viewCount, "view-count", "", "filter by view count")
	cmd.Flags().StringVar(&flags.created, "created", "", "filter by created time expression")
	cmd.Flags().StringVar(&flags.modified, "modified", "", "filter by modified time expression")
}

func newGetWikisRequest(flags wikiQueryFlags) *tapd.GetWikisRequest {
	request := &tapd.GetWikisRequest{
		WorkspaceID: tapd.Ptr(flags.workspaceID),
		Limit:       tapd.Ptr(flags.limit),
		Page:        tapd.Ptr(flags.page),
		Fields:      fieldsMulti(flags.fields),
	}
	if flags.id > 0 {
		request.ID = tapd.Ptr(flags.id)
	}
	if flags.name != "" {
		request.Name = tapd.Ptr(flags.name)
	}
	if flags.modifier != "" {
		request.Modifier = tapd.Ptr(flags.modifier)
	}
	if flags.creator != "" {
		request.Creator = tapd.Ptr(flags.creator)
	}
	if flags.note != "" {
		request.Note = tapd.Ptr(flags.note)
	}
	if flags.viewCount != "" {
		request.ViewCount = tapd.Ptr(flags.viewCount)
	}
	if flags.created != "" {
		request.Created = tapd.Ptr(flags.created)
	}
	if flags.modified != "" {
		request.Modified = tapd.Ptr(flags.modified)
	}
	return request
}

func newGetWikisCountRequest(flags wikiQueryFlags) *tapd.GetWikisCountRequest {
	request := &tapd.GetWikisCountRequest{WorkspaceID: tapd.Ptr(flags.workspaceID)}
	if flags.name != "" {
		request.Name = tapd.Ptr(flags.name)
	}
	if flags.modifier != "" {
		request.Modifier = tapd.Ptr(flags.modifier)
	}
	if flags.creator != "" {
		request.Creator = tapd.Ptr(flags.creator)
	}
	if flags.note != "" {
		request.Note = tapd.Ptr(flags.note)
	}
	if flags.viewCount != "" {
		request.ViewCount = tapd.Ptr(flags.viewCount)
	}
	if flags.created != "" {
		request.Created = tapd.Ptr(flags.created)
	}
	if flags.modified != "" {
		request.Modified = tapd.Ptr(flags.modified)
	}
	return request
}

type wikiFollowerFlags struct {
	workspaceID int
	limit       int
	page        int
	fields      string
	id          int64
	created     string
	wikiID      int64
	user        string
}

func addWikiFollowerFlags(cmd *cobra.Command, flags *wikiFollowerFlags, withPaging bool) {
	newWorkspaceFlag(cmd, &flags.workspaceID)
	if withPaging {
		newListFlags(cmd, &flags.limit, &flags.page)
		cmd.Flags().StringVar(&flags.fields, "fields", "", "comma separated fields")
	}
	cmd.Flags().Int64Var(&flags.id, "id", 0, "follower record ID")
	cmd.Flags().StringVar(&flags.created, "created", "", "filter by created time expression")
	cmd.Flags().Int64Var(&flags.wikiID, "wiki-id", 0, "wiki ID")
	cmd.Flags().StringVar(&flags.user, "user", "", "filter by user")
}

func newGetWikiFollowersRequest(flags wikiFollowerFlags) *tapd.GetWikiFollowersRequest {
	request := &tapd.GetWikiFollowersRequest{
		WorkspaceID: tapd.Ptr(flags.workspaceID),
		Limit:       tapd.Ptr(flags.limit),
		Page:        tapd.Ptr(flags.page),
		Fields:      fieldsMulti(flags.fields),
	}
	if flags.id > 0 {
		request.ID = tapd.Ptr(flags.id)
	}
	if flags.created != "" {
		request.Created = tapd.Ptr(flags.created)
	}
	if flags.wikiID > 0 {
		request.WikiID = tapd.Ptr(flags.wikiID)
	}
	if flags.user != "" {
		request.User = tapd.Ptr(flags.user)
	}
	return request
}

func newGetWikiFollowersCountRequest(flags wikiFollowerFlags) *tapd.GetWikiFollowersCountRequest {
	request := &tapd.GetWikiFollowersCountRequest{WorkspaceID: tapd.Ptr(flags.workspaceID)}
	if flags.id > 0 {
		request.ID = tapd.Ptr(flags.id)
	}
	if flags.created != "" {
		request.Created = tapd.Ptr(flags.created)
	}
	if flags.wikiID > 0 {
		request.WikiID = tapd.Ptr(flags.wikiID)
	}
	if flags.user != "" {
		request.User = tapd.Ptr(flags.user)
	}
	return request
}

type wikiTagFlags struct {
	workspaceID int
	limit       int
	page        int
	wikiID      int64
	tag         string
	creator     string
	created     string
}

func addWikiTagFlags(cmd *cobra.Command, flags *wikiTagFlags, withPaging bool) {
	newWorkspaceFlag(cmd, &flags.workspaceID)
	if withPaging {
		newListFlags(cmd, &flags.limit, &flags.page)
	}
	cmd.Flags().Int64Var(&flags.wikiID, "wiki-id", 0, "wiki ID")
	cmd.Flags().StringVar(&flags.tag, "tag", "", "filter by tag")
	cmd.Flags().StringVar(&flags.creator, "creator", "", "filter by creator")
	cmd.Flags().StringVar(&flags.created, "created", "", "filter by created time expression")
}

func newGetWikiTagsRequest(flags wikiTagFlags) *tapd.GetWikiTagsRequest {
	request := &tapd.GetWikiTagsRequest{
		WorkspaceID: tapd.Ptr(flags.workspaceID),
		Limit:       tapd.Ptr(flags.limit),
		Page:        tapd.Ptr(flags.page),
	}
	if flags.wikiID > 0 {
		request.WikiID = tapd.Ptr(flags.wikiID)
	}
	if flags.tag != "" {
		request.Tag = tapd.Ptr(flags.tag)
	}
	if flags.creator != "" {
		request.Creator = tapd.Ptr(flags.creator)
	}
	if flags.created != "" {
		request.Created = tapd.Ptr(flags.created)
	}
	return request
}

func newGetWikiTagsCountRequest(flags wikiTagFlags) *tapd.GetWikiTagsCountRequest {
	request := &tapd.GetWikiTagsCountRequest{WorkspaceID: tapd.Ptr(flags.workspaceID)}
	if flags.wikiID > 0 {
		request.WikiID = tapd.Ptr(flags.wikiID)
	}
	if flags.tag != "" {
		request.Tag = tapd.Ptr(flags.tag)
	}
	if flags.creator != "" {
		request.Creator = tapd.Ptr(flags.creator)
	}
	if flags.created != "" {
		request.Created = tapd.Ptr(flags.created)
	}
	return request
}

func wikiHeaders() []string {
	return []string{"ID", "Name", "Creator", "Modifier", "ViewCount", "Modified"}
}

func wikiRows(wikis []*tapd.Wiki) [][]string {
	rows := make([][]string, 0, len(wikis))
	for _, item := range wikis {
		if item == nil {
			continue
		}
		rows = append(rows, []string{item.ID, item.Name, item.Creator, item.Modifier, item.ViewCount, item.Modified})
	}
	return rows
}

func wikiFollowerRows(followers []*tapd.WikiFollower) [][]string {
	rows := make([][]string, 0, len(followers))
	for _, item := range followers {
		if item == nil {
			continue
		}
		rows = append(rows, []string{item.ID, item.WikiID, item.User, item.Created})
	}
	return rows
}

func wikiTagRows(tags []*tapd.WikiTag) [][]string {
	rows := make([][]string, 0, len(tags))
	for _, item := range tags {
		if item == nil {
			continue
		}
		rows = append(rows, []string{item.WikiID, item.Tag, item.Creator, item.Created})
	}
	return rows
}
