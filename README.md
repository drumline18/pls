# pls

`pls` is a natural-language shell command suggester.

Current MVP goals:
- command generation only
- no quotes required around the prompt
- Linux-first, shell-aware
- provider support early: Ollama and OpenAI
- explain command + risk level
- JSON output for future integrations

## Examples

```bash
pls show me all dotfiles in this directory
pls find files bigger than 500mb
pls --provider openai --model gpt-4.1-mini why is port 3000 busy
pls --provider ollama --model qwen2.5-coder:7b show hidden files here
pls --json list the 10 biggest files under the current directory
```

## Provider configuration

### Ollama

Defaults:
- provider: `ollama` if no OpenAI API key is present
- host: `http://127.0.0.1:11434`
- model: `qwen2.5-coder:7b-instruct-q4_K_M`

Environment variables:

```bash
export OLLAMA_HOST=http://127.0.0.1:11434
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

## Run locally

```bash
cd pls
chmod +x ./bin/pls.js
node ./bin/pls.js show hidden files here
```

## Install locally into PATH

From this folder:

```bash
npm link
```

Then:

```bash
pls show hidden files here
```

## Notes

- Everything after `pls` is treated as the request unless parsed as a known flag.
- This version does not execute commands yet.
- Safety policy can escalate risky commands for manual review.
