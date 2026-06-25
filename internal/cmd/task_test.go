package cmd

import (
	"strings"
	"testing"

	"github.com/go-tapd/tapd"
)

func TestApplyTaskCreateFlagsSetsCustomFields(t *testing.T) {
	request := new(tapd.CreateTaskRequest)

	err := applyTaskCreateFlags(request, taskMutationFlags{
		customFields: []string{
			"custom_field_one=开发阶段",
			"custom_field_9=phase-9",
		},
	})
	if err != nil {
		t.Fatalf("applyTaskCreateFlags returned error: %v", err)
	}

	if request.CustomFieldOne == nil || *request.CustomFieldOne != "开发阶段" {
		t.Fatalf("CustomFieldOne = %v, want 开发阶段", request.CustomFieldOne)
	}
	if request.CustomField9 == nil || *request.CustomField9 != "phase-9" {
		t.Fatalf("CustomField9 = %v, want phase-9", request.CustomField9)
	}
}

func TestApplyTaskCustomFieldsRejectsInvalidAssignment(t *testing.T) {
	err := applyTaskCustomFields(new(tapd.CreateTaskRequest), []string{"custom_field_one"})
	if err == nil {
		t.Fatal("applyTaskCustomFields returned nil error")
	}
	if !strings.Contains(err.Error(), "expected key=value") {
		t.Fatalf("error = %q, want key=value message", err)
	}
}

func TestApplyTaskCustomFieldsRejectsUnknownField(t *testing.T) {
	err := applyTaskCustomFields(new(tapd.CreateTaskRequest), []string{"custom_field_51=value"})
	if err == nil {
		t.Fatal("applyTaskCustomFields returned nil error")
	}
	if !strings.Contains(err.Error(), "unsupported task custom field") {
		t.Fatalf("error = %q, want unsupported field message", err)
	}
}
