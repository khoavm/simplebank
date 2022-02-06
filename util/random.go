package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

//Random Int 
func RandomInt(min, max int64) int64{
	return min + rand.Int63n(max - min + 1) 
}

//RandomString generates  a random strig of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// RandomOwner generates
func RandomOwner() string {
	return RandomString(6)
}
// RandomMoney generates
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency generates a ramdom currency code
func RandomCurrency() string{
	currencies := []string{"EUR", "USD", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}