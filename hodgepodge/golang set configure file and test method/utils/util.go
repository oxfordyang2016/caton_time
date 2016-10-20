package utils

import (
	"fmt"
)

type SizeSymbol struct {
	Size uint64
	Name string
}

var (
	SizeSymbols = []SizeSymbol{
		{uint64(1024 * 1024 * 1024 * 1024), "T"},
		{uint64(1024 * 1024 * 1024), "G"},
		{uint64(1024 * 1024), "M"},
		{uint64(1024), "K"},
	}
)

func GetHumanSize(v uint64) string {
	var s *SizeSymbol
	for _, ss := range SizeSymbols {
		if v >= ss.Size {
			s = &ss
			break
		}
	}
	if s != nil {
		ret := float64(v) / float64(s.Size)
		return fmt.Sprintf("%.1f%s", ret, s.Name)
	}
	return fmt.Sprintf("%dB", v)
}
