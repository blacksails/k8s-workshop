---
title: 📦 Introduction
weight: 1
draft: false
---

# Introduction to Helm

## What is Helm?

Helm is often called "the package manager for Kubernetes" - and for good reason.
Just as you use package managers like `apt`, `brew`, or `npm` to install and
manage software on your system, Helm helps you install and manage applications
on Kubernetes clusters.

At its core, Helm solves a fundamental challenge: deploying complex applications
to Kubernetes often requires multiple YAML files with intricate dependencies,
configurations, and relationships. Managing these manually becomes unwieldy as
applications grow in complexity or when you need to deploy the same application
across different environments with varying configurations.

## How Helm Works

Helm introduces the concept of **Charts** - packages of pre-configured Kubernetes
resources that can be installed as a unit. Think of a chart as a blueprint that
defines how an application should be deployed, including all its components,
dependencies, and configuration options.

Here's how the Helm ecosystem works:

### Charts
A Helm chart is a collection of files that describe a related set of Kubernetes
resources. Charts use templates to generate valid Kubernetes manifests based on
configuration values, making them reusable across different environments and
scenarios.

### Repositories
Helm repositories are collections of charts that can be shared and distributed.
The most popular public repository is [Artifact Hub](https://artifacthub.io/),
which hosts thousands of community-maintained charts for common applications
like databases, web servers, monitoring tools, and more.

### Releases
When you install a chart into a cluster, Helm creates a **release** - a specific
instance of that chart with a unique name. You can have multiple releases of the
same chart in your cluster, each with different configurations.

### Values
Charts are customizable through **values** - configuration parameters that can
be modified without changing the chart itself. This allows the same chart to be
deployed in development, staging, and production with environment-specific
settings.

## Why Use Helm?

**Simplified Deployment**: Instead of managing dozens of YAML files manually,
you can deploy complex applications with a single `helm install` command.

**Reusability**: Charts can be shared and reused across teams and projects,
eliminating the need to recreate deployment configurations from scratch.

**Configuration Management**: Values allow you to customize deployments for
different environments while maintaining a single source of truth for the
application structure.

**Community Ecosystem**: Thousands of pre-built charts are available for popular
applications, meaning you often don't need to create charts from scratch.

## Helm Architecture

Helm v3 (the current version) consists of:

- **Helm CLI**: The command-line tool used to interact with charts and releases
- **Chart Repository**: Storage for packaged charts that can be shared and distributed
- **Kubernetes API**: Helm communicates directly with the Kubernetes API server to manage resources

Unlike Helm v2, there's no server-side component (Tiller) required, making Helm v3
more secure and easier to deploy.

## Exercise: Installing ingress-nginx with Helm

Let's get hands-on experience with Helm by installing the popular ingress-nginx
controller using a Helm chart.

### Prerequisites

You need to have the `helm` CLI installed. If you don't have it, head over to
the [essential tools](../preparation/essential_tools) article, and find the
installation instructions.

### Step 1: Add the Helm Repository

Before we can install charts, we need to add repositories that contain them:

```bash
# Add the ingress-nginx repository
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx

# Update repository information
helm repo update
```

### Step 2: Search for Available Charts

You can search for charts in your configured repositories:

```bash
# Search for nginx-related charts
helm search repo nginx

# Get detailed information about the ingress-nginx chart
helm show chart ingress-nginx/ingress-nginx
```

### Step 3: Install the Chart

Install the ingress-nginx controller:

```bash
# Install ingress-nginx in the ingress-nginx namespace
helm install ingress-nginx ingress-nginx/ingress-nginx \
  --create-namespace \
  --namespace ingress-nginx
```

### Step 4: Verify the Installation

Check that the installation was successful:

```bash
# List Helm releases
helm list --all-namespaces

# Check the pods
kubectl get pods -n ingress-nginx

# Check the service
kubectl get svc -n ingress-nginx
```

Notice that the `Service` is of type `LoadBalancer` this is because the ingress
controller will act as our entry point for HTTP traffic, so it needs a public
IP.

### Step 5: Explore the Release

Helm provides several commands to inspect your release:

```bash
# Get release status
helm status ingress-nginx -n ingress-nginx

# View the values used for this release
helm get values ingress-nginx -n ingress-nginx

# See all manifests that were applied
helm get manifest ingress-nginx -n ingress-nginx
```

### Optional: Customize the Installation

You can customize the installation using values. Create a file called `values.yaml`:

```yaml
controller:
  replicaCount: 2
```

Then upgrade the release:

```bash
helm upgrade ingress-nginx ingress-nginx/ingress-nginx \
  -f values.yaml \
  -n ingress-nginx
```

This should scale the `Deployment` to two replicas.

{{% hint info %}}
The upgrade command has a `--install` flag. This enables us to use the `upgrade`
command to install charts, which means that the command will install the chart
if it isn't installed or upgrade it otherwise.
{{% /hint %}}

**Expected Result**: You should have a working ingress-nginx controller running
in your cluster, and you should understand how Helm simplifies the installation
and management of complex Kubernetes applications.

## Key Takeaways

- Helm packages Kubernetes applications into reusable charts
- Charts use templates and values to create customizable deployments
- Helm repositories provide a way to share and distribute charts
- Releases represent instances of charts deployed to your cluster

In the next section, we'll dive deeper into how Helm templating works and learn
how to create your own charts.
