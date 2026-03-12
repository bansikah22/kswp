# Local Development

## Prerequisites

- [Go](https://golang.org/) version 1.18 or higher
- [Docker](https://www.docker.com/) (optional, for running a local Kubernetes cluster)
- [minikube](https://minikube.sigs.k8s.io/docs/start/) or other local Kubernetes cluster

## Running the application

1. **Clone the repository:**
   ```bash
   git clone https://github.com/bansikah22/kswp.git
   cd kswp
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Run the application:**
   You can run the application using `go run`.

   ```bash
   go run main.go --help
   ```

   This will show you the available commands.

   **Scan for unused resources:**
   ```bash
   go run main.go scan
   ```

   **Clean unused resources:**
   ```bash
   go run main.go clean
   ```

   **Sweep unused resources:**
   The `sweep` command is a combination of `scan` and `clean` with an additional `--older-than` flag to filter resources by age.
   ```bash
   go run main.go sweep --older-than 7d
   ```

   **Run in dry-run mode:**
   The `scan`, `clean`, and `sweep` commands support a `--dry-run` flag. This will simulate the operation without making any changes to your cluster.

   ```bash
   go run main.go scan --dry-run
   go run main.go clean --dry-run
   go run main.go sweep --dry-run
   ```

   **Run the TUI:**
   ```bash
   go run main.go tui
   ```

   **Specify a namespace:**
   All commands support a `--namespace` (or `-n`) flag to specify a single namespace to scan. If no namespace is provided, the tool will scan all namespaces.

   ```bash
   go run main.go scan -n my-namespace
   go run main.go clean -n my-namespace
   go run main.go sweep -n my-namespace
   go run main.go tui -n my-namespace
   ```

## Testing

To run the tests, you can use the `go test` command.

```bash
go test ./...
```
