package pkg

// job指每个用户下载的任务作为一个job或者自己的上传任务作为job.
// 用于以Pkg为单位的上传或者下载状态的跟踪和记录

import (
	"./../transfer/task"
	"./models"
	"cydex"
	"cydex/transfer"
	"fmt"
	clog "github.com/cihub/seelog"
	"strings"
	"sync"
	"time"
)

const (
	// JD数据同步进数据库的间隔
	DEFAULT_CACHE_SYNC_TIMEOUT = 30 * time.Second
)

var (
	JobMgr *JobManager
)

func init() {
	JobMgr = NewJobManager()
}

func HashJob(uid, pid string, typ int) string {
	var s string
	switch typ {
	case cydex.UPLOAD:
		s = "U"
	case cydex.DOWNLOAD:
		s = "D"
	}
	return fmt.Sprintf("%s:%s:%s", uid, s, pid)
}

func getMetaFromTask(t *task.Task) (uid, pid, fid string, typ int) {
	if t == nil {
		return
	}
	uid = t.Uid
	pid = t.Pid
	fid = t.Fid
	typ = t.Type
	return
}

func getHashIdFromTask(t *task.Task) string {
	uid, pid, _, typ := getMetaFromTask(t)
	return HashJob(uid, pid, typ)
}

// func getFidFromTask(t *task.Task) (fid string) {
// 	_, _, fid, _ = getMetaFromTask(t)
// 	return
// }

// 获取运行时segs的信息, 没有就添加进runtime
// size表示接收到的size, 不是seg的size
func getSegRuntime(jd *models.JobDetail, sid string) (s *models.Seg) {
	if jd.Segs == nil {
		clog.Trace("new segs map")
		jd.Segs = make(map[string]*models.Seg)
	}
	s, _ = jd.Segs[sid]
	if s == nil {
		s = new(models.Seg)
		s.Sid = sid
		jd.Segs[sid] = s
	}
	clog.Trace("seg num: ", len(jd.Segs))
	return
}

func updateJobDetail(jd *models.JobDetail, state *transfer.TaskState, seg_state int) (save bool) {
	jd.Bitrate = state.Bitrate
	if seg_state == cydex.TRANSFER_STATE_DONE {
		jd.FinishedSize += state.TotalBytes
		jd.NumFinishedSegs++
		save = true
	} else {
		jd.State = seg_state
	}

	clog.Tracef("%s update: %d %d", jd, jd.NumFinishedSegs, jd.FinishedSize)

	// jd is finished?
	if jd.NumFinishedSegs == jd.File.NumSegs {
		clog.Infof("%s is finished", jd)
		jd.State = cydex.TRANSFER_STATE_DONE
		jd.FinishTime = time.Now()
		save = true
	}

	return
}

// func updateJobDetail(jd *models.JobDetail) {
// 	// finished size
// 	total := uint64(0)
// 	finished_segs := 0
// 	for sid, s := range jd.Segs {
// 		clog.Tracef("%s %s %d", sid, s.Sid, s.Size)
// 		total += s.Size
// 		if s.State == cydex.TRANSFER_STATE_DONE {
// 			finished_segs++
// 		}
// 	}
// 	clog.Tracef("finished size: %d, finished segs:%d", total, finished_segs)
// 	jd.FinishedSize = total
// 	jd.NumFinishedSegs = finished_segs
//
// 	for _, s := range jd.Segs {
// 		if s.State == cydex.TRANSFER_STATE_DOING {
// 			jd.State = cydex.TRANSFER_STATE_DOING
// 			break
// 		}
// 		if s.State == cydex.TRANSFER_STATE_PAUSE {
// 			jd.State = cydex.TRANSFER_STATE_PAUSE
// 			break
// 		}
// 	}
// }

type Track struct {
	// int没啥用,这里当作set使用
	Uploads   map[string]int
	Downloads map[string]int
}

func NewTrack() *Track {
	t := new(Track)
	t.Uploads = make(map[string]int)
	t.Downloads = make(map[string]int)
	return t
}

// type JobRuntime struct {
// 	*models.Job
// 	NumFinishedDetails int
// }

type JobManager struct {
	lock               sync.Mutex
	cache_sync_timeout time.Duration //cache同步超时时间

	jobs          map[string]*models.Job // jobid->job, cache
	track_users   map[string]*Track      // uid->track, track里记录上传下载的pid
	track_pkgs    map[string]*Track      // pid->track, track里记录上传下载的uid
	track_deletes map[string]*Track      // uid->track, track里记录被删除的pid
}

