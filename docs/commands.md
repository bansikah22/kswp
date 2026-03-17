# kswp Commands

This document describes all the available commands in `kswp`.

## Global Flags

These flags are available for all commands:

- `--dry-run`: run in dry-run mode
- `-n, --namespace`: specify the namespace to scan
- `--exclude-namespaces`: comma-separated list of namespaces to exclude (e.g., `kube-system,kube-public`)
- `--label`: filter resources by label (e.g., 'app=nginx')
- `--name`: filter resources by name

## Resource Exclusion

You can protect critical resources from being scanned or deleted by marking them with the `kswp.io/exclude: "true"` annotation. See [Resource Exclusion via Annotations](resource-exclusion.md) for detailed information and examples.

## Namespace Filtering

### Scanning Specific Namespace
Scan resources in a single namespace:
```bash
kswp scan -n my-namespace
kswp clean -n my-namespace
```

### Excluding Namespaces
Skip system-critical namespaces:
```bash
kswp scan --exclude-namespaces kube-system,kube-public
kswp clean --exclude-namespaces kube-system,kube-public,kube-node-lease
```

When using both `-n` and `--exclude-namespaces`, if the namespace matches an excluded namespace, no scanning occurs.

## `kswp scan`

Scan for unused resources.

## `kswp clean`

Clean unused resources.
Add support for TTL annotations, e.g., cleaner/ttl: 24h.
You can specify which resource types to clean by using the following flags. If no flags are provided, all unused resources will be cleaned.

- `--all`: clean all unused resources (default)
- `--configmaps`: clean unused configmaps
- `--secrets`: clean unused secrets
- `--services`: clean unused services
- `--replicasets`: clean unused replicasets
- `--jobs`: clean unused jobs
- `--pods`: clean unused pods
- `--pvcs`: clean unused persistentvolumeclaims
- `--ttl`: clean expired resources based on the cleaner/ttl annotation

**Examples:**

Clean only unused ConfigMaps and Secrets:
```bash
kswp clean --configmaps --secrets
```

Clean a specific PVC by name:
```bash
kswp clean --pvcs --name my-pvc
```

## `kswp sweep`

Sweep unused resources.

- `--older-than`: filter resources older than a duration (e.g., 7d, 24h)

## `kswp tui`

Terminal UI for kswp.

## `kswp graph`

Display a dependency graph of resources.

## `kswp apply`

Apply a Lua script to filter and delete resources.

- `-f, --file`: path to the Lua script file

## `kswp doctor`

Check the health of the Kubernetes cluster.
