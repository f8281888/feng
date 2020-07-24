package chain

import (
	"bytes"
	controllerenum "feng/enum/controller-enum"
	"feng/internal/chainbase"
	"feng/internal/fc/crypto"
	"feng/internal/log"
	"time"
)

const (
	//Full ..
	Full = iota
	//Light ..
	Light = iota
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
	ReversibleGuardSize             uint64
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
	blockValidationMode            uint32
	trustedProducers               []AccountName
	//named_thread_pool              thread_pool;
	//platform_timer                 timer;
	// typedef pair<scope_name,action_name>                   handler_key;
	// map< account_name, map<handler_key, apply_handler> >   apply_handlers;
	// unordered_map< builtin_protocol_feature_t, std::function<void(controller_impl&)>, enum_hash<builtin_protocol_feature_t> > protocol_feature_activation_handlers;
}

//MaybeSession ..
type MaybeSession struct {
	session chainbase.Session
}

//BuildingBlock ..
type BuildingBlock struct {
	PendingTrxMetas                         []*TransactionMetadata
	PendingBlockHeaderState                 PendingBlockHeaderState
	NewPendingProducerSchedule              ProducerAuthoritySchedule
	NewProtocolFeatureActivations           []DigestType
	NumNewProtocolFeaturesThatHaveActivated uint32
	PendingTrxReceipts                      []TransactionReceipt
	Actions                                 []ActionReceipt
	TransactionMroot                        Checksum256Type
}

//AssembledBlock ..
type AssembledBlock struct {
	ID                      BlockIDType
	TrxMetas                []*TransactionMetadata
	PendingBlockHeaderState PendingBlockHeaderState
	UnsignedBlock           *SignedBlock
}

//CompledBlock ..
type CompledBlock struct {
	blockState *BlockState
}

//BlockStageType ..
type BlockStageType struct {
	buildingBlock   *BuildingBlock
	assembledBlock  *AssembledBlock
	compledBlock    *CompledBlock
	blockStatus     uint32
	producerBlockID BlockIDType
}

//PendingState ..
type PendingState struct {
	dbSession  MaybeSession
	blockState BlockStageType
}

//GetPendingBlockHeaderState ..
func (p PendingState) GetPendingBlockHeaderState() PendingBlockHeaderState {
	if p.blockState.buildingBlock != nil {
		return p.blockState.buildingBlock.PendingBlockHeaderState
	}

	return p.blockState.assembledBlock.PendingBlockHeaderState
}

//ExtractTrxMetas ..
func (p *PendingState) ExtractTrxMetas() []*TransactionMetadata {
	if p.blockState.buildingBlock != nil {
		if len(p.blockState.buildingBlock.PendingTrxMetas) > 0 {
			return p.blockState.buildingBlock.PendingTrxMetas
		}
	}

	if p.blockState.assembledBlock != nil {
		if len(p.blockState.assembledBlock.TrxMetas) > 0 {
			return p.blockState.assembledBlock.TrxMetas
		}
	}

	return p.blockState.compledBlock.blockState.ExtractTrxsMetas()
}

//ControllerImpl ..
type ControllerImpl struct {
	DB                             chainbase.DataBase
	reversibleBlocks               chainbase.DataBase
	self                           *Controller
	ForkDB                         *ForkDatabase
	conf                           ControllerConfig
	inTrxRequiringChecks           bool
	TrustedProducerLightValidation bool
	snapshotHeadBlock              uint32
	ChainID                        CIDType
	readMode                       uint16
	head                           *BlockState
	Blog                           BlockLog
	pending                        *PendingState
	protocolFeatures               ProtocolFeatureManager
	preAcceptedBlock               chan *SignedBlock
	acceptedBlockHeader            chan *BlockState
	resourceLimits                 ResourceLimitsManager
}

//Controller ..
type Controller struct {
	ControllerConfig
	ControllerImpl
}

