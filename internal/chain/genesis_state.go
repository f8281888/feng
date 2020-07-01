package chain

import (
	"feng/internal/fc/crypto"
	"time"
)

var eosioRootKey string

//GenesisState ..
type GenesisState struct {
	initialConfiguration Config
	//time_point
	initialTimestamp int64
	initialKey       crypto.PublicKey
}

//Init ..
func (g *GenesisState) Init() {
	t, _ := time.Parse("2020-06-01T12:00:00", "2020-06-01T12:00:00")
	g.initialTimestamp = t.Unix()
}

//ComputeChainID ..
func (g GenesisState) ComputeChainID() CIDType {
	c := CIDType{}
	c.sha256 = &crypto.Sha256{}
	c.sha256.Hash = []byte("1")
	return c
}
