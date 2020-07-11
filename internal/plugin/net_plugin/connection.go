package netplugin

import (
	"feng/config"
	nodeenum "feng/enum/node-enum"
	"feng/internal/fc/common"
	"feng/internal/fc/crypto"
	"feng/internal/fc/io"
	"feng/internal/fc/network"
	"feng/internal/log"
	"net"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"
	"unsafe"
)

var unknown string
var messageBufLen uint32 = 1024 * 1024

//PeerSyncState ..
type PeerSyncState struct {
	StartBlock uint32
	EndBlock   uint32
	Last       uint32
	StartTime  time.Duration
}

//Connection ..
type Connection struct {
	peerRequested  PeerSyncState
	socketOpen     bool
	peerAddr       string
	connectionType byte
	//boost::asio::io_context::strand的主要作用,定义了事件处理程序的严格顺序调用
	//Strand                              context
	Conn                                net.Conn
	PendingMessageBuffer                network.MessageBuffer
	OutStandingReadBytes                int
	BufferQueued                        QueueBuffer
	TrxInProgressSize                   uint32
	ConnectionID                        uint32
	SentHandshakeCount                  int16
	Connecting                          bool
	Syncing                             bool
	ProtocolVersion                     uint16
	ConsecutiveRejectedBlocks           uint16
	ConsecutiveImmediateConnectionClose uint16
	ResponseExpectedTimeMtx             sync.Mutex
	ResponseExpectedTime                time.Ticker
	NoRetry                             int
	ConnMtx                             sync.Mutex
	LastReq                             RequestMessage
	LastHandshakeRecv                   HandshakeMessage
	LastHandshakeSent                   HandshakeMessage
	ForkHead                            crypto.Sha256
	ForkHeadNum                         uint32
	LastClose                           time.Time
	ConnNodeID                          crypto.Sha256
	RemoteEndpointIP                    string
	RemoteEndpointPort                  string
	LocalEndpointIP                     string
	LocalEndpointPort                   string
	myNetPlugin                         *NetPlugin
}

//SocketIsOpen ..
func (c *Connection) SocketIsOpen() bool {
	return c.socketOpen
}

//PeerAddress ..
func (c *Connection) PeerAddress() string {
	return c.peerAddr
}

//StartSession 开始处理业务
func (c *Connection) StartSession() bool {
	c.init()
	c.updateEndpoints()
	log.AppLog().Debugf("connected to %s", c.PeerName())
	c.socketOpen = true
	c.startReadMessage()
	return true
}

//
func (c *Connection) updateEndpoints() {
	rep := c.Conn.RemoteAddr()
	lep := c.Conn.LocalAddr()
	c.ConnMtx.Lock()
	defer c.ConnMtx.Unlock()
	c.RemoteEndpointIP = rep.String()
	c.RemoteEndpointPort = rep.Network()
	c.LocalEndpointIP = lep.String()
	c.LocalEndpointPort = lep.Network()

}

var messageHeaderSize = 4
var defSendBufferSizeMb = 4
var defSendBufferSize = 1024 * 1024 * defSendBufferSizeMb
var defMaxWriteQueueSize = defSendBufferSize * 10
var signedBlockWhich uint32 = 7

//startReadMessage 开始读
func (c *Connection) startReadMessage() {
	//读取并修改被封装的值 与atomic::exchange() 成员函数等价
	//std::atomic_exchange<decltype(outstanding_read_bytes.load())>( &outstanding_read_bytes, 0 );
	c.OutStandingReadBytes = 0
	minimumRead := c.OutStandingReadBytes
	if minimumRead == 0 {
		minimumRead = messageHeaderSize
	}

	if config.NodeConf.UseSocketReadWatermark {
		maxSocketReadWatermark := 4096
		socketReadWatermark := common.Min(minimumRead, maxSocketReadWatermark)
		//设置处理socket输入的最小的字节数
		//boost::asio::socket_base::receive_low_watermark read_watermark_opt(socket_read_watermark);
		//socket->set_option( read_watermark_opt, ec );
		//可以跳过
		println(socketReadWatermark)
	}

	//auto completion_handler = [minimum_read](boost::system::error_code ec, std::size_t bytes_transferred) -> std::size_t {};
	completionHandler := func(bytesTransferred int) int {
		if bytesTransferred >= minimumRead {
			return 0
		}

		return minimumRead - bytesTransferred
	}

	writeQueueSize := c.BufferQueued.WriteQueueSize()
	if writeQueueSize > uint32(defMaxWriteQueueSize) {
		log.AppLog().Errorf("write queue full %d bytes, giving up on connection, closing connection to:%s", writeQueueSize, c.PeerName())
		c.close(false, false)
	}

	bytesTransferred := 0
	go c.AsyncRead(completionHandler(bytesTransferred))
}

