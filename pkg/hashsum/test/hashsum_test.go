package test

import (
	"sha256sum/internal/repository"
	"sha256sum/internal/service"
	"testing"
)

var s service.Service

func init() {
	s = *service.NewService(repository.Repository{}, "sha256")
}

func TestFileHash(t *testing.T) {

	cases := []struct {
		path string
		want string
	}{
		{"./testfiles/1.txt", "1785cfc3bc6ac7738e8b38cdccd1af12563c2b9070e07af336a1bf8c0f772b6a"},
		{"./testfiles/2.txt", "fcec91509759ad995c2cd14bcb26b2720993faf61c29d379b270d442d92290eb"},
		{"./testfiles/3.txt", "54d626e08c1c802b305dad30b7e54a82f102390cc92c7d4db112048935236e9c"},
	}

	for _, c := range cases {
		got, err := s.FileHash(c.path)
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
