package cli

import (
	"feng/internal/log"
	"net"
	"time"
	"unsafe"
)

//Message ..
type Message struct {
	len  int
	body []byte
}

//SliceMock .. 跟slice 底层一样，转换去互转
type SliceMock struct {
	addr uintptr
	len  int
	cap  int
}

func structTobyte(m Message) []byte {
	lenM := unsafe.Sizeof(m)
	tmp := SliceMock{
		addr: uintptr(unsafe.Pointer(&m)),
		len:  int(lenM),
		cap:  int(lenM),
	}

	return *(*[]byte)(unsafe.Pointer(&tmp))
}

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

	go func() {
		for {
			body := []byte("i am client")
			a := Message{len: len(body), body: body}
			len, err := coon.Write(structTobyte(a))
			if err != nil {
				log.Assert("can't write")
			}

			log.AppLog().Infof("ip:%s,len:%d", laddr.IP.String(), len)
			time.Sleep(3600 * time.Second)
		}
	}()

	// go func() {
	// 	for {
	// 		println("nihao")
	// 		// b := make([]byte, 1024)
	// 		// len, err := coon.Read(b)
	// 		// if err != nil {
	// 		// 	log.AppLog().Errorf("can't read")
	// 		// }

	// 		// log.AppLog().Infof("b:%s,len:%d", string(b), len)
	// 		time.Sleep(1000)
	// 	}
	// }()

	select {}
}
