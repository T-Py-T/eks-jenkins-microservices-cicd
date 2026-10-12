import sys
import tempfile
import unittest
from pathlib import Path


REPOSITORY_ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(REPOSITORY_ROOT / "scripts"))

import pin_manifest_images as pin  # noqa: E402


class PinManifestHelperTests(unittest.TestCase):
    def test_split_tag_preserves_registry_port(self) -> None:
        repository, tag = pin.split_tag(
            "registry.example.com:5000/team/emailservice:1.42"
        )
        self.assertEqual(repository, "registry.example.com:5000/team/emailservice")
        self.assertEqual(tag, "1.42")

    def test_pin_reference_keeps_valid_digest(self) -> None:
        digest = "sha256:" + ("c" * 64)
        reference = f"redis:8.10.1-alpine@{digest}"
        self.assertEqual(pin.pin_reference(reference), reference)

    def test_pin_reference_rejects_invalid_digest(self) -> None:
        with self.assertRaises(ValueError):
            pin.pin_reference("redis:8.10.1-alpine@sha256:not-a-digest")

    def test_pin_reference_requires_explicit_tag(self) -> None:
        with self.assertRaises(ValueError):
            pin.pin_reference("registry.example.com/team/emailservice")

    def test_pin_manifest_requires_at_least_one_tagged_image(self) -> None:
        digest = "sha256:" + ("d" * 64)
        manifest = f"containers:\n  - name: redis\n    image: redis:8.10.1-alpine@{digest}\n"
        with tempfile.TemporaryDirectory() as temporary_directory:
            workspace = Path(temporary_directory)
            input_path = workspace / "deployment.yml"
            output_path = workspace / "deployment.pinned.yml"
            input_path.write_text(manifest)
            with self.assertRaises(ValueError):
                pin.pin_manifest(input_path, output_path)


if __name__ == "__main__":
    unittest.main()
