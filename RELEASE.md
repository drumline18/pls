# Release & Open-Source Prep

This repo is being prepared for public release, but should **not** be published yet.

## Chosen public settings

- public repo: `github.com/drumline18/pls`
- module path: `github.com/drumline18/pls`
- license: MIT
- recommended first tag: `v0.1.0`

## Already prepared

- `go.mod` uses `github.com/drumline18/pls`
- MIT `LICENSE` added
- CI workflow for tests on Go 1.26.1
- GoReleaser config for cross-platform archives
- local GoReleaser snapshot dry-run succeeded
- contributor / security / issue / PR templates
- README notes for install and release expectations
- Homebrew / Scoop packaging templates under `packaging/`

## Remaining steps before publishing

### 1) Create / push the public GitHub repo

Target:

```text
drumline18/pls
```

Once that repo exists and this code is pushed there, the public module path and GitHub URLs in the repo will line up.

### 2) Tag the first release

Recommended starting point:

```text
v0.1.0
```

### 3) Create the first public release artifacts

Use GoReleaser to produce draft release artifacts after the repo is pushed and tagged.

### 4) Create package-manager repos when ready

Planned names:
- Homebrew tap: `drumline18/homebrew-tap`
- Scoop bucket: `drumline18/scoop-bucket`

## Suggested first public-release sequence

1. Create the GitHub repo `drumline18/pls`
2. Push the repo and verify CI
3. Create tag `v0.1.0`
4. Run GoReleaser to create draft release artifacts
5. Review release notes, checksums, and archives
6. Publish the draft release
7. Create `drumline18/homebrew-tap` and `drumline18/scoop-bucket`
8. Fill in the templates from `packaging/`
9. Publish package-manager entries

## Public install target

Once the repo is live, the intended Go install command is:

```bash
go install github.com/drumline18/pls/cmd/pls@latest
```

## Local dry-run commands

Run tests:

```bash
make test
```

Build locally:

```bash
make build
```

Create a local snapshot build:

```bash
make release-snapshot
```

or directly:

```bash
goreleaser release --snapshot --clean --config .goreleaser.yaml --skip=publish,announce,sign
```

## Packaging templates

See:
- `packaging/README.md`
- `packaging/homebrew/pls.rb.tmpl`
- `packaging/scoop/pls.json.tmpl`
