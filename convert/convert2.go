package main

import (
	"encoding/json"
	"fmt"
	//"strconv"
)

//test finish1
type Base struct {
	From  string `json:"from"`
	Cmd   string `json:"cmd"`
	Type  string `json:"type"`
	Seq   uint32 `json:"seq"`
	Token string `json:"token,omitempty"`
	//	Yangming string `json':"yangming"`
}

type Message struct {
	Base
	Req *Request  `json:"req,omitempty"`
	Rsp *Response `json:"rsp,omitempty"`
}

//
type TransferNotify struct {
	TaskStateList []TaskState `json:"task_list"`
}

type Request struct {
	Register     *RegisterReq     `json:"register,omitempty"`
	Login        *LoginReq        `json:"login,omitempty"`
	Keepalive    *KeepaliveReq    `json:"keepalive,omitempty"`
	UploadTask   *UploadTaskReq   `json:"uploadtask,omitempty"`
	DownloadTask *DownloadTaskReq `json:"downloadtask,omitempty"`
	//	TransferNotify *TransferNotifyReq `json:"transfernotify,omitempty"`
}
type Response struct {
	Code         int              `json:"code"`
	Reason       string           `json:"reason"`
	Register     *RegisterRsp     `json:"register,omitempty"`
	Login        *LoginRsp        `json:"login,omitempty"`
	Keepalive    *KeepaliveRsp    `json:"keepalive,omitempty"`
	UploadTask   *UploadTaskRsp   `json:"uploadtask,omitempty"`
	DownloadTask *DownloadTaskRsp `json:"downloadtask,omitempty"`
}
type RegisterReq struct {
	MachineCode string `json:"machine_code"`
}
type RegisterRsp struct {
	Tnid string `json:"tnid"`
}
type LoginReq struct {
	Version           string        `json:"version"`
	NetAddr           string        `json:"net_addr"`
	OS                string        `json:"os"`
	NetSpeed          uint32        `json:"net_speed"`
	Storage           []StorageInfo `json:"storage"`
	TotalStorage      uint64        `json:"total_storage"`
	FreeStorage       uint64        `json:"free_storage"`
	CpuUsage          uint32        `json:"cpu_usage"`
	TotalMem          uint64        `json:"total_mem"`
	FreeMem           uint64        `json:"free_mem"`
	UploadBandwidth   uint64        `json:"upload_bandwidth"`
	DownloadBandwidth uint64        `json:"download_bandwidth"`
}

type StorageInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Net  string `json:"net,omitempty"`
}
type LoginRsp struct {
	Token                  string `json:"token"`
	ZoneId                 string `json:"zone_id,omitempty"`
	AliveInterval          uint32 `json:"alive_interval"`
	TransferNotifyInterval uint32 `json:"transfer_notify_interval"`
	Time                   string `json:"time,omitempty"`
	Version                string `json:"version,omitempty"`
}

type KeepaliveReq struct {
	CpuUsage          uint32 `json:"cpu_usage"`
	TotalStorage      uint64 `json:"total_storage"`
	FreeStorage       uint64 `json:"free_storage"`
	TotalMem          uint64 `json:"total_mem"`
	FreeMem           uint64 `json:"free_mem"`
	UploadBandwidth   uint64 `json:"upload_bandwidth"`
	DownloadBandwidth uint64 `json:"download_bandwidth"`
}

type KeepaliveRsp struct {
	Token string `json:"token,omitempty"`
	Time  string `json:"time,omitempty"`
}
type UploadTaskReq struct {
	TaskId  string   `json:"task_id"`
	Uid     string   `json:"uid"`
	Fid     string   `json:"fid"`
	SidList []string `json:"sid_list"`
}

type UploadTaskRsp struct {
	SidList         []string `json:"sid_list"`
	SidStorage      []string `json:"sid_sorage"`
	Port            uint32   `json:"port"`
	RecomendBitrate uint32   `json:"recomend_bitrate"`
}

type DownloadTaskReq struct {
	TaskId     string   `json:"task_id"`
	Uid        string   `json:"uid"`
	Fid        string   `json:"fid"`
	SidList    []string `json:"sid_list"`
	SidStorage []string `json:"sid_storage"`
	MaxBitrate uint32   `json:"max_bitrate,omitempty"`
}

type DownloadTaskRsp struct {
	SidList         []string `json:"sid_list"`
	Port            uint32   `json:"port"`
	RecomendBitrate uint32   `json:"recomend_bitrate"`
}
type TaskState struct {
	TaskId     string `json:"task_id"`
	State      string `json:"state"`
	TotalBytes uint64 `json:"total_bytes"`
	Bitrate    uint64 `json:"bitrate"`
}

type MyUser struct {
	Base
	//ID  int64    `json:"id"`
	Req *Request `json:"req,omitempty"`
	Rsp *Response
	//	Rsp *Response `json:"rsp,omitempty"`
}

var regireq = Message{Base{"kksk", "KKD", "JSJ", 563, "73HJSKA_K"},
	&Request{&RegisterReq{MachineCode: " i am machine code"}, &LoginReq{}, &KeepaliveReq{}, &UploadTaskReq{}, &DownloadTaskReq{}},
	&Response{1, "ok", &RegisterRsp{}, &LoginRsp{}, &KeepaliveRsp{}, &UploadTaskRsp{}, &DownloadTaskRsp{}}}

//var regireq1 = Message{Base{"kksk", "KKD", "JSJ", 563, "73HJSKA_K"},
//	&Request{&RegisterReq{MachineCode: " i am machine code"}}}

var regirsp = Message{Base{"kksk", "KKD", "JSJ", 563, "73HJSKA_K"},
	&Request{&RegisterReq{}, &LoginReq{}, &KeepaliveReq{}, &UploadTaskReq{}, &DownloadTaskReq{}},
	&Response{1, "ok", &RegisterRsp{"i am tnid"}, &LoginRsp{}, &KeepaliveRsp{}, &UploadTaskRsp{}, &DownloadTaskRsp{}}}

var loginrsp = Message{Base{"kksk", "KKD", "JSJ", 563, "73HJSKA_K"}, nil,
	&Response{1, "ok", &RegisterRsp{}, &LoginRsp{}, &KeepaliveRsp{}, &UploadTaskRsp{}, &DownloadTaskRsp{}}}

func main() {
	//_ = json.NewEncoder(os.Stdout).Encode(
	//	&MyUser{Base{"kksk"}, 1, &Request{"req"}, &Response{}})

	//k3, _ := json.Marshal(&k2)
	//fmt.Println(string(k3))
	//fmt.Println(Change(&k2))
	Change(&loginrsp)
	Change(&regireq)
	//	Change(&regireq1)
	Change(&regirsp)
}

func Change(p *Message) {
	//_ = json.NewEncoder(os.Stdout).Encode(
	//	&MyUser{Base{"kksk"}, 1, &Request{"req"}, &Response{}})

	k3, _ := json.Marshal(p)
	fmt.Println(string(k3))
	fmt.Println("==========================================================================================================")
	//fmt.Println(Change(&k2))
	//Change1(&k2)
}
