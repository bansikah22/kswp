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

Clean unused resources. Notifications will be sent automatically if configured in `~/.kswp/config.yaml`.

## `kswp sweep`

Sweep unused resources.

- `--older-than`: filter resources older than a duration (e.g., 7d, 24h)

Notifications will be sent automatically if configured in `~/.kswp/config.yaml`.

## `kswp tui`

Terminal UI for kswp.

## `kswp graph`

Display a dependency graph of resources.

## `kswp apply`

Apply a Lua script to filter and delete resources.

- `-f, --file`: path to the Lua script file

## `kswp doctor`

Check the health of the Kubernetes cluster.

## `kswp notify slack`

Send a Slack notification.

- `--webhook-url`: Slack webhook URL
- `--message`: Message to send

## `kswp notify email`

Send an email notification.

- `--message`: Message to send

## Configuring Notifications

To automate notifications, create a configuration file at `~/.kswp/config.yaml` with the following format:

```yaml
notifications:
  - type: slack
    webhook_url: "your-slack-webhook-url"
  - type: email
    email:
      from: "your-email@example.com"
      to:
        - "recipient1@example.com"
        - "recipient2@example.com"
      password: "your-email-password"
      host: "smtp.example.com"
      port: "587"
```

You can configure multiple notifiers of different types. When you run `kswp clean` or `kswp sweep`, a notification will be sent to each configured notifier.

You can test your notification setup with the `kswp notify` command. For example, to test your Slack integration, run:

```bash
kswp notify slack --webhook-url "your-slack-webhook-url" --message "Hello from kswp!"
```

To test your email integration, make sure you have it configured in `~/.kswp/config.yaml` and then run:

```bash
kswp notify email --message "Hello from kswp!"
```
