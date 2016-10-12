package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
)

func hello(w http.ResponseWriter, req *http.Request) {

	io.WriteString(w, "Hello world,i am from  server!")

	body, _ := ioutil.ReadAll(req.Body)
	//header, _ := ioutil.ReadAll(req.Header)
	fmt.Println(reflect.TypeOf(body))
	fmt.Println(reflect.TypeOf(req.Body))
	fmt.Println(reflect.TypeOf(string(body)))
	fmt.Println(req.Body)
	fmt.Println(body) //if only here,it will occur  byte
	fmt.Println(string(body))

}

func main() {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":8000", nil)
}
