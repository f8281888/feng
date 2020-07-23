package producerplugin

import (
	"encoding/json"
	"feng/config"
	controllerenum "feng/enum/controller-enum"
	nodeenum "feng/enum/node-enum"
	"feng/internal/app"
	"feng/internal/chain"
	"feng/internal/fc/common"
	"feng/internal/fc/crypto"
	"feng/internal/log"
	chainplugin "feng/internal/plugin/chain_plugin"
	"feng/internal/pool"
	"fmt"
	"os"
	"strings"
	"time"
)

// //多键索引容器
// using transaction_id_with_expiry_index = multi_index_container<
//    transaction_id_with_expiry,
//    indexed_by<
//       hashed_unique<tag<by_id>, BOOST_MULTI_INDEX_MEMBER(transaction_id_with_expiry, transaction_id_type, trx_id)>,
//       ordered_non_unique<tag<by_expiry>, BOOST_MULTI_INDEX_MEMBER(transaction_id_with_expiry, fc::time_point, expiry)>
//    >
// >;

//TransactionIDWithExpiry ..
type TransactionIDWithExpiry struct {
	TrxID  *chain.TransactionIDType
	Expiry time.Duration
}

//TransactionIDWithExpiryIndex ..
type TransactionIDWithExpiryIndex = map[TransactionIDWithExpiry]TransactionIDWithExpiry

// BuildTransactionIndex 构建查询索引
func BuildTransactionIndex(t TransactionIDWithExpiryIndex, list []TransactionIDWithExpiry) {
	// 遍历所有数据
	for _, profile := range list {
		// 构建查询键
		key := TransactionIDWithExpiry{
			TrxID:  profile.TrxID,
			Expiry: profile.Expiry,
		}

		// 保存查询键
		t[key] = profile
	}
}

// QueryTransactionData 根据条件查询数据
func QueryTransactionData(t TransactionIDWithExpiryIndex, trxID chain.TransactionIDType, expiry time.Duration) {

	// 根据查询条件构建查询键
	key := TransactionIDWithExpiry{&trxID, expiry}

	// 根据键值查询数据
	result, ok := t[key]

	// 找到数据打印出来
	if ok {
		fmt.Println(result)
	} else {
		fmt.Println("no found")
	}
}

//SnapshotInformation ..
type SnapshotInformation struct {
	HeadBlockID  chain.BlockIDType
	SnapshotName string
}

// template<typename T>
// using next_function = std::function<void(const fc::static_variant<fc::exception_ptr, T>&)>;
//using next_t = producer_plugin::next_function<producer_plugin::snapshot_information>;

//NextType ..
type NextType func(s SnapshotInformation)

//PendingSnapshot ..
type PendingSnapshot struct {
	BlockID      chain.BlockIDType
	next         NextType
	PendgingPath string
	FinalPath    string
}

//GetHeight ..
func (p PendingSnapshot) GetHeight() uint32 {
	return uint32(common.BytesToInt(p.BlockID.Hash))
}

//SnapshotIndex ..
type SnapshotIndex struct {
	BlockID *chain.BlockIDType
	Height  uint32
}

//PendingSnapshotIndexIndex ..
type PendingSnapshotIndexIndex = map[SnapshotIndex]PendingSnapshot

// BuildSnapshotIndex 构建查询索引
func BuildSnapshotIndex(t PendingSnapshotIndexIndex, list []PendingSnapshot) {
	// 遍历所有数据
	for _, profile := range list {
		// 构建查询键
		key := SnapshotIndex{
			BlockID: &profile.BlockID,
			Height:  profile.GetHeight(),
		}

		// 保存查询键
		t[key] = profile
	}
}

// QuerySnapshotData 根据条件查询数据
func QuerySnapshotData(t PendingSnapshotIndexIndex, blkID chain.BlockIDType, height uint32) {

	// 根据查询条件构建查询键
	key := SnapshotIndex{&blkID, height}

	// 根据键值查询数据
	result, ok := t[key]

	// 找到数据打印出来
	if ok {
		fmt.Println(result)
	} else {
		fmt.Println("no found")
	}
}

//SignatureProviderType ..
type SignatureProviderType = func(digist chain.DigestType) chain.SignatureType

