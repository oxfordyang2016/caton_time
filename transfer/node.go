package transfer

import (
	"./models"
	"cydex"
	"cydex/transfer"
	"errors"
	"fmt"
	clog "github.com/cihub/seelog"
	"github.com/pborman/uuid"
	"golang.org/x/net/websocket"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	// 响应消息的超时时间,如果超过则需要删除
	MESSAGE_TIMEOUT = 5 * time.Minute
)

var (
	NodeMgr *NodeManager
)

func init() {
	NodeMgr = NewNodeManager()
}

// 注册Node, 分配tnid
func registerNode(req *transfer.RegisterReq) (code int, tnid string, err error) {
	var node *models.Node
	node, err = models.GetNodeByMachineCode(req.MachineCode)
	if node != nil {
		// 已经存在
		return cydex.OK, node.Nid, err
	}
	tnid = uuid.New()
	if _, err = models.CreateNode(req.MachineCode, tnid); err != nil {
		code = cydex.ErrInnerServer
		return
	}
	code = cydex.OK
	return
}

type NodeObserver interface {
	AddNode(n *Node)
	UpdateNode(n *Node, req *transfer.KeepaliveReq)
	DelNode(n *Node)
}

// TransferNode管理
type NodeManager struct {
	StateChan chan []*transfer.TaskState
	mux       sync.Mutex
	id_map    map[string]*Node // id->node
	observers []NodeObserver
}

func NewNodeManager() *NodeManager {
	nm := new(NodeManager)
	nm.id_map = make(map[string]*Node)
	nm.StateChan = make(chan []*transfer.TaskState)
	return nm
}

func (self *NodeManager) IsOnline(nid string) bool {
	defer self.mux.Unlock()
	self.mux.Lock()
	_, ok := self.id_map[nid]
	return ok
}

func (self *NodeManager) GetByNid(id string) *Node {
	defer self.mux.Unlock()
	self.mux.Lock()
	node, _ := self.id_map[id]
	return node
}

func (self *NodeManager) AddNode(node *Node) {
	clog.Infof("Node Add: %+v\n", node)
	defer self.mux.Unlock()
	self.mux.Lock()
	self.id_map[node.Nid] = node
	for _, o := range self.observers {
		o.AddNode(node)
	}
}

func (self *NodeManager) DelNode(nid string) {
	clog.Infof("Node Delete: %s\n", nid)
	defer self.mux.Unlock()
	self.mux.Lock()
	node, ok := self.id_map[nid]
	if ok {
		delete(self.id_map, nid)
		for _, o := range self.observers {
			o.DelNode(node)
		}
	}
}

func (self *NodeManager) AddObserver(observer NodeObserver) {
	if observer == nil {
		return
	}
	defer self.mux.Unlock()
	self.mux.Lock()
	self.observers = append(self.observers, observer)
}

type NodeInfo struct {
	Version           string
	NetAddr           string
	OS                string
	NetSpeed          uint32
	Storage           []*transfer.StorageInfo
	TotalStorage      uint64
	FreeStorage       uint64
	CpuUsage          uint32
	TotalMem          uint64
	FreeMem           uint64
	UploadBandwidth   uint64
	DownloadBandwidth uint64
}

type TimeMessage struct {
	*transfer.Message
	ts time.Time
}

func NewTimeMessage(msg *transfer.Message) *TimeMessage {
	return &TimeMessage{
		msg, time.Now(),
	}
}

// TransferNode
type Node struct {
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
	server   *WSServer
	closed   bool
}

func NewNode(ws *websocket.Conn, server *WSServer) *Node {
	n := new(Node)
	n.server = server
	n.SetWSConn(ws)
	// n.rsp_chan = make(chan *TimeMessage)
	n.rsp_sem = make(chan int)
	n.rsp_msgs = make(map[uint32]*TimeMessage)
	return n
}

func (self *Node) Verify(nid, token string) bool {
	return self.Token == token && self.Nid == nid
}

func (self *Node) SetWSConn(ws *websocket.Conn) {
	self.ws = ws
	if ws != nil {
		addr := ws.Request().RemoteAddr
		host, _, err := net.SplitHostPort(addr)
		if err == nil {
			self.Host = host
		}
	}
}

