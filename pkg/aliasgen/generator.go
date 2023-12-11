package aliasgen

import (
	"math/rand"
	"strings"
	"time"
)

var (
	genAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_"
	genLen      = 10
)

func Generate() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	builder := strings.Builder{}
	builder.Grow(genLen)
	for i := 0; i < genLen; i++ {
		index := rand.Intn(len(genAlphabet))
		builder.WriteByte(genAlphabet[index])
	}

	return builder.String()
}
