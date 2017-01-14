package main

import (
	"fmt"
	//"io/ioutil"
	"os/exec"
	"strconv"
)

func main() {
	/*
	   ls := exec.Command("bash", "-c", "ls -a -l -h ")
	   lsout, _ := ls.Output()
	   fmt.Println(string(lsout))

	   iptable := exec.Command("bash", "-c", "iptables -I INPUT -p tcp --dport 8000 -j DROP")
	   ipout, _ := iptable.Output()
	   fmt.Print(ipout)
	*/
	Portmanage(true, "INPUT", "udp", 1200)
}

func Portmanage(action bool, chain string, proc string, port int) {
	var ru string
	if action == true {
		ru = "iptables  " + "-" + "I" + " " + chain + " -p " + proc + " --dport " + strconv.Itoa(port) + " -j " + "ACCEPT"
	}
	if action == false {
		ru = "iptables  " + "-" + "D" + " " + chain + " -p " + proc + " --dport " + strconv.Itoa(port) + " -j " + "ACCEPT"
	}
	iptable := exec.Command("bash", "-c", ru)
	out, _ := iptable.Output()
	fmt.Println(out)
}
