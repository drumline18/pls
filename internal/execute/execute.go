package execute

import (
	"bufio"
	"fmt"
	"io"
	"os"
	osexec "os/exec"
	"path/filepath"
	"strings"

	"pls/internal/types"
)

func MaybePromptAndRun(result types.Suggestion, runtimeContext types.RuntimeContext, flags types.Flags) (bool, int, error) {
	if flags.NoExec || result.Refused || result.NeedsClarification || strings.TrimSpace(result.Command) == "" {
		return false, 0, nil
	}

	if flags.Yes && !flags.NoExec && !result.RequiresConfirmation && !isHighRisk(result.Risk) {
		fmt.Fprint(os.Stdout, "\nRunning without prompt because auto-run mode is enabled.\n\n")
		exitCode, err := runCommand(result.Command, runtimeContext)
		return true, exitCode, err
	}

	if !canPrompt(runtimeContext) {
		return false, 0, nil
	}

	prompt := "\nRun it? [y/N] "
	if result.RequiresConfirmation || isHighRisk(result.Risk) {
		prompt = "\nHigh-risk command. Run it? [y/N] "
	}

	fmt.Fprint(os.Stdout, prompt)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return false, 0, err
	}

	answer := strings.ToLower(strings.TrimSpace(line))
	if answer != "y" && answer != "yes" {
		fmt.Fprintln(os.Stdout, "Not running.")
		return false, 0, nil
	}

	fmt.Fprintln(os.Stdout)
	exitCode, err := runCommand(result.Command, runtimeContext)
	return true, exitCode, err
}

func runCommand(command string, runtimeContext types.RuntimeContext) (int, error) {
	shellPath := resolveShellPath(runtimeContext)
	name, args := shellInvocation(shellPath)
	args = append(args, command)

	cmd := osexec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err == nil {
		return 0, nil
	}
	if exitErr, ok := err.(*osexec.ExitError); ok {
		return exitErr.ExitCode(), nil
	}
	return 1, err
}

func resolveShellPath(runtimeContext types.RuntimeContext) string {
	candidate := strings.TrimSpace(runtimeContext.ShellPath)
	if candidate == "" || candidate == "unknown" {
		candidate = strings.TrimSpace(os.Getenv("SHELL"))
	}
	if candidate == "" || candidate == "unknown" {
		candidate = "/bin/bash"
	}
	return candidate
}

func shellInvocation(shellPath string) (string, []string) {
	base := strings.ToLower(filepath.Base(shellPath))
	switch base {
	case "fish":
		return shellPath, []string{"-c"}
	case "pwsh", "powershell", "powershell.exe":
		return shellPath, []string{"-Command"}
	case "cmd", "cmd.exe":
		return shellPath, []string{"/C"}
	default:
		return shellPath, []string{"-lc"}
	}
}

func canPrompt(runtimeContext types.RuntimeContext) bool {
	if !runtimeContext.IsTTY {
		return false
	}
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) != 0
}

func isHighRisk(risk string) bool {
	switch strings.ToLower(strings.TrimSpace(risk)) {
	case "high", "critical":
		return true
	default:
		return false
	}
}
