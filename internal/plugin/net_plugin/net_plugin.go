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
	Coon                               net.Conn
	CurrentConnectionID                atomic.Uint32
	SyncMaster                         *SyncMaster
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
	a.SyncMaster.New(config.NodeConf.SyncFetchSpan)
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

	a.Listener, err = net.ListenTCP("tcp", laddr)
	if err != nil {
		log.AppLog().Error("prot:%s", laddr.Port)
		log.Assert("net_plugin::plugin_startup failed to bind to port")
	} else {
		a.Coon, err = a.Listener.Accept()
		if err != nil {
			log.Assert("can't accept")
		} else {
			a.StartListenLoop()
		}
	}
}

//StartListenLoop ..
func (a *NetPlugin) StartListenLoop() {
	println("startListenLoop")
}
