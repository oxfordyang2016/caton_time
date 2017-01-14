package utils

import (
	"github.com/go-ini/ini"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_Config(t *testing.T) {
	Convey("Test Set", t, func() {

		Convey("initest", func() {
			k := &Config{" ", "/tmp/ts.ini"} //note:you can set ts.ini 's location
			k.SaveDefaultUnpackerArgs(65, 65)
			k.SaveFileSlice(true)
			cfg, _ := ini.LooseLoad(k.cfgfile)

			val := cfg.Section("pkg").Key("file_slice").String()
			val2 := cfg.Section("default_unpacker").Key("max_seg_size").String()
			val3 := cfg.Section("default_unpacker").Key("max_seg_num").String()
			So(val, ShouldEqual, "true")
			So(val2, ShouldEqual, "65")
			So(val3, ShouldEqual, "65")
		})

	})
}
