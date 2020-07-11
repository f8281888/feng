package netplugin

import (
	"feng/config"
	controllerenum "feng/enum/controller-enum"
	nodeenum "feng/enum/node-enum"
	"feng/internal/app"
	"feng/internal/fc/crypto"
	"feng/internal/log"
	chainplugin "feng/internal/plugin/chain_plugin"
	producerplugin "feng/internal/plugin/producer_plugin"
	"feng/internal/pool"
	"net"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"go.uber.org/atomic"
)

//ConnectionStauts ..
type ConnectionStauts struct {
	peer          string
	connecting    bool
	syncing       bool
	lastHandshake HandshakeMessage
}

//netPluginImpl ..
type netPluginImpl struct {
	//直接跳过accept,用来监听
	Listener                           *net.TCPListener
	CurrentConnectionID                atomic.Uint32
	SyncManager                        *SyncManager
	Dispatcher                         *DispatcherManager
	P2pAddress                         string
	P2pServerAddress                   string
	SuppliedPeers                      []string
	AllowedPeers                       []crypto.PublicKey
	PrivateKeys                        map[crypto.PublicKey]crypto.PrivateKey
	AllowedConnections                 byte
	ConnectorPeriod                    time.Duration
	TxnExpPeriod                       time.Duration
	RespExpectedPeriod                 time.Duration
	KeepaliveInterval                  time.Duration
	MaxCleanupTimeMs                   int32
	MaxClientCount                     uint32
	MaxNodesPerHost                    uint32
	P2pAcceptTransactions              bool
	PeerAuthenticationInterval         time.Duration
	ChainID                            crypto.Sha256
	NodeID                             crypto.Sha256
	UserAgentName                      string
	ChainPlugin                        *chainplugin.ChainPlugin
	ProducerPlugin                     *producerplugin.ProducerPlugin
	UseSocketReadWatermark             bool
	ConnectionsMtx                     sync.Mutex
	Connections                        []*Connection
	ConnectorCheckTimerMtx             sync.Mutex
	ConnectorCheckTimer                *time.Timer
	ConnectorChecksInFlight            int32
	ExpireTimerMtx                     sync.Mutex
	ExpireTimer                        *time.Timer
	KeepaliveTimerMtx                  sync.Mutex
	KeepaliveTimer                     *time.Timer
	InShutdown                         atomic.Bool
	IncomingTransactionAckSubscription chan interface{}
	ThreadPoolSize                     uint16
	ThreadPool                         *pool.WorkPool
	chainInfoMtx                       sync.Mutex
	chainLibNum                        uint32
	chainHeadBlkNum                    uint32
	chainForkHeadBlkNum                uint32
	chainLibID                         crypto.Sha256
	chainHeadBlkID                     crypto.Sha256
	chainForkHeadBlkID                 crypto.Sha256
}

//GetChainInfo ..
func (n *netPluginImpl) GetChainInfo() (uint32, uint32, uint32, crypto.Sha256, crypto.Sha256, crypto.Sha256) {
	n.chainInfoMtx.Lock()
	defer n.chainInfoMtx.Unlock()
	return n.chainLibNum, n.chainHeadBlkNum, n.chainForkHeadBlkNum, n.chainLibID, n.chainHeadBlkID, n.chainForkHeadBlkID
}

//GetAuthenticationKey ..
func (n *netPluginImpl) GetAuthenticationKey() crypto.PublicKey {
	if len(n.PrivateKeys) != 0 {
		for k, v := range n.PrivateKeys {
			println(reflect.TypeOf(v).Name())
			return k
		}
	}

	return crypto.PublicKey{}
}

//SignCompact ..
func (n *netPluginImpl) SignCompact(signer crypto.PublicKey, digest crypto.Sha256) crypto.Signature {
	// privateKeyItr, ok := n.PrivateKeys[signer]
	// if !ok {
	// 	privateKeyItr[signer].
	// }
	//签名的东西先跳过
	// 	auto private_key_itr = private_keys.find(signer);
	// 	if(private_key_itr != private_keys.end())
	// 	   return private_key_itr->second.sign(digest);
	// 	if(producer_plug != nullptr && producer_plug->get_state() == abstract_plugin::started)
	// 	   return producer_plug->sign_compact(signer, digest);
	// 	return chain::signature_type();
	return crypto.Signature{}
}

//NetPlugin ..
type NetPlugin struct {
	netPluginImpl
}

func init() {
	netPlugin := &NetPlugin{}
	app.App().RegisterPlugin("NetPlugin", netPlugin)
}

