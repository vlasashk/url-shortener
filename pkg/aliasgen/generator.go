package aliasgen

import (
	"math/rand"
	"strings"
	"time"
)

const (
	GenAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_"
	GenLen      = 10
)

func Generate() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	builder := strings.Builder{}
	builder.Grow(GenLen)
	for i := 0; i < GenLen; i++ {
		index := rand.Intn(len(GenAlphabet))
		builder.WriteByte(GenAlphabet[index])
	}

	return builder.String()
}
