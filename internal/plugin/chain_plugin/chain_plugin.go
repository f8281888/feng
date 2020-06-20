package chainplugin

import "feng/internal/app"

//ChainPlugin ..
type ChainPlugin struct {
}

func init() {
	chainPlugin := &ChainPlugin{}
	app.App().RegisterPlugin("ChainPlugin", chainPlugin)
}

//Initialize ..
func (a *ChainPlugin) Initialize() {
	println("ChainPlugin Initialize")
}
