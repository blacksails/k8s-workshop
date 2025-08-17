---
title: 📦 kind
weight: 1
draft: false
---

# KIND: Kubernetes in Docker

<img src="kind.png" style="float: right; width: 250px; margin-left: 20px"/>

`kind` is an awesome tool which lets us quickly spin up a k8s cluster where each
server in the cluster is simulated as a docker container. This lets us very
quickly spin up local clusters for testing applications and other purposes.

We haven't interacted with any Kubernetes clusters yet (yes I know we went
through a lot of theory 😅). This ends now, I assume you already got the tools
installed which was mentioned in the chapter about [Essential
Tools](/docs/preparation/essential_tools), if not hurry up and get them
installed.

## Creating a cluster

A cluster can be created using the following command.

```bash
kind create cluster
```

This will start a docker container which runs Kubernetes. In addition it will
also set up the Kubernetes CLI `kubectl`. This enables us to use the CLI to
interact with our local cluster.

If you want multiple separate clusters, you can run the command again and supply
the `--name` flag which will let you give the cluster a name.

```bash
kind create cluster --name=my-cluster
```

## Listing clusters

To get a list of your current kind clusters run

```bash
kind get clusters
```

## Deleting clusters

To delete a cluster supply the name flag:

```bash
kind delete cluster --name=my-cluster
```

## Loading images

The container runtime inside your kind cluster doesn't is in fact not the same
as on your host machine. This comes to show if you build and image and try to
run it in your cluster. Because the cluster isen't running on the same
container runtime as your host docker, it wont be able to find the images you
build on your host. If running something from Dockerhub, this isen't a problem,
as the container runtime will simply fetch the image and run it.

But kind provides a small mechanism which lets us copy images from our host
docker runtime to the container runtime within the kind cluster. To copy a
docker image to your cluster run

```
kind load docker-image --name=my-cluster <image-name>
```

{{% hint warning %}}
There is currently a bug, which causes trouble when using Docker Desktop with
default settings. See the [tracking
issue](https://github.com/kubernetes-sigs/kind/issues/3795) for a workaround.
{{% /hint %}}

## Exercises

### Creating a cluster with multiple nodes

Kind supports a config file format which can be used for more advanced
configurations. One of the options in the config file lets you specify what
nodes your cluster should be made up of.

You can find the documentation for the kind config file
[here](https://kind.sigs.k8s.io/docs/user/configuration/)

- Create a config file which specifies a multi node cluster.
- Create a cluster based on the config file.
- Run `docker ps` to see that you have the nodes you expect.