//Impl ..
type Impl struct {
	ProductionEnabled                     bool
	PauseProduction                       bool
	SignatureProviders                    map[crypto.PublicKey]SignatureProviderType
	Producers                             []chain.AccountName
	Timer                                 time.Time
	ProducerWatermarks                    map[chain.AccountName]producerWatermark
	PendingBlockMode                      uint32
	UnappliedTransactions                 chain.UnappliedTransactionQueue
	ThreadPool                            pool.WorkPool
	MaxTransactionTimeMs                  uint32
	MaxIrreversibleBlockAgeUs             time.Duration
	ProduceTimeOffsetUs                   int32
	LastBlockTimeOffsetUs                 int32
	MaxBlockCPUUsageThresholdUs           uint32
	MaxBlockNetUsageThresholdBytes        uint32
	MaxScheduledTransactionTimePerBlockMs int32
	IrreversibleBlockTime                 time.Duration
	KeosdProviderTimeoutUs                time.Duration
	ProtocolFeaturesToActivate            []chain.DigestType
	PotocolFeaturesSignaled               bool
	ChainPlugin                           *chainplugin.ChainPlugin
	IncomingBlockSubscription             func(chan chain.SignedBlock)
	IncomingTransactionSubscription       chan chain.SignedTransaction
	TransactionAckChannel                 chan chain.TransactionMetadata
	IncomingBlockSyncProvider             chan chain.SignedBlock
	IncomingTransactionAsyncProvider      chan chain.PackedTransaction
	BlacklistedTransactions               TransactionIDWithExpiryIndex
	PendingSnapshotIndex                  PendingSnapshotIndexIndex
	//signals2基于Boost的另一个库signals，实现了线程安全的观察者模式   fc::optional<scoped_connection> _accepted_block_connection;
	AcceptedBlockConnection       func()
	AcceptedBlockHeaderConnection func()
	IrreversibleBlockConnection   func()
	TimerCorelationID             uint32
	IncomingDeferRatio            float64
	SnapshotsDir                  string
	PendingIncomingTransactions   IncomingTransactionQueue
}

//ProducerPlugin ..
type ProducerPlugin struct {
	Impl
}

func init() {
	producerPlugin := &ProducerPlugin{}
	app.App().RegisterPlugin("ProducerPlugin", producerPlugin)
}

//Key ..
type Key struct {
	PublicKey  crypto.PublicKey  `json:"publicKey"`
	PrivateKey crypto.PrivateKey `json:"privateKey"`
}