var peerLogFormat string
var defTxnExpireWait = time.Second * 3
var defRespExpectedWait = time.Second * 5
var maxP2pAddressLength int = 258
var maxHandshakeStrLength int = 384

//Initialize ..
func (a *NetPlugin) Initialize() {
	//从配置里面设置一些初始值
	println("NetPlugin Initialize")
	peerLogFormat = config.NodeConf.PeerLogFormat
	a.SyncManager.New(config.NodeConf.SyncFetchSpan)
	a.ConnectorPeriod = time.Duration(config.NodeConf.ConnectionCleanupPeriod)
	println(a.ConnectorPeriod)
	a.MaxCleanupTimeMs = config.NodeConf.MaxCleanupTimeMsec
	a.TxnExpPeriod = defTxnExpireWait
	a.RespExpectedPeriod = defRespExpectedWait
	a.MaxClientCount = config.NodeConf.MaxClients
	a.MaxNodesPerHost = config.NodeConf.P2pMaxNodesPerHost
	a.P2pAcceptTransactions = config.NodeConf.P2pAcceptTransactions
	a.UseSocketReadWatermark = config.NodeConf.UseSocketReadWatermark
	a.P2pAddress = config.NodeConf.P2pListenEndpoint
	if a.P2pAddress != "" && len(a.P2pAddress) > maxP2pAddressLength {
		log.Assert("p2p-listen-endpoint to long, must be less than %s", string(maxP2pAddressLength))
	}

	a.P2pServerAddress = config.NodeConf.P2pServerAddress
	if a.P2pServerAddress != "" && len(a.P2pServerAddress) > maxP2pAddressLength {
		log.Assert("p2p_server_address to long, must be less than %s", string(maxP2pAddressLength))
	}

	a.ThreadPoolSize = config.NodeConf.NetThreads
	if a.ThreadPoolSize <= 0 {
		log.Assert("net-threads %s must be greater than 0", string(a.ThreadPoolSize))
	}

	if len(config.NodeConf.P2pPeerAddress) > 0 {
		a.SuppliedPeers = config.NodeConf.P2pPeerAddress
	}

	if config.NodeConf.AgentName != "" {
		a.UserAgentName = config.NodeConf.AgentName
		if len(config.NodeConf.AgentName) > maxHandshakeStrLength {
			log.Assert("agent-name to long, must be less than %s", string(config.NodeConf.AgentName))
		}
	}

	if len(config.NodeConf.AllowedConnection) > 0 {
		for _, allowedRemote := range config.NodeConf.AllowedConnection {
			if allowedRemote == "any" {
				a.AllowedConnections |= nodeenum.Any
			} else if allowedRemote == "producers" {
				a.AllowedConnections |= nodeenum.Producers
			} else if allowedRemote == "specified" {
				a.AllowedConnections |= nodeenum.Specified
			} else if allowedRemote == "none" {
				a.AllowedConnections |= nodeenum.None
			}
		}
	}

	if a.AllowedConnections&nodeenum.Specified == byte(0) {
		log.Assert("At least one peer-key must accompany 'allowed-connection=specified'")
	}

	if len(config.NodeConf.PeerKey) > 0 {
		for _, keyString := range config.NodeConf.PeerKey {
			a.AllowedPeers = append(a.AllowedPeers, crypto.NewPublicKey(keyString))
		}
	}

	if len(config.NodeConf.PeerPrivateKey) > 0 {
		for first, second := range config.NodeConf.PeerPrivateKey {
			a.PrivateKeys[crypto.NewPublicKey(first)] = crypto.NewPrivateKey(second)
		}
	}

	a.ChainPlugin = app.App().FindPlugin("ChainPlugin").(*chainplugin.ChainPlugin)
	if a.ChainPlugin == nil {
		log.Assert("ChainPlugin is nil")
	}

	a.ChainID = a.ChainPlugin.GetChainID()
	crypto.RandPseudoBytes(&a.NodeID.Hash, len(a.NodeID.Hash))
	cc := a.ChainPlugin.GetChain()
	if cc.GetReadMode() == controllerenum.Irreversible || cc.GetReadMode() == controllerenum.ReadOnly {
		if a.P2pAcceptTransactions {
			a.P2pAcceptTransactions = false
			var m string
			if cc.GetReadMode() == controllerenum.Irreversible {
				m = "irreversible"
			} else {
				m = "read-only"
			}
			log.AppLog().Infof("p2p-accept-transactions set to false due to read-mode:%s", m)
		}
	}

	if a.P2pAcceptTransactions {
		a.ChainPlugin.EnableAcceptTransactions()
	}

}

//HandleSighup ..
func (a *NetPlugin) HandleSighup() {
	println("NetPlugin HandleSighup")
}