//StartUp ..
func (c Controller) StartUp(b func() bool, genesis GenesisState) {
	if c.DB == (chainbase.DataBase{}) {
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

	if c.ForkDB == nil {
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
	var producers []ProducerAuthority
	producer := Name{}
	producer.Name(systemAccountName)
	log.AppLog().Debugf("producer :%s", producer.ToString())
	producers = append(producers, ProducerAuthority{ProducerName: producer})
	initialSchedule := ProducerAuthoritySchedule{version: 0, Producers: producers}
	producerkey := ProducerKey{producerName: producer, blockSigningKey: genesis.initialKey}
	var producerKeys []ProducerKey
	producerKeys = append(producerKeys, producerkey)
	initialLegacySchedule := ProducerScheduleType{version: 0, producers: producerKeys}
	genheader := BlockHeaderState{}
	genheader.ActiveSchedule = initialSchedule
	genheader.pendingSchedule.schedule = initialSchedule
	genheader.pendingSchedule.scheduleHash = crypto.Sha256{Data: initialLegacySchedule}
	genheader.Header.Timestamp = BlockTimestamp{Slot: uint32(genesis.initialTimestamp)}
	genheader.Header.actionMroot = *genesis.ComputeChainID().GetSha256()
	genheader.ID = &crypto.Sha256{}
	*genheader.ID = genheader.Header.ID()
	genheader.BlockNum = genheader.Header.BlockNum()
	c.head = &BlockState{}
	c.head.Copy(genheader)
	c.head.activatedProtocolFeatures = ProtocolFeatureActivationSet{}
	c.head.Block = &SignedBlock{}
	c.head.Block.Copy(genheader.Header)
	c.DB.SetRversion(uint64(c.head.BlockNum))
	c.initializeDatabase(&genesis)
}

//TODO
func (c Controller) initializeDatabase(genesis *GenesisState) {

}

//StartUpSingle ..
func (c Controller) StartUpSingle(b func() bool) {

}

//GetReadMode ..
func (c Controller) GetReadMode() uint16 {
	return c.DBReadMode
}

//SetSubjectiveCPULeeway ..
func (c *Controller) SetSubjectiveCPULeeway(t time.Duration) {
	c.subjectiveCPULeeway = uint64(t)
}

//FetchBlockByID ..
func (c *Controller) FetchBlockByID(id BlockIDType) *SignedBlock {
	state := c.ForkDB.GetBlock(id)
	if state != nil && state.Block != nil {
		return state.Block
	}

	return nil
}

//FetchBlockByNumer ..
func (c *Controller) FetchBlockByNumer(blockNum uint32) *BlockState {
	return nil
}

//FetchBlockStateByNumer ..
func (c *Controller) FetchBlockStateByNumer(blockNum uint32) *BlockState {
	//revBlocks := c.reversibleBlocks
	return nil
}

//CreateBlockStateFuture ..
func (c *Controller) CreateBlockStateFuture(b *SignedBlock) *BlockState {
	//revBlocks := c.reversibleBlocks
	return nil
}

//AbortBlock ..
func (c *Controller) AbortBlock() []*TransactionMetadata {
	var appliedTrxs []*TransactionMetadata
	if c.pending != nil {
		appliedTrxs = c.pending.ExtractTrxMetas()
		//pending.reset();
		c.protocolFeatures.PoppedBlocksTo(c.head.BlockNum)
	}

	return appliedTrxs
}

//HeadBlockState ..
func (c Controller) HeadBlockState() *BlockState {
	return c.head
}

//HeadBlockTime ..
func (c Controller) HeadBlockTime() uint64 {
	return c.head.Header.Timestamp.Time()
}

//ForkedBranchCallback ..
type ForkedBranchCallback func(*BranchType)

//TrxMetaCacheLookup ..
type TrxMetaCacheLookup = func(*TransactionIDType) TransactionMetadata

//GetDB ..
func (c Controller) GetDB() chainbase.DataBase {
	return c.DB
}

func (c *Controller) validateDBAvailableSize() {
	//get_free_memory() 获取段内存，用来？
	//const auto free = my->reversible_blocks.get_segment_manager()->get_free_memory();
	free := c.GetDB().GetSegmentManager().Len()
	guard := c.ControllerImpl.conf.ReversibleGuardSize
	if uint64(free) < guard {
		log.Assert("database free: %d, guard size: %d", free, guard)
	}
}

func (c *Controller) validateReversibleAvailableSize() {
	free := c.reversibleBlocks.GetSegmentManager().Len()
	guard := c.conf.ReversibleGuardSize
	if uint64(free) < guard {
		log.Assert("reversible_guard_exception free: %d, guard size: %d", free, guard)
	}
}

//PushBlock ..
func (c *Controller) PushBlock(blockState *BlockState, forkedBranchCb ForkedBranchCallback, trxLookup TrxMetaCacheLookup) {
	c.validateDBAvailableSize()
	c.validateReversibleAvailableSize()
	c.pushBlock(blockState, forkedBranchCb, trxLookup)
}

const (
	//Irreversible ..
	Irreversible = iota
	//Validated ..
	Validated = iota
	//Complete ..
	Complete = iota
	//Incomplete ..
	Incomplete = iota
)

//PushBlock ..
func (c *Controller) pushBlock(blockState *BlockState, forkedBranchCb ForkedBranchCallback, trxLookup TrxMetaCacheLookup) {
	var s uint32 = Complete
	if c.pending == nil {
		log.Assert("it is not valid to push a block when there is a pending block")
	}

	//oldValue := c.TrustedProducerLightValidation
	// resetProdLightValidation := common.MakeScopeExit(func() {
	// 	c.TrustedProducerLightValidation = oldValue
	// })

	bsp := blockState
	b := bsp.Block
	c.emitSingedBlock(b)
	c.ForkDB.Add(bsp, false)

	if c.isTrustedProducer(b.Producer) {
		c.trustedProducerLightValidation = true
	}

	c.emitBlockState(bsp)

	if c.readMode != Irreversible {
		c.maybeSwitchForks(c.ForkDB.pendingHead(), s, forkedBranchCb, trxLookup)
	}

}

//TOD
func (c *Controller) maybeSwitchForks(newHead *BlockState, s uint32, forkedBranchCb ForkedBranchCallback, trxLoopup TrxMetaCacheLookup) {
	//headChanged := true
	if &newHead.Header.previous == c.head.ID {
		c.applyBlock(newHead, s, trxLoopup)
	}
}

//TODO
func (c *Controller) applyBlock(bsp *BlockState, s uint32, trxLoop TrxMetaCacheLookup) {
	b := bsp.Block
	newProtocolFeatureActivations := bsp.GetNewProtocolFeatureActivations()
	producerBlockID := b.ID()
	c.startBlock(b.Timestamp, b.Confirmed, newProtocolFeatureActivations, s, producerBlockID)
}

//TODO
func (c *Controller) startBlock(when BlockTimestamp, confirmBlockCount uint16, newProtocolFeatureActivations []DigestType, s uint32, producerBlockID BlockIDType) {
	if c.pending == nil {
		log.Assert("pending block already exists")
	}
}

//监听发出信号，有信号过来就触发
func (c *Controller) emitSingedBlock(k *SignedBlock) {
	c.preAcceptedBlock = make(chan *SignedBlock)
	c.preAcceptedBlock <- k
}

func (c *Controller) emitBlockState(k *BlockState) {
	c.acceptedBlockHeader = make(chan *BlockState)
	c.acceptedBlockHeader <- k
}

func (c Controller) isTrustedProducer(producer AccountName) bool {
	var ok bool = false
	for _, i := range c.conf.trustedProducers {
		if producer == i {
			ok = true
			break
		}
	}

	return c.getValidationMode() == Light || ok
}

func (c Controller) getValidationMode() uint32 {
	return c.conf.blockValidationMode
}

//LastIrreversibleBlockNum ..
func (c Controller) LastIrreversibleBlockNum() uint32 {
	return c.ForkDB.root.BlockNum
}

//IsBuildingBlock ..
func (c Controller) IsBuildingBlock() bool {
	return c.pending != nil
}

//PendingBlockTime ..
func (c Controller) PendingBlockTime() time.Duration {
	if c.pending == nil {
		log.Assert("no pending block")
	}

	if c.pending.blockState.compledBlock != nil {
		return time.Duration(c.pending.blockState.compledBlock.blockState.Header.Timestamp.Slot)
	}

	return time.Duration(c.pending.GetPendingBlockHeaderState().Timestamp.Slot)
}

//HeadBlockNum ..
func (c Controller) HeadBlockNum() uint32 {
	return c.head.BlockNum
}

//GetResourceLimitsManager ..
func (c Controller) GetResourceLimitsManager() ResourceLimitsManager {
	return c.resourceLimits
}
