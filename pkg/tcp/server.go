package tcp

import (
	"fmt"
	"log"
	"net"
	"wheels/unit"
)

type IServer interface {
	start()
	Stop()
	Serve()
}

type Server struct {
	Name    string
	Version string
	IP      string
	Port    int
	msg     chan Message
	conns   []*Connection
}

func NewServer() IServer {
	return &Server{
		IP:      unit.GlabalObject.IP,
		Port:    unit.GlabalObject.Port,
		Version: "tcp4",
		Name:    unit.GlabalObject.Name,
		msg:     make(chan Message, 100),
	}
}

func (s *Server) start() {
	fmt.Println("Starting...")
	addr, err := net.ResolveTCPAddr(s.Version, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		log.Println(err)
		return
	}
	listener, err := net.ListenTCP(s.Version, addr)
	if err != nil {
		log.Println(err)
		return
	}
	printlog()
	go s.Broadcast()
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Printf("Connection at %s \n", conn.RemoteAddr().String())
		c := NewConnection(s, conn, func(m Message) {
			//log.Println(string(m.data))
			s.msg <- m
		})
		s.conns = append(s.conns, c)
		go c.Start()
	}

}
func (s *Server) Stop() {

}
func (s *Server) Serve() {
	s.start()
	select {}
}

func (s *Server) Broadcast() {
	for {
		select {
		case str := <-s.msg:
			for i := range s.conns {
				if str.author == s.conns[i].conn.RemoteAddr().String() {
					continue
				}
				s.conns[i].SendChan <- str
			}
		}

	}

}

func printlog() {
	fmt.Println("Server Name:", unit.GlabalObject.Name)
	fmt.Printf("Listening at %s:%d \n", unit.GlabalObject.IP, unit.GlabalObject.Port)
}
