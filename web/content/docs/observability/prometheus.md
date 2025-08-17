---
title: 📊 Prometheus
weight: 2
draft: false
---

# Introduction to Prometheus

## What is Prometheus?

Prometheus is an open-source monitoring and alerting toolkit originally built
at SoundCloud. Since its inception in 2012, it has become the de facto standard
for monitoring in cloud-native environments, particularly in Kubernetes
ecosystems.

Prometheus is designed around a **pull-based model** where it actively scrapes
metrics from configured targets at given intervals. This approach differs
significantly from traditional push-based monitoring systems and provides
several advantages in dynamic, containerized environments.

## How Prometheus Works

Prometheus follows a simple yet powerful architecture:

1. **Discovery**: Prometheus discovers targets to monitor through service
   discovery mechanisms (Kubernetes API, DNS, static configuration, etc.)
2. **Scraping**: It pulls metrics from these targets over HTTP at regular
   intervals
3. **Storage**: Metrics are stored in a time-series database with efficient
   compression
4. **Querying**: PromQL (Prometheus Query Language) allows complex queries and
   aggregations
5. **Alerting**: Rules can trigger alerts based on metric thresholds or patterns

## Prometheus vs Traditional Approaches

### Traditional: StatsD and Push-Based Systems

Traditional monitoring often relies on push-based systems like StatsD:

- **Applications push metrics** to a central collector
- **Requires instrumentation** in every application to send metrics
- **Network dependencies** - if the collector is down, metrics are lost
- **Difficult service discovery** in dynamic environments
- **Less efficient** for high-cardinality data

### Prometheus: Pull-Based Approach

Prometheus takes a different approach:

- **Prometheus pulls metrics** from application endpoints
- **Simple HTTP endpoints** expose metrics in text format
- **Built-in service discovery** works seamlessly with Kubernetes
- **Efficient storage** optimized for time-series data
- **Rich query language** (PromQL) for analysis and alerting
- **Better for dynamic environments** where services come and go frequently

### Key Advantages of Prometheus

**Service Discovery**: Automatically discovers new services in Kubernetes
without manual configuration.

**Pull Model Benefits**: Prometheus controls when and how often to scrape,
making it more reliable and allowing for better rate limiting.

**Rich Metadata**: Metrics include labels that provide rich context and enable
powerful aggregations.

**Efficient Storage**: Custom time-series database optimized for monitoring
workloads.

**Ecosystem Integration**: Deep integration with Kubernetes and cloud-native
tools.

## Exercise 1: Install sealed-secrets

We are now starting to install applications which will need secrets. For this we
will use sealed-secrets as we looked at in an earlier exercise.

Sealed secrets should be installed using `HelmRepository` and `HelmRelease` like
we have done for other components as well. We need this installed as early as
possible because our helm release of the kube-prometheus-stack chart will depend
on secrets being available.

## Exercise 2: Install the Kube-Prometheus-Stack

The kube-prometheus-stack is a collection of Kubernetes manifests, Grafana
dashboards, and Prometheus rules combined with documentation and scripts to
provide easy to operate end-to-end Kubernetes cluster monitoring.

### Step 1: Create the HelmRepository

Create a Flux HelmRepository for the Prometheus community charts:

```yaml
apiVersion: source.toolkit.fluxcd.io/v1beta2
kind: HelmRepository
metadata:
  name: prometheus-community
  namespace: flux-system
spec:
  interval: 1h
  url: https://prometheus-community.github.io/helm-charts
```

Save this into the existing `infrastructure/helm-repositories`

### Step 2: Create a secret containing the grafana credentials

In order to create a secret through gitops we need to use sealed-secrets. If you
didn't get to try out sealed secrets yet, now is the time. You need to create a
sealed secret based on the following secret:

```
apiVersion: v1
kind: Secret
metadata:
  name: grafana-credentials
  namespace: monitoring
data:
  admin-user: <base64 encoded value>
  admin-password: <base 64 encoded value>
```

You will need to pick a username and a password. Don't change the keys in the
secret as that will break the grafana installation. Verify that the sealed
secret is decrypted in the monitoring namespace.

