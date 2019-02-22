package log

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

type FluentdClient struct {
	config  *Config
	conn    net.Conn
	ticker  *time.Ticker
	msgChan chan Message
	buffer  []byte
}

type Message struct {
	Tag    string      `json:"tag"`
	Time   int64       `json:"time"`
	Record interface{} `json:"record"`
}

func (m *Message) ToMsgpack() ([]byte, error) {
	record, err := json.Marshal(m.Record)
	return []byte(fmt.Sprintf("[\"%s\",%d,%s,null]", m.Tag, m.Time, record)), err
}

func NewFluentdClient(c *Config) *FluentdClient {
	if c == nil {
		return nil
	}

	cli := &FluentdClient{
		config:  c,
		ticker:  time.NewTicker(c.FlushTimeout),
		msgChan: make(chan Message, c.BufferLength),
		buffer:  make([]byte, 0, c.BufferLength),
	}
	err := cli.connect()
	if err != nil {
		return nil
	}
	go cli.loop()

	return cli
}

func (c *FluentdClient) connect() error {
	if c.conn != nil {
		return nil
	}
	var err error
	for i := 0; i < c.config.MaxConnTrial; i++ {
		c.conn, err = net.DialTimeout(
			"tcp",
			fmt.Sprintf("%s:%d", c.config.Host, c.config.Port),
			c.config.ConnTimeout,
		)
		if err == nil {
			return nil
		}
	}
	return nil
}

func (c *FluentdClient) Post(tag string, data interface{}, async bool) error {
	if async {
		tag = c.prependTag(tag)
		c.msgChan <- Message{Tag: tag, Time: time.Now().Unix(), Record: data}
	} else {
		tag = c.prependTag(tag)
		msg := Message{Tag: tag, Time: time.Now().Unix(), Record: data}
		raw, err := msg.ToMsgpack()
		if err != nil {
			return err
		}
		c.buffer = append(c.buffer, raw...)
		c.send()
	}
	return nil
}

func (c *FluentdClient) loop() {
	for {
		select {
		case msg := <-c.msgChan:
			raw, err := msg.ToMsgpack()
			if err != nil {
				log.Printf("[FluentdClient] failed to marshal message " + err.Error())
				continue
			}

			c.buffer = append(c.buffer, raw...)
			log.Printf("[FluentdClient] append\t%s\t%s", c.buffer, raw)
			if len(c.buffer) >= c.config.BufferLength {
				c.send()
			}
		case <-c.ticker.C:
			c.send()
		}
	}
}

func (c *FluentdClient) prependTag(tag string) string {
	if c.config.TagPrefix != "" {
		tag = c.config.TagPrefix + "." + tag
	}
	return tag
}

func (c *FluentdClient) send() error {
	if len(c.buffer) <= 0 {
		return errors.New("[FluentdClient] Buffer is empty")
	}

	if c.conn == nil {
		if err := c.connect(); err != nil {
			log.Printf("[FluentdClient] Not connected with fluentd " + err.Error())
			return errors.New("[FluentdClient] Not connected with fluentd")
		}
	}

	log.Printf("[FluentdClient] Send: %s\n", string(c.buffer))
	_, err := c.conn.Write(c.buffer)
	if err == nil {
		c.buffer = c.buffer[0:0]
	} else {
		log.Printf("[FluentdClient] failed to send message " + err.Error())
	}
	return nil
}
