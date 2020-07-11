package crypto

import (
	"crypto/sha256"
)

//Sha256 Sha256
type Sha256 struct {
	Hash []byte
}

//New ..
func (s *Sha256) New(input string) Sha256 {
	h := sha256.New()
	h.Write([]byte(input))
	s.Hash = h.Sum(nil)
	return *s
}

func (s *Sha256) String() string {
	return string(s.Hash)
}
