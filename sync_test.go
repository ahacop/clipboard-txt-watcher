package main

import (
	"errors"
	"testing"
)

type mockClipboard struct {
	content      string
	readErr      error
	writeErr     error
	writeCalled  bool
	writeContent string
}

func (m *mockClipboard) Read() (string, error) {
	return m.content, m.readErr
}

func (m *mockClipboard) Write(content string) error {
	m.writeCalled = true
	m.writeContent = content
	return m.writeErr
}

func TestSyncToClipboard_UpdatesWhenDifferent(t *testing.T) {
	cb := &mockClipboard{content: "old content"}

	err := SyncToClipboard(cb, "new content")
	if err != nil {
		t.Fatalf("SyncToClipboard failed: %v", err)
	}

	if !cb.writeCalled {
		t.Error("expected Write to be called")
	}
	if cb.writeContent != "new content" {
		t.Errorf("expected writeContent %q, got %q", "new content", cb.writeContent)
	}
}

func TestSyncToClipboard_SkipsWhenSame(t *testing.T) {
	cb := &mockClipboard{content: "same content"}

	err := SyncToClipboard(cb, "same content")
	if err != nil {
		t.Fatalf("SyncToClipboard failed: %v", err)
	}

	if cb.writeCalled {
		t.Error("expected Write NOT to be called when content is same")
	}
}

func TestSyncToClipboard_ReturnsReadError(t *testing.T) {
	cb := &mockClipboard{readErr: errors.New("read failed")}

	err := SyncToClipboard(cb, "content")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSyncToClipboard_ReturnsWriteError(t *testing.T) {
	cb := &mockClipboard{
		content:  "old content",
		writeErr: errors.New("write failed"),
	}

	err := SyncToClipboard(cb, "new content")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
