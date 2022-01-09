package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var SECRET string

type connection struct {
	conn     *websocket.Conn
	connType string
}

var connections []connection
var upgrader = websocket.Upgrader{
	HandshakeTimeout: 0,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	WriteBufferPool:  nil,
	Subprotocols:     []string{},
	Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
		var error strings.Builder
		error.WriteString("Error! Status: ")
		error.WriteString(string(status))
		errorString := error.String()
		fmt.Fprintf(w, errorString)
	},
	CheckOrigin: func(r *http.Request) bool {
		if SECRET == r.URL.Query().Get("a") {
			return true
		}
		return false
	},
	EnableCompression: false,
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Use /ws route to subscribe my websockets")
}

func reader(conn *websocket.Conn) {
	if conn != nil {
		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			log.Println((string(p)))
			for _, elem := range connections {
				if elem.connType == "listener" {
					if err := elem.conn.WriteMessage(messageType, p); err != nil {
						log.Println(err)
					}
				}
			}
		}
	}
}

func indexOfConn(element connection, data []connection) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1
}

func removeConn(s []connection, i int) []connection {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	connType := r.URL.Query().Get("t")
	ws, err := upgrader.Upgrade(w, r, nil)
	conn := connection{connType: connType, conn: ws}

	ws.SetCloseHandler(func(code int, text string) error {
		fmt.Println(text)
		indexOfConnection := indexOfConn(conn, connections)
		connections = removeConn(connections, indexOfConnection)
		var returned error
		return returned
	})

	if err != nil {
		log.Println(err)
	} else {
		log.Println("Client Succesfully Connected")
		connections = append(connections, conn)
		reader(ws)
	}
}
func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	if len(os.Args) > 1 {
		usesDotenv, err := strconv.ParseBool(os.Args[1])
		if usesDotenv && err == nil {
			err := godotenv.Load()
			if err != nil {
				log.Fatal("Error loading .env file")
			}
		}
	}
	SECRET = os.Getenv("SECRET")
	fmt.Println("Go WebSockets")
	setupRoutes()
	//log.Fatal(http.ListenAndServe(":80", nil))
	log.Fatal(http.ListenAndServeTLS(":443", "/etc/letsencrypt/live/modularizar.com/fullchain.pem", "/etc/letsencrypt/live/modularizar.com/privkey.pem", nil))
}
