package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/go-tapd/tapd"
)

func TestNormalizeAPIEndpoint(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "relative path", input: "tasks/get_fields_info", want: "tasks/get_fields_info"},
		{name: "leading slash", input: "/tasks", want: "tasks"},
		{name: "absolute URL", input: "https://example.com/tasks", wantErr: true},
		{name: "network path", input: "//example.com/tasks", wantErr: true},
		{name: "parent path", input: "../tasks", wantErr: true},
		{name: "encoded parent path", input: "%2e%2e/tasks", wantErr: true},
		{name: "query string", input: "tasks?workspace_id=1", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeAPIEndpoint(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("normalizeAPIEndpoint returned nil error")
				}
				return
			}
			if err != nil {
				t.Fatalf("normalizeAPIEndpoint returned error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("normalizeAPIEndpoint = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAPIFieldsData(t *testing.T) {
	got, err := apiFieldsData(
		[]string{"count=3", "enabled=true", "empty=null", "ratio=1.5"},
		[]string{"code=007"},
	)
	if err != nil {
		t.Fatalf("apiFieldsData returned error: %v", err)
	}

	if got["count"] != int64(3) {
		t.Fatalf("count = %#v, want int64(3)", got["count"])
	}
	if got["enabled"] != true {
		t.Fatalf("enabled = %#v, want true", got["enabled"])
	}
	if got["empty"] != nil {
		t.Fatalf("empty = %#v, want nil", got["empty"])
	}
	if got["ratio"] != 1.5 {
		t.Fatalf("ratio = %#v, want 1.5", got["ratio"])
	}
	if got["code"] != "007" {
		t.Fatalf("code = %#v, want 007 string", got["code"])
	}
}

func TestAPIRequestOptionsRejectsAuthorization(t *testing.T) {
	if _, err := apiRequestOptions([]string{"Authorization: Bearer token"}); err == nil {
		t.Fatal("apiRequestOptions returned nil error")
	}
}

func TestAppendAPIQuery(t *testing.T) {
	u, err := url.Parse("https://api.tapd.cn/tasks/get_fields_info")
	if err != nil {
		t.Fatalf("parse URL: %v", err)
	}

	err = appendAPIQuery(
		u,
		[]string{"workspace_id=123", "enabled=true", "empty=null"},
		[]string{"code=007"},
	)
	if err != nil {
		t.Fatalf("appendAPIQuery returned error: %v", err)
	}

	q := u.Query()
	if got := q.Get("workspace_id"); got != "123" {
		t.Fatalf("workspace_id query = %q, want 123", got)
	}
	if got := q.Get("enabled"); got != "true" {
		t.Fatalf("enabled query = %q, want true", got)
	}
	if got := q.Get("empty"); got != "null" {
		t.Fatalf("empty query = %q, want null", got)
	}
	if got := q.Get("code"); got != "007" {
		t.Fatalf("code query = %q, want 007", got)
	}
}

func TestAPIRequestDataFromInput(t *testing.T) {
	got, err := apiRequestData(
		strings.NewReader(`{"workspace_id":123,"name":"Task"}`),
		http.MethodPost,
		apiFlags{input: "-"},
	)
	if err != nil {
		t.Fatalf("apiRequestData returned error: %v", err)
	}

	raw, ok := got.(json.RawMessage)
	if !ok {
		t.Fatalf("apiRequestData = %T, want json.RawMessage", got)
	}
	if string(raw) != `{"workspace_id":123,"name":"Task"}` {
		t.Fatalf("raw JSON = %s", raw)
	}
}

func TestWriteAPIResponseIncludesHeaders(t *testing.T) {
	var stdout bytes.Buffer
	resp := &tapd.Response{Response: &http.Response{
		Proto:      "HTTP/1.1",
		Status:     "200 OK",
		StatusCode: http.StatusOK,
		Header: http.Header{
			"X-Test": []string{"value"},
		},
	}}

	err := writeAPIResponse(&stdout, resp, json.RawMessage(`{"ok":true}`), true, false)
	if err != nil {
		t.Fatalf("writeAPIResponse returned error: %v", err)
	}

	got := stdout.String()
	if !strings.Contains(got, "HTTP/1.1 200 OK") {
		t.Fatalf("stdout = %q, want status line", got)
	}
	if !strings.Contains(got, "X-Test: value") {
		t.Fatalf("stdout = %q, want header", got)
	}
	if !strings.Contains(got, `"ok": true`) {
		t.Fatalf("stdout = %q, want pretty JSON body", got)
	}
}
