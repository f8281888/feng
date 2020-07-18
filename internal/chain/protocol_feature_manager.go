package chain

import (
	"feng/internal/log"
)

//ProtocolFeatureManager ..
type ProtocolFeatureManager struct {
	pfm         *ProtocolFeatureManager
	index       uint32
	initialized bool
}

func (p ProtocolFeatureManager) isInitialized() bool {
	return p.initialized
}

//PoppedBlocksTo ..
func (p *ProtocolFeatureManager) PoppedBlocksTo(blockNum uint32) {
	if p.isInitialized() != true {
		log.Assert("protocol_feature_manager is not yet initialized")
	}

	//TODO
}
