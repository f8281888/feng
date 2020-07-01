package producerplugin

import "feng/internal/app"

//ProducerPlugin ..
type ProducerPlugin struct {
}

func init() {
	producerPlugin := &ProducerPlugin{}
	app.App().RegisterPlugin("ProducerPlugin", producerPlugin)
}

//Initialize ..
func (a *ProducerPlugin) Initialize() {
	println("ProducerPlugin Initialize")
}

//HandleSighup ..
func (a *ProducerPlugin) HandleSighup() {
	println("ProducerPlugin HandleSighup")
}

//StartUp ..
func (a *ProducerPlugin) StartUp() {
	println("ProducerPlugin StartUp")
}

//PluginStartUp ..
func (a *ProducerPlugin) PluginStartUp() {
	println("ProducerPlugin PluginStartUp")
}
