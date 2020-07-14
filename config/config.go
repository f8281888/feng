package config

import (
	"feng/internal/log"

	"github.com/spf13/viper"
)

//NodeConf ..
var NodeConf NodeConfig

//NodeConfig 配置表信息
type NodeConfig struct {
	Version        uint64 `json:"version"`
	VersionStr     string `json:"versionStr"`
	FullVersionStr string `json:"fullversionStr"`
	//peer的日志格式
	PeerLogFormat string `json:"peerLogFormat"`
	//在同步期间从任何单个对等方在块中检索的块数
	SyncFetchSpan uint64 `json:"syncFetchSpan"`
	//清除不可用连接的时间
	ConnectionCleanupPeriod int64 `json:"connectionCleanupPeriod"`
	//最大清除时间
	MaxCleanupTimeMsec int32 `json:"maxCleanupTimeMsec"`
	//最大连接数量
	MaxClients uint32 `json:"maxClients"`
	//一个ip地址最多可以连接的节点数量
	P2pMaxNodesPerHost                    uint32            `json:"p2pMaxNodesPerHost"`
	P2pAcceptTransactions                 bool              `json:"p2pAcceptTransactions"`
	UseSocketReadWatermark                bool              `json:"useSocketReadWatermark"`
	P2pListenEndpoint                     string            `json:"p2pListenEndpoint"`
	P2pServerAddress                      string            `json:"p2pServerAddress"`
	NetThreads                            uint16            `json:"netThreads"`
	P2pPeerAddress                        []string          `json:"p2pPeerAddress"`
	AgentName                             string            `json:"agentName"`
	AllowedConnection                     []string          `json:"allowedConnection"`
	PeerKey                               []string          `json:"peerKey"`
	PeerPrivateKey                        map[string]string `json:"peerPrivateKey"`
	HTTPValidateHost                      bool              `json:"httpValidateHost"`
	HTTPAlias                             []string          `json:"httpAlias"`
	HTTPServerAddress                     string            `json:"httpServerAddress"`
	HTTPSCertificateChainFile             string            `json:"httpsCertificateChainFile"`
	HTTPSServerAddress                    string            `json:"httpsServerAddress"`
	HTTPSPrivateKeyFile                   string            `json:"httpsPrivateKeyFile"`
	MaxBodySize                           uint32            `json:"maxBodySize"`
	VerboseHTTPErrors                     bool              `json:"verboseHttpErrors"`
	HTTPThreads                           uint16            `json:"httpThreads"`
	HTTPMaxBytesInFlightMb                uint32            `json:"httpMaxBytesInFlightMb"`
	HTTPMaxResponseTimeMs                 uint32            `json:"httpMaxResponseTimeMs"`
	PrivateKey                            []string          `json:"privateKey"`
	SignatureProvider                     []string          `json:"signatureProvider"`
	KeosdProviderTimeout                  uint32            `json:"keosdProviderTimeout"`
	ProduceTimeOffsetUs                   int32             `json:"produceTimeOffsetUs"`
	LastBlockTimeOffsetUs                 int32             `json:"lastBlockTimeOffsetUs"`
	CPUEffortPercent                      uint32            `json:"cpuEffortPercent"`
	LastBlockCPUEffortPct                 uint32            `json:"lastBlockCPUEffortPct"`
	MaxBlockCPUUsageThresholdUs           uint32            `json:"maxBlockCPUUsageThresholdUs"`
	MaxBlockNetUsageThresholdBytes        uint32            `json:"maxBlockNetUsageThresholdBytes"`
	MaxScheduledTransactionTimePerBlockMs uint32            `json:"maxScheduledTransactionTimePerBlockMs"`
	SubjectiveCPULeewayUs                 uint32            `json:"subjectiveCPULeewayUs"`
	MaxTransactionTime                    uint32            `json:"maxTransactionTime"`
	MaxIrreversibleBlockAge               int32             `json:"maxIrreversibleBlockAge"`
	IncomingTransactionQueueSizeMb        int16             `json:"incomingTransactionQueueSizeMb"`
	IncomingDeferRatio                    float64           `json:"incomingDeferRatio"`
	ProducerThreads                       int16             `json:"producerThreads"`
	SnapshotsDir                          string            `json:"snapshotsDir"`
}

//InitConfig 初始化配置
func InitConfig(configPath, pre string, value interface{}) {
	viper.AddConfigPath(configPath)
	viper.SetConfigType("json")
	configName := "config"
	if pre != "" {
		configName = pre + "-" + configName
	}

	viper.SetConfigName(configName)

	if err := viper.ReadInConfig(); err != nil {
		log.Assert(err.Error())
	}

	if err := viper.Unmarshal(&value); err != nil {
		log.Assert(err.Error())
	}
}
