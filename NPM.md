# npm Distribution Plan for `pls`

This document covers a practical npm strategy for `pls` as a **distribution channel**, not as a rewrite of the CLI into JavaScript.

## Goal

Offer a familiar cross-platform install command for users who already have Node.js:

```bash
npm install -g @drumline18/pls
```

That command should install the correct prebuilt `pls` binary for the user's platform.

## Why do this?

The current distribution story is technically good:
- `go install`
- GitHub Releases
- Homebrew
- Scoop

But npm can reduce user friction because many people already have:
- `node`
- `npm`
- a habit of using `npm install -g ...`

## Recommended architecture

Treat npm as a **binary installer wrapper**.

That means:
- keep the real CLI implemented in Go
- keep GitHub Releases as the canonical binary source
- let npm download/expose the right binary for the current platform

## Good implementation patterns

### Option A — single installer package

Package:
- `@drumline18/pls`

Behavior:
- postinstall script detects OS/arch
- downloads the correct binary from GitHub Releases
- places it into the package's bin path
- exposes `pls`

Pros:
- simple user-facing package name
- easiest install story

Cons:
- more postinstall logic
- more moving parts in one package

### Option B — platform-specific helper packages + thin meta package

Packages like:
- `@drumline18/pls-linux-x64`
- `@drumline18/pls-linux-arm64`
- `@drumline18/pls-darwin-x64`
- `@drumline18/pls-darwin-arm64`
- `@drumline18/pls-win32-x64`
- `@drumline18/pls-win32-arm64`
- plus a thin `@drumline18/pls` meta package

Pros:
- cleaner package contents
- easier per-platform debugging
- closer to how some established CLIs distribute native binaries via npm

Cons:
- more packages to manage

## My recommendation

Start with **Option A** unless packaging complexity becomes annoying.

Why:
- faster to ship
- easier to explain
- good enough for an early-stage CLI

## Release flow with npm added

Current release flow:
1. push code
2. tag release
3. run GoReleaser
4. publish GitHub Release
5. update Homebrew/Scoop

With npm:
6. publish/update `@drumline18/pls` to point at the new GitHub release binaries

## What the npm package should do

### Expected user experience

```bash
npm install -g @drumline18/pls
pls --help
```

### Internals

At install time:
- detect platform (`process.platform`)
- detect architecture (`process.arch`)
- map to release artifact name
- download the right archive from:
  - `https://github.com/drumline18/pls/releases/download/vX.Y.Z/...`
- extract binary
- expose executable via package `bin`

## Platform mapping

Release artifacts currently exist for:
- `pls_0.1.0_linux_amd64.tar.gz`
- `pls_0.1.0_linux_arm64.tar.gz`
- `pls_0.1.0_darwin_amd64.tar.gz`
- `pls_0.1.0_darwin_arm64.tar.gz`
- `pls_0.1.0_windows_amd64.zip`
- `pls_0.1.0_windows_arm64.zip`

The npm installer should map platform/arch to those exact artifact patterns.

## Safety and support notes

The npm package should:
- fail clearly on unsupported platforms
- print the GitHub Releases URL as fallback
- verify checksums if practical
- avoid pretending the tool is JS-native

## README positioning

If npm is added, I would present install options like this:

### Recommended for your OS
- macOS → Homebrew
- Windows → Scoop
- Linux → release binary or Go install

### Works everywhere
- `npm install -g @drumline18/pls`

### Other options
- `go install`
- manual release download

## Suggested next implementation steps

1. create a small `npm/` distribution package folder
2. add artifact download + extraction logic
3. test on Linux first
4. publish a first `@drumline18/pls` version manually
5. then wire npm publish into the release flow if it proves worthwhile
