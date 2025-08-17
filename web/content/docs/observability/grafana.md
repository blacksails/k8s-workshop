---
title: 📈 Grafana
weight: 3
draft: false
---

# Introduction to Grafana

## What is Grafana?

Grafana is an open-source analytics and interactive visualization web application.
It provides charts, graphs, and alerts when connected to supported data sources.
Grafana is widely used for monitoring and observability, particularly in
combination with time-series databases like Prometheus.

Grafana excels at taking raw metrics and transforming them into meaningful,
actionable dashboards that help teams understand system behavior, identify
trends, and respond to issues quickly.

## How Grafana Works

Grafana follows a simple architecture:

1. **Data Sources**: Grafana connects to various data sources (Prometheus,
   InfluxDB, Elasticsearch, etc.)
2. **Queries**: It queries these data sources using their native query languages
3. **Visualization**: Raw data is transformed into various visual formats
   (graphs, tables, heatmaps, etc.)
4. **Dashboards**: Multiple visualizations are organized into dashboards
5. **Alerting**: Alerts can be configured based on query results and thresholds

## Key Features

**Multi-Data Source Support**: Connect to dozens of different data sources
including Prometheus, InfluxDB, Elasticsearch, and cloud services.

**Rich Visualizations**: Create line graphs, bar charts, heatmaps, tables,
gauges, and more to represent your data.

**Dashboard Management**: Organize visualizations into dashboards that can be
shared, exported, and version controlled.

**Alerting**: Set up alerts that notify you when metrics cross thresholds or
meet specific conditions.

**User Management**: Control access with organizations, teams, and role-based
permissions.

## Exercise 1: Access Grafana from the Kube-Prometheus-Stack

The Grafana instance was installed as part of the kube-prometheus-stack in the
previous section. Let's access it and explore its features.

### Prerequisites

Ensure you have completed the Prometheus exercises and have the
kube-prometheus-stack running in your cluster.

### Step 1: Verify Grafana is Running

```bash
# Check that Grafana pod is running
kubectl get pods -n monitoring -l app.kubernetes.io/name=grafana

# Check the Grafana service
kubectl get svc -n monitoring -l app.kubernetes.io/name=grafana
```

### Step 2: Access the Grafana UI

```bash
# Port-forward to access Grafana UI
kubectl port-forward -n monitoring svc/kube-prometheus-stack-grafana 3000:80
```

Open your browser to `http://localhost:3000`.

If you got ingress working, you can also opt to use the public URL.

Log into Grafana using the credentials you put in the sealed secret.

### Step 3: Explore the Pre-configured Data Source

1. Click on **Configuration** (gear icon) in the left sidebar
2. Click on **Data Sources**
3. You should see a Prometheus data source already configured
4. Click on it to see the configuration details

Notice that the Prometheus URL is automatically configured to point to the
Prometheus service within the cluster.

**Expected Result**: You should be logged into Grafana with a Prometheus data
source already configured and ready to use.

## Exercise 2: Explore Pre-built Dashboards

The kube-prometheus-stack comes with many pre-built dashboards that provide
comprehensive monitoring for Kubernetes.

### Step 1: Browse Available Dashboards

1. Click on **Dashboards** (four squares icon) in the left sidebar
2. Click **Browse**
3. Explore the available folders and dashboards

You'll see folders like:
- **General** - Basic Grafana dashboards
- **Kubernetes / Compute Resources** - Pod and container metrics
- **Kubernetes / Networking** - Network-related metrics
- **Node Exporter** - Host-level metrics

### Step 2: Explore the Kubernetes Cluster Overview

1. Navigate to **General** > **Kubernetes / Compute Resources / Cluster**
2. This dashboard provides a high-level view of cluster resource usage
3. Explore the different panels:
   - CPU utilization across the cluster
   - Memory usage patterns
   - Network I/O
   - Storage metrics

### Step 3: Dive into Node-Level Metrics

1. Go to **Node Exporter** > **Node Exporter / Nodes**
2. This dashboard shows detailed metrics for each node:
   - CPU usage by core
   - Memory breakdown
   - Disk I/O and space usage
   - Network interface statistics

### Step 4: Examine Pod-Level Metrics

1. Navigate to **Kubernetes / Compute Resources** > **Kubernetes / Compute Resources / Pod**
2. Use the namespace dropdown to filter for the `monitoring` namespace
3. Explore metrics for your monitoring stack components

**Expected Result**: You should understand how pre-built dashboards provide
comprehensive monitoring coverage for different levels of the Kubernetes stack.

## Exercise 3: Create a Custom Dashboard

Let's create a custom dashboard focused on the monitoring stack itself.

### Step 1: Create a New Dashboard

1. Click the **+** (plus) icon in the left sidebar
2. Select **Dashboard**
3. Click **Add new panel**

