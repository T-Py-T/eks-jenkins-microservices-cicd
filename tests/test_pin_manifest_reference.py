import sys
import unittest
from pathlib import Path


REPOSITORY_ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(REPOSITORY_ROOT / "scripts"))

import pin_manifest_images as pin  # noqa: E402


class PinManifestReferenceTests(unittest.TestCase):
    def test_split_digest_separates_repository_and_digest(self) -> None:
        repository, digest = pin.split_digest(
            "registry.example.com/app/service:1.2.3@sha256:" + ("a" * 64)
        )
        self.assertEqual(repository, "registry.example.com/app/service:1.2.3")
        self.assertEqual(digest, "sha256:" + ("a" * 64))

    def test_split_tag_respects_registry_port(self) -> None:
        repository, tag = pin.split_tag("registry.example.com:5000/team/service:9.8.7")
        self.assertEqual(repository, "registry.example.com:5000/team/service")
        self.assertEqual(tag, "9.8.7")

    def test_pin_reference_rejects_latest_tag(self) -> None:
        with self.assertRaisesRegex(ValueError, "mutable image tag"):
            pin.pin_reference("paymentservice:latest")

    def test_pin_reference_rejects_build_number_placeholder(self) -> None:
        with self.assertRaisesRegex(ValueError, "BUILD_NUMBER"):
            pin.pin_reference("paymentservice:1.0.BUILD_NUMBER")

    def test_pin_reference_preserves_valid_digest(self) -> None:
        digest = "sha256:" + ("c" * 64)
        reference = f"redis:8.10.1-alpine@{digest}"
        self.assertEqual(pin.pin_reference(reference), reference)


if __name__ == "__main__":
    unittest.main()
