package main

import (
	"fmt"
	"goxy/handler"
	"goxy/network"
	"net"
	"os"

	"github.com/google/uuid"
)

func main() {

	localAddr4, _ := network.GetLocalInetAddress(network.AddrTypeIpv4)

	localAddr4 += ":1080"
	listener, err := net.Listen("tcp", localAddr4)
	if err != nil {
		fmt.Printf("main: failed to listern %s port err={%s}", localAddr4, err.Error())
		os.Exit(1)
	}

	sessions := make(map[uuid.UUID]*handler.Session)
	sch := make(chan uuid.UUID, 1000)
	go func() {
		for {
			uid := <-sch
			_, ok := sessions[uid]
			if ok {
				delete(sessions, uid)
				fmt.Printf("main: session count %d\n", len(sessions))
			}
		}
	}()

	fmt.Printf("main: start listering on %s\n", localAddr4)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("main: failed to accept connection on %s port err={%s}", localAddr4, err.Error())
			os.Exit(1)
		}

		uuid := uuid.New()
		session := handler.MakeSession(conn, uuid)
		sessions[uuid] = session
		go session.Run(sch)
	}
}
