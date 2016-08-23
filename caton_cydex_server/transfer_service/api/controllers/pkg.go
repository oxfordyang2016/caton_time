package controllers

import (
	"./../../pkg"
	pkg_model "./../../pkg/models"
	"cydex"
	"cydex/transfer"
	"encoding/json"
	// "fmt"
	"./../../transfer/task"
	"errors"
	clog "github.com/cihub/seelog"
	"io/ioutil"
	"strconv"
	"time"
)

func fillTransferState(state *cydex.TransferState, uid string, size uint64, jd *pkg_model.JobDetail) {
	state.Uid = uid
	state.Username = "placehold.cydex"
	state.State = jd.State
	if size > 0 {
		state.Percent = int(jd.FinishedSize * 100 / size)
	} else {
		state.Percent = 100
	}
	if !jd.StartTime.IsZero() {
		t, _ := jd.StartTime.MarshalText()
		s := string(t)
		state.StartTime = &s
	}
	if !jd.FinishTime.IsZero() {
		t, _ := jd.FinishTime.MarshalText()
		s := string(t)
		state.FinishTime = &s
	}
}

// 根据model里的pkg,得到消息响应
func aggregate(pkg_m *pkg_model.Pkg) (pkg_c *cydex.Pkg, err error) {
	if pkg_m == nil {
		return nil, errors.New("nil pkg model")
	}
	clog.Trace(pkg_m.Pid)
	pkg_c = new(cydex.Pkg)

	pkg_c.Pid = pkg_m.Pid
	pkg_c.Title = pkg_m.Title
	pkg_c.Notes = pkg_m.Notes
	pkg_c.NumFiles = int(pkg_m.NumFiles)
	pkg_c.Date = pkg_m.CreateAt
	pkg_c.EncryptionType = pkg_m.EncryptionType
	pkg_c.MetaData = pkg_m.MetaData

	if err = pkg_m.GetFiles(true); err != nil {
		return nil, err
	}

	// 根据pid获取下载该pkg的信息
	download_jobs, err := pkg_model.GetJobsByPid(pkg_m.Pid, cydex.DOWNLOAD, nil)
	if err != nil {
		return nil, err
	}
	upload_jobs, err := pkg_model.GetJobsByPid(pkg_m.Pid, cydex.UPLOAD, nil)
	if err != nil {
		return nil, err
	}

	for _, file_m := range pkg_m.Files {
		file := new(cydex.File)
		file.Fid = file_m.Fid
		file.Filename = file_m.Name
		file.Path = file_m.Path
		file.Type = file_m.Type
		// file.Chara = file_m.EigenValue
		// fake
		file.Chara = "NULL"
		file.PathAbs = file_m.PathAbs
		file.Mode = file_m.Mode
		file.SetSize(file_m.Size)

		// 获取该文件的关联信息
		// 获取文件下载信息
		for _, d_job := range download_jobs {
			// 使用cache中的jd值,cache里没有则使用数据库的, 通过JobMgr接口
			jd := pkg.JobMgr.GetJobDetail(d_job.JobId, file.Fid)
			if jd == nil {
				return nil, errors.New("get job detail failed")
			}
			state := new(cydex.TransferState)
			fillTransferState(state, d_job.Uid, file_m.Size, jd)
			file.DownloadState = append(file.DownloadState, state)
		}
		// 获取文件上传信息
		for _, u_job := range upload_jobs {
			// 使用cache中的jd值,cache里没有则使用数据库的, 通过JobMgr接口
			jd := pkg.JobMgr.GetJobDetail(u_job.JobId, file.Fid)
			if jd == nil {
				return nil, errors.New("get job detail failed")
			}
			state := new(cydex.TransferState)
			fillTransferState(state, u_job.Uid, file_m.Size, jd)
			file.UploadState = append(file.UploadState, state)
		}

		//jzh: 每次都要从数据库里取
		file_m.GetSegs()
		for _, seg_m := range file_m.Segs {
			seg := new(cydex.Seg)
			seg.Sid = seg_m.Sid
			seg.SetSize(seg_m.Size)
			seg.Status = seg_m.State
			file.Segs = append(file.Segs, seg)
		}

		// fill default if nil
		if file.Segs == nil {
			file.Segs = make([]*cydex.Seg, 0)
		}
		if file.DownloadState == nil {
			file.DownloadState = make([]*cydex.TransferState, 0)
		}
		if file.UploadState == nil {
			file.UploadState = make([]*cydex.TransferState, 0)
		}

		pkg_c.Files = append(pkg_c.Files, file)
	}

	return
}

