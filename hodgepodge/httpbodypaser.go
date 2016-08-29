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
}

func main() {
	http.HandleFunc("/test", test)
	log.Fatal(http.ListenAndServe(":8082", nil))
}
