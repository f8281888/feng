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

//ProducerAuthority ..
type ProducerAuthority struct {
	ProducerName Name
	Authority    BlockSigningAuthorityV0
}

//ProducerAuthoritySchedule ..
type ProducerAuthoritySchedule struct {
	version   uint32
	Producers []ProducerAuthority
}

//ProducerScheduleChangeExtension ..
type ProducerScheduleChangeExtension struct {
	ProducerAuthoritySchedule
}
