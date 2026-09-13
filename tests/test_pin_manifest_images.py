import os
import subprocess
import sys
import tempfile
import textwrap
import unittest
from pathlib import Path


REPOSITORY_ROOT = Path(__file__).resolve().parents[1]
PIN_MANIFEST = REPOSITORY_ROOT / "scripts" / "pin_manifest_images.py"
RESOLVED_DIGEST = "sha256:" + ("a" * 64)


class PinManifestImagesTests(unittest.TestCase):
    def setUp(self) -> None:
        self.temporary_directory = tempfile.TemporaryDirectory()
        self.workspace = Path(self.temporary_directory.name)
        self.bin_directory = self.workspace / "bin"
        self.bin_directory.mkdir()
        skopeo = self.bin_directory / "skopeo"
        skopeo.write_text(
            "#!/bin/sh\n"
            f"printf '%s\\n' '{{\"Digest\":\"{RESOLVED_DIGEST}\"}}'\n"
        )
        skopeo.chmod(0o755)

    def tearDown(self) -> None:
        self.temporary_directory.cleanup()

    def run_pin(
        self, manifest: str
    ) -> tuple[subprocess.CompletedProcess[str], Path]:
        input_path = self.workspace / "deployment.yml"
        output_path = self.workspace / "deployment.pinned.yml"
        input_path.write_text(textwrap.dedent(manifest))
        environment = os.environ.copy()
        environment["PATH"] = f"{self.bin_directory}:{environment['PATH']}"

        result = subprocess.run(
            [
                sys.executable,
                str(PIN_MANIFEST),
                "--input",
                str(input_path),
                "--output",
                str(output_path),
            ],
            capture_output=True,
            check=False,
            env=environment,
            text=True,
        )
        return result, output_path

    def test_resolves_versioned_tags_and_preserves_existing_digests(self) -> None:
        existing_digest = "sha256:" + ("b" * 64)
        result, output_path = self.run_pin(
            f"""
            containers:
              - name: emailservice
                image: registry.example.com:5000/team/emailservice:1.42
              - name: redis
                image: redis:8.10.1-alpine@{existing_digest}
            """
        )

        self.assertEqual(result.returncode, 0, result.stderr)
        output = output_path.read_text()
        self.assertIn(
            f"registry.example.com:5000/team/emailservice@{RESOLVED_DIGEST}",
            output,
        )
        self.assertIn(f"redis:8.10.1-alpine@{existing_digest}", output)

    def test_rejects_latest_tag(self) -> None:
        result, _ = self.run_pin(
            "containers:\n  - name: paymentservice\n"
            "    image: example/paymentservice:latest\n"
        )

        self.assertNotEqual(result.returncode, 0)
        self.assertIn("mutable image tag", result.stderr)

    def test_rejects_unresolved_build_placeholder(self) -> None:
        result, _ = self.run_pin(
            "containers:\n  - name: paymentservice\n"
            "    image: example/paymentservice:1.BUILD_NUMBER\n"
        )

        self.assertNotEqual(result.returncode, 0)
        self.assertIn("replace BUILD_NUMBER", result.stderr)


if __name__ == "__main__":
    unittest.main()
