---
title: 🎯 Templating
weight: 10
draft: false
---

# Helm Charts and Templating

## Understanding Helm Charts

A Helm chart is more than just a collection of Kubernetes YAML files - it's a
sophisticated templating system that allows you to create flexible, reusable
application packages. Charts use the Go template language combined with Sprig
functions to generate Kubernetes manifests dynamically based on configuration
values.

## Chart Structure

Every Helm chart follows a standard directory structure:

```
mychart/
├── Chart.yaml          # Chart metadata
├── values.yaml         # Default configuration values
├── charts/             # Chart dependencies
├── templates/          # Template files
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── configmap.yaml
│   ├── _helpers.tpl    # Template helpers
│   └── NOTES.txt       # Post-install notes
└── .helmignore         # Files to ignore when packaging
```

### Chart.yaml

The `Chart.yaml` file contains metadata about your chart:

```yaml
apiVersion: v2
name: webapp
description: A simple web application
version: 0.1.0
appVersion: "1.0"
dependencies:
  - name: redis
    version: "^17.0.0"
    repository: "https://charts.bitnami.com/bitnami"
```

Even though you can specify dependencies in your Chart I would generally advise
against using that feature of Helm. Instead I found it way less troublesome to
just install the chart that your chart depends in a parallel helm release.

### values.yaml

The `values.yaml` file defines the default configuration values for your chart:

```yaml
replicaCount: 1

image:
  repository: nginx
  tag: "1.21"
  pullPolicy: IfNotPresent

service:
  type: ClusterIP
  port: 80

ingress:
  enabled: false
  className: ""
  annotations: {}
  hosts:
    - host: example.local
      paths:
        - path: /
          pathType: Prefix

resources:
  limits:
    cpu: 100m
    memory: 128Mi
  requests:
    cpu: 100m
    memory: 128Mi
```

## Templating with Values

Helm templates use the `{{ }}` syntax to inject values and execute template
functions. Under the hood this is actually the go templating engine, so any
familiarity with that will transfer 1 to 1.

Here's how values from `values.yaml` are used in template files:

### Basic Value Substitution

```yaml
# templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Chart.Name }}
spec:
  replicas: {{ .Values.replicaCount }}
  selector:
    matchLabels:
      app: {{ .Chart.Name }}
  template:
    metadata:
      labels:
        app: {{ .Chart.Name }}
    spec:
      containers:
      - name: {{ .Chart.Name }}
        image: {{ .Values.image.repository }}:{{ .Values.image.tag }}
        imagePullPolicy: {{ .Values.image.pullPolicy }}
        ports:
        - containerPort: 80
        resources:
          {{- toYaml .Values.resources | nindent 10 }}
```

### Conditional Templates

Use `if` statements to conditionally include resources:

```yaml
# templates/ingress.yaml
{{- if .Values.ingress.enabled -}}
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {{ .Chart.Name }}-ingress
  annotations:
    {{- toYaml .Values.ingress.annotations | nindent 4 }}
spec:
  {{- if .Values.ingress.className }}
  ingressClassName: {{ .Values.ingress.className }}
  {{- end }}
  rules:
  {{- range .Values.ingress.hosts }}
  - host: {{ .host }}
    http:
      paths:
      {{- range .paths }}
      - path: {{ .path }}
        pathType: {{ .pathType }}
        backend:
          service:
            name: {{ $.Chart.Name }}-service
            port:
              number: {{ $.Values.service.port }}
      {{- end }}
  {{- end }}
{{- end }}
```

### Template Functions

Helm provides many built-in functions for manipulating values:

```yaml
metadata:
  name: {{ .Chart.Name | lower }}
  labels:
    app: {{ .Chart.Name }}
    version: {{ .Chart.AppVersion | quote }}
    environment: {{ .Values.environment | default "production" }}
    created: {{ now | date "2006-01-02" }}
```

### Helper Templates

Create reusable template snippets in `_helpers.tpl`:

```yaml
{{/* Common labels */}}
{{- define "webapp.labels" -}}
app.kubernetes.io/name: {{ .Chart.Name }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/* Create a default fully qualified app name */}}
{{- define "webapp.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}
```

Use helpers in your templates:

```yaml
metadata:
  name: {{ include "webapp.fullname" . }}
  labels:
    {{- include "webapp.labels" . | nindent 4 }}
```

## Exercise 1: Understanding Template Rendering

Let's explore how templates work by examining an existing chart.

### Step 1: Create a Sample Chart

```bash
# Create a new chart
helm create webapp

# Examine the chart structure
ls -la webapp/
```

### Step 2: Render Templates Locally

When writing or configuring helm charts it is very useful, to be able to see how
the end result changes when different values are provided as input to the chart.
The `template` command renders the helm chart without applying it to your
cluster.

```bash
# Render templates with default values
helm template webapp ./webapp

# Render with custom values
helm template webapp ./webapp --set replicaCount=3
```

### Step 3: Inspect Template Logic

Examine the generated `webapp/templates/deployment.yaml` and identify:
- How values are injected
- Conditional logic usage
- Template functions being used

**Expected Result**: You should understand how Helm transforms templates into
valid Kubernetes manifests using values.

## Exercise 2: Creating a Simple Chart from Scratch

Now let's create a simple chart for a web application from scratch.

### Step 1: Initialize Chart Structure

```bash
# Create directories
mkdir -p simple-webapp/{templates,charts}

# Create Chart.yaml
cat <<EOF > simple-webapp/Chart.yaml
apiVersion: v2
name: simple-webapp
description: A simple web application chart
version: 0.1.0
appVersion: "1.0"
EOF
```

