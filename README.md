# pls

`pls` is a natural-language shell command suggester written in Go.

Status:
- actively being prepared for open source and distribution
- planned public repo: `github.com/drumline18/pls`
- not published yet
- module path and MIT license are set; public `go install` will work once the repo exists publicly

Current MVP goals:
- command generation first, with optional confirmed execution
- no quotes required around the prompt
- Linux-first, shell-aware
- platform-aware prompting for macOS and Windows PowerShell
- broader provider support via any-llm-go
- explain command + risk level
- JSON output for future integrations

## Examples

```bash
pls doctor
pls setup
pls config init
pls config local init
pls config show
pls config path
pls show me all dotfiles in this directory
pls --yes show hidden files here
pls --no-exec prefix all jpgs with vacation-
pls find files bigger than 500mb
pls check if jellyfin is running
pls prefix all the mp3s with their lengths in seconds
pls prefix all jpgs with vacation-
pls replace spaces in all filenames here with underscores
pls move all srt files into a subtitles folder
pls --provider openai --model gpt-4.1-mini why is port 3000 busy
pls --provider ollama --model qwen2.5-coder:7b-instruct-q4_K_M show hidden files here
pls --json list the 10 biggest files under the current directory
pls -- show me files named --json
```

## Platform support

Current support status:
- Linux: **primary**
- macOS: **beta**
- Windows PowerShell: **beta**
- Windows cmd.exe: **limited**

Right now Linux has the strongest normalization and post-processing rules. macOS and PowerShell now get platform-aware prompt examples and instructions, but they still have less hand-tuned coverage than Linux.

## Config file

Default path:

```bash
~/.config/pls/config.json
```

First-time setup wizard:

```bash
pls setup
# or
pls config init
```

Built-in config commands:

```bash
pls config show
pls config path
pls config local init
```

`pls setup`, `pls config init`, and `pls config local init` are treated as strict built-in commands. Longer phrases like `pls config init my project` or `pls setup my repo` still go through the normal natural-language path.

You can print the resolved config path with:

```bash
pls --print-config-path
```

You can override the config path with either:

```bash
pls --config /path/to/config.json ...
```

or:

```bash
export PLS_CONFIG=/path/to/config.json
```

Config precedence is:

```text
flags > environment > local pls.json > global config > built-in defaults
```

Global config example (`~/.config/pls/config.json`):

```json
{
  "provider": "ollama",
  "model": "qwen2.5-coder:7b-instruct-q4_K_M",
  "host": "http://192.168.2.166:11434",
  "yoloMode": false
}
```

Local override example (`pls.json` in a project directory):

```json
{
  "yoloMode": true
}
```

`pls` looks for `pls.json` in the current directory, then walks upward until it finds one.

A copyable example also lives at `examples/config.example.json`.

## Provider configuration

Currently wired providers:
- Ollama
- OpenAI
- Anthropic
- Gemini
- Groq
- DeepSeek
- Mistral
- Z.ai
- llama.cpp server
- llamafile server

The setup wizard now covers all currently wired providers.

Behavior:
- the global wizard can set provider/model/host for any supported provider
- OpenAI can optionally store an API key in global config
- non-OpenAI hosted providers use their normal environment variable for credentials
- the local wizard can override provider/model/host/yolo mode, but never stores API keys locally

### Ollama

Defaults:
- provider: `ollama` if no OpenAI API key is present
- host: `http://127.0.0.1:11434`
- model: `qwen2.5-coder:7b-instruct-q4_K_M`

Environment variables:

```bash
export OLLAMA_HOST=http://127.0.0.1:11434
export PLS_OLLAMA_HOST=http://127.0.0.1:11434
export PLS_PROVIDER=ollama
export PLS_MODEL=qwen2.5-coder:7b-instruct-q4_K_M
```

### OpenAI

Environment variables:

```bash
export OPENAI_API_KEY=your_key_here
export PLS_PROVIDER=openai
export PLS_MODEL=gpt-4.1-mini
```

### Additional providers

Examples:

```bash
export PLS_PROVIDER=anthropic
export PLS_MODEL=claude-3-5-haiku-latest
export ANTHROPIC_API_KEY=your_key_here
```

```bash
export PLS_PROVIDER=gemini
export PLS_MODEL=gemini-2.5-flash
export GEMINI_API_KEY=your_key_here
```

```bash
export PLS_PROVIDER=groq
export PLS_MODEL=llama-3.3-70b-versatile
export GROQ_API_KEY=your_key_here
```

```bash
export PLS_PROVIDER=deepseek
export PLS_MODEL=deepseek-chat
export DEEPSEEK_API_KEY=your_key_here
```

For local OpenAI-compatible servers, `provider=openai` with `host=<base-url>` is still useful.

### Execution / yolo mode

Optional environment override:

```bash
export PLS_YOLO_MODE=true
```

Accepted values are: `true`, `false`, `yes`, `no`, `on`, `off`, `1`, `0`.

