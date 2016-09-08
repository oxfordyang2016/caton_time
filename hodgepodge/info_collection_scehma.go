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
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"strings"
	"time"
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
	version_g   string
	db          *sql.DB
	err         error
	sn_g        string
)

const (
	mySigningKey = "caton"
)

func main() {
	db, err = sql.Open("mysql", "root:123456@tcp(192.168.0.81:3306)/hello") //first configure a database
	checkErr(err)
	http.HandleFunc("/api/v1/login", login)
	http.HandleFunc("/api/v1/test", test2)
	http.HandleFunc("/api/v1/report", report)
	log.Fatal(http.ListenAndServe(":8082", nil))

}
func test2(rw http.ResponseWriter, req *http.Request) {
	io.WriteString(rw, "hallo") //write to body
}
func login(rw http.ResponseWriter, req *http.Request) {

	//-----------------------------------paser request body json-------------
	/*
	   	package main

	   import (
	   	"encoding/json"
	   	"fmt"
	   	"net/http"
	   )

	   type test_struct struct {
	   	Test string
	   }

	   func parseGhPost(rw http.ResponseWriter, request *http.Request) {
	   	decoder := json.NewDecoder(request.Body)

	   	var t test_struct
	   	err := decoder.Decode(&t)

	   	if err != nil {
	   		panic(err)
	   	}

	   	fmt.Println(t.Test)
	   }

	   func main() {
	   	http.HandleFunc("/", parseGhPost)
	   	http.ListenAndServe(":8080", nil)
	   }
	*/

	//-----------json  to struct---------------------------------------------------------
	/*
		        text := []byte(`{"from":"a", "cmd": "register", "type": "req",
				"seq": 11,"task_id":"67","uid":"6732","token": "EF02JLGFA09GVNG21F","req":{"register":{"machine_code":"mac"}}}`)

					var msg Message
					err := json.Unmarshal(text, &msg)
		             	So(msg.From, ShouldEqual, "a")

		}
	*/

	//------------------------paser http request body's json-------------------------------------
	decoder := json.NewDecoder(req.Body)

	type Info struct {
		Sn      string
		Model   string
		Version string
	}
	var t Info
	err1 := decoder.Decode(&t)
	fmt.Println(err1)
	fmt.Println("+json paser success-> ", t.Sn)

	//-----------------------------------retrive info from header------------------------------
	var sn, model, version, date, password string
	for name, headers := range req.Header {
		name := strings.ToLower(name)
		for _, h := range headers {
			//fmt.Println("name++++", name)

			if name == "sn" {
				sn = h
				sn_g = h
			}
			if name == "date" {
				date = h
			}
			if name == "model" {
				model = h
			}
			if name == "version" {
				version = h
				version_g = h
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
	fmt.Println(sn, version, date, model, version)

	//-------------------------------------md5 verify login--------------------------------------------------------
	snk := []byte(t.Sn)                     //origin sn from header,but changed ,not test!
	md51 := fmt.Sprintf("%x", md5.Sum(snk)) //lowcase
	//fmt.Print("snk........", snk, "     after computer.....", md51)
	//about passwd's computing method:password=MD5(MD5(sn):secrect:date)
	//fmt.Println(reflect.TypeOf(string(a)))
	//k := md51 + "caton" + date
	k := md51 + mySigningKey + date
	data1 := []byte(k)
	verify_data := fmt.Sprintf("%x", md5.Sum(data1))
	if verify_data == password {
		fmt.Println("login ok--")

		//return
	}
	if verify_data != password {
		fmt.Println("login  error--")
		//io.WriteString(rw, "password error")
		//io.WriteString(rw, string(b1))
		io.WriteString(rw, "login error")
		return
	}
	//-------------------------read http body----------------------------------------
	body, _ := ioutil.ReadAll(req.Body)
	//header, _ := ioutil.ReadAll(req.Header)
	fmt.Println(reflect.TypeOf(body))
	fmt.Println(reflect.TypeOf(req.Body))
	fmt.Println(reflect.TypeOf(string(body)))
	fmt.Println(req.Body)
	fmt.Println(string(body))

	fmt.Println(body)
	fmt.Println(formatRequest(req))

	//--------------return jwt in http response of login_ok------------------------------------->

	mySigningKey := "caton"
	to, _ := ExampleNew([]byte(mySigningKey)) //give jwt token
	fmt.Println(to)
	first_token = to
	fmt.Println(reflect.TypeOf(to))
	rep := &Response{Token: to, Status: 1}
	b, _ := json.Marshal(rep)

	io.WriteString(rw, string(b)) //write to body

	fmt.Println(b)

	//-----------------------write who to databases table who----------------------------------------
	/*

		CREATE TABLE potluck (id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
		name VARCHAR(20),
		food VARCHAR(30),
		confirmed CHAR(1),
		signup_date DATE);

	*/

	/*

	   //Query
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
	/*

	   age := 27
	   rows, err := db.Query("SELECT name FROM users WHERE age=?", age)
	   if err != nil {
	           log.Fatal(err)
	   }
	   defer rows.Close()
	   for rows.Next() {
	           var name string
	           if err := rows.Scan(&name); err != nil {
	                   log.Fatal(err)
	           }
	           fmt.Printf("%s is %d\n", name, age)
	   }
	   if err := rows.Err(); err != nil {
	           log.Fatal(err)
	   }
	*/
	//query

	rows, err_tmp := db.Query("SELECT * FROM who WHERE sn=?", t.Sn)
	fmt.Println("-------------------------------------", rows)
	checkErr(err)
	fmt.Println("------------------------------------------------------query database---------------------")
	for rows.Next() {
		var id1 string
		var sn string
		var model string
		var version string
		var last_login string
		err = rows.Scan(&id1, &sn, &model, &version, &last_login)
		checkErr(err)

		fmt.Println(sn)

		if sn == t.Sn {
			fmt.Println("-----------------------------------into loop-------------------------------")
			//-------------------------------update-------------------------------
			stmt, err := db.Prepare("update who set last_login=? where sn=?")
			checkErr(err)
			last_login := time.Now()
			fmt.Println(last_login)
			res, err := stmt.Exec(last_login, t.Sn)
			checkErr(err)

			affect, err := res.RowsAffected()
			checkErr(err)

			fmt.Println(affect)
			fmt.Println("--------------------------------update data exp----------------------------")
			return
		}

	}

	//if sn is not in table ,inser it

	fmt.Println("query------------------------fail------------------>", err_tmp, t.Sn)
	stmt, err := db.Prepare("INSERT who SET sn=?,model=?,version=?,last_login=?") //in database create a table called who
	defer stmt.Close()
	checkErr(err)
	time_login := time.Now()
	res, err := stmt.Exec(t.Sn, t.Model, t.Version, time_login)
	checkErr(err)

	id, err := res.LastInsertId()
	fmt.Println(id)
	checkErr(err)

	/*
			if sn == t.Sn {
				//-------------------------------update-------------------------------
				stmt, err := db.Prepare("update who set last_login=? where sn=?")
				checkErr(err)
				last_login := time.Now()
				res, err := stmt.Exec(last_login, t.Sn)
				checkErr(err)

				affect, err := res.RowsAffected()
				checkErr(err)

				fmt.Println(affect)

			}

			fmt.Println(version)
			fmt.Println(model)

		}


	*/

	// -----------------------------insert-----------------------------
	/*
		stmt, err := db.Prepare("INSERT who SET sn=?,model=?,version=?,last_login=?") //in database create a table called who
		defer stmt.Close()
		checkErr(err)
		time_login := time.Now()
		res, err := stmt.Exec(t.Sn, t.Model, t.Version, time_login)
		checkErr(err)

		id, err := res.LastInsertId()
		fmt.Println(id)
		checkErr(err)
	*/
	//fmt.Println("-------------------------------------test success---------------", rows, err)

	//fmt.Println(id)

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

//-----------------------------------------------------------controller------------------------------------------
func report(rw http.ResponseWriter, req *http.Request) {
	var authorization, ip, content_length string
	time := time.Now()
	fmt.Println(time.Format("20060102150405"))
	ip = strings.Split(req.RemoteAddr, ":")[0] //get ip
	for name, headers := range req.Header {
		name := strings.ToLower(name)
		for _, h := range headers {
			//fmt.Println("name++++", name)

			if name == "authorization" {
				authorization = h

			}
			if name == "content-length" {
				content_length = h
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
	//------------------------------------------report verify---------string verify----------------------------------
	final_token := "Bearer" + " " + first_token
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++=", final_token)
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++", authorization)
	//----------second verify---------------
	fmt.Println(final_token, authorization)
	/*
		if authorization != final_token {
			io.WriteString(rw, "jwt eror+++++++++++++++++++++++++++++++++++++++++++++")
			return
		}

	*/
	//--------------------------------------report verify--------------jwt verify---------------------------------

	report_sn := ExampleParse(first_token, "caton")
	if report_sn == "error" {
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
	fmt.Println(content_type)

	//-----------------------------------------------write database-----------------------------------------------
	//https://github.com/go-sql-driver/mysql/wiki/Examples
	// insert
	stmt, err := db.Prepare("INSERT report SET sn=?,ip=?,time=?,version=?,content_type=?,content_length=?,msg_body=?")
	defer stmt.Close()
	checkErr(err)

	//ko := "string"
	res, err := stmt.Exec(report_sn, ip, time, version_g, content_type, content_length, body)
	checkErr(err)

	id, err := res.LastInsertId()
	fmt.Println(id)
	checkErr(err)
	//db.Close()

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

//--------------------------------global error handing mechnism---------------------------
func checkErr(err error) {
	if err != nil {
		//panic(err)
		fmt.Println(err)

	}
}

//-------------------------------generate a  jwt token-------------------------------------
func ExampleNew(mySigningKey []byte) (string, error) {

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.

	/*
		token := jwt.New(jwt.SigningMethodRS512)
		claims := make(jwt.MapClaims)
		claims["foo"] = "halloaskkakslkkaskljaskjjashhjasjajshjas"
		claims["iat"] = time.Now().Unix()
		token.Claims = claims
	*/

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = jwt.MapClaims{
		"sn":  sn_g,
		"iat": time.Now().Unix(),
	}

	tokenString, err := token.SignedString(mySigningKey)
	return tokenString, err

}

//----------------------------------parser jwt token------------------------------------------

//--------jwt token create---------
//http://dghubble.com/blog/posts/json-web-tokens-and-go/
/*
func ExampleNew(mySigningKey []byte) (string, error) {
	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	token.Claims["sn"] = sn_g
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(mySigningKey)
	return tokenString, err
}
*/
func ExampleParse(myToken string, myKey string) string {

	token, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(myKey), nil
	})

	// if err1 != nil {
	//  panic(err1)
	// }
	//fmt.Println(err1)
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		b := fmt.Sprint(claims["sn"])
		return b
		//	return b
	} else {
		fmt.Println(err)
	}
	//fmt.Println(token.Claims["exp"])

	if err == nil && token.Valid {
		fmt.Println("Your token is valid.  I like your style.")
	} else {
		fmt.Println("This token is terrible!  I cannot accept this.")
		return "error"
	}
	// sample token string taken from the New example
	//tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.u1riaD1rW97opCoAuRCTy4w58Br-Zk-bh7vLiRIsrpU"

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	return "allok"

}
