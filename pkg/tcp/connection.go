package tcp

import (
	"fmt"
	"io"
	"log"
	"net"
)

type IConnection interface {
	Start()
	Read()
	Handler()
}
type HandlerFunc func(Message)

type Connection struct {
	conn        *net.TCPConn
	Server      IServer
	HandlerFunc HandlerFunc

	//RecvChan chan<- Message
	SendChan chan Message
}

func NewConnection(Server IServer, conn *net.TCPConn, HandlerFunc HandlerFunc) *Connection {
	ch := make(chan Message, 100)
	return &Connection{
		conn:        conn,
		Server:      Server,
		HandlerFunc: HandlerFunc,
		//RecvChan:    ch,
		SendChan: ch,
	}
}

func (c *Connection) Start() {
	go c.Handler()
	for {
		buf := make([]byte, 256)
		_, err := c.conn.Read(buf)
		if err != nil && err == io.EOF {
			log.Println(err)
			return
		}
		r := Message{
			data:   buf,
			author: c.conn.RemoteAddr().String(),
		}
		c.HandlerFunc(r)
	}
}

func (c *Connection) Read() {
	// for {
	// 	buf := make([]byte, 256)
	// 	_, err := c.conn.Read(buf)
	// 	if err != nil {
	// 		log.Println(err)
	// 		return
	// 	}
	// 	r := Message{
	// 		data: buf,
	// 	}
	// 	//c.RecvChan <- r
	// }
}

func (c *Connection) Handler() {
	for str := range c.SendChan {
		c.conn.Write([]byte(fmt.Sprintf("%s:%s \n", c.conn.RemoteAddr().String(), str.data)))
	}
}
