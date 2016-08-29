package routers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"quickstart/controllers"
	//"strings"
)

func init() {
	beego.Router("/api/v1/login", &controllers.Controller2{})
	beego.Router("/hello", &controllers.Controller1{})
	beego.Router("/login", &controllers.MainController{})
	beego.Router("/:uid/pkg", &controllers.Controller1{})

	beego.InsertFilter("/*", beego.FinishRouter, FilterUser, false)
}

var FilterUser = func(ctx *context.Context) {

	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++")
	fmt.Println(ctx.Output.IsEmpty())
	//fmt.Println(ctx.Output.Header(key, val))
	//fmt.Println(ctx.Output.JSON(data, hasIndent, coding))
}