// 包控制器
type PkgsController struct {
	BaseController
}

func (self *PkgsController) Get() {
	query := self.GetString("query")
	filter := self.GetString("filter")
	list := self.GetString("list")

	// sec 3.2 in api doc
	if list == "list" {
		self.getLitePkgs()
		return
	}

	switch query {
	case "all":
		if filter == "change" {
			// sec 3.1.2 in api doc
			self.getActive()
		}
	case "sender":
		// 3.1
		self.getJobs(cydex.UPLOAD)
	case "receiver":
		// 3.1
		self.getJobs(cydex.DOWNLOAD)
	case "admin":
		// 3.1
		self.getAllJobs()
	}
}

func (self *PkgsController) getAllJobs() {
	// uid := self.GetString(":uid")
	page := new(cydex.Pagination)
	page.PageSize, _ = self.GetInt("page_size")
	page.PageNum, _ = self.GetInt("page_num")
	if !page.Verify() {
		page = nil
	}

	rsp := new(cydex.QueryPkgRsp)
	rsp.Error = cydex.OK

	defer func() {
		if rsp.Pkgs == nil {
			rsp.Pkgs = make([]*cydex.Pkg, 0)
		}
		self.Data["json"] = rsp
		self.ServeJSON()
	}()

	// 非admin用户无权限
	if self.UserLevel != cydex.USER_LEVEL_ADMIN {
		rsp.Error = cydex.ErrNotAllowed
		return
	}

	jobs, err := pkg_model.GetJobs(cydex.UPLOAD, page)
	if err != nil {
		// error
		rsp.Error = cydex.ErrInnerServer
		return
	}

	for _, job := range jobs {
		if err = job.GetPkg(true); err != nil {
			// error
			rsp.Error = cydex.ErrInnerServer
			return
		}

		pkg_c, err := aggregate(job.Pkg)
		if err != nil {
			rsp.Error = cydex.ErrInnerServer
			return
		}

		rsp.Pkgs = append(rsp.Pkgs, pkg_c)
	}
}

func (self *PkgsController) getJobs(typ int) {
	uid := self.GetString(":uid")
	page := new(cydex.Pagination)
	page.PageSize, _ = self.GetInt("page_size")
	page.PageNum, _ = self.GetInt("page_num")
	if !page.Verify() {
		page = nil
	}

	rsp := new(cydex.QueryPkgRsp)
	rsp.Error = cydex.OK

	defer func() {
		if rsp.Pkgs == nil {
			rsp.Pkgs = make([]*cydex.Pkg, 0)
		}
		self.Data["json"] = rsp
		self.ServeJSON()
	}()

	// admin返回空
	if self.UserLevel == cydex.USER_LEVEL_ADMIN {
		return
	}

	// 按照uid和type得到和用户相关的jobs
	jobs, err := pkg_model.GetJobsByUid(uid, typ, page)
	if err != nil {
		// error
		rsp.Error = cydex.ErrInnerServer
		return
	}
	for _, job := range jobs {
		if err = job.GetPkg(true); err != nil {
			// error
			rsp.Error = cydex.ErrInnerServer
			return
		}

		pkg_c, err := aggregate(job.Pkg)
		if err != nil {
			rsp.Error = cydex.ErrInnerServer
			return
		}

		rsp.Pkgs = append(rsp.Pkgs, pkg_c)
	}
}

