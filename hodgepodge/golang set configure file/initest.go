package main

import (
	"fmt"

	"github.com/go-ini/ini"
)

func main() {

	//---------------load ini.file--------------
	cfg, err := ini.InsensitiveLoad("/tmp/initest.ini")
	section, err := cfg.GetSection("yangming")
	fmt.Println("i love golang")

	//----------get section and key---------------
	sections := cfg.Sections()
	names := cfg.SectionStrings()
	keynames := cfg.Section("yangming").KeyStrings()
	val := cfg.Section("yangming").Key("i").String()
	fmt.Println("get a string key _value   ", val)
	//err2 := cfg.Section("").NewKey("qa", "vae")

	//-------------get key ,chanr key's value-------
	key1, err1 := cfg.Section("yangming").GetKey("i")
	fmt.Println("get a key    ", key1)
	key1.SetValue("i have change key's value")
	fmt.Println("get a  change key    ", key1)

	//-------save change to a ini.file-0--------
	err = cfg.SaveTo("/tmp/initest.ini")

	fmt.Println("get cfg and err ", cfg, err, err1)
	fmt.Println("i will get section  ", section, " get section names ", names, " will gwt key name", keynames, "  will get sections", sections)

}
