package chain

//BlockState ..
type BlockState struct {
	BlockHeaderState
	block           *SignedBlock
	validated       bool
	pubKeysRecoverd bool
}
