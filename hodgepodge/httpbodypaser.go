package main

/*


http client  , in body write all kinds of infomation
      |
      |
      |
      v
"dajkjadkkalda"
dskklsflk
    |
    |
    |


*/
import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	//"os"
	"crypto/md5"
	"reflect"
	"strings"
	//"time"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	//"github.com/ziutek/mymysql/mysql"
	//_ "github.com/ziutek/mymysql/native" // Native engine
)

/* about http's header
1.
about http header,header's filed is string ,even though,it is
the type: id:643888348
2.
in server ,resolving http header to all kinds of type is another thing

*/

type test_struct struct {
	Test string
}
type Response struct {
	Token  string
	Status int
}

var (
	first_token string
)

//-----------------------------------------------------------controller------------------------------------------
func report(rw http.ResponseWriter, req *http.Request) {
	var authorization string

	for name, headers := range req.Header {
		name := strings.ToLower(name)
		for _, h := range headers {
			//fmt.Println("name++++", name)

			if name == "authorization" {
				authorization = h
			}

		}
	}
	//------------------------------split string----------
	//http://www.dotnetperls.com/split-go
	//result := strings.Split(authorization, " ")

	// Display all elements.
	// for i := range result {
	// 	fmt.Println(result[i])
	// }
	//split := result[1]
	final_token := "Bearer" + " " + first_token
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++=", final_token)
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++", authorization)

	if authorization != final_token {
		io.WriteString(rw, "jwt eror")
		return
	}

	io.WriteString(rw, "report ok") //write to body
	fmt.Println(formatRequest(req))
	body, _ := ioutil.ReadAll(req.Body)

	//header, _ := ioutil.ReadAll(req.Header)
	fmt.Println("body's type", reflect.TypeOf(body))
	fmt.Println(reflect.TypeOf(req.Body))
	fmt.Println(reflect.TypeOf(string(body)))
	fmt.Println(req.Body)
	fmt.Println("================================++++++++++++++++++++++++++++++++======================================")
	fmt.Println("request body bynary=>", body) //if only here,it will occur  byte
	fmt.Println(string(body))
	//fmt.Println(string(header))
	fmt.Println(formatRequest(req))
	//---------------------database access------------------------------
	var content_type string
	for name, headers := range req.Header {
		name := strings.ToLower(name)
		for _, h := range headers {
			//fmt.Println("name++++", name)
			if name == "content-type" {
				content_type = h
			}

			//fmt.Println("content+++", h)
			//request := append(request, fmt.Sprintf(name, h))
		}
	}

	db, err := sql.Open("mysql", "root:123456@tcp(192.168.0.81:3306)/hello") //first configure a database
	checkErr(err)

	// insert
	stmt, err := db.Prepare("INSERT report SET sn=?,content_type=?,msg_body=?")
	checkErr(err)
	//ko := "string"
	res, err := stmt.Exec("yang", content_type, body)
	checkErr(err)

	id, err := res.LastInsertId()
	fmt.Println(id)
	checkErr(err)
	db.Close()
}
func test2(rw http.ResponseWriter, req *http.Request) {
	io.WriteString(rw, "hallo") //write to body
}
func test(rw http.ResponseWriter, req *http.Request) {
	//io.WriteString(rw, "hallo") //write to body
	//data := []byte("hello")
	//a := fmt.Sprintf("%x", md5.Sum(data))
	//about passwd's computing method:password=MD5(MD5(sn):secrect:date)
	//fmt.Println(reflect.TypeOf(string(a)))
	//k := a + "caton" + date
	//------------------------------------------database connettion----------------------------

	db, err := sql.Open("mysql", "root:123456@tcp(192.168.0.81:3306)/hello") //first configure a database

	checkErr(err)

	// insert
	stmt, err := db.Prepare("INSERT who SET sn=?,model=?,version=?") //in database create a table called who
	/*

		CREATE TABLE potluck (id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
		name VARCHAR(20),
		food VARCHAR(30),
		confirmed CHAR(1),
		signup_date DATE);

	*/
	//retrive info from header
	var sn, model, version, date, password string

	for name, headers := range req.Header {
		name := strings.ToLower(name)
		for _, h := range headers {
			//fmt.Println("name++++", name)

			if name == "sn" {
				sn = h
			}
			if name == "date" {
				date = h
			}
			if name == "model" {
				model = h
			}
			if name == "version" {
				version = h
			}
			if name == "password" {
				password = h
			}
			//fmt.Println("content+++", h)
			//request := append(request, fmt.Sprintf(name, h))
			fmt.Println("name++++", name)
			fmt.Println("content+++", h)
		}
	}
	//----------md5 verify login-----------------------
	snk := []byte(sn)
	md51 := fmt.Sprintf("%x", md5.Sum(snk)) //lowcase
	//fmt.Print("snk........", snk, "     after computer.....", md51)
	//about passwd's computing method:password=MD5(MD5(sn):secrect:date)
	//fmt.Println(reflect.TypeOf(string(a)))
	k := md51 + "caton" + date
	data1 := []byte(k)
	verify_data := fmt.Sprintf("%x", md5.Sum(data1))
	if verify_data == password {
		fmt.Println("string ok--")

		//return
	}
	if verify_data != password {
		fmt.Println("string no--")
		io.WriteString(rw, "password error")
		return
	}
	fmt.Println(reflect.TypeOf(verify_data))
	fmt.Println(reflect.TypeOf(password))
	checkErr(err)
	//-----------------------write who to databases table who----------------------------------------
	res, err := stmt.Exec(sn, model, version)
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Println(id)
	/*
		// update
		stmt, err = db.Prepare("update userinfo set username=? where uid=?")
		checkErr(err)

		res, err = stmt.Exec("astaxieupdate", id)
		checkErr(err)

		affect, err := res.RowsAffected()
		checkErr(err)

		fmt.Println(affect)

		// query
		rows, err := db.Query("SELECT * FROM userinfo")
		checkErr(err)

		for rows.Next() {
			var uid int
			var username string
			var department string
			var created string
			err = rows.Scan(&uid, &username, &department, &created)
			checkErr(err)
			fmt.Println(uid)
			fmt.Println(username)
			fmt.Println(department)
			fmt.Println(created)

		}
	*/
	// // delete
	// stmt, err = db.Prepare("delete from userinfo where uid=?")
	// checkErr(err)

	// res, err = stmt.Exec(id)
	// checkErr(err)

	// affect, err = res.RowsAffected()
	// checkErr(err)

	// fmt.Println(affect)

	// db.Close()

	//-------------------------read http body----------------------------------------
	body, _ := ioutil.ReadAll(req.Body)
	//header, _ := ioutil.ReadAll(req.Header)
	fmt.Println(reflect.TypeOf(body))
	fmt.Println(reflect.TypeOf(req.Body))
	fmt.Println(reflect.TypeOf(string(body)))
	fmt.Println(req.Body)
	fmt.Println(body) //if only here,it will occur  byte
	fmt.Println(string(body))
	//fmt.Println(string(header))
	fmt.Println(formatRequest(req))
	//--------------http response------------------------------------->
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	mySigningKey := "caton"
	to, _ := ExampleNew([]byte(mySigningKey))
	fmt.Println(to)
	first_token = to
	fmt.Println(reflect.TypeOf(to))
	rep := &Response{Token: to, Status: 1}
	b, _ := json.Marshal(rep)
	io.WriteString(rw, string(b)) //write to body
	fmt.Println(b)
	//--------------------http server --------------------------
	/*http server

	package main

	import (
		"io"
		"net/http"
	)

	func hello(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello world!")
	}

	func main() {
		http.HandleFunc("/", hello)
		http.ListenAndServe(":8000", nil)
	}

	*/
	//----------------------------------------------struct to json-----------------------------
	/*
	   package main

	   import (
	       "fmt"
	       "encoding/json"
	   )

	   type User struct {
	       Name string
	   }

	   func main() {
	       user := &User{Name: "Frank"}
	       b, err := json.Marshal(user)
	       if err != nil {
	           fmt.Println(err)
	           return
	       }
	       fmt.Println(string(b))
	   }
	*/

	//output:{"Name":"Frank"}
}

func main() {
	http.HandleFunc("/api/v1/login", test)
	http.HandleFunc("/api/v1/test", test2)
	http.HandleFunc("/api/v1/report", report)
	log.Fatal(http.ListenAndServe(":8083", nil))

}

//------formatRequest generates ascii representation of a request--------------read header--------------------------------
func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("method==>", r.Method, "url===>", r.URL, "r.proto==>", r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf(r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name := strings.ToLower(name)
		for _, h := range headers {
			fmt.Println("name++++", name)
			fmt.Println("content+++", h)
			request = append(request, fmt.Sprintf(name, h))
		}
	}

	// // If this is a POST, add post data
	// if r.Method == "POST" {
	// 	r.ParseForm()
	// 	request = append(request, "\n")
	// 	request = append(request, r.Form.Encode())
	// }
	// // Return the request as a string
	// return strings.Join(request, "\n")

	return "header paser ok"
}

//------------------------jwt token create------------- create jwt token-------------------
func ExampleNew(mySigningKey []byte) (string, error) {
	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	//token.Claims["foo"] = "bar"
	//token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(mySigningKey)
	return tokenString, err
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
