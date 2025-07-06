package shamir

import (
	"crypto/rand"
	"errors"
	"fmt"
)

// Share представляет одну часть секрета
type Share struct {
	ID    byte   `json:"id"`
	Value []byte `json:"value"`
}

// Готовые таблицы для арифметики в GF(2^8)
var gfMulTable [256][256]byte
var gfInvTable [256]byte

func init() {
	initGF()
}

// initGF инициализирует таблицы для арифметики в GF(2^8)
func initGF() {
	// Инициализация таблицы умножения
	for a := 0; a < 256; a++ {
		for b := 0; b < 256; b++ {
			gfMulTable[a][b] = gfMulPrimitive(byte(a), byte(b))
		}
	}

	// Инициализация таблицы обратных элементов
	gfInvTable[0] = 0
	for i := 1; i < 256; i++ {
		gfInvTable[i] = gfInvPrimitive(byte(i))
	}
}

// gfMulPrimitive выполняет умножение в GF(2^8) без использования таблиц
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
			a ^= 0x1B // неприводимый полином x^8 + x^4 + x^3 + x + 1
		}
		b >>= 1
	}
	return result
}

// gfInvPrimitive вычисляет обратный элемент в GF(2^8) методом перебора
func gfInvPrimitive(a byte) byte {
	if a == 0 {
		return 0
	}
	
	// Перебираем все возможные значения
	for i := 1; i < 256; i++ {
		if gfMulPrimitive(a, byte(i)) == 1 {
			return byte(i)
		}
	}
	return 0
}

// gfMul выполняет умножение в GF(2^8) используя таблицы
func gfMul(a, b byte) byte {
	return gfMulTable[a][b]
}

// gfInv вычисляет обратный элемент в GF(2^8) используя таблицы
func gfInv(a byte) byte {
	return gfInvTable[a]
}

// gfAdd выполняет сложение в GF(2^8) (XOR)
func gfAdd(a, b byte) byte {
	return a ^ b
}

// gfSub выполняет вычитание в GF(2^8) (XOR)
func gfSub(a, b byte) byte {
	return a ^ b
}

// evaluatePolynomial вычисляет значение полинома в точке x
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

// Split разделяет секрет на n частей, где k частей необходимо для восстановления
func Split(secret []byte, n, k int) ([]Share, error) {
	if k < 2 {
		return nil, errors.New("k должно быть не менее 2")
	}
	if n < k {
		return nil, errors.New("n должно быть не менее k")
	}
	if n > 255 {
		return nil, errors.New("n не может быть больше 255")
	}
	
	shares := make([]Share, n)
	
	// Для каждого байта секрета создаем отдельный полином
	for byteIndex := 0; byteIndex < len(secret); byteIndex++ {
		// Создаем случайные коэффициенты для полинома степени k-1
		coeffs := make([]byte, k)
		coeffs[0] = secret[byteIndex] // свободный член - это байт секрета
		
		// Генерируем случайные коэффициенты для остальных степеней
		for i := 1; i < k; i++ {
			randomBytes := make([]byte, 1)
			rand.Read(randomBytes)
			coeffs[i] = randomBytes[0]
		}
		
		// Вычисляем значения полинома для каждой части
		for i := 0; i < n; i++ {
			shareID := byte(i + 1) // ID части (от 1 до n)
			shareValue := evaluatePolynomial(coeffs, shareID)
			
			if byteIndex == 0 {
				shares[i] = Share{
					ID:    shareID,
					Value: make([]byte, len(secret)),
				}
			}
			shares[i].Value[byteIndex] = shareValue
		}
	}
	
	return shares, nil
}

// Combine восстанавливает секрет из частей
func Combine(shares []Share) ([]byte, error) {
	if len(shares) < 2 {
		return nil, errors.New("необходимо минимум 2 части")
	}
	
	// Проверяем, что все части имеют одинаковую длину
	secretLen := len(shares[0].Value)
	for i := 1; i < len(shares); i++ {
		if len(shares[i].Value) != secretLen {
			return nil, errors.New("все части должны иметь одинаковую длину")
		}
	}
	
	secret := make([]byte, secretLen)
	
	// Восстанавливаем каждый байт секрета отдельно
	for byteIndex := 0; byteIndex < secretLen; byteIndex++ {
		// Собираем точки для интерполяции
		xs := make([]byte, len(shares))
		ys := make([]byte, len(shares))
		
		for i, share := range shares {
			xs[i] = share.ID
			ys[i] = share.Value[byteIndex]
		}
		
		// Используем интерполяцию Лагранжа для восстановления свободного члена
		secret[byteIndex] = lagrangeInterpolation(xs, ys)
	}
	
	return secret, nil
}

// lagrangeInterpolation восстанавливает свободный член полинома (значение в точке 0)
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

// ShareToString преобразует Share в строковое представление
func ShareToString(share Share) string {
	return fmt.Sprintf("%d:%x", share.ID, share.Value)
}

// StringToShare преобразует строковое представление в Share
func StringToShare(s string) (Share, error) {
	var share Share
	var hexValue string
	
	n, err := fmt.Sscanf(s, "%d:%s", &share.ID, &hexValue)
	if err != nil || n != 2 {
		return Share{}, errors.New("неверный формат части")
	}
	
	value := make([]byte, len(hexValue)/2)
	for i := 0; i < len(hexValue); i += 2 {
		var b byte
		n, err := fmt.Sscanf(hexValue[i:i+2], "%02x", &b)
		if err != nil || n != 1 {
			return Share{}, errors.New("неверный формат hex")
		}
		value[i/2] = b
	}
	
	share.Value = value
	return share, nil
}