package main

import (
	"errors"
	"testing"
)

func TestNewClipboard_ReturnsWaylandByDefault(t *testing.T) {
	cb := NewClipboard("wayland")

	_, ok := cb.(*WaylandClipboard)
	if !ok {
		t.Errorf("expected *WaylandClipboard, got %T", cb)
	}
}

func TestNewClipboard_ReturnsX11WhenSpecified(t *testing.T) {
	cb := NewClipboard("x11")

	_, ok := cb.(*X11Clipboard)
	if !ok {
		t.Errorf("expected *X11Clipboard, got %T", cb)
	}
}

func TestNewClipboard_ReturnsDarwinWhenSpecified(t *testing.T) {
	cb := NewClipboard("darwin")

	_, ok := cb.(*DarwinClipboard)
	if !ok {
		t.Errorf("expected *DarwinClipboard, got %T", cb)
	}
}

func TestDarwinClipboard_Read_CallsPbpaste(t *testing.T) {
	var calledCmd string
	var calledArgs []string

	cb := &DarwinClipboard{
		execCommand: func(cmd string, args ...string) ([]byte, error) {
			calledCmd = cmd
			calledArgs = args
			return []byte("darwin content"), nil
		},
	}

	content, err := cb.Read()
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if calledCmd != "pbpaste" {
		t.Errorf("expected command %q, got %q", "pbpaste", calledCmd)
	}
	if len(calledArgs) != 0 {
		t.Errorf("expected no args, got %v", calledArgs)
	}
	if content != "darwin content" {
		t.Errorf("expected content %q, got %q", "darwin content", content)
	}
}

func TestDarwinClipboard_Write_CallsPbcopy(t *testing.T) {
	var calledCmd string
	var stdinContent string

	cb := &DarwinClipboard{
		execCommandWithStdin: func(cmd string, stdin string, args ...string) error {
			calledCmd = cmd
			stdinContent = stdin
			return nil
		},
	}

	err := cb.Write("darwin test content")
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	if calledCmd != "pbcopy" {
		t.Errorf("expected command %q, got %q", "pbcopy", calledCmd)
	}
	if stdinContent != "darwin test content" {
		t.Errorf("expected stdin %q, got %q", "darwin test content", stdinContent)
	}
}

func TestWaylandClipboard_Read_CallsWlPaste(t *testing.T) {
	var calledCmd string
	var calledArgs []string

	cb := &WaylandClipboard{
		execCommand: func(cmd string, args ...string) ([]byte, error) {
			calledCmd = cmd
			calledArgs = args
			return []byte("clipboard content"), nil
		},
	}

	content, err := cb.Read()
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if calledCmd != "wl-paste" {
		t.Errorf("expected command %q, got %q", "wl-paste", calledCmd)
	}
	if len(calledArgs) != 1 || calledArgs[0] != "-n" {
		t.Errorf("expected args %v, got %v", []string{"-n"}, calledArgs)
	}
	if content != "clipboard content" {
		t.Errorf("expected content %q, got %q", "clipboard content", content)
	}
}

func TestWaylandClipboard_Write_CallsWlCopy(t *testing.T) {
	var calledCmd string
	var calledArgs []string
	var stdinContent string

	cb := &WaylandClipboard{
		execCommandWithStdin: func(cmd string, stdin string, args ...string) error {
			calledCmd = cmd
			calledArgs = args
			stdinContent = stdin
			return nil
		},
	}

	err := cb.Write("test content")
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	if calledCmd != "wl-copy" {
		t.Errorf("expected command %q, got %q", "wl-copy", calledCmd)
	}
	if stdinContent != "test content" {
		t.Errorf("expected stdin %q, got %q", "test content", stdinContent)
	}
	_ = calledArgs // no args expected for wl-copy
}

