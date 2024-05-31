package network

import (
	"net"
	"strings"
	"sync/atomic"
)

func NewTcpListener(addr string, newCodec func() ICodec, newHanlder func() INetHandler, options ...func(*TcpListener) error) INetListener {
	l := &TcpListener{
		addr:       addr,
		newCodec:   newCodec,
		newHanlder: newHanlder,
	}

	for _, f := range options {
		if err := f(l); err != nil {
			return nil
		}
	}
	return l
}

type TcpListener struct {
	addr       string
	listener   net.Listener
	running    int32
	newCodec   func() ICodec
	newHanlder func() INetHandler
}

func (s *TcpListener) Listen() error {
	err := s.listen()
	if err == nil {
		return nil
	}
	return err
}

func (s *TcpListener) listen() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.listener = listener

	for {
		if conn, err := listener.Accept(); err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				break
			}
		} else {
			newTcpSession(conn, s.newCodec(), s.newHanlder()).start()
		}
	}
	s.Stop()
	return nil
}

func (s *TcpListener) Stop() {
	if atomic.CompareAndSwapInt32(&s.running, 1, 0) {
		if l := s.listener; l != nil {
			_ = s.listener.Close()
		}
	}
}
