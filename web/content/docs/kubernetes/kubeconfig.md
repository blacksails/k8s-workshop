---
title: ⚙️ CLI
weight: 3
---

# Kubernetes CLI: `kubectl`

The CLI which we use to interact with Kubernetes is called `kubectl`. How to
pronounce the name of the CLI is splitting waters. Some of the most common ones
are

- kubeCTL
- kubecuddle
- kubecontrol

I don't care what you call it, personally I usually go with "kubeCTL".

## Kube Config

The configuration file for the `kubectl` CLI is commonly reffered to as  "kube
config". By default the file is located at `~/.kube/config`.

If you at some point are provided with a "kube config" which you want to use
instead of the what you have in `~/.kube/config`, you can do one of three
things.

- Merge it with the contents of `~/.kube/config`. This can be a bit tricky but
  is certainly doable, as `kubectl` supports multiple contexts (more on this in
  next section).
- Set the `KUBECONFIG` environment variable to the file path of the kube config
  file. This will make `kubectl` use that instead of the default.
- Set the `--kubeconfig` flag to the file path of the kube config file. This
  will make `kubectl` use that instead of the default.

To view your current kube config run

```bash
kubectl config view
```

## Kube Context

The k8s CLI has a context system, which let's us store multiple k8s cluster
connections. Why might we need this? Different environments are often isolated
to different clusters. When we use the `kubectl` CLI, the commands will reach
the Kubernetes cluster which our context currently points to.

To list your contexts run

```bash
kubectl config get-contexts
```

The context which is currently active will be marked by an asterisk in the
*CURRENT* column. If your current context isn't the one you expected you can
change it using

```bash
kubectl config use-context <context-name>
```

You view the kube config only showing what's relevant to your current context

```bash
kubectl config view --minify
```
