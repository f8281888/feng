package node

import (
	"feng/config"
	"feng/internal/app"
	"feng/internal/log"
	_ "feng/internal/plugin/chain_plugin"
	httpplugin "feng/internal/plugin/http_plugin"
	_ "feng/internal/plugin/net_plugin"
	_ "feng/internal/plugin/producer_plugin"
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
	app.App().SetDataDir(filepath.Join(viper.GetString("data-path")))

	httpplugin.SetDefaults(httpplugin.HTTPluginDefaults{
		DefaultUnixSocketPath: "",
		DefaulthttpPort:       8888,
	})

	println(app.App().GetVersionStr())
	println(httpplugin.HTTPluginDefault.DefaulthttpPort)

	plugins := []string{"ChainPlugin", "NetPlugin", "ProducerPlugin"}
	if err := app.App().Initialize(plugins); err != nil {
		log.Assert(err.Error())
	}

	log.AppLog().Infof("node %s %s", app.App().GetVersionStr(), app.App().GetFullVersionStr())
	log.AppLog().Infof("node using configuration file %s", app.App().GetConfiDir())
	log.AppLog().Infof("node data directory is %s", app.App().GetDataDir())
	app.App().StartUp()
}
