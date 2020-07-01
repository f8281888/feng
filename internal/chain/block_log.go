package chain

import (
	"feng/internal/fc/crypto"
	"os"
)

//BlockLog ..
type BlockLog struct {
	myBlockLogImpl *BlockLogImpl
}

//BlockLogImpl ..
type BlockLogImpl struct {
	Head                     *SignedBlock
	HeadID                   crypto.Sha256
	blockFile                os.File
	indexFile                os.File
	openFiles                bool
	genesisWrittenToBlockLog bool
	version                  uint32
	firstBlockNum            uint32
}

//Head ..
func (b BlockLog) Head() *SignedBlock {
	return b.myBlockLogImpl.Head
}

//Reset ..
func (b BlockLog) Reset(gs GenesisState, firstBlock SignedBlock) {
	println("BlockLog Reset")
	b.myBlockLogImpl.Reset(gs, firstBlock, 1)
}

//Reset ..
func (b BlockLogImpl) Reset(gs GenesisState, firstBlock SignedBlock, firstBnum uint32) {
	println("BlockLogImpl Reset")
}
