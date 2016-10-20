package utils

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_GetHumanSize(t *testing.T) {
	Convey("Test GetHumanSize", t, func() {
		Convey("T", func() {
			var v = uint64(25 * 1024 * 1024 * 1024 * 1024)
			ret := GetHumanSize(v)
			So(ret, ShouldEqual, "25.0T")
		})
		Convey("G", func() {
			var v = uint64(33880356 * 1024)
			ret := GetHumanSize(v)
			So(ret, ShouldEqual, "32.3G")
		})
		Convey("M", func() {
			var v = uint64(5812 * 1024)
			ret := GetHumanSize(v)
			So(ret, ShouldEqual, "5.7M")
		})
		Convey("K", func() {
			var v = uint64(123 * 1024)
			ret := GetHumanSize(v)
			So(ret, ShouldEqual, "123.0K")
		})
		Convey("B", func() {
			var v = uint64(786)
			ret := GetHumanSize(v)
			So(ret, ShouldEqual, "786B")
		})
		Convey("0", func() {
			var v = uint64(0)
			ret := GetHumanSize(v)
			So(ret, ShouldEqual, "0B")
		})
		Convey("initest", func() {
			k := &Config{" ", "/tmp/ts.ini"}
			k.SaveDefaultUnpackerArgs(65, 65)
			k.SaveFileSlice(true)
			cfg, _ := ini.LooseLoad(self.cfgfile)

			val := cfg.Section("pkg").Key("file_slice").String()
			val2 := cfg.Section("default_unpacker").Key("max_seg_size").String()
			val3 := cfg.Section("default_unpacker").Key("max_seg_num").String()
			So(val, ShouldEqual, "true")
			So(val2, ShouldEqual, "608")
			So(val3, ShouldEqual, "65")
		})

	})
}
