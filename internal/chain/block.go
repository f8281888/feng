package chain

//TransactionReceipt ..
type TransactionReceipt struct {
}

//SignedBlock ..
type SignedBlock struct {
	transactions []TransactionReceipt
}

//BlockNum ..
func (s SignedBlock) BlockNum() uint32 {
	return 0
}
