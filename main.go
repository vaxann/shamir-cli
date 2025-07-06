package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"shamir-cli/shamir"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "shamir-cli",
	Short: "CLI application for secret sharing using Shamir's algorithm",
	Long: `Application for splitting a string into parts with the ability to recover
from fewer parts using Shamir's secret sharing algorithm.`,
}

var splitCmd = &cobra.Command{
	Use:   "split [string] [total_parts] [threshold]",
	Short: "Split a string into parts",
	Long: `Splits the input string into the specified number of parts, where a minimum
number of parts (threshold) is required for recovery.`,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		secret := args[0]
		n, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("Error: invalid number of parts '%s'\n", args[1])
			os.Exit(1)
		}
		
		k, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Printf("Error: invalid threshold '%s'\n", args[2])
			os.Exit(1)
		}
		
		if k < 2 {
			fmt.Println("Error: minimum number of parts for recovery must be at least 2")
			os.Exit(1)
		}
		
		if n < k {
			fmt.Println("Error: total number of parts cannot be less than threshold")
			os.Exit(1)
		}
		
		if n > 255 {
			fmt.Println("Error: total number of parts cannot be greater than 255")
			os.Exit(1)
		}
		
		shares, err := shamir.Split([]byte(secret), n, k)
		if err != nil {
			fmt.Printf("Error during splitting: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("Secret split into %d parts, %d parts required for recovery:\n\n", n, k)
		for i, share := range shares {
			fmt.Printf("Part %d: %s\n", i+1, shamir.ShareToString(share))
		}
		
		fmt.Printf("\nTo recover the secret use the command:\n")
		fmt.Printf("shamir-cli combine \"[parts_separated_by_commas]\"\n")
		fmt.Printf("Example: shamir-cli combine \"%s,%s\"\n", 
			shamir.ShareToString(shares[0]), shamir.ShareToString(shares[1]))
	},
}

var combineCmd = &cobra.Command{
	Use:   "combine [parts_separated_by_commas]",
	Short: "Recover a string from parts",
	Long: `Recovers the original string from parts separated by commas.
Each part must be in the format "ID:hex_value".`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		shareStrings := strings.Split(args[0], ",")
		if len(shareStrings) < 2 {
			fmt.Println("Error: minimum 2 parts required for recovery")
			os.Exit(1)
		}
		
		shares := make([]shamir.Share, 0, len(shareStrings))
		for i, shareStr := range shareStrings {
			shareStr = strings.TrimSpace(shareStr)
			if shareStr == "" {
				continue
			}
			
			share, err := shamir.StringToShare(shareStr)
			if err != nil {
				fmt.Printf("Error parsing part %d ('%s'): %v\n", i+1, shareStr, err)
				os.Exit(1)
			}
			shares = append(shares, share)
		}
		
		if len(shares) < 2 {
			fmt.Println("Error: minimum 2 valid parts required for recovery")
			os.Exit(1)
		}
		
		secret, err := shamir.Combine(shares)
		if err != nil {
			fmt.Printf("Error during recovery: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("Recovered secret: %s\n", string(secret))
	},
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run algorithm test",
	Long:  `Runs a simple test to verify the secret sharing algorithm functionality.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Testing Shamir's algorithm...")
		
		// Test 1: Simple string
		secret := "Hello, World!"
		fmt.Printf("Original secret: %s\n", secret)
		
		shares, err := shamir.Split([]byte(secret), 5, 3)
		if err != nil {
			fmt.Printf("Error during splitting: %v\n", err)
			return
		}
		
		fmt.Printf("Split into %d parts (threshold: 3):\n", len(shares))
		for i, share := range shares {
			fmt.Printf("  Part %d: %s\n", i+1, shamir.ShareToString(share))
		}
		
		// Test recovery with minimum number of parts
		testShares := shares[:3]
		recovered, err := shamir.Combine(testShares)
		if err != nil {
			fmt.Printf("Error during recovery: %v\n", err)
			return
		}
		
		fmt.Printf("Recovered from 3 parts: %s\n", string(recovered))
		
		if string(recovered) == secret {
			fmt.Println("✓ Test passed successfully!")
		} else {
			fmt.Println("✗ Test failed!")
		}
		
		// Test 2: Recovery with more parts
		testShares = shares[:4]
		recovered, err = shamir.Combine(testShares)
		if err != nil {
			fmt.Printf("Error during recovery: %v\n", err)
			return
		}
		
		fmt.Printf("Recovered from 4 parts: %s\n", string(recovered))
		
		if string(recovered) == secret {
			fmt.Println("✓ Test with 4 parts passed successfully!")
		} else {
			fmt.Println("✗ Test with 4 parts failed!")
		}
	},
}

func init() {
	rootCmd.AddCommand(splitCmd)
	rootCmd.AddCommand(combineCmd)
	rootCmd.AddCommand(testCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}