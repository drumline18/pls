package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"pls/internal/app"
	"pls/internal/cli"
	"pls/internal/config"
	"pls/internal/configinit"
	"pls/internal/configshow"
	"pls/internal/doctor"
	"pls/internal/execute"
	"pls/internal/render"
	runtimeinfo "pls/internal/runtimeinfo"
)

func main() {
	exitCode := run(os.Args[1:])
	os.Exit(exitCode)
}

func run(args []string) int {
	parsed, err := cli.ParseArgs(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pls: %v\n", err)
		return 1
	}

	if parsed.Flags.PrintConfigPath {
		path, err := config.ResolvePath(parsed.Flags.ConfigPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "pls: %v\n", err)
			return 1
		}
		fmt.Fprintln(os.Stdout, path)
		return 0
	}

	if parsed.Help || len(parsed.RequestParts) == 0 {
		fmt.Fprint(os.Stdout, cli.HelpText)
		if len(parsed.RequestParts) == 0 && !parsed.Help {
			return 1
		}
		return 0
	}

	if len(parsed.RequestParts) == 1 && parsed.RequestParts[0] == "doctor" {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		report, err := doctor.Run(ctx, parsed.Flags)
		if err != nil {
			fmt.Fprintf(os.Stderr, "pls: %v\n", err)
			return 1
		}

		if parsed.Flags.JSON {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(report); err != nil {
				fmt.Fprintf(os.Stderr, "pls: %v\n", err)
				return 1
			}
			return 0
		}

		fmt.Fprintln(os.Stdout, doctor.Human(report))
		return 0
	}

	if len(parsed.RequestParts) == 1 && parsed.RequestParts[0] == "setup" {
		if err := configinit.Run(parsed.Flags); err != nil {
			fmt.Fprintf(os.Stderr, "pls: %v\n", err)
			return 1
		}
		return 0
	}

	if len(parsed.RequestParts) == 2 && parsed.RequestParts[0] == "config" && parsed.RequestParts[1] == "init" {
		if err := configinit.Run(parsed.Flags); err != nil {
			fmt.Fprintf(os.Stderr, "pls: %v\n", err)
			return 1
		}
		return 0
	}

	if len(parsed.RequestParts) == 3 && parsed.RequestParts[0] == "config" && parsed.RequestParts[1] == "local" && parsed.RequestParts[2] == "init" {
		if err := configinit.RunLocal(parsed.Flags); err != nil {
			fmt.Fprintf(os.Stderr, "pls: %v\n", err)
			return 1
		}
		return 0
	}

	if len(parsed.RequestParts) == 2 && parsed.RequestParts[0] == "config" && parsed.RequestParts[1] == "path" {
		path, err := config.ResolvePath(parsed.Flags.ConfigPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "pls: %v\n", err)
			return 1
		}
		fmt.Fprintln(os.Stdout, path)
		return 0
	}

	if len(parsed.RequestParts) == 2 && parsed.RequestParts[0] == "config" && parsed.RequestParts[1] == "show" {
		report, err := configshow.Run(parsed.Flags)
		if err != nil {
			fmt.Fprintf(os.Stderr, "pls: %v\n", err)
			return 1
		}
		if parsed.Flags.JSON {
			payload, err := configshow.JSON(report)
			if err != nil {
				fmt.Fprintf(os.Stderr, "pls: %v\n", err)
				return 1
			}
			fmt.Fprintln(os.Stdout, string(payload))
			return 0
		}
		fmt.Fprintln(os.Stdout, configshow.Human(report))
		return 0
	}

	cfg, err := config.Load(parsed.Flags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pls: %v\n", err)
		return 1
	}

	runtimeContext, err := runtimeinfo.Get(parsed.Flags.Shell, parsed.Flags.OS)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pls: %v\n", err)
		return 1
	}

	request := strings.TrimSpace(strings.Join(parsed.RequestParts, " "))
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	result, err := app.GenerateSuggestion(ctx, request, runtimeContext, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pls: %v\n", err)
		return 1
	}

	if cfg.OutputJSON {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(result); err != nil {
			fmt.Fprintf(os.Stderr, "pls: %v\n", err)
			return 1
		}
		return 0
	}

	fmt.Fprintln(os.Stdout, render.Human(result))

	executionFlags := parsed.Flags
	if cfg.YoloMode {
		executionFlags.Yes = true
	}

	runCommand, exitCode, err := execute.MaybePromptAndRun(result, runtimeContext, executionFlags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pls: %v\n", err)
		return 1
	}
	if runCommand {
		return exitCode
	}

	return 0
}
