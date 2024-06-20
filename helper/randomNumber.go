package helper

import "math/rand"

func RandomNumber(n int) string {
	var charset = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	temp := make([]byte, n)
	for i := range temp {
		temp[i] = charset[rand.Intn(len(charset))]
	}
	return string(temp)
}
