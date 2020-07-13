package producerplugin

import (
	"encoding/json"
	"feng/config"
	controllerenum "feng/enum/controller-enum"
	"feng/internal/app"
	"feng/internal/chain"
	"feng/internal/fc/common"
	"feng/internal/fc/crypto"
	"feng/internal/fc/stl"
	"feng/internal/log"
	chainplugin "feng/internal/plugin/chain_plugin"
	"feng/internal/pool"
	"fmt"
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

//Impl ..
type Impl struct {
	ProductionEnabled                     bool
	PauseProduction                       bool
	SignatureProviders                    map[crypto.PublicKey]crypto.Signature
	Producers                             []chain.AccountName
	Timer                                 time.Time
	ProducerWatermarks                    map[chain.AccountName]stl.Pair
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
	IncomingBlockSubscription             chan chain.SignedBlock
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
			keyIDToWifPair := json.Unmarshal([]byte(keyIDToWifPairString), &key)
			a.SignatureProviders[key.PublicKey] = a.MakeKeySignatureProvider(key.PrivateKey)

		}
	}

}

//SignatureProviderType ..
type SignatureProviderType = func(digist chain.DigestType) chain.SignatureType

//MakeKeySignatureProvider ..
func (a *ProducerPlugin) MakeKeySignatureProvider(key chain.PrivateKeType) SignatureProviderType {
	return func(digist chain.DigestType) chain.SignatureType {
		return key.Sign(digist)
	}
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
