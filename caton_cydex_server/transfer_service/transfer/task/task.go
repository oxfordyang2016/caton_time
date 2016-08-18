package task

import (
	trans "./.."
	"cydex"
	"cydex/transfer"
	"errors"
	"fmt"
	clog "github.com/cihub/seelog"
	"github.com/pborman/uuid"
	"sync"
	"time"
)

const (
	TASK_TIMEOUT = 5 * time.Minute
)

var (
	TaskMgr *TaskManager
)

func init() {
	TaskMgr = NewTaskManager()
}

func GenerateTaskId() string {
	return uuid.New()
}

type TaskObserver interface {
	AddTask(t *Task)
	DelTask(t *Task)
	TaskStateNotify(t *Task, state *transfer.TaskState)
}

type UploadReq struct {
	*transfer.UploadTaskReq
}

// func NewUploadReq(r *transfer.UploadTaskReq, pid string) (req *UploadReq, err error) {
// 	if r == nil || pid == "" {
// 		return nil, errors.New("Invalid params")
// 	}
// 	req = &UploadReq{
// 		UploadTaskReq: r,
// 		Pid:           pid,
// 	}
// 	return req, nil
// }

type DownloadReq struct {
	*transfer.DownloadTaskReq
	// Pid             string
	FinishedSidList []string //已经下完的片段
	meta            interface{}
}

// func NewDownloadReq(r *transfer.DownloadTaskReq, pid string, finished_sid_list []string) (req *DownloadReq, err error) {
// 	if r == nil || pid == "" {
// 		return nil, errors.New("Invalid params")
// 	}
// 	req = &DownloadReq{
// 		DownloadTaskReq: r,
// 		Pid:             pid,
// 		FinishedSidList: finished_sid_list,
// 	}
// 	return req, nil
// }

type Task struct {
	TaskId      string // 任务ID
	Uid         string
	Pid         string
	Fid         string
	Type        int // 类型, U or D
	UploadReq   *UploadReq
	DownloadReq *DownloadReq
	SegsState   map[string]*transfer.TaskState
	CreateAt    time.Time
	UpdateAt    time.Time
	Nid         string
}

func NewTask(msg *transfer.Message) *Task {
	t := new(Task)
	t.CreateAt = time.Now()
	t.UpdateAt = time.Now()
	t.SegsState = make(map[string]*transfer.TaskState)

	var sid_list []string

	if msg.Cmd == "uploadtask" {
		req := msg.Req.UploadTask
		t.TaskId = req.TaskId
		t.Uid = req.Uid
		t.Pid = req.Pid
		t.Fid = req.Fid
		t.Type = cydex.UPLOAD
		sid_list = req.SidList
	} else {
		req := msg.Req.DownloadTask
		t.TaskId = req.TaskId
		t.Uid = req.Uid
		t.Pid = req.Pid
		t.Fid = req.Fid
		t.Type = cydex.DOWNLOAD
		sid_list = req.SidList
	}

	for _, sid := range sid_list {
		t.SegsState[sid] = &transfer.TaskState{
			TaskId: t.TaskId,
			Sid:    sid,
		}
	}
	return t
}

func (self *Task) IsDispatched() bool {
	return self.Nid != ""
}

func (self *Task) Stop() {
	clog.Infof("%s stop", self)
	if !self.IsDispatched() {
		return
	}

	const timeout = 5 * time.Second

	msg := transfer.NewReqMessage("", "stoptask", "", 0)
	msg.Req.StopTask = &transfer.StopTaskReq{
		TaskId: self.TaskId,
	}
	node := trans.NodeMgr.GetByNid(self.Nid)
	if node != nil {
		if _, err := node.SendRequestSync(msg, timeout); err != nil {
			return
		}
	}
}

func (self *Task) IsOver() bool {
	finished_segs := 0
	for _, state := range self.SegsState {
		if state.State == "interrupted" || state.State == "end" {
			finished_segs++
		}
	}
	if finished_segs == len(self.SegsState) {
		return true
	}
	return false
}

func (self *Task) String() string {
	var s string
	switch self.Type {
	case cydex.UPLOAD:
		s = "U"
	case cydex.DOWNLOAD:
		s = "D"
	default:
		s = "?"
	}
	return fmt.Sprintf("<Task(id:%s %s segs:%d)>", self.TaskId[:8], s, len(self.SegsState))
}

// 任务管理器
type TaskManager struct {
	scheduler    TaskScheduler
	lock         sync.Mutex
	tasks        map[string]*Task // task_id -> node
	observer_idx uint32
	observers    map[uint32]TaskObserver
	// task_state_routine_run bool
	// task_state_chan        chan *transfer.TaskState
}

func NewTaskManager() *TaskManager {
	t := new(TaskManager)
	t.tasks = make(map[string]*Task)
	t.observers = make(map[uint32]TaskObserver)
	// t.task_state_chan = make(chan *transfer.TaskState, 10)
	return t
}

func (self *TaskManager) SetScheduler(s TaskScheduler) {
	clog.Infof("Task Manager Set scheduler: %s", s.String())
	self.scheduler = s
}

