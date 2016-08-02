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
	// This JSON contains an int array.
	//text1 := `{ "cmd": "Login", "type": "req", "seq": "1", "rsp": { "code": 0, "reason": "OK", "body": { "token": "EF02JLGFA09GVNG21F", "alive_interval": 100 } } }`
	//text := `{ "from":"a", "cmd": "Download", "type": "req", "seq": 4294967295,  "token": "EF02JLGFA09GVNG21F" ,"yangming":890}`
	/*loginreq:=`{"from":"a", "cmd": "Login", "type": "req", "seq": 4294967295,  "token": "EF02JLGFA09GVNG21F" ,"yangming":890,"version":"3.0","net_addr":"172.168.0.12","os":"windows","storage":[{"name":"s1","type":"local","net":"china_unicom"}],"total_storage":63,"free_storage":45
	  ,
	  "cpu_usage":65677,"total_mem":88993,"free_mem":563,"upload_bandwidth":83888,"download_bandwidth":728882

	  }`*/

	/*loginrsp:=`{"from":"a", "cmd": "Login", "type": "rsp", "seq": 4294967295,  "token": "EF02JLGFA09GVNG21F" ,"yangming":890,"version":"uis","code":54,
	  "alive_interval":546673,"free_mem":563,"upload_bandwidth":83888,"transfer_notify_interval":5366,"time":"653","download_bandwidth":728882

	  }`
	*/
	/*regsterreq:=`{"from":"a", "cmd": "Register", "type": "rsp", "seq": 4294967295,  "token": "EF02JLGFA09GVNG21F" ,"yangming":890,"machine_code":"x0-6dghsjjjjs-*idii","tnid":"656"}`
	 */

	/*
	   keepalivereq:=`{"from":"a", "cmd": "keepalive", "type": "req", "seq": 4294967295,  "token": "EF02JLGFA09GVNG21F","free_mem":6783,"cpu_usage":60543
	   ,"total_storage":63,"free_storage":45,"total_mem":53663773,"free_mem":563,"upload_bandwidth":83888,"download_bandwidth":728882}`

	*/
	//keepaliversp:=`{"from":"a", "cmd": "keepalive", "type": "rsp", "seq": 4294967295,  "token": "EF02JLGFA09GVNG21F","time":"6783"}`

	/*
	   uploadrsp:=`{"from":"a", "cmd": "Upload", "type": "rsp", "seq": 4294967295,  "token": "EF02JLGFA09GVNG21F"
	   ,"port":838,"recomend_bitrate":677,"sid_storage":["hshh","0.1.1.2"],"sid_list":["1","2","3","4"]}`
	*/
	/*uploadreq:=`{"from":"a", "cmd": "Upload", "type": "req", "seq": 4294967295,  "token": "EF02JLGFA09GVNG21F","task_id":"67","uid":"6732","fid":"hshh"
	  ,"sid_list":["1","2","3","4"]}`
	*/

	/*downloadrsp:=`{"from":"a", "cmd": "Download", "type": "rsp", "seq": 4294967295,  "token": "EF02JLGFA09GVNG21F"
	  ,"port":838,"recomend_bitrate":677,"sid_list":["1","2","3","4"]}`
	*/
	downloadreq := `{"from":"a", "cmd": "Download", "type": "req", "seq": 4294967295,"task_id":"67","uid":"6732","fid":"hshh"
	,"sid_list":["1","2","3","4"],"max_bitrate":7883,"sid_storage":["a","b","c","10.0.2.3"]}`

	/*downloadreq := `{"from":"a", "cmd": "Download", "type": "req", "seq": 4294967295,"req": {
	  "login": {
	      "version": "0.0.1",
	      "os": "linux centos 6.5 x86_64",
	      "..."
	  }}`*/
	fmt.Println("*is printing return value***", Convert(downloadreq))
}

//================================================================================================================convert function==================>

