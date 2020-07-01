package netplugin

import (
	"sync"
	"sync/atomic"
)

const (
	//LibCatchup ..
	LibCatchup = 0
	//HeadCatchup ..
	HeadCatchup = 1
	//InSync ..
	InSync = 2
)

//SyncMaster ..
type SyncMaster struct {
	syncMtx              sync.Mutex
	syncKnownLibNum      uint32
	syncLastRequestedNum uint32
	syncNextExpectedNum  uint32
	syncReqSpan          uint64
	syncSource           *Connection
	syncState            int32
}

var onceMaster sync.Once

//New ..
func (s *SyncMaster) New(reqSpan uint64) {
	onceMaster.Do(func() {
		s = new(SyncMaster)
		s.syncKnownLibNum = 0
		s.syncLastRequestedNum = 0
		s.syncReqSpan = reqSpan
		atomic.StoreInt32(&s.syncState, InSync)
		s.syncSource = &Connection{}
	})
}
