# Release Process

This document describes how to create releases for the Shamir CLI project.

## Automated Releases

The project uses GitHub Actions for automated building and releasing.

### Development Releases (Automatic)

Every push to the `main` or `master` branch automatically:
1. Runs all tests
2. Builds binaries for all supported platforms
3. Creates a pre-release with auto-generated version
4. Uploads binaries as release assets

**Auto-generated version format:** `v1.0.YYYYMMDDHHMMSS-<commit-hash>`

Example: `v1.0.20240706123456-abc12345`

### Tagged Releases (Manual)

To create a proper versioned release:

1. **Create and push a tag:**
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. **GitHub Actions will automatically:**
   - Run all tests
   - Build for all platforms
   - Create a release with the tag version
   - Upload binaries

## Supported Platforms

The release workflow builds for:

| Platform | Architecture | Binary Name |
|----------|-------------|-------------|
| Linux | AMD64 | `shamir-cli-*-linux-amd64.tar.gz` |
| Linux | ARM64 | `shamir-cli-*-linux-arm64.tar.gz` |
| macOS | Intel | `shamir-cli-*-darwin-amd64.tar.gz` |
| macOS | Apple Silicon | `shamir-cli-*-darwin-arm64.tar.gz` |
| Windows | AMD64 | `shamir-cli-*-windows-amd64.zip` |

## Release Contents

Each release archive includes:
- Compiled binary (`shamir-cli` or `shamir-cli.exe`)
- `README.md` - Main documentation
- `examples.md` - Usage examples
- `TESTING.md` - Testing information

## Version Embedding

The version is embedded into the binary during build using Go's ldflags:
```bash
go build -ldflags="-X main.version=v1.0.0" .
```

Users can check the version with:
```bash
./shamir-cli --version
```

## Manual Release Process

If you need to create a release manually:

1. **Test everything:**
   ```bash
   go test ./...
   go fmt ./...
   go vet ./...
   ```

2. **Build for all platforms:**
   ```bash
   # Linux AMD64
   GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.version=v1.0.0" -o shamir-cli-linux-amd64 .
   
   # Linux ARM64
   GOOS=linux GOARCH=arm64 go build -ldflags="-s -w -X main.version=v1.0.0" -o shamir-cli-linux-arm64 .
   
   # macOS Intel
   GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X main.version=v1.0.0" -o shamir-cli-darwin-amd64 .
   
   # macOS Apple Silicon
   GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -X main.version=v1.0.0" -o shamir-cli-darwin-arm64 .
   
   # Windows
   GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X main.version=v1.0.0" -o shamir-cli-windows-amd64.exe .
   ```

3. **Create archives:**
   ```bash
   # Unix platforms (tar.gz)
   tar -czf shamir-cli-v1.0.0-linux-amd64.tar.gz shamir-cli-linux-amd64 README.md examples.md TESTING.md
   tar -czf shamir-cli-v1.0.0-linux-arm64.tar.gz shamir-cli-linux-arm64 README.md examples.md TESTING.md
   tar -czf shamir-cli-v1.0.0-darwin-amd64.tar.gz shamir-cli-darwin-amd64 README.md examples.md TESTING.md
   tar -czf shamir-cli-v1.0.0-darwin-arm64.tar.gz shamir-cli-darwin-arm64 README.md examples.md TESTING.md
   
   # Windows (zip)
   zip shamir-cli-v1.0.0-windows-amd64.zip shamir-cli-windows-amd64.exe README.md examples.md TESTING.md
   ```

4. **Create GitHub release manually with these assets**

## Versioning Strategy

The project follows [Semantic Versioning](https://semver.org/):

- `MAJOR.MINOR.PATCH` (e.g., `v1.0.0`)
- `MAJOR`: Breaking changes
- `MINOR`: New features (backward compatible)
- `PATCH`: Bug fixes (backward compatible)

## Pre-release Testing

Before creating a tagged release:

1. Test on multiple platforms if possible
2. Verify cross-compilation works
3. Check that all tests pass
4. Ensure documentation is up to date
5. Test the CLI functionality manually

## Release Checklist

- [ ] All tests pass (`go test ./...`)
- [ ] Code is formatted (`go fmt ./...`)
- [ ] No lint issues (`go vet ./...`)
- [ ] Documentation updated
- [ ] Version number decided
- [ ] Tag created and pushed
- [ ] GitHub Actions completed successfully
- [ ] Release assets uploaded
- [ ] Release notes written