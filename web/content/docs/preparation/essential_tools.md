---
title: 🛠️ Essential tools
weight: 7
---

# Essential tools

In order for us to get started and be efficient from the get go we need to
prepare our environment with the tools we need in order to learn and explore
Kubernetes (will from here on be abbreviated as k8s).

## Text Editor

You will need a text editor, use the one you know the best. If you don't have
any in particular I would recommend using [Visual Studio
Code](https://code.visualstudio.com/).

## Docker

The most important tool we will be using is Docker. Docker is used for building
and running container images. We will also be using docker for running k8s
locally.

For most people [Docker
Desktop](https://www.docker.com/products/docker-desktop/) is the best and
easiest way of running docker. On some corporate (windows) laptops it might be
easier to run docker directly in WSL (I won't go into detail about how to set
this up as instructions might differ from place to place)

If you can I would recommend running docker desktop as this abstracts the
complexity of configuring docker with correct settings, networking etc. You can
verify that docker has been installed and is running correctly using this
command: `docker ps` this should return the following output:

```
CONTAINER ID   IMAGE     COMMAND   CREATED   STATUS    PORTS     NAMES
```

This might look a bit strange but this is a list of the containers you are
currently running. In the example above, the list is empty as we haven't
started any containers yet.

If you get an error like

```
Cannot connect to the Docker daemon at unix:///Users/blacksails/.docker/run/docker.sock. Is the docker daemon running?
```

This means that docker is not running and you will need to start it. For Docker
Desktop this is simply a matter of opening the application.

## Kubectl

In order to interact with kubernetes clusters the main tool is called
`kubectl`. This is a CLI which we use to query state and make changes to our
cluster. The `kubectl` installation instructions can be found
[here](https://kubernetes.io/docs/tasks/tools/).

In order to test wether you got `kubectl` installed, you can run the following
command: `kubectl version`. It should show an output similar to the following:

```
Client Version: v1.33.1
Kustomize Version: v5.6.0
The connection to the server localhost:8080 was refused - did you specify the right host or port?
```

Don't mind the message about connection to the server. This is simply because
we haven't set it up to connect to a kubernetes cluster yet (more on this
later).

## Helm

Helm is a "package manager" for k8s which aims to make it easier to install
things in your k8s cluster. We will be using it through out the workshop and
also look into how helm works during the workshop.

You can find the installation instructions
[here](https://helm.sh/docs/intro/install/).

## Flux

Flux is an implementation of the GitOps pattern (more on this later).
[Installation instructions](https://fluxcd.io/flux/installation/#install-the-flux-cli)

## Go toolchain

Go is a programming language widely used in the k8s ecosystem. `docker` and
`kubectl` are actually written in go. But don't worry you will not need to
learn the language in order to complete this course, we will mainly use it in
order to easily install tools that we need.

The go installation instructions can be found [here](https://go.dev/doc/install).

In addition to the official installation instructions I would recommend that
you explicitly set the `GOPATH` and `GOBIN` environment variables and add
`GOBIN` to the `PATH` environment variable.

For most people this means adding a few lines to `~/.zshrc`. Open the file and
paste the following lines at the bottom:

```bash
export GOPATH="$HOME/go"
export GOBIN="$GOPATH/bin"
export PATH="$PATH:$GOBIN"
```

This tells the `go` tool chain where to find source code and where to install
binaries. By adding `GOBIN` to the `PATH` environment variable we ensure that
our shell is able to find the tools we install using the `go` toolchain.

When installing tools using the `go` tool chain, we will be using the `go
install <tool>` command This will download the source code, compile a binary
and place it in the `GOBIN` directory.

## KIND

In order to easily run kubernetes locally we will be using a tool called
`kind`. The name of this tool is an abbreviation of *Kubernetes in Docker*
(more on this later).

With the `go` toolchain installed we can easily install `kind` by running the
following command:

```bash
go install sigs.k8s.io/kind@v0.29.0
```

## Cluster API

Cluster API is an API for managing kubernetes clusters as kubernetes resources
(more on this later).

```bash
go install sigs.k8s.io/cluster-api/cmd/clusterctl@latest
```

## Additional tools

The following are additional optional tools which are easy to install using the
`go` tool chain.

## yq

`yq` is a small tool which can be used to transform `yaml`

```bash
go install github.com/mikefarah/yq/v4@latest
```

## stern

`stern` is a tool which makes it easier to view logs in k8s clusters.

```bash
go install github.com/stern/stern@latest
```

## k9s

`k9s` is a terminal UI (TUI) for viewing and and manipulating kubernetes
clusters. I don't personally use this, but I know a lot of people like it a
lot. Think of it as a dynamic user interface compared to the CLI nature that
`kubectl` has.A

```bash
go install github.com/derailed/k9s@latest
```

## kubeseal (Sealed Secrets)

`kubeseal` is a CLI for which can be used to encrypt secrets, so that they are
safe to store in git. Wether we will have time to dig into this is a bit
uncertain.

```bash
go install github.com/bitnami-labs/sealed-secrets/cmd/kubeseal@main
```
