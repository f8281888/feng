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

//SelectIDInterface ..
type SelectIDInterface interface{}

//SelectIDs ..
type SelectIDs struct {
	Pending uint32
	Mode    uint32
	ids     SelectIDInterface
}

//RequestMessage ..
type RequestMessage struct {
	ReqTrx    *SelectIDs
	ReqBlocks *SelectIDs
}

//IsEmpty ..
func (r *RequestMessage) IsEmpty() bool {
	return r.ReqTrx != nil && r.ReqBlocks != nil
}

//syncRequestMessage ..
type syncRequestMessage struct {
	startBlock uint32
	endBlock   uint32
}

//NetMessage 网络请求集合
type NetMessage interface {
}

type goAwayMessage struct {
	reason int
	nodeID crypto.Sha256
}
