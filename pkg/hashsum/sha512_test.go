package hashsum

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSha512_Hash(t *testing.T) {
	hasher := NewSha512()

	r := strings.NewReader("hello")

	hash, err := hasher.Hash(r)
	if assert.Nil(t, err) {
		assert.Equal(t, "9b71d224bd62f3785d96d46ad3ea3d73319bfbc2890caadae2dff72519673ca72323c3d99ba5c11d7c7acc6e14b8c5da0c4663475c2e5c3adef46f73bcdec043", hash)
	}
}
