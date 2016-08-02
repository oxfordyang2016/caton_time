//package convert
package main

import (
	"encoding/json"
	"fmt"
	//"strconv"
)

//test finish1
type Base struct {
	From     string `json:"from"`
	Cmd      string `json:"cmd"`
	Type     string `json:"type"`
	Seq      uint32 `json:"seq"`
	Token    string `json:"token,omitempty"`
	Yangming string `json':"yangming"`
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

func main() {

	regireq := `{"from":"a", "cmd": "register", "type": "req",
		"seq": 4294967295,"task_id":"67","uid":"6732","token": "EF02JLGFA09GVNG21F","req":{"register":{"machine_code":"i am machine code"}}}`

	regirsp := `{"from":"a", "cmd": "register", "type": "rsp",
		"seq": 4294967295,"task_id":"67","uid":"6732","token": "EF02JLGFA09GVNG21F","rsp":{"register":{"tnid":"56shusj"}}}`

	loginrsp := `{"from":"a", "cmd": "Login", "type": "rsp",
		"seq": 4294967295,"task_id":"67","uid":"6732","token": "EF02JLGFA09GVNG21F","rsp":{"code":0,"reason":"ok","login":{
	   "alive_interval":546673,"free_mem":563,"upload_bandwidth":83888,
	   "transfer_notify_interval":5366,"time":"653","download_bandwidth":728882,"version":"uis"
		,"total_storage":63,"free_storage":45}}}`

	loginreq := `{"from":"a", "cmd": "Login", "type": "req",
			"seq": 4294967295,"req":{"login":{"version":"890","os":"windows","net_speed":432
			,"net_addr":"172.168.0.12","os":"windows","cpu_usage":54,"total_mem":78,"free_mem":78,
              "download_bandwidth":728882,"upload_bandwidth":83888,
			"token": "EF02JLGFA09GVNG21F",
		   "alive_interval":546673,"free_mem":563,"upload_bandwidth":83888,"storage":[{"name":"s1","type":"local","net":"china_unicom"}],
		   "total_storage":63,"free_storage":45
		   }}}`

	uploadreq := `{"from":"a", "cmd": "upload", "type": "req",
			"seq": 4294967295,"req":{"uploadtask":{
         "taski_d":"i am task_id",
		    "uid":"i am uid",
		    "fid":"i am fid",
		    "sid_list":["1","2","3","4"],
		    "sid_storage":["hshh","0.1.1.2"],
		    "max_bitrate":673

		   }}}`
	uploadrsp := `{"from":"a", "cmd": "upload", "type": "rsp",
			"seq": 4294967295,"rsp":{"uploadtask":{
		    "sid_list":["1","2","3","4","i am sid list"],
		    "sid_storage":["hshh","0.1.1.2"],
		    "recomend_bitrate":673,
		    "port":564,
		    "max_bitrate":673

		   }}}`
	downloadreq := `{"from":"a", "cmd": "download", "type": "req",
			"seq": 4294967295,"req":{"downloadtask":{
         "taski_d":"i am task_id",
		    "uid":"i am uid",
		    "fid":"i am fid",
		    "sid_list":["1","2","3","4"],
		    "sid_storage":["hshh","0.1.1.2"],
		    "max_bitrate":673

		   }}}`
	downloadrsp := `{"from":"a", "cmd": "download", "type": "rsp",
			"seq": 4294967295,"rsp":{"downloadtask":{
		    "sid_list":["1","2","3","4","i am sid list"],
		    "recomend_bitrate":673,
		    "port":564

		   }}}`

	alivereq := `{"from":"a", "cmd": "keepalive", "type": "req",
			"seq": 4294967295,"req":{"keepalive":{"cpu_usage":54,"total_mem":78,"free_mem":78,
              "download_bandwidth":728882,"upload_bandwidth":83888,
		   "total_storage":63,"free_storage":45
		   }}}`

	/*aliversp := `{"from":"a", "cmd": "keepalive", "type": "rsp",
		"seq": 4294967295,"rsp":{"keepalive":{
	    "token":"i am alive token"
	    "time":"564"

	   }}}`
	*/
	aliversp := `{"from":"a", "cmd": "keepalive", "type": "rsp",
			"seq": 4294967295,"rsp":{"keepalive":{
		    "time":"564",
		    "token":"i am alive token"

		   }}}`

	var tag1 Message
	fmt.Println("*is printing return value***", tag1.Convert(loginreq))
	fmt.Println("*is printing return value***", tag1.Convert(loginrsp))
	fmt.Println("*is printing return value***", tag1.Convert(regireq))
	fmt.Println("*is printing return value***", tag1.Convert(regirsp))
	fmt.Println("*is printing return value***", tag1.Convert(downloadreq))
	fmt.Println("*is printing return value***", tag1.Convert(downloadrsp))
	fmt.Println("*is printing return value***", tag1.Convert(uploadreq))
	fmt.Println("*is printing return value***", tag1.Convert(uploadrsp))
	fmt.Println("*is printing return value***", tag1.Convert(alivereq))
	fmt.Println("*is printing return value***", tag1.Convert(aliversp))
}

func (r *Message) Convert(text string) []string {
	fmt.Println()
	var tag1 Message
	var req Request

	var rsp Response
	var down DownloadTaskReq
	var down1 DownloadTaskRsp
	var upload UploadTaskReq
	var upload1 UploadTaskRsp
	var login LoginReq
	var login1 LoginRsp
	var keepalive KeepaliveReq
	var keepalive1 KeepaliveRsp
	var register RegisterReq
	var register1 RegisterRsp
	req.Keepalive = &keepalive
	rsp.Keepalive = &keepalive1
	req.Login = &login
	rsp.Login = &login1
	req.Register = &register
	rsp.Register = &register1
	req.UploadTask = &upload
	rsp.UploadTask = &upload1
	req.DownloadTask = &down
	rsp.DownloadTask = &down1
	tag1.Req = &req
	tag1.Rsp = &rsp

	bytes := []byte(text) //========json input position======>
	// Unmarshal JSON to Result struct.
	json.Unmarshal(bytes, &tag1)
	fmt.Println("cmd is=====>", tag1.Cmd, "================", "type is  ========================>", tag1.Type)

	if tag1.Type == "req" {
		k := tag1.Cmd
		if k == "Login" {
			return []string{tag1.From, tag1.Cmd, tag1.Type, fmt.Sprint(tag1.Seq), tag1.Token, //basic info
				fmt.Sprint(tag1.Req.Login.NetAddr), tag1.Req.Login.OS, fmt.Sprint(tag1.Req.Login.NetSpeed), fmt.Sprint(tag1.Req.Login.Storage),
				fmt.Sprint(tag1.Req.Login.TotalStorage), tag1.Req.Login.Version, fmt.Sprint(tag1.Req.Login.FreeStorage),
				fmt.Sprint(tag1.Req.Login.CpuUsage), fmt.Sprint(tag1.Req.Login.TotalMem), fmt.Sprint(tag1.Req.Login.FreeMem),
				fmt.Sprint(tag1.Req.Login.UploadBandwidth),
				fmt.Sprint(tag1.Req.Login.DownloadBandwidth)}
		}
		if k == "register" {
			return []string{tag1.From, tag1.Cmd, tag1.Type, fmt.Sprint(tag1.Seq), tag1.Token, //basic info
				tag1.Req.Register.MachineCode}
			//there test have problem
		}
		if k == "keepalive" {
			return []string{tag1.From, tag1.Cmd, tag1.Type, fmt.Sprint(tag1.Seq), tag1.Token, //basic info

				fmt.Sprint(tag1.Req.Keepalive.CpuUsage), fmt.Sprint(tag1.Req.Keepalive.DownloadBandwidth), fmt.Sprint(tag1.Req.Keepalive.FreeMem), fmt.Sprint(tag1.Req.Keepalive.FreeStorage),
				fmt.Sprint(tag1.Req.Keepalive.UploadBandwidth)}
		}
		if k == "upload" {
			return []string{tag1.From, tag1.Cmd, tag1.Type, fmt.Sprint(tag1.Seq), tag1.Token, //basic info
				fmt.Sprint(tag1.Req.UploadTask.Fid), fmt.Sprint(tag1.Req.UploadTask.SidList),
				fmt.Sprint(tag1.Req.UploadTask.TaskId), fmt.Sprint(tag1.Req.UploadTask.Uid)}
		}
		if k == "download" {
			return []string{tag1.From, tag1.Cmd, tag1.Type, fmt.Sprint(tag1.Seq), tag1.Token, //basic info
				tag1.Req.DownloadTask.Fid, fmt.Sprint(tag1.Req.DownloadTask.MaxBitrate),
				fmt.Sprint(tag1.Req.DownloadTask.SidList), fmt.Sprint(tag1.Req.DownloadTask.SidStorage),
				tag1.Req.DownloadTask.TaskId, tag1.Req.DownloadTask.Uid}
		}
	}
	if tag1.Type == "rsp" {
		k := tag1.Cmd
		if k == "Login" {
			return []string{tag1.From, tag1.Cmd, tag1.Type, fmt.Sprint(tag1.Seq), tag1.Token,
				fmt.Sprint(tag1.Rsp.Code), tag1.Rsp.Reason, tag1.Rsp.Login.Token, tag1.Rsp.Login.ZoneId,
				fmt.Sprint(tag1.Rsp.Login.AliveInterval), tag1.Rsp.Login.Version, tag1.Rsp.Login.Time}
		}
		if k == "register" {
			return []string{tag1.From, tag1.Cmd, tag1.Type, fmt.Sprint(tag1.Seq), tag1.Token,
				fmt.Sprint(tag1.Rsp.Code), tag1.Rsp.Reason, tag1.Rsp.Register.Tnid}
		}
		if k == "keepalive" {
			return []string{tag1.From, tag1.Cmd, tag1.Type, fmt.Sprint(tag1.Seq), tag1.Token,
				fmt.Sprint(tag1.Rsp.Code), tag1.Rsp.Reason, tag1.Rsp.Keepalive.Token, tag1.Rsp.Keepalive.Time}
		}
		if k == "upload" {
			return []string{tag1.From, tag1.Cmd, tag1.Type, fmt.Sprint(tag1.Seq), tag1.Token,
				fmt.Sprint(tag1.Rsp.Code), tag1.Rsp.Reason, fmt.Sprint(tag1.Rsp.UploadTask.SidStorage),
				fmt.Sprint(tag1.Rsp.UploadTask.Port), fmt.Sprint(tag1.Rsp.UploadTask.SidList),
				fmt.Sprint(tag1.Rsp.UploadTask.RecomendBitrate)}
		}
		if k == "download" {
			return []string{tag1.From, tag1.Cmd, tag1.Type, fmt.Sprint(tag1.Seq), tag1.Token,
				fmt.Sprint(tag1.Rsp.Code), tag1.Rsp.Reason, fmt.Sprint(tag1.Rsp.DownloadTask.SidList), fmt.Sprint(tag1.Rsp.DownloadTask.Port),
				fmt.Sprint(tag1.Rsp.DownloadTask.RecomendBitrate)}
		}
	}
	return []string{"test finish,but have exception occur"}
}
