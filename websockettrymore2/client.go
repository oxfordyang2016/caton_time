/* EchoClient
 */
package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], "ws://host:port")
		os.Exit(1)
	}
	service := os.Args[1]

	conn, err := websocket.Dial(service, "", "http://localhost")
	//client in local(the third arg)
	checkError(err)
	var msg string
	for {
		//info write to msg from server
		err := websocket.Message.Receive(conn, &msg)

		if err != nil {
			if err == io.EOF {
				// graceful shutdown by server
				break
			}
			fmt.Println("Couldn't receive msg " + err.Error())
			break
		}
		fmt.Println("Received from server:=================> " + msg)
		// return the msg
		//========================send info=================================//
		//sendinfo := "i love go"
		sendjson := `{"from":"a", "cmd": "register", "type": "rsp",
		"seq": 11,"task_id":"67","uid":"6732","token": "EF02JLGFA09GVNG21F","rsp":{"register":{"tnid":"56"}}}`
		err = websocket.Message.Send(conn, sendjson)
		if err != nil {
			fmt.Println("Coduln't return msg")
			break
		}
	}
	os.Exit(0)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
