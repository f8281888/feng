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
	transactions []TransactionReceipt
}

//BlockNum ..
func (s SignedBlock) BlockNum() uint32 {
	return 0
}