//AsyncRead 异步读 先写再处理读
func (c *Connection) AsyncRead(bytesTransferred int) {
	//async_read通常用户读取指定长度的数据，读完或出错才返回,读完之后才去处理handle
	//std::vector<boost::asio::mutable_buffer> async_read 第二个参数，把读到的数据写到里面，然后解析数据
	//async_read 最大究竟能读多少？ 读取用户指定的长度
	readBytes := make([]byte, 1024)
	var closeConnection bool
	_, err := c.Conn.Read(readBytes)
	if err != nil {
		log.AppLog().Errorf("Closing connection to: %s", c.PeerAddress())
		closeConnection = true
		c.Conn.Close()
	}

	c.PendingMessageBuffer.AsyncReadToWrite(readBytes)
	result := strings.Replace(string(readBytes), "\n", "", 1)
	log.AppLog().Debugf("AsyncRead read buf :%s", result)
	if !c.SocketIsOpen() {
		return
	}

	closeConnection = false
	if bytesTransferred > int(c.PendingMessageBuffer.BytesToWrite()) {
		log.AppLog().Errorf("async_read_some callback: bytes_transfered = %d, buffer.bytes_to_write = %d", bytesTransferred, c.PendingMessageBuffer.BytesToWrite())
		c.close(true, false)
		log.Assert("async_read_some err")
	}

	//上面AsyncReadToWrite 实现
	//c.PendingMessageBuffer.AdvanceWritePtr((uint32)(bytesTransferred))
	for c.PendingMessageBuffer.BytesToRead() > 0 {
		bytesInBuffer := c.PendingMessageBuffer.BytesToRead()
		if int(bytesInBuffer) < messageHeaderSize {
			c.OutStandingReadBytes = messageHeaderSize - int(bytesInBuffer)
			break
		} else {
			var messageLength uint32
			index := c.PendingMessageBuffer.ReadIndex()
			a := make([]byte, unsafe.Sizeof(messageLength))
			//从写中读取，先peek 一下有多少字节，C++ 里 用uint32 取地址传进去 void *就有点搞不懂了
			c.PendingMessageBuffer.Peek(&a, uint32(unsafe.Sizeof(messageLength)), &index)
			if a == nil {
				log.Assert("body is error")
			}

			messageLength = uint32(common.BytesToInt(a))
			if messageLength > uint32(defSendBufferSize*2) || messageLength == 0 {
				log.AppLog().Errorf("ncoming message length unexpected %d", messageLength)
				closeConnection = true
				break
			}

			totalMessageBytes := messageLength + uint32(messageHeaderSize)
			if bytesInBuffer >= totalMessageBytes {
				c.PendingMessageBuffer.AdvanceWritePtr(uint32(messageHeaderSize))
				c.ConsecutiveImmediateConnectionClose = 0
				if !c.ProcessNextMessage(messageLength) {
					return
				}
			} else {
				outstandingMessageBytes := totalMessageBytes - bytesInBuffer
				availableBufferBytes := c.PendingMessageBuffer.BytesToWrite()
				if outstandingMessageBytes > availableBufferBytes {
					c.PendingMessageBuffer.AddSpace(outstandingMessageBytes - availableBufferBytes)
				}

				c.OutStandingReadBytes = int(outstandingMessageBytes)
				break
			}
		}
	}

	closeConnection = true
	if !closeConnection {
		c.startReadMessage()
	}
}

//ProcessNextMessage ..
//这里的函数应该就是来处理，把包头读走后，读取包体的内容，包体需要解包
func (c *Connection) ProcessNextMessage(messageLength uint32) bool {
	//TODO 有点复杂
	peekDs := c.PendingMessageBuffer.CreatePeekDatastream()
	var which uint32
	//一开始进来，是first 0 ,second 4
	io.Unpack(&peekDs, which)
	if which == signedBlockWhich {
		// var bh chain.BlockHeader
		// io.Unpack(&peekDs, which)
		// blkID := bh.ID()
		// blkNum := bh.BlockNum()

		// if c.myNetPlugin.Dispatcher.haveBlock(blkID) {
		// 	log.AppLog().Infof("canceling wait on %s, already received block ${num}, id %d...", c.PeerName(), common.BytesToInt(blkID.Hash[8:16]))
		// 	c.myNetPlugin.SyncManager.syncrec
		// }
	}

	return true
}

