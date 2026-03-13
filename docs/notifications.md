# Notifications

The `kswp notify` command can be used to send notifications to various channels. The following channels are currently supported:

- Slack
- Email

## Configuration

Notifications are configured in the `~/.kswp/config.yaml` file. The following is an example of a config file that configures a slack and an email notifier:

```yaml
notifications:
  - type: slack
    webhook_url: "https://hooks.slack.com/services/..."
  - type: email
    email:
      from: "user@example.com"
      to:
        - "user@example.com"
      password: "password"
      host: "smtp.gmail.com"
      port: "587"
```

## Usage

Once you have configured your notifiers, you can send a notification by running the following command:

```bash
kswp notify slack --message "Hello from kswp"
```

or

```bash
kswp notify email --message "Hello from kswp"
```
