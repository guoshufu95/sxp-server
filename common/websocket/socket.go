package websocket

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sxp-server/common/logger"
)

var (
	conn       *websocket.Conn
	SocketChan = make(chan struct{}, 1)
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func StartSocket() {
	http.HandleFunc("/taskSocket", Handler)
	http.ListenAndServe(":8001", nil)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln(err)
	}
	for {
		select {
		case _ = <-SocketChan:
			err = conn.WriteMessage(websocket.TextMessage, []byte("start"))
			if err != nil {
				return
			}
		}

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
