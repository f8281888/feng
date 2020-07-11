package app

import (
	"errors"
	"feng/internal/log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type ApplicationImpl struct {
	IsQuiting bool
}

//Appplication ..
type Appplication struct {
	version            uint64
	versionStr         string
	fullVersionStr     string
	configDir          string
	dataDir            string
	plugins            map[string]AbstractPlugin
	initializedPlugins map[string]AbstractPlugin
	runningPlugins     map[string]AbstractPlugin
	myApplicationImpl  *ApplicationImpl
}

var app *Appplication
var onceApp sync.Once
var oncePlugin sync.Once
var onceRunningPlugins sync.Once

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

//GetFullVersionStr 获取名称
func (a *Appplication) GetFullVersionStr() string {
	return a.fullVersionStr
}

//GetConfiDir 获取名称
func (a *Appplication) GetConfiDir() string {
	return a.configDir
}

//GetDataDir 获取名称
func (a *Appplication) GetDataDir() string {
	return a.dataDir
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
		log.Assert("pluginName is error")
	}

	oncePlugin.Do(func() {
		App().plugins = make(map[string]AbstractPlugin)
	})

	App().plugins[pluginName] = pluginStruct
}

//FindPlugin ..
func (a *Appplication) FindPlugin(pluginName string) AbstractPlugin {
	if len(a.plugins) == 0 {
		log.AppLog().Errorf("plugins is null")
	}

	plugin, ok := a.plugins[pluginName]
	if !ok {
		log.Assert("can't find plugin %s", pluginName)
	}

	return plugin
}

//Initialize 初始化插件
func (a *Appplication) Initialize(plugins []string) error {
	println("Initialize start size:", len(App().plugins))
	App().initializedPlugins = make(map[string]AbstractPlugin)
	for _, plugin := range plugins {
		startPlugin, ok := App().plugins[plugin]
		if !ok {
			return errors.New("-1")
		}

		startPlugin.Initialize()
		App().initializedPlugins[plugin] = startPlugin
	}

	return nil
}

var exitChan chan os.Signal

//StartUp 启动
func (a *Appplication) StartUp() {
	println("StartUp")
	exitChan = make(chan os.Signal)
	signal.Notify(exitChan, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGTERM)
	go startSighupHandler()
	var wg sync.WaitGroup
	wg.Add(1)
	for pluginName, plugin := range App().plugins {
		defer wg.Done()
		plugin.PluginStartUp()
		plugin.StartUp()
		a.setRuntingPlugin(pluginName, plugin)
	}

	wg.Wait()
}

func startSighupHandler() {
	s := <-exitChan
	log.AppLog().Infof("receive signal :%s", s.String())
	for _, plugin := range App().plugins {
		plugin.HandleSighup()
	}

	log.Assert("startSighupHandler")
}

//setRuntingPlugin 设置
func (a *Appplication) setRuntingPlugin(pluginName string, pluginStruct AbstractPlugin) {
	onceRunningPlugins.Do(func() {
		a.runningPlugins = make(map[string]AbstractPlugin)
	})
	a.runningPlugins[pluginName] = pluginStruct
}

//IsQuiting ..
func (a *Appplication) IsQuiting() bool {
	return a.myApplicationImpl.IsQuiting
}
