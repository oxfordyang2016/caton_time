package main

/*
ABOUT HTTP REQUEST

ABOUT URL AND REQUEST
          |
          |
          v
i can let argvs(or info  that i wnat send ) into url or body or header
note:as client,for example ,i get a page (jd login page https://passport.jd.com/new/login.aspx),it will
guide data's flow_direction
          |
          |
          |
          v
a.urls ---->http://192.168.0.68:8080/5643/pkg?query=86mskksjkdsfjdsfkkse&yang=65
b.header---->but get_method's body is empty
c.body-----
          |
          |
          v
          uuid:5e77eec3-ea37-4c12-9ed0-798ecbdf0b1b
machineNet:
machineCpu:
machineDisk:
eid:4B74E9DBE609F53DEC6F6241B9883F85AAC1F10B9764168D600402A280C3934DBE02B216AE4F01F57512AE9018256459
fp:91295ba049b53e3a4ccf068abcfe1ad2
_t:_ntBwnKD
loginType:f
USqADyILTD:OoQWd
loginname:15099667237
nloginpwd:reterdfsd
loginpwd:reterdfsd
chkRememberMe:on
authcode:SFDFSD


*/

import (
	"github.com/astaxie/beego"
	_ "quickstart/routers"
)

func main() {
	//fmt.Println("hjskjkl")
	beego.Run()

}
