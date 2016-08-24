package api

//http://aslanbakan.com/en/blog/33-essential-sublime-text-plugins-for-all-developers/
//it is likely http server api
import (
	c "./controllers"
	"github.com/astaxie/beego"
)

/*
about controller
beego.Router("/hello", &controllers.Controller1{})
          |
          v
          v
          v
          v
type MainController struct {
	beego.Controller
}
          |
          v
          v
          v
          v
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
	//url========>/3627327/pkg?query==547543 will be captured (according to mrs ran)
	//url=====>/https://www.google.com/search?q=shanghai&oq=shanghai&aqs=chrome..69i57j69i60l2j69i65l2j69i59.1549j0j7&sourceid=chrome&ie=UTF-8
	//	GET /tickets?fields=id,subject,customer_name,updated_at&state=open&sort=-updated_at
	//http://vinaysahni.us6.list-manage.com/subscribe/post?u=005ede9b44df2622ce536ab88&id=8edc074683

	/*
	                         |
	                         |
	                         |
	                         v
	                         v
	                         v
	   func init() {

	   	beego.Router("/hello", &controllers.Controller1{})
	   	beego.Router("/login", &controllers.MainController{})
	   	beego.Router("/:uid/pkg", &controllers.Controller1{})

	   	beego.InsertFilter("/*", beego.FinishRouter, FilterUser, false)
	   }
	                     |
	                     |
	                     |
	                     v
	                     v
	   func (c *Controller1) Get() {
	   	// c.Data["Website"] = "beego.me"
	   	// c.Data["Email"] = "astaxie@gmail.com"
	   	// c.TplName = "index.tpl"
	   	// fmt.Println(c.Data["Email"])
	   	// c.Data["Json"]="======================================>"
	   	// fmt.Println(c.Data["Json"])
	   	query := c.GetString("query")
	   	query1 := c.GetString("yang")
	   	//fmt.Println(get)
	   	fmt.Println(query)
	   	fmt.Println(query1)
	   	c.Ctx.WriteString(query)
	   }

	*/
	/*
		                                |
		                                |
		                                |
		                                v
		                                v
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

		                |
		                |----------->switch
		                v
		                v

		package main
		import "fmt"
		import "time"
		func main() {
		                    i := 2
		    fmt.Print("write ", i, " as ")
		    switch i {
		    case 1:
		        fmt.Println("one")
		    case 2:
		        fmt.Println("two")
		    case 3:
		        fmt.Println("three")
		    }
		  switch time.Now().Weekday() {
		    case time.Saturday, time.Sunday:
		        fmt.Println("it's the weekend")
		    default:
		        fmt.Println("it's a weekday")
		    }
		     t := time.Now()
		    switch {
		    case t.Hour() < 12:
		        fmt.Println("it's before noon")
		    default:
		        fmt.Println("it's after noon")
		    }
		}

		$ go run switch.go
		write 2 as two
		it's the weekend
		it's before noon

	*/
	//use regular exp http://code.tutsplus.com/tutorials/8-regular-expressions-you-should-know--net-6149

	/*
					   In order to make the router settings easier, Beego references the router implementation approach found in Sinatra. It supports many router types.

					   beego.Router(“/api/?:id”, &controllers.RController{})
					   default matching /api/123 :id = 123 can match /api/

					   beego.Router(“/api/:id”, &controllers.RController{})
					   default matching /api/123 :id = 123 can’t match /api/

					   beego.Router(“/api/:id([0-9]+)“, &controllers.RController{})
					   Customized regex matching /api/123 :id = 123

					   beego.Router(“/user/:username([\w]+)“, &controllers.RController{})
					   Regex string matching /user/astaxie :username = astaxie

					   beego.Router(“/download/*.*”, &controllers.RController{})
					   matching /download/file/api.xml :path= file/api :ext=xml

					   beego.Router(“/download/ceshi/*“, &controllers.RController{})
					   full matching /download/ceshi/file/api.json :splat=file/api.json

					   beego.Router(“/:id:int”, &controllers.RController{})
					   int type matching :id is int type. Beego implements ([0-9]+) for you

					   beego.Router(“/:hello:string”, &controllers.RController{})
					   string type matching :hello is string type. Beego implements ([\w]+) for you

					   beego.Router(“/cms_:id([0-9]+).html”, &controllers.CmsController{})
					   has prefix regex :id is the regex. matching cms_123.html :id = 123

					   In controller, you can get the variables like this:
					   this.Ctx.Input.Param(":id")
		        		this.Ctx.Input.Param(":username")
			        	this.Ctx.Input.Param(":splat")
				       this.Ctx.Input.Param(":path")
				      this.Ctx.Input.Param(":ext")

	*/

	beego.Router("/:uid/pkg/:pid", &c.PkgController{}) //this client 's entrance
	beego.Router("/:uid/pkg/:pid/:fid", &c.FileController{})
	beego.Router("/:uid/transfer", &c.TransferController{})
}

//now time, i have a doubt that what json does????
//because all argv pas to api ,even if a priod of messages
//about json ,there is something that when  all argvs was passed to url,and then ,build conenction
//with server ,server deal with request ,and send json message to client
func init() {
	setup()
}
