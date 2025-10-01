package ws

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebsocketMI interface {
	CheckAuth(r *http.Request) (*http.Request, error)
	SetData(r *http.Request) (*http.Request, error)
	GetData(r *http.Request) (any, error)
	NewConn(conn *websocket.Conn, data any) error
	ReceiveMessage(conn *websocket.Conn, messageType int, p []byte) error
	WriteMessageTo(conn *websocket.Conn, messageType int, p []byte) error

	GetConnByID(userID int64) (*websocket.Conn, error)
}


var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

func Handle(w http.ResponseWriter, r *http.Request, ws WebsocketMI) {
	r, err := ws.CheckAuth(r)
	if err != nil {
		fmt.Println(err)
		return
	}
	r, _ = ws.SetData(r)

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
	}

	data, _ := ws.GetData(r)
	fmt.Println(data)
	err = ws.NewConn(c, data)
	if err != nil {
		return
	}
	go handleConn(c, ws)
}

func handleConn(c *websocket.Conn, ws WebsocketMI) {
	defer func(c *websocket.Conn) {
		err := c.Close()
		if err != nil {
			log.Println("close err:", err)
		}
	}(c)

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		err = ws.ReceiveMessage(c, mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
