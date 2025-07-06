# Shamir CLI - Secret Sharing with Shamir's Algorithm

CLI application for splitting a string into parts with the ability to recover from fewer parts using Shamir's secret sharing algorithm.

## What is it?

Shamir's Secret Sharing is a cryptographic method that allows splitting a secret into n parts, where k parts are required for recovery (where k < n). This creates a threshold secret sharing scheme.

**Key principles:**
- The secret is divided into multiple parts
- A minimum number of parts (threshold) is needed for recovery
- Fewer parts than the threshold cannot recover the secret
- Each part looks like random data

## Installation

### Download Pre-built Binaries

Download the latest release for your platform from the [GitHub Releases page](../../releases).

Available platforms:
- Linux (AMD64, ARM64)
- macOS (Intel, Apple Silicon)
- Windows (AMD64)

### Build from Source

```bash
git clone <repository-url>
cd shamir-cli
go mod tidy
go build -o shamir-cli
```

## Usage

### Splitting a secret

```bash
./shamir-cli split "My secret" 5 3
```

This will split the string "My secret" into 5 parts, where a minimum of 3 parts will be required for recovery.

### Recovering a secret

```bash
./shamir-cli combine "1:d2b8c1a5,2:f4e3d2c1,3:a6b5c4d3"
```

Recovers the secret from the specified parts.



## Commands

- `split [string] [total_parts] [threshold]` - split a secret
- `combine [parts_separated_by_commas]` - recover a secret
- `help` - show help

## Examples

```bash
# Split into 7 parts, recover with 4
./shamir-cli split "Secret password" 7 4

# Recover from parts
./shamir-cli combine "1:a1b2c3,2:d4e5f6,3:g7h8i9,4:j1k2l3"

# Minimal scheme
./shamir-cli split "test" 3 2
```

## Practical Applications

1. **Secure password storage** - split passwords between multiple people
2. **Key backup** - protect cryptographic keys
3. **Corporate security** - require participation of multiple employees for access
4. **Family security** - access to important data only with participation of multiple family members

## Technical Features

- Maximum 255 parts
- Minimum 2 parts for recovery
- Uses arithmetic in finite field GF(2^8)
- Based on Lagrange polynomial interpolation
- Each run creates new random parts

## Security

- Parts look like random data
- Impossible to recover secret with fewer parts than threshold
- Cryptographically secure algorithm
- No additional keys or passwords required

## Development

The project uses GitHub Actions for continuous integration and automated releases:

- **Pull Requests**: Automatically run tests, formatting checks, and builds
- **Releases**: Automatic cross-platform builds on every commit to main branch
- **Tagged Releases**: Create proper releases by pushing git tags (e.g., `v1.0.0`)

### Testing

```bash
go test ./shamir -v          # Run tests
go test ./shamir -bench=.    # Run benchmarks
go test ./shamir -cover      # Run with coverage
```

## Documentation

- [`examples.md`](examples.md) - Detailed usage examples
- [`TESTING.md`](TESTING.md) - Testing information and results
- [`RELEASE.md`](RELEASE.md) - Release process documentation
- [`.github/README.md`](.github/README.md) - GitHub Actions documentation
