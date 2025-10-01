package ws

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

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

type WebsocketM struct {
	sync.RWMutex
	UserMap map[int64]*websocket.Conn
}

func (ws *WebsocketM) SetData(r *http.Request) (*http.Request, error) {
	return r.WithContext(context.WithValue(r.Context(), "user_id", int64(1))), nil
}

func (ws *WebsocketM) GetData(r *http.Request) (any, error) {
	return r.Context().Value("user_id"), nil
}

func (ws *WebsocketM) NewConn(conn *websocket.Conn, data any) error {
	fmt.Println(data)
	ws.Lock()
	defer ws.Unlock()
	ws.UserMap[data.(int64)] = conn
	return nil
}

func (ws *WebsocketM) CheckAuth(r *http.Request) (*http.Request, error) {

	return r, nil
}

func (ws *WebsocketM) WriteMessageTo(conn *websocket.Conn, messageType int, p []byte) error {
	return conn.WriteMessage(messageType, p)
}

func (ws *WebsocketM) GetConnByID(userID int64) (*websocket.Conn, error) {
	ws.RLock()
	defer ws.RUnlock()

	return ws.UserMap[userID], nil
}

func (ws *WebsocketM) ReceiveMessage(conn *websocket.Conn, messageType int, p []byte) error {
	err := ws.WriteMessageTo(conn, messageType, p)
	return err
}

var _ WebsocketMI = (*WebsocketM)(nil)

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
