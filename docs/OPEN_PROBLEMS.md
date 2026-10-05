# Open problems and held decisions

This repository documents Jenkins delivery paths for an Amazon EKS-hosted
microservices application. The items below identify evidence and operating
boundaries that remain open or held. Repository configuration, offline checks,
and screenshots do not substitute for current, authorized live evidence.

## Open items

### Live EKS and Jenkins evidence is not maintained here

The Jenkinsfiles, Kubernetes manifest, local validation commands, and documented
operator steps describe how delivery can be exercised. They do not establish a
current Jenkins run, registry scan result, EKS rollout, service health result,
or production outcome. Current run IDs, logs, scan receipts, and rollout
evidence must be retained separately when an authorized operator runs the path;
this inventory does not infer them from source or documentation.

### Secrets and authentication remain an external boundary — BLOCKED-AUTH held

AWS credentials and IAM permissions, EKS cluster access, Jenkins credentials,
registry credentials, and any credential bindings are external operator
prerequisites. This repository must not carry those secrets, and this repository
does not unlock, validate, or bypass them. A live authenticated execution is
therefore held at the authentication boundary until the authorized environment
supplies and records the relevant evidence.

### Environment-specific configuration still requires operator decisions

The AWS account and region, EKS cluster and context, namespaces, registry or
image repositories, Jenkins agents and credential IDs, image versions, and
apply choice are environment-specific. The documented placeholders and example
names are not a portable environment contract and do not prove that a
particular deployment can be applied without that configuration.

### Promotion and deployment evidence is separate from build evidence

A successful build or image scan, if produced by a future Jenkins run, would not
by itself prove that the reviewed manifest was promoted and applied to the
intended EKS environment. Manifest review, apply authorization, rollout
observation, and retained receipts remain separate operational steps.

### No comparative results

The repository is an implementation and documentation example. It doesn't
provide comparative performance, reliability or security-certification
results.

## Held decisions

- Do not commit credentials, tokens, kubeconfigs, Jenkins secrets, or other
  secret material; keep authentication in the approved external systems.
- Do not relabel YAML, screenshots, offline validation, or documentation as
  current live-cluster, Jenkins, or registry evidence.
- Do not apply a manifest or promote images without the intended environment's
  authorized review and retained evidence.
- Do not invent metrics or scores, or infer a production conclusion from this
  page.
- Do not treat this inventory as a substitute for the relevant live evidence
  and operational owner.

## Updating this page

Keep repository, Jenkins and registry, EKS, and authentication evidence
distinct. When an item changes, link the authoritative run record, log, scan
result or rollout record rather than summarizing it as a score. No live EKS,
Jenkins or registry result is inferred from repository artifacts.
