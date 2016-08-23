package api

//it is likely http server api
import (
	c "./controllers"
	"github.com/astaxie/beego"
)

/*
about controller
beego.Router("/hello", &controllers.Controller1{})
          |
          |
type MainController struct {
	beego.Controller
}
          |
          |
func (c *MainController) Get() {
	// c.Data["Website"] = "beego.me"
	// c.Data["Email"] = "astaxie@gmail.com"
	// c.TplName = "index.tpl"
	// fmt.Println(c.Data["Email"])
	// c.Data["Json"]="======================================>"
	// fmt.Println(c.Data["Json"])

	c.Ctx.WriteString("hello,world")//write data to connect
}

*/

func setup() {
	// TODO 正则限制长度
	beego.Router("/:uid/pkg", &c.PkgsController{})
	/*
	   type PkgsController struct {
	   	BaseController
	   }

	   func (self *PkgsController) Get() {
	   	query := self.GetString("query")
	   	filter := self.GetString("filter")
	   	list := self.GetString("list")

	   	// sec 3.2 in api doc
	   	if list == "list" {
	   		self.getLitePkgs()
	   		return
	   	}

	   	switch query {
	   	case "all":
	   		if filter == "change" {
	   			// sec 3.1.2 in api doc
	   			self.getActive()
	   		}
	   	case "sender":
	   		// 3.1
	   		self.getJobs(cydex.UPLOAD)
	   	case "receiver":
	   		// 3.1
	   		self.getJobs(cydex.DOWNLOAD)
	   	case "admin":
	   		// 3.1
	   		self.getAllJobs()
	   	}
	   }

	*/

	beego.Router("/:uid/pkg/:pid", &c.PkgController{})
	beego.Router("/:uid/pkg/:pid/:fid", &c.FileController{})
	beego.Router("/:uid/transfer", &c.TransferController{})
}

func init() {
	setup()
}
