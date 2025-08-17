---
title: 🍼 First Deployment
weight: 4
---

# First Deployment

Now that we can create a Kubernetes cluster using `kind` we have everything we
need to make our first actual Kubernetes `Deployment`. We will use `kubectl` for
this.

## Applying YAML to a cluster

If we create a file `nginx-deployment.yaml` with the following content

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
```

And then run

```bash
kubectl apply -f nginx-deployment.yaml
```

Then we are telling Kubernetes to make that desired state a reality. In the YAML
we see that we told it to run 3 replicas of the pod. Let's check if it managed
to make that the actual state:

```bash
kubectl get pods
```

This should give an output similar to

```
NAME                                READY   STATUS    RESTARTS   AGE
nginx-deployment-647677fc66-jpp5k   1/1     Running   0          11s
nginx-deployment-647677fc66-tfqbn   1/1     Running   0          11s
nginx-deployment-647677fc66-xqdqc   1/1     Running   0          11s
```

To check that the nginx application is actually running we can use
port-forwarding to port-forward one of the pods to our host machine

```bash
kubectl port-forward deployments/nginx-deployment 8080:80
```

This will make requests to `localhost:8080` on our host machine hit one of the
pods of the deployment, and you should see the default "Welcome to nginx!"
message in the browser.

## Exercises

### Deploy nginx with HTML from configmap

- Run the above example, but where you mount HTML from a configmap.
- Port forward and verify in the browser.

### Deploy the go application

- Use `kind load docker-image` to copy the go application we built earlier to
  your kind cluster.
- Duplicate `exercises/kubernetes/first_deployment/nginx-deployment.yaml` to a
  file called `go-deployment.yaml`.
- Open the file and edit fields so that it will deploy our go application.
- Apply it to the cluster and use port-forwarding to verify that it works.

{{% hint info %}}
The Kubernetes default pull policy is `IfNotPresent` unless the image tag is
`:latest` or omitted (and implicitly `:latest`) in which case the default policy
is `Always`. `IfNotPresent` causes the Kubelet to skip pulling an image if it
already exists. If you want those images loaded into node to work as expected,
please:

- don’t use a `:latest` tag

and / or:

- specify `imagePullPolicy: IfNotPresent` or `imagePullPolicy: Never` on your
  container(s).

See [Kubernetes
imagePullPolicy](https://kubernetes.io/docs/concepts/containers/images/#updating-images)
for more information.
{{% /hint %}}

