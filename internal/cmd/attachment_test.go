package cmd

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"
)

func TestAttachmentImageBase64Data(t *testing.T) {
	t.Run("uses direct base64 data", func(t *testing.T) {
		got, err := attachmentImageBase64Data("", "encoded")
		if err != nil {
			t.Fatalf("attachmentImageBase64Data returned error: %v", err)
		}
		if got != "encoded" {
			t.Fatalf("attachmentImageBase64Data = %q, want %q", got, "encoded")
		}
	})

	t.Run("encodes image file", func(t *testing.T) {
		dir := t.TempDir()
		file := filepath.Join(dir, "image.png")
		if err := os.WriteFile(file, []byte("image"), 0o600); err != nil {
			t.Fatalf("write image file: %v", err)
		}

		got, err := attachmentImageBase64Data(file, "")
		if err != nil {
			t.Fatalf("attachmentImageBase64Data returned error: %v", err)
		}
		want := base64.StdEncoding.EncodeToString([]byte("image"))
		if got != want {
			t.Fatalf("attachmentImageBase64Data = %q, want %q", got, want)
		}
	})

	t.Run("rejects missing input", func(t *testing.T) {
		if _, err := attachmentImageBase64Data("", ""); err == nil {
			t.Fatal("attachmentImageBase64Data returned nil error")
		}
	})

	t.Run("rejects conflicting input", func(t *testing.T) {
		if _, err := attachmentImageBase64Data("image.png", "encoded"); err == nil {
			t.Fatal("attachmentImageBase64Data returned nil error")
		}
	})
}
