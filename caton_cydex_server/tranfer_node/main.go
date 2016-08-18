package main
//this is really node that manage downloading and upload
import (
	"cydex"
	"cydex/transfer"
	"encoding/json"
	"fmt"
	clog "github.com/cihub/seelog"
	"net"
	"time"
)

func initLog() {
	cfgfiles := []string{
		"/opt/cydex/etc/ts_seelog.xml",
		"seelog.xml",
	}
	for _, file := range cfgfiles {
		logger, err := clog.LoggerFromConfigAsFile(file)
		if err != nil {
			println(err.Error())
			continue
		}
		clog.ReplaceLogger(logger)
		break
	}
}

type Application struct {
	n                *Node
	tasks            map[string]*Task
	new_task_chan    chan *Task
	end_task_chan    chan string
	notify_task_chan chan *transfer.TaskState
	task_notify_conn *net.UDPConn
}

func NewApplication() *Application {
	a := new(Application)
	a.tasks = make(map[string]*Task)
	a.new_task_chan = make(chan *Task)
	a.end_task_chan = make(chan string)
	a.notify_task_chan = make(chan *transfer.TaskState)
	return a
}

func (self *Application) taskServe() {
	for {
		select {
		case t := <-self.new_task_chan:
			self.tasks[t.Tid] = t
			go t.Run()
		case tid := <-self.end_task_chan:
			// t := self.tasks[tid]
			// if t != nil {
			// 	var stats []*transfer.TaskState
			// 	for _, s := range t.segs_state {
			// 		if s.State != "" {
			// 			stats = append(stats, s)
			// 		}
			// 	}
			// 	self.n.Notify(stats, 5*time.Second)
			// }
			delete(self.tasks, tid)
		// case <-time.Tick(3 * time.Second):
		// 	// 定时上报
		// 	var stats []*transfer.TaskState
		// 	for _, t := range self.tasks {
		// 		if t.cur_state == nil {
		// 			continue
		// 		}
		// 		stats = append(stats, t.cur_state)
		// 	}
		// 	self.n.Notify(stats, 5*time.Second)
		case state := <-self.notify_task_chan:
			var stats []*transfer.TaskState
			stats = append(stats, state)
			self.n.Notify(stats, 5*time.Second)
		}
	}
}

func (self *Application) taskEventServe() {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:3322")
	if err != nil {
		clog.Error(err)
		return
	}
	udp_conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		clog.Error(err)
		return
	}

	buf := make([]byte, 1024)
	for {
		n, err := udp_conn.Read(buf)
		if n == 0 || err != nil {
			clog.Error(n, err)
			break
		}
		var event Event
		err = json.Unmarshal(buf[:n], &event)
		if err != nil {
			clog.Error(err)
			continue
		}
		clog.Tracef("%+v", event)
		t, _ := self.tasks[event.Tid]
		if t != nil {
			t.event_chan <- &event
		}
	}
}

func (self *Application) RunNode(url string) error {
	self.n = NewNode()
	self.n.SetHandler(self)
	err := self.n.Dial(url)
	if err != nil {
		clog.Error(err)
		return err
	}
	self.n.Run()
	return err
}

func (self *Application) OnNodeMessage(n *Node, msg *transfer.Message, rsp *transfer.Message) {
	if msg.IsRsp() {
		return
	}

	switch msg.Cmd {
	case "uploadtask":
		port := getPort(cydex.UPLOAD)
		rsp.Rsp.UploadTask = &transfer.UploadTaskRsp{
			SidList: msg.Req.UploadTask.SidList,
			Port:    uint32(port),
		}
		// FIXME 目前是固定的
		rsp.Rsp.UploadTask.RecomendBitrate = 10 * 1024 * 1024
		// TODO 根据sid生成seg的url, 寻找可用路径等?
		for _, sid := range rsp.Rsp.UploadTask.SidList {
			url := fmt.Sprintf("file://%s/storage_1/%s", self.n.Nid, sid)
			rsp.Rsp.UploadTask.SidStorage = append(rsp.Rsp.UploadTask.SidStorage, url)
		}

		t := NewTask(msg, port, self.end_task_chan, self.notify_task_chan)
		self.new_task_chan <- t
	case "downloadtask":
		port := getPort(cydex.DOWNLOAD)

		rsp.Rsp.DownloadTask = &transfer.DownloadTaskRsp{
			SidList: msg.Req.DownloadTask.SidList,
			Port:    uint32(port),
		}
		// TODO 当前带宽和maxbitrate取小的那个, 供f2tp使用
		// FIXME 目前是固定的
		rsp.Rsp.DownloadTask.RecomendBitrate = 10 * 1024 * 1024

		t := NewTask(msg, port, self.end_task_chan, self.notify_task_chan)
		self.new_task_chan <- t
	default:
	}
}

func main() {
	initLog()

	app := NewApplication()
	go app.taskEventServe()
	go app.taskServe()
	app.RunNode("ws://127.0.0.1:12345/ts")
}
