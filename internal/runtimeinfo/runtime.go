package runtimeinfo

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/drumline18/pls/internal/types"
)

func Get(shellOverride, osOverride string) (types.RuntimeContext, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return types.RuntimeContext{}, err
	}

	shellPath := firstNonEmpty(shellOverride, os.Getenv("SHELL"), os.Getenv("ComSpec"), "unknown")
	homeDirectory, _ := os.UserHomeDir()

	return types.RuntimeContext{
		CWD:           cwd,
		OS:            normalizeOS(firstNonEmpty(osOverride, runtime.GOOS)),
		Shell:         normalizeShell(shellPath),
		ShellPath:     shellPath,
		HomeDirectory: homeDirectory,
		IsTTY:         isTTY(),
	}, nil
}

func normalizeOS(value string) string {
	switch strings.ToLower(value) {
	case "darwin":
		return "macos"
	case "linux":
		return "linux"
	case "windows", "win32":
		return "windows"
	default:
		return value
	}
}

func normalizeShell(shellPath string) string {
	lower := strings.ToLower(shellPath)
	switch {
	case strings.Contains(lower, "fish"):
		return "fish"
	case strings.Contains(lower, "zsh"):
		return "zsh"
	case strings.Contains(lower, "bash"):
		return "bash"
	case strings.Contains(lower, "powershell"), strings.Contains(lower, "pwsh"):
		return "powershell"
	case strings.Contains(lower, "cmd.exe"):
		return "cmd"
	default:
		return filepath.Base(shellPath)
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func isTTY() bool {
	info, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) != 0
}
