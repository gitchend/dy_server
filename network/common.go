package network

import (
	"net"
	"sync/atomic"
	"time"
)

type SessionType int

const (
	TYPE_TCP SessionType = 1
	TYPE_UDP SessionType = 2
	TYPE_WS  SessionType = 3
)

type INetListener interface {
	Listen() error
	Stop()
}

type INetClient interface {
	SendMsg(interface{}) error
	Connect(reconect bool) error
	Stop()
}

type INetSession interface {
	Id() uint32
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	RemoteIP() string
	SendMsg(interface{})
	Stop()
	StoreKV(interface{}, interface{})
	DeleteKV(interface{})
	Load(interface{}) (interface{}, bool)
	Type() SessionType
}

type INetHandler interface {
	OnSessionCreated()
	OnSessionClosed()
	OnRecv(int32, interface{})

	setSession(session INetSession)
}

type BaseNetHandler struct {
	INetSession
}

func (s *BaseNetHandler) setSession(session INetSession) { s.INetSession = session }

var genNetSessionId = _gen_net_session_id()

func _gen_net_session_id() func() uint32 {
	now := time.Now()
	_session_gen_id := uint32(now.Hour()*100000000 + now.Minute()*1000000 + now.Second()*10000)
	return func() uint32 { return atomic.AddUint32(&_session_gen_id, 1) }
}
