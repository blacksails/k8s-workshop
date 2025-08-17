---
title: 🛂 Namespaces
weight: 5
---

# Namespaces

Kubernetes provides a way of logically separating resources. Kubernetes
resources are either namespaced or cluster-scoped. Namespaces themselves are a
resource, but namespaces are not a namespaced resource, which means that
namespaces can not be put into namespaces. `Namespace` is a cluster-scoped
resource. On the other hand all of resources we went through in [the
introduction](./introduction) are namespaced. Meaning that they will always live
in the scope of a namespace.

A namespace can quickly be created using `kubectl`

```bash
kubectl create namespace my-namespace
```

A namespace can also be created by writing out the YAML resource and applying it
to the cluster like we did with the `Deployment` earlier.

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: my-namespace
```

## Specifying namespace

The `Deployment` we applied earlier didn't specify a namespace, which caused
Kubernetes to put it in the default namespace. If we want the deployment to live
in a namespace other than the default, we can do one of two things.

- Set the namespace explicitly in the metadata
- Tell `kubectl` to use a specific namespace

If we set the namespace explicitly in the deployment resource, it will look
something like this

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  namespace: my-namespace
  labels:
    app: nginx
spec:
  replicas: 3
  # ...
```

Instead explicitly setting the namespace in the `Deployment` we could also have
used the `--namespace` flag when applying the resource.

```bash
kubectl apply --namespace my-namespace -f nginx-deployment.yaml
```

## Getting the namespace name inside a pod

The following will serve as an example of three things:

- [Kubernetes downward
API](https://kubernetes.io/docs/concepts/workloads/pods/downward-api/)
- Configuring a docker container with dynamic input

We will use the downward API to set an environment variable in the pod
containing the name of the namespace in which the pod is running.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: echo-namespace
  labels:
    app: echo-namespace
spec:
  replicas: 3
  selector:
    matchLabels:
      app: echo-namespace
  template:
    metadata:
      labels:
        app: echo-namespace
    spec:
      containers:
      - name: echosvc
        image: blacksails/echosvc
        args:
        - "Hello from the namespace called '$(POD_NAMESPACE)'"
        env:
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        ports:
        - containerPort: 8080
```

The downward API let's us read fields in the resource and set their value as
environment variables. The image we are running is a simple service which
replies with the message that it is configured with. You can find the source
code under `exercises/kubernetes/namespaces`.

Depending on what namespace we deploy the above example to, we will get
different results when we send a request.

## Exercises

### Echo namespace in different namespaces

- Apply the echo-namespace `Deployment` to two different namespaces.
- Port-forward each of the deployments to a different port on your host machine
  and verify that you get a different result depending on which of the ports you
  hit.

{{% hint info %}}
When running `kubectl port-forward` be sure to specify the namespace using the
`-n`/`--namespace` flag when you reference the deployment. Otherwise `kubectl`
will assume the default namespace.
{{% /hint %}}

### Echo a message from a configmap

Instead of echoing the namespace message, change the message so that it echos
another environment variable set from a configmap.

{{% hint info %}}
You can reuse the way that the `POD_NAMESPACE` variable is set. `valueFrom`
excepts a configmap source. Take a look at the Kubernetes API reference.
{{% /hint %}}