func Convert(text string) []string {
	fmt.Println()
	// This JSON contains an int array.
	/*loginreq:=`{"from":"a", "cmd": "Download", "type": "req", "seq": 4294967295,  "token": "EF02JLGFA09GVNG21F" ,"yangming":890,"version":"3.0","net_addr":"172.168.0.12","os":"windows","storage":[{"name":"s1","type":"local","net":"china_unicom"}],"total_storage":63,"free_storage":45
	  ,"cpu_usage":65677,"total_mem":88993,"free_mem":563,"upload_bandwidth":83888,"download_bandwidth":728882

	  }`*/
	var tag1 Message

	//init a var ===allocate memory locatin
	bytes := []byte(text) //=======================================================================================================================>json input position======>
	// Unmarshal JSON to Result struct.
	json.Unmarshal(bytes, &tag1)
	fmt.Println("cmd is=====>", tag1.Cmd, "================", "type is  ========================>", tag1.Type)
	fmt.Println("")

	fmt.Println("                     ")

	//=================================================login=======================
	if tag1.Cmd == "Login" {
		//another line
		if tag1.Type == "req" {

			var tag2 LoginReq
			logintest := []byte(text)
			json.Unmarshal(logintest, &tag2)
			return []string{"info_more=============>", tag1.From, tag1.Cmd, tag1.Token, tag2.Version, tag2.NetAddr, tag2.OS, fmt.Sprint(tag2.CpuUsage),
				fmt.Sprint(tag2.TotalMem), fmt.Sprint(tag2.FreeMem), fmt.Sprint(tag2.UploadBandwidth), fmt.Sprint(tag2.DownloadBandwidth)}

		} else if tag1.Type == "rsp" {
			var tag3 LoginRsp
			logintest1 := []byte(text)
			json.Unmarshal(logintest1, &tag3)
			return []string{tag3.Token, tag3.ZoneId, fmt.Sprint(tag3.AliveInterval), fmt.Sprint(tag3.TransferNotifyInterval), tag3.Time, tag3.Version}

		} else {
			return []string{"illegal info"}

		}
		//return []string{"go go gog"}
		//================================================register================have no else============>
	} else if tag1.Cmd == "Register" {
		if tag1.Type == "rsp" {
			var register RegisterRsp
			json.Unmarshal(bytes, &register)
			return []string{"register=======>rsp", register.Tnid}
		} else if tag1.Type == "req" {
			var register RegisterReq
			json.Unmarshal(bytes, &register)
			return []string{"register req=====>", register.MachineCode}
			//==================================================upload=============================================>
		}
	} else if tag1.Cmd == "Upload" {
		if tag1.Type == "rsp" {
			var tag4 UploadTaskRsp
			json.Unmarshal(bytes, &tag4)
			return []string{"Upload=======>rsp", fmt.Sprint(tag4.Port), fmt.Sprint(tag4.RecomendBitrate), fmt.Sprint(tag4.SidList)}
		} else if tag1.Type == "req" {
			var tag4 UploadTaskReq
			json.Unmarshal(bytes, &tag4)
			return []string{tag4.TaskId, tag4.Uid, tag4.Fid, fmt.Sprint(tag4.SidList)}
			//==================================================keepalive========================================>
		}
	} else if tag1.Cmd == "keepalive" {
		if tag1.Type == "rsp" {
			var keepalive KeepaliveRsp
			json.Unmarshal(bytes, &keepalive)
			return []string{"keepalive===========>", keepalive.Token, keepalive.Time}
		} else if tag1.Type == "req" {
			var keepalive KeepaliveReq
			json.Unmarshal(bytes, &keepalive)
			return []string{fmt.Sprint(keepalive.FreeMem), fmt.Sprint(keepalive.UploadBandwidth),
				fmt.Sprint(keepalive.TotalMem), fmt.Sprint(keepalive.FreeMem), fmt.Sprint(keepalive.TotalStorage),
				fmt.Sprint(keepalive.FreeStorage), fmt.Sprint(keepalive.UploadBandwidth), fmt.Sprint(keepalive.DownloadBandwidth), fmt.Sprint(keepalive.CpuUsage)}

		}

		//==================================================download===========================================>
	} else if tag1.Cmd == "Download" {
		if tag1.Type == "req" {
			//var tag5 Req.DownloadTask
			//json.Unmarshal(bytes, &tag1)
			return []string{"i love bayby"}

		} else if tag1.Type == "rsp" {
			var tag5 DownloadTaskRsp
			json.Unmarshal(bytes, &tag5)
			return []string{"DownloadTask=======>rsp", fmt.Sprint(tag5.Port), fmt.Sprint(tag5.RecomendBitrate), fmt.Sprint(tag5.SidList)}
		}

		//================================================test finish ,but no get result===================================>
	} else {

		fmt.Println("these line will be launched")
		fmt.Println("i lovehsjjjsjjjjj")
		fmt.Println("")
		return []string{"what is fuls"}

	}
	return []string{"test finish"}
}
