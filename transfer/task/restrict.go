package task

import (
	trans "./../"
	"cydex"
	"cydex/transfer"
	"fmt"
	"sync"
	"time"
)

const (
	TASK_RESTRICT_NONE   = iota // 无约束
	TASK_RESTRICT_BY_PID        // 按照PkgId约束, 分配在同一个Node上
	TASK_RESTRICT_BY_FID        // 按照FileId约束, 分配在同一个Node上
)

const (
	DEFAULT_RESOUCE_EXPIRE = 5 * time.Minute
)

type GetIdFunc func(req *UploadReq) string

func GetPid(req *UploadReq) string {
	return req.Pid
}

func GetFid(req *UploadReq) string {
	return req.Fid
}

// 有约束的上传任务分配器
type RestrictUploadScheduler struct {
	restrict_mode int
	resource      *XidResource
	getId         GetIdFunc
	lock          sync.Mutex
}

func NewRestrictUploadScheduler(restrict_mode int) *RestrictUploadScheduler {
	n := new(RestrictUploadScheduler)
	n.resource = NewXidResource()
	if err := n.SetRestrict(restrict_mode); err != nil {
		return nil
	}
	return n
}

func (self *RestrictUploadScheduler) SetRestrict(mode int) error {
	defer self.lock.Unlock()
	self.lock.Lock()

	if self.restrict_mode == mode {
		return nil
	}
	switch mode {
	case TASK_RESTRICT_NONE:
		self.getId = nil
	case TASK_RESTRICT_BY_PID:
		self.getId = GetPid
	case TASK_RESTRICT_BY_FID:
		self.getId = GetFid
	default:
		return fmt.Errorf("Unknown restrict mode %d", mode)
	}

	self.restrict_mode = mode
	self.resource.Reset()
	return nil
}

// implement UploadScheduler
func (self *RestrictUploadScheduler) DispatchUpload(req *UploadReq) (n *trans.Node, err error) {
	if self.resource == nil || self.getId == nil {
		return nil, nil
	}

	defer self.lock.Unlock()
	self.lock.Lock()

	self.resource.DelExpired()
	r := self.resource.Get(self.getId(req))
	if r != nil && r.node != nil {
		if r.node.Info.FreeStorage >= req.Size {
			r.Update()
			n = r.node
		}
	}
	return
}

// implement task observer
func (self *RestrictUploadScheduler) AddTask(t *Task) {
	if t.Type != cydex.UPLOAD || t.UploadReq == nil {
		return
	}

	defer self.lock.Unlock()
	self.lock.Lock()

	xid := self.getId(t.UploadReq)
	self.resource.Add(xid, t.Node, DEFAULT_RESOUCE_EXPIRE)
}

func (self *RestrictUploadScheduler) DelTask(t *Task) {
	// Do nothing
}

func (self *RestrictUploadScheduler) TaskStateNotify(t *Task, state *transfer.TaskState) {
	if t == nil {
		return
	}
	if t.Type != cydex.UPLOAD || t.UploadReq == nil {
		return
	}

	defer self.lock.Unlock()
	self.lock.Lock()

	// xid := self.getId(t.UploadReq)
	// self.resource.Update(xid, t.Node)
}

type Resource struct {
	node      *trans.Node
	timestamp time.Time
	expire    time.Duration
}

func (self *Resource) Update() {
	self.timestamp = time.Now()
}

// 各约束都是Xid->Node的关系
// 这些关系会过期, 例如pkg或者file传输完毕,这里不侦测具体是否完毕,通过超时来释放资源
type XidResource struct {
	maps map[string]*Resource
}

func NewXidResource() *XidResource {
	n := new(XidResource)
	n.maps = make(map[string]*Resource)
	return n
}

func (self *XidResource) Add(xid string, n *trans.Node, expire time.Duration) {
	if n == nil {
		return
	}
	self.Update(xid, n)
	self.DelExpired()

	// add new
	r := &Resource{
		node:      n,
		timestamp: time.Now(),
		expire:    expire,
	}
	self.maps[xid] = r
}

func (self *XidResource) Update(xid string, n *trans.Node) {
	if n == nil {
		return
	}
	r := self.maps[xid]
	// 如果有则更新timestamp
	if r != nil && r.node == n {
		r.Update()
	}
}

func (self *XidResource) Get(xid string) (r *Resource) {
	r, _ = self.maps[xid]
	return
}

func (self *XidResource) Del(xid string) {
	delete(self.maps, xid)
}

func (self *XidResource) DelExpired() {
	for k, v := range self.maps {
		if time.Since(v.timestamp) > v.expire {
			delete(self.maps, k)
		}
	}
}

func (self *XidResource) Reset() {
	self.maps = make(map[string]*Resource)
}

func (self *XidResource) Len() int {
	return len(self.maps)
}