package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"time"
)

func main() {
	origin := "http://192.168.0.53:8080/"
	url := "ws://192.168.0.53:8080/ws"

	var err error
	var ws *websocket.Conn
	for {
		ws, err = websocket.Dial(url, "", origin)
		if err != nil {
			fmt.Println("Connection fails, is being re-connection")
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}
	if _, err := ws.Write([]byte("hallo wangyan")); err != nil {
		log.Fatal(err)
	}

}
