package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-tapd/cli/internal/app"
	"github.com/go-tapd/tapd"
	"github.com/spf13/cobra"
)

func requireTableOrJSON(format string) error {
	switch format {
	case "table", "json":
		return nil
	default:
		return fmt.Errorf("unsupported format %q, expected table or json", format)
	}
}

func writeOutput(cmd *cobra.Command, format string, headers []string, rows [][]string, data any) error {
	if err := requireTableOrJSON(format); err != nil {
		return err
	}

	switch format {
	case "json":
		return app.WriteJSON(cmd.OutOrStdout(), data)
	default:
		return app.WriteTable(cmd.OutOrStdout(), headers, rows)
	}
}

func newListFlags(cmd *cobra.Command, limit, page *int) {
	cmd.Flags().IntVar(limit, "limit", 30, "page size")
	cmd.Flags().IntVar(page, "page", 1, "page number")
}

func newWorkspaceFlag(cmd *cobra.Command, workspaceID *int) {
	cmd.Flags().IntVarP(workspaceID, "workspace-id", "w", 0, "workspace ID")
	_ = cmd.MarkFlagRequired("workspace-id")
}

func fieldsMulti(fields string) *tapd.Multi[string] {
	items := splitCSV(fields)
	if len(items) == 0 {
		return nil
	}
	return tapd.NewMulti(items...)
}

func int64Multi(csv string) *tapd.Multi[int64] {
	items := splitCSV(csv)
	if len(items) == 0 {
		return nil
	}

	values := make([]int64, 0, len(items))
	for _, item := range items {
		if item == "" {
			continue
		}
		v, err := strconv.ParseInt(item, 10, 64)
		if err != nil {
			continue
		}
		values = append(values, v)
	}
	if len(values) == 0 {
		return nil
	}

	return tapd.NewMulti(values...)
}

func splitCSV(csv string) []string {
	if csv == "" {
		return nil
	}
	rawItems := strings.Split(csv, ",")
	items := make([]string, 0, len(rawItems))
	for _, item := range rawItems {
		item = strings.TrimSpace(item)
		if item != "" {
			items = append(items, item)
		}
	}
	return items
}
