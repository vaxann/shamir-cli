# Testing Shamir CLI

This document describes the testing strategy and results for the Shamir secret sharing implementation.

## Test Coverage

The algorithm is thoroughly tested with comprehensive Go unit tests located in `shamir/shamir_test.go`.

### Test Categories

1. **Galois Field Operations**
   - Tests basic arithmetic operations in GF(2^8)
   - Verifies inverse element calculations
   - Ensures mathematical correctness

2. **Core Algorithm Tests**
   - Basic split and combine functionality
   - Recovery with different subsets of shares
   - Recovery with more than threshold shares
   - Various combinations of n and k parameters

3. **Validation Tests**
   - Input parameter validation for Split function
   - Error handling for invalid inputs
   - Edge cases and boundary conditions

4. **String Conversion Tests**
   - Share to string conversion
   - String to share parsing
   - Error handling for malformed inputs

5. **Edge Cases**
   - Empty secrets
   - Large secrets (1KB)
   - Minimum viable parameters
   - Randomness verification

6. **Performance Tests**
   - Benchmarks for split operations
   - Benchmarks for combine operations

## Running Tests

```bash
# Run all tests with verbose output
go test ./shamir -v

# Run specific test
go test ./shamir -run TestBasicSplitAndCombine

# Run benchmarks
go test ./shamir -bench=.

# Run tests with coverage
go test ./shamir -cover
```

## Test Results

All tests pass successfully, demonstrating:

- ✅ Correct implementation of Shamir's secret sharing
- ✅ Proper error handling and validation
- ✅ Robustness across various input sizes and parameters
- ✅ Strong performance characteristics

## Performance Benchmarks

Recent benchmark results on AMD EPYC 7R13:

```
BenchmarkSplit-4     42978    83351 ns/op   (~83μs per split)
BenchmarkCombine-4   1489624  2415 ns/op    (~2.4μs per combine)
```

The algorithm shows excellent performance:
- Split operations: ~83 microseconds for a 39-byte secret into 10 parts
- Combine operations: ~2.4 microseconds for recovery from 5 parts

## Security Validation

The tests verify:
- Shares appear random and provide no information about the secret
- Recovery requires exactly the threshold number of shares
- Different splits of the same secret produce different shares
- Mathematical properties of the Galois field operations

## Code Quality

- All code is formatted with `go fmt`
- Tests follow Go testing best practices
- Comprehensive error checking and validation
- Clear test names and documentation