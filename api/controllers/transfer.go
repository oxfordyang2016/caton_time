package controllers

import (
	"./../../pkg"
	pkg_model "./../../pkg/models"
	// trans "./../../transfer"
	"./../../transfer/task"
	"cydex"
	"cydex/transfer"
	"encoding/json"
	"github.com/astaxie/beego"
	clog "github.com/cihub/seelog"
	"io/ioutil"
	"strconv"
	"time"
)

const (
	DISPATCH_TIMEOUT = 10 * time.Second
)

type TransferController struct {
	beego.Controller
}

func (self *TransferController) Post() {
	req := new(cydex.TransferReq)
	rsp := new(cydex.TransferRsp)

	clog.Trace("transfer")
	defer func() {
		self.Data["json"] = rsp
		self.ServeJSON()
	}()

	r := self.Ctx.Request
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rsp.Error = cydex.ErrInvalidParam
		return
	}
	clog.Trace(string(body))
	// 获取请求
	if err := json.Unmarshal(body, req); err != nil {
		clog.Error(err)
		rsp.Error = cydex.ErrInvalidParam
		return
	}
	req.FinishedSize, _ = strconv.ParseUint(req.FinishedSizeStr, 10, 64)

	if len(req.SegIds) == 0 || req.Fid == "" {
		clog.Error("invalid param")
		rsp.Error = cydex.ErrInvalidParam
		return
	}

	switch req.From {
	case "receiver":
		self.processDownload(req, rsp)
	case "sender":
		self.processUpload(req, rsp)
	default:
		rsp.Error = cydex.ErrInvalidParam
		return
	}
}

func (self *TransferController) processDownload(req *cydex.TransferReq, rsp *cydex.TransferRsp) {
	clog.Trace("download")
	uid := self.GetString(":uid")
	pid := pkg.GetUnpacker().GetPidFromFid(req.Fid)
	task_req := new(task.DownloadReq)
	task_req.DownloadTaskReq = &transfer.DownloadTaskReq{
		TaskId:  task.GenerateTaskId(),
		Uid:     uid,
		Pid:     pid,
		Fid:     req.Fid,
		SidList: req.SegIds,
	}

	// update jd process
	jobid := pkg.HashJob(uid, pid, cydex.DOWNLOAD)
	pkg.UpdateJobDetailProcess(jobid, req.Fid, req.FinishedSize, req.NumFinishedSegs)

	// get storages
	for _, sid := range req.SegIds {
		seg, _ := pkg_model.GetSeg(sid)
		storage := ""
		if seg != nil {
			storage = seg.Storage
		}
		task_req.DownloadTaskReq.SidStorage = append(task_req.DownloadTaskReq.SidStorage, storage)
	}

	task_rsp, node, err := task.TaskMgr.DispatchDownload(task_req, DISPATCH_TIMEOUT)
	if err != nil || node == nil || task_rsp == nil {
		clog.Error(err)
		rsp.Error = cydex.ErrInnerServer
		return
	}

	rsp.Host = getHost(node.Host)
	rsp.Port = task_rsp.Port
	rsp.RecomendBitrate = task_rsp.RecomendBitrate
	for _, sid := range task_rsp.SidList {
		seg_m, _ := pkg_model.GetSeg(sid)
		if seg_m != nil {
			seg := new(cydex.Seg)
			seg.Sid = seg_m.Sid
			seg.SetSize(seg_m.Size)
			seg.Status = seg_m.State
			rsp.Segs = append(rsp.Segs, seg)
		}
		if rsp.Segs == nil {
			rsp.Segs = make([]*cydex.Seg, 0)
		}
	}
}

func (self *TransferController) processUpload(req *cydex.TransferReq, rsp *cydex.TransferRsp) {
	uid := self.GetString(":uid")
	task_req := new(task.UploadReq)
	pid := pkg.GetUnpacker().GetPidFromFid(req.Fid)
	task_req.UploadTaskReq = &transfer.UploadTaskReq{
		TaskId:  task.GenerateTaskId(),
		Uid:     uid,
		Pid:     pid,
		Fid:     req.Fid,
		SidList: req.SegIds,
	}
	for _, sid := range req.SegIds {
		seg, _ := pkg_model.GetSeg(sid)
		if seg != nil {
			task_req.UploadTaskReq.Size += seg.Size
		}
	}

	clog.Tracef("pid:%s, segs size:%d", task_req.Pid, task_req.UploadTaskReq.Size)

	task_rsp, node, err := task.TaskMgr.DispatchUpload(task_req, DISPATCH_TIMEOUT)
	clog.Trace(task_rsp, node)
	if err != nil || node == nil || task_rsp == nil {
		clog.Error(err)
		rsp.Error = cydex.ErrInnerServer
		return
	}
	clog.Trace("dispatch upload ok")

	update_storage := true
	if len(task_rsp.SidStorage) != len(task_rsp.SidList) {
		clog.Error("storage len is not match")
		update_storage = false
	}

	rsp.Host = getHost(node.Host)
	rsp.Port = task_rsp.Port
	rsp.RecomendBitrate = task_rsp.RecomendBitrate
	for i, sid := range task_rsp.SidList {
		seg_m, _ := pkg_model.GetSeg(sid)
		if seg_m != nil {
			if update_storage {
				seg_m.SetStorage(task_rsp.SidStorage[i])
			}
			seg := new(cydex.Seg)
			seg.Sid = seg_m.Sid
			seg.SetSize(seg_m.Size)
			seg.Status = seg_m.State
			rsp.Segs = append(rsp.Segs, seg)
		}
		if rsp.Segs == nil {
			rsp.Segs = make([]*cydex.Seg, 0)
		}
	}
}

// 如果是"127.0.0.1", 则返回空,兼容单机版
func getHost(host string) string {
	if host == "127.0.0.1" {
		return ""
	}
	return host
}
