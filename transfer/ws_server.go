package transfer

import (
	"cydex/transfer"
	"encoding/json"
	"fmt"
	clog "github.com/cihub/seelog"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"time"
)

type WSServerConfig struct {
	// 连接上来后多长时间没有响应就关闭连接
	ConnDeadline time.Duration
	// 心跳周期 in second
	KeepaliveInterval uint
	// 任务状态上传周期 in second
	TransferNotifyInterval uint
}

var DefaultConfig WSServerConfig = WSServerConfig{
	ConnDeadline:           10 * time.Second,
	KeepaliveInterval:      300,
	TransferNotifyInterval: 3,
}

type WSServer struct {
	Version string
	config  *WSServerConfig
	url     string
	port    int
}

func NewWSServer(url string, port int, cfg *WSServerConfig) *WSServer {
	if cfg == nil {
		cfg = &DefaultConfig
	}
	return &WSServer{
		config: cfg,
		url:    url,
		port:   port,
	}
}

func (self *WSServer) SetConfig(cfg *WSServerConfig) {
	if cfg != nil {
		self.config = cfg
	}
}

func (self *WSServer) SetVersion(v string) {
	self.Version = v
}

func (s *WSServer) Serve() {
	http.Handle(s.url, websocket.Handler(s.connHandle))
	addr := fmt.Sprintf(":%d", s.port)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func (s *WSServer) connHandle(ws *websocket.Conn) {
	var node *Node
	var msgstring string
	var msg transfer.Message
	var rsp *transfer.Message
	var err error

	// cleanup
	defer func() {
		if node != nil {
			clog.Warnf("Node disconnected: %+v", node)
			node.Close(true)
		}
	}()

	ws.SetDeadline(time.Now().Add(s.config.ConnDeadline))
	node = NewNode(ws, s)
	for {
		if err = websocket.Message.Receive(ws, &msgstring); err != nil {
			log.Print(err)
			break
		}
		clog.Trace(msgstring)
		if err = json.Unmarshal([]byte(msgstring), &msg); err != nil {
			clog.Warnf("json unmarshal error:%s", err)
			continue
		}

		rsp, err = node.HandleMsg(&msg)
		if rsp != nil {
			clog.Trace(rsp)
			node.SendMessage(rsp)
		}
		if err != nil {
			clog.Errorf("%s handle msg error: %s", node, err)
			break
		}
	}
}