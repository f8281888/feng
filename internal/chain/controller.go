package chain

import (
	"bytes"
	controllerenum "feng/enum/controller-enum"
	"feng/internal/chainbase"
	"feng/internal/fc/crypto"
	"feng/internal/log"
)

//ControllerConfig ..
type ControllerConfig struct {
	senderBypassWhiteblacklist      []Name
	actorWhitelist                  []Name
	actorBlacklist                  []Name
	contractWhitelist               []Name
	contractBlacklist               []Name
	actionBlacklist                 map[Name]Name
	keyBlacklist                    []crypto.PublicKey
	blocksDir                       string
	stateDir                        string
	stateSize                       uint64
	stateGuardSize                  uint64
	reversibleCacheSize             uint64
	reversibleGuardSize             uint64
	sigCPUBillPct                   uint32
	threadPoolSize                  uint16
	readonly                        bool
	forceAllChecks                  bool
	disableReplayOpts               bool
	contractsConsole                bool
	allowRAMBillingInNotify         bool
	maximumVariableSignatureLength  uint32
	disableAllSubjectiveMitigations bool
	//wasm_interface::vm_type  wasm_runtime
	//eosvmoc::config          eosvmoc_config
	eosvmocTierup bool
	//enum class
	//DBReadMode ..
	DBReadMode           uint16
	inTrxRequiringChecks bool
	//microseconds
	subjectiveCPULeeway            uint64
	trustedProducerLightValidation bool
	snapshotHeadBlock              uint32
	//named_thread_pool              thread_pool;
	//platform_timer                 timer;
	// typedef pair<scope_name,action_name>                   handler_key;
	// map< account_name, map<handler_key, apply_handler> >   apply_handlers;
	// unordered_map< builtin_protocol_feature_t, std::function<void(controller_impl&)>, enum_hash<builtin_protocol_feature_t> > protocol_feature_activation_handlers;
}

//ControllerImpl ..
type ControllerImpl struct {
	DB                             chainbase.ChainBase
	reversibleBlocks               chainbase.ChainBase
	self                           *Controller
	ForkDB                         ForkDatabase
	conf                           ControllerConfig
	inTrxRequiringChecks           bool
	trustedProducerLightValidation bool
	snapshotHeadBlock              uint32
	ChainID                        CIDType
	readMode                       uint16
	head                           *BlockState
	Blog                           BlockLog
}

//Controller ..
type Controller struct {
	ControllerConfig
	ControllerImpl
}

//StartUp ..
func (c Controller) StartUp(b func() bool, genesis GenesisState) {
	if c.DB == (chainbase.ChainBase{}) {
		log.Assert("ChainBase is null")
	}

	if c.DB.GetRversion() > 1 {
		log.Assert("This version of controller.StartUp only works with a fresh state database.")
	}

	genesisChainID := genesis.ComputeChainID()
	if ok := bytes.Equal(c.ChainID.GetSha256().Hash, genesisChainID.GetSha256().Hash); !ok {
		log.Assert("genesis state provided to startup corresponds to a chain ID %s that does not match the chain ID that controller was constructed with %s)",
			string(genesisChainID.sha256.Hash), string(genesisChainID.sha256.Hash))
	}

	if c.ForkDB == (ForkDatabase{}) {
		log.Assert("ForkDB is null")
	}

	if c.ForkDB.Head() != nil {
		if c.readMode == controllerenum.Irreversible && c.ForkDB.Head() != c.ForkDB.Root() {
			c.ForkDB.RollbackHeadToRoot()
		}

		log.AppLog().Infof("No existing chain state. Initializing fresh blockchain state.")
	} else {
		log.AppLog().Infof("No existing chain state or fork database. Initializing fresh blockchain state and resetting fork database.")
	}

	c.InitializeBlockchainState(genesis)
	if c.ForkDB.Head() == nil {
		c.ForkDB.Reset(c.head)
	}

	// if c.Blog == (BlockLog{}) {
	// 	log.Assert("blog IS NULL")
	// }

	// if c.Blog.Head() == nil {
	// 	if c.Blog.myBlockLogImpl.firstBlockNum != 1 {
	// 		log.Assert("block log does not start with genesis block")
	// 	}
	// } else {
	// 	c.Blog.Reset(genesis, *c.head.block)
	// }

}

//Init ..
func (c *Controller) Init(func() bool) {
	var libNum uint32 = 0
	if c.Blog.Head() != nil {
		libNum = c.Blog.Head().BlockNum()
	} else {
		libNum = c.ForkDB.Root().BlockNum
	}

	//auto header_itr = validate_db_version( db );
	//...一些数据库操作
	if c.DB.GetRversion() > uint64(c.head.BlockNum) {
		log.AppLog().Infof("database revision %s is greater than head block number %s, ", c.DB.GetRversion(),
			c.head.BlockNum)
	}

	for {
		if c.DB.GetRversion() < uint64(c.head.BlockNum) {
			break
		}

		c.DB.Undo()
	}

	//protocol_features.init( db );
	println(libNum)
}

//InitializeBlockchainState ..
func (c *Controller) InitializeBlockchainState(genesis GenesisState) {
	log.AppLog().Infof("Initializing new blockchain with genesis state")
	// producer_authority_schedule initial_schedule = { 0, { producer_authority{config::system_account_name, block_signing_authority_v0{ 1, {{genesis.initial_key, 1}} } } } };
	// legacy::producer_schedule_type initial_legacy_schedule{ 0, {{config::system_account_name, genesis.initial_key}} };
	genheader := BlockHeaderState{}
	genheader.id = genheader.header.ID()
}

//StartUpSingle ..
func (c Controller) StartUpSingle(b func() bool) {

}

//GetReadMode ..
func (c Controller) GetReadMode() uint16 {
	return c.DBReadMode
}