//Initialize ..
func (a *ProducerPlugin) Initialize() {
	println("ProducerPlugin Initialize")
	a.ChainPlugin = (app.App().FindPlugin("ChainPlugin")).(*chainplugin.ChainPlugin)
	if a.ChainPlugin == nil {
		log.Assert("chain_plugin not found")
	}

	var myChain chain.Controller = a.ChainPlugin.GetChain()
	var unappliedMode uint32
	if myChain.GetReadMode() != controllerenum.Speculative {
		unappliedMode = chain.NonSpeculative
	} else {
		if len(a.Producers) <= 0 {
			unappliedMode = chain.SpeculativeNonProducer
		} else {
			unappliedMode = chain.SpeculativeProducer
		}
	}

	a.UnappliedTransactions.SetMode(unappliedMode)
	if len(config.NodeConf.PrivateKey) > 0 {
		keyIDToWifPairStrings := config.NodeConf.PrivateKey
		for _, keyIDToWifPairString := range keyIDToWifPairStrings {
			//KEY  是一个json 格式的
			key := Key{}
			json.Unmarshal([]byte(keyIDToWifPairString), &key)
			a.SignatureProviders[key.PublicKey] = a.MakeKeySignatureProvider(key.PrivateKey)
			var blankedPrivkey string
			for _, i := range key.PrivateKey.Tostring() {
				blankedPrivkey += string("*")
				i++
			}

			log.AppLog().Infof("\"private-key\" is DEPRECATED, use \"signature-provider=%s=KEY:%s\"", key.PublicKey.ToString(), blankedPrivkey)
		}
	}

	//EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV=KEY:5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3"
	if len(config.NodeConf.SignatureProvider) > 0 {
		keySpecPairs := config.NodeConf.SignatureProvider
		for _, keySpecPair := range keySpecPairs {
			delim := strings.Index(keySpecPair, "=")
			if delim <= 0 {
				log.Assert("Missing \"=\" in the key spec pair")
			}

			pubKeyStr := keySpecPair[0:delim]
			specStr := keySpecPair[delim+1:]
			specDelim := strings.Index(specStr, "=")
			if specDelim <= 0 {
				log.Assert("Missing \":\" in the key spec pair")
			}

			specTypeStr := specStr[0:specDelim]
			specData := specStr[specDelim+1:]
			pubkey := crypto.NewPublicKey(pubKeyStr)

			specDataKey := crypto.NewPrivateKey(specData)
			if specTypeStr == "KEY" {
				a.SignatureProviders[pubkey] = a.MakeKeySignatureProvider(specDataKey)
			} else if specTypeStr == "KEOSD" {
				a.SignatureProviders[pubkey] = a.MakeKeySignatureProvider(specDataKey)
				//TODO跳过
				//my->_signature_providers[pubkey] = make_keosd_signature_provider(my, spec_data, pubkey);
			}
		}
	}

	a.KeosdProviderTimeoutUs = time.Millisecond * time.Duration(config.NodeConf.KeosdProviderTimeout)
	a.ProduceTimeOffsetUs = config.NodeConf.ProduceTimeOffsetUs
	if a.ProduceTimeOffsetUs > 0 && int(a.ProduceTimeOffsetUs) < -chain.BlockIntervalUs {
		log.Assert("produce-time-offset-us %d must be 0 .. -%d", chain.BlockIntervalUs, a.ProduceTimeOffsetUs)
	}

	a.LastBlockTimeOffsetUs = config.NodeConf.ProduceTimeOffsetUs
	if a.LastBlockTimeOffsetUs > 0 && int(a.LastBlockTimeOffsetUs) < -chain.BlockIntervalUs {
		log.Assert("last-block-time-offset-us %d must be 0 .. -%d", chain.BlockIntervalUs, a.ProduceTimeOffsetUs)
	}

	cpuEffortPct := config.NodeConf.CPUEffortPercent
	if cpuEffortPct == 0 {
		cpuEffortPct = chain.DefaultBlockCPUEffortPct / uint32(chain.Percent1)
	}

	if cpuEffortPct < 0 && cpuEffortPct > 100 {
		log.Assert("cpu-effort-percent %d must be 0 - 100", cpuEffortPct)
	}

	cpuEffortPct *= uint32(chain.Percent1)

	CPUEffortOffsetUs := -chain.FengPrercent(uint64(chain.BlockIntervalUs), uint32(chain.Percent100)-uint32(cpuEffortPct))

	lastBlockCPUEffortPct := config.NodeConf.LastBlockCPUEffortPct
	if lastBlockCPUEffortPct == 0 {
		lastBlockCPUEffortPct = chain.DefaultBlockCPUEffortPct / uint32(chain.Percent1)
	}

	if lastBlockCPUEffortPct < 0 && lastBlockCPUEffortPct > 100 {
		log.Assert("last-block-cpu-effort-percent %d must be 0 - 100", lastBlockCPUEffortPct)
	}

	lastBlockCPUEffortPct *= uint32(chain.Percent1)
	lastBlockCPUEffortOffsetUs := -chain.FengPrercent(uint64(chain.BlockIntervalUs), uint32(chain.Percent100)-uint32(lastBlockCPUEffortPct))
	a.ProduceTimeOffsetUs = int32(common.Min(int(a.ProduceTimeOffsetUs), int(CPUEffortOffsetUs)))
	a.LastBlockTimeOffsetUs = int32(common.Min(int(a.LastBlockTimeOffsetUs), int(lastBlockCPUEffortOffsetUs)))
	a.MaxBlockCPUUsageThresholdUs = config.NodeConf.MaxBlockCPUUsageThresholdUs
	if a.MaxBlockCPUUsageThresholdUs == 0 {
		a.MaxBlockCPUUsageThresholdUs = 500
	}

	if int(a.MaxBlockCPUUsageThresholdUs) >= chain.BlockIntervalUs {
		log.Assert("max-block-cpu-usage-threshold-us %d  must be 0  %d", chain.BlockIntervalUs, a.MaxBlockCPUUsageThresholdUs)
	}

	a.MaxBlockNetUsageThresholdBytes = config.NodeConf.MaxBlockNetUsageThresholdBytes
	if a.MaxBlockNetUsageThresholdBytes == 0 {
		a.MaxBlockNetUsageThresholdBytes = 1024
	}

	a.MaxScheduledTransactionTimePerBlockMs = int32(config.NodeConf.MaxScheduledTransactionTimePerBlockMs)

	if a.MaxScheduledTransactionTimePerBlockMs == 0 {
		a.MaxScheduledTransactionTimePerBlockMs = 100
	}

	subjectiveCPULeewayUs := config.NodeConf.SubjectiveCPULeewayUs
	if subjectiveCPULeewayUs != 0 && subjectiveCPULeewayUs != chain.DefaultSubjectiveCPULeewayUs {
		myChain.SetSubjectiveCPULeeway(time.Duration(subjectiveCPULeewayUs) * time.Microsecond)
	}

	a.MaxTransactionTimeMs = config.NodeConf.MaxTransactionTime
	a.MaxIrreversibleBlockAgeUs = time.Duration(config.NodeConf.MaxIrreversibleBlockAge) * time.Second
	incomingTransactionQueueSizeMb := uint64(config.NodeConf.IncomingTransactionQueueSizeMb)
	if incomingTransactionQueueSizeMb == 0 {
		incomingTransactionQueueSizeMb = 1024
	}

	var maxIncomingTransactionQueueSize uint64 = 1024 * 1024 * incomingTransactionQueueSizeMb

	if maxIncomingTransactionQueueSize <= 0 {
		log.Assert("incoming-transaction-queue-size-mb %d must be greater than 0", maxIncomingTransactionQueueSize)
	}

	a.PendingIncomingTransactions.SetMaxIncomingTransactionQueueSize(uint64(maxIncomingTransactionQueueSize))

	a.IncomingDeferRatio = config.NodeConf.IncomingDeferRatio
	threadPoolSize := config.NodeConf.ProducerThreads
	if threadPoolSize == 0 {
		threadPoolSize = chain.DefaultControllerThreadPoolSize
	}

	a.ThreadPool.PoolSize = int(threadPoolSize)
	if config.NodeConf.SnapshotsDir == "" {
		config.NodeConf.SnapshotsDir = "snapshots"
	}

	a.SnapshotsDir = app.App().GetDataDir() + "/" + config.NodeConf.SnapshotsDir
	s, err := os.Stat(a.SnapshotsDir)
	if err == nil {
		if s.IsDir() {
			log.Assert("No such directory %s", a.SnapshotsDir)
		}
	} else {
		if !os.IsExist(err) {
			os.Mkdir(a.SnapshotsDir, 666)
		}
	}

	a.IncomingBlockSubscription = func(block chan chain.SignedBlock) {
		t := <-block
		a.OnIncomingBlock(t, nil)
	}
}

