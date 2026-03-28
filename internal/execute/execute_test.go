package execute

import (
	"testing"

	"pls/internal/types"
)

func TestShellInvocationDefaultsToLoginShellStyle(t *testing.T) {
	name, args := shellInvocation("/bin/bash")
	if name != "/bin/bash" {
		t.Fatalf("unexpected shell name: %s", name)
	}
	if len(args) != 1 || args[0] != "-lc" {
		t.Fatalf("unexpected shell args: %#v", args)
	}
}

func TestShellInvocationHandlesFish(t *testing.T) {
	_, args := shellInvocation("/usr/bin/fish")
	if len(args) != 1 || args[0] != "-c" {
		t.Fatalf("unexpected shell args: %#v", args)
	}
}

func TestMaybePromptAndRunSkipsNonTTY(t *testing.T) {
	run, code, err := MaybePromptAndRun(types.Suggestion{Command: "echo hi", Explanation: "x", Risk: "low"}, types.RuntimeContext{IsTTY: false}, types.Flags{})
	if err != nil {
		t.Fatalf("MaybePromptAndRun returned error: %v", err)
	}
	if run {
		t.Fatalf("expected command not to run in non-TTY mode")
	}
	if code != 0 {
		t.Fatalf("unexpected exit code: %d", code)
	}
}

func TestMaybePromptAndRunSkipsWhenNoExecIsSet(t *testing.T) {
	run, code, err := MaybePromptAndRun(types.Suggestion{Command: "echo hi", Explanation: "x", Risk: "low"}, types.RuntimeContext{IsTTY: true}, types.Flags{NoExec: true})
	if err != nil {
		t.Fatalf("MaybePromptAndRun returned error: %v", err)
	}
	if run {
		t.Fatalf("expected command not to run when --no-exec is set")
	}
	if code != 0 {
		t.Fatalf("unexpected exit code: %d", code)
	}
}

func TestMaybePromptAndRunSkipsPromptWhenYesIsSetAndCommandIsSafe(t *testing.T) {
	run, _, err := MaybePromptAndRun(types.Suggestion{Command: "true", Explanation: "x", Risk: "low"}, types.RuntimeContext{IsTTY: false, ShellPath: "/bin/bash"}, types.Flags{Yes: true})
	if err != nil {
		t.Fatalf("MaybePromptAndRun returned error: %v", err)
	}
	if !run {
		t.Fatalf("expected command to run immediately when --yes is set for a safe command")
	}
}
