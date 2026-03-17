# CLI Branding

KSWP displays a stylized ASCII banner on startup to provide a professional CLI experience.

## Banner Display

When you run `kswp` without any flags, you'll see:

```
  _                       
 | |                      
 | | _______      ___ __  
 | |/ / __\ \ /\ / / '_ \ 
 |   <\__ \\ V  V /| |_) |
 |_|\_\___/ \_/\_/ | .__/ 
                   | |    
                   |_|  
  v0.1.0
  Kubernetes Cluster Hygiene Tool
```

## Version Information

The banner displays:
- **KSWP Logo**: Stylized ASCII art representation
- **Version**: Current application version (e.g., v0.1.0)
- **Tagline**: "Kubernetes Cluster Hygiene Tool"

## Suppressing the Banner

The banner is automatically suppressed when:
- Using `--help` flag: `kswp --help`
- Using `--version` flag: `kswp --version`
- Requesting help for specific commands: `kswp scan --help`

This ensures that flag-based queries remain concise and focused.

## Version Management

The version is managed in the `internal/branding/banner.go` file and should be updated during releases to reflect the current application version.

### Current Version
- **Version**: 0.1.0
- **Status**: Initial release

Version updates follow [Semantic Versioning](https://semver.org/) (MAJOR.MINOR.PATCH).