func NewJobManager() *JobManager {
	jm := new(JobManager)
	jm.cache_sync_timeout = DEFAULT_CACHE_SYNC_TIMEOUT
	jm.jobs = make(map[string]*models.Job)
	jm.track_users = make(map[string]*Track)
	jm.track_pkgs = make(map[string]*Track)
	jm.track_deletes = make(map[string]*Track)
	return jm
}

// 创建一个新任务, 因为是活动任务,会加入cache
func (self *JobManager) CreateJob(uid, pid string, typ int) (err error) {
	clog.Infof("create job: u[%s], p[%s], t[%d]", uid, pid, typ)
	hashid := HashJob(uid, pid, typ)
	jobid := hashid

	// 已经存在的状态要复位,并重新加入track
	if job_m, _ := models.GetJob(jobid, false); job_m != nil {
		clog.Warnf("%s is existed, add to track again", jobid)
		self.lock.Lock()
		defer self.lock.Unlock()

		// 已经结束的job需要reset, 重新开始
		if job_m.State == cydex.TRANSFER_STATE_DONE {
			job_m.Reset()
			job_m.GetDetails()
			for _, jd := range job_m.Details {
				jd.GetFile()
				if jd.File.Size > 0 {
					jd.Reset()
				}
			}
		}

		// add track
		self.AddTrack(uid, pid, typ, false)
		// issue-1, 上传用户要监控下载用户状态,上传完的要加入track
		if typ == cydex.DOWNLOAD {
			upload_jobs, _ := models.GetJobsByPid(pid, cydex.UPLOAD, nil)
			for _, u_job := range upload_jobs {
				self.AddTrack(u_job.Uid, u_job.Pid, u_job.Type, false)
			}
		}
		return nil
	}

	session := models.DB().NewSession()
	defer func() {
		models.SessionRelease(session)
		if err != nil {
			clog.Errorf("create job failed: %s", err)
		}
	}()
	if err = session.Begin(); err != nil {
		return
	}

	j := &models.Job{
		JobId: jobid,
		Uid:   uid,
		Pid:   pid,
		Type:  typ,
	}
	if _, err = session.Insert(j); err != nil {
		return err
	}
	clog.Debugf("insert a new Job: %s", jobid)
	// j.Details = make(map[string]*models.JobDetail)
	// create details
	pkg, err := models.GetPkg(pid, true)
	if err != nil || pkg == nil {
		return err
	}
	j.Pkg = pkg
	for _, f := range pkg.Files {
		jd := &models.JobDetail{
			JobId: jobid,
			Fid:   f.Fid,
		}
		// jzh: 如果是0的文件或者文件夹,则状态就置为DONE, 因为客户端不会发送传输命令
		if f.Size == 0 {
			jd.State = cydex.TRANSFER_STATE_DONE
		}
		if _, err := session.Insert(jd); err != nil {
			return err
		}
		// jd.Segs = make(map[string]*models.Seg)
		// j.Details[f.Fid] = jd
		clog.Tracef("insert job_detail fid:%s", f.Fid)
	}
	if err = session.Commit(); err != nil {
		return
	}

	self.lock.Lock()
	defer self.lock.Unlock()
	// add track
	self.AddTrack(uid, pid, typ, false)
	// issue-1, 上传用户要监控下载用户状态,上传完的要加入track
	if typ == cydex.DOWNLOAD {
		upload_jobs, _ := models.GetJobsByPid(pid, cydex.UPLOAD, nil)
		for _, u_job := range upload_jobs {
			self.AddTrack(u_job.Uid, u_job.Pid, u_job.Type, false)
		}
	}

	return nil
}

// 删除job
func (self *JobManager) DeleteJob(uid, pid string, typ int) (err error) {
	hashid := HashJob(uid, pid, typ)
	clog.Infof("delete job: %s", hashid)
	job, _ := models.GetJob(hashid, false)
	if job == nil {
		return
	}
	if err = models.DeleteJob(hashid); err != nil {
		return
	}

	self.lock.Lock()
	defer self.lock.Unlock()

	self.DelTrack(uid, pid, typ, false)
	delete(self.jobs, hashid)
	// 要加入delete track, client通过增量接口获取
	self.AddTrackOfDelete(uid, pid, typ, false)

	clog.Debugf("delete job: %s over", hashid)

	return
}

