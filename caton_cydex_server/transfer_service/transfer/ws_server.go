//this is websocket server
//wsserver 's function is that communicating with transfer node

package transfer

import (
	"cydex/transfer"//this is from src-->cydex
	"encoding/json"
	"fmt"
	clog "github.com/cihub/seelog"//log package 	
	"golang.org/x/net/websocket"//websockets package
	"log"
	"net/http"
	"time"
)


//===========================define some websocket-time arguments=================
type WSServerConfig struct {
	// 连接上来后多长时间没有响应就关闭连接
	ConnDeadline time.Duration
	// 心跳周期 in second
	KeepaliveInterval uint
	// 任务状态上传周期 in second
	TransferNotifyInterval uint
}
//set defualt arg of time
var DefaultConfig WSServerConfig = WSServerConfig{
	ConnDeadline:           10 * time.Second,
	KeepaliveInterval:      300,
	TransferNotifyInterval: 3,
}
//add stuffs to  a websocket server
//===================================define server config struct=============
type WSServer struct {//server config
	Version string
	config  *WSServerConfig//from above
	url     string
	port    int
}
//finish almost option unless version
//there, pass arg ,initial a wsserver ,if server config is empty ,the funtion  will pass
//defualt args
func NewWSServer(url string, port int, cfg *WSServerConfig) *WSServer {//from  above
	if cfg == nil {
		cfg = &DefaultConfig
	}
	//even if it lack  version it is senseless
	return &WSServer{
		config: cfg,
		url:    url,
		port:   port,
	}
}

//===================pass WSServer config paraments==================
//below,they are respectively set all kind of argvs
func (self *WSServer) SetConfig(cfg *WSServerConfig) {
	if cfg != nil {//configure config option
		self.config = cfg
	}
}
//config version
func (self *WSServer) SetVersion(v string) {
	self.Version = v
}

func (s *WSServer) Serve() {
//this is to start  a websocket server and prepare to receive info
	http.Handle(s.url, websocket.Handler(s.connHandle))//it is route(include url
	//connhandle is from below
	addr := fmt.Sprintf(":%d", s.port)//generate a addr :567
	log.Fatal(http.ListenAndServe(addr, nil))
	//funtion does ,although it is paraments//it is  likely 
	//launch a websocket server
}
//===============================deal with Node============
//this is route's maniputation function
func (s *WSServer) connHandle(ws *websocket.Conn) {
	var node *Node//this node define from  another file,why??because they  are from same dir
	var msgstring string
	var msg transfer.Message
	var rsp *transfer.Message
	var err error

	// cleanup
	defer func() {//having_non-name function
		if node != nil {
			clog.Warnf("Node disconnected: %+v", node)
			node.Close(true)
		}
	}()//what is brackets?


//=============================control websocket connetc================================
	ws.SetDeadline(time.Now().Add(s.config.ConnDeadline))
	//new a node ,but you need to konw it is special.it is used to deal with receive info
	node = NewNode(ws, s)//ws is websocket connect,s is websocket server
//nowtime ,node is not likely  downloading and uploading node 
	for {//keep loop to receive info from tranfer node
		if err = websocket.Message.Receive(ws, &msgstring); err != nil {
			log.Print(err)
			break
		}
		clog.Trace(msgstring)
		//trnsfer json to struct
		if err = json.Unmarshal([]byte(msgstring), &msg); err != nil {
			clog.Warnf("json unmarshal error:%s", err)
			continue
		}

		rsp, err = node.HandleMsg(&msg)//now msg is written
		if rsp != nil {//warn :it is not err
			clog.Trace(rsp)
			node.SendMessage(rsp)//figure out it send to where
		}
		if err != nil {
			clog.Errorf("%s handle msg error: %s", node, err)
			break
		}
	}
}
