package chain

import (
	"feng/internal/fc/crypto"
)

//CIDType ..
type CIDType struct {
	sha256 *crypto.Sha256
}

//SetCIDType ..
func (c *CIDType) SetCIDType(s *crypto.Sha256) {
	c.sha256 = s
}

//GetSha256 ..
func (c CIDType) GetSha256() *crypto.Sha256 {
	return c.sha256
}
