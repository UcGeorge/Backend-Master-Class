package util

import "math/rand"

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	// This function is intentionally left blank.
}

// RandomInt generates a random integer between min and max.
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n.
func RandomString(n int) string {
	bytes := make([]byte, n)
	for i := range bytes {
		bytes[i] = letterBytes[rand.Int63n(int64(len(letterBytes)))]
	}
	return string(bytes)
}

// RandomOwner generates a random owner name.
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney generates a random amount of money.
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomAmount() int64 {
	return RandomInt(-1000, 1000)
}

// RandomCurrency generates a random currency code.
func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD"}
	n := len(currencies)
	return currencies[rand.Int63n(int64(n))]
}
