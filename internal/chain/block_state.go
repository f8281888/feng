package chain

import (
	"feng/internal/fc/crypto"
)

//BlockState ..
type BlockState struct {
	BlockHeaderState
	Block           *SignedBlock
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

//IsValid ..
func (b BlockState) IsValid() bool {
	return b.validated
}

//BranchType ..
type BranchType []*BlockState

//Copy ..
func (b *BlockState) Copy(h BlockHeaderState) {
	b.ID = &crypto.Sha256{}
	b.ID = h.ID
	b.Header = h.Header
	b.pendingSchedule = h.pendingSchedule
	b.headerExts = h.headerExts
	b.activatedProtocolFeatures = h.activatedProtocolFeatures
	b.BlockNum = h.BlockNum
	b.dposIrreversibleBlocknum = h.dposIrreversibleBlocknum
	b.DposProposedIrreversibleBlocknum = h.DposProposedIrreversibleBlocknum
	b.ActiveSchedule = h.ActiveSchedule
}

//TrxsMetas ..
func (b BlockState) TrxsMetas() []*TransactionMetadata {
	return b.cachedTrxs
}
