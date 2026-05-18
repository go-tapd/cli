package cmd

import (
	"strconv"

	"github.com/go-tapd/cli/internal/app"
	"github.com/go-tapd/tapd"
	"github.com/spf13/cobra"
)

func newAttachmentCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attachment",
		Short: "Work with TAPD attachments",
	}

	cmd.AddCommand(
		newAttachmentListCmd(rt),
		newAttachmentDownloadURLCmd(rt),
		newAttachmentImageURLCmd(rt),
		newAttachmentDocumentURLCmd(rt),
	)
	return cmd
}

func newAttachmentListCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		id          int
		entryID     int
		entryType   string
		filename    string
		owner       string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List attachments",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetAttachmentsRequest{
				WorkspaceID: tapd.Ptr(workspaceID),
			}
			if id > 0 {
				request.ID = tapd.Ptr(id)
			}
			if entryID > 0 {
				request.EntryID = tapd.Ptr(entryID)
			}
			if entryType != "" {
				request.Type = tapd.Ptr(entryType)
			}
			if filename != "" {
				request.Filename = tapd.Ptr(filename)
			}
			if owner != "" {
				request.Owner = tapd.Ptr(owner)
			}

			attachments, _, err := client.AttachmentService.GetAttachments(cmd.Context(), request)
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(attachments))
			for _, item := range attachments {
				rows = append(rows, []string{item.ID, item.Filename, item.Type, item.EntryID, item.Owner, item.Created})
			}

			return writeOutput(cmd, rt.OutputFormat, []string{"ID", "Filename", "Type", "EntryID", "Owner", "Created"}, rows, attachments)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().IntVar(&id, "id", 0, "attachment ID")
	cmd.Flags().IntVar(&entryID, "entry-id", 0, "attached object ID")
	cmd.Flags().StringVar(&entryType, "type", "", "attached object type")
	cmd.Flags().StringVar(&filename, "filename", "", "filter by filename")
	cmd.Flags().StringVar(&owner, "owner", "", "filter by uploader")
	return cmd
}

func newAttachmentDownloadURLCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		id          int
	)

	cmd := &cobra.Command{
		Use:   "download-url",
		Short: "Get an attachment download URL",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			attachment, _, err := client.AttachmentService.GetAttachmentDownloadURL(
				cmd.Context(),
				&tapd.GetAttachmentDownloadURLRequest{
					WorkspaceID: tapd.Ptr(workspaceID),
					ID:          tapd.Ptr(id),
				},
			)
			if err != nil {
				return err
			}

			rows := [][]string{{attachment.ID, attachment.Filename, attachment.DownloadURL}}
			return writeOutput(cmd, rt.OutputFormat, []string{"ID", "Filename", "DownloadURL"}, rows, attachment)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().IntVar(&id, "id", 0, "attachment ID")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func newAttachmentImageURLCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		imagePath   string
	)

	cmd := &cobra.Command{
		Use:   "image-url",
		Short: "Get an image download URL",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			attachment, _, err := client.AttachmentService.GetImageDownloadURL(
				cmd.Context(),
				&tapd.GetImageDownloadURLRequest{
					WorkspaceID: tapd.Ptr(workspaceID),
					ImagePath:   tapd.Ptr(imagePath),
				},
			)
			if err != nil {
				return err
			}

			rows := [][]string{{attachment.Filename, strconv.Itoa(attachment.WorkspaceID), attachment.DownloadURL}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Filename", "WorkspaceID", "DownloadURL"}, rows, attachment)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().StringVar(&imagePath, "image-path", "", "image path or URL")
	_ = cmd.MarkFlagRequired("image-path")
	return cmd
}

func newAttachmentDocumentURLCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		id          int
	)

	cmd := &cobra.Command{
		Use:   "document-url",
		Short: "Get a document download URL",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			document, _, err := client.AttachmentService.GetDocumentDownloadURL(
				cmd.Context(),
				&tapd.GetDocumentDownloadURLRequest{
					WorkspaceID: tapd.Ptr(workspaceID),
					ID:          tapd.Ptr(id),
				},
			)
			if err != nil {
				return err
			}

			rows := [][]string{{document.ID, document.Name, document.Type, document.DownloadURL}}
			return writeOutput(cmd, rt.OutputFormat, []string{"ID", "Name", "Type", "DownloadURL"}, rows, document)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().IntVar(&id, "id", 0, "document ID")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}
