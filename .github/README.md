# GitHub Actions Workflows

This directory contains GitHub Actions workflows for the Shamir CLI project.

## Workflows

### 1. Test Workflow (`test.yml`)

**Triggers:**
- Pull requests to `main` or `master` branches
- Pushes to `main` or `master` branches

**What it does:**
- Sets up Go 1.21 environment
- Downloads and verifies dependencies
- Checks code formatting with `gofmt`
- Runs `go vet` for static analysis
- Executes all tests with race detection
- Runs benchmarks
- Generates code coverage report
- Uploads coverage to Codecov (optional)
- Builds the application to verify compilation

**Purpose:** Ensures code quality and functionality before merging changes.

### 2. Release Workflow (`release.yml`)

**Triggers:**
- Pushes to `main` or `master` branches
- Git tags starting with `v*` (e.g., `v1.0.0`)

**What it does:**
1. **Test Phase:** Runs all tests to ensure quality
2. **Build Phase:** 
   - Builds binaries for multiple platforms:
     - Linux (amd64, arm64)
     - macOS/Darwin (amd64, arm64)
     - Windows (amd64)
   - Creates archives (`.tar.gz` for Unix, `.zip` for Windows)
   - Includes documentation files in archives
3. **Release Phase:**
   - Creates a GitHub release
   - Uploads all built binaries as release assets
   - Generates automatic version numbers for non-tagged commits

**Platforms Supported:**
- `linux-amd64` - Standard Linux 64-bit
- `linux-arm64` - Linux ARM64 (e.g., Raspberry Pi, ARM servers)
- `darwin-amd64` - macOS Intel
- `darwin-arm64` - macOS Apple Silicon (M1, M2, etc.)
- `windows-amd64` - Windows 64-bit

## Version Management

The project uses automatic version generation:
- **Tagged releases:** Uses the git tag as version (e.g., `v1.0.0`)
- **Development builds:** Generates version like `v1.0.20240706123000-abc12345`
  - Includes timestamp and commit hash
  - Marked as pre-release

## Release Assets

Each release includes:
- Cross-platform binaries
- Documentation (`README.md`, `examples.md`, `TESTING.md`)
- Detailed release notes with usage instructions

## Security

- All workflows use pinned action versions for security
- No secrets are exposed in logs
- GitHub token permissions are limited to releases only

## Local Testing

To test builds locally before pushing:

```bash
# Test formatting
go fmt ./...

# Run tests
go test -v ./...

# Test cross-compilation
GOOS=linux GOARCH=amd64 go build -o shamir-cli-linux .
GOOS=darwin GOARCH=amd64 go build -o shamir-cli-darwin .
GOOS=windows GOARCH=amd64 go build -o shamir-cli-windows.exe .
```