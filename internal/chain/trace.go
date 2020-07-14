package chain

import (
	"time"
)

//AccountDelta ..
type AccountDelta struct {
	Account AccountName
	Delta   int64
}

//ActionTrace ..
type ActionTrace struct {
	ActionOrdinal                          uint32
	CreatorActionOrdinal                   uint32
	ClosestUnnotifiedAncestorActionOrdinal uint32
	//Receipt ..
	Receipt          ActionReceipt
	Receiver         AccountName
	Act              Action
	ContextFree      bool
	Elapsed          time.Duration
	Console          string
	TrxID            TransactionIDType
	BlockNum         uint32
	BlockTime        BlockTimestamp
	ProducerBlockID  BlockIDType
	AccountRAMDeltas []AccountDelta
	//Except
	ErrorCode uint64
}

//TransactionTrace ..
type TransactionTrace struct {
	ID              TransactionIDType
	BlockNum        uint32
	BlockTime       BlockTimestamp
	ProducerBlockID BlockIDType
	Receipt         TransactionReceiptHeader
	Elapsed         time.Duration
	NetUsage        uint64
	Scheduled       bool
	ActionTraces    []ActionTrace
	AccountRAMDelta AccountDelta
	FailedDtrxTrace *TransactionTrace
	//Except except_ptr
	ErrorCode uint64
}
