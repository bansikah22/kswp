# kswp Deployment

This document describes the different ways you can deploy and use `kswp`.

## 1. Precompiled Binary

This is the most common way to install `kswp`. We provide precompiled binaries for Linux, macOS, and Windows for different architectures.

You can download the latest binary for your system from the [GitHub Releases](https://github.com/bansikah22/kswp/releases) page.

**Example installation on Linux:**

```bash
wget https://github.com/bansikah22/kswp/releases/download/v0.1.0/kswp_0.1.0_linux_amd64.tar.gz
tar -xvf kswp_0.1.0_linux_amd64.tar.gz
chmod +x kswp
sudo mv kswp /usr/local/bin/kswp
```

**Uninstallation:**

```bash
sudo rm /usr/local/bin/kswp
```

## 2. Docker Container

We provide a Docker image for `kswp` on Docker Hub. This is a great way to run `kswp` in CI/CD pipelines or in an isolated environment.

**Example usage:**

```bash
docker run --rm -v ~/.kube:/root/.kube bansikah/kswp scan
```

**Uninstallation:**

Since the Docker container is run with the `--rm` flag, it is automatically removed after execution. There are no artifacts to uninstall.

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

**Uninstallation:**

```bash
kubectl krew uninstall kswp
```

## 4. Homebrew

If you are on macOS or Linux, you can install `kswp` using [Homebrew](https://brew.sh/).

**Installation:**

```bash
brew tap bansikah22/kswp
brew install kswp
```

**Uninstallation:**

```bash
brew uninstall kswp
```

## 5. Helm Chart

You can deploy `kswp` as a CronJob in your Kubernetes cluster using our Helm chart. This is useful for running scheduled cleanups of your cluster.

**Installation:**

```bash
helm repo add kswp https://bansikah22.github.io/kswp/
helm install kswp-cleaner kswp/kswp
```

**Uninstallation:**

```bash
helm uninstall kswp-cleaner
```

## 7. Apply a Lua script

You can apply a Lua script to filter and delete resources.

**Usage:**

```bash
kswp apply -f <script.lua>
```
