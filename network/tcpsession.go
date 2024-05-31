package network

import (
	"net"
	"sync"
	"time"

	"app/container/safe/mpsc"
	"app/goroutine"
	"go.uber.org/atomic"
)

func newTcpSession(conn net.Conn, coder ICodec, handler INetHandler) *TcpSession {
	session := &TcpSession{
		id:         genNetSessionId(),
		conn:       conn,
		stopped:    make(chan struct{}),
		coder:      coder,
		handler:    handler,
		sendQue:    mpsc.New[interface{}](),
		notifySend: make(chan struct{}, 1),
	}
	handler.setSession(session)
	return session
}

type TcpSession struct {
	id       uint32
	conn     net.Conn
	storage  sync.Map
	stopOnce sync.Once
	stopped  chan struct{}

	coder   ICodec
	handler INetHandler

	sendQue    *mpsc.Queue[interface{}]
	notifySend chan struct{}
	sending    atomic.Int32
}

func (s *TcpSession) Type() SessionType    { return TYPE_TCP }
func (s *TcpSession) Id() uint32           { return s.id }
func (s *TcpSession) LocalAddr() net.Addr  { return s.conn.LocalAddr() }
func (s *TcpSession) RemoteAddr() net.Addr { return s.conn.RemoteAddr() }
func (s *TcpSession) RemoteIP() string {
	addr := s.RemoteAddr()
	switch v := addr.(type) {
	case *net.UDPAddr:
		if ip := v.IP.To4(); ip != nil {
			return ip.String()
		}
	case *net.TCPAddr:
		if ip := v.IP.To4(); ip != nil {
			return ip.String()
		}
	case *net.IPAddr:
		if ip := v.IP.To4(); ip != nil {
			return ip.String()
		}
	}
	return ""
}

func (s *TcpSession) StoreKV(key, value interface{}) { s.storage.Store(key, value) }
func (s *TcpSession) DeleteKV(key interface{})       { s.storage.Delete(key) }
func (s *TcpSession) Load(key interface{}) (value interface{}, ok bool) {
	return s.storage.Load(key)
}

func (s *TcpSession) start() {
	goroutine.Try(s.handler.OnSessionCreated, nil)
	goroutine.GoLogic(s.read, func(ex interface{}) { s.Stop() })
	goroutine.GoLogic(s.write, func(ex interface{}) { s.Stop() })
}

func (s *TcpSession) Stop() {
	s.stopOnce.Do(func() {
		s.conn.Close()
		close(s.stopped)
		goroutine.Try(s.handler.OnSessionClosed, nil)
	})
}

func (s *TcpSession) SendMsg(msg interface{}) {
	s.sendQue.Push(msg)
	if s.sending.CAS(0, 1) {
		s.notifySend <- struct{}{}
	}
}

func (s *TcpSession) read() {
	buffer := make([]byte, 1024)
	for {
		if err := s.conn.SetReadDeadline(time.Now().Add(time.Second * 15)); err != nil {
			break
		}

		n, err := s.conn.Read(buffer)
		if err != nil {
			if e, ok := err.(net.Error); !ok || !e.Timeout() {
				break
			}
		}

		if n == 0 {
			continue
		}

		msgIds, datas, err := s.coder.Decode(buffer[:n])
		if err != nil {
			break
		}

		for i, msgId := range msgIds {
			goroutine.Try(func() { s.handler.OnRecv(msgId, datas[i]) }, nil)
		}
	}
	s.Stop()
}

func (s *TcpSession) write() {
	var err error
loop:
	for {
		select {
		case <-s.notifySend:
			s.sending.Store(0)
			for data := s.sendQue.Pop(); data != nil; data = s.sendQue.Pop() {
				if edata := s.coder.Encode(data); edata != nil {
					if err = s.conn.SetWriteDeadline(time.Now().Add(time.Second * 5)); err == nil {
						_, err = s.conn.Write(edata)
					}
					if err != nil {
						break loop
					}
				}
			}
		case <-s.stopped:
			break loop
		}
	}
	s.Stop()
}
