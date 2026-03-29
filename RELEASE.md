# Release & Open-Source Prep

This repo is being prepared for public release, but should **not** be published yet.

## Already prepared

- CI workflow for tests on Go 1.26.1
- GoReleaser config for cross-platform archives
- Contributor / security / issue / PR templates
- README notes for install and release expectations
- System and local toolchain aligned on Go 1.26.1 during development

## Still needs final decisions before publishing

### 1) Pick the public repo path

`go.mod` still says:

```go
module pls
```

Before public release, change it to the final module path, for example:

```go
module github.com/<owner>/pls
```

This is required for `go install github.com/<owner>/pls/cmd/pls@latest` style installs.

### 2) Choose a license

Do **not** publish without an explicit license file.

Typical options:
- MIT — simple and permissive
- Apache-2.0 — permissive with explicit patent grant
- GPL-3.0 — copyleft

### 3) Decide initial versioning

Suggested starting point:
- `v0.1.0` if still MVP / fast-moving
- `v1.0.0` only if CLI behavior is intentionally stable

### 4) Decide distribution channels

Recommended first wave:
- GitHub Releases archives via GoReleaser
- `go install` once module path is public

Recommended second wave:
- Homebrew tap
- Scoop for Windows
- maybe `nfpm` packages (`.deb`, `.rpm`) later if demand exists

## Suggested first public-release sequence

1. Create the public repo
2. Set the final module path in `go.mod`
3. Add the chosen `LICENSE`
4. Push the repo and verify CI
5. Create a first tag such as `v0.1.0`
6. Use GoReleaser to create draft release artifacts
7. Review release notes and binaries
8. Publish the draft release
9. Add package-manager channels after the first release works cleanly

## Local dry-run commands

Run tests:

```bash
GO=~/.local/bin/go1.26 make test
```

Build locally:

```bash
GO=~/.local/bin/go1.26 make build
```

If GoReleaser is installed, create a local snapshot build:

```bash
goreleaser release --snapshot --clean --config .goreleaser.yaml --skip=publish,announce,sign
```

## Homebrew later

Once the public repo exists, add a Homebrew tap or formula repo and wire a `brews:` section into `.goreleaser.yaml`.

## Scoop later

Once the first Windows release exists, add a Scoop manifest repo and point it at the GitHub Release archive/checksum.
