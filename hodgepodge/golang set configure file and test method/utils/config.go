package utils

import (
	"fmt"
	"github.com/go-ini/ini"
	"strconv"
	"strings"
	"syscall"
)

const (
	CFG_SEP = ","
)

var (
	SizeStrPostfix = []string{"k", "m", "g"}
	SizeUnits      = []uint64{1024, 1024 * 1024, 1024 * 1024 * 1024}
	config         *Config
)

type Config struct {
	profile string
	cfgfile string
}

func NewConfig(profile, cfgfile string) *Config {
	o := &Config{
		profile, cfgfile,
	}
	return o
}

func (self *Config) SaveFileSlice(Switch bool) error {
	cfg, err := ini.LooseLoad(self.cfgfile)
	if err != nil {
		return err
	}
	sec := cfg.Section("pkg")

	_, err = sec.NewKey("file_slice", fmt.Sprint(Switch))
	if err != nil {
		return err
	}
	return cfg.SaveTo(self.cfgfile)
}

func (self *Config) SaveDefaultUnpackerArgs(maxofslice uint64, numberofslice uint64) error {
	cfg, err := ini.LooseLoad(self.cfgfile)
	if err != nil {
		return err
	}

	sec := cfg.Section("default_unpacker")

	_, err5 := sec.NewKey("max_seg_size", fmt.Sprint(maxofslice))
	if err5 != nil {
		return err5
	}
	_, err6 := sec.NewKey("max_seg_num", fmt.Sprint(numberofslice))

	if err6 != nil {
		return err6
	}
	return cfg.SaveTo(self.cfgfile)
}

func (self *Config) Load() (*ini.File, error) {
	if err := syscall.Access(self.profile, syscall.F_OK); err != nil {
		return nil, fmt.Errorf("%s: %s", err.Error(), self.profile)
	}
	return ini.LooseLoad(self.profile, self.cfgfile)
}

// sizestr, 例如100, 100K/k, 100M/m, 100G/g等
func ConfigGetSize(sizestr string) uint64 {
	if sizestr == "" {
		return 0
	}

	sep_idx := -1
	lower_str := strings.ToLower(sizestr)
	for i, sep := range SizeStrPostfix {
		r := strings.Index(lower_str, sep)
		if r != -1 {
			// 字母必须是最后
			if r != len(lower_str)-1 {
				return 0
			}
			lower_str = lower_str[:r]
			sep_idx = i
			break
		}
	}
	n, err := strconv.ParseUint(lower_str, 10, 64)
	if err != nil {
		return 0
	}
	if sep_idx >= 0 {
		n = n * SizeUnits[sep_idx]
	}
	return n
}

func MakeDefaultConfig(cfg *Config) {
	if cfg != nil {
		config = cfg
	}
}

func DefaultConfig() *Config {
	return config
}
