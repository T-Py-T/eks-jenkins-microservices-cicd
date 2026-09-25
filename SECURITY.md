# Security policy

> Tip-cite: base main `80f8a882` + PR #48. Steward resolves after merge; this pointer is not approval and never `READY`.

## Supported code

The current `main` branch is the only supported version. This repository is an EKS Jenkins microservices CI/CD lab for local and lab infrastructure; it does not operate a hosted service.

## Report a vulnerability

Do not open a public issue for an unpatched vulnerability.

Email the repository owner at [tnt850910@aol.com](mailto:tnt850910@aol.com). When the repository Security tab offers it, you may also use GitHub's private vulnerability reporting.

Include the affected commit, the vulnerable path, the impact, and the smallest reproduction that does not expose sensitive data. You can expect an acknowledgment within seven days. A fix schedule depends on the severity and the affected component.

## Keep reports and evidence safe

- Do not send or commit API keys, Docker registry credentials, AWS keys, Kubernetes tokens, Jenkins secrets, or personal data.
- Do not attach unredacted cluster manifests, kubeconfig files, Jenkins credential exports, or production registry contents.
- Use synthetic credentials, local fixtures, and isolated lab clusters when reproducing pipeline or deployment defects.
- Treat captured third-party output under its original license and terms.

## Repository boundary

Registry, AWS, and Kubernetes credentials belong in Jenkins or your local environment, not in the repository. Never commit a secret value, kubeconfig, private registry token, or unredacted deployment artifact.

The root `Jenkinsfile`, `deploy/eks/Jenkinsfile`, and local validation scripts are operator tools for building, scanning, and optionally applying manifests on infrastructure you control. They do not certify third-party container images, cloud providers, or generated artifacts as secure.

The GitHub pull-request workflow and local `pre-commit` checks use dependency audits, service tests, and manifest validation. Local and CI checks do not certify a Jenkins agent, EKS cluster, container image, or generated change as secure.

## Evidence boundaries

The [open problems and held decisions](docs/OPEN_PROBLEMS.md) inventory records
active evidence gaps and held decisions; it is not a scorecard, `READY` gate,
or security certification. This policy makes no `READY` claim and does not
infer live security evidence from repository artifacts.