### Step 2: Create a Panel for Prometheus Query Rate

1. In the query editor, enter:
   ```promql
   rate(prometheus_http_requests_total{job="kube-prometheus-stack-prometheus"}[5m])
   ```
2. Set the panel title to "Prometheus HTTP Request Rate"
3. Change the visualization type to **Time series**
4. In the right panel, set:
   - **Unit**: ops (operations per second)
   - **Legend**: {{method}} {{handler}}
5. Click **Apply**

### Step 3: Add a Panel for Grafana Metrics

1. Click **Add panel**
2. Enter this query:
   ```promql
   grafana_http_request_duration_seconds_bucket{job="kube-prometheus-stack-grafana"}
   ```
3. Set the panel title to "Grafana Request Duration"
4. Change visualization to **Heatmap**
5. Click **Apply**

### Step 4: Add a Panel for AlertManager Status

1. Click **Add panel**
2. Enter this query:
   ```promql
   up{job="kube-prometheus-stack-kube-prom-alertmanager"}
   ```
3. Set the panel title to "AlertManager Status"
4. Change visualization to **Stat**
5. In the right panel, configure:
   - **Unit**: short
   - **Color scheme**: Green-Red
   - **Thresholds**: 0.5 (red below, green above)
6. Click **Apply**

### Step 5: Save Your Dashboard

1. Click the save icon (disk) at the top
2. Name your dashboard "Monitoring Stack Health"
3. Select or create a folder called "Custom"
4. Click **Save**

**Expected Result**: You should have created a custom dashboard with three
panels monitoring the health of your monitoring stack components.

## Exercise 4: Set Up a Simple Alert

Let's create an alert that fires when any component in the monitoring stack
goes down.

### Step 1: Create an Alert Rule

1. In your custom dashboard, click on the "AlertManager Status" panel
2. Click **Edit** (pencil icon)
3. Go to the **Alert** tab
4. Click **Create Alert**

### Step 2: Configure the Alert Condition

1. Set the alert name: "Monitoring Component Down"
2. In the query section, the query should already be populated:
   ```promql
   up{job="kube-prometheus-stack-kube-prom-alertmanager"}
   ```
3. In the conditions section:
   - **Is Below**: 0.5
   - **For**: 1m (alert after being down for 1 minute)

### Step 3: Configure Alert Notifications

1. In the **Notifications** section
2. Add a message: "AlertManager component is down!"
3. For now, leave notification channels empty (in production, you'd configure
   email, Slack, etc.)

### Step 4: Save the Alert

1. Click **Save** at the top
2. Go back to the dashboard view

### Step 5: Test the Alert (Optional)

To test the alert, you could temporarily scale down the AlertManager:

```bash
# Scale down AlertManager (this will trigger the alert)
kubectl scale statefulset alertmanager-kube-prometheus-stack-kube-prom-alertmanager -n monitoring --replicas=0

# Wait a moment, then scale it back up
kubectl scale statefulset alertmanager-kube-prometheus-stack-kube-prom-alertmanager -n monitoring --replicas=1
```

**Expected Result**: You should have a working alert that monitors the health
of your monitoring components.

## Key Concepts

**Data Sources**: External systems that Grafana queries for data (Prometheus,
InfluxDB, etc.).

**Dashboards**: Collections of panels that display visualizations of your data.

**Panels**: Individual visualizations within a dashboard (graphs, tables,
gauges, etc.).

**Queries**: Requests to data sources written in the data source's native query
language (PromQL for Prometheus).

**Variables**: Dynamic values that can be used in queries to make dashboards
interactive and reusable.

**Alerts**: Notifications triggered when metrics meet specific conditions.

## Key Takeaways

- Grafana transforms raw metrics into meaningful visualizations
- The kube-prometheus-stack provides a complete monitoring solution with
  pre-configured integration between Prometheus and Grafana
- Pre-built dashboards cover most common Kubernetes monitoring needs
- Custom dashboards allow you to focus on metrics specific to your applications
- Grafana's alerting system helps you stay informed about system health
- The combination of Prometheus and Grafana is a foundational observability
  stack for Kubernetes

## Next Steps

With Prometheus collecting metrics and Grafana visualizing them, you now have a
solid foundation for Kubernetes observability. Consider exploring:

- **Custom Metrics**: Instrument your applications to expose custom metrics
- **Advanced Alerting**: Set up more sophisticated alerts with notification
  channels (email, Slack, PagerDuty)
- **Log Integration**: Add log aggregation with tools like Loki
- **Distributed Tracing**: Implement tracing with Jaeger or Zipkin
- **Dashboard Templates**: Use variables and templating for more dynamic
  dashboards
- **RBAC**: Configure role-based access control for different team members
