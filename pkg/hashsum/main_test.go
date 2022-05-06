package hashsum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	hasher, _, err := New("something")
	if assert.Nil(t, err) {
		assert.IsType(t, &Sha256{}, hasher)
	}
}

func TestNewMD5(t *testing.T) {
	hasher, _, err := New("md5")
	if assert.Nil(t, err) {
		assert.IsType(t, &MD5{}, hasher)
	}
}

func TestNewSha512(t *testing.T) {
	hasher, _, err := New("sha512")
	if assert.Nil(t, err) {
		assert.IsType(t, &Sha512{}, hasher)
	}
}
