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