func (self *PkgsController) getActive() {
	uid := self.GetString(":uid")
	rsp := cydex.NewQueryActivePkgRsp()

	clog.Debugf("query active")

	defer func() {
		self.Data["json"] = rsp
		self.ServeJSON()
	}()

	// jzh: 管理员不实现该接口
	if self.UserLevel == cydex.USER_LEVEL_ADMIN {
		return
	}

	// upload
	{
		jobs, err := pkg.JobMgr.GetJobsByUid(uid, cydex.UPLOAD)
		if err != nil {
			clog.Error("get jobs failed, u[%s] t[%d]", uid, cydex.UPLOAD)
			rsp.Error = cydex.ErrInnerServer
			return
		}

		for _, job_m := range jobs {
			if job_m.Pkg == nil || job_m.Pkg.Files == nil {
				// err
				rsp.Error = cydex.ErrInnerServer
				return
			}
			pkg_u := new(cydex.PkgUpload)
			pkg_u.Pid = job_m.Pid
			pkg_u.MetaData = job_m.Pkg.MetaData
			pkg_u.EncryptionType = job_m.Pkg.EncryptionType

			download_jobs, err := pkg.JobMgr.GetJobsByPid(pkg_u.Pid, cydex.DOWNLOAD)
			if err != nil {
				clog.Error("get jobs failed, p[%s] t[%d]", pkg_u.Pid, cydex.DOWNLOAD)
				rsp.Error = cydex.ErrInnerServer
				return
			}

			for _, file_m := range job_m.Pkg.Files {
				file := new(cydex.FileReceiver)
				file.Fid = file_m.Fid
				file.Filename = file_m.Name
				for _, d_job := range download_jobs {
					jd := pkg.JobMgr.GetJobDetail(d_job.JobId, file.Fid)
					if jd != nil {
						state := new(cydex.TransferState)
						fillTransferState(state, d_job.Uid, file_m.Size, jd)
						file.Receivers = append(file.Receivers, state)
					}
				}

				pkg_u.Files = append(pkg_u.Files, file)
			}
			rsp.Uploads = append(rsp.Uploads, pkg_u)
		}
	}

	// download
	{
		jobs, err := pkg.JobMgr.GetJobsByUid(uid, cydex.DOWNLOAD)
		if err != nil {
			clog.Error("get jobs failed, u[%s] t[%d]", uid, cydex.DOWNLOAD)
			rsp.Error = cydex.ErrInnerServer
			return
		}
		for _, job_m := range jobs {
			if job_m.Pkg == nil || job_m.Pkg.Files == nil {
				// err
				rsp.Error = cydex.ErrInnerServer
				return
			}
			pkg_d := new(cydex.PkgDownload)
			pkg_d.Pid = job_m.Pid
			pkg_d.MetaData = job_m.Pkg.MetaData
			pkg_d.EncryptionType = job_m.Pkg.EncryptionType

			uploads_jobs, err := pkg.JobMgr.GetJobsByPid(pkg_d.Pid, cydex.UPLOAD)
			if err != nil || uploads_jobs == nil {
				clog.Error("get jobs failed, p[%s] t[%d]", pkg_d.Pid, cydex.UPLOAD)
				rsp.Error = cydex.ErrInnerServer
				return
			}

			for _, file_m := range job_m.Pkg.Files {
				file := new(cydex.FileSender)
				file.Fid = file_m.Fid
				file.Filename = file_m.Name
				for _, u_job := range uploads_jobs {
					// clog.Trace(file.Fid)
					// clog.Trace(u_job.JobId)
					jd := pkg.JobMgr.GetJobDetail(u_job.JobId, file.Fid)
					if jd != nil {
						state := new(cydex.TransferState)
						fillTransferState(state, u_job.Uid, file_m.Size, jd)
						file.Sender = append(file.Sender, state)
					}
				}

				//jzh: 每次都要从数据库里取
				file_m.GetSegs()
				for _, seg_m := range file_m.Segs {
					seg := new(cydex.Seg)
					seg.Sid = seg_m.Sid
					seg.SetSize(seg_m.Size)
					seg.Status = seg_m.State
					file.Segs = append(file.Segs, seg)
				}

				pkg_d.Files = append(pkg_d.Files, file)
			}
			rsp.Downloads = append(rsp.Downloads, pkg_d)
		}
	}

	//delete_sender_pkg_list
	pids := pkg.JobMgr.GetTrackOfDelete(uid, cydex.UPLOAD, true, true)
	for _, pid := range pids {
		rsp.DelSenderPkgs = append(rsp.DelSenderPkgs, &cydex.PkgId{pid})
	}

	//delete_receive_pkg_list
	pids = pkg.JobMgr.GetTrackOfDelete(uid, cydex.DOWNLOAD, true, true)
	for _, pid := range pids {
		rsp.DelReceiverPkgs = append(rsp.DelReceiverPkgs, &cydex.PkgId{pid})
	}
}

