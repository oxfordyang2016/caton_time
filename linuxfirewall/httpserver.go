package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	//      "reflect"
)

func hello(w http.ResponseWriter, req *http.Request) {

	io.WriteString(w, "Hello world,i am from  server!")
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(body))

}

func main() {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":8000", nil)
}