func TestX11Clipboard_Read_CallsXclip(t *testing.T) {
	var calledCmd string
	var calledArgs []string

	cb := &X11Clipboard{
		execCommand: func(cmd string, args ...string) ([]byte, error) {
			calledCmd = cmd
			calledArgs = args
			return []byte("x11 content"), nil
		},
	}

	content, err := cb.Read()
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if calledCmd != "xclip" {
		t.Errorf("expected command %q, got %q", "xclip", calledCmd)
	}
	expectedArgs := []string{"-selection", "clipboard", "-o"}
	if len(calledArgs) != len(expectedArgs) {
		t.Errorf("expected args %v, got %v", expectedArgs, calledArgs)
	}
	for i, arg := range expectedArgs {
		if calledArgs[i] != arg {
			t.Errorf("expected args %v, got %v", expectedArgs, calledArgs)
			break
		}
	}
	if content != "x11 content" {
		t.Errorf("expected content %q, got %q", "x11 content", content)
	}
}

func TestX11Clipboard_Write_CallsXclip(t *testing.T) {
	var calledCmd string
	var calledArgs []string
	var stdinContent string

	cb := &X11Clipboard{
		execCommandWithStdin: func(cmd string, stdin string, args ...string) error {
			calledCmd = cmd
			calledArgs = args
			stdinContent = stdin
			return nil
		},
	}

	err := cb.Write("x11 test content")
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	if calledCmd != "xclip" {
		t.Errorf("expected command %q, got %q", "xclip", calledCmd)
	}
	expectedArgs := []string{"-selection", "clipboard"}
	if len(calledArgs) != len(expectedArgs) {
		t.Errorf("expected args %v, got %v", expectedArgs, calledArgs)
	}
	for i, arg := range expectedArgs {
		if calledArgs[i] != arg {
			t.Errorf("expected args %v, got %v", expectedArgs, calledArgs)
			break
		}
	}
	if stdinContent != "x11 test content" {
		t.Errorf("expected stdin %q, got %q", "x11 test content", stdinContent)
	}
}

func TestNewClipboard_ReturnsWaylandForUnknownBackend(t *testing.T) {
	cb := NewClipboard("unknown")

	_, ok := cb.(*WaylandClipboard)
	if !ok {
		t.Errorf("expected *WaylandClipboard for unknown backend, got %T", cb)
	}
}

func TestWaylandClipboard_Read_ReturnsError(t *testing.T) {
	cb := &WaylandClipboard{
		execCommand: func(cmd string, args ...string) ([]byte, error) {
			return nil, errors.New("wl-paste failed")
		},
	}

	_, err := cb.Read()
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestWaylandClipboard_Write_ReturnsError(t *testing.T) {
	cb := &WaylandClipboard{
		execCommandWithStdin: func(cmd string, stdin string, args ...string) error {
			return errors.New("wl-copy failed")
		},
	}

	err := cb.Write("content")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestX11Clipboard_Read_ReturnsError(t *testing.T) {
	cb := &X11Clipboard{
		execCommand: func(cmd string, args ...string) ([]byte, error) {
			return nil, errors.New("xclip failed")
		},
	}

	_, err := cb.Read()
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestX11Clipboard_Write_ReturnsError(t *testing.T) {
	cb := &X11Clipboard{
		execCommandWithStdin: func(cmd string, stdin string, args ...string) error {
			return errors.New("xclip failed")
		},
	}

	err := cb.Write("content")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestDarwinClipboard_Read_ReturnsError(t *testing.T) {
	cb := &DarwinClipboard{
		execCommand: func(cmd string, args ...string) ([]byte, error) {
			return nil, errors.New("pbpaste failed")
		},
	}

	_, err := cb.Read()
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestDarwinClipboard_Write_ReturnsError(t *testing.T) {
	cb := &DarwinClipboard{
		execCommandWithStdin: func(cmd string, stdin string, args ...string) error {
			return errors.New("pbcopy failed")
		},
	}

	err := cb.Write("content")
	if err == nil {
		t.Error("expected error, got nil")
	}
}
