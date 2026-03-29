# Publishing `pls`

This file is the operational handoff for future releases of `pls`.

Current state:
- public repo: `drumline18/pls`
- default branch: `main`
- first public release: `v0.1.0`
- Homebrew tap: `drumline18/homebrew-tap`
- Scoop bucket: `drumline18/scoop-bucket`

## Release sequence

### 1) Sync the latest code

The public repo is maintained from the local `pls/` subtree inside the workspace. Push updates before tagging a release.

### 2) Create the next tag

Example:

```bash
git tag v0.1.1
git push origin v0.1.1
```

### 3) Publish release artifacts

```bash
GITHUB_TOKEN="$(gh auth token)" goreleaser release --clean --config .goreleaser.yaml
```

## Verify release outputs

You should see release archives for:
- Linux amd64/arm64
- macOS amd64/arm64
- Windows amd64/arm64
- `checksums.txt`

## Regenerate Homebrew / Scoop package files

After the tagged release artifacts exist and checksums are known:

```bash
python3 scripts/render_packaging.py \
  --version v0.1.1 \
  --checksums dist/checksums.txt
```

That writes:
- `packaging/generated/homebrew/pls.rb`
- `packaging/generated/scoop/pls.json`

## Publish package-manager updates

### Homebrew tap

Repo:
- `drumline18/homebrew-tap`

Copy:
- `packaging/generated/homebrew/pls.rb` → `Formula/pls.rb`

### Scoop bucket

Repo:
- `drumline18/scoop-bucket`

Copy:
- `packaging/generated/scoop/pls.json` → `pls.json`

## Notes

- `go install github.com/drumline18/pls/cmd/pls@latest` already works because the repo is public.
- Homebrew and Scoop are live, so future releases should update those repos after every tagged release.
