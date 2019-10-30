package castcenter

import (
	"net"
	"time"
)

const (
	datagramChannelBufferSize = 4 * 1024
	datagramReadBufferSize    = 4 * 1024
)

// Handler .
type Handler func(*UDPEvent)

// UDPServer .
type UDPServer struct {
	handler    Handler
	ch         chan *UDPEvent
	connection *net.UDPConn
}

// UDPEvent .
type UDPEvent struct {
	ip  string
	buf []byte
}

// NewUDPServer returns a new UDPServer
func NewUDPServer(handler Handler, chanSize int) *UDPServer {
	return &UDPServer{
		handler: handler,
		ch:      make(chan *UDPEvent, chanSize),
	}
}

// ListenUDP Configure the UDPServer for listen on an UDP addr
func (s *UDPServer) ListenUDP(conn *net.UDPConn) error {
	s.connection = conn
	conn.SetReadBuffer(datagramReadBufferSize)

	go recoverLoop(s.receiveDatagrams)
	recoverLoop(s.parseDatagrams)

	return nil
}

func (s *UDPServer) receiveDatagrams() {
	for {
		buf := GetBytes()
		n, addr, err := s.connection.ReadFromUDP(buf)
		if err == nil {
			if n > 0 {
				s.ch <- &UDPEvent{
					ip:  addr.IP.String(),
					buf: buf[:n],
				}
			}
		} else {
			// there has been an error. Either the UDPServer has been killed
			// or may be getting a transitory error due to (e.g.) the
			// interface being shutdown in which case sleep() to avoid busy wait.
			opError, ok := err.(*net.OpError)
			if (ok) && !opError.Temporary() && !opError.Timeout() {
				return
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func (s *UDPServer) parseDatagrams() {
	for {
		select {
		case msg, ok := <-s.ch:
			if !ok {
				return
			}
			s.handle(msg)
		}
	}
}

func (s *UDPServer) handle(msg *UDPEvent) {
	defer PutBytes(msg.buf)
	s.handler(msg)
}
