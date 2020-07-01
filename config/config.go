package config

import (
	"fmt"
	"os"

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
	P2pMaxNodesPerHost     uint32            `json:"p2pMaxNodesPerHost"`
	P2pAcceptTransactions  bool              `json:"p2pAcceptTransactions"`
	UseSocketReadWatermark bool              `json:"useSocketReadWatermark"`
	P2pListenEndpoint      string            `json:"p2pListenEndpoint"`
	P2pServerAddress       string            `json:"p2pServerAddress"`
	NetThreads             uint16            `json:"netThreads"`
	P2pPeerAddress         []string          `json:"p2pPeerAddress"`
	AgentName              string            `json:"agentName"`
	AllowedConnection      []string          `json:"allowedConnection"`
	PeerKey                []string          `json:"peerKey"`
	PeerPrivateKey         map[string]string `json:"peerPrivateKey"`
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
		fmt.Println(err)
		os.Exit(1)
	}

	if err := viper.Unmarshal(&value); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
