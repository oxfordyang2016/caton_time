package main

// #cgo CFLAGS: -I ./../libf2tp/include
// #cgo CFLAGS: -I ./../libf2tp_storage/include
// #cgo CFLAGS: -I ./../libf2tp_mate/include
// #cgo LDFLAGS: -L./../libf2tp_storage/ -lf2tp_storage
// #cgo LDFLAGS: -L./../libf2tp_mate/ -lf2tp_mate
// #cgo LDFLAGS: -L./../libf2tp/libs/linux/static -lf2tp
// #cgo LDFLAGS: -lpthread -lcrypto
//
// #include <stdio.h>
// #include <stdlib.h>
// #include <string.h>
// #include "f2tp.h"
// #include "f2tp_storage.h"
// #include "udp_notify.h"
//
// int on_event(F2tpHandle handle, int event, void *param, void *userdata)
// {
//  fprintf(stderr, "call me !!!!\n");
//	UdpNotify * un = (UdpNotify *)userdata;
//	udp_notify_event(un, event, param);
//  fprintf(stderr, "call me over !!!!\n");
// }
import "C"

import (
	"fmt"
	clog "github.com/cihub/seelog"
	"os"
	"time"
	"unsafe"
)

func (self *Task) RunF2tpServer() {
	clog.Tracef("%s run f2tp server", self)
	// C.f2tp_setLogLevel(C.LOG_DEBUG)

	var info C.struct_StdStorageInfo
	path := C.CString("/tmp/cydex_data/")
	C.strcpy(&info.path[0], path)
	C.free(unsafe.Pointer(path))
	io_interface := C.std_storage_new(&info)
	if io_interface == nil {
		clog.Errorf("%s new std storage failed", self)
		return
	}
	defer C.std_storage_free(io_interface)

	// event notify
	notify := C.udp_notify_new(nil, 5)
	if notify == nil {
		clog.Errorf("%s udp notify new failed", self)
		return
	}
	host := C.CString("127.0.0.1")
	C.udp_notify_add_target(notify, host, C.int(3322))
	C.free(unsafe.Pointer(host))
	defer C.udp_notify_free(notify)

	// 创建侦听端口
	listen_handle := C.f2tp_open(C.F2TP_SERVER)
	if listen_handle == nil {
		clog.Errorf("%s f2tp open failed", self)
		return
	}
	defer func() {
		fmt.Fprintf(os.Stderr, "%s close f2tp\n", self)
		C.f2tp_close(listen_handle)
		fmt.Fprintf(os.Stderr, "%s close f2tp 2\n", self)
	}()

	// 设置io interface
	C.f2tp_setCallbacks(listen_handle, io_interface)

	// 绑定地址和端口
	host = C.CString("0.0.0.0")
	ret := C.f2tp_bind(listen_handle, host, C.ushort(self.Port))
	C.free(unsafe.Pointer(host))
	if ret != 0 {
		clog.Errorf("%s f2tp bind error\n", self)
		return
	}

	// listen
	ret = C.f2tp_listen(listen_handle, C.int(1))
	if ret != 0 {
		clog.Errorf("%s f2tp listen error\n", self)
		return
	}

	for {
		fmt.Println("before accept")
		conn_handle := C.f2tp_accept(listen_handle, C.int(10*1000*1000))
		// TODO 出错处理, 超过一定次数就退出了
		if conn_handle == nil {
			err := C.f2tp_getLastError(listen_handle)
			clog.Errorf("%s f2tp accept failed, err:%d", self, err)
			return
		}

		// 设置加密方式, FIXME 目前是固定的
		encryptoion_type := C.int(2)
		C.f2tp_setOpt(conn_handle, C.O_ENCRYPTTYPE, unsafe.Pointer(&encryptoion_type), 4)
		// 设置限速, FIXME 目前是固定的
		max_bitrate := C.int(10 * 1024 * 1024)
		C.f2tp_setOpt(conn_handle, C.O_MAXBITRATE, unsafe.Pointer(&max_bitrate), 4)

		// get name
		sid_a := [256]C.char{}
		sid := &sid_a[0]
		C.f2tp_getName(conn_handle, sid, 256)
		// fill info
		go_sid := C.GoString(sid)
		cur_state := self.segs_state[go_sid]
		if cur_state == nil {
			// FIXME
			panic("nnnnn")
		}
		self.cur_state = cur_state
		cur_state.State = "transferring"

		// 设置事件回调
		var info C.struct_F2tpNotifyInfo
		taskid := C.CString(self.Tid)
		C.strcpy(&info.taskid[0], taskid)
		C.strcpy(&info.sid[0], sid)
		C.udp_notify_set_info(notify, &info)
		C.free(unsafe.Pointer(taskid))

		callback := C.EventHandler(C.on_event)
		C.f2tp_registNotifier(conn_handle, callback, unsafe.Pointer(notify))

	over_label:
		for {
			select {
			case <-time.Tick(3 * time.Second):
				var stat C.struct_F2tpStat_t
				C.f2tp_getOpt(conn_handle, C.O_STATISTICS, unsafe.Pointer(&stat), C.int(unsafe.Sizeof(stat)))
				clog.Tracef("%s total bytes: %d, bitrate:%d", self, stat.totalBytes, stat.bitrate)
				cur_state.TotalBytes = uint64(stat.totalBytes)
				cur_state.Bitrate = uint64(stat.bitrate)

				self.notify_chan <- cur_state

			case event := <-self.event_chan:
				fmt.Fprintf(os.Stderr, "%s from udp notify: %+v\n", self, event.Event)
				// clog.Tracef("%s from udp notify: %+v", self, event.Event)
				switch event.Event {
				case int(C.EVT_COMPLETE):
					cur_state.State = "end"
				case int(C.EVT_INTERRUPT):
					cur_state.State = "interrupted"
				}

				var stat C.struct_F2tpStat_t
				C.f2tp_getOpt(conn_handle, C.O_STATISTICS, unsafe.Pointer(&stat), C.int(unsafe.Sizeof(stat)))
				cur_state.TotalBytes = uint64(stat.totalBytes)
				cur_state.Bitrate = uint64(stat.bitrate)

				// notify
				self.notify_chan <- cur_state
				self.cur_state = nil

				break over_label
			}
		}

		fmt.Fprintf(os.Stderr, "%s before conn close\n", self)
		C.f2tp_close(conn_handle)
		fmt.Fprintf(os.Stderr, "%s after conn close\n", self)

		// clog.Infof("%s sid:%s is over", self, go_sid)

		if self.IsOver() {
			break
		}
	}

	fmt.Fprintf(os.Stderr, "%s over\n", self)
	// clog.Infof("%s is over", self)

	// C.udp_notify_free(notify)
	// C.std_storage_free(io_interface)
	// C.f2tp_close(listen_handle)
}