func (self *PkgsController) getLitePkgs() {
	uid := self.GetString(":uid")
	query := self.GetString("query")
	rsp := new(cydex.QueryPkgLiteRsp)

	defer func() {
		if rsp.Pkgs == nil {
			rsp.Pkgs = make([]*cydex.PkgLite, 0)
		}
		self.Data["json"] = rsp
		self.ServeJSON()
	}()

	clog.Debugf("get lite pkgs query is %s", query)

	var pkgs []*pkg_model.Pkg
	var err error

	switch query {
	case "admin":
		if self.UserLevel != cydex.USER_LEVEL_ADMIN {
			rsp.Error = cydex.ErrNotAllowed
			return
		}
		jobs, err := pkg_model.GetJobs(cydex.UPLOAD, nil)
		if err != nil {
			clog.Error(err)
			rsp.Error = cydex.ErrInnerServer
			return
		}
		for _, j := range jobs {
			if j.Pkg == nil {
				j.GetPkg(false)
			}
			pkgs = append(pkgs, j.Pkg)
		}
	case "sender":
		jobs, err := pkg_model.GetJobsByUid(uid, cydex.UPLOAD, nil)
		if err != nil {
			clog.Error(err)
			rsp.Error = cydex.ErrInnerServer
			return
		}
		for _, j := range jobs {
			if j.Pkg == nil {
				j.GetPkg(false)
			}
			pkgs = append(pkgs, j.Pkg)
		}
	case "receiver":
		jobs, err := pkg_model.GetJobsByUid(uid, cydex.DOWNLOAD, nil)
		if err != nil {
			clog.Error(err)
			rsp.Error = cydex.ErrInnerServer
			return
		}
		for _, j := range jobs {
			if j.Pkg == nil {
				j.GetPkg(false)
			}
			pkgs = append(pkgs, j.Pkg)
		}
	default:
		clog.Warnf("Invalid query:%s", query)
		rsp.Error = cydex.ErrInvalidParam
		return
	}

	if err != nil {
		rsp.Error = cydex.ErrInnerServer
		return
	}

	for _, p := range pkgs {
		rsp.Pkgs = append(rsp.Pkgs, &cydex.PkgLite{
			Pid:      p.Pid,
			Title:    p.Title,
			Date:     p.CreateAt,
			Notes:    p.Notes,
			NumFiles: int(p.NumFiles),
		})
	}
}

func (self *PkgsController) Post() {
	rsp := new(cydex.CreatePkgRsp)
	req := new(cydex.CreatePkgReq)

	defer func() {
		self.Data["json"] = rsp
		self.ServeJSON()
	}()

	if self.UserLevel == cydex.USER_LEVEL_ADMIN {
		clog.Error("admin not allowed to create pkg")
		rsp.Error = cydex.ErrNotAllowed
		return
	}

	// 获取拆包器
	unpacker := pkg.GetUnpacker()
	if unpacker == nil {
		rsp.Error = cydex.ErrInnerServer
		return
	}

	// 获取请求
	r := self.Ctx.Request
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rsp.Error = cydex.ErrInvalidParam
		return
	}
	clog.Debug(string(body))
	if err := json.Unmarshal(body, req); err != nil {
		clog.Error(err)
		rsp.Error = cydex.ErrInvalidParam
		return
	}
	for _, f := range req.Files {
		f.Size, _ = strconv.ParseUint(f.SizeStr, 10, 64)
		clog.Trace(f.Size)
	}

	self.createPkg(req, rsp)
}

