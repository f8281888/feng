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
