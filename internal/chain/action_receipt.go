package chain

//ActionReceipt ..
type ActionReceipt struct {
	Receiver       AccountName
	ActDigist      DigestType
	GlobalSequence uint64
	RecvSequence   uint64
	AuthSequence   map[AccountName]uint64
	CodeSequence   uint32
	AbiSequence    uint32
}
