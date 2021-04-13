package tcp

import (
	"log"
	"net"
	"testing"
	"time"
	"wheels/unit"
)

func ClientTest() {

	conn, err := net.Dial("tcp4", "127.0.0.1:8888")
	if err != nil {
		log.Println(err)
		return
	}
	for {
		_, err = conn.Write([]byte("abcccc"))
		if err != nil {
			log.Println(err)
			return
		}
		buf := make([]byte, 1024)
		_, err = conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(buf))
		time.Sleep(1 * time.Second)
	}

}

func TestServer(t *testing.T) {
	unit.PrintLogo()
	s := NewServer()
	s.Serve()
	//go ClientTest()
}
