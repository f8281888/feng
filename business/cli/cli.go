package cli

import (
	"feng/internal/log"
	"net"
)

//Start ..
func Start() {
	println("cli start")
	addr := "0.0.0.0:9876"
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Assert("ResolveTCPAddr is error")
	}

	var laddr net.TCPAddr
	coon, err := net.DialTCP("tcp", &laddr, tcpAddr)
	defer coon.Close()

	coon.Write([]byte("i am a client"))
	log.AppLog().Infof(laddr.IP.String())
}
