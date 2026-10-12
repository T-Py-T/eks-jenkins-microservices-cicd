import subprocess
import unittest
from pathlib import Path


REPOSITORY_ROOT = Path(__file__).resolve().parents[1]
SMOKE_SCRIPT = REPOSITORY_ROOT / "tests" / "node_service_smoke.js"


class NodeServiceSmokeTests(unittest.TestCase):
    def test_requires_four_arguments(self) -> None:
        result = subprocess.run(
            ["node", str(SMOKE_SCRIPT)],
            capture_output=True,
            text=True,
            check=False,
        )
        self.assertEqual(result.returncode, 2)
        self.assertIn("usage:", result.stderr)


if __name__ == "__main__":
    unittest.main()
