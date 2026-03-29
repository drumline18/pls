# Release Notes

`pls` is now public and has an initial release:

- repo: `github.com/drumline18/pls`
- current release: `v0.1.0`
- license: MIT

## Live distribution channels

- GitHub repo: <https://github.com/drumline18/pls>
- GitHub Releases: <https://github.com/drumline18/pls/releases>
- Homebrew tap: <https://github.com/drumline18/homebrew-tap>
- Scoop bucket: <https://github.com/drumline18/scoop-bucket>

## Install targets

### Go install

```bash
go install github.com/drumline18/pls/cmd/pls@latest
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

## What exists in-repo for future releases

- GoReleaser config for cross-platform archives
- CI workflow for Go 1.26.1
- packaging templates under `packaging/`
- package rendering helper: `scripts/render_packaging.py`

## Future release flow

For the next release after `v0.1.0`, the rough flow is:

1. push the latest `pls` subtree to `drumline18/pls`
2. create and push a new tag
3. run GoReleaser against the public repo
4. publish the GitHub release
5. regenerate Homebrew/Scoop package files from the new checksums
6. update the tap and bucket repos

## Local dry-run commands

```bash
make test
make build
make release-snapshot
```

or directly:

```bash
goreleaser release --snapshot --clean --config .goreleaser.yaml --skip=publish,announce,sign
```
