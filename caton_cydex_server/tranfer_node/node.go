package main

import (
	"cydex"
	"cydex/transfer"
	"encoding/json"
	"errors"
	"fmt"
	clog "github.com/cihub/seelog"
	"golang.org/x/net/websocket"
	"sync"
	"time"
)

type NodeHandler interface {
	// 信息收集回调
	// OnNodeRegister(req *transfer.RegisterReq)
	// OnNodeLogin(req *transfer.TransferReq)
	// OnNodeKeepalive(req *transfer.KeepaliveReq)
	// 消息回调
	OnNodeMessage(n *Node, msg *transfer.Message, rsp *transfer.Message)
	// node断开连接
	// OnNodeClose(n *Node)
}

type Node struct {
	Nid      string
	handler  NodeHandler
	Token    string
	seq_lock sync.Mutex
	seq      uint32
	rsp_chan chan *transfer.Message
	ws       *websocket.Conn
}

func NewNode() *Node {
	n := new(Node)
	n.rsp_chan = make(chan *transfer.Message)
	return n
}

func (self *Node) SetHandler(h NodeHandler) {
	self.handler = h
}

func (self *Node) Dial(url string) error {
	origin := "http://localhost"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		return err
	}
	self.ws = ws
	return err
}

func (self *Node) connHandle() {
	var msgstring string
	var msg transfer.Message
	// var rsp *transfer.Message
	var err error

	// cleanup
	defer func() {
	}()

	for {
		if err = websocket.Message.Receive(self.ws, &msgstring); err != nil {
			clog.Error(err)
			break
		}
		clog.Trace(msgstring)
		if err = json.Unmarshal([]byte(msgstring), &msg); err != nil {
			clog.Warnf("json unmarshal error:%s", err)
			continue
		}

		if msg.IsRsp() {
			self.rsp_chan <- &msg
		}
		if self.handler != nil {
			var rsp *transfer.Message
			if msg.IsReq() {
				rsp = msg.BuildRsp()
			}
			self.handler.OnNodeMessage(self, &msg, rsp)
			self.SendMessage(rsp)
		}
	}
}

func (self *Node) SendRequestSync(msg *transfer.Message, timeout time.Duration) (rsp *transfer.Message, err error) {
	if err = self.SendRequest(msg); err != nil {
		return
	}

	select {
	case rsp = <-self.rsp_chan:
		if rsp.Seq != msg.Seq || rsp.Cmd != msg.Cmd {
			err = fmt.Errorf("%s rsp is not match, %s %s", self, msg, rsp)
			rsp = nil
		}
	case <-time.After(timeout):
		err = fmt.Errorf("%s msg %s wait rsp timeout", self, msg.Cmd)
	}
	return
}

func (self *Node) SendRequest(msg *transfer.Message) error {
	if !msg.IsReq() {
		return errors.New("msg is not request")
	}
	self.seq_lock.Lock()
	msg.Seq = self.seq
	self.seq++
	self.seq_lock.Unlock()
	return self.SendMessage(msg)
}

func (self *Node) Register(timeout time.Duration) (rsp *transfer.Message, err error) {
	req := transfer.NewReqMessage("", "register", "", self.seq)
	req.Req.Register = &transfer.RegisterReq{
		MachineCode: "n1_machine_code",
	}
	rsp, err = self.SendRequestSync(req, timeout)
	return
}

func (self *Node) Login(timeout time.Duration) (rsp *transfer.Message, err error) {
	req := transfer.NewReqMessage(self.Nid, "login", "", self.seq)
	req.Req.Login = &transfer.LoginReq{
		Version:      "1.0.0-fake",
		OS:           "centos",
		TotalStorage: 50 * 1024 * 1024 * 1024,
		FreeStorage:  50 * 1024 * 1024 * 1024,
	}
	rsp, err = self.SendRequestSync(req, timeout)
	return
}

func (self *Node) Notify(stats []*transfer.TaskState, timeout time.Duration) (rsp *transfer.Message, err error) {
	if stats == nil || len(stats) == 0 {
		return nil, errors.New("state is empty")
	}
	req := transfer.NewReqMessage(self.Nid, "transfernotify", self.Token, self.seq)
	req.Req.TransferNotify = &transfer.TransferNotifyReq{
		TaskStateList: stats,
	}
	rsp, err = self.SendRequestSync(req, timeout)
	return
}

func (self *Node) Keepalive(timeout time.Duration) (rsp *transfer.Message, err error) {
	req := transfer.NewReqMessage(self.Nid, "keepalive", self.Token, self.seq)
	req.Req.Keepalive = &transfer.KeepaliveReq{
		CpuUsage:          21,
		TotalMem:          1 * 1024 * 1024 * 1024,
		FreeMem:           500 * 1024 * 1024,
		TotalStorage:      50 * 1024 * 1024 * 1024,
		FreeStorage:       50 * 1024 * 1024 * 1024,
		UploadBandwidth:   100 * 1024 * 1024,
		DownloadBandwidth: 100 * 1024 * 1024,
	}
	rsp, err = self.SendRequestSync(req, timeout)
	return
}

func (self *Node) SendMessage(msg *transfer.Message) error {
	if msg == nil {
		return errors.New("msg is nil")
	}
	if self.ws != nil {
		websocket.JSON.Send(self.ws, *msg)
	}
	return nil
}

func (self *Node) Run() {
	const (
		STATE_IDLE = iota
		STATE_REGISTER
		STATE_LOGIN
	)
	go self.connHandle()
	state := STATE_IDLE

	for {
		switch state {
		case STATE_IDLE:
			rsp, err := self.Register(10 * time.Second)
			if err != nil {
				clog.Error(err)
			} else {
				if rsp.Rsp.Code == cydex.OK {
					self.Nid = rsp.Rsp.Register.Tnid
					state = STATE_REGISTER
					break
				}
			}
			time.Sleep(3 * time.Second)
		case STATE_REGISTER:
			rsp, err := self.Login(10 * time.Second)
			if err != nil {
				clog.Error(err)
			} else {
				if rsp.Rsp.Code == cydex.OK {
					self.Token = rsp.Rsp.Login.Token
					state = STATE_LOGIN
					break
				}
			}
			time.Sleep(3 * time.Second)
		case STATE_LOGIN:
			rsp, err := self.Keepalive(10 * time.Second)
			rsp = rsp
			err = err
			time.Sleep(60 * time.Second)
		}
	}
}
