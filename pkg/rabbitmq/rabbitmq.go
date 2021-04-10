package rabbitmq

import (
	"fmt"
	"log"
	"net/url"
	"wheels/pkg/pool"

	"github.com/streadway/amqp"
)

type RabbitMQConfig struct {
	UserName    string
	Password    string
	Host        string
	Port        int
	VirtualHost string
}

func (cof RabbitMQConfig) getConnString() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d//%s", url.QueryEscape(cof.UserName), url.QueryEscape(cof.Password), cof.Host, cof.Port, cof.VirtualHost)
}

type RabbitMQ struct {
	connPool  pool.Pool
	consumers []*Consumer
}
type consumeFn func([]byte) bool

type Consumer struct {
	name      string
	msg       <-chan amqp.Delivery
	consumeFn consumeFn
	ch        *amqp.Channel
	_closed   chan struct{}
}

func getConn(cof RabbitMQConfig) (*amqp.Connection, error) {
	return amqp.Dial(cof.getConnString())
}

func NewRabbitMQ(cof RabbitMQConfig) (*RabbitMQ, error) {
	//*amqp.Connection
	p := pool.NewPool(pool.PoolConfig{
		Min: 5,
		Max: 10,
		InitFn: func() interface{} {
			conn, err := getConn(cof)
			if err != nil {
				log.Panicln(err)
				return nil
			}
			return conn
		},
	})
	return &RabbitMQ{
		connPool:  p,
		consumers: []*Consumer{},
	}, nil
}

func (rb *RabbitMQ) NewConsumer(queueName string, consumeFn consumeFn) (*Consumer, error) {
	conn, _ := rb.connPool.Get().(*amqp.Connection)
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	err = ch.Qos(100, 0, false)
	if err != nil {
		return nil, err
	}
	c := &Consumer{
		consumeFn: consumeFn,
		name:      queueName,
		ch:        ch,
		_closed:   make(chan struct{}, 1),
	}
	rb.consumers = append(rb.consumers, c)
	if consumeFn != nil {
		msg, err := ch.Consume(queueName, "", false, false, false, false, nil)
		if err != nil {
			return nil, err
		}
		c.msg = msg
		go c.doConsume()
	}
	return c, nil
}

func (c *Consumer) GetMsg(limit int) []string {
	ans := []string{}
	for i := 0; i < limit; i++ {
		msg, ok, err := c.ch.Get(c.name, false)
		if err != nil {
			log.Panicln(err)
			return ans
		}
		if !ok {
			return ans
		}
		ans = append(ans, string(msg.Body))
	}
	return ans
}

func (c *Consumer) Close() {
	c._closed <- struct{}{}
}

func (c *Consumer) doConsume() {
	defer log.Println("Exit")
	if c.consumeFn == nil {
		log.Println("Consume method is nil")
	}
	for {
		select {
		case item := <-c.msg:
			//log.Printf("Consume Do:%s \n", string(item.Body))
			if c.consumeFn(item.Body) {
				item.Ack(false)
			}
		case <-c._closed:
			log.Println("Consume closed")
			return
		}
	}
}

func (c *Consumer) Publish(msg string) error {
	return c.ch.Publish("", c.name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msg),
	})
}
