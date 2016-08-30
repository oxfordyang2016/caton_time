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

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

type test_struct struct {
	Test string
}

func test(rw http.ResponseWriter, req *http.Request) {
//-----------------------http body deal--------------------------------
	body, err := ioutil.ReadAll(req.Body)
	fmt.Println(reflect.TypeOf(body))
	fmt.Println(reflect.TypeOf(req.Body))
	fmt.Println(reflect.TypeOf(string(body))
	fmt.Println(req.Body)
	fmt.Println(body) //if only here,it will occur  byte
	fmt.Println(string(body))
	if err != nil {
		//panic()
		fmt.Println("erorr")
	}

	log.Println(string(body))

	var t test_struct
	err = json.Unmarshal(body, &t)
	if err != nil {
		//panic()
		fmt.Println("error")
	}
	log.Println(t.Test)

	fmt.Println(formatRequest(req))//request header deal

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


}

func main() {
	http.HandleFunc("/test", test)
	log.Fatal(http.ListenAndServe(":8082", nil))
}
//=---------------------------------------http header deal-----------------------------------------
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
//-----------------------------------------------------------------------------------------