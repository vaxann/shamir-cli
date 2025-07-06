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
	Short: "CLI приложение для разделения секретов по алгоритму Шамира",
	Long: `Приложение для разделения строки на части с возможностью восстановления 
по меньшему количеству частей, используя алгоритм Шамира.`,
}

var splitCmd = &cobra.Command{
	Use:   "split [строка] [общее количество частей] [минимальное количество для восстановления]",
	Short: "Разделить строку на части",
	Long: `Разделяет входную строку на указанное количество частей, где для восстановления 
требуется минимальное количество частей (порог).`,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		secret := args[0]
		n, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("Ошибка: неверное количество частей '%s'\n", args[1])
			os.Exit(1)
		}
		
		k, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Printf("Ошибка: неверный порог '%s'\n", args[2])
			os.Exit(1)
		}
		
		if k < 2 {
			fmt.Println("Ошибка: минимальное количество частей для восстановления должно быть не менее 2")
			os.Exit(1)
		}
		
		if n < k {
			fmt.Println("Ошибка: общее количество частей не может быть меньше минимального")
			os.Exit(1)
		}
		
		if n > 255 {
			fmt.Println("Ошибка: общее количество частей не может быть больше 255")
			os.Exit(1)
		}
		
		shares, err := shamir.Split([]byte(secret), n, k)
		if err != nil {
			fmt.Printf("Ошибка при разделении: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("Секрет разделен на %d частей, для восстановления требуется %d частей:\n\n", n, k)
		for i, share := range shares {
			fmt.Printf("Часть %d: %s\n", i+1, shamir.ShareToString(share))
		}
		
		fmt.Printf("\nДля восстановления секрета используйте команду:\n")
		fmt.Printf("shamir-cli combine \"[части через запятую]\"\n")
		fmt.Printf("Например: shamir-cli combine \"%s,%s\"\n", 
			shamir.ShareToString(shares[0]), shamir.ShareToString(shares[1]))
	},
}

var combineCmd = &cobra.Command{
	Use:   "combine [части через запятую]",
	Short: "Восстановить строку из частей",
	Long: `Восстанавливает оригинальную строку из частей, разделенных запятыми.
Каждая часть должна быть в формате "ID:hex_value".`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		shareStrings := strings.Split(args[0], ",")
		if len(shareStrings) < 2 {
			fmt.Println("Ошибка: необходимо минимум 2 части для восстановления")
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
				fmt.Printf("Ошибка при разборе части %d ('%s'): %v\n", i+1, shareStr, err)
				os.Exit(1)
			}
			shares = append(shares, share)
		}
		
		if len(shares) < 2 {
			fmt.Println("Ошибка: необходимо минимум 2 корректные части для восстановления")
			os.Exit(1)
		}
		
		secret, err := shamir.Combine(shares)
		if err != nil {
			fmt.Printf("Ошибка при восстановлении: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("Восстановленный секрет: %s\n", string(secret))
	},
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Запустить тест алгоритма",
	Long:  `Запускает простой тест для проверки работы алгоритма разделения секретов.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Тестирование алгоритма Шамира...")
		
		// Тест 1: Простая строка
		secret := "Привет, мир!"
		fmt.Printf("Исходный секрет: %s\n", secret)
		
		shares, err := shamir.Split([]byte(secret), 5, 3)
		if err != nil {
			fmt.Printf("Ошибка при разделении: %v\n", err)
			return
		}
		
		fmt.Printf("Разделено на %d частей (порог: 3):\n", len(shares))
		for i, share := range shares {
			fmt.Printf("  Часть %d: %s\n", i+1, shamir.ShareToString(share))
		}
		
		// Тест восстановления с минимальным количеством частей
		testShares := shares[:3]
		recovered, err := shamir.Combine(testShares)
		if err != nil {
			fmt.Printf("Ошибка при восстановлении: %v\n", err)
			return
		}
		
		fmt.Printf("Восстановлено из 3 частей: %s\n", string(recovered))
		
		if string(recovered) == secret {
			fmt.Println("✓ Тест пройден успешно!")
		} else {
			fmt.Println("✗ Тест провален!")
		}
		
		// Тест 2: Восстановление с большим количеством частей
		testShares = shares[:4]
		recovered, err = shamir.Combine(testShares)
		if err != nil {
			fmt.Printf("Ошибка при восстановлении: %v\n", err)
			return
		}
		
		fmt.Printf("Восстановлено из 4 частей: %s\n", string(recovered))
		
		if string(recovered) == secret {
			fmt.Println("✓ Тест с 4 частями пройден успешно!")
		} else {
			fmt.Println("✗ Тест с 4 частями провален!")
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