// 从cache里取; 没有的话从数据库取; 如果是非finished,则加入cache
func (self *JobManager) GetJob(hashid string) *models.Job {
	var err error

	self.lock.Lock()
	defer self.lock.Unlock()

	j, _ := self.jobs[hashid]
	if j != nil {
		return j
	}

	if j, err = models.GetJob(hashid, true); err != nil {
		return nil
	}
	if j != nil {
		j.Details = make(map[string]*models.JobDetail)
		if j.State != cydex.TRANSFER_STATE_DONE {
			j.GetDetails()
			// issue-6: 需要计数已经完成的jd
			for _, jd := range j.Details {
				if jd.State == cydex.TRANSFER_STATE_DONE {
					j.NumFinishedDetails++
				}
			}
			// save to cache
			self.jobs[hashid] = j
		}
	}
	return j
}

// 从cache里取,没有则从数据库取
func (self *JobManager) GetJobDetail(jobid, fid string) (jd *models.JobDetail) {
	var err error
	j := self.GetJob(jobid)
	if j == nil {
		return
	}
	var ok bool
	jd, ok = j.Details[fid]
	if !ok {
		if jd, err = models.GetJobDetail(jobid, fid); err != nil {
			return nil
		}
		j.Details[fid] = jd
	}
	return
}

// implement task.TaskObserver
func (self *JobManager) AddTask(t *task.Task) {
	hashid := getHashIdFromTask(t)
	job := self.GetJob(hashid)
	if job == nil {
		return
	}
	jd := self.GetJobDetail(hashid, t.Fid)
	if jd == nil {
		return
	}
	if jd.StartTime.IsZero() {
		jd.SetStartTime(time.Now())
	}
	// // 上传需要更新seg storage
	// if t.Type == cydex.UPLOAD {
	// }

	if job.State != cydex.TRANSFER_STATE_DOING {
		if job.State == cydex.TRANSFER_STATE_DONE {
			clog.Warnf("%s transfer again", job)
		} else {
			job.SetState(cydex.TRANSFER_STATE_DOING)
		}
	}

	// jzh:不清楚是续传还是补传还是重新下载, 不好处理, api协议有缺陷
	// // TODO: 要处理已经finished的,然后重新下载的
	// if t.DownlaodReq != nil && t.DownloadReq.FinishedSidList != nil {
	// 	//TODO 断点续传, 需要重新计算FinishedSize和NumFinishedSeg
	// }
}

func (self *JobManager) DelTask(t *task.Task) {

}

func (self *JobManager) TaskStateNotify(t *task.Task, state *transfer.TaskState) {
	uid, pid, _, typ := getMetaFromTask(t)
	sid := state.Sid
	if sid == "" || uid == "" || pid == "" {
		return
	}
	hashid := getHashIdFromTask(t)
	j := self.GetJob(hashid)
	if j == nil {
		return
	}
	jd := self.GetJobDetail(hashid, t.Fid)
	if jd == nil {
		return
	}
	if jd.File == nil {
		jd.GetFile()
	}
	// seg_rt := getSegRuntime(jd, sid)

	// 更新JobDetails状态, 根据判断更新Job状态, 是否finished?
	// seg_rt.Size = state.TotalBytes
	seg_state := 0
	// jd.Bitrate = state.Bitrate

	s := strings.ToLower(state.State)
	switch s {
	case "transferring":
		seg_state = cydex.TRANSFER_STATE_DOING
	case "interrupted":
		seg_state = cydex.TRANSFER_STATE_PAUSE
	case "end":
		seg_state = cydex.TRANSFER_STATE_DONE
	default:
		return
	}

	force_save := updateJobDetail(jd, state, seg_state)

	if jd.State == cydex.TRANSFER_STATE_DONE {
		j.NumFinishedDetails++
	}

	// job is finished?
	if j.NumFinishedDetails == len(j.Details) {
		clog.Infof("%s is finished", j)
		j.SetState(cydex.TRANSFER_STATE_DONE)
		self.lock.Lock()
		delete(self.jobs, j.JobId)
		self.DelTrack(uid, pid, typ, false)
		self.lock.Unlock()
	}

	//jzh: 将上传seg的发生变化的状态更新进数据库
	if typ == cydex.UPLOAD {
		seg_m, _ := models.GetSeg(sid)
		// clog.Tracef("sid:%s model_s:%d runtime_s:%d", sid, seg_m.State, seg_rt.State)
		if seg_m != nil && seg_m.State != seg_state {
			// clog.Trace(sid, "set state ", seg_rt.State)
			seg_m.SetState(seg_state)
		}
	}

	if force_save || time.Since(jd.UpdateAt) >= self.cache_sync_timeout {
		jd.Save()
	}
}

