package websocket

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sxp-server/common/logger"
)

var (
	conn *websocket.Conn
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Handler
//
//	@Description: websocket连接
//	@param c
func Handler(c *gin.Context) {
	cc, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatalln(err)
	}
	conn = cc
}

// W
//
//	@Description: 通知前端页面刷新
func W() {
	err := conn.WriteMessage(websocket.TextMessage, []byte("start"))
	fmt.Println("err: ", err)
	if err != nil {
		return
	}
}

// CloseSocket
//
//	@Description: 关闭socket
func CloseSocket() {
	l := logger.GetLogger()
	err := conn.Close()
	if err != nil {
		l.Errorf("websocket关闭异常:%s", err.Error())
	}
}
