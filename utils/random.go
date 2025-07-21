package utils

import (
	"fmt"
	"math/rand"
	"time"
)

var letters = []rune("abcdefghiiklmnopqrrstuvvwxzABCDEFGHIJKLMNOPQRSTUVWXYZ")

var initialConsonants = []string{
	"b", "c", "d", "g", "h", "k", "l", "m", "n", "ph", "qu", "s", "t", "th", "tr", "v",
}

var vowels = []string{
	"an", "ang", "inh", "ong", "i", "ai", "u", "o", "y", "ien", "uan", "anh",
}

var finalConsonants = []string{
	"", "n", "ng", "nh", "m", "t", "ch",
}

func init() {
	rand.NewSource(time.Now().UnixNano())
}

func RandInt(min, max int64) int64 {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return min + r1.Int63n(max-min+1) //
}

func RandomString(n int) string {
	b := make([]rune, n)

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	for i := range b {
		b[i] = letters[r1.Intn(99999)%len(letters)]
	}

	return string(b)
}

func RandomOwner() string {
	return RandomString(10)
}

func RandomMoney() int64 {
	return RandInt(100, 1000)
}

func RandomCurrency() string {
	currencies := []string{"USD", "USD", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

func GenerateVietnameseStyleUsername() string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	initial := initialConsonants[r1.Intn(len(initialConsonants))]
	vowel := vowels[r1.Intn(len(vowels))]
	final := finalConsonants[r1.Intn(len(finalConsonants))]

	base := initial + vowel + final

	// Thêm số 3 hoặc 4 chữ số phía sau
	suffix := r1.Intn(9000) + 100

	return fmt.Sprintf("%s%d", base, suffix)
}
