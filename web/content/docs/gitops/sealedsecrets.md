---
title: 🔐 Sealed Secrets
weight: 20
draft: false
---

# Introduction to Sealed Secrets

One of the biggest challenges in GitOps is managing secrets securely. While we
can store all our Kubernetes manifests in Git repositories, we cannot store
plain-text secrets there due to security concerns. This creates a gap in our
"everything in Git" philosophy of GitOps.

Sealed Secrets, developed by Bitnami, solves this problem by providing a way to
encrypt secrets that can be safely stored in Git repositories and automatically
decrypted by your Kubernetes cluster.

## The Secret Management Problem in GitOps

In traditional GitOps workflows, you might have configuration like this:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: database-credentials
data:
  username: YWRtaW4=  # base64 encoded 'admin'
  password: UGFzc3dvcmQxMjM=  # base64 encoded 'Password123'
```

The problem is that base64 encoding is not encryption - anyone with access to
your Git repository can easily decode these values. This means you either:

1. **Store secrets outside Git** - breaking GitOps principles
2. **Use a separate secret management tool** - adding complexity
3. **Accept the security risk** - not recommended

## How Sealed Secrets Works

Sealed Secrets introduces two key components:

1. **Controller**: Runs in your cluster and holds a private key
2. **kubeseal CLI**: Encrypts secrets using the controller's public key

Here's the workflow:

1. You create a regular Secret manifest
2. Use `kubeseal` to encrypt it into a SealedSecret
3. Store the SealedSecret in Git (it's safe - only your cluster can decrypt it)
4. The controller watches for SealedSecrets and automatically creates the
   corresponding Secret resources

The encrypted SealedSecret looks like this:

```yaml
apiVersion: bitnami.com/v1alpha1
kind: SealedSecret
metadata:
  name: database-credentials
spec:
  encryptedData:
    username: AgBy3i4OJSWK+PiTySYZZA9rO43cGDEQAx...
    password: AgAKAoiQm0XBmKUCOTfdPGdKnVL4n4OTU...
  template:
    metadata:
      name: database-credentials
```

## Why Sealed Secrets is Perfect for GitOps

Sealed Secrets aligns perfectly with GitOps principles:

- **Everything in Git**: Both application config and secrets are version controlled
- **Declarative**: SealedSecrets are Kubernetes resources like any other
- **Immutable**: Changes require creating new SealedSecrets (proper audit trail)
- **Secure**: Only the target cluster can decrypt the secrets
- **Simple**: No external dependencies or complex PKI infrastructure

The controller automatically handles the lifecycle - when you update a
SealedSecret in Git, Flux detects the change and applies it, the controller
decrypts it and updates the corresponding Secret.

## Key Concepts

### SealedSecret Resource

A SealedSecret is a Kubernetes Custom Resource that contains encrypted data:

```yaml
apiVersion: bitnami.com/v1alpha1
kind: SealedSecret
metadata:
  name: my-secret
  namespace: default
spec:
  encryptedData:
    key1: encrypted_value_1
    key2: encrypted_value_2
  template:
    metadata:
      name: my-secret
      labels:
        app: myapp
```

### Scoping

Sealed Secrets supports different scoping levels:

- **strict** (default): Secret can only be unsealed in the same namespace and with the same name
- **namespace-wide**: Secret can be unsealed in the same namespace with any name
- **cluster-wide**: Secret can be unsealed in any namespace with any name

### Key Rotation

The controller automatically handles key rotation to maintain security while
ensuring existing SealedSecrets continue to work.

## Exercises

### Install Sealed Secrets Controller with Flux

**Objective**: Install the Sealed Secrets controller in your cluster using Flux
to manage it via GitOps.

**Steps**:

1. Create a new directory in your GitOps repository for sealed secrets:
   `mkdir -p clusters/c0x/sealed-secrets`
1. Download the sealed-secrets installation manifest from the [releases
   page](https://github.com/bitnami-labs/sealed-secrets/releases) and place it
   in the directory
4. Commit and push the files to your GitOps repository
5. Verify the controller is installed:
   ```bash
   kubectl get pods -n kube-system | grep sealed-secrets
   ```

**Expected Result**: The sealed-secrets-controller pod should be running in the
kube-system namespace.

### Create and Deploy a Sealed Secret

**Objective**: Create a sealed secret for a database connection and deploy it
using GitOps.

**Steps**:

1. If you haven't installed the `kubeseal` tool already, head over to the
   [essential tools](../preparation/essential_tools#kubeseal-sealed-secrets)
   article for installation instructions.
2. Create a regular Secret manifest
   ```yaml
   # temp-secret.yaml
   apiVersion: v1
   kind: Secret
   metadata:
     name: database-credentials
     namespace: default
   data:
     username: YWRtaW4=          # admin
     password: UGFzc3dvcmQxMjM=  # Password123
   ```
   {{% hint warning %}}
   ⚠️ This file should not be committed to your git repo. Never commit a
   `Secret` to your git repository.
   {{% /hint %}}
3. Encrypt the secret using `kubeseal`:
   ```bash
   kubeseal -f temp-secret.yaml -w clusters/c0x/database-sealed-secret.yaml
   ```
4. Delete the temporary unencrypted file:
   ```bash
   rm temp-secret.yaml
   ```
5. Examine the generated `SealedSecret`:
   ```bash
   cat clusters/c0x/database-sealed-secret.yaml
   ```
6. Commit and push the `SealedSecret` to your repository
7. Verify the secret was created in your cluster:
   ```bash
   kubectl get sealedsecrets
   kubectl get secrets
   kubectl get secret database-credentials -o yaml
   ```

**Expected Result**:
- A SealedSecret should exist in your cluster
- A corresponding Secret should be automatically created
- The Secret should contain the correct decrypted values

## Troubleshooting

Common commands for debugging Sealed Secrets issues:

```bash
# Check controller status
kubectl get pods -n kube-system | grep sealed-secrets
kubectl logs -n kube-system -l name=sealed-secrets-controller

# Verify SealedSecret resources
kubectl get sealedsecrets
kubectl describe sealedsecret <name>

# Check if secrets are being created
kubectl get secrets
kubectl describe secret <name>

# Test encryption manually
kubeseal --fetch-cert > cert.pem
echo -n "test" | kubeseal --raw --from-file=/dev/stdin --name myname --namespace default
```
