package chain

import (
	"time"
)

const (
	//Input ..
	Input = iota
	//Implicit ..
	Implicit = iota
	//Scheduled ..
	Scheduled = iota
)

//TransactionMetadata 与上下文无关的缓存，如打包/解包/压缩和回收钥匙
type TransactionMetadata struct {
	packedTrx        *PackedTransaction
	sigCPUUsage      time.Duration
	recoveredPubKeys []PublicKeyType

	Implicit        bool
	Scheduled       bool
	Accepted        bool
	BilledCPUTimeUs uint32
}

//ID ..
func (t TransactionMetadata) ID() *TransactionIDType {
	return t.packedTrx.ID()
}

//PackedTrx ..
func (t TransactionMetadata) PackedTrx() *PackedTransaction {
	return t.packedTrx
}
