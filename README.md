# pls

`pls` is a natural-language shell command suggester written in Go.

Current MVP goals:
- command generation only
- no quotes required around the prompt
- Linux-first, shell-aware
- provider support early: Ollama and OpenAI
- explain command + risk level
- JSON output for future integrations

## Examples

```bash
pls doctor
pls show me all dotfiles in this directory
pls find files bigger than 500mb
pls check if jellyfin is running
pls prefix all the mp3s with their lengths in seconds
pls prefix all jpgs with vacation-
pls replace spaces in all filenames here with underscores
pls move all srt files into a subtitles folder
pls --provider openai --model gpt-4.1-mini why is port 3000 busy
pls --provider ollama --model qwen2.5-coder:7b-instruct-q4_K_M show hidden files here
pls --json list the 10 biggest files under the current directory
```

## Config file

Default path:

```bash
~/.config/pls/config.json
```

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
flags > environment > config file > built-in defaults
```

Example config:

```json
{
  "provider": "ollama",
  "model": "qwen2.5-coder:7b-instruct-q4_K_M",
  "host": "http://192.168.2.166:11434"
}
```

A copyable example also lives at `examples/config.example.json`.

## Provider configuration

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

## Run

```bash
pls show hidden files here
```

In an interactive terminal, `pls` now prompts before execution:

```text
Run it? [y/N]
```

If you answer `y`, it runs the suggested command through your shell. In non-interactive mode and `--json` mode, it stays suggestion-only.

## Doctor

Run a quick local sanity check:

```bash
pls doctor
```

That checks things like:
- resolved config path
- whether the config file exists
- current runtime OS/shell/cwd
- whether `pls` is in `PATH`
- provider basics and a lightweight health check

It also opens with a bad joke, because `pls doctor` kind of deserves one.

## Development

```bash
make test
make build
make print-config-path
```

## Notes

- Everything after `pls` is treated as the request unless parsed as a known flag.
- In a real TTY, `pls` can prompt with `Run it? [y/N]` and execute the suggested command through your shell.
- Safety policy can escalate risky commands for manual review.
- Style normalization prefers boring direct commands over parsing `ls` output for common listing tasks.
- More advanced prompts can return concise shell loops for batch file operations when that is the clearest single command.
- Bulk rename and move commands are treated as high-risk suggestions and should be reviewed before execution.
- The previous Node prototype is preserved under `legacy/node-prototype/`.