func (self *Node) Update(update_login_time bool) {
	if !self.IsLogined() {
		return
	}
	if self.server != nil && self.ws != nil {
		t := time.Now().Add(time.Duration(self.server.config.KeepaliveInterval) * time.Second * 3)
		self.ws.SetDeadline(t)
	}
	if update_login_time && self.Node != nil {
		self.Node.UpdateLoginTime(time.Now())
	}
}

func (self *Node) HandleMsg(msg *transfer.Message) (rsp *transfer.Message, err error) {
	if msg.IsReq() {
		rsp = msg.BuildRsp()
		rsp.Rsp.Code = cydex.OK
		if msg == nil {
			rsp.Rsp.Code = cydex.ErrInvalidParam
			rsp.Rsp.Reason = "Invalid Param"
			return
		}
	}

	lower_cmd := strings.ToLower(msg.Cmd)
	if msg.IsReq() {
		switch lower_cmd {
		case "register":
			err = self.handleRegister(msg, rsp)
		case "login":
			err = self.handleLogin(msg, rsp)
		case "keepalive":
			err = self.handleKeepAlive(msg, rsp)
		case "transfernotify":
			err = self.handleTransferNotify(msg, rsp)
		default:
			rsp.Rsp.Code = cydex.ErrInvalidParam
			rsp.Rsp.Reason = fmt.Sprintf("Unsupport command %s", msg.Cmd)
		}
		if err == nil {
			self.Update(false)
		}
	} else {
		self.rsp_lock.Lock()
		self.rsp_msgs[msg.Seq] = NewTimeMessage(msg)
		self.rsp_lock.Unlock()
		self.rsp_sem <- 1
	}
	return
}

func (self *Node) handleRegister(msg, rsp *transfer.Message) (err error) {
	if msg.Req == nil || msg.Req.Register == nil {
		err = fmt.Errorf("Invalid Param")
		rsp.Rsp.Code = cydex.ErrInvalidParam
		rsp.Rsp.Reason = err.Error()
		return
	}
	code, tnid, err := registerNode(msg.Req.Register)
	if code == cydex.OK && tnid != "" {
		rsp.Rsp.Code = code
		rsp.Rsp.Register = &transfer.RegisterRsp{
			Tnid: tnid,
		}
	} else {
		err = errors.New("Register node failed")
	}
	return
}

func (self *Node) IsRegisted() bool {
	return self.Node != nil
}

func (self *Node) IsLogined() bool {
	return self.Token != ""
}

func (self *Node) handleLogin(msg, rsp *transfer.Message) (err error) {
	if msg.Req == nil || msg.Req.Login == nil {
		err = fmt.Errorf("Invalid Param")
		rsp.Rsp.Code = cydex.ErrInvalidParam
		rsp.Rsp.Reason = err.Error()
		return
	}

	nid := msg.From
	n := NodeMgr.GetByNid(nid)
	if n != nil {
		// node with same nid is logined, and should kickout
		n.Close(true)
	}
	if self.Node, err = models.GetNode(nid); err != nil {
		return
	}
	if self.Node == nil {
		err = fmt.Errorf("%s is not registed", nid)
		rsp.Rsp.Code = cydex.ErrInvalidParam
		rsp.Rsp.Reason = err.Error()
		return
	}

	self.Token = uuid.New()
	self.Info.Version = msg.Req.Login.Version
	self.Info.NetAddr = msg.Req.Login.NetAddr
	self.Info.OS = msg.Req.Login.OS
	self.Info.NetSpeed = msg.Req.Login.NetSpeed
	self.Info.Storage = msg.Req.Login.Storage
	self.Info.TotalStorage = msg.Req.Login.TotalStorage
	self.Info.FreeStorage = msg.Req.Login.FreeStorage
	self.Info.CpuUsage = msg.Req.Login.CpuUsage
	self.Info.TotalMem = msg.Req.Login.TotalMem
	self.Info.FreeMem = msg.Req.Login.FreeMem
	self.Info.UploadBandwidth = msg.Req.Login.UploadBandwidth
	self.Info.DownloadBandwidth = msg.Req.Login.DownloadBandwidth
	self.login_at = time.Now()
	self.Update(true)

	NodeMgr.AddNode(self)

	var (
		alive_interval           uint32
		transfer_notify_interval uint32
		version                  string
	)
	if self.server != nil {
		alive_interval = uint32(self.server.config.KeepaliveInterval)
		transfer_notify_interval = uint32(self.server.config.TransferNotifyInterval)
		version = self.server.Version
	}

	t, _ := time.Now().MarshalText()
	rsp.Rsp.Login = &transfer.LoginRsp{
		Token:                  self.Token,
		ZoneId:                 self.Zid,
		AliveInterval:          alive_interval,
		TransferNotifyInterval: transfer_notify_interval,
		Time:    string(t),
		Version: version,
	}
	return
}

