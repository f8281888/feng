package socket

import (
	"feng/internal/log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/graarh/golang-socketio/transport"

	socketio "github.com/graarh/golang-socketio"
)

type subscribeFunc func(key, channel, pushType, payload string, msgTime int64)

//SendMsg ..
type SendMsg struct {
	Channel  string `json:"channel"`   // 频道
	PushType string `json:"push_type"` // 推送类型
	Payload  string `json:"payload"`   // 消息内容
	Time     int64  `json:"time"`      // 消息发送的时间
}

//MySocketManger ..
type MySocketManger struct {
	server        *socketio.Server
	HTTPListener  *net.TCPListener
	HTTPSListener *net.TCPListener
	mux           sync.Mutex
	roomMap       map[string]string
	subscribeMap  map[string]map[string]subscribeFunc
	sendChan      chan SendMsg
	syncChan      chan struct{}
	MyConn        *Connection
}

var defaultOption = &transport.WebsocketTransport{
	PingInterval:   time.Second * 10,
	PingTimeout:    time.Minute,
	ReceiveTimeout: time.Minute,
	SendTimeout:    time.Minute,
	BufferSize:     1024 * 1024,
}

//NewSocketManager ..
func NewSocketManager(opts *transport.WebsocketTransport) *MySocketManger {
	if opts == nil {
		opts = defaultOption
	}

	sm := MySocketManger{
		server:       socketio.NewServer(opts),
		roomMap:      make(map[string]string, 0),
		subscribeMap: make(map[string]map[string]subscribeFunc, 0),
		sendChan:     make(chan SendMsg, 0),
		syncChan:     make(chan struct{}, 0),
	}

	sm.addConsume()
	// go sm.listeningChannel()
	// go sm.publish()
	return &sm
}

func (s *MySocketManger) addConsume() {
}

//ListeningChannel ..
func (s *MySocketManger) ListeningChannel() {

}

func (s *MySocketManger) publish() {

}

//StartAccpet ..
func (s *MySocketManger) StartAccpet() {
	if s.HTTPListener != nil {
		s.HandletHTTP()
	}

	if s.HTTPSListener != nil {
		s.HandleHTTPS()
	}
}

//SetTLSInitHandler ..
func (s *MySocketManger) SetTLSInitHandler() {

}

//ClearAccessChannels ..
func (s *MySocketManger) ClearAccessChannels() {

}

//Init ..
func (s *MySocketManger) Init() {

}

//SetMaxHTTPBodySize ..
func (s *MySocketManger) SetMaxHTTPBodySize(newValue uint32) {

}

//HandletHTTP ..
func (s *MySocketManger) HandletHTTP() {
	newConnection := new(Connection)
	if s.HTTPListener == nil {
		log.Assert("Listener is nil %s", "nihao")
	}

	newConnection.Conn, _ = s.HTTPListener.Accept()
	if newConnection.Conn != nil {
		s.MyConn = newConnection
		s.StartSession()
	} else {
		log.Assert("HandletHTTP listener is error")
	}
}

//StartSession ..
func (s *MySocketManger) StartSession() {
	if s.MyConn.Conn == nil {
		return
	}

	readBytes := make([]byte, 1024)
	// var closeConnection bool
	_, err := s.MyConn.Conn.Read(readBytes)
	if err != nil {
		log.Assert("StartSession read is error")
	}

	result := strings.Replace(string(readBytes), "\n", "", 1)
	log.AppLog().Debugf("http read :%s", result)
}

//HandleHTTPS ..
func (s *MySocketManger) HandleHTTPS() {

}
