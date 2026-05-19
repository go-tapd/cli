package cmd

import (
	"strconv"

	"github.com/go-tapd/cli/internal/app"
	"github.com/go-tapd/tapd"
	"github.com/spf13/cobra"
)

func newCommentCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "comment",
		Short: "Work with TAPD comments",
	}

	cmd.AddCommand(
		newCommentCreateCmd(rt),
		newCommentListCmd(rt),
		newCommentCountCmd(rt),
		newCommentUpdateCmd(rt),
	)
	return cmd
}

func newCommentCreateCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		title       string
		description string
		author      string
		entryType   string
		entryID     int64
		replyID     int64
		rootID      int64
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a comment",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.CreateCommentRequest{
				Description: tapd.Ptr(description),
				Author:      tapd.Ptr(author),
				EntryType:   tapd.Ptr(tapd.CommentEntryType(entryType)),
				EntryID:     tapd.Ptr(entryID),
				WorkspaceID: tapd.Ptr(workspaceID),
			}
			if title != "" {
				request.Title = tapd.Ptr(title)
			}
			if replyID > 0 {
				request.ReplyID = tapd.Ptr(replyID)
			}
			if rootID > 0 {
				request.RootID = tapd.Ptr(rootID)
			}

			comment, _, err := client.CommentService.CreateComment(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, commentHeaders(), commentRows([]*tapd.Comment{comment}), comment)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().StringVar(&title, "title", "", "comment title")
	cmd.Flags().StringVar(&description, "description", "", "comment content")
	cmd.Flags().StringVar(&author, "author", "", "comment author")
	cmd.Flags().StringVar(&entryType, "entry-type", "", "entry type: bug, bug_remark, stories, tasks, wiki, or mini_items")
	cmd.Flags().Int64Var(&entryID, "entry-id", 0, "entry ID")
	cmd.Flags().Int64Var(&replyID, "reply-id", 0, "reply comment ID")
	cmd.Flags().Int64Var(&rootID, "root-id", 0, "root comment ID")
	_ = cmd.MarkFlagRequired("description")
	_ = cmd.MarkFlagRequired("author")
	_ = cmd.MarkFlagRequired("entry-type")
	_ = cmd.MarkFlagRequired("entry-id")
	return cmd
}

func newCommentListCmd(rt *app.Runtime) *cobra.Command {
	var flags commentQueryFlags

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List comments",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request, err := newGetCommentsRequest(flags)
			if err != nil {
				return err
			}

			comments, _, err := client.CommentService.GetComments(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, commentHeaders(), commentRows(comments), comments)
		},
	}

	addCommentQueryFlags(cmd, &flags, true)
	return cmd
}

func newCommentCountCmd(rt *app.Runtime) *cobra.Command {
	var flags commentQueryFlags

	cmd := &cobra.Command{
		Use:   "count",
		Short: "Count comments",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request, err := newGetCommentsCountRequest(flags)
			if err != nil {
				return err
			}

			count, _, err := client.CommentService.GetCommentsCount(cmd.Context(), request)
			if err != nil {
				return err
			}

			data := map[string]any{"resource": "comment", "workspace_id": flags.workspaceID, "count": count}
			rows := [][]string{{"comment", strconv.Itoa(flags.workspaceID), strconv.Itoa(count)}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Resource", "WorkspaceID", "Count"}, rows, data)
		},
	}

	addCommentQueryFlags(cmd, &flags, false)
	return cmd
}

func newCommentUpdateCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID   int
		description   string
		changeCreator string
	)

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a comment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseInt64Arg("comment id", args[0])
			if err != nil {
				return err
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.UpdateCommentRequest{
				WorkspaceID: tapd.Ptr(workspaceID),
				ID:          tapd.Ptr(id),
				Description: tapd.Ptr(description),
			}
			if changeCreator != "" {
				request.ChangeCreator = tapd.Ptr(changeCreator)
			}

			comment, _, err := client.CommentService.UpdateComment(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, commentHeaders(), commentRows([]*tapd.Comment{comment}), comment)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().StringVar(&description, "description", "", "comment content")
	cmd.Flags().StringVar(&changeCreator, "change-creator", "", "change creator")
	_ = cmd.MarkFlagRequired("description")
	return cmd
}

type commentQueryFlags struct {
	workspaceID int
	limit       int
	page        int
	fields      string
	ids         string
	title       string
	description string
	author      string
	entryType   string
	entryID     int64
	created     string
	modified    string
	rootID      int64
	replyID     int64
}