func (self *TaskManager) AddObserver(o TaskObserver) (uint32, error) {
	if o == nil {
		return 0, errors.New("Observer is invalid")
	}
	defer func() {
		self.observer_idx++
		self.lock.Unlock()
	}()
	self.lock.Lock()
	self.observers[self.observer_idx] = o
	return self.observer_idx, nil
}

func (self *TaskManager) DelObserver(id uint32) {
	defer self.lock.Unlock()
	self.lock.Lock()
	delete(self.observers, id)
}

func (self *TaskManager) AddTask(t *Task) {
	// delete timeouted task
	for _, t := range self.tasks {
		if time.Since(t.UpdateAt) > TASK_TIMEOUT {
			self.DelTask(t.TaskId)
		}
	}

	self.lock.Lock()
	defer self.lock.Unlock()

	self.tasks[t.TaskId] = t
	for _, o := range self.observers {
		o.AddTask(t)
	}
	clog.Infof("Add task %s", t)
}

func (self *TaskManager) GetTask(taskid string) (t *Task) {
	defer self.lock.Unlock()
	self.lock.Lock()
	t, _ = self.tasks[taskid]
	return
}

func (self *TaskManager) DelTask(taskid string) {
	defer self.lock.Unlock()
	self.lock.Lock()
	t, _ := self.tasks[taskid]
	delete(self.tasks, taskid)
	if t != nil {
		for _, o := range self.observers {
			o.DelTask(t)
		}
	}
	clog.Infof("Del task(%s) %+v", taskid, t)
}

func (self *TaskManager) DispatchUpload(req *UploadReq, timeout time.Duration) (rsp *transfer.UploadTaskRsp, node *trans.Node, err error) {
	if req == nil || req.TaskId == "" {
		err = errors.New("Invalid req")
		return
	}
	if self.scheduler == nil {
		err = errors.New("No Scheduler, please set first")
		return
	}

	if node, err = self.scheduler.DispatchUpload(req); err != nil {
		return
	}
	if node == nil {
		err = errors.New("Scheduler didn't find suitable node!")
		return
	}

	var rsp_msg *transfer.Message
	msg := transfer.NewReqMessage("", "uploadtask", "", 0)
	msg.Req.UploadTask = req.UploadTaskReq
	if rsp_msg, err = node.SendRequestSync(msg, timeout); err != nil {
		// log
		clog.Warnf("%s send upload task failed, err:%s", node, err)
		return
	}
	// TODO 判断rsp_msg的Code值
	rsp = rsp_msg.Rsp.UploadTask

	t := NewTask(msg)
	t.UploadReq = req
	t.Nid = node.Nid
	TaskMgr.AddTask(t)
	return
}

func (self *TaskManager) DispatchDownload(req *DownloadReq, timeout time.Duration) (rsp *transfer.DownloadTaskRsp, node *trans.Node, err error) {
	if req == nil || req.TaskId == "" {
		err = errors.New("Invalid req")
		return
	}
	if self.scheduler == nil {
		err = errors.New("No Scheduler, please set first")
		return
	}

	if node, err = self.scheduler.DispatchDownload(req); err != nil {
		return
	}
	if node == nil {
		// Log
		return
	}

	var rsp_msg *transfer.Message
	msg := transfer.NewReqMessage("", "downloadtask", "", 0)
	msg.Req.DownloadTask = req.DownloadTaskReq
	if rsp_msg, err = node.SendRequestSync(msg, timeout); err != nil {
		// log
		return
	}
	rsp = rsp_msg.Rsp.DownloadTask

	t := NewTask(msg)
	t.DownloadReq = req
	t.Nid = node.Nid
	TaskMgr.AddTask(t)
	return
}

func (self *TaskManager) handleTaskState(state *transfer.TaskState) (err error) {
	// clog.Tracef("TaskState: %+v", state)
	t := self.GetTask(state.TaskId)
	// clog.Trace(t)
	if t != nil {
		t.UpdateAt = time.Now()
		s := t.SegsState[state.Sid]
		if s != nil {
			s.State = state.State
			s.TotalBytes = state.TotalBytes
			s.Bitrate = state.Bitrate
		}
	}

	self.lock.Lock()
	for _, o := range self.observers {
		o.TaskStateNotify(t, state)
	}
	self.lock.Unlock()

	if t != nil {
		if t.IsOver() {
			clog.Tracef("%s is over", t)
			self.DelTask(state.TaskId)
		}
	}
	return
}

// 侦听node接收的任务状态消息
func (self *TaskManager) ListenTaskState() {
	go func() {
		for states := range trans.NodeMgr.StateChan {
			for _, state := range states {
				self.handleTaskState(state)
			}
		}
	}()
}

// 根据相关信息停止任务
func (self *TaskManager) StopTasks(uid, pid string, typ int) {
	self.lock.Lock()
	defer self.lock.Unlock()

	for _, t := range self.tasks {
		if t.Uid == uid && t.Pid == pid && t.Type == typ {
			go t.Stop()
		}
	}
}

func (self *TaskManager) Scheduler() TaskScheduler {
	return self.scheduler
}
