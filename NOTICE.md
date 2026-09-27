# Notice

> Tip-cite bank: base main `f2f8c983` + this PR. Steward resolves after merge;
> provenance only; no `READY` claim.

**Status:** attribution and provenance record. This page names upstream sources,
repository boundaries, and related notice documents for this Jenkins and Amazon
EKS microservices CI/CD lab. It is not a scorecard, security certification,
readiness declaration, or invented score.

## Upstream application

The service code under `services/` is derived from
[GoogleCloudPlatform/microservices-demo](https://github.com/GoogleCloudPlatform/microservices-demo)
(Online Boutique). Google-authored files retain their copyright notices and
Apache License 2.0 terms. See [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md)
for license scope and third-party attribution detail.

## Repository-specific work

Pipeline definitions, Kubernetes manifests, deployment automation, tests, and
documentation in this repository are repository-specific work under the
[MIT License](LICENSE). That license applies to this delivery path; it does not
replace upstream or dependency licenses attached to application code or
third-party packages.

## Tooling and platform references

This lab documents a Jenkins pipeline path, container image build and scan
steps, and an operator-controlled EKS deployment workflow. References to
Jenkins, Amazon EKS, Docker, Grype, Skopeo, Kubernetes, `kubectl`, AWS, and
Podman name platforms and tools used in the documented example. Mention of a
platform or tool is descriptive only; it is not an endorsement, certification,
or `READY` claim for any provider, registry, cluster, or scan result.

## Provenance boundary

Repository files, offline validation, and documentation establish what is
present in source control. They do not establish current Jenkins build numbers,
registry receipts, image scan outcomes, EKS rollout health, or operator
authorization. Live evidence remains with the authorized environment and must be
recorded separately. See [docs/OPEN_PROBLEMS.md](docs/OPEN_PROBLEMS.md) for
active evidence gaps and held decisions.

## Secrets and sensitive material

Do not commit or request AWS credentials, Docker registry passwords, Jenkins
credential exports, `kubectl` kubeconfigs, pipeline tokens, or personal data in
issues or pull requests. This notice does not unlock, validate, or bypass
authentication boundaries.

## Related documents

- [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md) — upstream and dependency
  license attribution.
- [LICENSE](LICENSE) — MIT terms for repository-specific work.
- [SECURITY.md](SECURITY.md) — vulnerability reporting boundary; separate from
  this attribution record.
- [docs/OPEN_PROBLEMS.md](docs/OPEN_PROBLEMS.md) — evidence gaps and held
  decisions; not a scorecard or `READY` gate.

## Explicit non-claims

- No attribution or provenance statement here is a `READY` declaration.
- No score, metric, production outcome, or security certification is asserted.
- No live Jenkins run, registry scan, or EKS rollout is inferred from this page.
- The tip-cite above is only a trace pointer; the Steward resolves it against
  `main` after merge. It is not approval and never implies `READY`.
