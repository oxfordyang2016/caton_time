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
	"reflect"
	"strings"
	//"time"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type test_struct struct {
	Test string
}
type Response struct {
	Token  string
	Status int
}

func test(rw http.ResponseWriter, req *http.Request) {
	//------------------------------------------database connettion----------------------------
	db, err := sql.Open("mysql",
		"root:123456@tcp(192.168.0.81:3306)/hello") //first configure a database
	if err != nil {
		fmt.Println("database error")
		log.Fatal(err)
	}
	fmt.Println("database connect success")
	defer db.Close()
	//---------------------------------------------------------------------------------------------------
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
	fmt.Println(reflect.TypeOf(to))
	rep := &Response{Token: to, Status: 1}
	b, _ := json.Marshal(rep)
	io.WriteString(rw, string(b)) //write to body

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
	//--------------------------------------------------------------
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
	log.Fatal(http.ListenAndServe(":8082", nil))

}

// formatRequest generates ascii representation of a request
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
			fmt.Println(name, h)
			request = append(request, fmt.Sprintf(name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
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