func (c *Connection) close(recoonect, shutdown bool) {
	c.socketOpen = false
	//找不到这个用法
	// if( self->socket->is_open() ) {
	// 	self->socket->shutdown( tcp::socket::shutdown_both, ec );
	// 	self->socket->close( ec );
	//  }
	//socket 和这里的coon区别呢？
	//self->socket.reset( new tcp::socket( my_impl->thread_pool->get_executor() ) );
	c.flushQueue()
	c.Connecting = false
	c.Syncing = false
	c.ConsecutiveRejectedBlocks = 0
	c.ConsecutiveImmediateConnectionClose++
	hasLastReq := true
	c.ConnMtx.Lock()
	hasLastReq = c.LastReq.IsEmpty()
	c.LastHandshakeRecv = HandshakeMessage{}
	c.LastHandshakeSent = HandshakeMessage{}
	c.LastClose = time.Now()
	c.ConnNodeID = crypto.Sha256{}
	c.ConnMtx.Unlock()

	if hasLastReq && shutdown {
		c.myNetPlugin.Dispatcher.retryFetch(c)
	}

	c.peerRequested = PeerSyncState{}
	c.SentHandshakeCount = 0
	if !shutdown {
		c.myNetPlugin.SyncManager.SyncResetLibNum(c)
	}
}

func (c *Connection) flushQueue() {
	c.BufferQueued.ClearWriteQueue()
}

//PeerName ..
func (c *Connection) PeerName() string {
	c.ConnMtx.Lock()
	defer c.ConnMtx.Unlock()
	if c.LastHandshakeRecv.P2pAddress != "" {
		return c.LastHandshakeRecv.P2pAddress
	}

	if c.peerAddress() != "" {
		return c.PeerAddress()
	}

	if c.RemoteEndpointPort != unknown {
		return c.RemoteEndpointIP + ":" + c.RemoteEndpointPort
	}

	return "connection client"
}

func (c *Connection) peerAddress() string {
	return c.peerAddr
}

//TODO  const net_message& m 是一个结构体
func (c *Connection) enqueue(n NetMessage) {
	//跳过
}

//
func (c *Connection) fetchWait() {
	newC := c
	c.ResponseExpectedTimeMtx.Lock()
	defer c.ResponseExpectedTimeMtx.Unlock()
	go newC.fetchTimeout(c.myNetPlugin.RespExpectedPeriod)
	//定时器异步去执行超时处理 TODO 应该不是这么处理的
	// response_expected_timer.async_wait(
	// boost::asio::bind_executor( c->strand, [c]( boost::system::error_code ec ) {
	// 	c->fetch_timeout(ec);
	//  } ) );
}

func (c *Connection) fetchTimeout(d time.Duration) {
	// c.time.NewTimer(d)
	// c.dispatcher.retryFetch(c)
	println("fetchTimeout")
	return
}

//IsTransactionsOnlyConnection ..
func (c *Connection) IsTransactionsOnlyConnection() bool {
	return c.connectionType == nodeenum.BlocksOnly
}

//Connected ..
func (c *Connection) Connected() bool {
	return c.SocketIsOpen() && !c.Connecting
}

//Current ..
func (c *Connection) Current() bool {
	return c.Connected() && !c.Syncing
}

//RequestSyncBlocks ..
func (c *Connection) RequestSyncBlocks(start, end uint32) {
	s := syncRequestMessage{startBlock: start, endBlock: end}
	c.enqueue(s)
	c.SyncWait()
}

//SyncWait ..
func (c *Connection) SyncWait() {
	newC := c
	c.ResponseExpectedTimeMtx.Lock()
	defer c.ResponseExpectedTimeMtx.Unlock()
	go newC.SyncTimeOut()
}

//SyncTimeOut ..
func (c *Connection) SyncTimeOut() {
	c.myNetPlugin.SyncManager.syncReassignFetch(c, nodeenum.BenignOther)
}

