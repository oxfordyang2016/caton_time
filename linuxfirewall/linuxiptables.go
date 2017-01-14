package main

import (
	"fmt"
	//"io/ioutil"
	"os/exec"
)

func main() {

	ls := exec.Command("bash", "-c", "ls -a -l -h ")
	lsout, _ := ls.Output()
	fmt.Println(string(lsout))

	iptable := exec.Command("bash", "-c", "iptables -I INPUT -p tcp --dport 8000 -j DROP")
	ipout, _ := iptable.Output()
	fmt.Print(ipout)

}
