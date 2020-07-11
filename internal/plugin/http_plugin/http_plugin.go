package httpplugin

import (
	"feng/config"
	"feng/internal/app"
	"feng/internal/fc/common"
	"feng/internal/log"
	"feng/internal/pool"
	"feng/internal/socket"
	"net"
	"strings"
	"time"
)

//HTTPluginDefaults ..
type HTTPluginDefaults struct {
	DefaultUnixSocketPath string
	DefaulthttpPort       uint16
}

//HTTPluginDefault ..
var HTTPluginDefault HTTPluginDefaults

type abstractConn interface {
	verifyMaxBytesInFlight() bool
	handleException()
}

type urlResponseCallback func(int, common.Variant)

const (
	//SECP384R1 ..
	SECP384R1 = 0
	//PRIME256V1 ..
	PRIME256V1 = 1
)

//Impl ..
type Impl struct {
	URLHandlers                   map[string]func(abstractConn, string, string, urlResponseCallback)
	HTTPAddress                   *net.TCPAddr
	HTTPSAddress                  *net.TCPAddr
	AccessControlAllowOrigin      string
	AccessControlAllowHeaders     string
	AccessControlMaxAge           string
	AccessControlAllowCredentials bool
	MaxBodySize                   uint32
	Server                        *socket.MySocketManger
	ThreadPoolSize                uint16
	ThreadPool                    *pool.WorkPool
	BytesInFlight                 uint32
	MaxBytesInFlight              uint32
	MaxResponseTime               time.Duration
	HTTPSCertChain                string
	HTTPSKey                      string
	HTTPSEcdhCurve                uint32
	HTTPSServer                   *socket.MySocketManger
	ValidateHost                  bool
	ValidHosts                    []string
}

//HTTPPlugin ..
type HTTPPlugin struct {
	Impl
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

var verboseHTTPErrors bool = false

//Initialize ..
func (a *HTTPPlugin) Initialize() {
	println("HTTPPlugin Initialize")
	//启用一个普通的连接和一个https 连接
	a.Server = socket.NewSocketManager(nil)
	a.HTTPSServer = socket.NewSocketManager(nil)
	a.ValidateHost = config.NodeConf.HTTPValidateHost
	a.ValidHosts = config.NodeConf.HTTPAlias
	if len(config.NodeConf.HTTPServerAddress) > 0 {
		lipStr := config.NodeConf.HTTPServerAddress
		laddr, err := net.ResolveTCPAddr("tcp", lipStr)
		a.HTTPAddress = laddr
		if err != nil {
			log.Assert("failed to configure http to listen on %s", lipStr)
		}

		log.AppLog().Infof("configured http to listen on %s", lipStr)
		arr := strings.Split(lipStr, ":")
		a.AddAliasesForEndpoint(laddr, arr[0], arr[1])
	}

	//unix-socket-path
	if len(config.NodeConf.HTTPSServerAddress) > 0 {
		if len(config.NodeConf.HTTPSCertificateChainFile) <= 0 {
			log.Assert("https-certificate-chain-file is required for HTTPS")
		}

		if len(config.NodeConf.HTTPSPrivateKeyFile) < 0 {
			log.Assert("https-private-key-file is required for HTTPS")

		}

		lipStr := config.NodeConf.HTTPSServerAddress
		laddr, err := net.ResolveTCPAddr("tcp", lipStr)
		a.HTTPSAddress = laddr
		if err != nil {
			log.Assert("failed to configure https to listen on %s", lipStr)
		}

		log.AppLog().Infof("configured https to listen on %s", lipStr)
		a.HTTPSCertChain = config.NodeConf.HTTPSCertificateChainFile
		a.HTTPSKey = config.NodeConf.HTTPSPrivateKeyFile
		arr := strings.Split(lipStr, ":")
		a.AddAliasesForEndpoint(laddr, arr[0], arr[1])
	}

	a.MaxBodySize = config.NodeConf.MaxBodySize
	verboseHTTPErrors = config.NodeConf.VerboseHTTPErrors
	a.ThreadPoolSize = config.NodeConf.HTTPThreads
	if a.ThreadPoolSize <= 0 {
		log.Assert("http-threads %d must be greater than 0", a.ThreadPoolSize)
	}

	a.MaxBytesInFlight = config.NodeConf.HTTPMaxBytesInFlightMb * 1024 * 1024
	a.MaxResponseTime = time.Microsecond * time.Duration(config.NodeConf.HTTPMaxResponseTimeMs)
}

//AddAliasesForEndpoint ..
func (a *HTTPPlugin) AddAliasesForEndpoint(addr *net.TCPAddr, host, port string) {
	a.ValidHosts = append(a.ValidHosts, host+":"+port)
	a.ValidHosts = append(a.ValidHosts, addr.IP.String()+":"+port)
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
	a.Initialize()
	println("HTTPPlugin PluginStartUp")
	a.ThreadPool = pool.NewPool(int(a.ThreadPoolSize))
	var err error
	if a.HTTPAddress != nil {
		go func() {
			log.AppLog().Debugf("go HTTP handle")
			a.Server.HTTPListener, err = net.ListenTCP("tcp", a.HTTPAddress)
			if err != nil {
				log.Assert("http_plugin::plugin_startup failed to http bind to port")
			}

			a.Server.StartAccpet()
		}()
	}

	if a.HTTPSAddress != nil {
		go func() {
			a.HTTPSServer.HTTPSListener, err = net.ListenTCP("tcp", a.HTTPSAddress)
			if err != nil {
				log.Assert("http_plugin::plugin_startup failed to https bind to port")
			}

			a.HTTPSServer.StartAccpet()
		}()
	}
}
