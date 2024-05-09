package websocket

import (
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
	l := logger.GetLogger()
	if err != nil {
		l.Errorf("websocket推送错误： %s", err.Error())
		return
	}
}

// CloseSocket
//
//	@Description: 关闭socket
func CloseSocket() {
	l := logger.GetLogger()
	l.Info("###################### websocket连接关闭 ########################")
	err := conn.Close()
	if err != nil {
		l.Errorf("websocket关闭异常:%s", err.Error())
	}
}
