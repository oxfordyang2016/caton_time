package main

import (
	"fmt"
	//"io/ioutil"
	"os/exec"
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
	Portmanage("D", "INPUT", "tcp", "ACCEPT", "8000")
}
func Portmanage(action string, chain string, proc string, rule string, port string) {

	ru := "iptables  " + "-" + action + " " + chain + " -p " + proc + " --dport " + port + " -j " + rule
	iptable := exec.Command("bash", "-c", ru)
	out, _ := iptable.Output()
	fmt.Println(out)
}
