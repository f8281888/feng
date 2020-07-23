package chain

import (
	"feng/internal/fc/crypto"
)

//BlockHeaderStateCommon ..
type BlockHeaderStateCommon struct {
	BlockNum                         uint32
	dposProposedIrreversibleBlocknum uint32
	dposIrreversibleBlocknum         uint32
	ActiveSchedule                   ProducerAuthoritySchedule
}

//ScheduleInfo ..
type ScheduleInfo struct {
	scheduleLibNum uint32
	scheduleHash   crypto.Sha256
	schedule       ProducerAuthoritySchedule
}

//BlockHeaderExtensionTypes ..
type BlockHeaderExtensionTypes struct {
	protocolFeatureActivation       *ProtocolFeatureActivation
	producerScheduleChangeExtension *ProducerScheduleChangeExtension
}

//BlockHeaderState ..
type BlockHeaderState struct {
	BlockHeaderStateCommon
	id              *crypto.Sha256
	header          SignedBlockHeader
	pendingSchedule ScheduleInfo
	//flat_multimap<uint16_t, block_header_extension>
	headerExts                map[uint16]BlockHeaderExtensionTypes
	activatedProtocolFeatures ProtocolFeatureActivationSet
}

//GetScheduledProducer ..
func (b BlockHeaderState) GetScheduledProducer(t BlockTimestamp) ProducerAuthority {
	index := int(t.Slot) % len(b.ActiveSchedule.Producers) * int(ProducerRepetitions)
	index /= int(ProducerRepetitions)
	return b.ActiveSchedule.Producers[index]
}

//GetNewProtocolFeatureActivations ..
func (b BlockHeaderState) GetNewProtocolFeatureActivations() []DigestType {
	noActivations := []DigestType{}
	if b.headerExts[0].protocolFeatureActivation == nil {
		return noActivations
	}

	return b.headerExts[0].protocolFeatureActivation.ProtocolFeatures
}
