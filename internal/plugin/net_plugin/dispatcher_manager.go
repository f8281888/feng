package netplugin

import (
	"context"
	"sync"
)

type peerBlockStateIndex struct {
}

type nodeTransactionIndex struct {
}

//DispatcherManager ..
type DispatcherManager struct {
	blkStateMtx  sync.Mutex
	blkState     peerBlockStateIndex
	localtxnsMtx sync.Mutex
	locakTxns    nodeTransactionIndex
	strand       context.Context
}

//DispatcherReset ..
func DispatcherReset(ioc context.Context) *DispatcherManager {
	d := &DispatcherManager{
		strand: ioc,
	}

	return d
}
