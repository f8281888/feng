package chain

//BlockState ..
type BlockState struct {
	BlockHeaderState
	block           *SignedBlock
	validated       bool
	pubKeysRecoverd bool
	cachedTrxs      []*TransactionMetadata
}

//ExtractTrxsMetas ..
func (b *BlockState) ExtractTrxsMetas() []*TransactionMetadata {
	b.pubKeysRecoverd = false
	result := b.cachedTrxs
	b.cachedTrxs = b.cachedTrxs[:0:0]
	return result
}