//CancelSync ..
func (c *Connection) CancelSync(reason int) {
	log.AppLog().Debugf("cancel sync reason = %d, write queue size %d bytes peer %s", reason, c.BufferQueued.WriteQueueSize(), c.PeerName())
	c.CancelWait()
	c.FlushQueues()
	switch reason {
	case nodeenum.Validation:
	case nodeenum.FatalOther:
		c.NoRetry = reason
		g := goAwayMessage{reason: reason}
		c.enqueue(g)
		break
	default:
		log.AppLog().Infof("sending empty request but not calling sync wait on %s", c.PeerName())
		c.enqueue(syncRequestMessage{startBlock: 0, endBlock: 0})
	}
}

//CancelWait ..
func (c *Connection) CancelWait() {
	c.ResponseExpectedTimeMtx.Lock()
	c.ResponseExpectedTimeMtx.Unlock()
	c.ResponseExpectedTime.Stop()
}

//FlushQueues ..
func (c *Connection) FlushQueues() {
	c.BufferQueued.ClearWriteQueue()
}

var int16Max int16 = 32767

//SendHandshake ..
func (c *Connection) SendHandshake(force bool) {
	c.ConnMtx.Lock()
	defer c.ConnMtx.Unlock()
	if c.PopulateHandshake(&c.LastHandshakeSent, force) {
		//static_assert( std::is_same_v<decltype( c->sent_handshake_count ), int16_t>, "INT16_MAX based on int16_t" );
		if reflect.TypeOf(c.SentHandshakeCount).Name() != reflect.Int16.String() {
			log.Assert("INT16_MAX based on int16_t")
		}

		if c.SentHandshakeCount == int16Max {
			c.SentHandshakeCount = 1
		}

		c.LastHandshakeSent.generation = c.SentHandshakeCount + 1
		lastHandshakeSent := c.LastHandshakeSent
		log.AppLog().Infof("Sending handshake generation %d to %s, lib %d, head %d, id %s", lastHandshakeSent.generation, c.PeerName(), lastHandshakeSent.HeadNum, lastHandshakeSent.LastIrreversibleBlockNum, lastHandshakeSent.HeadID.String()[8:16])
		c.enqueue(lastHandshakeSent)
	}
}

var netVersionBase uint16 = 0x04b5
var netVersionRange uint16 = 106

//PopulateHandshake ..
func (c *Connection) PopulateHandshake(hello *HandshakeMessage, force bool) bool {
	var send bool = force
	hello.NetworkVersion = netVersionBase + netVersionRange
	prevHeadID := hello.HeadID
	var lib, head uint32
	lib, _, head, hello.LastIrreversibleBlockID, _, hello.HeadID = c.myNetPlugin.GetChainInfo()
	send = send || (lib != hello.LastIrreversibleBlockNum)
	send = send || (head != hello.HeadNum)
	send = send || prevHeadID.String() != hello.HeadID.String()
	if !send {
		return false
	}

	hello.LastIrreversibleBlockNum = lib
	hello.HeadNum = head
	hello.ChainID = c.myNetPlugin.ChainID
	hello.NodeID = c.myNetPlugin.NodeID
	hello.Key = c.myNetPlugin.GetAuthenticationKey()
	//sc::duration_cast<sc::nanoseconds>(sc::system_clock::now().time_since_epoch()).count();
	hello.Time = time.Now().Sub(time.Now())
	hello.Token = crypto.Sha256{}
	hello.Token.New(hello.Time.String())
	hello.Sig = c.myNetPlugin.SignCompact(hello.Key, hello.Token)
	if hello.Sig == (crypto.Signature{}) {
		hello.Token.New(hello.Time.String())
	}

	hello.P2pAddress = c.myNetPlugin.P2pAddress
	if c.IsTransactionsOnlyConnection() {
		hello.P2pAddress += ":trx"
	}

	if c.IsBlocksOnlyConnection() {
		hello.P2pAddress += ":blk"
	}

	hello.P2pAddress += " - " + hello.NodeID.String()[0:7]
	hello.os, _ = os.Hostname()
	hello.agent = c.myNetPlugin.UserAgentName
	return true
}

//IsBlocksOnlyConnection ..
func (c *Connection) IsBlocksOnlyConnection() bool {
	return c.connectionType == nodeenum.BlocksOnly
}

func (c *Connection) init() {
	c.PendingMessageBuffer.BufferLen = messageBufLen
	c.PendingMessageBuffer.InitIndex()
}
