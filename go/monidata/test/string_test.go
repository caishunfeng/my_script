package test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func Test_S(t *testing.T) {
	fmt.Println(GetRandomString(3))
}

func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
