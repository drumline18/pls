#!/usr/bin/env python3
import argparse
import pathlib
import re
import sys


PLACEHOLDERS = {
    "darwin_amd64": "__DARWIN_AMD64_SHA256__",
    "darwin_arm64": "__DARWIN_ARM64_SHA256__",
    "linux_amd64": "__LINUX_AMD64_SHA256__",
    "linux_arm64": "__LINUX_ARM64_SHA256__",
    "windows_amd64": "__WINDOWS_AMD64_SHA256__",
    "windows_arm64": "__WINDOWS_ARM64_SHA256__",
}


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Render Homebrew and Scoop packaging files from checksums.")
    parser.add_argument("--version", required=True, help="Release version, with or without leading v (example: v0.1.0)")
    parser.add_argument("--checksums", required=True, help="Path to GoReleaser checksums.txt")
    parser.add_argument(
        "--output-dir",
        default="packaging/generated",
        help="Output directory for rendered package files (default: packaging/generated)",
    )
    return parser.parse_args()


def normalize_version(version: str) -> str:
    version = version.strip()
    return version[1:] if version.startswith("v") else version


def load_checksums(path: pathlib.Path, version: str) -> dict[str, str]:
    pattern = re.compile(rf"^([a-f0-9]{{64}})\s+pls_{re.escape(version)}_(darwin|linux|windows)_(amd64|arm64)\.(tar\.gz|zip)$")
    values: dict[str, str] = {}
    for raw_line in path.read_text().splitlines():
        line = raw_line.strip()
        if not line:
            continue
        match = pattern.match(line)
        if not match:
            continue
        checksum, os_name, arch, _ext = match.groups()
        values[f"{os_name}_{arch}"] = checksum
    return values


def render_template(template_path: pathlib.Path, replacements: dict[str, str]) -> str:
    content = template_path.read_text()
    for old, new in replacements.items():
        content = content.replace(old, new)
    return content


def write_output(path: pathlib.Path, content: str) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(content)


def main() -> int:
    args = parse_args()
    repo_root = pathlib.Path(__file__).resolve().parent.parent
    version = normalize_version(args.version)
    checksums_path = pathlib.Path(args.checksums)
    if not checksums_path.is_absolute():
        checksums_path = repo_root / checksums_path
    output_dir = pathlib.Path(args.output_dir)
    if not output_dir.is_absolute():
        output_dir = repo_root / output_dir

    if not checksums_path.exists():
        print(f"checksums file not found: {checksums_path}", file=sys.stderr)
        return 1

    checksums = load_checksums(checksums_path, version)
    missing = [key for key in PLACEHOLDERS if key not in checksums]
    if missing:
        print("missing checksum entries for: " + ", ".join(missing), file=sys.stderr)
        return 1

    replacements = {"__VERSION__": version}
    for key, placeholder in PLACEHOLDERS.items():
        replacements[placeholder] = checksums[key]

    homebrew_template = repo_root / "packaging" / "homebrew" / "pls.rb.tmpl"
    scoop_template = repo_root / "packaging" / "scoop" / "pls.json.tmpl"

    write_output(output_dir / "homebrew" / "pls.rb", render_template(homebrew_template, replacements))
    write_output(output_dir / "scoop" / "pls.json", render_template(scoop_template, replacements))

    print(f"wrote {output_dir / 'homebrew' / 'pls.rb'}")
    print(f"wrote {output_dir / 'scoop' / 'pls.json'}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
