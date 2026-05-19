package app

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

func WriteJSON(w io.Writer, v any) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(v)
}

func WriteTable(w io.Writer, headers []string, rows [][]string) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	if len(headers) > 0 {
		if _, err := fmt.Fprintln(tw, joinRow(headers)); err != nil {
			return err
		}
	}
	for _, row := range rows {
		if _, err := fmt.Fprintln(tw, joinRow(row)); err != nil {
			return err
		}
	}

	return tw.Flush()
}

func joinRow(parts []string) string {
	var result strings.Builder
	for i, part := range parts {
		if i > 0 {
			result.WriteString("\t")
		}
		result.WriteString(part)
	}
	return result.String()
}
