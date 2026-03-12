# kswp Deployment

This document describes the different ways you can deploy and use `kswp`.

## 1. Precompiled Binary

This is the most common way to install `kswp`. We provide precompiled binaries for Linux, macOS, and Windows for different architectures.

You can download the latest binary for your system from the [GitHub Releases](https://github.com/bansikah22/kswp/releases) page.

**Example installation on Linux:**

```bash
wget https://github.com/bansikah22/kswp/releases/latest/download/kswp-linux-amd64
chmod +x kswp-linux-amd64
sudo mv kswp-linux-amd64 /usr/local/bin/kswp
```

## 2. Docker Container

We provide a Docker image for `kswp` on Docker Hub. This is a great way to run `kswp` in CI/CD pipelines or in an isolated environment.

**Example usage:**

```bash
docker run --rm -v ~/.kube:/root/.kube bansikah22/kswp scan
```

## 3. kubectl Plugin (Krew)

You can install `kswp` as a `kubectl` plugin using [Krew](https://krew.sigs.k8s.io/).

**Installation:**

```bash
kubectl krew install kswp
```

**Usage:**

```bash
kubectl kswp scan
```

## 4. Homebrew

If you are on macOS or Linux, you can install `kswp` using [Homebrew](https://brew.sh/).

**Installation:**

```bash
brew tap bansikah22/kswp
brew install kswp
```

## 5. Helm Chart

You can deploy `kswp` as a CronJob in your Kubernetes cluster using our Helm chart. This is useful for running scheduled cleanups of your cluster.

**Installation:**

```bash
helm repo add kswp https://bansikah22.github.io/kswp/
helm install kswp-cleaner kswp/kswp
```