## Toolchain

`pls` currently targets **Go 1.25+** because of the provider stack behind `any-llm-go`.

If your default `go` is older, either use a newer Go directly:

```bash
~/.local/bin/go1.26 test ./...
~/.local/bin/go1.26 build -o bin/pls ./cmd/pls
```

or point the Makefile at it:

```bash
cd pls
GO=~/.local/bin/go1.26 make test
GO=~/.local/bin/go1.26 make build
```

## Build

```bash
cd pls
make build
```

Or directly:

```bash
go build -o bin/pls ./cmd/pls
```

## Install into PATH

Recommended local install:

```bash
cd pls
make install
```

That installs `pls` to:

```bash
~/.local/bin/pls
```

If `~/.local/bin` is not already in your PATH, add this to your shell rc file:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

Alternative:

```bash
GOBIN=$HOME/.local/bin go install ./cmd/pls
```

If you are using a non-default Go binary, the same install works with that toolchain too:

```bash
GOBIN=$HOME/.local/bin ~/.local/bin/go1.26 install ./cmd/pls
```

Planned public install command once the repo is live:

```bash
go install github.com/drumline18/pls/cmd/pls@latest
```

## Distribution status

Planned public distribution path:
- GitHub Releases archives via GoReleaser
- `go install github.com/drumline18/pls/cmd/pls@latest`
- Homebrew via a future `drumline18/homebrew-tap`
- Scoop via a future `drumline18/scoop-bucket`

Already prepared:
- public module path in `go.mod`
- MIT `LICENSE`
- GoReleaser release archives/checksums
- packaging templates under `packaging/`

Still not done:
- actual public repo push/publish
- first tagged release artifacts
- live Homebrew tap / Scoop bucket repos

See `RELEASE.md` for the current pre-publish checklist.

## Run

```bash
pls show hidden files here
```

In an interactive terminal, `pls` now prompts before execution:

```text
Run it? [y/N]
```

Execution flags fit cleanly with no quotes because `pls` only parses **leading** known flags. As soon as the natural-language request starts, the rest is treated as request text.

Examples:

```bash
pls --yes show hidden files here
pls --no-exec prefix all jpgs with vacation-
pls show me files named --json
pls -- show me files named --json
```

Behavior:
- `--yes` auto-runs low/medium-risk commands without prompting
- `yoloMode: true` acts like a config-backed `--yes`
- `PLS_YOLO_MODE=true` can override yolo mode from the environment
- high-risk commands still ask for confirmation
- `--no-exec` forces suggestion-only behavior even in a TTY
- in non-interactive mode and `--json` mode, `pls` stays suggestion-only

## Doctor

Run a quick local sanity check:

```bash
pls doctor
```

That checks things like:
- resolved global config path
- whether a local `pls.json` override is active
- whether yolo mode is enabled and where it came from
- current runtime OS/shell/cwd
- the current platform support tier
- whether `pls` is in `PATH`
- provider basics and a lightweight health check

It also opens with a bad joke, because `pls doctor` kind of deserves one.

## Development

```bash
make test
make build
make print-config-path
make release-snapshot
```

Useful packaging/release prep files:
- `RELEASE.md`
- `packaging/README.md`
- `packaging/homebrew/pls.rb.tmpl`
- `packaging/scoop/pls.json.tmpl`

If your default Go is too old:

```bash
GO=~/.local/bin/go1.26 make test
GO=~/.local/bin/go1.26 make build
```

## Setup and config commands

`pls setup` is a friendly alias for `pls config init`.

`pls config init` walks through global provider setup for all supported providers, writes the global config file, and can enable yolo mode.

`pls config local init` writes `./pls.json` for the current project and focuses on local overrides like provider/model/host/yolo mode. It never stores API keys locally.

`pls config show` prints the effective config state, including the active global path, any local override, provider/model/host, and yolo mode.

`pls config path` prints the resolved global config path.

Current scope:
- global config wizard for first-run setup across all supported providers
- project-local wizard for `pls.json` overrides
- provider-aware host/model prompts
- Ollama model discovery when the target host is reachable
- OpenAI API key prompt for global config only
- environment-variable credential guidance for other hosted providers
- optional yolo mode toggle in both wizards
- local wizard avoids storing API keys

## Notes

- Everything after `pls` is treated as the request unless parsed as a known flag.
- In a real TTY, `pls` can prompt with `Run it? [y/N]` and execute the suggested command through your shell.
- `--yes` auto-runs only low/medium-risk suggestions; high-risk ones still require confirmation.
- `--no-exec` forces suggestion-only behavior.
- Safety policy can escalate risky commands for manual review.
- Style normalization prefers boring direct commands over parsing `ls` output for common listing tasks.
- More advanced prompts can return concise shell loops for batch file operations when that is the clearest single command.
- Bulk rename and move commands are treated as high-risk suggestions and should be reviewed before execution.
- The previous Node prototype is preserved under `legacy/node-prototype/`.
