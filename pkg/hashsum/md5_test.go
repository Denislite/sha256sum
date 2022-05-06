package hashsum

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestMD5_Hash(t *testing.T) {
	hasher := NewMD5()

	r := strings.NewReader("hello")

	hash, err := hasher.Hash(r)
	if assert.Nil(t, err) {
		assert.Equal(t, "5d41402abc4b2a76b9719d911017c592", hash)
	}
}
