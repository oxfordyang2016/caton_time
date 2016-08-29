package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}
type Controller1 struct {
	beego.Controller //in beego dir ,some file contain Controllers
}

/*type Controller1 struct {
	beego.Controller //in beego dir ,some file contain Controllers
}
*/
func (c *MainController) Get() {
	// c.Data["Website"] = "beego.me"
	// c.Data["Email"] = "astaxie@gmail.com"
	// c.TplName = "index.tpl"
	// fmt.Println(c.Data["Email"])
	// c.Data["Json"]="======================================>"
	// fmt.Println(c.Data["Json"])
	query := c.GetString("query")
	fmt.Println(query)
	c.Ctx.WriteString("hello,world")
}
func (c *Controller1) Get() {
	// c.Data["Website"] = "beego.me"
	// c.Data["Email"] = "astaxie@gmail.com"
	// c.TplName = "index.tpl"
	// fmt.Println(c.Data["Email"])
	// c.Data["Json"]="======================================>"
	// fmt.Println(c.Data["Json"])
	query := c.GetString("query")
	query1, _ := c.GetInt("yang")
	//fmt.Println(get)
	fmt.Println(query)
	fmt.Println(query1)
	c.Ctx.WriteString(query)
}
