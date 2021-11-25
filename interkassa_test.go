package interkassa

import (
	"os"
)

var instance = New(Config{
	ID:            requiredEnv("KASSA_ID"),
	SignAlgorithm: signAlgorithm(requiredEnv("KASSA_SIGN_ALGORITHM")),
	SignKey:       requiredEnv("KASSA_SIGN_KEY"),
	SignTestKey:   requiredEnv("KASSA_SIGN_TEST_KEY"),
})

func requiredEnv(name string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		panic("required environment not provided: " + name)
	}

	return value
}
