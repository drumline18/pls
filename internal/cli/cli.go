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
}

const HelpText = `pls — natural-language shell command suggester

Usage:
  pls <request>
  pls --provider openai --model gpt-4.1-mini show hidden files here
  pls --provider ollama --model qwen2.5-coder:7b-instruct-q4_K_M list large files in this directory
  pls --json find files bigger than 500mb

Flags:
  --provider <openai|ollama>   LLM provider
  --model <name>               Model name to use
  --json                       Emit JSON only
  --shell <name>               Override shell detection
  --os <name>                  Override OS detection
  --host <url>                 Override provider host/base URL
  --help, -h                   Show help

Environment:
  OPENAI_API_KEY               OpenAI API key
  PLS_OPENAI_API_KEY           Alternate OpenAI API key env var
  OLLAMA_HOST                  Ollama base URL (default: http://127.0.0.1:11434)
  PLS_OLLAMA_HOST              Alternate Ollama host env var
  PLS_PROVIDER                 Default provider when flag is omitted
  PLS_MODEL                    Default model when flag is omitted
`

func ParseArgs(args []string) (types.ParsedArgs, error) {
	parsed := types.ParsedArgs{}

	for index := 0; index < len(args); index++ {
		arg := args[index]

		switch arg {
		case "--help", "-h":
			parsed.Help = true
			parsed.Flags.Help = true
			continue
		case "--json":
			parsed.Flags.JSON = true
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

		parsed.RequestParts = append(parsed.RequestParts, arg)
	}

	return parsed, nil
}
