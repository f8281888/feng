package chain

import (
	"feng/internal/log"
	"time"
)

const (
	// NonSpeculative HEAD, READ_ONLY, IRREVERSIBLE
	NonSpeculative = iota
	//SpeculativeNonProducer will never produce
	SpeculativeNonProducer = iota
	// SpeculativeProducer can produce
	SpeculativeProducer = iota
)

const (
	//Unknown ..
	Unknown = iota
	//Persisted ..
	Persisted = iota
	//Forked ..
	Forked = iota
	//Aborted ..
	Aborted = iota
)

//UnappliedTransaction ..
type UnappliedTransaction struct {
	TrxMeta     *TransactionMetadata
	Expirty     time.Duration
	TrxEnumType uint32
}

//ID ..
func (u UnappliedTransaction) ID() *TransactionIDType {
	return u.TrxMeta.ID()
}

//UnappliedKey ..
type UnappliedKey struct {
	ID          *TransactionIDType
	Expirty     time.Duration
	TrxEnumType uint32
}

//UnappliedTrxQueueType ..
type UnappliedTrxQueueType map[*TransactionIDType]UnappliedTransaction

//type UnappliedTrxQueueType map[UnappliedKey]UnappliedTransaction

// // BuildSnapshotIndex 构建查询索引
// func BuildSnapshotIndex(t UnappliedTrxQueueType, list []UnappliedTransaction) {
// 	// 遍历所有数据
// 	for _, profile := range list {
// 		// 构建查询键
// 		key := UnappliedKey{
// 			ID:          profile.TrxMeta.ID(),
// 			Expirty:     profile.Expirty,
// 			TrxEnumType: profile.TrxEnumType,
// 		}

// 		// 保存查询键
// 		t[key] = profile
// 	}
// }

// // QuerySnapshotData 根据条件查询数据
// func QuerySnapshotData(u UnappliedTrxQueueType, key UnappliedKey) {
// 	// 根据键值查询数据
// 	result, ok := u[key]

// 	// 找到数据打印出来
// 	if ok {
// 		fmt.Println(result)
// 	} else {
// 		fmt.Println("no found")
// 	}
// }

//UnappliedTransactionQueue ..
type UnappliedTransactionQueue struct {
	mode  uint32
	queue UnappliedTrxQueueType
}

//SetMode ..
func (u *UnappliedTransactionQueue) SetMode(newMode uint32) {
	if u.mode == newMode {
		if !u.empty() {
			log.Assert("set_mode, queue required to be empty")
		}
	}

	u.mode = newMode
}

func (u *UnappliedTransactionQueue) empty() bool {
	return len(u.queue) <= 0
}

//AddAborted ..
func (u *UnappliedTransactionQueue) AddAborted(abortedTrxs []*TransactionMetadata) {
	if u.mode == NonSpeculative || u.mode == SpeculativeNonProducer {
		return
	}

	for _, trx := range abortedTrxs {
		expiry := trx.PackedTrx().Expiration()
		tmp := UnappliedTransaction{TrxMeta: trx, Expirty: expiry, TrxEnumType: Aborted}
		// queue.insert( { std::move( trx ), expiry, trx_enum_type::aborted } ); 多索引容器插入一条数据。。
		u.queue[trx.ID()] = tmp
	}
}
