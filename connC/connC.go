package connC

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/PTFS/connsocket/models"
	"github.com/PTFS/connsocket/socket"
)

// ConnC client-conn
type ConnC struct {
	socket.Socket
	Port    string
	Addr    string
	ComChan chan models.Command
	rDo     reconnectDo
}

type reconnectDo func()

// NewConnC return instance of ConnC
func NewConnC(port string, addr string) (client *ConnC, err error) {
	c := ConnC{
		Port: port,
		Addr: addr,
	}

	conn, err := net.Dial("tcp", addr+":"+port)
	c.ComChan = make(chan models.Command, 16)
	c.ReadChan = make(chan []byte, 16)
	c.WriteChan = make(chan []byte, 16)
	c.CloseChan = make(chan bool, 1)
	c.Conn = conn
	c.rDo = func() {
		fmt.Println("This is test RDo")
	}

	if err != nil {
		log.Println(err)
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	c.Ctx = ctx
	c.Closef = cancel
	go c.ReadMsg()
	go c.WriteMsg()
	go c.handleReadChan(ctx, cancel)
	return &c, nil
}

// SetReconnect set the function when reconnect
func (c *ConnC) SetReconnect(fn reconnectDo) {
	c.rDo = fn
}

// Write write msg to writeChan
func (c *ConnC) Write(data *models.Report) error {
	val, err := data.String()
	if err != nil {
		return err
	}
	c.WriteChan <- val
	return nil
}

// GetCmdChan return CmdChan
func (c *ConnC) GetCmdChan() chan models.Command {
	return c.ComChan
}

// Close close ConnC
func (c *ConnC) Close() {
	c.Closef()
}

func (c *ConnC) handleReadChan(ctx context.Context, cancel context.CancelFunc) {
	for {
		select {
		case <-ctx.Done():
			return
		case data := <-c.ReadChan:
			// com := models.Command{}
			// err := json.Unmarshal(data, &com)
			var obj json.RawMessage
			com := models.Command{Content: &obj}
			err := json.Unmarshal(data, &com)
			com.Unmarshal(obj)
			if err == nil {
				c.ComChan <- com
			} else {
				log.Println("Receive Msg Error", err)
			}
		case <-c.CloseChan:
			go c.reconnect()
		}
	}
}

// reconnect to server
func (c *ConnC) reconnect() {
	for {
		conn, err := net.Dial("tcp", c.Addr+":"+c.Port)
		if err != nil {
			fmt.Println("Fatal error:", err.Error())
			time.Sleep(3 * time.Second)
		} else {
			c.Conn = conn
			ctx, cancel := context.WithCancel(context.Background())
			c.Ctx = ctx
			c.Closef = cancel
			go c.ReadMsg()
			go c.WriteMsg()
			go c.handleReadChan(ctx, cancel)
			c.rDo()
			return
		}
	}
}
