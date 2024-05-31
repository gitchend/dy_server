package game_mgr

import (
	"app/in_obj"
	"app/message"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"net/http"
	"time"
)

type HttpMgr struct {
	r   *gin.Engine
	mgr in_obj.IGameMgr
}

func NewHttpMgr(mgr in_obj.IGameMgr) *HttpMgr {
	return &HttpMgr{
		r:   gin.Default(),
		mgr: mgr,
	}
}

func (s *HttpMgr) Run() error {
	//公共-websocket
	s.r.GET("/socket", s.OnWebSocket)

	//抖音
	s.r.GET("/v1/ping", s.OnDouyinPing)
	s.r.HEAD("/douyin/push/:appId", s.OnDouyinDataPush)
	s.r.POST("/douyin/push/:appId", s.OnDouyinDataPush)
	return s.r.Run(":8000")
}

func (s *HttpMgr) OnWebSocket(c *gin.Context) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		CheckOrigin:      func(r *http.Request) bool { return true },
		HandshakeTimeout: 5 * time.Second,
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	conn.SetPingHandler(func(appData string) error {
		return conn.WriteMessage(websocket.PongMessage, []byte(appData))
	})
	closeChan := make(chan struct{})
	sendChan := make(chan proto.Message)
	gameUser := s.mgr.NewUser()
	gameUser.OnSessionCreated(sendChan)
	defer func() {
		gameUser.OnSessionClosed()
		conn.Close()
	}()
	//写协程
	go func() {
		for {
			select {
			case msg := <-sendChan:
				err := message.SendMessage(conn, msg)
				if err != nil {
					gameUser.Log("SendMessage err", err)
					return
				}
			case <-closeChan:
				return
			}
		}
	}()
	for {
		msg, err := message.ReadMessage(conn)
		if err != nil {
			closeChan <- struct{}{}
			close(closeChan)
			return
		}
		gameUser.OnRecv(msg)
	}
}
