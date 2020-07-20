package chain

import (
	"feng/internal/fc/crypto"
)

//BlockHeaderStateCommon ..
type BlockHeaderStateCommon struct {
	BlockNum                         uint32
	dposProposedIrreversibleBlocknum uint32
	dposIrreversibleBlocknum         uint32
	ActiveSchedule                   ProduceAuthoritySchedule
}

//ScheduleInfo ..
type ScheduleInfo struct {
	scheduleLibNum uint32
	scheduleHash   crypto.Sha256
	schedule       ProducerScheduleType
}

//BlockHeaderState ..
type BlockHeaderState struct {
	BlockHeaderStateCommon
	id              crypto.Sha256
	header          SignedBlockHeader
	pendingSchedule ScheduleInfo
}

//GetScheduledProducer ..
func (b BlockHeaderState) GetScheduledProducer(t BlockTimestamp) ProduceAuthority {
	index := int(t.Slot) % len(b.ActiveSchedule.producers) * int(ProducerRepetitions)
	index /= int(ProducerRepetitions)
	return b.ActiveSchedule.producers[index]
}
