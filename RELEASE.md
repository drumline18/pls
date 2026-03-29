# Release & Open-Source Prep

This repo is now public, but packaged releases and distribution channels are still being prepared.

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

## Remaining steps before packaged release

### 1) Tag the first release

Recommended starting point:

```text
v0.1.0
```

### 2) Create the first public release artifacts

Use GoReleaser to produce draft or published release artifacts after the repo is tagged.

### 3) Create package-manager repos when ready

Planned names:
- Homebrew tap: `drumline18/homebrew-tap`
- Scoop bucket: `drumline18/scoop-bucket`

## Suggested first packaged-release sequence

1. Verify CI on `drumline18/pls`
2. Create tag `v0.1.0`
3. Run GoReleaser to create draft release artifacts
4. Review release notes, checksums, and archives
5. Publish the release
6. Create `drumline18/homebrew-tap` and `drumline18/scoop-bucket`
7. Fill in the templates from `packaging/`
8. Publish package-manager entries

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
