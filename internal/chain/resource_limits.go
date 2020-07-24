package chain

import (
	"feng/internal/chainbase"
)

//AccountResourceLimit ..
type AccountResourceLimit struct {
	Used      uint64
	Available uint64
	Max       uint64
}

//ResourceLimitsManager ..
type ResourceLimitsManager struct {
	db chainbase.DataBase
}

//GetBlockCPULimit ..TODO
func (r ResourceLimitsManager) GetBlockCPULimit() uint64 {
	//state := r.db.r
	return 0
}
