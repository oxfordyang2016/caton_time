//this is websocket server
//wsserver 's function is that communicating with transfer node

package transfer

import (
	"cydex/transfer" //this is from src-->cydex
	"encoding/json"
	"fmt"
	clog "github.com/cihub/seelog" //log package
	"golang.org/x/net/websocket"   //websockets package
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
type WSServer struct { //server config
	Version string
	config  *WSServerConfig //from above
	url     string
	port    int
}

//finish almost option unless version
//there, pass arg ,initial a wsserver ,if server config is empty ,the funtion  will pass
//defualt args
func NewWSServer(url string, port int, cfg *WSServerConfig) *WSServer { //from  above
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
	/*
			type WSServer struct {//server config
			Version string
			config  *WSServerConfig//from above
			url     string
			port    int
		}
	*/
	if cfg != nil { //configure config option
		self.config = cfg
	}
}

//config version
func (self *WSServer) SetVersion(v string) {
	self.Version = v
}

func (s *WSServer) Serve() {
	//this is to start  a websocket server and prepare to receive info
	http.Handle(s.url, websocket.Handler(s.connHandle)) //it is route(include url
	//connhandle is from below
	addr := fmt.Sprintf(":%d", s.port) //generate a addr :567
	log.Fatal(http.ListenAndServe(addr, nil))
	//funtion does ,although it is paraments//it is  likely
	//launch a websocket server
}

//===============================deal with Node============
//this is route's maniputation function
func (s *WSServer) connHandle(ws *websocket.Conn) {

	//when pass handles  argvs,it will generate a Connect
	var node *Node //define a var
	//this node define from  another file,why??

	//because they  are from same dir-->node.go
	/*
			type Node struct {
			//from tranfer/models
			*models.Node

			// 运行时数据
			Host  string
			Token string
			Info  NodeInfo

			// private
			// alive_interval uint32
			login_at time.Time
			ws       *websocket.Conn
			lock     sync.Mutex
			seq      uint32
			// rsp_chan chan *TimeMessage
			rsp_sem  chan int
			rsp_lock sync.Mutex
			rsp_msgs map[uint32]*TimeMessage
			server   *WSServer//ws
			closed   bool
		}

	*/
	var msgstring string
	var msg transfer.Message
	var rsp *transfer.Message

	var err error

	// cleanup
	defer func() { //having_non-name function
		if node != nil {
			clog.Warnf("Node disconnected: %+v", node)
			node.Close(true)
		}
	}() //what is brackets?it is that donot pass argvs

	//=============================control websocket connetc================================
	ws.SetDeadline(time.Now().Add(s.config.ConnDeadline))
	//new a node ,but you need to konw it is special.it is used to deal with receive info
	node = NewNode(ws, s) //ws is websocket connect,s is websocket server
	//nowtime ,node is not likely  downloading and uploading node
	//this node is used to hold some server important info

	////wsserver instru->https://github.com/oxfordyang2016/caton_time/blob/master/images/wsserverinstru.jpg

	/*
	   func NewNode(ws *websocket.Conn, server *WSServer) *Node {
	   	n := new(Node)//it is allocate memory and initiaze it
	   	n.server = server//node 's server
	   	n.SetWSConn(ws)//from below
	   	// n.rsp_chan = make(chan *TimeMessage)
	   	n.rsp_sem = make(chan int)//make a channel that can send ,recieve and deliver a int
	   	//make return  a same type other than pointer
	   	n.rsp_msgs = make(map[uint32]*TimeMessage)//in bracket ,it is map type
	   	return n
	   =================================================================
	   	func (self *Node) SetWSConn(ws *websocket.Conn) {
	   	self.ws = ws
	   	if ws != nil {
	   		addr := ws.Request().RemoteAddr//for instance:resolve addr 192,168.0.21:90
	   		//get node ip
	   		host, _, err := net.SplitHostPort(addr)
	   		if err == nil {
	   			self.Host = host
	   		}


	   	}
	   }

	   }

	*/
	for { //keep loop to receive info from tranfer node
		if err = websocket.Message.Receive(ws, &msgstring); err != nil {
			//ws is server_node's connection
			log.Print(err)
			break
		}
		clog.Trace(msgstring)
		//trnsfer json to struct
		if err = json.Unmarshal([]byte(msgstring), &msg); err != nil { //var msg transfer.Message
			/*
							type Message struct {
					Base
					Req *Request  `json:"req,omitempty"`
					Rsp *Response `json:"rsp,omitempty"`
				}
			*/
			clog.Warnf("json unmarshal error:%s", err)
			continue
		}

		rsp, err = node.HandleMsg(&msg) //now msg is written
		//now time node has been written

		//there is node handler   to do somthings accoding to resloved json-struct info
		//func (self *Node) HandleMsg(msg *transfer.Message) (rsp *transfer.Message, err error) {
		//transfer from cydex/tranfer
		//this function return many values and format is caution
		if msg.IsReq() { //msg has been resolved to msg struct

			rsp = msg.BuildRsp() //msg is infomation by resolved

			/*func (m *Message) BuildRsp() *Message {
			  	return NewRspMessage("", m.Cmd, "", m.Seq)
			  }
			*/
			/*
							  func NewRspMessage(from, cmd, token string, seq uint32) *Message {
					return NewMessage(from, cmd, token, "rsp", seq)
				}
			*/
			/*
			   //generate a new message struct initialization
			   func NewMessage(from, cmd, token, typ string, seq uint32) *Message {
			   	m := new(Message)//return a pointer
			   	//intilizia m and return it
			   	m.From = from
			   	m.Cmd = cmd
			   	m.Token = token
			   	m.Seq = seq
			   	m.Type = typ

			   	switch typ {
			   	case "req":
			   		m.Req = new(Request)
			   	case "rsp":
			   		m.Rsp = new(Response)
			   	default:
			   		return nil
			   	}
			   	return m
			   }
			*/
			//if the info received is request,it will build response
			rsp.Rsp.Code = cydex.OK
			//there logic is not starnge

			//it is likely erorr !!!!!!
			if msg == nil { //if receive info is empty
				rsp.Rsp.Code = cydex.ErrInvalidParam
				rsp.Rsp.Reason = "Invalid Param"
				return
			}
		}

		if rsp != nil { //warn :it is not err
			clog.Trace(rsp)
			node.SendMessage(rsp) //figure out it send to where
			/*
							func (self *Node) SendMessage(msg *transfer.Message) error {
					if msg == nil {
						return errors.New("msg is nil")
					}
					self.lock.Lock()
					defer self.lock.Unlock()

					if msg.IsReq() {//generate msg seq
						msg.Seq = self.seq
						self.seq++
					}
					if self.closed {
						return fmt.Errorf("%s send msg failed because closed", self)
					}
					if self.ws != nil {//sendinfo
						websocket.JSON.Send(self.ws, *msg)
					}
					return nil
				}
			*/
		}
		if err != nil {
			clog.Errorf("%s handle msg error: %s", node, err)
			break
		}
	}
}
