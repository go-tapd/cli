package cmd

import (
	"strconv"

	"github.com/go-tapd/cli/internal/app"
	"github.com/go-tapd/tapd"
	"github.com/spf13/cobra"
)

const liteCommentEntryType = tapd.CommentEntryTypeMiniItems

func newLiteCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lite",
		Short: "Work with TAPD Lite resources",
	}

	cmd.AddCommand(newLiteCommentCmd(rt))
	return cmd
}

func newLiteCommentCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "comment",
		Short: "Work with TAPD Lite workitem comments",
	}

	cmd.AddCommand(
		newLiteCommentCreateCmd(rt),
		newLiteCommentListCmd(rt),
		newLiteCommentCountCmd(rt),
		newLiteCommentUpdateCmd(rt),
	)
	return cmd
}

func newLiteCommentCreateCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		title       string
		description string
		author      string
		entryID     int64
		replyID     int64
		rootID      int64
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a Lite workitem comment",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.CreateCommentRequest{
				Description: new(description),
				Author:      new(author),
				EntryType:   new(liteCommentEntryType),
				EntryID:     new(entryID),
				WorkspaceID: new(workspaceID),
			}
			if title != "" {
				request.Title = new(title)
			}
			if replyID > 0 {
				request.ReplyID = new(replyID)
			}
			if rootID > 0 {
				request.RootID = new(rootID)
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
	cmd.Flags().Int64Var(&entryID, "entry-id", 0, "Lite workitem ID")
	cmd.Flags().Int64Var(&replyID, "reply-id", 0, "reply comment ID")
	cmd.Flags().Int64Var(&rootID, "root-id", 0, "root comment ID")
	_ = cmd.MarkFlagRequired("description")
	_ = cmd.MarkFlagRequired("author")
	_ = cmd.MarkFlagRequired("entry-id")
	return cmd
}

func newLiteCommentListCmd(rt *app.Runtime) *cobra.Command {
	var flags commentQueryFlags

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Lite workitem comments",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			flags.entryType = liteCommentEntryType.String()
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

	addLiteCommentQueryFlags(cmd, &flags, true)
	return cmd
}

func newLiteCommentCountCmd(rt *app.Runtime) *cobra.Command {
	var flags commentQueryFlags

	cmd := &cobra.Command{
		Use:   "count",
		Short: "Count Lite workitem comments",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			flags.entryType = liteCommentEntryType.String()
			request, err := newGetCommentsCountRequest(flags)
			if err != nil {
				return err
			}

			count, _, err := client.CommentService.GetCommentsCount(cmd.Context(), request)
			if err != nil {
				return err
			}

			data := map[string]any{"resource": "lite_comment", "workspace_id": flags.workspaceID, "count": count}
			rows := [][]string{{"lite_comment", strconv.Itoa(flags.workspaceID), strconv.Itoa(count)}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Resource", "WorkspaceID", "Count"}, rows, data)
		},
	}

	addLiteCommentQueryFlags(cmd, &flags, false)
	return cmd
}

func newLiteCommentUpdateCmd(rt *app.Runtime) *cobra.Command {
	cmd := newCommentUpdateCmd(rt)
	cmd.Short = "Update a Lite workitem comment"
	return cmd
}

func addLiteCommentQueryFlags(cmd *cobra.Command, flags *commentQueryFlags, withPaging bool) {
	newWorkspaceFlag(cmd, &flags.workspaceID)
	if withPaging {
		newListFlags(cmd, &flags.limit, &flags.page)
		cmd.Flags().StringVar(&flags.fields, "fields", "", "comma separated fields")
	}
	cmd.Flags().StringVar(&flags.ids, "ids", "", "comma separated comment IDs")
	cmd.Flags().StringVar(&flags.title, "title", "", "filter by title")
	cmd.Flags().StringVar(&flags.description, "description", "", "filter by content")
	cmd.Flags().StringVar(&flags.author, "author", "", "filter by author")
	cmd.Flags().Int64Var(&flags.entryID, "entry-id", 0, "Lite workitem ID")
	cmd.Flags().StringVar(&flags.created, "created", "", "filter by created time expression")
	cmd.Flags().StringVar(&flags.modified, "modified", "", "filter by modified time expression")
	cmd.Flags().Int64Var(&flags.rootID, "root-id", 0, "root comment ID")
	cmd.Flags().Int64Var(&flags.replyID, "reply-id", 0, "reply comment ID")
}
