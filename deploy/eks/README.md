# EKS deployment

`deployment-service.yml` defines the eleven Online Boutique services, Redis,
their service discovery names, resource requests, and health probes.

The Jenkins pipeline validates the manifest on every run. It contacts a
cluster only when an operator starts it with `APPLY=true` and supplies a
Jenkins Kubernetes credential. No endpoint or credential is stored here.

Before deployment, replace every `BUILD_NUMBER` placeholder with a versioned
image produced by a successful service build. Validate the template locally:

```bash
go run github.com/yannh/kubeconform/cmd/kubeconform@v0.8.0 \
  -strict -summary deployment-service.yml
```

The deployment pipeline uses Skopeo to resolve those tags to digests. To run
the same step from an authenticated workstation:

```bash
python3 ../../scripts/pin_manifest_images.py \
  --input deployment-service.yml \
  --output deployment-service.resolved.yml
kubectl apply --dry-run=server -f deployment-service.resolved.yml
kubectl apply -f deployment-service.resolved.yml
kubectl rollout status deployment/frontend
```
