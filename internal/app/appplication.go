package app

import (
	"errors"
	"feng/internal/log"
	"os"
	"sync"
)

//Appplication ..
type Appplication struct {
	version        uint64
	versionStr     string
	fullVersionStr string
	configDir      string
	dataDir        string
	plugins        map[string]AbstractPlugin
}

var app *Appplication
var onceApp sync.Once
var oncePlugin sync.Once

//App ..
func App() *Appplication {
	onceApp.Do(func() {
		app = &Appplication{}
	})

	return app
}

//GetVersionStr 获取名称
func (a *Appplication) GetVersionStr() string {
	return a.versionStr
}

//SetVersion 设置版本号
func (a *Appplication) SetVersion(version uint64) {
	a.version = version
}

//SetVersionStr 设置
func (a *Appplication) SetVersionStr(versionStr string) {
	a.versionStr = versionStr
}

//SetFullVersionStr 设置
func (a *Appplication) SetFullVersionStr(fullVersionStr string) {
	a.fullVersionStr = fullVersionStr
}

//SetConfigDir 设置
func (a *Appplication) SetConfigDir(configDir string) {
	a.configDir = configDir
}

//SetDataDir 设置
func (a *Appplication) SetDataDir(dataDir string) {
	a.dataDir = dataDir
}

//RegisterPlugin ..
func (a *Appplication) RegisterPlugin(pluginName string, pluginStruct AbstractPlugin) {
	if pluginStruct == nil {
		log.AppLog().Errorf("pluginName is error")
		os.Exit(-1)
	}

	oncePlugin.Do(func() {
		App().plugins = make(map[string]AbstractPlugin)
	})

	App().plugins[pluginName] = pluginStruct
}

//Initialize 初始化插件
func (a *Appplication) Initialize(plugins []string) error {
	println("Initialize start size:", len(App().plugins))
	for _, plugin := range plugins {
		startPlugin, ok := App().plugins[plugin]
		if !ok {
			return errors.New("-1")
		}

		startPlugin.Initialize()
	}

	return nil
}