// 创建包裹
func (self *PkgsController) createPkg(req *cydex.CreatePkgReq, rsp *cydex.CreatePkgRsp) {
	if !req.Verify() {
		clog.Error("create pkg request is not verified!")
		rsp.Error = cydex.ErrInvalidParam
		return
	}

	var (
		err           error
		uid, fid, pid string
		file_o        *cydex.File
	)
	uid = self.GetString(":uid")
	if uid == "" {
		rsp.Error = cydex.ErrInvalidParam
		return
	}
	session := pkg_model.DB().NewSession()
	session.Begin()

	defer func() {
		if err != nil {
			session.Rollback()
		}
		session.Close()
	}()

	unpacker := pkg.GetUnpacker()
	if unpacker == nil {
		rsp.Error = cydex.ErrInnerServer
		return
	}

	pkg_o := new(cydex.Pkg)
	pid = unpacker.GeneratePid(uid, req.Title, req.Notes)
	size := uint64(0)
	for _, f := range req.Files {
		size += f.Size
	}

	// 创建pkg数据库记录
	pkg_m := &pkg_model.Pkg{
		Pid:            pid,
		Title:          req.Title,
		Notes:          req.Notes,
		NumFiles:       uint64(len(req.Files)),
		Size:           size,
		EncryptionType: req.EncryptionType,
		MetaData:       req.MetaData,
	}
	if _, err = session.Insert(pkg_m); err != nil {
		clog.Error(err)
		rsp.Error = cydex.ErrInnerServer
		return
	}

	pkg_o.Pid = pid
	pkg_o.Title = req.Title
	pkg_o.Notes = req.Notes
	pkg_o.Date = time.Now()
	pkg_o.NumFiles = len(req.Files)
	pkg_o.EncryptionType = pkg_m.EncryptionType
	pkg_o.MetaData = pkg_m.MetaData

	for i, f := range req.Files {
		if fid, err = unpacker.GenerateFid(pid, i, f); err != nil {
			rsp.Error = cydex.ErrInnerServer
			return
		}
		file_o = new(cydex.File)
		file_o.Fid = fid
		file_o.SimpleFile = *f
		file_o.Segs, err = unpacker.GenerateSegs(fid, f)
		if err != nil {
			rsp.Error = cydex.ErrInnerServer
			return
		}
		file_o.DownloadState = make([]*cydex.TransferState, 0)
		file_o.UploadState = make([]*cydex.TransferState, 0)
		file_m := &pkg_model.File{
			Fid:        fid,
			Pid:        pid,
			Name:       f.Filename,
			Path:       f.Path,
			Size:       f.Size,
			Mode:       f.Mode,
			Type:       f.Type,
			PathAbs:    f.PathAbs,
			EigenValue: f.Chara,
			NumSegs:    len(file_o.Segs),
		}
		if _, err = session.Insert(file_m); err != nil {
			rsp.Error = cydex.ErrInnerServer
			return
		}
		for _, s := range file_o.Segs {
			seg_m := &pkg_model.Seg{
				Sid:  s.Sid,
				Fid:  fid,
				Size: s.Size,
			}
			if _, err = session.Insert(seg_m); err != nil {
				rsp.Error = cydex.ErrInnerServer
				return
			}
		}
		pkg_o.Files = append(pkg_o.Files, file_o)
	}
	session.Commit()
	rsp.Pkg = pkg_o

	// Dispatch jobs
	pkg.JobMgr.CreateJob(uid, pid, cydex.UPLOAD)
	if req.Receivers != nil {
		for _, r := range req.Receivers {
			pkg.JobMgr.CreateJob(r.Uid, pid, cydex.DOWNLOAD)
		}
	}
}

// 单包控制器
type PkgController struct {
	BaseController
}

func (self *PkgController) Get() {
	rsp := new(cydex.QuerySinglePkgRsp)

	defer func() {
		self.Data["json"] = rsp
		self.ServeJSON()
	}()

	uid := self.GetString(":uid")
	pid := self.GetString(":pid")

	if self.UserLevel != cydex.USER_LEVEL_ADMIN {
		types := []int{cydex.UPLOAD, cydex.DOWNLOAD}
		for _, t := range types {
			hashid := pkg.HashJob(uid, pid, t)
			job, err := pkg_model.GetJob(hashid, true)
			if err == nil && job != nil {
				rsp.Pkg, _ = aggregate(job.Pkg)
			}
		}
	} else {
		// 管理员用户
		clog.Trace("admin get pkg")
		jobs, _ := pkg_model.GetJobsByPid(pid, cydex.UPLOAD, nil)
		if len(jobs) > 0 {
			job := jobs[0]
			job.GetPkg(true)
			rsp.Pkg, _ = aggregate(job.Pkg)
		}
	}
}

