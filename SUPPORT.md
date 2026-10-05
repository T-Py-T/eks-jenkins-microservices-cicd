# Support

This repository is a Jenkins and Amazon EKS microservices delivery example, not a
hosted service. Support is limited to the repository's documented pipelines,
configuration guidance, tests, and reproducible local checks on `main`.

## Before opening an issue

- Confirm the problem still occurs on the current `main` branch.
- Include the relevant path, commit, command, and a minimal reproduction.
- Redact credentials, tokens, kubeconfigs, personal data, and private endpoint
  details from logs and configuration.
- For a Jenkins or EKS result, include only sanitized evidence that you are
  authorized to share; repository artifacts do not substitute for live
  environment evidence.

## Operational boundary

The repository does not provide an operated EKS cluster, Jenkins controller or
project, registry, or production support channel. Environment-specific
credentials, IAM permissions, credential bindings, and deployment decisions
remain with the authorized operator. Do not commit secrets or request them in
an issue.

## Limits

A support response or local validation result doesn't certify a pipeline,
cluster, image, provider or deployment.

For known gaps and held decisions, see
[Open problems and held decisions](docs/OPEN_PROBLEMS.md).
