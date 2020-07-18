package chainplugin

import (
	controllerenum "feng/enum/controller-enum"
	"feng/internal/app"
	"feng/internal/chain"
	"feng/internal/chainbase"
	"feng/internal/fc/crypto"
	"feng/internal/log"
	chaininterface "feng/internal/plugin/chain_interface"
)

//ChainPlugin ..
type ChainPlugin struct {
	chainPluginImpl
}

type chainPluginImpl struct {
	blocksDir         string
	readonly          bool
	loadedCheckpoints map[uint32]crypto.Sha256
	//AcceptTransactions..
	AcceptTransactions    bool
	apiAcceptTransactions bool
	forkDB                *chain.ForkDatabase
	//fc::optional<block_log>          block_logger;
	chainConfig chain.ControllerConfig
	chain       chain.Controller
	genesis     *chain.GenesisState
	//fc::optional<vm_type>            wasm_runtime
	//microseconds
	abiSerializerMaxTimeUs uint32
	//bfs::path
	snapshotPath            string
	preAcceptedBlockChannel chan chaininterface.PreAcceptedBlock
}

func init() {
	chainPlugin := &ChainPlugin{}
	chainPlugin.chain = chain.Controller{}
	app.App().RegisterPlugin("ChainPlugin", chainPlugin)
}

//Initialize ..
func (a *ChainPlugin) Initialize() {
	if a == nil {
		a = &ChainPlugin{}
	}

	log.AppLog().Infof("initializing chain plugin")
	gs := chain.GenesisState{}
	a.chainConfig = chain.ControllerConfig{}
	a.genesis = &gs
	a.chain = chain.Controller{}
	a.chainConfig = chain.ControllerConfig{}
	a.chain.DB = chainbase.DataBase{}
	a.chain.DB.SetRversion(1)
	a.chain.ChainID = chain.CIDType{}
	s := &crypto.Sha256{}
	s.Hash = []byte("1")
	a.chain.ChainID.SetCIDType(s)
	a.chain.ForkDB = new(chain.ForkDatabase)
	fork := &chain.ForkDatabaseImpl{}
	fork.DataDir = "/data"
}

//HandleSighup ..
func (a *ChainPlugin) HandleSighup() {
	println("ChainPlugin HandleSighup")
}

//StartUp ..
func (a *ChainPlugin) StartUp() {
	println("ChainPlugin StartUp")
}

//PluginStartUp ..
func (a *ChainPlugin) PluginStartUp() {
	println("ChainPlugin PluginStartUp")
	if a.chainConfig.DBReadMode == controllerenum.Irreversible || a.getAcceptTransactions() {
		log.Assert("read-mode = irreversible. transactions should not be enabled by enable_accept_transactions")
	}

	shotdown := func() bool {
		return app.App().IsQuiting()
	}

	if a.snapshotPath != "" {
		//IstreamSnapshotReader 从快照里面读取，先跳过
		//i := chain.IstreamSnapshotReader{}
		null := chain.GenesisState{}
		a.chain.StartUp(shotdown, null)
	} else if a.genesis != nil {
		a.chain.StartUp(shotdown, *a.genesis)
	} else {
		a.chain.StartUpSingle(shotdown)
	}
}

func (a *ChainPlugin) getAcceptTransactions() bool {
	return a.AcceptTransactions
}

//GetChainID ..
func (a ChainPlugin) GetChainID() crypto.Sha256 {
	return crypto.Sha256{}
}

//GetChain ..
func (a ChainPlugin) GetChain() chain.Controller {
	return a.chain
}

//EnableAcceptTransactions ..
func (a ChainPlugin) EnableAcceptTransactions() {
	println("EnableAcceptTransactions")
}
