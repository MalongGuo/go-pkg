package ws

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 使用 Gin 注册一个 GET /ws 路由示例
func TestGinWebsocketRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 初始化测试用的 Websocket 管理器
	m := &WebsocketM{UserMap: make(map[int64]*websocket.Conn)}

	// 仅注册一个 GET /ws 路由
	r.GET("/ws", func(c *gin.Context) {
		Handle(c.Writer, c.Request, m)
	})

	// 不发起普通 HTTP 请求, 以避免非 WebSocket 请求导致升级失败
	_ = r
	r.Run() // listen and serve on 0.0.0.0:8080
}
