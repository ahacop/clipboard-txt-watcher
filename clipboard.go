package main

import (
	"os/exec"
	"strings"
)

type (
	CommandExecutor          func(cmd string, args ...string) ([]byte, error)
	CommandWithStdinExecutor func(cmd string, stdin string, args ...string) error
)

func defaultExec(cmd string, args ...string) ([]byte, error) {
	return exec.Command(cmd, args...).Output()
}

func defaultExecWithStdin(cmd string, stdin string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Stdin = strings.NewReader(stdin)
	return c.Run()
}

type Clipboard interface {
	Read() (string, error)
	Write(content string) error
}

type WaylandClipboard struct {
	execCommand          CommandExecutor
	execCommandWithStdin CommandWithStdinExecutor
}

func (w *WaylandClipboard) Read() (string, error) {
	executor := w.execCommand
	if executor == nil {
		executor = defaultExec
	}
	out, err := executor("wl-paste", "-n")
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func (w *WaylandClipboard) Write(content string) error {
	executor := w.execCommandWithStdin
	if executor == nil {
		executor = defaultExecWithStdin
	}
	return executor("wl-copy", content)
}

type X11Clipboard struct {
	execCommand          CommandExecutor
	execCommandWithStdin CommandWithStdinExecutor
}

type DarwinClipboard struct {
	execCommand          CommandExecutor
	execCommandWithStdin CommandWithStdinExecutor
}

func (d *DarwinClipboard) Read() (string, error) {
	executor := d.execCommand
	if executor == nil {
		executor = defaultExec
	}
	out, err := executor("pbpaste")
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func (d *DarwinClipboard) Write(content string) error {
	executor := d.execCommandWithStdin
	if executor == nil {
		executor = defaultExecWithStdin
	}
	return executor("pbcopy", content)
}

func (x *X11Clipboard) Read() (string, error) {
	executor := x.execCommand
	if executor == nil {
		executor = defaultExec
	}
	out, err := executor("xclip", "-selection", "clipboard", "-o")
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func (x *X11Clipboard) Write(content string) error {
	executor := x.execCommandWithStdin
	if executor == nil {
		executor = defaultExecWithStdin
	}
	return executor("xclip", content, "-selection", "clipboard")
}

func NewClipboard(backend string) Clipboard {
	switch backend {
	case "x11":
		return &X11Clipboard{}
	case "darwin":
		return &DarwinClipboard{}
	default:
		return &WaylandClipboard{}
	}
}
