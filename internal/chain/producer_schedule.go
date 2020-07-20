package chain

import "crypto"

//ProducerKey ..
type ProducerKey struct {
	producerName    Name
	blockSigningKey crypto.PublicKey
}

//ProducerScheduleType ..
type ProducerScheduleType struct {
	version   uint32
	producers []ProducerKey
}

//BlockSigningAuthorityV0 ..
type BlockSigningAuthorityV0 struct {
	Threshold uint32
	keys      []KeyWeight
}

//ProduceAuthority ..
type ProduceAuthority struct {
	ProducerName Name
	Authority    BlockSigningAuthorityV0
}

//ProduceAuthoritySchedule ..
type ProduceAuthoritySchedule struct {
	version   uint32
	Producers []ProduceAuthority
}