{{% hint info %}}
You might need to create the namespace before you can create the sealed secret.
Creating both of them at the same time is OK, as Flux understands that it should
always create namespaces before other resources.
{{% /hint %}}

### Step 3: Install the Kube-Prometheus-Stack

Create a HelmRelease to install the complete monitoring stack:

```bash
apiVersion: helm.toolkit.fluxcd.io/v2beta1
kind: HelmRelease
metadata:
  name: kube-prometheus-stack
  namespace: flux-system
spec:
  interval: 10m
  targetNamespace: monitoring
  chart:
    spec:
      chart: kube-prometheus-stack
      version: '>=45.0.0'
      sourceRef:
        kind: HelmRepository
        name: prometheus-community
        namespace: flux-system
  values:
    grafana:
      admin:
        existingSecret: grafana-credentials
      ingress:
        enabled: true
        annotations:
          # after having made sure that stating certificate can be issued this
          # can be changed to production
          cert-manager.io/cluster-issuer: letsencrypt-staging
        hosts:
        # set the below hostnames to match your cluster number
        - grafana.c0x.k8s-workshop.dev
        tls:
        - hosts:
          - grafana.c0x.k8s-workshop.dev
          secretName: grafana-tls
```

### Step 4: Verify the Installation

Check that all components are running:

```bash
# Check the HelmRelease status
kubectl get helmrelease kube-prometheus-stack -n flux-system

# Check all pods in the monitoring namespace
kubectl get pods -n monitoring

# Check the services
kubectl get svc -n monitoring
```

**Expected Result**: You should see Prometheus, Grafana, AlertManager, and the
Prometheus Operator all running in the monitoring namespace.

## Exercise 2: Explore the Prometheus Instance

Now let's explore the Prometheus instance that was created by the stack.

### Step 1: Access the Prometheus UI

```bash
# Port-forward to access Prometheus UI
kubectl port-forward -n monitoring svc/kube-prometheus-stack-prometheus 9090:9090
```

Open your browser to `http://localhost:9090` to access the Prometheus web UI.

### Step 2: Explore Targets and Service Discovery

In the Prometheus UI:

1. Go to **Status > Targets** to see what endpoints Prometheus is scraping
2. Notice how many targets are automatically discovered:
   - Kubernetes API server
   - Node exporters
   - kube-state-metrics
   - Prometheus itself
   - And many more!

### Step 3: Try Some Basic Queries

Go to **Graph** and try these queries:

```promql
# Check which targets are up
up

# CPU usage across all nodes
100 - (avg by (instance) (rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)

# Memory usage percentage
(1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100

# Number of pods by namespace
count by (namespace) (kube_pod_info)

# Container restart rate
rate(kube_pod_container_status_restarts_total[5m])
```

### Step 4: Explore Service Monitors

Check what ServiceMonitors were created:

```bash
# List all ServiceMonitors
kubectl get servicemonitor -n monitoring

# Examine a specific ServiceMonitor
kubectl describe servicemonitor kube-prometheus-stack-kube-state-metrics -n monitoring
```

**Expected Result**: You should see a fully functional Prometheus instance with
automatic discovery of Kubernetes components and rich metrics collection.

## Key Concepts

**Prometheus Operator**: A Kubernetes operator that simplifies the deployment
and management of Prometheus instances using custom resources.

**ServiceMonitor**: Defines how Prometheus should scrape metrics from Kubernetes
services.

**Pull Model**: Prometheus actively fetches metrics from targets rather than
having targets push metrics to it.

**Time Series**: Each metric is stored as a time series with labels that provide
context and enable powerful queries.

## Key Takeaways

- Prometheus uses a pull-based model that's well-suited for dynamic environments
- The Prometheus Operator makes it easy to manage Prometheus in Kubernetes
- ServiceMonitors define how services should be scraped for metrics
- Prometheus provides a powerful query language (PromQL) for analysis
- The operator approach follows Kubernetes best practices for managing complex
  applications

In the next section, we'll set up Grafana to visualize the metrics that
Prometheus is collecting.