const (
	//Succeeded ..
	Succeeded = iota
	//Failed ..
	Failed = iota
	//WaitingForBlock ..
	WaitingForBlock = iota
	//WaitingForProduction ..
	WaitingForProduction = iota
	//Exhausted ..
	Exhausted = iota
)

func (a *ProducerPlugin) calculatePendingBlockTime() time.Duration {
	myChain := a.ChainPlugin.GetChain()
	now := time.Now().Unix()
	base := common.MaxUint64(uint64(now), myChain.HeadBlockTime())
	minTimeToNextBlock := uint64(chain.BlockIntervalUs) - base%uint64(chain.BlockIntervalUs)
	blockTime := time.Duration(base) + time.Microsecond*time.Duration(minTimeToNextBlock)
	return blockTime
}

const (
	//Producing ..
	Producing = iota
	//Speculating ..
	Speculating = iota
)

//StartBlock ..
func (a *ProducerPlugin) StartBlock() uint32 {
	myChain := a.ChainPlugin.GetChain()
	if !a.ChainPlugin.AcceptTransactions {
		return WaitingForBlock
	}

	hbs := myChain.HeadBlockState()
	now := time.Now()
	blockTime := a.calculatePendingBlockTime()
	previousPendingMode := a.PendingBlockMode
	println(previousPendingMode)
	a.PendingBlockMode = Producing
	b := chain.BlockTimestamp{Slot: uint32(blockTime)}
	scheduledProducer := hbs.GetScheduledProducer(b)
	currentWatermark := a.getWatermark(scheduledProducer.ProducerName)
	numRelevantSignatures := 0
	// scheduled_producer.for_each_key([&](const public_key_type& key){
	// 	const auto& iter = _signature_providers.find(key);
	// 	if(iter != _signature_providers.end()) {
	// 	   num_relevant_signatures++;
	// 	}
	//  });

	var findFlag bool = false
	for _, i := range a.Producers {
		if i == scheduledProducer.ProducerName {
			findFlag = true
			break
		}
	}

	irreversibleBlockAge := a.getIrreversibleBlockAge()
	if !a.ProductionEnabled {
		a.PendingBlockMode = Speculating
	} else if !findFlag {
		a.PendingBlockMode = Speculating
	} else if numRelevantSignatures == 0 {
		log.AppLog().Errorf("Not producing block because I don't have any private keys relevant to authority: %s", scheduledProducer.Authority)
		a.PendingBlockMode = Speculating
	} else if a.PauseProduction {
		log.AppLog().Errorf("Not producing block because production is explicitly paused")
		a.PendingBlockMode = Speculating
	} else if a.MaxIrreversibleBlockAgeUs >= 0 && irreversibleBlockAge >= a.MaxIrreversibleBlockAgeUs {
		log.AppLog().Errorf("Not producing block because the irreversible block is too old [age:%ds, max:%ds]", irreversibleBlockAge/10000000, a.MaxIrreversibleBlockAgeUs/10000000)
		a.PendingBlockMode = Speculating
	}

	if a.PendingBlockMode == Producing {
		blockTimestamp := chain.BlockTimestamp{Slot: uint32(blockTime)}
		if currentWatermark.First > hbs.BlockNum {
			log.AppLog().Errorf("Not producing block because %s signed a block at a higher block number %d than the current fork's head %d", scheduledProducer.ProducerName, currentWatermark.First, hbs.BlockNum)
			a.PendingBlockMode = Speculating
		} else if currentWatermark.Second.IsGreater(blockTimestamp) {
			log.AppLog().Errorf("Not producing block because %s signed a block at the next block time or later %d than the pending block time %d", scheduledProducer.ProducerName, currentWatermark.First, hbs.BlockNum)
		}
	}

	if a.PendingBlockMode == Speculating {
		headBlockAge := now.Unix() - int64(myChain.HeadBlockTime())
		if headBlockAge > int64(time.Second*5) {
			return WaitingForBlock
		}
	}

	if a.PendingBlockMode == Producing {
		startBlockTime := blockTime - time.Microsecond*time.Duration(chain.BlockIntervalUs)
		if now.Unix() < int64(startBlockTime) {
			log.AppLog().Debugf("Not producing block waiting for production window %d %d", hbs.BlockNum+1, blockTime)
			//a.scheduleDelayedProductionLoop(a, a.ca)
		}
	}

	return 0
}

