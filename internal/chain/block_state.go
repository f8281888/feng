package chain

import (
	"feng/internal/fc/crypto"
)

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

//IsValid ..
func (b BlockState) IsValid() bool {
	return b.validated
}

//BranchType ..
type BranchType []*BlockState

//Copy ..
func (b *BlockState) Copy(h BlockHeaderState) {
	b.id = &crypto.Sha256{}
	b.id = h.id
	b.header = h.header
	b.pendingSchedule = h.pendingSchedule
	b.headerExts = h.headerExts
	b.activatedProtocolFeatures = h.activatedProtocolFeatures
	b.BlockNum = h.BlockNum
	b.dposIrreversibleBlocknum = h.dposIrreversibleBlocknum
	b.dposProposedIrreversibleBlocknum = h.dposProposedIrreversibleBlocknum
	b.ActiveSchedule = h.ActiveSchedule
}
