package core

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/google/uuid"
)

type Resolver interface {
	Resolve(doamin string) (string, error)
}

type Server struct {
	ProxyAddress string
	LocalAddress string
}

func (server Server) Run(ctx context.Context, doaminResolver Resolver) {
	log.Printf("server: start forwarding on %s via %s\n", server.ProxyAddress, server.LocalAddress)

	listener, err := net.Listen("tcp", server.ProxyAddress)
	if err != nil {
		log.Printf("server: failed to %s\n", err)
		return
	}

	sessions := make(map[uuid.UUID]*Session)

	wg := &sync.WaitGroup{}
	sessMutex := sync.Mutex{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		<-ctx.Done()

		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("server: failed to accept connection on %s port err={%s}", server.ProxyAddress, err.Error())
			break
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			uid := uuid.New()

			session := NewSession(uid.String())

			sessMutex.Lock()

			sessions[uid] = session

			sessMutex.Unlock()

			log.Printf("server: add session, count %d\n", len(sessions))

			session.Run(ctx, server.LocalAddress, conn, doaminResolver)

			sessMutex.Lock()

			delete(sessions, uid)

			sessMutex.Unlock()

			log.Printf("server: del session, count %d\n", len(sessions))
		}()
	}

	wg.Wait()
}
