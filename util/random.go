package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijkmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Generatore di INT random
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// Generatore di STRINGHE random di lunghezza n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// Generatore di proprietari di conto
func RandomOwner() string {
	return RandomString(6)
}

// Generatore di amount
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// Generatore di currency
func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "YEN"}

	n := len(currencies)

	return currencies[rand.Intn(n)]

}
