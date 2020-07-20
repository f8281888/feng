package chain

//BlockTimestamp BlockTimestamp
type BlockTimestamp struct {
	time                uint64
	blockIntervalMs     uint16
	blockTimestampEpoch uint64
	Slot                uint32
}

//Next 下一秒
func (b *BlockTimestamp) Next() {
	b.time++
}

//Time ..
func (b *BlockTimestamp) Time() uint64 {
	return b.time
}

//IsGreater ..
func (b BlockTimestamp) IsGreater(rhs BlockTimestamp) bool {
	return b.Slot > rhs.Slot
}
