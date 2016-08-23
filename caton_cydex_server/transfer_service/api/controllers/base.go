package controllers

import (
	"github.com/astaxie/beego"
	"strconv"
)

type BaseController struct {
	beego.Controller
	UserLevel int
}

func (self *BaseController) fetchUserLevel() {
	ul := self.Ctx.Input.Header("x-cydex-userlevel")
	if ul != "" {
		i, err := strconv.Atoi(ul)
		if err == nil {
			self.UserLevel = i
		}
	}
}

func (self *BaseController) Prepare() {
	self.fetchUserLevel()
}
