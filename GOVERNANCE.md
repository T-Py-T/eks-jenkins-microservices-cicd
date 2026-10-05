# Governance

This repository is a Jenkins and Amazon EKS delivery example with a single
owner, [@T-Py-T](https://github.com/T-Py-T) (see [MAINTAINERS.md](MAINTAINERS.md)).

## How changes are made

1. Changes arrive as pull requests against `main`, one concern per pull
   request.
2. The pull-request checks must pass before merge.
3. The owner reviews and decides whether to merge.
4. Notable changes are recorded in [CHANGELOG.md](CHANGELOG.md). Planned work
   lives in [ROADMAP.md](ROADMAP.md), and known gaps in
   [docs/OPEN_PROBLEMS.md](docs/OPEN_PROBLEMS.md).

## What the repository can and can't show

Repository files, checks, manifests and pull requests show the intended
implementation. They aren't approval for, or proof of, a live deployment.
Jenkins runs, registry provenance, EKS context, rollout health and operator
authorization must be verified in the environment being changed.

Credentials stay in Jenkins or other approved external systems, never in this
repository.
