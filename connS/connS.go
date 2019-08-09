package connS

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	"github.com/LiMoMoMo/go-connSocket/models"
	"github.com/LiMoMoMo/go-connSocket/socket"
)

// RecieveMsg from client
type RecieveMsg struct {
	Rep  models.Report
	Conn *socket.Socket
}

// WriteMsg to client
type WriteMsg struct {
	Comm models.Command
	ID   string
}

// ConnS server-conn
type ConnS struct {
	Port         string
	Addr         string
	server       net.Listener
	msgChan      chan models.Report
	readmsgChan  chan RecieveMsg
	writemsgChan chan WriteMsg
	cancel       context.CancelFunc
	ctx          context.Context

	container map[string]*socket.Socket
}

// NewConnS return instance of ConnS
func NewConnS(port string, addr string) *ConnS {
	c := ConnS{
		Port:         port,
		Addr:         addr,
		msgChan:      make(chan models.Report, 16),
		readmsgChan:  make(chan RecieveMsg, 16),
		writemsgChan: make(chan WriteMsg, 16),
		container:    make(map[string]*socket.Socket),
	}
	ctx, cancel := context.WithCancel(context.Background())
	c.ctx = ctx
	c.cancel = cancel
	return &c
}

// GetRepChan return chann instance
func (c *ConnS) GetRepChan() chan models.Report {
	return c.msgChan
}

// WriteTo write msg to client.
func (c *ConnS) WriteTo(id string, data *models.Command) {
	wmsg := WriteMsg{
		Comm: *data,
		ID:   id,
	}
	c.writemsgChan <- wmsg
}

// Stop stop server
func (c *ConnS) Stop() {
	c.server.Close()
	c.cancel()
}

// Start start server
func (c *ConnS) Start() error {
	server, err := net.Listen("tcp", c.Addr+":"+c.Port)
	if err != nil {
		return err
	}
	c.server = server
	fmt.Println("Start listening on Tcp " + c.Addr + ":" + c.Port)

	go func() {
		for {
			conn, err := c.server.Accept()
			if err != nil {
				fmt.Println(err)
				break
			}
			soc := socket.Socket{
				Conn:      conn,
				ReadChan:  make(chan []byte, 16),
				WriteChan: make(chan []byte, 16),
				CloseChan: make(chan bool, 1),
			}
			ctx, cancel := context.WithCancel(context.Background())
			soc.Ctx = ctx
			soc.Closef = cancel
			go soc.ReadMsg()
			go soc.WriteMsg()
			go func() {
				for {
					select {
					case val := <-soc.ReadChan:
						//
						var obj json.RawMessage
						re := models.Report{Content: &obj}
						err := json.Unmarshal(val, &re)
						re.Unmarshal(obj)
						//
						if err == nil {
							smg := RecieveMsg{
								Rep:  re,
								Conn: &soc,
							}
							c.readmsgChan <- smg
						}
					case <-ctx.Done():
						re := models.Report{
							Type: models.Type_Logout,
						}
						smg := RecieveMsg{
							Rep:  re,
							Conn: &soc,
						}
						c.readmsgChan <- smg
						fmt.Println("Quit ConnS Loop Socket")
						return
					}
				}
			}()
		}
	}()

	go c.run()

	return nil
}

func (c *ConnS) run() {
	for {
		select {
		case <-c.ctx.Done():
			for k, v := range c.container {
				v.Closef()
				delete(c.container, k)
				v = nil
			}
			return
		case smg := <-c.readmsgChan:
			switch smg.Rep.Type {
			case models.Type_Register:
				//
				// login := models.Login{}
				// err := json.Unmarshal([]byte(smg.Rep.Content), &login)
				login := (smg.Rep.Content).(*models.Register)
				if _, ok := c.container[login.ID]; ok {
					fmt.Println("已经登陆,", login.ID)
				} else {
					c.container[login.ID] = smg.Conn
					fmt.Println("登陆成功,", login.ID)
				}

			case models.Type_Logout:
				for k, v := range c.container {
					if v == smg.Conn {
						fmt.Println("退出登陆,", k)
						delete(c.container, k)
						v = nil
						smg.Rep.Content = models.Register{ID: k}
						break
					}
				}
			}
			c.msgChan <- smg.Rep
		case msg := <-c.writemsgChan:
			if conn, ok := c.container[msg.ID]; ok {
				val, err := msg.Comm.String()
				if err == nil {
					conn.WriteChan <- val
				}
			} else {
				fmt.Println("写消息失败")
			}
		}
	}
}
