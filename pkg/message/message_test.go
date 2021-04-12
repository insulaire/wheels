package message

import (
	"fmt"
	"io"
	"log"
	"net"
	"testing"
	"time"
)

func Test(t *testing.T) {

	go func() {
		lis, err := net.Listen("tcp4", "0.0.0.0:7777")
		if err != nil {
			return
		}
		conn, err := lis.Accept()
		if err != nil {
			return
		}
		pack := NewPack()
		for {

			bufHead := make([]byte, pack.GetHeaderLength())
			_, err = io.ReadFull(conn, bufHead)
			if err != nil {
				fmt.Println(err)
				return
			}

			msg, err := pack.UnPack(bufHead)
			if err != nil {
				log.Println(err)
				return
			}
			if msg.GetLength() > 0 {

				body := make([]byte, msg.GetLength())
				_, err = io.ReadFull(conn, body)
				if err != nil {
					fmt.Println(err)
					return
				}
				msg.SetBody(body)
				fmt.Println("abc")
				fmt.Println(string(msg.GetBody()))
			}

		}
	}()
	time.Sleep(3 * time.Second)
	conn, err := net.Dial("tcp4", "0.0.0.0:7777")
	if err != nil {
		return
	}
	pack := NewPack()
	msg := NewMessage(123, []byte("abcdefg"))
	msg1 := NewMessage(122, []byte("abcdefg"))
	buf1, err := pack.Pack(msg)
	if err != nil {
		return
	}
	buf2, err := pack.Pack(msg1)
	if err != nil {
		return
	}
	conn.Write(append(buf1, buf2...))
	select {}
}
