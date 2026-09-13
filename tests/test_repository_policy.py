import re
import unittest
from pathlib import Path


REPOSITORY_ROOT = Path(__file__).resolve().parents[1]
SHA256_DIGEST = re.compile(r"@sha256:[0-9a-f]{64}$")
PINNED_ACTION = re.compile(r"\buses:\s+[^\s@]+@[0-9a-f]{40}(?:\s+#.*)?$")
DISALLOWED_GITHUB_EVENTS = re.compile(
    r"^  (?:push|schedule|workflow_dispatch|workflow_call|workflow_run|repository_dispatch):",
    re.MULTILINE,
)


def external_parent_image(from_line: str) -> str | None:
    tokens = from_line.split()
    if not tokens or tokens[0].upper() != "FROM":
        return None

    image_index = 2 if len(tokens) > 1 and tokens[1].startswith("--platform=") else 1
    image = tokens[image_index]
    if image.lower() in {"base", "build", "builder", "publish"}:
        return None
    return image


class RepositoryPolicyTests(unittest.TestCase):
    def test_github_workflows_only_run_for_pull_requests(self) -> None:
        workflow_paths = sorted((REPOSITORY_ROOT / ".github" / "workflows").glob("*.y*ml"))
        self.assertTrue(workflow_paths, "at least one GitHub workflow is required")

        for path in workflow_paths:
            contents = path.read_text()
            self.assertRegex(contents, r"(?m)^on:\s*$")
            self.assertRegex(contents, r"(?m)^  pull_request:\s*$")
            self.assertNotRegex(contents, DISALLOWED_GITHUB_EVENTS)

    def test_github_actions_are_pinned_to_commits(self) -> None:
        workflow_paths = sorted((REPOSITORY_ROOT / ".github" / "workflows").glob("*.y*ml"))
        action_lines = [
            line.strip()
            for path in workflow_paths
            for line in path.read_text().splitlines()
            if "uses:" in line
        ]
        self.assertTrue(action_lines, "at least one GitHub Action is required")
        for line in action_lines:
            self.assertRegex(line, PINNED_ACTION)

    def test_external_docker_parents_are_digest_pinned(self) -> None:
        dockerfiles = sorted(REPOSITORY_ROOT.glob("services/**/Dockerfile*"))
        self.assertEqual(len(dockerfiles), 12)

        for path in dockerfiles:
            for line in path.read_text().splitlines():
                image = external_parent_image(line)
                if image is not None:
                    self.assertRegex(image, SHA256_DIGEST, f"unpinned parent in {path}: {image}")

    def test_container_os_packages_do_not_float(self) -> None:
        dockerfiles = sorted(REPOSITORY_ROOT.glob("services/**/Dockerfile*"))
        for path in dockerfiles:
            contents = path.read_text()
            self.assertNotRegex(
                contents,
                r"\b(?:apk|apt-get)\s+upgrade\b",
                f"blanket OS upgrade in {path}",
            )
            for line in contents.splitlines():
                stripped = line.strip()
                if "apk add --no-cache" not in stripped:
                    continue
                packages = stripped.split("apk add --no-cache", 1)[1].strip().rstrip("\\").split()
                self.assertTrue(packages, f"missing package list in {path}")
                self.assertTrue(
                    all("=" in package for package in packages),
                    f"unpinned Alpine package in {path}: {stripped}",
                )

    def test_manifest_has_reproducible_and_native_dependencies(self) -> None:
        manifest = (REPOSITORY_ROOT / "deploy" / "eks" / "deployment-service.yml").read_text()
        self.assertRegex(manifest, r"image: redis:\d+\.\d+\.\d+-alpine@sha256:[0-9a-f]{64}")
        self.assertNotIn("/bin/grpc_health_probe", manifest)
        self.assertNotIn(":latest", manifest)

    def test_jenkins_side_effects_require_explicit_parameters(self) -> None:
        build_pipeline = (REPOSITORY_ROOT / "Jenkinsfile").read_text()
        deploy_pipeline = (REPOSITORY_ROOT / "deploy" / "eks" / "Jenkinsfile").read_text()
        readme = (REPOSITORY_ROOT / "README.md").read_text()
        self.assertIn("PUBLISH_IMAGE", build_pipeline)
        self.assertIn("defaultValue: false", build_pipeline)
        self.assertNotIn("git push", build_pipeline)
        self.assertIn("grype", build_pipeline)
        self.assertIn("--config .grype.yaml", build_pipeline)
        self.assertNotIn("trivy", build_pipeline.lower())
        self.assertIn("APPLY", deploy_pipeline)
        self.assertIn("defaultValue: false", deploy_pipeline)
        self.assertNotIn("serverUrl", deploy_pipeline)
        self.assertIn("pin_manifest_images.py", deploy_pipeline)
        self.assertIn("deployment-service.resolved.yml", deploy_pipeline)
        self.assertNotIn(
            "kubectl apply -f deploy/eks/deployment-service.yml", deploy_pipeline
        )
        self.assertNotIn("kubectl apply -f deploy/eks/deployment-service.yml", readme)
        self.assertIn(".jenkins/deployment-service.resolved.yml", readme)
        self.assertFalse(list(REPOSITORY_ROOT.glob("services/*/Jenkinsfile")))

    def test_security_scan_exceptions_are_narrow(self) -> None:
        config = (REPOSITORY_ROOT / ".grype.yaml").read_text()
        self.assertEqual(config.count("  - vulnerability:"), 4)
        self.assertEqual(config.count("      name: python"), 4)
        self.assertEqual(config.count("      version: 3.14.7"), 4)
        self.assertEqual(config.count("      type: binary"), 4)


if __name__ == "__main__":
    unittest.main()
