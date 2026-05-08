package service

import (
	"fmt"
	"math"
	"strings"
)

var satuan = []string{
	"", "satu", "dua", "tiga", "empat", "lima",
	"enam", "tujuh", "delapan", "sembilan", "sepuluh",
	"sebelas",
}

func Terbilang(angka float64) (string, error) {
	if math.IsNaN(angka) || math.IsInf(angka, 0) {
		return "", fmt.Errorf("angka tidak valid")
	}

	if angka < 0 {
		t, err := Terbilang(-angka)
		if err != nil {
			return "", err
		}
		return "negatif " + t, nil
	}

	intPart := int64(angka)
	fracPart := angka - float64(intPart)

	if math.Abs(angka) > 999_999_999_999_999 {
		return "", fmt.Errorf("angka terlalu besar (maksimum 999 triliun)")
	}

	result := terbilangInt(intPart)

	if fracPart > 0.0001 {
		fracStr := fmt.Sprintf("%.10f", fracPart)[2:]
		fracStr = strings.TrimRight(fracStr, "0")
		digits := make([]string, len(fracStr))
		for i, c := range fracStr {
			digits[i] = terbilangDigit(int(c - '0'))
		}
		result += " koma " + strings.Join(digits, " ")
	}

	if result == "" {
		return "nol", nil
	}
	return strings.TrimSpace(result), nil
}

func terbilangDigit(d int) string {
	if d == 0 {
		return "nol"
	}
	return satuan[d]
}

func terbilangInt(n int64) string {
	if n == 0 {
		return ""
	}
	if n < 12 {
		return satuan[n]
	}
	if n < 20 {
		return satuan[n-10] + " belas"
	}
	if n < 100 {
		prefix := satuan[n/10]
		if n/10 == 1 {
			prefix = "se"
		}
		rem := n % 10
		if rem == 0 {
			return prefix + "puluh"
		}
		return prefix + "puluh " + satuan[rem]
	}
	if n < 200 {
		rem := n % 100
		if rem == 0 {
			return "seratus"
		}
		return "seratus " + terbilangInt(rem)
	}
	if n < 1000 {
		prefix := terbilangInt(n / 100)
		rem := n % 100
		if rem == 0 {
			return prefix + " ratus"
		}
		return prefix + " ratus " + terbilangInt(rem)
	}
	if n < 2000 {
		rem := n % 1000
		if rem == 0 {
			return "seribu"
		}
		return "seribu " + terbilangInt(rem)
	}
	if n < 1_000_000 {
		prefix := terbilangInt(n / 1000)
		rem := n % 1000
		if rem == 0 {
			return prefix + " ribu"
		}
		return prefix + " ribu " + terbilangInt(rem)
	}
	if n < 1_000_000_000 {
		prefix := terbilangInt(n / 1_000_000)
		rem := n % 1_000_000
		if rem == 0 {
			return prefix + " juta"
		}
		return prefix + " juta " + terbilangInt(rem)
	}
	if n < 1_000_000_000_000 {
		prefix := terbilangInt(n / 1_000_000_000)
		rem := n % 1_000_000_000
		if rem == 0 {
			return prefix + " miliar"
		}
		return prefix + " miliar " + terbilangInt(rem)
	}
	prefix := terbilangInt(n / 1_000_000_000_000)
	rem := n % 1_000_000_000_000
	if rem == 0 {
		return prefix + " triliun"
	}
	return prefix + " triliun " + terbilangInt(rem)
}