func (a ProducerPlugin) calculateProducerWakeUpTime(refBlockTime chain.BlockTimestamp) {
	// var wakeUpTime time.Time
	// for p := range a.Producers {
	// 	//nextProducerBlockTime : = a.ca
	// }
}

func (a ProducerPlugin) calculateNextBlockTime(produceName chain.AccountName, currentBlockTime chain.BlockTimestamp) time.Time {
	myChain := a.ChainPlugin.GetChain()
	hbs := myChain.HeadBlockState()
	activeSchedule := hbs.ActiveSchedule.Producers
	var findFlag bool = false
	var producerIndex uint32
	for k, b := range activeSchedule {
		if b.ProducerName == produceName {
			findFlag = true
			producerIndex = uint32(k)
			break
		}
	}

	if !findFlag {
		return time.Time{}
	}

	println(producerIndex)
	return time.Time{}
}

func (a ProducerPlugin) scheduleDelayedProductionLoop(weakThis *ProducerPlugin, wakeUpTime time.Duration) {
	if wakeUpTime > 0 {
		log.AppLog().Debugf("Scheduling Speculative/Production Change at %d", wakeUpTime)
		timer := time.NewTimer(wakeUpTime)
		<-timer.C
		go func() {
			a.TimerCorelationID++
			if a.TimerCorelationID == weakThis.TimerCorelationID {
				weakThis.ScheduleProductionLoop()
			}
		}()
	}

}

