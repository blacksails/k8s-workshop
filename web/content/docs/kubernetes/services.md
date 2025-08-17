---
title: 🫡 Services
weight: 6
---

# Services

We briefly touched upon services earlier in [the introduction](./introduction).
We will now explore services a little bit deeper.

## Endpoints

When you create a service, you will be able to see that it picked up pods by
looking at the endpoints. The following command should show that your service
has found a list of IPs which correspond to the IPs of the pods that it is
matching.

```bash
kubectl get endpointslice
```

## Service type

The Kubernetes `Service` resource can be of different types. If we don't specify
the type it will default to `ClusterIP`. You can read more about other service
types in [the sevice
documentation](https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types)

## Exercises

### Service load balancing in action

In this exercise we will deploy an application which simply writes a log each
time it gets a request. The source code for the application can be found in
`exercises/kubernetes/services`. I made it easy for you to deploy by supplying
the `Deployment` resource for you at
`exercises/kubernetes/services/deployment.yaml`

- Apply the `deployment.yaml` to the cluster
- Fill out the specification of the service in
  `exercises/kubernetes/services/service.yaml`, and apply it to the cluster.
- Use the `stern` CLI to watch the logs of the deployment
  `stern deployment/http-logger`.
- Run `kubectl proxy`. This enables us to send requests to the service from our
  host machine.
- Send requests to the service using
  `curl http://localhost:8001/api/v1/namespaces/default/services/http-logger:80/proxy/`
- Notice how different pods are hit in the stern log output.
