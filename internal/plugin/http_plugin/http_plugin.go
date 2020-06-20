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
