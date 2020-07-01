package app

//AbstractPlugin ..
type AbstractPlugin interface {
	Initialize()
	HandleSighup()
	StartUp()
	PluginStartUp()
}
