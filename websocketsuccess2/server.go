/* EchoServer
 */
package main

import (
	"fmt"
	"net/http"
	"os"
	// "io"
	"golang.org/x/net/websocket"
)

func Echo(ws *websocket.Conn) {// in this connect,do some actions 
	fmt.Println("Echoing")

	for n := 0; n < 10; n++ {
		msg := "Hello client " + string(n+48)
		fmt.Println("Sending to client: " + msg)
		err := websocket.Message.Send(ws, msg) 
		//ws is websocket connect,send infomation

		if err != nil {
			fmt.Println("Can't send")
			break
		}

		var reply string
		err = websocket.Message.Receive(ws, &reply) //receive info
		if err != nil {
			fmt.Println("Can't receive")
			break
		}
		fmt.Println("Received from client: " + reply)
	}
}

func main() {

	http.Handle("/", websocket.Handler(Echo))
	//this is route,accept info and deal with it
	err := http.ListenAndServe(":88", nil) //this is to launch  server
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
