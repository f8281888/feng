package netplugin

import (
	"feng/internal/fc/crypto"
	"time"
)

//HandshakeMessage ..
type HandshakeMessage struct {
	NetworkVersion           uint16
	ChainID                  crypto.Sha256
	NodeID                   crypto.Sha256
	Key                      crypto.PublicKey
	Time                     time.Duration
	Token                    crypto.Sha256
	Sig                      crypto.Signature
	P2pAddress               string
	LastIrreversibleBlockNum uint32
	LastIrreversibleBlockID  crypto.Sha256
	HeadNum                  uint32
	HeadID                   crypto.Sha256
	os                       string
	agent                    string
	generation               int16
}
