# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-03-13

### Added

- **Scan for unused resources**: `kswp` can scan for unused ConfigMaps, Secrets, Services, ReplicaSets, Jobs, and Pods.
- **Reasoning Engine**: `kswp` will tell you *why* a resource is considered unused.
- **Interactive UI**: A terminal UI to visualize and manage unused resources.
- **Cluster Hygiene Score**: Get a score that represents the health of your cluster.
- **Dependency Graph**: Visualize the dependencies between your Kubernetes resources.
