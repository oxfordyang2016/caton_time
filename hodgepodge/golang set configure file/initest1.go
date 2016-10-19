package main

import (
	"fmt"

	"github.com/go-ini/ini"
)

func main() {

	//---------------load ini.file--------------
	cfg, err := ini.InsensitiveLoad("/tmp/initest.ini")

	//----------------get a section---------------------
	section, err := cfg.GetSection("yangming") //get a section that does not create a section
	section1 := cfg.Section("yangmingming")

	fmt.Println("i love golang", section1)

	//----------get section and key---------------
	sections := cfg.Sections()
	names := cfg.SectionStrings()
	keynames := cfg.Section("yangming").KeyStrings()
	val := cfg.Section("yangming").Key("i").String()
	fmt.Println("-----------------------------i will get key exp---------------------")
	fmt.Println("get a string key _value   ", val)
	fmt.Println("-----------------------------i will get key exp---------------------")
	//err2 := cfg.Section("").NewKey("qa", "vae")

	//-------------get key ,change key's value-------
	key1, err1 := cfg.Section("yangming").GetKey("i")
	fmt.Println("get a key    ", key1)
	key1.SetValue("i have change key's value")
	fmt.Println("get a  change key    ", key1)

	//--------------get a not exist key tp test if it will be create

	key2, _ := cfg.GetSection("yangming3")
	key3, _ := cfg.Section("yangming").GetKey("key name")
	fmt.Println("test key2 key3   ", key2, key3)

	//--------------delte a section---------
	cfg.DeleteSection("yangming")

	//----------get section and delete a key and new a key ------
	section2, err := cfg.GetSection("yangming4")
	section2.DeleteKey("dir")
	section2.NewKey("yangming", "i wan to say")

	//-----------test new a exisit key if chnage a key's value-----
	section2.NewKey("yangming", "i have no to say") //yes it will change origin key's value

	val6 := cfg.Section("yangming4").Key("yangming").String()
	val8 := cfg.Section("yangming4").Key("yangming2").String()
	fmt.Println("two test%%%%%%%%%%%%%%%%%%%%%   ", val6, "******donot exist key********** ", val8)

	//------------detect a key if exist--------
	exist := section2.HasKey("yang1ming")
	fmt.Println("detect exist of yangmingkey ,result is --->", exist)

	//---------------new a section-----------
	cfg.NewSection("yangming3")

	//--------note:you can reload a file----------

	//-------save change to a ini.file-0--------
	err = cfg.SaveTo("/tmp/initest.ini")

	/*
	   i have to say
	   1.cfg.Section will create  a section if it does not exist
	   2.section.NewKey ,it will a key if key doesnot exist ,it will update key's value if key  exists



	*/

	fmt.Println("get cfg and err ", cfg, err, err1)
	fmt.Println("i will get section  ", section, " get section names ", names, " will gwt key name", keynames, "  will get sections", sections)
	version := ini.Version()
	fmt.Println("i want to get  ini versison ", version)

}
