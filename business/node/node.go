package node

import (
	"feng/config"
	"feng/internal/app"
	"feng/internal/log"
	_ "feng/internal/plugin/chain_plugin"
	httpplugin "feng/internal/plugin/http_plugin"
	_ "feng/internal/plugin/net_plugin"
	_ "feng/internal/plugin/producer_plugin"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

//Start 启动节点
func Start() {
	log.AppLog().Infof("node business start")
	app.App().SetVersion(config.NodeConf.Version)
	app.App().SetVersionStr(config.NodeConf.VersionStr)
	app.App().SetFullVersionStr(config.NodeConf.FullVersionStr)
	app.App().SetConfigDir(filepath.Join(viper.GetString("config-path")))
	app.App().SetConfigDir(filepath.Join(viper.GetString("data-path")))

	httpplugin.SetDefaults(httpplugin.HTTPluginDefaults{
		DefaultUnixSocketPath: "",
		DefaulthttpPort:       8888,
	})

	println(app.App().GetVersionStr())
	println(httpplugin.HTTPluginDefault.DefaulthttpPort)

	plugins := []string{"ChainPlugin", "NetPlugin", "ProducerPlugin"}
	if err := app.App().Initialize(plugins); err != nil {
		println(err.Error())
		os.Exit(-1)
	}

	log.AppLog().Infof("node %s", app.App().GetVersionStr())
}
