package netplugin

import "feng/internal/app"

//NetPlugin ..
type NetPlugin struct {
}

func init() {
	netPlugin := &NetPlugin{}
	app.App().RegisterPlugin("NetPlugin", netPlugin)
}

//Initialize ..
func (a *NetPlugin) Initialize() {
	println("NetPlugin Initialize")
}
