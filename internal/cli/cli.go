package cli

import (
	"fmt"
	"strings"

	"pls/internal/types"
)

var flagsWithValues = map[string]func(*types.Flags, string){
	"--provider": func(f *types.Flags, v string) { f.Provider = v },
	"--model":    func(f *types.Flags, v string) { f.Model = v },
	"--shell":    func(f *types.Flags, v string) { f.Shell = v },
	"--os":       func(f *types.Flags, v string) { f.OS = v },
	"--host":     func(f *types.Flags, v string) { f.Host = v },
	"--config":   func(f *types.Flags, v string) { f.ConfigPath = v },
}

const HelpText = `pls — natural-language shell command suggester

Usage:
  pls <request>
  pls doctor
  pls setup
  pls config init
  pls config local init
  pls config show
  pls config path
  pls --yes show hidden files here
  pls --no-exec prefix all jpgs with vacation-
  pls --provider openai --model gpt-4.1-mini show hidden files here
  pls --provider ollama --model qwen2.5-coder:7b-instruct-q4_K_M list large files in this directory
  pls --json find files bigger than 500mb
  pls --print-config-path
  pls -- show me files named --json

Built-ins:
  pls doctor                  Show environment/config/provider health
  pls setup                   Friendly alias for pls config init
  pls config init             Interactive global config setup wizard
  pls config local init       Interactive project-local override wizard
  pls config show             Show the effective config and active overrides
  pls config path             Print the global config path

Flags:
  --provider <name>            LLM provider (openai, ollama, anthropic, gemini, groq, deepseek, mistral, zai, llamacpp, llamafile)
  --model <name>               Model name to use
  --json                       Emit JSON only
  --yes                        Auto-run low/medium-risk commands without prompting
  --no-exec                    Never execute; suggestion only even in a TTY
  --shell <name>               Override shell detection
  --os <name>                  Override OS detection
  --host <url>                 Override provider host/base URL
  --config <path>              Override config file path
  --print-config-path          Print resolved config path and exit
  --help, -h                   Show help

Config precedence:
  flags > environment > local pls.json > global config > built-in defaults

Environment:
  OPENAI_API_KEY               OpenAI API key
  PLS_OPENAI_API_KEY           Alternate OpenAI API key env var
  OLLAMA_HOST                  Ollama base URL (default: http://127.0.0.1:11434)
  PLS_OLLAMA_HOST              Alternate Ollama host env var
  PLS_HOST                     Generic provider host override
  PLS_PROVIDER                 Default provider when flag is omitted
  PLS_MODEL                    Default model when flag is omitted
  PLS_YOLO_MODE                true/false auto-run default for safe commands
  PLS_CONFIG                   Override config file path
`

func ParseArgs(args []string) (types.ParsedArgs, error) {
	parsed := types.ParsedArgs{}
	parsingFlags := true

	for index := 0; index < len(args); index++ {
		arg := args[index]

		if parsingFlags {
			switch arg {
			case "--":
				parsingFlags = false
				continue
			case "--help", "-h":
				parsed.Help = true
				parsed.Flags.Help = true
				continue
			case "--json":
				parsed.Flags.JSON = true
				continue
			case "--yes":
				parsed.Flags.Yes = true
				continue
			case "--no-exec":
				parsed.Flags.NoExec = true
				continue
			case "--print-config-path":
				parsed.Flags.PrintConfigPath = true
				continue
			}

			if setter, ok := flagsWithValues[arg]; ok {
				if index+1 >= len(args) || strings.HasPrefix(args[index+1], "--") {
					return types.ParsedArgs{}, fmt.Errorf("missing value for %s", arg)
				}
				setter(&parsed.Flags, args[index+1])
				index++
				continue
			}

			parsingFlags = false
		}

		parsed.RequestParts = append(parsed.RequestParts, arg)
	}

	return parsed, nil
}
