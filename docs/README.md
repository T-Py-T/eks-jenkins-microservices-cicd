# Documentation index

> Tip-cite bank: base main `9374c43a` + Ship 244 pending Steward resolve; provenance only; never `READY`.

**Status:** index page. This directory collects supporting documentation for the
repository. It is not a scorecard, readiness declaration, or invented score.

## Purpose of `/docs`

The `/docs` directory holds evidence-boundary and planning records that
supplement the root README. Use it for held decisions, open evidence gaps, and
documentation that should stay separate from pipeline and deployment source.
Diagrams under `docs/img/` illustrate the documented path; they are not live
deployment evidence.

## Repository documents

Links below point to documents at the repository root or under `.github/` when
they exist in this repository:

| Document | Path | Role |
| --- | --- | --- |
| Contributing | [`CONTRIBUTING.md`](../CONTRIBUTING.md) | Contribution guidance; not a `READY` gate. |
| Security | [`SECURITY.md`](../SECURITY.md) | Vulnerability reporting boundary; not a readiness declaration. |
| Maintainers | [`MAINTAINERS.md`](../MAINTAINERS.md) | Factual owner record; not a scorecard. |
| Notice | [`NOTICE.md`](../NOTICE.md) | Attribution and provenance; not a `READY` claim. |
| CODEOWNERS | [`.github/CODEOWNERS`](../.github/CODEOWNERS) | Review routing; not a scorecard. |
| Funding | [`.github/FUNDING.yml`](../.github/FUNDING.yml) | Sponsorship pointer; not a `READY` gate. |
| Roadmap | [`ROADMAP.md`](../ROADMAP.md) | Planned evidence-backed work; not a release declaration. |
| Pull request template | [`.github/PULL_REQUEST_TEMPLATE.md`](../.github/PULL_REQUEST_TEMPLATE.md) | PR scaffold; not a `READY` gate. |

## In this directory

- [`OPEN_PROBLEMS.md`](OPEN_PROBLEMS.md) — active evidence gaps and held decisions;
  not a scorecard or `READY` gate.

## Explicit non-claims

- No index entry is a `READY` declaration.
- No score, metric, production outcome, or security certification is asserted.
- No live Jenkins run, registry scan, or EKS rollout is inferred from this
  page.
- The tip-cite above is only a trace pointer; the Steward resolves it against
  `main` after merge. It is not approval and never implies `READY`.
