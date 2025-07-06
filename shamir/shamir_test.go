package shamir

import (
	"bytes"
	"fmt"
	"testing"
)

func TestGaloisFieldOperations(t *testing.T) {
	// Test basic GF operations
	tests := []struct {
		name     string
		a, b     byte
		expected byte
		op       string
	}{
		{"Add 0+0", 0, 0, 0, "add"},
		{"Add 1+1", 1, 1, 0, "add"},
		{"Add 5+3", 5, 3, 6, "add"},
		{"Mul 0*5", 0, 5, 0, "mul"},
		{"Mul 1*5", 1, 5, 5, "mul"},
		{"Mul 2*3", 2, 3, 6, "mul"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result byte
			switch tt.op {
			case "add":
				result = gfAdd(tt.a, tt.b)
			case "mul":
				result = gfMul(tt.a, tt.b)
			}
			if result != tt.expected {
				t.Errorf("%s: got %d, want %d", tt.name, result, tt.expected)
			}
		})
	}
}

func TestGaloisFieldInverse(t *testing.T) {
	// Test that a * inv(a) = 1 for non-zero elements
	for a := 1; a < 256; a++ {
		inv := gfInv(byte(a))
		if inv == 0 {
			t.Errorf("Inverse of %d should not be 0", a)
			continue
		}

		product := gfMul(byte(a), inv)
		if product != 1 {
			t.Errorf("a * inv(a) != 1: %d * %d = %d", a, inv, product)
		}
	}
}

func TestBasicSplitAndCombine(t *testing.T) {
	secret := []byte("Hello, World!")
	n, k := 5, 3

	// Split the secret
	shares, err := Split(secret, n, k)
	if err != nil {
		t.Fatalf("Split failed: %v", err)
	}

	if len(shares) != n {
		t.Fatalf("Expected %d shares, got %d", n, len(shares))
	}

	// Verify each share has correct structure
	for i, share := range shares {
		if share.ID != byte(i+1) {
			t.Errorf("Share %d has wrong ID: got %d, want %d", i, share.ID, i+1)
		}
		if len(share.Value) != len(secret) {
			t.Errorf("Share %d has wrong length: got %d, want %d", i, len(share.Value), len(secret))
		}
	}

	// Test recovery with minimum threshold
	testShares := shares[:k]
	recovered, err := Combine(testShares)
	if err != nil {
		t.Fatalf("Combine failed: %v", err)
	}

	if !bytes.Equal(recovered, secret) {
		t.Errorf("Recovery failed: got %q, want %q", string(recovered), string(secret))
	}
}

func TestRecoveryWithDifferentSubsets(t *testing.T) {
	secret := []byte("Test secret 123")
	n, k := 7, 4

	shares, err := Split(secret, n, k)
	if err != nil {
		t.Fatalf("Split failed: %v", err)
	}

	// Test different combinations of k shares
	testCombinations := [][]int{
		{0, 1, 2, 3},
		{0, 2, 4, 6},
		{1, 3, 5, 6},
		{0, 1, 4, 5},
		{2, 3, 4, 5},
	}

	for i, combination := range testCombinations {
		t.Run(fmt.Sprintf("Combination_%d", i), func(t *testing.T) {
			testShares := make([]Share, len(combination))
			for j, idx := range combination {
				testShares[j] = shares[idx]
			}

			recovered, err := Combine(testShares)
			if err != nil {
				t.Fatalf("Combine failed for combination %v: %v", combination, err)
			}

			if !bytes.Equal(recovered, secret) {
				t.Errorf("Recovery failed for combination %v: got %q, want %q",
					combination, string(recovered), string(secret))
			}
		})
	}
}

func TestRecoveryWithMoreThanThreshold(t *testing.T) {
	secret := []byte("Extended test")
	n, k := 6, 3

	shares, err := Split(secret, n, k)
	if err != nil {
		t.Fatalf("Split failed: %v", err)
	}

	// Test recovery with more than threshold shares
	for numShares := k; numShares <= n; numShares++ {
		t.Run(fmt.Sprintf("With_%d_shares", numShares), func(t *testing.T) {
			testShares := shares[:numShares]
			recovered, err := Combine(testShares)
			if err != nil {
				t.Fatalf("Combine failed with %d shares: %v", numShares, err)
			}

			if !bytes.Equal(recovered, secret) {
				t.Errorf("Recovery failed with %d shares: got %q, want %q",
					numShares, string(recovered), string(secret))
			}
		})
	}
}

