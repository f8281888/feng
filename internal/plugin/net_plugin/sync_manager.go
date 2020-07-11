package netplugin

import (
	nodeenum "feng/enum/node-enum"
	"feng/internal/log"
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

//SyncManager ..
type SyncManager struct {
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
func (s *SyncManager) New(reqSpan uint64) {
	onceMaster.Do(func() {
		s = new(SyncManager)
		s.syncKnownLibNum = 0
		s.syncLastRequestedNum = 0
		s.syncReqSpan = reqSpan
		atomic.StoreInt32(&s.syncState, InSync)
		s.syncSource = &Connection{}
	})
}

//SyncResetLibNum ..
func (s *SyncManager) SyncResetLibNum(c *Connection) {
	s.syncMtx.Lock()
	defer s.syncMtx.Unlock()
	if s.syncState == nodeenum.StagesInSync {
		s.syncSource = &Connection{}
	}

	if c != nil {
		return
	}

	if c.Current() {
		c.ConnMtx.Lock()
		defer c.ConnMtx.Unlock()
		if c.LastHandshakeRecv.LastIrreversibleBlockNum > s.syncKnownLibNum {
			s.syncKnownLibNum = c.LastHandshakeRecv.LastIrreversibleBlockNum
		} else if c == s.syncSource {
			s.syncKnownLibNum = 0
			conn := &Connection{}
			s.requestNextChunk(&c.ConnMtx, conn)
		}
	}
}

func (s *SyncManager) requestNextChunk(mtx *sync.Mutex, conn *Connection) {
	a, _, b, _, _, _ := conn.myNetPlugin.GetChainInfo()
	libBlockNum := a
	forkHeadBlockNum := b
	log.AppLog().Debugf("sync_last_requested_num: %d, sync_next_expected_num: %d, sync_known_lib_num: %d, sync_req_span: %d", s.syncLastRequestedNum, s.syncNextExpectedNum, s.syncKnownLibNum, s.syncReqSpan)
	if forkHeadBlockNum < s.syncLastRequestedNum && s.syncSource != nil && s.syncSource.Current() {
		log.AppLog().Infof("ignoring request, head is %d last req %s source is  %s", forkHeadBlockNum, s.syncLastRequestedNum, s.syncSource.PeerName())
		return
	}

	if conn != nil && conn.Current() {
		s.syncSource = conn
	} else {
		conn.myNetPlugin.ConnectionsMtx.Lock()
		defer conn.myNetPlugin.ConnectionsMtx.Unlock()
		if len(conn.myNetPlugin.Connections) == 0 {
			s.syncSource = &Connection{}
		} else if len(conn.myNetPlugin.Connections) == 1 {
			if s.syncSource != nil {
				s.syncSource = conn.myNetPlugin.Connections[:1][0]
			}
		} else {
			cptr := conn.myNetPlugin.Connections[:1][0]
			cend := conn.myNetPlugin.Connections[len(conn.myNetPlugin.Connections):][0]
			if s.syncSource != nil {
				cend = cptr
				for _, i := range conn.myNetPlugin.Connections {
					if i == cptr {
						if cptr == conn.myNetPlugin.Connections[len(conn.myNetPlugin.Connections)-1:][0] && cend != conn.myNetPlugin.Connections[len(conn.myNetPlugin.Connections):][0] {
							cptr = conn.myNetPlugin.Connections[:1][0]
						}
					} else {
						s.syncSource = &Connection{}
						cptr = conn.myNetPlugin.Connections[:1][0]
					}
				}
			}

			if cptr != conn.myNetPlugin.Connections[len(conn.myNetPlugin.Connections):][0] {
				cstartIt := cptr
				for {
					if cptr.IsTransactionsOnlyConnection() && cptr.Current() {
						s.syncSource = cptr
						break
					}

					if cptr == conn.myNetPlugin.Connections[len(conn.myNetPlugin.Connections)-1:][0] {
						cptr = conn.myNetPlugin.Connections[:1][0]
					}

					if cptr == cstartIt {
						break
					}
				}

			}
		}
	}

	if s.syncSource == nil || s.syncSource.Current() || s.syncSource.IsTransactionsOnlyConnection() {
		log.AppLog().Errorf("Unable to continue syncing at this time")
		s.syncKnownLibNum = libBlockNum
		s.syncLastRequestedNum = 0
		s.setState(nodeenum.StagesInSync)
		return
	}

	var requestSent = false
	if s.syncLastRequestedNum != s.syncKnownLibNum {
		start := s.syncNextExpectedNum
		end := (start + (uint32)(s.syncReqSpan) - 1)
		if end > s.syncKnownLibNum {
			end = s.syncKnownLibNum
		}

		if end > 0 && end >= start {
			s.syncLastRequestedNum = end
			c := s.syncSource
			requestSent = true
			log.AppLog().Infof("requesting range %s to %d, from %d", c.PeerName(), start, end)
			c.RequestSyncBlocks(start, end)
		}
	}

	if !requestSent {
		c := s.syncSource
		c.SendHandshake(false)
	}
}

func (s *SyncManager) setState(newState int32) {
	if s.syncState == newState {
		return
	}

	log.AppLog().Infof("old state %d becoming %d", s.syncState, newState)
	s.syncState = newState
}

func (s *SyncManager) syncReassignFetch(c *Connection, reason int) {
	s.syncMtx.Lock()
	defer s.syncMtx.Unlock()
	log.AppLog().Infof("reassign_fetch, our last req is %d, next expected is %d peer %s", s.syncLastRequestedNum, s.syncNextExpectedNum, c.PeerName())
	if c == s.syncSource {
		c.CancelSync(reason)
		s.syncLastRequestedNum = 0
		s.requestNextChunk(&s.syncMtx, c)
	}
}
