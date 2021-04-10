package rabbitmq

import (
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"
)

func Test(t *testing.T) {

	//close := make(chan struct{})
	rb, err := NewRabbitMQ(RabbitMQConfig{
		UserName:    "admin",
		Password:    "admin",
		Host:        "localhost",
		Port:        5672,
		VirtualHost: "",
	})
	if err != nil {
		log.Panicln(err)
	}

	// _, err = rb.NewQueue("test", func(b []byte) bool {
	// 	//fmt.Println(string(b))
	// 	time.Sleep(time.Microsecond * 100)
	// 	return true
	// })
	// if err != nil {
	// 	log.Panicln(err)
	// }
	for i := 0; i < 7; i++ {
		_, err = rb.NewConsumer("test", func(b []byte) bool {
			fmt.Println(string(b))
			time.Sleep(time.Microsecond * 100)
			return true
		})
		if err != nil {
			log.Panicln(err)
		}
	}

	go func() {

		qe, err := rb.NewConsumer("test", nil)
		if err != nil {
			log.Panicln(err)
		}
		for i := 0; i < 10000; i++ {
			if err := qe.Publish(strconv.Itoa(i)); err != nil {
				log.Println(err)
				return
			}
			//time.Sleep(time.Second * 100)
			log.Println(i)
		}
		//close <- struct{}{}
	}()
	select {}
	//<-close

	//log.Println(qe.GetMsg(100))
}
