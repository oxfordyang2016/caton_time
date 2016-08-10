package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"net/http"
	//	"strconv"
	"reflect" //get some type of var
	"time"
)

func main() {
	fmt.Println("i love server")
	e := echo.New()
	fmt.Println(reflect.TypeOf("i love"))
	//------get VAR TYPE---------
	fmt.Println("get echo.new() type", reflect.TypeOf(e))

	//----------deal with get method----------------------------------------->
	//when you run this .go ,you input 192.168.0.58:1323/user,you will get time
	e.GET("/", func(c echo.Context) error { //this line is get info from client
		daytime := time.Now().String()
		fmt.Println("get echo.Context() type", reflect.TypeOf(c))
		return c.String(http.StatusOK, "Hello, World, i am yangming !"+daytime)
	})
	e.GET("/user", func(c echo.Context) error { //this line is get info from client
		daytime := time.Now().String()
		return c.String(http.StatusOK, "Hello, World, i am yangming !"+daytime)
	})
	//this callback function
	e.GET("/users/:id", getUser)

	//-------------deal with post method-------------------------------------------------------------->
	e.GET("/show?team=x-men&member=wolverine", show)

	e.Run(standard.New(":1323"))

}

//--------there ,you can storage data to database----------------------->
/*there,you can design call back fucntion to meet all kinds of needs
 */
func getUser(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("id")
	daytime := time.Now().String()
	fmt.Println("request  add==>", c.Request().(*standard.Request).Request)
	fmt.Println("request url===>", c.Request().URL().(*standard.URL).URL)
	fmt.Println("req header>", c.Request().Header().(*standard.Header).Header)
	return c.String(http.StatusOK, "Hello, World, i am yangming !"+daytime+"your is is"+id)
}

func show(c echo.Context) error {
	// Get team and member from the query string
	team := c.QueryParam("team")
	member := c.QueryParam("member")
	return c.String(http.StatusOK, member+"hallo yanmging"+team)
}
