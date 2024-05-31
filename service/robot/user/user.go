package user

import (
	"app/network"
	"fmt"
)

func NewUser() *User {
	return &User{}
}

type User struct {
	network.BaseNetHandler
}

func (s *User) OnSessionCreated() {
}

func (s *User) OnSessionClosed() {
}

func (s *User) OnRecv(msgId int32, data interface{}) {
	fmt.Println(data)
}