### Step 2: Create values.yaml

```bash
cat <<EOF > simple-webapp/values.yaml
replicaCount: 1

image:
  repository: nginx
  tag: "latest"
  pullPolicy: IfNotPresent

service:
  type: ClusterIP
  port: 80

env:
  APP_NAME: "Simple WebApp"
  DEBUG: "false"
EOF
```

### Step 3: Create Deployment Template

```bash
cat <<EOF > simple-webapp/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Chart.Name }}-{{ .Release.Name }}
  labels:
    app: {{ .Chart.Name }}
    release: {{ .Release.Name }}
spec:
  replicas: {{ .Values.replicaCount }}
  selector:
    matchLabels:
      app: {{ .Chart.Name }}
      release: {{ .Release.Name }}
  template:
    metadata:
      labels:
        app: {{ .Chart.Name }}
        release: {{ .Release.Name }}
    spec:
      containers:
      - name: {{ .Chart.Name }}
        image: {{ .Values.image.repository }}:{{ .Values.image.tag }}
        imagePullPolicy: {{ .Values.image.pullPolicy }}
        ports:
        - containerPort: 80
        env:
        {{- range \$key, \$value := .Values.env }}
        - name: {{ \$key }}
          value: {{ \$value | quote }}
        {{- end }}
EOF
```

### Step 4: Create Service Template

```bash
cat <<EOF > simple-webapp/templates/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: {{ .Chart.Name }}-{{ .Release.Name }}-service
  labels:
    app: {{ .Chart.Name }}
    release: {{ .Release.Name }}
spec:
  type: {{ .Values.service.type }}
  ports:
  - port: {{ .Values.service.port }}
    targetPort: 80
    protocol: TCP
  selector:
    app: {{ .Chart.Name }}
    release: {{ .Release.Name }}
EOF
```

### Step 5: Test Your Chart

```bash
# Validate the chart
helm lint simple-webapp

# Render templates to verify output
helm template test-release simple-webapp

# Install the chart
helm install test-release simple-webapp

# Verify deployment
kubectl get pods,svc -l app=simple-webapp
```

### Step 6: Customize with Values

Create a custom values file:

```bash
cat <<EOF > custom-values.yaml
replicaCount: 2
image:
  repository: httpd
  tag: "2.4"
env:
  APP_NAME: "Custom WebApp"
  DEBUG: "true"
  ENVIRONMENT: "development"
EOF

# Upgrade with custom values
helm upgrade test-release simple-webapp -f custom-values.yaml
```

**Expected Result**: You should have a working chart that deploys a web
application with customizable configuration through values.

## Exercise 3: Advanced Templating Features

Let's enhance our chart with more advanced templating features.

### Step 1: Add Conditional ConfigMap

```bash
cat <<EOF > simple-webapp/templates/configmap.yaml
{{- if .Values.configMap.enabled }}
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ .Chart.Name }}-{{ .Release.Name }}-config
  labels:
    app: {{ .Chart.Name }}
    release: {{ .Release.Name }}
data:
  {{- range \$key, \$value := .Values.configMap.data }}
  {{ \$key }}: {{ \$value | quote }}
  {{- end }}
{{- end }}
EOF
```

### Step 2: Update values.yaml

Add configMap section to values.yaml:

```bash
cat <<EOF >> simple-webapp/values.yaml

configMap:
  enabled: false
  data:
    app.properties: |
      app.name=Simple WebApp
      app.version=1.0
      log.level=INFO
EOF
```

### Step 3: Create Helper Template

```bash
cat <<EOF > simple-webapp/templates/_helpers.tpl
{{/*
Common labels
*/}}
{{- define "simple-webapp.labels" -}}
app: {{ .Chart.Name }}
release: {{ .Release.Name }}
version: {{ .Chart.AppVersion | quote }}
managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Create a default fully qualified app name
*/}}
{{- define "simple-webapp.fullname" -}}
{{- printf "%s-%s" .Chart.Name .Release.Name | trunc 63 | trimSuffix "-" }}
{{- end }}
EOF
```

### Step 4: Update Templates to Use Helpers

Update your deployment.yaml to use the helper:

```bash
# Replace the labels section in deployment.yaml
metadata:
  name: {{ include "simple-webapp.fullname" . }}
  labels:
    {{- include "simple-webapp.labels" . | nindent 4 }}
```

### Step 5: Test Advanced Features

```bash
# Test with configMap enabled
helm template test-release simple-webapp --set configMap.enabled=true

# Install with configMap
helm upgrade test-release simple-webapp --set configMap.enabled=true
```

**Expected Result**: Your chart now includes conditional resources and reusable
helper templates, demonstrating advanced Helm templating capabilities.

## Key Templating Concepts

**Values Hierarchy**: Values can come from multiple sources in order of precedence:
1. Command line `--set` flags
2. Values files specified with `-f`
3. Chart's default `values.yaml`

**Template Functions**: Helm provides 60+ template functions for string manipulation,
type conversion, date formatting, and more.

**Pipelines**: Use `|` to chain functions: `{{ .Values.name | upper | quote }}`

**Whitespace Control**: Use `{{-` and `-}}` to control whitespace in output.

**Scope**: Use `$` to access root scope from within loops: `{{ $.Chart.Name }}`

In the next section, we'll learn how to integrate Helm charts with Flux for
GitOps-based deployments.
