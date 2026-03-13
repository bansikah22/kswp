# kswp
[![Release](https://github.com/bansikah22/kswp/actions/workflows/release.yml/badge.svg)](https://github.com/bansikah22/kswp/actions/workflows/release.yml)

kswp is a Kubernetes cluster hygiene tool that detects and safely cleans unused resources.

## Installation

For detailed installation instructions, please see the [deployment documentation](docs/deployment.md).

## Features

- **Scan for unused resources**: `kswp` can scan for unused ConfigMaps, Secrets, Services, ReplicaSets, Jobs, and Pods.
- **Reasoning Engine**: `kswp` will tell you *why* a resource is considered unused.
- **Interactive UI**: A terminal UI to visualize and manage unused resources.
- **Cluster Hygiene Score**: Get a score that represents the health of your cluster.
- **Dependency Graph**: Visualize the dependencies between your Kubernetes resources.
- **Notifications**: Get notified on Slack or Email when a scan is complete.

For a full list of commands and their flags, see the [command documentation](docs/commands.md).

## Contributing

Contributions are welcome! Please see the [contributing guidelines](CONTRIBUTING.md) for more information.

## Release Process

Please see the [release process documentation](docs/release.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
