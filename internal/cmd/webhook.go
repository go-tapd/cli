package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-tapd/cli/internal/app"
	tapdwebhook "github.com/go-tapd/tapd/webhook"
	"github.com/spf13/cobra"
)

func newWebhookCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "webhook",
		Short: "Work with TAPD webhook payloads",
	}

	cmd.AddCommand(
		newWebhookServeCmd(rt),
		newWebhookValidateCmd(rt),
		newWebhookInspectCmd(rt),
	)
	return cmd
}

func newWebhookServeCmd(rt *app.Runtime) *cobra.Command {
	var (
		addr         string
		path         string
		maxBodyBytes int64
	)

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve a local TAPD webhook receiver",
		RunE: func(cmd *cobra.Command, _ []string) error {
			mux := http.NewServeMux()
			mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
				if req.Method != http.MethodPost {
					w.Header().Set("Allow", http.MethodPost)
					http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
					return
				}

				payload, err := io.ReadAll(http.MaxBytesReader(w, req.Body, maxBodyBytes))
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				summary, err := parseWebhookPayload(payload)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				logWebhookSummary(cmd.OutOrStdout(), summary)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"ok":           true,
					"event":        summary.Event,
					"workspace_id": summary.WorkspaceID,
					"id":           summary.ID,
				})
			})

			server := &http.Server{
				Addr:              addr,
				Handler:           mux,
				ReadHeaderTimeout: 5 * time.Second,
			}
			errCh := make(chan error, 1)
			go func() {
				errCh <- server.ListenAndServe()
			}()

			fmt.Fprintf(cmd.ErrOrStderr(), "Listening on http://%s%s\n", addr, path)
			select {
			case <-cmd.Context().Done():
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				return server.Shutdown(ctx)
			case err := <-errCh:
				if errors.Is(err, http.ErrServerClosed) {
					return nil
				}
				return err
			}
		},
	}

	cmd.Flags().StringVar(&addr, "addr", "127.0.0.1:8080", "listen address")
	cmd.Flags().StringVar(&path, "path", "/webhook", "webhook request path")
	cmd.Flags().Int64Var(&maxBodyBytes, "max-body-bytes", 1<<20, "maximum request body size")
	return cmd
}

func newWebhookValidateCmd(rt *app.Runtime) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate a TAPD webhook payload",
		RunE: func(cmd *cobra.Command, _ []string) error {
			payload, err := readWebhookPayload(cmd, file)
			if err != nil {
				return err
			}

			summary, err := parseWebhookPayload(payload)
			if err != nil {
				return err
			}

			data := map[string]any{
				"valid":        true,
				"event":        summary.Event,
				"workspace_id": summary.WorkspaceID,
				"id":           summary.ID,
			}
			rows := [][]string{{"true", summary.Event, summary.WorkspaceID, summary.ID}}
			return writeOutput(cmd, rt.OutputFormat, []string{"Valid", "Event", "WorkspaceID", "ID"}, rows, data)
		},
	}

	cmd.Flags().StringVar(&file, "file", "-", "payload file, or - for stdin")
	return cmd
}

func newWebhookInspectCmd(rt *app.Runtime) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "inspect",
		Short: "Inspect a TAPD webhook payload",
		RunE: func(cmd *cobra.Command, _ []string) error {
			payload, err := readWebhookPayload(cmd, file)
			if err != nil {
				return err
			}

			summary, err := parseWebhookPayload(payload)
			if err != nil {
				return err
			}

			data := map[string]any{
				"summary": summary,
				"payload": summary.Raw,
			}
			rows := [][]string{{
				summary.Event,
				summary.WorkspaceID,
				summary.ID,
				summary.CurrentUser,
				summary.EventID,
				summary.Created,
				summary.Subject,
			}}
			return writeOutput(
				cmd,
				rt.OutputFormat,
				[]string{"Event", "WorkspaceID", "ID", "CurrentUser", "EventID", "Created", "Subject"},
				rows,
				data,
			)
		},
	}

	cmd.Flags().StringVar(&file, "file", "-", "payload file, or - for stdin")
	return cmd
}

type webhookPayloadSummary struct {
	Event       string         `json:"event,omitempty"`
	WorkspaceID string         `json:"workspace_id,omitempty"`
	ID          string         `json:"id,omitempty"`
	CurrentUser string         `json:"current_user,omitempty"`
	EventID     string         `json:"event_id,omitempty"`
	Created     string         `json:"created,omitempty"`
	Subject     string         `json:"subject,omitempty"`
	Raw         map[string]any `json:"-"`
}

func readWebhookPayload(cmd *cobra.Command, file string) ([]byte, error) {
	if file == "" || file == "-" {
		return io.ReadAll(cmd.InOrStdin())
	}
	return os.ReadFile(file)
}

func parseWebhookPayload(payload []byte) (*webhookPayloadSummary, error) {
	eventType, _, err := tapdwebhook.ParseWebhookEvent(payload)
	if err != nil {
		return nil, err
	}

	var raw map[string]any
	if err := json.Unmarshal(payload, &raw); err != nil {
		return nil, err
	}

	return &webhookPayloadSummary{
		Event:       eventType.String(),
		WorkspaceID: stringField(raw, "workspace_id"),
		ID:          stringField(raw, "id"),
		CurrentUser: stringField(raw, "current_user"),
		EventID:     stringField(raw, "event_id"),
		Created:     stringField(raw, "created"),
		Subject:     webhookSubject(raw),
		Raw:         raw,
	}, nil
}

func logWebhookSummary(w io.Writer, summary *webhookPayloadSummary) {
	fmt.Fprintf(
		w,
		"event=%s workspace_id=%s id=%s current_user=%s event_id=%s subject=%q\n",
		summary.Event,
		summary.WorkspaceID,
		summary.ID,
		summary.CurrentUser,
		summary.EventID,
		summary.Subject,
	)
}

func webhookSubject(raw map[string]any) string {
	for _, key := range []string{"name", "title", "old_name", "old_title", "description"} {
		if value := stringField(raw, key); value != "" {
			return value
		}
	}
	return ""
}

func stringField(raw map[string]any, key string) string {
	value, ok := raw[key]
	if !ok || value == nil {
		return ""
	}
	switch v := value.(type) {
	case string:
		return v
	case json.Number:
		return v.String()
	default:
		return fmt.Sprint(v)
	}
}
