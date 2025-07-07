package shamir

import (
	"crypto/rand"
	"errors"
	"fmt"
)

// Share represents one part of the secret
type Share struct {
	ID    byte   `json:"id"`
	Value []byte `json:"value"`
}

// Lookup tables for arithmetic in GF(2^8)
var gfMulTable [256][256]byte
var gfInvTable [256]byte

func init() {
	initGF()
}

// initGF initializes tables for arithmetic in GF(2^8)
func initGF() {
	// Initialize multiplication table
	for a := 0; a < 256; a++ {
		for b := 0; b < 256; b++ {
			gfMulTable[a][b] = gfMulPrimitive(byte(a), byte(b))
		}
	}

	// Initialize inverse elements table
	gfInvTable[0] = 0
	for i := 1; i < 256; i++ {
		gfInvTable[i] = gfInvPrimitive(byte(i))
	}
}

// gfMulPrimitive performs multiplication in GF(2^8) without using tables
func gfMulPrimitive(a, b byte) byte {
	if a == 0 || b == 0 {
		return 0
	}

	var result byte
	for i := 0; i < 8; i++ {
		if (b & 1) == 1 {
			result ^= a
		}
		highBit := (a & 0x80) != 0
		a <<= 1
		if highBit {
			a ^= 0x1B // irreducible polynomial x^8 + x^4 + x^3 + x + 1
		}
		b >>= 1
	}
	return result
}

// gfInvPrimitive calculates the inverse element in GF(2^8) using brute force
func gfInvPrimitive(a byte) byte {
	if a == 0 {
		return 0
	}

	// Try all possible values
	for i := 1; i < 256; i++ {
		if gfMulPrimitive(a, byte(i)) == 1 {
			return byte(i)
		}
	}
	return 0
}

// gfMul performs multiplication in GF(2^8) using tables
func gfMul(a, b byte) byte {
	return gfMulTable[a][b]
}

// gfInv calculates the inverse element in GF(2^8) using tables
func gfInv(a byte) byte {
	return gfInvTable[a]
}

// gfAdd performs addition in GF(2^8) (XOR)
func gfAdd(a, b byte) byte {
	return a ^ b
}

// gfSub performs subtraction in GF(2^8) (XOR)
func gfSub(a, b byte) byte {
	return a ^ b
}

// evaluatePolynomial calculates the value of a polynomial at point x
func evaluatePolynomial(coeffs []byte, x byte) byte {
	if len(coeffs) == 0 {
		return 0
	}

	result := coeffs[0]
	xPow := byte(1)

	for i := 1; i < len(coeffs); i++ {
		xPow = gfMul(xPow, x)
		result = gfAdd(result, gfMul(coeffs[i], xPow))
	}

	return result
}

// calculateChecksum calculates XOR checksum of all bytes
func calculateChecksum(data []byte) byte {
	var checksum byte
	for _, b := range data {
		checksum ^= b
	}
	return checksum
}

// Split divides a secret into n parts, where k parts are needed for recovery
func Split(secret []byte, n, k int) ([]Share, error) {
	if k < 2 {
		return nil, errors.New("k must be at least 2")
	}
	if n < k {
		return nil, errors.New("n must be at least k")
	}
	if n > 255 {
		return nil, errors.New("n cannot be greater than 255")
	}

	// Add checksum to the secret
	checksum := calculateChecksum(secret)
	secretWithChecksum := append(secret, checksum)

	shares := make([]Share, n)

	// For each byte of the secret (including checksum), create a separate polynomial
	for byteIndex := 0; byteIndex < len(secretWithChecksum); byteIndex++ {
		// Create random coefficients for polynomial of degree k-1
		coeffs := make([]byte, k)
		coeffs[0] = secretWithChecksum[byteIndex] // constant term is the secret byte

		// Generate random coefficients for other degrees
		for i := 1; i < k; i++ {
			randomBytes := make([]byte, 1)
			rand.Read(randomBytes)
			coeffs[i] = randomBytes[0]
		}

		// Calculate polynomial values for each part
		for i := 0; i < n; i++ {
			shareID := byte(i + 1) // Share ID (from 1 to n)
			shareValue := evaluatePolynomial(coeffs, shareID)

			if byteIndex == 0 {
				shares[i] = Share{
					ID:    shareID,
					Value: make([]byte, len(secretWithChecksum)),
				}
			}
			shares[i].Value[byteIndex] = shareValue
		}
	}

	return shares, nil
}

// Combine recovers a secret from parts
func Combine(shares []Share) ([]byte, error) {
	if len(shares) < 2 {
		return nil, errors.New("minimum 2 parts required")
	}

	// Check that all parts have the same length
	secretLen := len(shares[0].Value)
	for i := 1; i < len(shares); i++ {
		if len(shares[i].Value) != secretLen {
			return nil, errors.New("all parts must have the same length")
		}
	}

	secretWithChecksum := make([]byte, secretLen)

	// Recover each byte of the secret separately
	for byteIndex := 0; byteIndex < secretLen; byteIndex++ {
		// Collect points for interpolation
		xs := make([]byte, len(shares))
		ys := make([]byte, len(shares))

		for i, share := range shares {
			xs[i] = share.ID
			ys[i] = share.Value[byteIndex]
		}

		// Use Lagrange interpolation to recover the constant term
		secretWithChecksum[byteIndex] = lagrangeInterpolation(xs, ys)
	}

	// Verify checksum
	if len(secretWithChecksum) < 1 {
		return nil, errors.New("recovered data is too short")
	}

	secret := secretWithChecksum[:len(secretWithChecksum)-1]
	expectedChecksum := secretWithChecksum[len(secretWithChecksum)-1]
	actualChecksum := calculateChecksum(secret)

	if expectedChecksum != actualChecksum {
		return nil, errors.New("checksum verification failed: unable to recover original string")
	}

	return secret, nil
}

// lagrangeInterpolation recovers the constant term of the polynomial (value at point 0)
func lagrangeInterpolation(xs, ys []byte) byte {
	var result byte

	for i := 0; i < len(xs); i++ {
		var numerator, denominator byte = 1, 1

		for j := 0; j < len(xs); j++ {
			if i != j {
				numerator = gfMul(numerator, xs[j])
				denominator = gfMul(denominator, gfAdd(xs[i], xs[j]))
			}
		}

		if denominator != 0 {
			lagrangeBasis := gfMul(numerator, gfInv(denominator))
			result = gfAdd(result, gfMul(ys[i], lagrangeBasis))
		}
	}

	return result
}

// ShareToString converts a Share to string representation
func ShareToString(share Share) string {
	return fmt.Sprintf("%d:%x", share.ID, share.Value)
}

// StringToShare converts string representation to Share
func StringToShare(s string) (Share, error) {
	var share Share
	var hexValue string

	n, err := fmt.Sscanf(s, "%d:%s", &share.ID, &hexValue)
	if err != nil || n != 2 {
		return Share{}, errors.New("invalid part format")
	}

	// Check if hex string has even length
	if len(hexValue)%2 != 0 {
		return Share{}, errors.New("invalid hex format")
	}

	value := make([]byte, len(hexValue)/2)
	for i := 0; i < len(hexValue); i += 2 {
		var b byte
		n, err := fmt.Sscanf(hexValue[i:i+2], "%02x", &b)
		if err != nil || n != 1 {
			return Share{}, errors.New("invalid hex format")
		}
		value[i/2] = b
	}

	share.Value = value
	return share, nil
}
