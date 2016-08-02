package main

import (
	"code.google.com/p/go.net/websocket"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	listenAddr = "localhost:9527" // server address
)

var (
	pwd, _        = os.Getwd()
	RootTemp      = template.Must(template.ParseFiles(pwd + "/chat.html"))
	JSON          = websocket.JSON              // codec for JSON
	Message       = websocket.Message           // codec for string, []byte
	ActiveClients = make(map[ClientConn]string) // map containing clients
	User          = make(map[string]string)
)

// Initialize handlers and websocket handlers
func init() {
	User["aaa"] = "aaa"
	User["bbb"] = "bbb"
	User["test"] = "test"
	User["test2"] = "test2"
	User["test3"] = "test3"
}

// Client connection consists of the websocket and the client ip
type ClientConn struct {
	websocket *websocket.Conn
	clientIP  string
}

// WebSocket server to handle chat between clients
func SockServer(ws *websocket.Conn) {
	var err error
	var clientMessage string
	// use []byte if websocket binary type is blob or arraybuffer
	// var clientMessage []byte

	// cleanup on server side
	defer func() {
		if err = ws.Close(); err != nil {
			log.Println("Websocket could not be closed", err.Error())
		}
	}()

	client := ws.Request().RemoteAddr
	log.Println("Client connected:", client)
	sockCli := ClientConn{ws, client}
	ActiveClients[sockCli] = ""
	log.Println("Number of clients connected:", len(ActiveClients))

	// for loop so the websocket stays open otherwise
	// it'll close after one Receieve and Send
	for {
		if err = Message.Receive(ws, &clientMessage); err != nil {
			// If we cannot Read then the connection is closed
			log.Println("Websocket Disconnected waiting", err.Error())
			// remove the ws client conn from our active clients
			delete(ActiveClients, sockCli)
			log.Println("Number of clients still connected:", len(ActiveClients))
			return
		}

		var msg_arr = strings.Split(clientMessage, "|")
		if msg_arr[0] == "login" && len(msg_arr) == 3 {
			name := msg_arr[1]
			pass := msg_arr[2]

			if pass == User[name] {
				ActiveClients[sockCli] = name

				if err = Message.Send(ws, "login|"+name); err != nil {
					log.Println("Could not send message to ", client, err.Error())
				}
			} else {
				log.Println("login faild:", clientMessage)
			}

		} else if msg_arr[0] == "msg" {
			if ActiveClients[sockCli] != "" {
				clientMessage = "msg|" + time.Now().Format("2006-01-02 15:04:05") + " " + ActiveClients[sockCli] + " Said: " + msg_arr[1]
				for cs, na := range ActiveClients {
					if na != "" {
						if err = Message.Send(cs.websocket, clientMessage); err != nil {
							log.Println("Could not send message to ", cs.clientIP, err.Error())
						}
					}
				}
			}
		}
	}
}

// RootHandler renders the template for the root page
func RootHandler(w http.ResponseWriter, req *http.Request) {
	err := RootTemp.Execute(w, listenAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", RootHandler)
	http.Handle("/socket", websocket.Handler(SockServer))
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
