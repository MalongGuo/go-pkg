package test

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/MalongGuo/go-pkg/ws"
	"github.com/gorilla/websocket"
)

// 使用自定义的上下文键类型, 避免与其他包的键冲突
type contextKey int

const (
	ctxKeyUserID contextKey = iota
)

// WebsocketM 测试用的 websocket 管理器实现
type WebsocketM struct {
	sync.RWMutex
	UserMap map[int64]*websocket.Conn
}

// SetData 设置请求数据
func (ws *WebsocketM) SetData(r *http.Request) (*http.Request, error) {
	return r.WithContext(context.WithValue(r.Context(), ctxKeyUserID, int64(1))), nil
}

// GetData 获取请求数据
func (ws *WebsocketM) GetData(r *http.Request) (any, error) {
	return r.Context().Value(ctxKeyUserID), nil
}

// NewConn 创建新连接
func (ws *WebsocketM) NewConn(conn *websocket.Conn, data any) error {
	fmt.Println(data)
	ws.Lock()
	defer ws.Unlock()
	ws.UserMap[data.(int64)] = conn
	return nil
}

// CheckAuth 检查认证
func (ws *WebsocketM) CheckAuth(r *http.Request) (*http.Request, error) {
	return r, nil
}

// WriteMessageTo 向连接写入消息
func (ws *WebsocketM) WriteMessageTo(conn *websocket.Conn, messageType int, p []byte) error {
	return conn.WriteMessage(messageType, p)
}

// GetConnByID 根据用户ID获取连接
func (ws *WebsocketM) GetConnByID(userID int64) (*websocket.Conn, error) {
	ws.RLock()
	defer ws.RUnlock()
	return ws.UserMap[userID], nil
}

// ReceiveMessage 接收消息处理
func (ws *WebsocketM) ReceiveMessage(conn *websocket.Conn, messageType int, p []byte) error {
	err := ws.WriteMessageTo(conn, messageType, p)
	return err
}

// 确保 WebsocketM 实现了 WebsocketMI 接口
var _ ws.WebsocketMI = (*WebsocketM)(nil)
