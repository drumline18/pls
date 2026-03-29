# Publishing `pls`

This file is the exact handoff for the first live public release of `pls`.

Current chosen settings:
- GitHub repo: `drumline18/pls`
- module path: `github.com/drumline18/pls`
- license: MIT
- recommended first version: `v0.1.0`

## Current state

- GitHub auth is active for `drumline18`
- `drumline18/pls` already exists and is public
- default branch is `main`

## 1) Sync the latest code

The public repo is maintained from the local `pls/` subtree inside the workspace. Push updates before tagging a release.

## 2) Create the first tag

```bash
git tag v0.1.0
git push origin v0.1.0
```

## 4) Create release artifacts

### Option A: local GoReleaser publish

```bash
GITHUB_TOKEN="$(gh auth token)" goreleaser release --clean --config .goreleaser.yaml
```

### Option B: dry-run locally first

```bash
make release-snapshot
```

## 5) Verify release outputs

You should see release archives for:
- Linux amd64/arm64
- macOS amd64/arm64
- Windows amd64/arm64
- `checksums.txt`

## 6) Generate Homebrew / Scoop package files

After the tagged release artifacts exist and checksums are known:

```bash
python3 scripts/render_packaging.py \
  --version v0.1.0 \
  --checksums dist/checksums.txt
```

That writes generated files to:
- `packaging/generated/homebrew/pls.rb`
- `packaging/generated/scoop/pls.json`

## 7) Publish package-manager repos

### Homebrew tap

Create:
- `drumline18/homebrew-tap`

Then copy:
- `packaging/generated/homebrew/pls.rb` → `Formula/pls.rb`

### Scoop bucket

Create:
- `drumline18/scoop-bucket`

Then copy:
- `packaging/generated/scoop/pls.json` → `pls.json`

## Notes

- `go install github.com/drumline18/pls/cmd/pls@latest` already works because the repo is public.
- GitHub Releases / packaged binaries are the next release task.
- If you want GitHub Releases to happen in CI later, we can add a tag-triggered publish workflow next.
