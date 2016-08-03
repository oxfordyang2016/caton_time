/* EchoServer
 */
package main

import (
	"fmt"
	"net/http"
	"os"
	// "io"
	"encoding/json"
	"golang.org/x/net/websocket"
)

//=========================start==========================================
type Base struct {
	From  string `json:"from,omitempty"`
	Cmd   string `json:"cmd"`
	Type  string `json:"type"`
	Seq   uint32 `json:"seq"`
	Token string `json:"token,omitempty"`
}
type Message struct {
	Base
	Req *Request  `json:"req,omitempty"`
	Rsp *Response `json:"rsp,omitempty"`
}
type Request struct {
	Register       *RegisterReq       `json:"register,omitempty"`
	Login          *LoginReq          `json:"login,omitempty"`
	Keepalive      *KeepaliveReq      `json:"keepalive,omitempty"`
	UploadTask     *UploadTaskReq     `json:"uploadtask,omitempty"`
	DownloadTask   *DownloadTaskReq   `json:"downloadtask,omitempty"`
	TransferNotify *TransferNotifyReq `json:"transfernotify,omitempty"`
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
	Version           string         `json:"version"`
	NetAddr           string         `json:"net_addr"`
	OS                string         `json:"os"`
	NetSpeed          uint32         `json:"net_speed"`
	Storage           []*StorageInfo `json:"storages"`
	TotalStorage      uint64         `json:"total_storage"`
	FreeStorage       uint64         `json:"free_storage"`
	CpuUsage          uint32         `json:"cpu_usage"`
	TotalMem          uint64         `json:"total_mem"`
	FreeMem           uint64         `json:"free_mem"`
	UploadBandwidth   uint64         `json:"upload_bandwidth"`
	DownloadBandwidth uint64         `json:"download_bandwidth"`
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
type TransferNotifyReq struct {
	TaskStateList []*TaskState `json:"task_list"`
}

func NewMessage(from, cmd, token, typ string, seq uint32) *Message {
	m := new(Message)
	m.From = from
	m.Cmd = cmd
	m.Token = token
	m.Seq = seq
	m.Type = typ
	switch typ {
	case "req":
		m.Req = new(Request)
	case "rsp":
		m.Rsp = new(Response)
	default:
		return nil
	}
	return m
}
func NewReqMessage(from, cmd, token string, seq uint32) *Message {
	return NewMessage(from, cmd, token, "req", seq)
}
func NewRspMessage(from, cmd, token string, seq uint32) *Message {
	return NewMessage(from, cmd, token, "rsp", seq)
}
func (m *Message) VerifyBase(check_from bool) bool {
	if m.Cmd == "" || (m.Type != "req" && m.Type != "rsp") || (check_from && m.From == "") {
		return false
	}
	return true
}
func (m *Message) IsReq() bool {
	return m.Type == "req"
}
func (m *Message) IsRsp() bool {
	return m.Type == "rsp"
}
func (m *Message) BuildRsp() *Message {
	return NewRspMessage("", m.Cmd, "", m.Seq)
}
func (m *Message) String() string {
	return fmt.Sprintf("<Msg %s %s %d>", m.Type, m.Cmd, m.Seq)
}

//========================end=================================
func Echo(ws *websocket.Conn) {
	fmt.Println("Echoing")

	for n := 0; n < 10; n++ {
		msg := "Hello client " + string(n+48)
		fmt.Println("Sending to client: " + msg)
		err := websocket.Message.Send(ws, msg)

		if err != nil {
			fmt.Println("Can't send")
			break
		}

		var reply string
		err = websocket.Message.Receive(ws, &reply)
		//=======================start=====================
		text := []byte(reply)

		var msg1 Message
		err := json.Unmarshal(text, &msg1)
		fmt.Println("==============from=================>", msg1.From)

		//=====================end==============
		if err != nil {
			fmt.Println("Can't receive")
			break
		}
		fmt.Println("Received from client: " + reply)
	}
}

func main() {

	http.Handle("/", websocket.Handler(Echo))
	err := http.ListenAndServe(":88", nil)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
