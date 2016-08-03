package transfer

import (
	"encoding/json"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_Message(t *testing.T) {
	Convey("Test Message", t, func() {
		Convey("Request msg", func() {
			req := NewReqMessage("from", "Login", "token", 1234)
			So(req.From, ShouldEqual, "from")
			So(req.Cmd, ShouldEqual, "Login")
			So(req.Type, ShouldEqual, "req")
			So(req.Token, ShouldEqual, "token")
			So(req.Seq, ShouldEqual, 1234)
			So(req.Req, ShouldNotBeNil)
			So(req.IsReq(), ShouldBeTrue)
			Convey("build rsp", func() {
				rsp := req.BuildRsp()
				So(rsp.IsRsp(), ShouldBeTrue)
				So(rsp.Cmd, ShouldEqual, req.Cmd)
				So(rsp.Seq, ShouldEqual, req.Seq)
				So(rsp.Rsp, ShouldNotBeNil)
			})
		})
		Convey("Test response msg", func() {
			m := NewRspMessage("from", "Keepalive", "token", 4321)
			So(m.IsRsp(), ShouldBeTrue)
			So(m.From, ShouldEqual, "from")
			So(m.Cmd, ShouldEqual, "Keepalive")
			So(m.Token, ShouldEqual, "token")
			So(m.Seq, ShouldEqual, 4321)
			So(m.Rsp, ShouldNotBeNil)
		})
	})
}
func Test_RegisterReq(t *testing.T) {
	Convey("Test Login Req", t, func() {
		Convey("Unmarshal", func() {

			text := []byte(`{"from":"a", "cmd": "register", "type": "req",
		"seq": 11,"task_id":"67","uid":"6732","token": "EF02JLGFA09GVNG21F","req":{"register":{"machine_code":"mac"}}}`)

			var msg Message
			err := json.Unmarshal(text, &msg)
			So(err, ShouldBeNil)
			So(msg.IsReq(), ShouldBeTrue)
			So(msg.From, ShouldEqual, "a")
			So(msg.Cmd, ShouldEqual, "register")
			So(msg.Seq, ShouldEqual, 11)
			So(msg.Req, ShouldNotBeNil)
			So(msg.Req.Register, ShouldNotBeNil)
			So(msg.Req.Register.MachineCode, ShouldEqual, "mac")

		})
		Convey("Marshal", func() {
		})
	})
}
func Test_Registerrsp(t *testing.T) {
	Convey("Test Login Req", t, func() {
		Convey("Unmarshal", func() {

			text1 := []byte(`{"from":"a", "cmd": "register", "type": "rsp",
		"seq": 11,"task_id":"67","uid":"6732","token": "EF02JLGFA09GVNG21F","rsp":{"register":{"tnid":"56"}}}`)

			var msg1 Message
			err1 := json.Unmarshal(text1, &msg1)
			So(err1, ShouldBeNil)
			So(msg1.IsRsp(), ShouldBeTrue)
			So(msg1.From, ShouldEqual, "a")
			So(msg1.Cmd, ShouldEqual, "register")
			So(msg1.Seq, ShouldEqual, 11)
			So(msg1.Rsp, ShouldNotBeNil)
			So(msg1.Rsp.Register, ShouldNotBeNil)
			So(msg1.Rsp.Register.Tnid, ShouldEqual, "56")
		})
		Convey("Marshal", func() {
		})
	})
}
func Test_Loginrsp(t *testing.T) {
	Convey("Test Login rsp", t, func() {
		Convey("Unmarshal", func() {

			text1 := []byte(`{"from":"a", "cmd": "Login", "type": "rsp",
		"seq": 11,"task_id":"67","uid":"6732","token": "EF02JLGFA09GVNG21F","rsp":{"code":0,"reason":"ok","login":{
	   "alive_interval":54,"free_mem":563,"upload_bandwidth":83888,
	   "transfer_notify_interval":5366,"time":"653","download_bandwidth":728882,"version":"uis","zone_id":"564"
		,"total_storage":63,"free_storage":45}}}`)

			var msg1 Message
			err1 := json.Unmarshal(text1, &msg1)
			So(err1, ShouldBeNil)
			So(msg1.IsRsp(), ShouldBeTrue)
			So(msg1.From, ShouldEqual, "a")
			So(msg1.Cmd, ShouldEqual, "Login")
			So(msg1.Seq, ShouldEqual, 11)
			So(msg1.Rsp, ShouldNotBeNil)
			So(msg1.Rsp.Login, ShouldNotBeNil)
			So(msg1.Rsp.Login.AliveInterval, ShouldEqual, 54)
			So(msg1.Rsp.Login.Time, ShouldEqual, "653")
			So(msg1.Rsp.Login.TransferNotifyInterval, ShouldEqual, 5366)
			So(msg1.Rsp.Login.Version, ShouldEqual, "uis")
			So(msg1.Rsp.Login.ZoneId, ShouldEqual, "564")
		})

	})
}
func Test_Loginreq(t *testing.T) {
	Convey("Test Login Req", t, func() {
		Convey("Unmarshal", func() {

			text := []byte(`{"from":"a", "cmd": "Login", "type": "req",
			"seq": 11,"req":{"login":{"version":"890","os":"windows","net_speed":432
			,"net_addr":"here","cpu_usage":54,"total_mem":78,"free_mem":78,
              "download_bandwidth":72,"upload_bandwidth":83888,
		  "free_mem":563,"upload_bandwidth":83888,"storages":[{"name":"s1","type":"local","net":"china_unicom"}],
		   "total_storage":63,"free_storage":45
		   }}}`)

			var msg Message
			err := json.Unmarshal(text, &msg)
			So(err, ShouldBeNil)
			So(msg.IsReq(), ShouldBeTrue)
			So(msg.From, ShouldEqual, "a")
			So(msg.Cmd, ShouldEqual, "Login")
			So(msg.Seq, ShouldEqual, 11)
			So(msg.Req, ShouldNotBeNil)
			So(msg.Req.Login, ShouldNotBeNil)
			login := msg.Req.Login
			So(login.Version, ShouldEqual, "890")
			So(login.OS, ShouldEqual, "windows")
			So(login.NetSpeed, ShouldEqual, 432)
			So(login.NetAddr, ShouldEqual, "here")
			So(login.CpuUsage, ShouldEqual, 54)
			So(login.TotalMem, ShouldEqual, 78)
			So(login.FreeMem, ShouldEqual, 563)
			So(login.DownloadBandwidth, ShouldEqual, 72)
			So(login.UploadBandwidth, ShouldEqual, 83888)
			So(login.TotalStorage, ShouldEqual, 63)
			So(login.FreeStorage, ShouldEqual, 45)
			var k = *(login.Storage[0]) //this way is influenced by address
			//k := login.Storage[0]
			So(k.Name, ShouldEqual, "s1")
			So(k.Type, ShouldEqual, "local")
			So(k.Net, ShouldEqual, "china_unicom")
		})

	})
}
func Test_Uploadreq(t *testing.T) {
	Convey("Test Login Req", t, func() {
		Convey("Unmarshal", func() {

			text := []byte(`{"from":"a", "cmd": "upload", "type": "req",
			"seq": 11,"req":{"uploadtask":{
         "task_id":"i",
		    "uid":"u",
		    "fid":"f",
		    "sid_list":["1"],
		    "sid_storage":["hshh","0.1.1.2"],
		    "max_bitrate":673

		   }}}`)

			var msg Message
			err := json.Unmarshal(text, &msg)
			So(err, ShouldBeNil)
			So(msg.IsReq(), ShouldBeTrue)
			So(msg.From, ShouldEqual, "a")
			So(msg.Cmd, ShouldEqual, "upload")
			So(msg.Seq, ShouldEqual, 11)
			So(msg.Req, ShouldNotBeNil)
			So(msg.Req.UploadTask, ShouldNotBeNil)
			upload := msg.Req.UploadTask
			So(upload.TaskId, ShouldEqual, "i")
			So(upload.Uid, ShouldEqual, "u")
			So(upload.SidList[0], ShouldEqual, "1")

			So(upload.Fid, ShouldEqual, "f")
		})

	})
}
func Test_uploadrsp(t *testing.T) {
	Convey("Test upload rsp", t, func() {
		Convey("Unmarshal", func() {

			text1 := []byte(`{"from":"a", "cmd": "upload", "type": "rsp",
			"seq": 42,"rsp":{"uploadtask":{
		    "sid_list":["1"],
		    "sid_storage":["hshh","0.1.1.2"],
		    "recomend_bitrate":673,
		    "port":564,
		    "max_bitrate":673

		   }}}`)

			var msg1 Message
			err1 := json.Unmarshal(text1, &msg1)
			So(err1, ShouldBeNil)
			So(msg1.IsRsp(), ShouldBeTrue)
			So(msg1.From, ShouldEqual, "a")
			So(msg1.Cmd, ShouldEqual, "upload")
			So(msg1.Seq, ShouldEqual, 42)
			So(msg1.Rsp, ShouldNotBeNil)
			So(msg1.Rsp.UploadTask, ShouldNotBeNil)
			upload := msg1.Rsp.UploadTask
			So(upload.SidList[0], ShouldEqual, "1")
			So(upload.Port, ShouldEqual, 564)
			//So(upload.SidStorage[0], ShouldEqual, "hshh")
			So(upload.RecomendBitrate, ShouldEqual, 673)

		})

	})
}
func Test_Downloadreq(t *testing.T) {
	Convey("Test Login Req", t, func() {
		Convey("Unmarshal", func() {

			text := []byte(`{"from":"a", "cmd": "download", "type": "req",
			"seq": 42,"req":{"downloadtask":{
         "task_id":"i",
		    "uid":"uid",
		    "fid":"fid",
		    "sid_list":["1","2"],
		    "sid_storage":["hshh","0.1.1.2"],
		    "max_bitrate":673

		   }}}`)

			var msg Message
			err := json.Unmarshal(text, &msg)
			So(err, ShouldBeNil)
			So(msg.IsReq(), ShouldBeTrue)
			So(msg.From, ShouldEqual, "a")
			So(msg.Cmd, ShouldEqual, "download")
			So(msg.Seq, ShouldEqual, 42)
			So(msg.Req, ShouldNotBeNil)
			So(msg.Req.DownloadTask, ShouldNotBeNil)
			download := msg.Req.DownloadTask
			So(download.TaskId, ShouldEqual, "i")
			So(download.Uid, ShouldEqual, "uid")
			So(download.SidList[0], ShouldEqual, "1")
			So(download.SidList[1], ShouldEqual, "2")
			So(download.MaxBitrate, ShouldEqual, 673)

		})

	})
}
func Test_downloadrsp(t *testing.T) {
	Convey("Test upload rsp", t, func() {
		Convey("Unmarshal", func() {

			text1 := []byte(`{"from":"a", "cmd": "download", "type": "rsp",
			"seq": 42,"rsp":{"downloadtask":{
		    "sid_list":["1","2"],
		    "recomend_bitrate":673,
		    "port":564

		   }}}`)

			var msg1 Message
			err1 := json.Unmarshal(text1, &msg1)
			So(err1, ShouldBeNil)
			So(msg1.IsRsp(), ShouldBeTrue)
			So(msg1.From, ShouldEqual, "a")
			So(msg1.Cmd, ShouldEqual, "download")
			So(msg1.Seq, ShouldEqual, 42)
			So(msg1.Rsp, ShouldNotBeNil)
			So(msg1.Rsp.DownloadTask, ShouldNotBeNil)
			download := msg1.Rsp.DownloadTask
			fmt.Println((download.SidList)[0])
			So(download.SidList[0], ShouldEqual, "1")
			So(download.SidList[1], ShouldEqual, "2")
			So(download.Port, ShouldEqual, 564)
			So(download.RecomendBitrate, ShouldEqual, 673)

		})

	})
}
func Test_Keepalivereq(t *testing.T) {
	Convey("Test keepalive Req", t, func() {
		Convey("Unmarshal", func() {

			text := []byte(`{"from":"a", "cmd": "keepalive", "type": "req",
			"seq": 42,"req":{"keepalive":{"cpu_usage":54,"total_mem":78,"free_mem":78,
              "download_bandwidth":728882,"upload_bandwidth":83888,
		   "total_storage":63,"free_storage":45
		   }}}`)

			var msg Message
			err := json.Unmarshal(text, &msg)
			So(err, ShouldBeNil)
			So(msg.IsReq(), ShouldBeTrue)
			So(msg.From, ShouldEqual, "a")
			So(msg.Cmd, ShouldEqual, "keepalive")
			So(msg.Seq, ShouldEqual, 42)
			So(msg.Req, ShouldNotBeNil)
			So(msg.Req.Keepalive, ShouldNotBeNil)
			alive := msg.Req.Keepalive
			So(alive.CpuUsage, ShouldEqual, 54)
			So(alive.TotalMem, ShouldEqual, 78)
			So(alive.FreeMem, ShouldEqual, 78)
			So(alive.DownloadBandwidth, ShouldEqual, 728882)
			So(alive.UploadBandwidth, ShouldEqual, 83888)
			So(alive.TotalStorage, ShouldEqual, 63)
			So(alive.FreeStorage, ShouldEqual, 45)

		})

	})
}
func Test_Keepaliversp(t *testing.T) {
	Convey("Test keepalive rsp", t, func() {
		Convey("Unmarshal", func() {

			text1 := []byte(`{"from":"a", "cmd": "keepalive", "type": "rsp",
			"seq": 42,"rsp":{"keepalive":{
		    "time":"564",
		    "token":"token"

		   }}}`)

			var msg1 Message
			err1 := json.Unmarshal(text1, &msg1)
			So(err1, ShouldBeNil)
			So(msg1.IsRsp(), ShouldBeTrue)
			So(msg1.From, ShouldEqual, "a")
			So(msg1.Cmd, ShouldEqual, "keepalive")
			So(msg1.Seq, ShouldEqual, 42)
			So(msg1.Rsp, ShouldNotBeNil)
			So(msg1.Rsp.Keepalive, ShouldNotBeNil)
			alive := msg1.Rsp.Keepalive
			So(alive.Time, ShouldEqual, "564")
			So(alive.Token, ShouldEqual, "token")

		})

	})
}

func Test_2loginrsp(t *testing.T) {
	Convey("Test keepalive rsp", t, func() {
		Convey("Unmarshal", func() {

			var loginrsp = Message{Base{"kksk", "KKD", "JSJ", 563, "73HJSKA_K"}, nil,
				&Response{1, "ok", nil, &LoginRsp{"token", "zoneid", 33, 45, "time", "version"}, nil, nil, nil}}

			k := &loginrsp

			k3, _ := json.Marshal(k)
			d1 := []byte(`{"from":"kksk","cmd":"KKD","type":"JSJ","seq":563,"token":"73HJSKA_K","rsp":{"code":1,"reason":"ok","login":{"token":"token","zone_id":"zoneid","alive_interval":33,"transfer_notify_interval":45,"time":"time","version":"version"}}}`)
			//warning:another line impact result
			So(string(k3), ShouldEqual, string(d1))

		})

	})
}
func Test_2loginreq(t *testing.T) {
	Convey("Testloginreq", t, func() {
		Convey("marshal", func() {

			var loginreq = Message{Base{"kksk", "KKD", "JSJ", 563, "73HJSKA_K"}, &Request{nil, &LoginReq{"version",
				"netaddrr", "os", 34, []*StorageInfo{&StorageInfo{"name", "a", "mob"}}, 34, 34, 34, 45, 32, 34, 54}, nil, nil, nil, nil}, nil}

			k := &loginreq

			k3, _ := json.Marshal(k)
			//d1 := []byte(`{"from":"kksk","cmd":"KKD","type":"JSJ","seq":563,"token":"73HJSKA_K","rsp":{"code":1,"reason":"ok","login":{"token":"token","zone_id":"zoneid","alive_interval":33,"transfer_notify_interval":45,"time":"time","version":"version"}}}`)
			d1 := []byte(`{"from":"kksk","cmd":"KKD","type":"JSJ","seq":563,"token":"73HJSKA_K","req":{"login":{"version":"version","net_addr":"netaddrr","os":"os","net_speed":34,"storages":[{"name":"name","type":"a","net":"mob"}],"total_storage":34,"free_storage":34,"cpu_usage":34,"total_mem":45,"free_mem":32,"upload_bandwidth":34,"download_bandwidth":54}}}`)
			So(string(k3), ShouldEqual, string(d1))

		})

	})
}
func Test_2Registereq(t *testing.T) {
	Convey("Testregisterreq", t, func() {
		Convey("marshal", func() {

			var req = Message{Base{"kksk", "KKD", "JSJ", 563, "73HJSKA_K"}, &Request{&RegisterReq{"macode"}, nil, nil, nil, nil, nil}, nil}

			k := &req

			k3, _ := json.Marshal(k)
			//d1 := []byte(`{"from":"kksk","cmd":"KKD","type":"JSJ","seq":563,"token":"73HJSKA_K","rsp":{"code":1,"reason":"ok","login":{"token":"token","zone_id":"zoneid","alive_interval":33,"transfer_notify_interval":45,"time":"time","version":"version"}}}`)
			d1 := []byte(`{"from":"kksk","cmd":"KKD","type":"JSJ","seq":563,"token":"73HJSKA_K","req":{"register":{"machine_code":"macode"}}}`)
			So(string(k3), ShouldEqual, string(d1))
		})

	})
}
func Test_2Registersp(t *testing.T) {
	Convey("Testregisterreq", t, func() {
		Convey("marshal", func() {

			var rsp = Message{Base{"kksk", "KKD", "JSJ", 563, "73HJSKA_K"}, nil, &Response{1, "ok", &RegisterRsp{""}, nil, nil, nil, nil}}

			k := &rsp

			k3, _ := json.Marshal(k)
			//d1 := []byte(`{"from":"kksk","cmd":"KKD","type":"JSJ","seq":563,"token":"73HJSKA_K","rsp":{"code":1,"reason":"ok","login":{"token":"token","zone_id":"zoneid","alive_interval":33,"transfer_notify_interval":45,"time":"time","version":"version"}}}`)
			d1 := []byte(`{"from":"kksk","cmd":"KKD","type":"JSJ","seq":563,"token":"73HJSKA_K","rsp":{"code":1,"reason":"ok","register":{"tnid":""}}}`)
			So(string(k3), ShouldEqual, string(d1))
		})

	})
}
func Test_2Uploadreq(t *testing.T) {
	Convey("Testregisterreq", t, func() {
		Convey("marshal", func() {

			var req = Message{Base{"kksk", "KKD", "JSJ", 563, "73HJSKA_K"}, &Request{nil, nil, nil, &UploadTaskReq{"taskid", "uid", "fid", []string{"sid"}}, nil, nil}, nil}

			k := &req

			k3, _ := json.Marshal(k)
			//d1 := []byte(`{"from":"kksk","cmd":"KKD","type":"JSJ","seq":563,"token":"73HJSKA_K","rsp":{"code":1,"reason":"ok","login":{"token":"token","zone_id":"zoneid","alive_interval":33,"transfer_notify_interval":45,"time":"time","version":"version"}}}`)
			d1 := []byte(`{"from":"kksk","cmd":"KKD","type":"JSJ","seq":563,"token":"73HJSKA_K","req":{"uploadtask":{"task_id":"taskid","uid":"uid","fid":"fid","sid_list":["sid"]}}}`)
			So(string(k3), ShouldEqual, string(d1))
		})
	})
}
func Test_2Uploadrsp(t *testing.T) {
	Convey("Testregisterreq", t, func() {
		Convey("marshal", func() {

			var rsp = Message{Base{"kksk", "KKD", "JSJ", 563, "73HJSKA_K"}, nil, &Response{1, "ok", nil, nil, nil, &UploadTaskRsp{[]string{"sidlist"}, []string{"sidstore"}, 43, 435}, nil}}

			k := &rsp

			k3, _ := json.Marshal(k)
			//d1 := []byte(`{"from":"kksk","cmd":"KKD","type":"JSJ","seq":563,"token":"73HJSKA_K","rsp":{"code":1,"reason":"ok","login":{"token":"token","zone_id":"zoneid","alive_interval":33,"transfer_notify_interval":45,"time":"time","version":"version"}}}`)
			d1 := []byte(`{"from":"kksk","cmd":"KKD","type":"JSJ","seq":563,"token":"73HJSKA_K","rsp":{"code":1,"reason":"ok","uploadtask":{"sid_list":["sidlist"],"sid_sorage":["sidstore"],"port":43,"recomend_bitrate":435}}}`)
			So(string(k3), ShouldEqual, string(d1))
		})

	})
}
func Test_2Downloadreq(t *testing.T) {
	Convey("Testregisterreq", t, func() {
		Convey("marshal", func() {

			var req = Message{Base{"kksk", "KKD", "JSJ", 563, "73HJSKA_K"}, &Request{nil, nil, nil, nil, &DownloadTaskReq{"taskid", "uid", "fid", []string{"sidlist"},
				[]string{"sidstore"}, 232}, nil}, nil}

			k := &req

			k3, _ := json.Marshal(k)
			//d1 := []byte(`{"from":"kksk","cmd":"KKD","type":"JSJ","seq":563,"token":"73HJSKA_K","rsp":{"code":1,"reason":"ok","login":{"token":"token","zone_id":"zoneid","alive_interval":33,"transfer_notify_interval":45,"time":"time","version":"version"}}}`)
			d1 := []byte(`{"from":"kksk","cmd":"KKD","type":"JSJ","seq":563,"token":"73HJSKA_K","req":{"downloadtask":{"task_id":"taskid","uid":"uid","fid":"fid","sid_list":["sidlist"],"sid_storage":["sidstore"],"max_bitrate":232}}}`)
			So(string(k3), ShouldEqual, string(d1))
		})
	})
}
func Test_2Downloadrsp(t *testing.T) {
	Convey("download", t, func() {
		Convey("marshal", func() {

			var rsp = Message{Base{"kksk", "KKD", "JSJ", 563, "73HJSKA_K"}, nil, &Response{1, "ok", nil, nil, nil, nil, &DownloadTaskRsp{[]string{"kid"}, 32, 322}}}

			k := &rsp

			k3, _ := json.Marshal(k)
			//d1 := []byte(`{"from":"kksk","cmd":"KKD","type":"JSJ","seq":563,"token":"73HJSKA_K","rsp":{"code":1,"reason":"ok","login":{"token":"token","zone_id":"zoneid","alive_interval":33,"transfer_notify_interval":45,"time":"time","version":"version"}}}`)
			d1 := []byte(`{"from":"kksk","cmd":"KKD","type":"JSJ","seq":563,"token":"73HJSKA_K","rsp":{"code":1,"reason":"ok","downloadtask":{"sid_list":["kid"],"port":32,"recomend_bitrate":322}}}`)
			So(string(k3), ShouldEqual, string(d1))
		})

	})
}

func Test_2Keepalivereq(t *testing.T) {
	Convey("Testregisterreq", t, func() {
		Convey("marshal", func() {

			var req = Message{Base{"kksk", "KKD", "JSJ", 563, "73HJSKA_K"}, &Request{nil, nil, &KeepaliveReq{
				21, 34, 45, 45, 44, 54, 45}, nil, nil, nil}, nil}

			k := &req

			k3, _ := json.Marshal(k)
			//d1 := []byte(`{"from":"kksk","cmd":"KKD","type":"JSJ","seq":563,"token":"73HJSKA_K","rsp":{"code":1,"reason":"ok","login":{"token":"token","zone_id":"zoneid","alive_interval":33,"transfer_notify_interval":45,"time":"time","version":"version"}}}`)
			d1 := []byte(`{"from":"kksk","cmd":"KKD","type":"JSJ","seq":563,"token":"73HJSKA_K","req":{"keepalive":{"cpu_usage":21,"total_storage":34,"free_storage":45,"total_mem":45,"free_mem":44,"upload_bandwidth":54,"download_bandwidth":45}}}`)
			So(string(k3), ShouldEqual, string(d1))
		})
	})
}
func Test_2keepaliversp(t *testing.T) {
	Convey("download", t, func() {
		Convey("marshal", func() {

			var rsp = Message{Base{"kksk", "KKD", "JSJ", 563, "73HJSKA_K"}, nil, &Response{1, "ok", nil, nil, &KeepaliveRsp{"token", "time"}, nil, nil}}

			k := &rsp

			k3, _ := json.Marshal(k)
			//d1 := []byte(`{"from":"kksk","cmd":"KKD","type":"JSJ","seq":563,"token":"73HJSKA_K","rsp":{"code":1,"reason":"ok","login":{"token":"token","zone_id":"zoneid","alive_interval":33,"transfer_notify_interval":45,"time":"time","version":"version"}}}`)
			d1 := []byte(`{"from":"kksk","cmd":"KKD","type":"JSJ","seq":563,"token":"73HJSKA_K","rsp":{"code":1,"reason":"ok","keepalive":{"token":"token","time":"time"}}}`)
			So(string(k3), ShouldEqual, string(d1))
		})

	})
}
func Test_Transfernotifyeq(t *testing.T) {
	Convey("Testregisterreq", t, func() {
		Convey("marshal", func() {

			var req = Message{Base{"kksk", "KKD", "JSJ", 563, "73HJSKA_K"}, &Request{nil, nil, nil, nil, nil, &TransferNotifyReq{[]*TaskState{&TaskState{"taskid", "state", 43, 45}}}}, nil}

			k := &req

			k3, _ := json.Marshal(k)
			//d1 := []byte(`{"from":"kksk","cmd":"KKD","type":"JSJ","seq":563,"token":"73HJSKA_K","rsp":{"code":1,"reason":"ok","login":{"token":"token","zone_id":"zoneid","alive_interval":33,"transfer_notify_interval":45,"time":"time","version":"version"}}}`)
			d1 := []byte(`{"from":"kksk","cmd":"KKD","type":"JSJ","seq":563,"token":"73HJSKA_K","req":{"transfernotify":{"task_list":[{"task_id":"taskid","state":"state","total_bytes":43,"bitrate":45}]}}}`)
			So(string(k3), ShouldEqual, string(d1))
		})
	})
}
