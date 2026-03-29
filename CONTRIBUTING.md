# Contributing to `pls`

Thanks for helping.

## Development setup

Requirements:
- Go 1.25+
- preferably Go 1.26.1 during development on this project

Run tests:

```bash
make test
```

If your default Go is older:

```bash
GO=~/.local/bin/go1.26 make test
```

Build locally:

```bash
make build
```

## Before opening a PR

- keep changes focused
- add or update tests when behavior changes
- update README/docs when user-facing behavior changes
- run `make test`

## PR style

Please include:
- what changed
- why it changed
- any UX or behavior impact
- any follow-up work intentionally left out

## Scope notes

`pls` is intentionally conservative about execution and risk handling. Changes that affect:
- confirmation behavior
- risk classification
- config precedence
- provider selection
- shell normalization

should come with tests and clear explanation.
