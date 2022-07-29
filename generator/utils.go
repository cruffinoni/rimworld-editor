package generator

import "math/rand"

const (
	stringFixedSize = 5
	letterBytes     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func generateRandomString() string {
	b := make([]byte, stringFixedSize)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func isListTag(tag string) bool {
	return tag == "li" || tag == "list"
}
