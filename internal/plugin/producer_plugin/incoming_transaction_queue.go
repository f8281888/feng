package producerplugin

import (
	"feng/internal/chain"
)

//IncomingTransactionsType ..
type IncomingTransactionsType struct {
	TransactionMetadata *chain.TransactionMetadata
	IsTrue              bool
	Trace               *chain.TransactionTrace
}

//IncomingTransactionQueue ..
type IncomingTransactionQueue struct {
	MaxIncomingTransactionQueueSize uint64
	SizeInBytes                     uint64
	IncomingTransactions            IncomingTransactionsType
}

//SetMaxIncomingTransactionQueueSize ..
func (i *IncomingTransactionQueue) SetMaxIncomingTransactionQueueSize(v uint64) {
	i.MaxIncomingTransactionQueueSize = v
}
