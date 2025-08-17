---
title: 🔄 Flux Integration
weight: 20
draft: false
---

# Helm and Flux Integration

## GitOps with Helm Charts

While Helm provides excellent package management for Kubernetes applications,
manually running `helm install` and `helm upgrade` commands doesn't align with
GitOps principles. This is where Flux's Helm Controller comes in - it allows
you to manage Helm releases declaratively through Git, combining the power of
Helm's templating with GitOps automation.

## How Flux Manages Helm Charts

Flux uses several Custom Resource Definitions (CRDs) to manage Helm charts:

### HelmRepository

Defines a Helm repository as a source for charts:

```yaml
apiVersion: source.toolkit.fluxcd.io/v1beta1
kind: HelmRepository
metadata:
  name: ingress-nginx
  namespace: flux-system
spec:
  interval: 1h
  url: https://kubernetes.github.io/ingress-nginx
```

### HelmRelease

Defines a Helm release to be installed and managed:

```yaml
apiVersion: helm.toolkit.fluxcd.io/v2beta1
kind: HelmRelease
metadata:
  name: ingress-nginx
  namespace: ingress-nginx
spec:
  interval: 5m
  chart:
    spec:
      chart: ingress-nginx
      version: "4.8.3"
      sourceRef:
        kind: HelmRepository
        name: ingress-nginx
        namespace: flux-system
  values:
    controller:
      replicaCount: 2
      service:
        type: NodePort
```

## Benefits of Flux + Helm Integration

**Declarative Management**: Helm releases are defined as Kubernetes resources
in Git, making them auditable and version-controlled.

**Automatic Updates**: Flux continuously monitors for changes and automatically
applies updates when configurations change in Git.

**Drift Detection**: Flux detects when the actual state differs from the desired
state and automatically reconciles differences.

**Multi-Environment Support**: Different environments can have different values
files while using the same base chart configuration.

## Exercise: Migrating from Manual Helm to Flux

Let's migrate the ingress-nginx installation from the previous chapter to be
managed by Flux instead.

### Prerequisites

Ensure you have:
- A working Flux installation (from the GitOps chapter)
- The ingress-nginx chart currently installed via manual Helm

### Step 1: Uninstall Manual Helm Release

First, let's remove the manually installed ingress-nginx:

```bash
# List current Helm releases
helm list --all-namespaces

# Uninstall the ingress-nginx release
helm uninstall ingress-nginx -n ingress-nginx

# Verify removal
kubectl get pods -n ingress-nginx
```

### Step 2: Create Helm Repository Resource

In your GitOps repository, create a directory structure for Helm configurations:

```bash
# In your GitOps repository
mkdir -p infrastructure/helm-repositories
mkdir -p infrastructure/helm-releases
```

Create the HelmRepository resource:

```bash
cat <<EOF > infrastructure/helm-repositories/ingress-nginx.yaml
apiVersion: source.toolkit.fluxcd.io/v1beta1
kind: HelmRepository
metadata:
  name: ingress-nginx
  namespace: flux-system
spec:
  interval: 1h
  url: https://kubernetes.github.io/ingress-nginx
EOF
```

### Step 3: Create HelmRelease Resource

Create the HelmRelease for ingress-nginx:

```bash
cat <<EOF > infrastructure/helm-releases/ingress-nginx.yaml
apiVersion: helm.toolkit.fluxcd.io/v2beta1
kind: HelmRelease
metadata:
  name: ingress-nginx
  namespace: ingress-nginx
spec:
  interval: 5m
  releaseName: ingress-nginx
  chart:
    spec:
      chart: ingress-nginx
      version: "4.8.3"
      sourceRef:
        kind: HelmRepository
        name: ingress-nginx
        namespace: flux-system
  install:
    createNamespace: true
  values:
    controller:
      replicaCount: 1
EOF
```

### Step 4: Configure Flux to Monitor These Resources

Create or update your Kustomization to include the new Helm resources:

```bash
cat <<EOF > infrastructure/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - helm-repositories/
  - helm-releases/
EOF
```

If you don't already have a Kustomization pointing to the infrastructure
directory, create one in your flux-system namespace:

```bash
cat <<EOF > clusters/c0x/infrastructure.yaml
apiVersion: kustomize.toolkit.fluxcd.io/v1
kind: Kustomization
metadata:
  name: infrastructure
  namespace: flux-system
spec:
  interval: 10m
  sourceRef:
    kind: GitRepository
    name: flux-system
  path: "./infrastructure"
  prune: true
  wait: true
EOF
```

### Step 5: Commit and Push Changes

```bash
# Add and commit the changes
git add .
git commit -m "Add ingress-nginx Helm release managed by Flux"
git push origin main
```

### Step 6: Monitor the Deployment

Watch Flux deploy the Helm release:

```bash
# Monitor Flux components
flux get all

# Watch the HelmRepository
flux get sources helm

# Watch the HelmRelease
flux get helmreleases

# Monitor the actual deployment
kubectl get pods -n ingress-nginx -w
```

### Step 7: Verify the Installation

```bash
# Check Helm releases managed by Flux
helm list --all-namespaces

# Verify the pods are running
kubectl get pods -n ingress-nginx

# Check the service
kubectl get svc -n ingress-nginx
```

**Expected Result**: The ingress-nginx controller should be running, now managed
by Flux instead of manual Helm commands.

## Troubleshooting Flux Helm Integration

Common troubleshooting commands:

```bash
# Check HelmRepository status
kubectl describe helmrepository ingress-nginx -n flux-system

# Check HelmRelease status
kubectl describe helmrelease ingress-nginx -n ingress-nginx

# View Flux Helm Controller logs
kubectl logs -n flux-system deployment/helm-controller

# Force reconciliation
flux reconcile helmrelease ingress-nginx -n ingress-nginx

# Suspend and resume releases
flux suspend helmrelease ingress-nginx -n ingress-nginx
flux resume helmrelease ingress-nginx -n ingress-nginx
```

## Exercise: Use cert-manager to provide SSL certificates for your ingress

[cert-manager](https://cert-manager.io/) is a very popular tool for managing
certificates within your Kubernetes cluster. It handles the whole certificate
life cycle, and makes certificates a thing that you don't really need to worry
about.

1. Install cert-manager using the same approach as we did in the last exercise
   with ingress-nginx.
1. Get DNS setup for your `LoadBalancer` service (Ask for help).
1. Add and `Ingress` with TLS to an application in your cluster.

## Key Takeaways

- Flux HelmController enables GitOps management of Helm charts
- HelmRepository and HelmRelease resources declaratively define chart sources and deployments
- All Helm operations are auditable through Git history

By combining Helm's powerful templating with Flux's GitOps capabilities, you get
the best of both worlds: flexible, reusable application packages with declarative,
auditable deployment management.