// 3.7.2
func (self *PkgController) Put() {
	uid := self.GetString(":uid")
	pid := self.GetString(":pid")

	req := new(cydex.ModifyPkgReq)
	rsp := new(cydex.ModifyPkgRsp)

	defer func() {
		self.Data["json"] = rsp
		self.ServeJSON()
	}()

	// 获取请求
	r := self.Ctx.Request
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rsp.Error = cydex.ErrInvalidParam
		return
	}
	clog.Debug(string(body))
	if err := json.Unmarshal(body, req); err != nil {
		clog.Error(err)
		rsp.Error = cydex.ErrInvalidParam
		return
	}
	var job_m *pkg_model.Job

	if self.UserLevel != cydex.USER_LEVEL_ADMIN {
		hashid := pkg.HashJob(uid, pid, cydex.UPLOAD)
		job_m, err = pkg_model.GetJob(hashid, false)
	} else {
		// admin
		clog.Trace("admin put pkg")
		jobs, _ := pkg_model.GetJobsByPid(pid, cydex.UPLOAD, nil)
		if len(jobs) > 0 {
			job_m = jobs[0]
			job_m.GetPkg(true)
		}
	}
	if err != nil {
		clog.Error(err)
		rsp.Error = cydex.ErrInvalidParam
		return
	}
	if job_m == nil {
		rsp.Error = cydex.ErrPackageNotExisted
		return
	}
	clog.Tracef("%+v", job_m)
	job_m.GetPkg(true)

	// add
	for _, uid := range req.AddList {
		if err := pkg.JobMgr.CreateJob(uid, pid, cydex.DOWNLOAD); err != nil {
			clog.Error(err)
		}
	}

	// delete
	for _, uid := range req.RemoveList {
		// stop transferring tasks
		clog.Trace("remove:", uid, pid)
		task.TaskMgr.StopTasks(uid, pid, cydex.DOWNLOAD)

		// delete resource in cache and db
		if err := pkg.JobMgr.DeleteJob(uid, pid, cydex.DOWNLOAD); err != nil {
			clog.Error(err)
			continue
		}
	}

	rsp.Pkg, _ = aggregate(job_m.Pkg)
}

// TODO 删除包裹, 释放空间
func (self *PkgController) Delete() {
	uid := self.GetString(":uid")
	pid := self.GetString(":pid")

	rsp := new(cydex.BaseRsp)

	defer func() {
		self.Data["json"] = rsp
		self.ServeJSON()
	}()

	var job_m *pkg_model.Job

	// 判断是否有此包
	if self.UserLevel != cydex.USER_LEVEL_ADMIN {
		hashid := pkg.HashJob(uid, pid, cydex.UPLOAD)
		clog.Trace(hashid)
		job_m, _ = pkg_model.GetJob(hashid, true)
	} else {
		// admin
		clog.Trace("admin delete pkg")
		jobs, _ := pkg_model.GetJobsByPid(pid, cydex.UPLOAD, nil)
		if len(jobs) > 0 {
			job_m = jobs[0]
			job_m.GetPkg(true)
		}
	}
	if job_m == nil || job_m.Pkg == nil {
		rsp.Error = cydex.ErrPackageNotExisted
		return
	}

	// 是否在传输
	transferring, err := pkg.PkgIsTransferring(pid, cydex.UPLOAD)
	if err != nil {
		rsp.Error = cydex.ErrInnerServer
		return
	}
	if transferring {
		rsp.Error = cydex.ErrActivePackage
		return
	}
	transferring, err = pkg.PkgIsTransferring(pid, cydex.DOWNLOAD)
	if err != nil {
		rsp.Error = cydex.ErrInnerServer
		return
	}
	if transferring {
		rsp.Error = cydex.ErrActivePackage
		return
	}

	//上传Job软删除
	job_m.SoftDelete(pkg_model.SOFT_DELETE_TAG)
	// 下载Job真删除
	download_jobs, _ := pkg_model.GetJobsByPid(pid, cydex.DOWNLOAD, nil)
	for _, j := range download_jobs {
		// delete resource in cache and db
		if err := pkg.JobMgr.DeleteJob(j.Uid, pid, cydex.DOWNLOAD); err != nil {
			clog.Error(err)
			continue
		}
	}

	go deletePkg(job_m)
}