func (self *JobManager) SetCacheSyncTimeout(d time.Duration) {
	self.cache_sync_timeout = d
}

func (self *JobManager) HasCachedJob(jobid string) bool {
	defer self.lock.Unlock()
	self.lock.Lock()
	_, ok := self.jobs[jobid]
	return ok
}

func (self *JobManager) AddTrack(uid, pid string, typ int, mutex bool) {
	if mutex {
		self.lock.Lock()
		defer self.lock.Unlock()
	}

	clog.Debugf("add track, u[%s], p[%s], t[%d]", uid, pid, typ)
	track, _ := self.track_pkgs[pid]
	if track == nil {
		track = NewTrack()
		self.track_pkgs[pid] = track
	}
	if typ == cydex.UPLOAD {
		track.Uploads[uid] = 1
	} else {
		track.Downloads[uid] = 1
	}

	track, _ = self.track_users[uid]
	if track == nil {
		track = NewTrack()
		self.track_users[uid] = track
	}
	if typ == cydex.UPLOAD {
		track.Uploads[pid] = 1
	} else {
		track.Downloads[pid] = 1
	}
}

// 删除track, issues-1, 上传用户要等下载完成后才能删除
func (self *JobManager) DelTrack(uid, pid string, typ int, mutex bool) {
	if mutex {
		self.lock.Lock()
		defer self.lock.Unlock()
	}

	if typ == cydex.UPLOAD {
		track, _ := self.track_pkgs[pid]
		if track != nil {
			// 如果还有下载则不退出
			if len(track.Downloads) > 0 {
				return
			}
		}
	}

	self.delTrack(uid, pid, typ)

	// 如果上传的pid, 无下载用户了,需要删除
	track, _ := self.track_pkgs[pid]
	if track != nil {
		if len(track.Downloads) == 0 {
			for uid, _ := range track.Uploads {
				self.delTrack(uid, pid, cydex.UPLOAD)
			}
		}
	}
}

func (self *JobManager) delTrack(uid, pid string, typ int) {
	clog.Debugf("del track, u[%s], p[%s], t[%d]", uid, pid, typ)
	track, _ := self.track_pkgs[pid]
	if track != nil {
		if typ == cydex.UPLOAD {
			delete(track.Uploads, uid)
		} else {
			delete(track.Downloads, uid)
		}
		if len(track.Uploads) == 0 && len(track.Downloads) == 0 {
			delete(self.track_pkgs, pid)
		}
	}

	track, _ = self.track_users[uid]
	if track != nil {
		if typ == cydex.UPLOAD {
			delete(track.Uploads, pid)
		} else {
			delete(track.Downloads, pid)
		}
		if len(track.Uploads) == 0 && len(track.Downloads) == 0 {
			delete(self.track_users, uid)
		}
	}
}

func (self *JobManager) GetPkgTrack(pid string, typ int) (uids []string) {
	self.lock.Lock()
	defer self.lock.Unlock()

	var m map[string]int
	track, _ := self.track_pkgs[pid]
	if track != nil {
		if typ == cydex.UPLOAD {
			m = track.Uploads
		} else {
			m = track.Downloads
		}
		for k, _ := range m {
			uids = append(uids, k)
		}
	}
	return
}

func (self *JobManager) GetUserTrack(uid string, typ int) (pids []string) {
	self.lock.Lock()
	defer self.lock.Unlock()

	var m map[string]int
	track, _ := self.track_users[uid]
	if track != nil {
		if typ == cydex.UPLOAD {
			m = track.Uploads
		} else {
			m = track.Downloads
		}
		for k, _ := range m {
			pids = append(pids, k)
		}
	}
	return
}