//StartUp ..
func (a *NetPlugin) StartUp() {
	println("NetPlugin StartUp")
}

//PluginStartUp ..
func (a *NetPlugin) PluginStartUp() {
	println("NetPlugin PluginStartUp")
	log.AppLog().Infof("my node_id is %d", a.NodeID)
	a.ProducerPlugin = app.App().FindPlugin("ProducerPlugin").(*producerplugin.ProducerPlugin)
	a.ThreadPool = pool.NewPool(int(a.ThreadPoolSize))
	a.Dispatcher = DispatcherReset(a.ThreadPool.GetExecutor())
	if !a.P2pAcceptTransactions && len(a.P2pAddress) > 0 {
		print(`
		"***********************************\n"
		"* p2p-accept-transactions = false *\n"
		"*    Transactions not forwarded   *\n"
		"***********************************\n" `)
		println("")
	}

	var laddr *net.TCPAddr
	var err error
	if len(a.P2pAddress) > 0 {
		laddr, err = net.ResolveTCPAddr("tcp", a.P2pAddress)
		if len(a.P2pServerAddress) > 0 {
			a.P2pAddress = a.P2pServerAddress
		} else {
			if err != nil {
				host, _ := os.Hostname()
				if host == "" {
					log.Assert("can't find hostname")
				}

				arr := strings.Split(a.P2pAddress, ":")
				a.P2pAddress = host + ":" + arr[1]
			}
		}
	}

	go func() {
		a.Listener, err = net.ListenTCP("tcp", laddr)
		if err != nil {
			log.AppLog().Error("prot:%s", laddr.Port)
			log.Assert("net_plugin::plugin_startup failed to bind to port")
		} else {
			a.StartListenLoop()
		}
	}()
}

//StartListenLoop ..
func (a *NetPlugin) StartListenLoop() {
	println("startListenLoop")
	newConnection := new(Connection)
	newConnection.Connecting = true
	//new_connection->strand.post( [this, new_connection = std::move( new_connection )](){
	//里面去处理accept，然后递归调用自己，这里的strand可能是要保证执行顺序
	//}
	go func() {
		var err error
		newConnection.Conn, err = a.Listener.Accept()
		newConnection.myNetPlugin = a
		//跟boost的区别就是绑定一个异步的执行函数去执行客户端的请求
		// acceptor->async_accept( *new_connection->socket,
		// 	boost::asio::bind_executor( new_connection->strand, [new_connection, socket=new_connection->socket, this]( boost::system::error_code ec ) {}
		if err != nil {
			log.AppLog().Error("Error accepting connection:%s", err.Error())
		} else {
			go a.handAccept(newConnection)
			a.StartListenLoop()
		}
	}()

	//a.Coon.Read(b)
}

func (a *NetPlugin) handAccept(conn *Connection) {
	visitors := 0
	fromAddr := 0
	paddrADD := conn.Conn.RemoteAddr()
	paddrStr := paddrADD.String()
	f := func() {
		//for_each_connection( [&visitors, &from_addr, &paddr_str]( auto& conn ) {}
		//不用像C++那样，捕获参数还要用[]去捕获,入参也不用，比较方便，比较不方便的是不能立即执行，还要f()去调用一下
		if conn.SocketIsOpen() {
			if conn.PeerAddress() != "" {
				visitors++
				conn.ConnMtx.Lock()
				if paddrStr == conn.RemoteEndpointIP {
					fromAddr++
				}
				conn.ConnMtx.Unlock()
			}
		}
	}

	f()

	if fromAddr < int(a.MaxNodesPerHost) && (a.MaxClientCount == 0 || visitors < int(a.MaxClientCount)) {
		log.AppLog().Infof("Accepted new connection:%s", paddrStr)
		if conn.StartSession() {
			a.ConnectionsMtx.Lock()
			defer a.ConnectionsMtx.Unlock()
			a.Connections = append(a.Connections, conn)
		}
	} else {
		if fromAddr >= int(a.MaxNodesPerHost) {
			log.AppLog().Debugf("Number of connections %d from %s exceeds limit %d", fromAddr, paddrStr, a.MaxNodesPerHost)
		} else {
			log.AppLog().Debugf("ax_client_count %d exceeded", a.MaxClientCount)
		}

		conn.close(true, false)
	}
}

//按顺序执行函数
// template<typename Function>
// void for_each_connection( Function f ) {
//    std::shared_lock<std::shared_mutex> g( my_impl->connections_mtx );
//    for( auto& c : my_impl->connections ) {
// 	  if( !f( c ) ) return;
//    }
// }
