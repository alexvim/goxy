package server

import (
	"context"
	"fmt"
	"goxy/internal/socks5"
	"log"
	"net"
	"sync"
	"time"
)

type Session struct {
	uuid     string
	bytesSnd int64
	bytesRcv int64
}

func NewSession(uuid string) *Session {
	return &Session{uuid: uuid}
}

func (s *Session) Run(ctxParent context.Context, localAddress string, sourceConn net.Conn) {
	log.Printf("%s: start\n", s)

	var targetConn net.Conn

	defer func() {
		sourceConn.Close()

		if targetConn != nil {
			targetConn.Close()
		}

		log.Printf("%s: stop snd=%d rcv=%d\n", s, s.bytesSnd, s.bytesRcv)
	}()

	flow := socks5.NewFlow(s.uuid, sourceConn, func(addr string, port uint16) (string, uint16, error) {
		dialer := net.Dialer{
			LocalAddr: &net.TCPAddr{
				IP:   net.ParseIP(localAddress),
				Port: 0,
			},
			Timeout: 3 * time.Second,
		}

		conn, err := dialer.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
		if err != nil {
			log.Printf("flow[%s]: failed to dial to %s err=%s", s.uuid, fmt.Sprintf("%s:%d", addr, port), err)
			return "", 0, err
		}

		targetConn = conn

		localAddr := targetConn.LocalAddr().(*net.TCPAddr)

		return localAddr.IP.String(), uint16(localAddr.Port), nil
	})

	if err := flow.Run(); err != nil {
		log.Printf("%s: failed to initiate socks5 tunnel err=%s", s, err)
		return
	}

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(ctxParent)

	wg.Add(1)
	go func() {
		defer wg.Done()

		s.bytesSnd, _ = stream(targetConn, sourceConn)

		targetConn.Close()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		s.bytesRcv, _ = stream(sourceConn, targetConn)

		sourceConn.Close()
	}()

	go func() {
		<-ctx.Done()

		sourceConn.Close()
		targetConn.Close()
	}()

	wg.Wait()

	cancel()
}

func (s Session) String() string {
	return fmt.Sprintf("session[%s]", s.uuid)
}
