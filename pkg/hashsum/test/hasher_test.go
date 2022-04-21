package test

import (
	"sha256sum/pkg/hashsum"
	"testing"
)

func TestFileHash(t *testing.T) {
	cases := []struct {
		path     string
		hashType string
		want     string
	}{
		{"./testfiles/1.txt", "sha256", "1785cfc3bc6ac7738e8b38cdccd1af12563c2b9070e07af336a1bf8c0f772b6a"},
		{"./testfiles/2.txt", "sha256", "fcec91509759ad995c2cd14bcb26b2720993faf61c29d379b270d442d92290eb"},
		{"./testfiles/3.txt", "sha256", "54d626e08c1c802b305dad30b7e54a82f102390cc92c7d4db112048935236e9c"},
	}

	for _, c := range cases {
		got, err := hashsum.FileHash(c.path, c.hashType)
		if err != nil {
			t.Error(err)
			return
		}
		if got.HashValue != c.want {
			t.Errorf(".Expected %s, got %s", c.want, got.HashValue)
			return
		}
	}
}