func (a ProducerPlugin) getIrreversibleBlockAge() time.Duration {
	now := time.Now()
	if now.Unix() < int64(a.IrreversibleBlockTime) {
		return time.Microsecond * 0
	}

	return time.Duration(now.Unix()-int64(a.IrreversibleBlockTime)) * time.Microsecond
}

type producerWatermark struct {
	First  uint32
	Second chain.BlockTimestamp
}

func (a ProducerPlugin) getWatermark(producer chain.AccountName) producerWatermark {
	itr, ok := a.ProducerWatermarks[producer]
	if !ok {
		return producerWatermark{}
	}

	return itr
}

//ScheduleProductionLoop ..
func (a *ProducerPlugin) ScheduleProductionLoop() {
	a.Timer.Clock()
	//result :=
}

//OnIncomingBlock ..
func (a *ProducerPlugin) OnIncomingBlock(block chain.SignedBlock, blockID *chain.BlockIDType) bool {
	myChain := a.ChainPlugin.GetChain()
	if a.PendingBlockMode == nodeenum.Porducing {
		var idString string
		idString = "UNKNOWN"
		if blockID != nil {
			idString = blockID.String()
		}

		log.AppLog().Infof("dropped incoming block %d id: %s", block.BlockNum(), idString)
		return false
	}

	var id *chain.BlockIDType
	if blockID != nil {
		id = blockID
	} else {
		idtmp := block.ID()
		id = &idtmp
	}

	blkNum := block.BlockNum()
	log.AppLog().Debugf("received incoming block blkNum %d  id:%d", blkNum, id)
	if block.Timestamp.Time() >= uint64(time.Now().Unix()+int64(time.Second*7)) {
		log.Assert("received a block from the future, ignoring it: id :%d", id)
	}

	existing := myChain.FetchBlockByID(*id)
	if existing != nil {
		return false
	}

	//bsf := myChain.CreateBlockStateFuture(&block)
	a.UnappliedTransactions.AddAborted(myChain.AbortBlock())
	ensure := func() {
		a.ScheduleProductionLoop()
	}

	ensure()
	//myChain.

	return true
}

//MakeKeySignatureProvider ..
func (a *ProducerPlugin) MakeKeySignatureProvider(key chain.PrivateKeType) SignatureProviderType {
	return func(digist chain.DigestType) chain.SignatureType {
		return key.Sign(digist)
	}
}

//MakeKeosdSignatureProvider TODO 先跳过
func (a *ProducerPlugin) MakeKeosdSignatureProvider(impl *Impl, urlStr string, pubkey chain.PublicKeyType) {
	//var keosdUrl url.URL
	// if(boost::algorithm::starts_with(url_str, "unix://"))
	//    //send the entire string after unix:// to http_plugin. It'll auto-detect which part
	//    // is the unix socket path, and which part is the url to hit on the server
	//    keosd_url = fc::url("unix", url_str.substr(7), ostring(), ostring(), ostring(), ostring(), ovariant_object(), fc::optional<uint16_t>());
	// else
	//    keosd_url = fc::url(url_str);
	// std::weak_ptr<producer_plugin_impl> weak_impl = impl;

	// return [weak_impl, keosd_url, pubkey]( const chain::digest_type& digest ) {
	//    auto impl = weak_impl.lock();
	//    if (impl) {
	// 	  fc::variant params;
	// 	  fc::to_variant(std::make_pair(digest, pubkey), params);
	// 	  auto deadline = impl->_keosd_provider_timeout_us.count() >= 0 ? fc::time_point::now() + impl->_keosd_provider_timeout_us : fc::time_point::maximum();
	// 	  return app().get_plugin<http_client_plugin>().get_client().post_sync(keosd_url, params, deadline).as<chain::signature_type>();
	//    } else {
	// 	  return signature_type();
	//    }
	// };
}

//HandleSighup ..
func (a *ProducerPlugin) HandleSighup() {
	println("ProducerPlugin HandleSighup")
}

//StartUp ..
func (a *ProducerPlugin) StartUp() {
	println("ProducerPlugin StartUp")
}

//PluginStartUp ..
func (a *ProducerPlugin) PluginStartUp() {
	println("ProducerPlugin PluginStartUp")
}
