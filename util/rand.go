package util

import "math/rand"

// RandString randomly generate a specified length string based on the given string
func RandString(base string, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = base[rand.Intn(len(base))]
	}
	return string(b)
}
