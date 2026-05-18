package cmd

import (
	"fmt"
	"strconv"

	"github.com/go-tapd/cli/internal/app"
	"github.com/go-tapd/tapd"
	"github.com/spf13/cobra"
)

func newSourceCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "source",
		Short: "Work with TAPD source code integrations",
	}

	cmd.AddCommand(newSourceCommitCmd(rt))
	return cmd
}

func newSourceCommitCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commit",
		Short: "Work with source commits",
	}

	cmd.AddCommand(
		newSourceCommitAddCmd(rt),
		newSourceCommitListCmd(rt),
		newSourceCommitObjectsCmd(rt),
	)
	return cmd
}

func newSourceCommitAddCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		commitID    string
		author      string
		message     string
		files       string
		repo        string
		repoID      string
		commitTime  string
		gitEnv      string
		repoURL     string
		commitURL   string
	)

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add source commit information",
		RunE: func(cmd *cobra.Command, _ []string) error {
			fileItems := splitCSV(files)
			if len(fileItems) == 0 {
				return fmt.Errorf("files cannot be empty")
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.AddCodeCommitInfoRequest{
				WorkspaceID: tapd.Ptr(workspaceID),
				CommitID:    tapd.Ptr(commitID),
				Author:      tapd.Ptr(author),
				Message:     tapd.Ptr(message),
				Files:       tapd.Ptr(fileItems),
				Repo:        tapd.Ptr(repo),
				RepoID:      tapd.Ptr(repoID),
				CommitTime:  tapd.Ptr(commitTime),
			}
			if gitEnv != "" {
				request.GitEnv = tapd.Ptr(gitEnv)
			}
			if repoURL != "" {
				request.RepoURL = tapd.Ptr(repoURL)
			}
			if commitURL != "" {
				request.CommitURL = tapd.Ptr(commitURL)
			}

			info, _, err := client.SourceService.AddCodeCommitInfo(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, sourceCommitHeaders(), sourceCommitRows([]*tapd.CodeCommitInfo{info}), info)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	cmd.Flags().StringVar(&commitID, "commit-id", "", "commit ID")
	cmd.Flags().StringVar(&author, "author", "", "commit author")
	cmd.Flags().StringVar(&message, "message", "", "commit message")
	cmd.Flags().StringVar(&files, "files", "", "comma separated changed files")
	cmd.Flags().StringVar(&repo, "repo", "", "repository name")
	cmd.Flags().StringVar(&repoID, "repo-id", "", "repository ID")
	cmd.Flags().StringVar(&commitTime, "commit-time", "", "commit time")
	cmd.Flags().StringVar(&gitEnv, "git-env", "", "git environment, such as github or gitlab")
	cmd.Flags().StringVar(&repoURL, "repo-url", "", "repository URL")
	cmd.Flags().StringVar(&commitURL, "commit-url", "", "commit URL")
	_ = cmd.MarkFlagRequired("commit-id")
	_ = cmd.MarkFlagRequired("author")
	_ = cmd.MarkFlagRequired("message")
	_ = cmd.MarkFlagRequired("files")
	_ = cmd.MarkFlagRequired("repo")
	_ = cmd.MarkFlagRequired("repo-id")
	_ = cmd.MarkFlagRequired("commit-time")
	return cmd
}

func newSourceCommitListCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		entityType  string
		objectID    int64
		commitTime  string
		relatedType string
		limit       int
		page        int
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List source commit information",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetCodeCommitInfosRequest{
				WorkspaceID: tapd.Ptr(workspaceID),
				Type:        tapd.Ptr(tapd.EntityType(entityType)),
				ObjectID:    tapd.Ptr(objectID),
				Limit:       tapd.Ptr(limit),
				Page:        tapd.Ptr(page),
			}
			if commitTime != "" {
				request.CommitTime = tapd.Ptr(commitTime)
			}
			if relatedType != "" {
				request.RelatedType = tapd.Ptr(tapd.CodeCommitRelatedType(relatedType))
			}

			infos, _, err := client.SourceService.GetCodeCommitInfos(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, sourceCommitHeaders(), sourceCommitRows(infos), infos)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	newListFlags(cmd, &limit, &page)
	cmd.Flags().StringVar(&entityType, "entity-type", "", "entity type: story, bug, or task")
	cmd.Flags().Int64Var(&objectID, "object-id", 0, "TAPD object ID")
	cmd.Flags().StringVar(&commitTime, "commit-time", "", "filter by commit time expression")
	cmd.Flags().StringVar(&relatedType, "related-type", "", "related type: all, branch, or source_code")
	_ = cmd.MarkFlagRequired("entity-type")
	_ = cmd.MarkFlagRequired("object-id")
	return cmd
}

func newSourceCommitObjectsCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		commitIDs   string
		entityType  string
		scmType     string
		limit       int
		page        int
		fields      string
	)

	cmd := &cobra.Command{
		Use:   "objects",
		Short: "List TAPD objects related to commits",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ids, err := strictStringMulti("commit IDs", commitIDs)
			if err != nil {
				return err
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetCommitObjectsRequest{
				WorkspaceID: tapd.Ptr(workspaceID),
				CommitID:    ids,
				EntityType:  tapd.Ptr(tapd.EntityType(entityType)),
				Limit:       tapd.Ptr(limit),
				Page:        tapd.Ptr(page),
				Fields:      fieldsMulti(fields),
			}
			if scmType != "" {
				request.SCMType = tapd.Ptr(scmType)
			}

			objects, _, err := client.SourceService.GetCommitObjects(cmd.Context(), request)
			if err != nil {
				return err
			}

			rows := make([][]string, 0, len(objects))
			for _, item := range objects {
				rows = append(rows, []string{
					objectID(item),
					objectName(item),
					objectType(item),
				})
			}

			return writeOutput(cmd, rt.OutputFormat, []string{"ID", "Name", "Type"}, rows, objects)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	newListFlags(cmd, &limit, &page)
	cmd.Flags().StringVar(&commitIDs, "commit-ids", "", "comma separated commit IDs")
	cmd.Flags().StringVar(&entityType, "entity-type", "", "entity type: story, bug, or task")
	cmd.Flags().StringVar(&scmType, "scm-type", "", "SCM type")
	cmd.Flags().StringVar(&fields, "fields", "", "comma separated fields")
	_ = cmd.MarkFlagRequired("commit-ids")
	_ = cmd.MarkFlagRequired("entity-type")
	return cmd
}

func sourceCommitHeaders() []string {
	return []string{"ID", "CommitID", "Author", "Repo", "CommitTime", "Related"}
}

func sourceCommitRows(infos []*tapd.CodeCommitInfo) [][]string {
	rows := make([][]string, 0, len(infos))
	for _, item := range infos {
		rows = append(rows, []string{
			item.ID,
			item.CommitID,
			item.HookUserName,
			item.HookProjectName,
			item.CommitTime,
			strconv.Itoa(len(item.Related)),
		})
	}
	return rows
}

func objectID(item *tapd.CommitObject) string {
	switch {
	case item.Story != nil:
		return item.Story.ID
	case item.Bug != nil:
		return item.Bug.ID
	case item.Task != nil:
		return item.Task.ID
	default:
		return ""
	}
}

func objectName(item *tapd.CommitObject) string {
	switch {
	case item.Story != nil:
		return item.Story.Name
	case item.Bug != nil:
		return item.Bug.Title
	case item.Task != nil:
		return item.Task.Name
	default:
		return ""
	}
}

func objectType(item *tapd.CommitObject) string {
	switch {
	case item.Story != nil:
		return "story"
	case item.Bug != nil:
		return "bug"
	case item.Task != nil:
		return "task"
	default:
		return ""
	}
}
