package chain

import "feng/internal/fc/io"

const (
	//Executed succeed, no error handler executed
	Executed = iota
	//SoftFail objectively failed (not executed), error handler executed
	SoftFail = iota
	//HardFail objectively failed and error handler objectively failed thus no state change
	HardFail = iota
	//Delayed transaction delayed/deferred/scheduled for future execution
	Delayed = iota
	//Expired transaction expired and storage space refuned to user
	Expired = iota
)

//TransactionReceiptHeader ..
type TransactionReceiptHeader struct {
	Status        io.EnumType
	CPUUsageUs    uint32
	NetUsageWords uint32
}

//TransactionReceipt ..
type TransactionReceipt struct {
}

//SignedBlock ..
type SignedBlock struct {
	SignedBlockHeader
	Transactions []TransactionReceipt
}

//BlockNum ..
func (s SignedBlock) BlockNum() uint32 {
	return 0
}

//Copy ..
func (s *SignedBlock) Copy(copy SignedBlockHeader) {
	s.producerSignature = copy.producerSignature
	s.Timestamp = copy.Timestamp
	s.Producer = copy.Producer
	s.Confirmed = copy.Confirmed
	s.previous = copy.previous
	s.transactionMroot = copy.transactionMroot
	s.actionMroot = copy.actionMroot
	s.scheduleVersion = copy.scheduleVersion
	s.newProducers = copy.newProducers
	s.extensionsType = copy.extensionsType
}