func deletePkg(job *pkg_model.Job) {
	// 删除文件, 释放空间
	if job == nil {
		return
	}
	if job.Pkg.Files == nil {
		job.Pkg.GetFiles(true)
	}
	pid := job.Pid
	c := make(chan int)
	const timeout = 5 * time.Second

	for _, file := range job.Pkg.Files {
		// async job
		go func() {
			// 注意循环参数'file'对闭包的影响
			file := file
			defer func() {
				c <- 1
			}()
			if file.Size == 0 {
				return
			}
			file.GetSegs()
			var sids []string
			for _, seg := range file.Segs {
				sids = append(sids, seg.Sid)
			}
			if sids == nil {
				return
			}
			task_req := buildTaskDownloadReq(job.Uid, pid, file.Fid, sids)
			clog.Tracef("%+v", task_req)
			node, err := task.TaskMgr.Scheduler().DispatchDownload(task_req)
			if node != nil {
				// 发送协议
				msg := transfer.NewReqMessage("", "removefile", "", 0)
				msg.Req.RemoveFile = &transfer.RemoveFileReq{
					RemoveId:   "",
					Uid:        job.Uid,
					Pid:        pid,
					Fid:        file.Fid,
					SidList:    sids,
					SidStorage: task_req.DownloadTaskReq.SidStorage,
				}
				clog.Tracef("%+v", msg.Req.RemoveFile)
				if _, err = node.SendRequestSync(msg, timeout); err != nil {
					// clog.Error("")
					return
				}
			}
		}()
	}

	// wait delete job over
	for i := 0; i < len(job.Pkg.Files); i++ {
		<-c
	}

	job.SoftDelete(pkg_model.SOFT_DELETE_FILES_REMOVED)
}

// 单文件控制器
type FileController struct {
	BaseController
}

func (self *FileController) Get() {
	uid := self.GetString(":uid")
	pid := self.GetString(":pid")
	fid := self.GetString(":fid")

	rsp := new(cydex.QuerySingleFileRsp)

	defer func() {
		self.Data["json"] = rsp
		self.ServeJSON()
	}()

	file_m, err := pkg_model.GetFile(fid)
	if err != nil || file_m == nil || file_m.Pid != pid {
		rsp.Error = cydex.ErrPackageNotExisted
		return
	}

	// 下载
	{
		hashid := pkg.HashJob(uid, pid, cydex.DOWNLOAD)
		job := pkg.JobMgr.GetJob(hashid)
		if job != nil {
			file := new(cydex.File)
			file.Fid = file_m.Fid
			file.Filename = file_m.Name
			file.Size = file_m.Size
			file.Path = file_m.Path
			file.Type = file_m.Type
			file.Chara = file_m.EigenValue
			file.PathAbs = file_m.PathAbs

			uploads_jobs, err := pkg_model.GetJobsByPid(pid, cydex.UPLOAD, nil)
			if err != nil {
				rsp.Error = cydex.ErrInnerServer
				return
			}
			for _, u_job := range uploads_jobs {
				jd := pkg.JobMgr.GetJobDetail(u_job.JobId, fid)
				if jd != nil {
					state := new(cydex.TransferState)
					fillTransferState(state, u_job.Uid, file_m.Size, jd)
					file.UploadState = append(file.UploadState, state)
				}
			}
			rsp.File = file
			return
		}
	}

	// 上传
	{
		hashid := pkg.HashJob(uid, pid, cydex.UPLOAD)
		job := pkg.JobMgr.GetJob(hashid)
		if job != nil {
			file := new(cydex.File)
			file.Fid = file_m.Fid
			file.Filename = file_m.Name
			file.Size = file_m.Size
			file.Path = file_m.Path
			file.Type = file_m.Type
			file.Chara = file_m.EigenValue
			file.PathAbs = file_m.PathAbs

			download_jobs, err := pkg_model.GetJobsByPid(pid, cydex.DOWNLOAD, nil)
			if err != nil {
				rsp.Error = cydex.ErrInnerServer
				return
			}
			for _, d_job := range download_jobs {
				jd := pkg.JobMgr.GetJobDetail(d_job.JobId, fid)
				if jd != nil {
					state := new(cydex.TransferState)
					fillTransferState(state, d_job.Uid, file_m.Size, jd)
					file.DownloadState = append(file.DownloadState, state)
				}
			}
			rsp.File = file
			return
		}
	}

	rsp.Error = cydex.ErrPackageNotExisted
}
