package cmd

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
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
		newAttachmentUploadCmd(rt),
		newAttachmentUploadImageBase64Cmd(rt),
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
		limit       int
		page        int
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
			if limit > 0 {
				request.Limit = tapd.Ptr(limit)
			}
			if page > 0 {
				request.Page = tapd.Ptr(page)
			}

			attachments, _, err := client.AttachmentService.GetAttachments(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeAttachmentOutput(cmd, rt.OutputFormat, attachments, attachments)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().IntVar(&id, "id", 0, "attachment ID")
	cmd.Flags().IntVar(&entryID, "entry-id", 0, "attached object ID")
	cmd.Flags().StringVar(&entryType, "type", "", "attached object type")
	cmd.Flags().StringVar(&filename, "filename", "", "filter by filename")
	cmd.Flags().StringVar(&owner, "owner", "", "filter by uploader")
	newListFlags(cmd, &limit, &page)
	return cmd
}

func newAttachmentUploadCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		entryID     int64
		entryType   string
		customField string
		owner       string
		file        string
		filename    string
	)

	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload an attachment file",
		RunE: func(cmd *cobra.Command, _ []string) error {
			body, err := os.Open(file)
			if err != nil {
				return fmt.Errorf("open attachment file: %w", err)
			}
			defer func() {
				_ = body.Close()
			}()

			if filename == "" {
				filename = filepath.Base(file)
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.UploadAttachmentRequest{
				WorkspaceID: tapd.Ptr(workspaceID),
				Type:        tapd.Ptr(entryType),
				CustomField: tapd.Ptr(customField),
				EntryID:     tapd.Ptr(entryID),
				Filename:    tapd.Ptr(filename),
				File:        body,
			}
			if owner != "" {
				request.Owner = tapd.Ptr(owner)
			}

			attachment, _, err := client.AttachmentService.UploadAttachment(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeAttachmentOutput(cmd, rt.OutputFormat, []*tapd.Attachment{attachment}, attachment)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	newAttachmentUploadFlags(cmd, &entryID, &entryType, &customField, &owner)
	cmd.Flags().StringVar(&file, "file", "", "file to upload")
	cmd.Flags().StringVar(&filename, "filename", "", "uploaded filename; defaults to the file basename")
	_ = cmd.MarkFlagRequired("file")
	return cmd
}

func newAttachmentUploadImageBase64Cmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		entryID     int64
		entryType   string
		customField string
		owner       string
		imageFile   string
		base64Data  string
	)

	cmd := &cobra.Command{
		Use:   "upload-image-base64",
		Short: "Upload a base64-encoded image",
		RunE: func(cmd *cobra.Command, _ []string) error {
			data, err := attachmentImageBase64Data(imageFile, base64Data)
			if err != nil {
				return err
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.UploadImageBase64Request{
				WorkspaceID: tapd.Ptr(workspaceID),
				Base64Data:  tapd.Ptr(data),
				Type:        tapd.Ptr(entryType),
				CustomField: tapd.Ptr(customField),
				EntryID:     tapd.Ptr(entryID),
			}
			if owner != "" {
				request.Owner = tapd.Ptr(owner)
			}

			attachment, _, err := client.AttachmentService.UploadImageBase64(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeAttachmentOutput(cmd, rt.OutputFormat, []*tapd.Attachment{attachment}, attachment)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	newAttachmentUploadFlags(cmd, &entryID, &entryType, &customField, &owner)
	cmd.Flags().StringVar(&imageFile, "image-file", "", "image file to read and base64 encode")
	cmd.Flags().StringVar(&base64Data, "base64-data", "", "base64-encoded image data")
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

func newAttachmentUploadFlags(cmd *cobra.Command, entryID *int64, entryType, customField, owner *string) {
	cmd.Flags().Int64Var(entryID, "entry-id", 0, "attached object ID")
	cmd.Flags().StringVar(entryType, "type", "story_custom_field", "attached object type")
	cmd.Flags().StringVar(customField, "custom-field", "", "custom field name")
	cmd.Flags().StringVar(owner, "owner", "", "uploader account")
	_ = cmd.MarkFlagRequired("entry-id")
	_ = cmd.MarkFlagRequired("custom-field")
}

func attachmentImageBase64Data(imageFile, base64Data string) (string, error) {
	switch {
	case imageFile != "" && base64Data != "":
		return "", fmt.Errorf("use either --image-file or --base64-data, not both")
	case imageFile == "" && base64Data == "":
		return "", fmt.Errorf("one of --image-file or --base64-data is required")
	case imageFile != "":
		data, err := os.ReadFile(imageFile)
		if err != nil {
			return "", fmt.Errorf("read image file: %w", err)
		}
		return base64.StdEncoding.EncodeToString(data), nil
	default:
		return base64Data, nil
	}
}

func writeAttachmentOutput(cmd *cobra.Command, format string, attachments []*tapd.Attachment, data any) error {
	rows := make([][]string, 0, len(attachments))
	for _, item := range attachments {
		if item == nil {
			continue
		}
		rows = append(rows, []string{item.ID, item.Filename, item.Type, item.EntryID, item.Owner, item.Created})
	}

	return writeOutput(cmd, format, []string{"ID", "Filename", "Type", "EntryID", "Owner", "Created"}, rows, data)
}
