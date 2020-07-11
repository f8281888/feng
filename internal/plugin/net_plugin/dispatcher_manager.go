package netplugin

import (
	"context"
	nodeenum "feng/enum/node-enum"
	"feng/internal/fc/crypto"
	"feng/internal/log"
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

func (d *DispatcherManager) retryFetch(c *Connection) {
	log.AppLog().Infof("retry fetch")
	var lastReq RequestMessage
	var bids []crypto.Sha256
	var bid crypto.Sha256
	c.ConnMtx.Lock()
	if c.LastReq.IsEmpty() {
		return
	}

	log.AppLog().Infof("failed to fetch from %s", c.PeerAddress())
	if c.LastReq.ReqBlocks.Mode == nodeenum.IDListModeNormal && c.LastReq.ReqBlocks.ids == nil {
		bids = c.LastReq.ReqBlocks.ids.([]crypto.Sha256)
		bid = bids[len(bids)-1:][0]
	} else {
		log.AppLog().Errorf("no retry, block mpde = %d trx mode = %s", c.LastReq.ReqBlocks.Mode, c.LastReq.ReqTrx.Mode)
		return
	}

	lastReq = c.LastReq
	c.ConnMtx.Unlock()

	k := func() bool {
		if c.LastReq.IsEmpty() {
			return true
		}

		sendit := d.peerHasBlock(bid, c.ConnectionID)
		if sendit {
			c.enqueue(lastReq)
			c.fetchWait()
			c.ConnMtx.Lock()
			c.LastReq = lastReq
			c.ConnMtx.Unlock()
			return false
		}

		return true
	}

	f := func(t func() bool) {
		c.myNetPlugin.ConnectionsMtx.Lock()
		defer c.myNetPlugin.ConnectionsMtx.Unlock()
		for _, c := range c.myNetPlugin.Connections {
			if c.IsTransactionsOnlyConnection() {
				continue
			}

			if !t() {
				return
			}
		}
	}

	f(k)
	log.AppLog().Infof("no peer has last_req")
	if c.Connected() {
		c.enqueue(lastReq)
		c.fetchWait()
	}
}

func (d *DispatcherManager) peerHasBlock(blkid crypto.Sha256, connectionID uint32) bool {
	d.blkStateMtx.Lock()
	defer d.blkStateMtx.Unlock()
	return true
	//查找，TODO
	// const auto blk_itr = blk_state.get<by_id>().find( std::make_tuple( connection_id, std::ref( blkid )));
	// return blk_itr != blk_state.end();
}

func (d *DispatcherManager) haveBlock(blkid crypto.Sha256) bool {
	// d.blkStateMtx.Lock()
	// d.blkStateMtx.Unlock()
	// index := d.blkState.

	// const auto& index = blk_state.get<by_block_id>();
	// auto blk_itr = index.find( blkid );
	// if( blk_itr != index.end() ) {
	//    return blk_itr->have_block;
	// }
	return false
}
