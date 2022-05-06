package hashsum

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSha256_Hash(t *testing.T) {
	hasher := NewSha256()

	r := strings.NewReader("hello")

	hash, err := hasher.Hash(r)
	if assert.Nil(t, err) {
		assert.Equal(t, "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824", hash)
	}
}
