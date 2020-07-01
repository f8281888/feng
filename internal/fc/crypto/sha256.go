package crypto

import (
	"crypto/sha256"
)

//Sha256 Sha256
type Sha256 struct {
	Hash []byte
}

//New ..
func (s *Sha256) New(input string) {
	h := sha256.New()
	h.Write([]byte(input))
	s.Hash = h.Sum(nil)
}
