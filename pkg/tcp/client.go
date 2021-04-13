package tcp

import (
	"fmt"
	"log"
	"net"
	"wheels/unit"
)

type Client struct {
	Name string
	conn net.Conn
}

func NewClient(Name string) *Client {
	conn, err := net.Dial("tcp4", fmt.Sprintf("%s:%d", unit.GlabalObject.IP, unit.GlabalObject.Port))
	if err != nil {
		log.Println(err)
		return nil
	}
	c := Client{
		Name: Name,
		conn: conn,
	}
	go c.Recv()
	go c.Send()

	return &c
}
func (c *Client) Send() {
	for {
		msg := ""
		_, err := fmt.Scanln(&msg)
		if err != nil {
			log.Println("input err:", err)
		}
		log.Printf("%s:%s \n", c.Name, msg)
		go c.SendMsg(Message{data: []byte(msg)})
	}
}

func (c *Client) SendMsg(msg Message) {
	_, err := c.conn.Write(msg.data)
	if err != nil {
		log.Printf("msg %s send err:%s \n", string(msg.data), err)
	}
}
func (c *Client) Recv() {
	for {
		buf := make([]byte, 256)
		_, err := c.conn.Read(buf)
		if err != nil {
			return
		}
		log.Printf("%s \n", string(buf))
	}

}
