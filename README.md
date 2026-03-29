# pls

`pls` is a natural-language shell command suggester written in Go.

You type what you want in plain English. `pls` turns that into a shell command, explains what it does, rates the risk, and can ask for confirmation before running it.

It is built for the kind of things people actually do in a terminal:
- inspection commands
- file operations
- service checks
- shell-friendly one-liners
- quick "what's the command for this?" moments

## Demo

![pls demo](assets/pls-demo.gif)

## Why use it?

Because a lot of shell work is really this:
- "show hidden files here"
- "check if jellyfin is running"
- "move all srt files into a subtitles folder"
- "find files bigger than 500mb"

You know what you want. You just do not always want to stop and reconstruct the exact command syntax from memory.

`pls` aims to be:
- direct
- readable
- reasonably safe
- boring in a good way

## Install

### Fastest install from GitHub

If you already have Go 1.25+:

```bash
go install github.com/drumline18/pls/cmd/pls@latest
```

That installs the `pls` binary into your Go bin directory.

If that directory is not already on your `PATH`, add one of these depending on your setup:

```bash
export PATH="$HOME/go/bin:$PATH"
```

or:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

### Install from a local clone

```bash
git clone https://github.com/drumline18/pls.git
cd pls
make install
```

That installs `pls` to:

```bash
~/.local/bin/pls
```

### Build manually

```bash
git clone https://github.com/drumline18/pls.git
cd pls
go build -o bin/pls ./cmd/pls
```

### Homebrew

```bash
brew tap drumline18/tap
brew install pls
```

### Scoop

```powershell
scoop bucket add drumline18 https://github.com/drumline18/scoop-bucket
scoop install pls
```

### Release archives

Prebuilt binaries are also available from GitHub Releases:

- <https://github.com/drumline18/pls/releases>

### Distribution status

Current state:
- GitHub repo is live
- `go install github.com/drumline18/pls/cmd/pls@latest` works
- GitHub Releases are live
- Homebrew tap is live
- Scoop bucket is live

Install links:
- Releases: <https://github.com/drumline18/pls/releases>
- Homebrew tap: <https://github.com/drumline18/homebrew-tap>
- Scoop bucket: <https://github.com/drumline18/scoop-bucket>

## Quick start

If you want a guided first-run setup:

```bash
pls setup
```

or:

```bash
pls config init
```

That lets you choose a provider, model, host, and optional execution behavior.

Then try something simple:

```bash
pls show hidden files here
pls check if jellyfin is running
pls find files bigger than 500mb under the current directory
```

## What it feels like

Examples:

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
pls --provider ollama --model qwen3.5:4b show hidden files here
pls --json list the 10 biggest files under the current directory
pls -- show me files named --json
```

## How execution works

In an interactive terminal, `pls` can suggest a command and then ask before running it:

```text
Run it? [y/N]
```

Useful flags:

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

## Platform support

Current support status:
- Linux: **primary**
- macOS: **beta**
- Windows PowerShell: **beta**
- Windows cmd.exe: **limited**

Linux currently has the strongest normalization and post-processing rules. macOS and PowerShell are supported, but have less hand-tuned coverage right now.

## Configuration

Global config path:

```bash
~/.config/pls/config.json
```

Built-in config commands:

```bash
pls config init
pls config local init
pls config show
pls config path
```

`pls setup` is just a friendly alias for `pls config init`.

`pls setup`, `pls config init`, and `pls config local init` are strict built-in commands. Longer phrases like `pls config init my project` still go through the normal natural-language path.

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

Config precedence:

```text
flags > environment > local pls.json > global config > built-in defaults
```

Example global config:

```json
{
  "provider": "ollama",
  "model": "qwen3.5:4b",
  "host": "http://127.0.0.1:11434",
  "yoloMode": false
}
```

Example local override (`pls.json` inside a project):

```json
{
  "yoloMode": true
}
```

`pls` looks for `pls.json` in the current directory, then walks upward until it finds one.

A copyable example also lives at `examples/config.example.json`.

## Providers

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

The setup wizard covers all of them.

Behavior:
- the global wizard can set provider / model / host for any supported provider
- OpenAI can optionally store an API key in global config
- non-OpenAI hosted providers use their usual environment variable for credentials
- the local wizard can override provider / model / host / yolo mode, but never stores API keys locally

### Ollama

Typical environment variables:

```bash
export OLLAMA_HOST=http://127.0.0.1:11434
export PLS_OLLAMA_HOST=http://127.0.0.1:11434
export PLS_PROVIDER=ollama
export PLS_MODEL=qwen3.5:4b
```

### OpenAI

```bash
export OPENAI_API_KEY=your_key_here
export PLS_PROVIDER=openai
export PLS_MODEL=gpt-4.1-mini
```

### Other providers

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

## Doctor

Run a quick sanity check:

```bash
pls doctor
```

That checks things like:
- resolved global config path
- whether a local `pls.json` override is active
- whether yolo mode is enabled and where it came from
- current runtime OS / shell / cwd
- current platform support tier
- whether `pls` is in `PATH`
- provider basics and a lightweight health check

It also opens with a bad joke, because `pls doctor` kind of deserves one.

## Development

Toolchain target:
- Go **1.25+**
- Go **1.26.1** is the current development baseline here

Common commands:

```bash
make test
make build
make print-config-path
make release-snapshot
```

If your default Go is too old:

```bash
GO=~/.local/bin/go1.26 make test
GO=~/.local/bin/go1.26 make build
```

Useful project files:
- `RELEASE.md`
- `PUBLISHING.md`
- `packaging/README.md`
- `packaging/homebrew/pls.rb.tmpl`
- `packaging/scoop/pls.json.tmpl`
- `scripts/render_packaging.py`
- `demo/readme.tape`
- `scripts/render_readme_demo.sh`

## Notes

- Everything after `pls` is treated as the request unless parsed as a known leading flag.
- `pls` prefers direct, readable commands over text-parsing hacks for common tasks.
- Style normalization intentionally rewrites some generated commands into safer or more boring equivalents.
- More advanced prompts can still return short shell loops for batch file operations when that is the clearest answer.
- Bulk rename and move commands are treated as high-risk suggestions and should be reviewed before execution.
