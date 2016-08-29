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
)

type test_struct struct {
	Test string
}

func test(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		//panic()
		fmt.Println("erorr")
	}
	log.Println(string(body))
	fmt.Println(string(body))
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
