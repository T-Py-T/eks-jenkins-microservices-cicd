#!/usr/bin/env python3

import argparse
import json
import re
import subprocess
import sys
import tempfile
from pathlib import Path


IMAGE_LINE = re.compile(
    r"^(?P<prefix>\s*image:\s*)(?P<reference>\S+)(?P<suffix>\s*(?:#.*)?)$"
)
DIGEST = re.compile(r"^sha256:[0-9a-f]{64}$")


def parse_arguments() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Resolve tagged container images in a Kubernetes manifest to digests."
    )
    parser.add_argument("--input", type=Path, required=True)
    parser.add_argument("--output", type=Path, required=True)
    return parser.parse_args()


def split_digest(reference: str) -> tuple[str, str | None]:
    if "@" not in reference:
        return reference, None
    repository, digest = reference.rsplit("@", 1)
    return repository, digest


def split_tag(reference: str) -> tuple[str, str | None]:
    last_slash = reference.rfind("/")
    last_colon = reference.rfind(":")
    if last_colon <= last_slash:
        return reference, None
    return reference[:last_colon], reference[last_colon + 1 :]


def inspect_digest(reference: str) -> str:
    result = subprocess.run(
        ["skopeo", "inspect", f"docker://{reference}"],
        capture_output=True,
        check=False,
        text=True,
    )
    if result.returncode != 0:
        detail = result.stderr.strip() or "registry inspection failed"
        raise ValueError(f"cannot resolve {reference}: {detail}")

    try:
        digest = json.loads(result.stdout)["Digest"]
    except (json.JSONDecodeError, KeyError, TypeError) as error:
        raise ValueError(f"registry returned no digest for {reference}") from error
    if not isinstance(digest, str) or not DIGEST.fullmatch(digest):
        raise ValueError(f"registry returned an invalid digest for {reference}")
    return digest


def pin_reference(reference: str) -> str:
    repository, digest = split_digest(reference)
    if digest is not None:
        if not DIGEST.fullmatch(digest):
            raise ValueError(f"invalid image digest: {reference}")
        return reference

    repository, tag = split_tag(repository)
    if tag is None:
        raise ValueError(f"image requires an explicit version tag: {reference}")
    if tag == "latest":
        raise ValueError(f"mutable image tag is not allowed: {reference}")
    if "BUILD_NUMBER" in tag:
        raise ValueError(f"replace BUILD_NUMBER before deployment: {reference}")

    return f"{repository}@{inspect_digest(reference)}"


def pin_manifest(input_path: Path, output_path: Path) -> int:
    resolved_count = 0
    output_lines: list[str] = []
    for line in input_path.read_text().splitlines(keepends=True):
        match = IMAGE_LINE.match(line.rstrip("\n"))
        if match is None:
            output_lines.append(line)
            continue

        reference = match.group("reference")
        pinned_reference = pin_reference(reference)
        if pinned_reference != reference:
            resolved_count += 1
        newline = "\n" if line.endswith("\n") else ""
        output_lines.append(
            f"{match.group('prefix')}{pinned_reference}{match.group('suffix')}{newline}"
        )

    if resolved_count == 0:
        raise ValueError("manifest contains no version-tagged images to resolve")

    output_path.parent.mkdir(parents=True, exist_ok=True)
    with tempfile.NamedTemporaryFile(
        "w", dir=output_path.parent, delete=False
    ) as temporary_file:
        temporary_file.writelines(output_lines)
        temporary_path = Path(temporary_file.name)
    temporary_path.replace(output_path)
    return resolved_count


def main() -> int:
    arguments = parse_arguments()
    try:
        resolved_count = pin_manifest(arguments.input, arguments.output)
    except (OSError, ValueError) as error:
        print(f"error: {error}", file=sys.stderr)
        return 1
    print(f"resolved {resolved_count} image references in {arguments.output}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
