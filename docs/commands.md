# kswp Commands

This document describes all the available commands in `kswp`.

## Global Flags

These flags are available for all commands:

- `--dry-run`: run in dry-run mode
- `-n, --namespace`: specify the namespace to scan
- `--label`: filter resources by label (e.g., 'app=nginx')

## `kswp scan`

Scan for unused resources.

## `kswp clean`

Clean unused resources.

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
