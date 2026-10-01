# Hireability and discoverability

> Tip-cite bank: base main `60b2798b` + this PR pending Steward; paired with aks-ado Ship 239;
> provenance only; never `READY`.

Pair note: mirrors aks-ado Ship 239 hireability lean; provenance only.

This page orients staffing readers to inspectable delivery-engineering choices.
It is not an acceptance gate, scorecard, release declaration, or `READY` signal.

## What / why / how

| Question | Short answer |
| --- | --- |
| **What** | Jenkins + EKS delivery example around Online Boutique with separate build and deploy paths. |
| **Why** | Show parameterized per-service CI and explicit operator-controlled apply without baking credentials into git. |
| **How** | Root `Jenkinsfile` tests/scans/builds one service; `deploy/eks/Jenkinsfile` validates a selected manifest and requires `APPLY=true`. |

## What this proves

This repository is a concrete, inspectable delivery-engineering example. It is
not a claim of employment readiness and does not substitute for an interview or
an evaluation rubric. See [`OPEN_PROBLEMS.md`](OPEN_PROBLEMS.md) for the active
evidence gaps and held decisions that bound this narrative; it is not a
scorecard or `READY` gate.

A reviewer can inspect evidence of the ability to:

- work across an eleven-service application spanning Go, C#, Node.js, Python,
  Java, and Locust;
- design a parameterized Jenkins path that tests and audits a selected service,
  scans source and the resulting image, builds a container, and keeps image
  publishing explicitly optional;
- separate build and deployment concerns, validate Kubernetes resources, pin
  registry tags to digests, and require an explicit `APPLY=true` before cluster
  changes;
- document repeatable local checks for pre-commit, service tests, dependency
  audits, manifest conformance, and an isolated integration smoke test; and
- keep registry, AWS, and Kubernetes credentials in Jenkins rather than in the
  repository.

What this proves is limited to the engineering decisions and repeatable checks
that this repository documents. It does not prove production uptime, capacity,
security certification, a live-cluster outcome, or any hiring/evaluation score.
No score is assigned or implied.

## Stack

- **CI:** root `Jenkinsfile` (parameterized service build/scan; publish optional)
- **Deploy:** `deploy/eks/Jenkinsfile` + EKS manifests (explicit `APPLY=true`)
- **App:** eleven-service Online Boutique
- **Local checks:** pre-commit, `scripts/test-service.sh`, kubeconform / policy tests

## Suggested GitHub topics

`aws`, `eks`, `jenkins`, `kubernetes`, `terraform`, `devsecops`, `cicd`,
`microservices`, `trivy`, `online-boutique`

Topics aid search only; they do not certify results or readiness.

## License

Repository-specific pipeline, deployment, and documentation work is under the
[MIT License](../LICENSE). Online Boutique source retains Google LLC
[Apache License 2.0](../LICENSE-APACHE-2.0) notices; see
[`THIRD_PARTY_NOTICES.md`](../THIRD_PARTY_NOTICES.md).

## Related docs

| Document | Role |
| --- | --- |
| [../README.md](../README.md) | Architecture, delivery flow, local validation |
| [../SECURITY.md](../SECURITY.md) | Vulnerability reporting |
| [../CONTRIBUTING.md](../CONTRIBUTING.md) | Contribution and tip-cite rules |
| [OPEN_PROBLEMS.md](OPEN_PROBLEMS.md) | Held evidence boundaries |
| [README.md](README.md) | Documentation index |
