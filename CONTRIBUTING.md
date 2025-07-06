# Contributing to Shamir CLI

Thank you for your interest in contributing to the Shamir CLI project!

## Development Setup

1. **Fork the repository** on GitHub
2. **Clone your fork:**
   ```bash
   git clone https://github.com/YOUR-USERNAME/shamir-cli.git
   cd shamir-cli
   ```
3. **Install dependencies:**
   ```bash
   go mod download
   ```

## Development Workflow

### Before Making Changes

1. **Create a new branch:**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Verify everything works:**
   ```bash
   go test ./...
   go fmt ./...
   go vet ./...
   go build .
   ```

### Making Changes

1. **Write tests** for new functionality
2. **Update documentation** if needed
3. **Follow Go conventions** and best practices
4. **Keep changes focused** and atomic

### Testing Your Changes

1. **Run all tests:**
   ```bash
   go test ./... -v
   go test ./... -race
   go test ./... -cover
   ```

2. **Run benchmarks:**
   ```bash
   go test ./shamir -bench=.
   ```

3. **Test cross-compilation:**
   ```bash
   GOOS=linux GOARCH=amd64 go build .
   GOOS=darwin GOARCH=amd64 go build .
   GOOS=windows GOARCH=amd64 go build .
   ```

4. **Test CLI functionality:**
   ```bash
   go build -o shamir-cli .
   ./shamir-cli split "test secret" 5 3
   ./shamir-cli combine "1:abc,2:def,3:xyz"
   ```

### Code Quality Checks

Run these commands before committing:

```bash
# Format code
go fmt ./...

# Static analysis
go vet ./...

# Check for common issues
golint ./... (if installed)

# Security check
gosec ./... (if installed)
```

## GitHub Actions

The project uses automated CI/CD:

### Pull Request Workflow
- Runs on every PR
- Executes all tests
- Checks formatting
- Builds the application
- Reports test coverage

### Release Workflow
- Runs on push to main/master
- Builds for all platforms
- Creates releases automatically
- Uploads binary assets

### Testing GitHub Actions Locally

You can test GitHub Actions locally using [act](https://github.com/nektos/act):

```bash
# Install act
curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash

# Test the test workflow
act pull_request

# Test the release workflow (dry run)
act push --dry-run
```

## Pull Request Guidelines

1. **Fill out the PR template** (if available)
2. **Link related issues** in the description
3. **Ensure all CI checks pass**
4. **Keep the PR focused** on a single feature/fix
5. **Write clear commit messages**

### Commit Message Format

Use conventional commits:
```
type(scope): description

Examples:
feat(cli): add support for custom separators
fix(shamir): handle edge case in polynomial evaluation
docs(readme): update installation instructions
test(shamir): add tests for large secrets
```

## Code Style

- Follow standard Go formatting (`go fmt`)
- Use meaningful variable and function names
- Add comments for complex logic
- Keep functions small and focused
- Write table-driven tests when appropriate

## Documentation

Update documentation when making changes:

- `README.md` - Main project documentation
- `examples.md` - Usage examples
- `TESTING.md` - Testing information
- Code comments for complex algorithms

## Release Process

Maintainers handle releases:

1. **Development releases:** Automatic on every commit to main
2. **Tagged releases:** Manual with semantic versioning
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

## Getting Help

- **Issues:** Report bugs or request features via GitHub Issues
- **Discussions:** Ask questions in GitHub Discussions
- **Code Review:** Maintainers will review PRs and provide feedback

## Security

If you discover a security vulnerability:
- **Do NOT** open a public issue
- **Email** the maintainers privately
- **Provide** detailed information about the vulnerability

## License

By contributing, you agree that your contributions will be licensed under the same license as the project.

Thank you for contributing! 🎉