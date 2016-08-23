package main
//nowtime,task.go is different from ts-tranfer-task-task.go
import (
	"cydex"
	"cydex/transfer"
	"fmt"
)

var dl_port = 18766

func getPort(typ int) int {
	switch typ {
	case cydex.UPLOAD:
		return 18765
	case cydex.DOWNLOAD:
		dl_port++
		return dl_port
	}
	return 18800
}

// f2tp通过回调上来的事件信息
type Event struct {
	Tid   string `json:"taskid"`
	Sid   string `json:"sid"`
	Event int    `json:"event"`
	Seq   uint32 `json:"seq"`
}

type Task struct {//why sometime ,it is struct ,but somtime it is pointer
	//thee all kinds of channel , i donnot understand
	Tid  string
	Pid  string
	Uid  string
	Fid  string
	Type int
	Port int
	// State           transfer.TaskState
	UploadTask      *transfer.UploadTaskReq//cydex transfer
	DownloadTask    *transfer.DownloadTaskReq
	FinishedNumSegs uint
	segs_state      map[string]*transfer.TaskState//segment state
	cur_state       *transfer.TaskState//what???
	event_chan      chan *Event//why need it

	over_chan   chan string
	notify_chan chan *transfer.TaskState
}

func NewTask(msg *transfer.Message, port int, over_chan chan string, notify_chan chan *transfer.TaskState) *Task {
	t := new(Task)
	t.over_chan = over_chan
	t.notify_chan = notify_chan
	t.segs_state = make(map[string]*transfer.TaskState)
	t.Port = port
	t.event_chan = make(chan *Event)

	var sid_list []string

	if msg.Cmd == "uploadtask" {
		req := msg.Req.UploadTask
		t.Tid = req.TaskId
		t.Uid = req.Uid
		t.Fid = req.Fid
		t.Type = cydex.UPLOAD
		t.UploadTask = req
		sid_list = req.SidList
	} else {
		req := msg.Req.DownloadTask
		t.Tid = req.TaskId
		t.Uid = req.Uid
		t.Fid = req.Fid
		t.Type = cydex.DOWNLOAD
		t.DownloadTask = req
		sid_list = req.SidList
	}

	for _, sid := range sid_list {
		t.segs_state[sid] = &transfer.TaskState{
			TaskId: t.Tid,
			Sid:    sid,
		}
	}
	return t
}

func (self *Task) Run() {
	self.RunF2tpServer()
	//通知app
	self.over_chan <- self.Tid
}

func (self *Task) IsOver() bool {
	finished_segs := 0
	for _, state := range self.segs_state {
		if state.State == "interrupted" || state.State == "end" {
			finished_segs++
		}
	}
	if finished_segs == len(self.segs_state) {
		return true
	}
	return false
}

func (self *Task) String() string {
	return fmt.Sprintf("<Task(%s type:%d port:%d)>", self.Tid[:8], self.Type, self.Port)
}