func addCommentQueryFlags(cmd *cobra.Command, flags *commentQueryFlags, withPaging bool) {
	newWorkspaceFlag(cmd, &flags.workspaceID)
	if withPaging {
		newListFlags(cmd, &flags.limit, &flags.page)
		cmd.Flags().StringVar(&flags.fields, "fields", "", "comma separated fields")
	}
	cmd.Flags().StringVar(&flags.ids, "ids", "", "comma separated comment IDs")
	cmd.Flags().StringVar(&flags.title, "title", "", "filter by title")
	cmd.Flags().StringVar(&flags.description, "description", "", "filter by content")
	cmd.Flags().StringVar(&flags.author, "author", "", "filter by author")
	cmd.Flags().StringVar(&flags.entryType, "entry-type", "", "entry type: bug, bug_remark, stories, tasks, wiki, or mini_items")
	cmd.Flags().Int64Var(&flags.entryID, "entry-id", 0, "entry ID")
	cmd.Flags().StringVar(&flags.created, "created", "", "filter by created time expression")
	cmd.Flags().StringVar(&flags.modified, "modified", "", "filter by modified time expression")
	cmd.Flags().Int64Var(&flags.rootID, "root-id", 0, "root comment ID")
	cmd.Flags().Int64Var(&flags.replyID, "reply-id", 0, "reply comment ID")
}

func newGetCommentsRequest(flags commentQueryFlags) (*tapd.GetCommentsRequest, error) {
	request := &tapd.GetCommentsRequest{
		WorkspaceID: tapd.Ptr(flags.workspaceID),
		Limit:       tapd.Ptr(flags.limit),
		Page:        tapd.Ptr(flags.page),
		Fields:      fieldsMulti(flags.fields),
	}
	if flags.ids != "" {
		ids, err := strictInt64Multi("comment IDs", flags.ids)
		if err != nil {
			return nil, err
		}
		request.ID = ids
	}
	if flags.title != "" {
		request.Title = tapd.Ptr(flags.title)
	}
	if flags.description != "" {
		request.Description = tapd.Ptr(flags.description)
	}
	if flags.author != "" {
		request.Author = tapd.Ptr(flags.author)
	}
	if flags.entryType != "" {
		request.EntryType = tapd.Ptr(tapd.CommentEntryType(flags.entryType))
	}
	if flags.entryID > 0 {
		request.EntryID = tapd.Ptr(flags.entryID)
	}
	if flags.created != "" {
		request.Created = tapd.Ptr(flags.created)
	}
	if flags.modified != "" {
		request.Modified = tapd.Ptr(flags.modified)
	}
	if flags.rootID > 0 {
		request.RootID = tapd.Ptr(flags.rootID)
	}
	if flags.replyID > 0 {
		request.ReplyID = tapd.Ptr(flags.replyID)
	}
	return request, nil
}

func newGetCommentsCountRequest(flags commentQueryFlags) (*tapd.GetCommentsCountRequest, error) {
	request := &tapd.GetCommentsCountRequest{WorkspaceID: tapd.Ptr(flags.workspaceID)}
	if flags.ids != "" {
		ids, err := strictInt64Multi("comment IDs", flags.ids)
		if err != nil {
			return nil, err
		}
		request.ID = ids
	}
	if flags.title != "" {
		request.Title = tapd.Ptr(flags.title)
	}
	if flags.description != "" {
		request.Description = tapd.Ptr(flags.description)
	}
	if flags.author != "" {
		request.Author = tapd.Ptr(flags.author)
	}
	if flags.entryType != "" {
		request.EntryType = tapd.Ptr(tapd.CommentEntryType(flags.entryType))
	}
	if flags.entryID > 0 {
		request.EntryID = tapd.Ptr(flags.entryID)
	}
	if flags.created != "" {
		request.Created = tapd.Ptr(flags.created)
	}
	if flags.modified != "" {
		request.Modified = tapd.Ptr(flags.modified)
	}
	if flags.rootID > 0 {
		request.RootID = tapd.Ptr(flags.rootID)
	}
	if flags.replyID > 0 {
		request.ReplyID = tapd.Ptr(flags.replyID)
	}
	return request, nil
}

func commentHeaders() []string {
	return []string{"ID", "EntryType", "EntryID", "Author", "Title", "Created", "Modified"}
}

func commentRows(comments []*tapd.Comment) [][]string {
	rows := make([][]string, 0, len(comments))
	for _, item := range comments {
		rows = append(rows, []string{
			item.ID,
			string(item.EntryType),
			item.EntryID,
			item.Author,
			item.Title,
			item.Created,
			item.Modified,
		})
	}
	return rows
}