func (self *Node) handleKeepAlive(msg, rsp *transfer.Message) (err error) {
	if !self.Verify(msg.From, msg.Token) {
		err = fmt.Errorf("%s verify failed, token:%s, remote:[nid:%s, token:%s]", self, self.Token, msg.From, msg.Token)
		rsp.Rsp.Code = cydex.ErrInvalidLicense
		rsp.Rsp.Reason = err.Error()
	}

	return
}

func (self *Node) handleTransferNotify(msg, rsp *transfer.Message) (err error) {
	if !self.Verify(msg.From, msg.Token) {
		err = fmt.Errorf("%s verify failed, token:%s, remote:[nid:%s, token:%s]", self, self.Token, msg.From, msg.Token)
		rsp.Rsp.Code = cydex.ErrInvalidLicense
		rsp.Rsp.Reason = err.Error()
	}

	if msg.Req == nil || msg.Req.TransferNotify == nil {
		rsp.Rsp.Code = cydex.ErrInvalidParam
		rsp.Rsp.Reason = "Invalid Param"
		return
	}

	NodeMgr.StateChan <- msg.Req.TransferNotify.TaskStateList
	// for _, r := range msg.Req.TransferNotify.TaskStateList {
	// 	NodeMgr.StateChan <- r
	// }

	return
}

func (self *Node) SendRequest(msg *transfer.Message) error {
	if !msg.IsReq() {
		return errors.New("msg is not request")
	}
	return self.SendMessage(msg)
}

func (self *Node) SendMessage(msg *transfer.Message) error {
	if msg == nil {
		return errors.New("msg is nil")
	}
	self.lock.Lock()
	defer self.lock.Unlock()

	if msg.IsReq() {
		msg.Seq = self.seq
		self.seq++
	}
	if self.closed {
		return fmt.Errorf("%s send msg failed because closed", self)
	}
	if self.ws != nil {
		websocket.JSON.Send(self.ws, *msg)
	}
	return nil
}

// 同步获取消息
func (self *Node) SendRequestSync(msg *transfer.Message, timeout time.Duration) (rsp *transfer.Message, err error) {

	if err = self.SendRequest(msg); err != nil {
		return
	}

	var alive bool

	select {
	case _, alive = <-self.rsp_sem:
		if !alive {
			err = fmt.Errorf("%s rsp_sem closed, maybe disconnected")
			return nil, err
		}

		self.rsp_lock.Lock()
		defer self.rsp_lock.Unlock()
		time_msg, _ := self.rsp_msgs[msg.Seq]
		if time_msg != nil && time_msg.Message != nil {
			rsp = time_msg.Message
		}
		delete(self.rsp_msgs, msg.Seq)

		// 删除超时的响应消息
		for _, m := range self.rsp_msgs {
			if time.Since(m.ts) >= MESSAGE_TIMEOUT {
				delete(self.rsp_msgs, m.Seq)
			}
		}
	case <-time.After(timeout):
		err = fmt.Errorf("%s msg %s wait rsp timeout", self, msg.Cmd)
	}
	return
}

func (self *Node) Close(close_conn bool) {
	self.lock.Lock()
	self.closed = true
	if self.rsp_sem != nil {
		close(self.rsp_sem)
	}
	if close_conn && self.ws != nil {
		self.ws.Close()
	}
	self.lock.Unlock()

	if self.Node != nil {
		NodeMgr.DelNode(self.Nid)
		self.Node.UpdateLogoutTime(time.Now())
	}
}

func (self *Node) String() string {
	nid := ""
	if self.Node != nil {
		nid = self.Node.Nid
	}
	return fmt.Sprintf("<Node(%s %s)>", nid, self.Host)
}

func (self *Node) OnlineDuration() time.Duration {
	if !self.IsLogined() {
		return time.Duration(0)
	}
	return time.Since(self.login_at)
}
