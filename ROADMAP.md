# Roadmap

Planned next steps for this Jenkins and Amazon EKS delivery example. Items come
from the [README](README.md), [`docs/OPEN_PROBLEMS.md`](docs/OPEN_PROBLEMS.md)
and the existing CI. Use [`CHANGELOG.md`](CHANGELOG.md) for work that has
already landed.

Finishing an item here changes the repository only. It doesn't show that a
live Jenkins controller, registry or cluster works; that needs a run on
infrastructure you control.

## Status key

- `UNTESTED`: not implemented yet, or not re-checked on `main`.
- `GAP`: the repository's coverage is incomplete or stale.
- `BLOCKED-AUTH`: needs authorized Jenkins, registry, AWS or EKS access that
  this repository doesn't have.

## Planned work

### Documentation

| Item | Status | Notes |
| --- | --- | --- |
| Refresh `docs/OPEN_PROBLEMS.md` when held boundaries change | **GAP** | Track new gaps there. Closing a roadmap item is not the same as closing an open problem. |
| Record future docs-only changes in `CHANGELOG.md` | **UNTESTED** | Short, factual entries. |

### CI and offline validation

| Item | Status | Notes |
| --- | --- | --- |
| Document GitHub PR-check scope vs Jenkins pipeline scope | **GAP** | [`.github/workflows/pr-checks.yml`](.github/workflows/pr-checks.yml) covers offline policy, language tests and kubeconform. The build, scan, publish and apply steps in the root [`Jenkinsfile`](Jenkinsfile) and [`deploy/eks/Jenkinsfile`](deploy/eks/Jenkinsfile) need an authorized Jenkins run. |
| Extend repository policy tests when pipeline contracts change | **UNTESTED** | [`tests/test_repository_policy.py`](tests/test_repository_policy.py) guards digest pins, workflow triggers and grype exceptions. New delivery constraints should add focused tests. |
| Keep pre-commit and PR-check commands aligned with the README | **GAP** | The README's local validation commands should stay copy-paste accurate as checks change. |

### Live delivery (outside the repository)

| Item | Status | Notes |
| --- | --- | --- |
| Template for recording an authorized Jenkins run | **BLOCKED-AUTH** | What a run should record (build number, scan summary, optional publish tag), without claiming current pipeline success from source. |
| EKS promotion and rollout checklist, separate from build results | **BLOCKED-AUTH** | Manifest review, apply authorization and rollout observation stay operator steps. README screenshots are illustrative, not current state. |
| Pipeline telemetry capture | **BLOCKED-AUTH** | No metrics are published until an authorized environment produces them. |

### Sibling repositories

| Item | Status | Notes |
| --- | --- | --- |
| Share the roadmap layout with [aks-ado-microservices-cicd](https://github.com/T-Py-T/aks-ado-microservices-cicd) | **UNTESTED** | Keep the structure simple so both repos can use it without sharing live-cluster claims. |

## Out of scope for this page

- No score, metric, production outcome, security certification or comparison
  result is claimed.
- No secrets, Jenkins credentials or cluster access are provided.
- Planned work and repository artifacts don't imply a live Jenkins run,
  registry scan or EKS rollout.
