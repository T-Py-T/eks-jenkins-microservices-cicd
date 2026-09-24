# Contributing

> Tip-cite: main `a36909d7` + PR #15.

Thanks for helping improve this Jenkins and Amazon EKS delivery example. Keep
changes focused on the repository-specific pipelines, deployment manifests,
tests, and documentation around the upstream Online Boutique application.

## Before opening a pull request

- Start from the current `main` branch and use a short branch name that
  describes the change.
- Keep one concern per pull request and explain which service, Jenkins path, or
  Kubernetes resource is affected.
- Do not commit credentials, registry passwords, generated images, local
  environments, or editor settings.
- When changing upstream-derived source or assets, preserve their notices and
  update `THIRD_PARTY_NOTICES.md` when provenance changes.

## Local validation

Install the repository hooks and run the policy checks:

```bash
pre-commit run --all-files
```

Run the focused service checks relevant to your change:

```bash
for service in \
  adservice cartservice checkoutservice currencyservice emailservice \
  frontend loadgenerator paymentservice productcatalogservice \
  recommendationservice shippingservice; do
  ./scripts/test-service.sh "$service"
done
```

Validate Kubernetes resources without contacting a cluster:

```bash
go run github.com/yannh/kubeconform/cmd/kubeconform@v0.8.0 \
  -strict -summary deploy/eks/deployment-service.yml
```

If you change the full-stack path, run the documented Podman integration
smoke test when the local container engine is available:

```bash
CONTAINER_ENGINE=podman ./scripts/integration-smoke.sh
```

## Pull requests

1. Describe the change, its operational assumptions, and the validation you
   ran.
2. Include tests or documentation updates for changed behavior.
3. For Jenkins or manifest changes, call out credentials, image references,
   registry settings, and any environment-specific values that reviewers must
   verify.
4. Keep secrets in Jenkins credentials or other approved runtime configuration,
   never in the repository.
5. Wait for the pull-request checks to pass before requesting merge.