// 从cache中获取jobs信息
func (self *JobManager) GetJobsByUid(uid string, typ int) (jobs []*models.Job, err error) {
	pids := self.GetUserTrack(uid, typ)
	for _, pid := range pids {
		hashid := HashJob(uid, pid, typ)
		job := self.GetJob(hashid)
		if job != nil {
			jobs = append(jobs, job)
		}
	}
	return
}

// 从cache中获取jobs信息
func (self *JobManager) GetJobsByPid(pid string, typ int) (jobs []*models.Job, err error) {
	uids := self.GetPkgTrack(pid, typ)
	for _, uid := range uids {
		hashid := HashJob(uid, pid, typ)
		job := self.GetJob(hashid)
		if job != nil {
			jobs = append(jobs, job)
		}
	}
	return
}

// 从数据库中同步track信息
func (self *JobManager) LoadTracks() error {
	clog.Debug("load tracks")
	jobs, err := models.GetUnFinishedJobs()
	if err != nil {
		return err
	}

	defer self.lock.Unlock()
	self.lock.Lock()
	self.ClearTracks(false)
	for _, j := range jobs {
		self.AddTrack(j.Uid, j.Pid, j.Type, false)
		// issue-1, 上传用户要监控下载用户状态,上传完的要加入track
		if j.Type == cydex.DOWNLOAD {
			upload_jobs, _ := models.GetJobsByPid(j.Pid, cydex.UPLOAD, nil)
			for _, u_job := range upload_jobs {
				self.AddTrack(u_job.Uid, u_job.Pid, u_job.Type, false)
			}
		}
	}
	return nil
}

// 增加delete_track的信息
func (self *JobManager) AddTrackOfDelete(uid string, pid string, typ int, mutex bool) {
	if mutex {
		self.lock.Lock()
		defer self.lock.Unlock()
	}

	clog.Debugf("add delete track, u[%s], p[%s], t[%d]", uid, pid, typ)
	track, _ := self.track_deletes[uid]
	if track == nil {
		track = NewTrack()
		self.track_deletes[uid] = track
	}
	if typ == cydex.UPLOAD {
		track.Uploads[pid] = 1
	} else {
		track.Downloads[pid] = 1
	}
}

// 获取delete_track的信息
// remove表示获取后是否删除,下次就取不到了
func (self *JobManager) getTrackOfDelete(uid string, typ int, remove bool) (pids []string) {
	var m map[string]int
	track, _ := self.track_deletes[uid]
	if track == nil {
		return
	}

	if typ == cydex.UPLOAD {
		m = track.Uploads
	} else {
		m = track.Downloads
	}
	for k, _ := range m {
		pids = append(pids, k)

		if remove {
			delete(m, k)
		}
	}

	if len(track.Uploads) == 0 && len(track.Downloads) == 0 {
		delete(self.track_deletes, uid)
	}

	return
}

// 获取delete_track的信息
// remove表示获取后是否删除,下次就取不到了
func (self *JobManager) GetTrackOfDelete(uid string, typ int, remove, mutex bool) (pids []string) {
	if mutex {
		self.lock.Lock()
		defer self.lock.Unlock()
	}
	pids = self.getTrackOfDelete(uid, typ, remove)
	return
}

func (self *JobManager) ClearTracks(mutex bool) {
	if mutex {
		defer self.lock.Unlock()
		self.lock.Lock()
	}
	self.track_users = make(map[string]*Track)
	self.track_pkgs = make(map[string]*Track)
	self.track_deletes = make(map[string]*Track)
}

func PkgIsTransferring(pid string, typ int) (bool, error) {
	jobs, err := JobMgr.GetJobsByPid(pid, typ)
	if err != nil {
		return false, err
	}
	if jobs == nil {
		jobs, err = models.GetJobsByPid(pid, typ, nil)
		if err != nil {
			return false, err
		}
	}
	for _, job := range jobs {
		if job.State == cydex.TRANSFER_STATE_DOING {
			return true, nil
		}
	}
	return false, nil
}

// 更新JD进度
func UpdateJobDetailProcess(job_id, fid string, finished_size uint64, num_finished_segs int) {
	jd := JobMgr.GetJobDetail(job_id, fid)
	if jd == nil {
		return
	}
	clog.Tracef("%s update process: %d %d", jd, finished_size, num_finished_segs)
	jd.FinishedSize = finished_size
	jd.NumFinishedSegs = num_finished_segs
	jd.Save()
}