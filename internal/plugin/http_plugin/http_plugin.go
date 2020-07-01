package httpplugin

import (
	"feng/internal/app"
)

//HTTPluginDefaults ..
type HTTPluginDefaults struct {
	DefaultUnixSocketPath string
	DefaulthttpPort       uint16
}

//HTTPluginDefault ..
var HTTPluginDefault HTTPluginDefaults

//HTTPPlugin ..
type HTTPPlugin struct {
}

func init() {
	httpPlugin := &HTTPPlugin{}
	app.App().RegisterPlugin("HTTPPlugin", httpPlugin)
}

//SetDefaults ..
func SetDefaults(config HTTPluginDefaults) {
	HTTPluginDefault = config
}

//Startup ..
func (a *HTTPPlugin) Startup() {

}

//Initialize ..
func (a *HTTPPlugin) Initialize() {
	println("HTTPPlugin Initialize")
}

//HandleSighup ..
func (a *HTTPPlugin) HandleSighup() {
	println("HTTPPlugin HandleSighup")
}

//StartUp ..
func (a *HTTPPlugin) StartUp() {
	println("HTTPPlugin StartUp")
}

//PluginStartUp ..
func (a *HTTPPlugin) PluginStartUp() {
	println("HTTPPlugin PluginStartUp")
}
