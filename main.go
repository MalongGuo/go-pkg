package main

import (
	"flag"
	"log"

	"github.com/MalongGuo/go-pkg/test"
	"github.com/MalongGuo/go-pkg/ws"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func main() {
	// 解析命令行参数: -p 指定端口, 默认 8080
	var port string
	flag.StringVar(&port, "p", "8080", "port to listen")
	flag.Parse()

	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 初始化测试用的 Websocket 管理器
	m := &test.WebsocketM{UserMap: make(map[int64]*websocket.Conn)}

	// 仅注册一个 GET /ws 路由
	r.GET("/ws", func(c *gin.Context) {
		ws.Handle(c.Writer, c.Request, m)
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, World!"})
	})

	addr := ":" + port
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	} // listen and serve on 0.0.0.0:8080
}
