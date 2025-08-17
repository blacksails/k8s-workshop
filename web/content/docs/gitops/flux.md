---
title: 🔥 Flux
weight: 10
draft: false
---

# Introduction to Flux

Flux is a CNCF graduated project that provides a set of Kubernetes controllers
for implementing GitOps workflows. It automatically synchronizes your cluster
state with configuration stored in Git repositories, ensuring your deployments
stay consistent and up-to-date.

## How Flux Works

Flux operates on a simple but powerful principle: it watches Git repositories
for changes and automatically applies those changes to your Kubernetes cluster.
Here's how the process works:

1. **Source Controller**: Monitors Git repositories, OCI artifacts, and Helm
   repositories for changes
2. **Kustomize Controller**: Applies Kubernetes manifests using Kustomize
3. **Helm Controller**: Manages Helm releases based on HelmRelease resources
4. **Image Automation**: Automatically updates container image tags in Git when
   new images are available

When you push changes to your Git repository, Flux detects these changes within
minutes and applies them to your cluster, maintaining the desired state
automatically.

## Core Components

Flux consists of several specialized controllers that work together:

- **Source Controller**: Handles different types of sources (Git, OCI, Helm)
- **Kustomize Controller**: Reconciles Kustomize configurations
- **Helm Controller**: Manages Helm chart deployments
- **Notification Controller**: Sends alerts about deployment status
- **Image Reflector/Automation Controllers**: Handle automatic image updates

## Key Concepts

### GitRepository

Defines a Git repository as a source for Kubernetes manifests. One instance of
this resource will be created during the flux bootstrap process. We could create
more instances of this if our cluster needed to install resources from different
git repositories.

```yaml
apiVersion: source.toolkit.fluxcd.io/v1
kind: GitRepository
metadata:
  name: webapp-source
  namespace: flux-system
spec:
  interval: 1m
  url: https://github.com/example/webapp-config
  ref:
    branch: main
```

### Kustomization

Applies Kubernetes manifests from a source. The name might be a bit misleading
(or maybe not depending on how you look at it), it refers to being able to
install [kustomizations](https://kustomize.io/) which is a built-in feature of
Kubernetes for installing multiple resources in different configurations.
However the `Kustomization` resource can also install regular Kubernetes YAML
files.

It works by referencing a `GitRepository` which holds the YAML, and then
specifing the path within the repository which holds the YAML files we want to
have applied in our cluster. If we add or remove files from this folder, flux
will pick it up and apply them to our cluster.

```yaml
apiVersion: kustomize.toolkit.fluxcd.io/v1
kind: Kustomization
metadata:
  name: webapp
  namespace: flux-system
spec:
  interval: 10m
  sourceRef:
    kind: GitRepository
    name: webapp-source
  path: "./manifests"
  prune: true
```

### HelmRelease

As we haven't dug into helm just yet, we wont go into much more detail about
this as of now, but just know that Flux supports installing Helm releases.

```yaml
apiVersion: helm.toolkit.fluxcd.io/v2
kind: HelmRelease
metadata:
  name: nginx
  namespace: default
spec:
  interval: 5m
  chart:
    spec:
      chart: nginx
      version: "1.0.0"
      sourceRef:
        kind: HelmRepository
        name: nginx-repo
```

## Installation

We cannot use gitops for installing flux because we would need flux to be
installed to use gitops to install flux. 🐥🥚

This means we will initially use the `flux` CLI to bootstrap our flux setup.

```bash
flux bootstrap github \
  --owner=your-username \
  --repository=your-repo \
  --branch=main \
  --path=./clusters/c0x
```

Before running this command, we need to ensure that our current kube config
points to the cluster in which we would like to install flux. In addition we
will also need a github repo which will be our gitops repo.

## Exercises

### Bootstrap Flux

In the first exercise we will bootstrap the flux installation. Bootstrap Flux in
your Kubernetes cluster and verify it's working.

**Steps**:
1. If you haven't already, create a repository on github for flux to fetch
   resources from. The repository can be public, as we won't be putting any IP
   in there. If you feel adventurous and/or want to get a bit of a deeper learning,
   you can also opt to make it a private repository. The information you need
   for this can be found
   [here](https://fluxcd.io/flux/installation/bootstrap/github/)
1. Run `flux check --pre` to verify your cluster is ready
1. Bootstrap Flux into your cluster
1. Verify all Flux controllers are running with `kubectl get pods -n flux-system`

**Expected Result**: All Flux pods should be in Running state.

### Deploy an Application with `Kustomization`

**Objective**: Use Flux to deploy a simple application from a Git repository.

**Steps**:
1. *Optional* Create a Kustomization resource that applies manifests from another folder in
   the repository.
1. Add Kubernetes manifests (YAML files containing resources), to a directory
   included by a `Kustomization`. You can reuse some of the manifests we applied
   to our local kind clusters earlier.
1. Verify the application is deployed using `kubectl`
1. Use `flux get kustomizations` to monitor reconciliation status

**Expected Result**: Your application should be running in the cluster, and Flux
should show the resources as reconciled.

<!--
### Exercise 4: Multi-Environment Setup

**Objective**: Configure Flux to manage multiple environments (staging and production).

**Steps**:
1. Create directory structure in Git: `clusters/staging/` and `clusters/production/`
2. Configure different Kustomization resources for each environment
3. Use Kustomize overlays to customize deployments per environment
4. Deploy the same application with different configurations to each environment

**Expected Result**: Same application running with environment-specific configurations in both staging and production namespaces.
-->

### *Optional:* GitOps Workflow with Pull Requests

**Objective**: Practice a complete GitOps workflow using Git branches and pull requests.

**Steps**:
1. Create a feature branch in your Git repository
1. Modify application configuration (e.g., change image tag or replicas)
1. Open a pull request and review the changes
1. Merge the pull request and observe Flux applying the changes
1. Use `flux get kustomizations` to monitor reconciliation status

**Expected Result**: Changes should be automatically deployed to the cluster after merging the pull request.

## Troubleshooting

Common commands for debugging Flux issues:

```bash
# Check Flux status
flux get all

# View reconciliation logs
flux logs --level=error

# Suspend and resume reconciliation
flux suspend kustomization webapp
flux resume kustomization webapp

# Force reconciliation
flux reconcile kustomization webapp --with-source
```
