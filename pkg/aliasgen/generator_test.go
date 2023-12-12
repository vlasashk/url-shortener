package aliasgen_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/vlasashk/url-shortener/pkg/aliasgen"
	"testing"
)

func TestGenerate(t *testing.T) {
	t.Run("GeneratedAliasLength", func(t *testing.T) {
		alias := aliasgen.Generate()
		assert.Equal(t, aliasgen.GenLen, len(alias))
	})

	t.Run("GeneratedAliasCharacters", func(t *testing.T) {
		alias := aliasgen.Generate()
		for _, char := range alias {
			assert.Contains(t, aliasgen.GenAlphabet, string(char))
		}
	})

	t.Run("GeneratedAliasUniqueness", func(t *testing.T) {
		const numAliases = 10000
		aliases := make(map[string]bool)

		for i := 0; i < numAliases; i++ {
			alias := aliasgen.Generate()
			aliases[alias] = true
		}

		assert.Equal(t, numAliases, len(aliases))
	})
}
