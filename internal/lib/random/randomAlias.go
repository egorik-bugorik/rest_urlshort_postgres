package random

import (
	"math/rand"
	"time"
)

func NewAlias(i int) string {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	letters := []rune("qwertyuiopasdfghjklzxcvbnm")

	b := make([]rune, i)

	for i, _ := range b {
		b[i] = letters[r.Intn(len(letters))]

	}

	return string(b)

}
