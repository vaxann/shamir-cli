# Shamir CLI - Secret Sharing with Shamir's Algorithm

CLI application for splitting secrets into parts with the ability to recover from fewer parts using Shamir's secret sharing algorithm.

## What is it?

Shamir's Secret Sharing is a cryptographic method that allows splitting a secret into n parts, where k parts are required for recovery (where k < n). This creates a threshold secret sharing scheme.

**Key principles:**
- The secret is divided into multiple parts (shares)
- A minimum number of parts (threshold) is needed for recovery
- Fewer parts than the threshold cannot recover the secret
- Each part looks like random data
- **Built-in checksum validation** ensures data integrity during recovery

## Features

- ✅ **Checksum validation** - Automatic detection of corrupted or invalid shares
- ✅ **Cross-platform** - Available for Linux, macOS, and Windows
- ✅ **Secure** - Uses arithmetic in finite field GF(2^8) with cryptographically secure random generation
- ✅ **Flexible** - Support for 2-255 parts with customizable threshold
- ✅ **CLI-friendly** - Simple command-line interface with clear output

## Installation

### Download Pre-built Binaries

Download the latest release for your platform from the [GitHub Releases page](../../releases).

Available platforms:
- **Linux**: AMD64, ARM64
- **macOS**: Intel (AMD64), Apple Silicon (ARM64)
- **Windows**: AMD64

### Build from Source

```bash
git clone <repository-url>
cd shamir-cli
go mod tidy
go build -o shamir-cli
```

**Requirements:**
- Go 1.21 or later
- No external dependencies (uses only standard library + cobra CLI)

## Usage

### Splitting a secret

```bash
./shamir-cli split "My secret password" 5 3
```

This will split the string "My secret password" into 5 parts, where a minimum of 3 parts will be required for recovery.

**Example output:**
```
Secret split into 5 parts, 3 parts required for recovery:

Part 1: 1:a1b2c3d4e5f6
Part 2: 2:f4e3d2c1b0a9  
Part 3: 3:a6b5c4d3e2f1
Part 4: 4:9f8e7d6c5b4a
Part 5: 5:3e4d5c6b7a89

To recover the secret use the command:
shamir-cli combine "[parts_separated_by_commas]"
Example: shamir-cli combine "1:a1b2c3d4e5f6,2:f4e3d2c1b0a9,3:a6b5c4d3e2f1"
```

### Recovering a secret

```bash
./shamir-cli combine "1:a1b2c3d4e5f6,2:f4e3d2c1b0a9,3:a6b5c4d3e2f1"
```

Recovers the secret from the specified parts with automatic checksum validation.

**Example output:**
```
Recovered secret: My secret password
```

If shares are corrupted or invalid, you'll see an error:
```
Error during recovery: checksum verification failed: unable to recover original string
```

## Commands

- `split [string] [total_parts] [threshold]` - Split a secret into parts
- `combine [parts_separated_by_commas]` - Recover a secret from parts
- `help` - Show help information
- `version` - Show version information

## Examples

```bash
# Split into 7 parts, recover with 4
./shamir-cli split "Secret password" 7 4

# Recover from parts (any 4 out of 7)
./shamir-cli combine "1:a1b2c3,2:d4e5f6,3:g7h8i9,4:j1k2l3"

# Minimal scheme (2 parts, both required)
./shamir-cli split "test" 2 2

# Maximum parts (255 parts, recover with 10)
./shamir-cli split "highly secure secret" 255 10
```

## Practical Applications

1. **Secure password storage** - Distribute passwords among multiple trusted parties
2. **Cryptographic key backup** - Protect important encryption keys
3. **Corporate security** - Require multiple employees for critical access
4. **Family security** - Secure access to important documents/accounts
5. **Disaster recovery** - Distribute backups across multiple locations

## Technical Details

### Security Features
- **Finite field arithmetic**: Uses GF(2^8) with irreducible polynomial x^8 + x^4 + x^3 + x + 1
- **Lagrange interpolation**: Recovers secrets using polynomial interpolation
- **Checksum validation**: XOR checksum prevents accepting corrupted shares
- **Cryptographic randomness**: Uses `crypto/rand` for secure coefficient generation
- **Information-theoretic security**: Shares reveal no information about the secret

### Limitations
- **Maximum 255 parts** (due to GF(2^8) field size)
- **Minimum 2 parts** required for recovery
- **Secret length**: No practical limit (each byte processed independently)
- **Performance**: Optimized with lookup tables for field operations

### Algorithm Details
1. Secret is padded with XOR checksum
2. Each byte creates a separate polynomial of degree k-1
3. Shares are evaluations of polynomials at distinct points
4. Recovery uses Lagrange interpolation to find polynomial constants
5. Checksum is validated to ensure data integrity

## Development

### Testing

```bash
go test ./shamir -v          # Run all tests
go test ./shamir -bench=.    # Run benchmarks
go test ./shamir -cover      # Run with coverage
```

### CI/CD

The project uses GitHub Actions for:
- **Pull Requests**: Automated testing, formatting checks, and cross-platform builds
- **Releases**: Automatic binary builds for all platforms on every push to main
- **Tagged Releases**: Create official releases by pushing git tags (e.g., `v1.0.0`)

### Project Structure

```
├── main.go              # CLI interface using Cobra
├── shamir/
│   ├── shamir.go        # Core Shamir's Secret Sharing implementation
│   └── shamir_test.go   # Comprehensive test suite
├── .github/workflows/   # GitHub Actions CI/CD
├── examples.md          # Detailed usage examples
├── TESTING.md           # Testing documentation
└── RELEASE.md           # Release process documentation
```

## Documentation

- [`examples.md`](examples.md) - Detailed usage examples and scenarios
- [`TESTING.md`](TESTING.md) - Testing information and performance results
- [`RELEASE.md`](RELEASE.md) - Release process and versioning
- [`CONTRIBUTING.md`](CONTRIBUTING.md) - How to contribute to the project

## Version

Current version includes checksum validation for enhanced security and data integrity verification.

Built with Go 1.21 and the Cobra CLI framework.
