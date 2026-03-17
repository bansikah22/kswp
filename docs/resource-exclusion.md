## Resource Exclusion via Annotations

This feature allows users to protect specific Kubernetes resources from being scanned or deleted by kswp using annotations.

### Usage

Mark any resource with the following annotation to exclude it from scanning:

```yaml
metadata:
  annotations:
    kswp.io/exclude: "true"
```

### Example

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: important-config
  namespace: default
  annotations:
    kswp.io/exclude: "true"
data:
  app.properties: |
    setting=value
```

Resources with this annotation will be:
- Skipped during `kswp scan`
- Skipped during `kswp clean`
- Skipped when running the TUI
- Skipped during TTL cleanup

### Annotation Details

- **Key**: `kswp.io/exclude`
- **Value**: Must be exactly `"true"` (case-sensitive)
- **Scope**: Applies to all resource types (ConfigMaps, Secrets, Services, Jobs, ReplicaSets, Pods, PVCs)

### Implementation Details

The exclusion check is implemented in `internal/scanner/exclude.go` and is called by all scanner functions before adding resources to the unused/expired lists. The implementation is:
- **Self-documenting**: Function name `ShouldExclude()` clearly indicates its purpose
- **Defensive**: Safely handles nil annotations
- **Efficient**: Simple string comparison with early returns
