package network

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"app/goroutine"
)

func NewTcpClient(addr string, newCodec func() ICodec, newHanlder func() INetHandler, op ...func(*TcpClient) error) INetClient {
	c := &TcpClient{
		addr:             addr,
		running:          1,
		newCodec:         newCodec,
		newHanlder:       newHanlder,
		sessionReconnect: make(chan struct{}, 1),
	}
	c.newHanlder = func() INetHandler { return &TcpClientHandler{client: c, handler: newHanlder()} }

	for _, f := range op {
		if err := f(c); err != nil {
			return nil
		}
	}
	return c
}

type TcpClient struct {
	addr    string
	running int32

	newCodec   func() ICodec
	newHanlder func() INetHandler

	session          *TcpSession
	reconnectTimes   int
	sessionReconnect chan struct{}

	connMux sync.Mutex
}

func (s *TcpClient) Connect(reconnect bool) error {
	if reconnect {
		goroutine.GoLogic(s.reconnect, nil)
	} else {
		return s.connect()
	}
	return nil
}

func (s *TcpClient) connect() error {
	s.connMux.Lock()
	defer s.connMux.Unlock()

	if !s.IsRunning() {
		return errors.New("tcp client has stopped")
	}

	conn, err := net.DialTimeout("tcp", s.addr, time.Second)
	if err != nil {
		return err
	}

	session := newTcpSession(conn, s.newCodec(), s.newHanlder())
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&s.session)), unsafe.Pointer(session))
	session.start()
	return nil
}

func (s *TcpClient) reconnect() {
	s.sessionReconnect <- struct{}{}
	for range s.sessionReconnect {
		if !s.IsRunning() {
			return
		}

		if err := s.connect(); err == nil {
			s.reconnectTimes = 0
			break
		}

		time.Sleep(time.Second * time.Duration(s.reconnectTimes))
		s.reconnectTimes++
		s.sessionReconnect <- struct{}{}
	}
}

func (s *TcpClient) SendMsg(data interface{}) error {
	session := (*TcpSession)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&s.session))))
	if session != nil {
		session.SendMsg(data)
		return nil
	}
	return errors.New("has not connect")
}

func (s *TcpClient) IsRunning() bool { return atomic.LoadInt32(&s.running) == 1 }
func (s *TcpClient) Stop() {
	if atomic.CompareAndSwapInt32(&s.running, 1, 0) {
		s.connMux.Lock()
		defer s.connMux.Unlock()

		s.session.Stop()
	}
}

type TcpClientHandler struct {
	BaseNetHandler
	client  *TcpClient
	handler INetHandler
}

func (s *TcpClientHandler) OnSessionCreated() {
	s.handler.setSession(s.INetSession)
	s.handler.OnSessionCreated()
}

func (s *TcpClientHandler) OnSessionClosed() {
	s.handler.OnSessionClosed()
	s.client.sessionReconnect <- struct{}{}
}

func (s *TcpClientHandler) OnRecv(msgId int32, data interface{}) {
	s.handler.OnRecv(msgId, data)
}
