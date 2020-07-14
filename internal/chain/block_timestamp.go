package chain

//BlockTimestamp BlockTimestamp
type BlockTimestamp struct {
	time                uint32
	blockIntervalMs     uint16
	blockTimestampEpoch uint64
}

//Next 下一秒
func (b *BlockTimestamp) Next() {
	b.time++
}
