# Packaging Scaffolding

These files are intentionally **templates/scaffolding only**.

They are here so `pls` is ready for package-manager distribution once the public repositories exist, but they should not be published automatically yet.

## Planned public locations

- main repo: `github.com/drumline18/pls`
- Homebrew tap: `drumline18/homebrew-tap`
- Scoop bucket: `drumline18/scoop-bucket`

## Files

- `homebrew/pls.rb.tmpl` — Homebrew formula template
- `scoop/pls.json.tmpl` — Scoop manifest template

## How to use after the first public release

### Homebrew

1. Create the tap repo, for example `drumline18/homebrew-tap`
2. Copy `homebrew/pls.rb.tmpl` to `Formula/pls.rb`
3. Replace:
   - `__VERSION__`
   - `__DARWIN_AMD64_SHA256__`
   - `__DARWIN_ARM64_SHA256__`
4. Commit to the tap repo

### Scoop

1. Create the bucket repo, for example `drumline18/scoop-bucket`
2. Copy `scoop/pls.json.tmpl` to `pls.json`
3. Replace:
   - `__VERSION__`
   - `__WINDOWS_AMD64_SHA256__`
4. Commit to the bucket repo

## Rendering generated package files

After a real tagged release exists and `dist/checksums.txt` is available, render the files with:

```bash
python3 scripts/render_packaging.py --version v0.1.0 --checksums dist/checksums.txt
```

That writes:
- `packaging/generated/homebrew/pls.rb`
- `packaging/generated/scoop/pls.json`

## Notes

The current GoReleaser setup already produces the release archives and checksums these templates expect.