func TestSplitValidation(t *testing.T) {
	secret := []byte("test")

	tests := []struct {
		name    string
		n, k    int
		wantErr bool
	}{
		{"Valid parameters", 5, 3, false},
		{"k too small", 5, 1, true},
		{"n less than k", 3, 5, true},
		{"n too large", 256, 2, true},
		{"Minimum valid", 2, 2, false},
		{"k equals n", 3, 3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Split(secret, tt.n, tt.k)
			if (err != nil) != tt.wantErr {
				t.Errorf("Split() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCombineValidation(t *testing.T) {
	// Test with insufficient shares
	shares := []Share{
		{ID: 1, Value: []byte{0x12}},
	}
	_, err := Combine(shares)
	if err == nil {
		t.Error("Combine should fail with only 1 share")
	}

	// Test with mismatched share lengths
	shares = []Share{
		{ID: 1, Value: []byte{0x12, 0x34}},
		{ID: 2, Value: []byte{0x56}},
	}
	_, err = Combine(shares)
	if err == nil {
		t.Error("Combine should fail with mismatched share lengths")
	}
}

func TestStringConversion(t *testing.T) {
	share := Share{
		ID:    1,
		Value: []byte{0x12, 0x34, 0xab, 0xcd},
	}

	// Test ShareToString
	str := ShareToString(share)
	expected := "1:1234abcd"
	if str != expected {
		t.Errorf("ShareToString() = %q, want %q", str, expected)
	}

	// Test StringToShare
	recovered, err := StringToShare(str)
	if err != nil {
		t.Fatalf("StringToShare() failed: %v", err)
	}

	if recovered.ID != share.ID {
		t.Errorf("Recovered ID = %d, want %d", recovered.ID, share.ID)
	}

	if !bytes.Equal(recovered.Value, share.Value) {
		t.Errorf("Recovered Value = %x, want %x", recovered.Value, share.Value)
	}

	// Test additional valid cases
	validCases := []struct {
		input    string
		expected Share
	}{
		{"1:ab", Share{ID: 1, Value: []byte{0xab}}},
		{"255:1234", Share{ID: 255, Value: []byte{0x12, 0x34}}},
		{"10:00ff", Share{ID: 10, Value: []byte{0x00, 0xff}}},
	}

	for _, testCase := range validCases {
		t.Run(testCase.input, func(t *testing.T) {
			result, err := StringToShare(testCase.input)
			if err != nil {
				t.Fatalf("StringToShare(%q) failed: %v", testCase.input, err)
			}

			if result.ID != testCase.expected.ID {
				t.Errorf("ID = %d, want %d", result.ID, testCase.expected.ID)
			}

			if !bytes.Equal(result.Value, testCase.expected.Value) {
				t.Errorf("Value = %x, want %x", result.Value, testCase.expected.Value)
			}
		})
	}
}

func TestStringConversionErrors(t *testing.T) {
	tests := []string{
		"invalid",
		"1-abcd",
		"256:abcd",
		"1:xyz",
		"1:abc", // odd length hex string
		"1:",    // empty hex string
	}

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			_, err := StringToShare(test)
			if err == nil {
				t.Errorf("StringToShare(%q) should fail", test)
			}
		})
	}
}

func TestEmptySecret(t *testing.T) {
	secret := []byte("")
	shares, err := Split(secret, 3, 2)
	if err != nil {
		t.Fatalf("Split failed with empty secret: %v", err)
	}

	recovered, err := Combine(shares[:2])
	if err != nil {
		t.Fatalf("Combine failed with empty secret: %v", err)
	}

	if !bytes.Equal(recovered, secret) {
		t.Errorf("Recovery failed for empty secret: got %q, want %q", string(recovered), string(secret))
	}
}

func TestLargeSecret(t *testing.T) {
	// Test with a larger secret
	secret := make([]byte, 1000)
	for i := range secret {
		secret[i] = byte(i % 256)
	}

	shares, err := Split(secret, 10, 5)
	if err != nil {
		t.Fatalf("Split failed with large secret: %v", err)
	}

	recovered, err := Combine(shares[:5])
	if err != nil {
		t.Fatalf("Combine failed with large secret: %v", err)
	}

	if !bytes.Equal(recovered, secret) {
		t.Error("Recovery failed for large secret")
	}
}

func TestRandomnessOfShares(t *testing.T) {
	secret := []byte("same secret")

	// Generate shares twice and verify they're different
	shares1, err := Split(secret, 5, 3)
	if err != nil {
		t.Fatalf("First split failed: %v", err)
	}

	shares2, err := Split(secret, 5, 3)
	if err != nil {
		t.Fatalf("Second split failed: %v", err)
	}

	// Shares should be different (except for very unlikely case)
	different := false
	for i := 0; i < len(shares1); i++ {
		if !bytes.Equal(shares1[i].Value, shares2[i].Value) {
			different = true
			break
		}
	}

	if !different {
		t.Error("Shares from two splits should be different due to randomness")
	}

	// But both should recover the same secret
	recovered1, err := Combine(shares1[:3])
	if err != nil {
		t.Fatalf("First combine failed: %v", err)
	}

	recovered2, err := Combine(shares2[:3])
	if err != nil {
		t.Fatalf("Second combine failed: %v", err)
	}

	if !bytes.Equal(recovered1, secret) || !bytes.Equal(recovered2, secret) {
		t.Error("Both splits should recover the original secret")
	}
}

func BenchmarkSplit(b *testing.B) {
	secret := []byte("benchmark secret for testing performance")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Split(secret, 10, 5)
		if err != nil {
			b.Fatalf("Split failed: %v", err)
		}
	}
}

func BenchmarkCombine(b *testing.B) {
	secret := []byte("benchmark secret for testing performance")
	shares, err := Split(secret, 10, 5)
	if err != nil {
		b.Fatalf("Split failed: %v", err)
	}

	testShares := shares[:5]

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Combine(testShares)
		if err != nil {
			b.Fatalf("Combine failed: %v", err)
		}
	}
}
