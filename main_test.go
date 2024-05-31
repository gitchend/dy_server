package main

import (
	"app/message"
	"app/message/pb"
	"fmt"
	"github.com/gorilla/websocket"
	"net/url"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	message.InitMessageParser("message", "MSG_TYPE")
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:8000", Path: "/socket"}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("handshake failed: ", err)
		return
	}
	defer conn.Close()

	done := make(chan struct{})
	defer close(done)

	go func() {
		for {
			msg, err := message.ReadMessage(conn)
			if err != nil {
				fmt.Println("read:", err)
				return
			}
			fmt.Println("recv:", msg)
		}
	}()

	msg := &pb.GetRank{
		TopCount: 10,
	}

	for i := 0; i < 10; i++ {
		fmt.Println(i)
		err = message.SendMessage(conn, msg)
		if err != nil {
			fmt.Println("write:", err)
			return
		}
		time.Sleep(time.Second)
	}
	_ = <-done
}
