package chain

import (
	"feng/internal/fc/io"
	"feng/internal/fc/stl"
	"time"
)

const (
	//None ..
	None = iota
	//Zlib ..
	Zlib = iota
)

//TransactionHeader 事务头
type TransactionHeader struct {
	//事务的到期时间
	Expiration time.Duration
	//指定最后2^16 的区块
	RefBlockNum uint16
	//指定get_ref_blocknum处blockid的低位32位
	RefBlockPrefix uint32
	//最大net 上限
	MaxNetUsageWords uint32
	//最大CPU 上限
	MaxCPUUsageMs uint8
	//延迟事务的秒数，在此期间可以取消事务
	delaySec uint32
}

//GetRefBlocknum ..
func (t TransactionHeader) GetRefBlocknum(headBlocknum BlockNumType) BlockNumType {
	return ((headBlocknum / 0xffff) * 0xffff) + headBlocknum%0xffff
}

//SetReferenceBlock ..
func (t *TransactionHeader) SetReferenceBlock(referenceBlock BlockIDType) {
	// t.RefBlockNum = referenceBlock.Hash[0]
	// t.RefBlockPrefix = referenceBlock.Hash[1]
}

//Transaction 一个事务可以由多个action 组成
type Transaction struct {
	ContextFreeActions    []Action
	Actions               []Action
	TransactionExtensions []stl.Pair
}

//SignedTransaction ..
type SignedTransaction struct {
	Signatures      []SignatureType
	ContextFreeData []byte
}

//PackedTransaction ..
type PackedTransaction struct {
	signatures            []SignatureType
	compression           io.EnumType
	packedContextFreeData []byte
	packedTrx             []byte
	unpackedTrx           SignedTransaction
	trxID                 TransactionIDType
}

//ID ..
func (p PackedTransaction) ID() *TransactionIDType {
	return &p.trxID
}
