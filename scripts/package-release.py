#!/usr/bin/env python3
"""Archive a staged release without macOS sidecar files or extended metadata."""
import hashlib
from pathlib import Path
import sys
import tarfile


def package_release(source: Path, output: Path):
    def public_files(member):
        if any(part.startswith("._") or part == ".DS_Store" for part in Path(member.name).parts):
            return None
        return member

    # Python's tar writer does not generate AppleDouble files from macOS xattrs.
    with output.open("xb") as stream:
        with tarfile.open(fileobj=stream, mode="w:gz") as archive:
            archive.add(source, arcname=source.name, filter=public_files)
    checksum = hashlib.sha256()
    with output.open("rb") as stream:
        for chunk in iter(lambda: stream.read(1024 * 1024), b""):
            checksum.update(chunk)
    digest = checksum.hexdigest()
    output.with_name(output.name + ".sha256").write_text(f"{digest}  {output.name}\n", encoding="utf-8")


if __name__ == "__main__":
    package_release(Path(sys.argv[1]), Path(sys.argv[2]))
