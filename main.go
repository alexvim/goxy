package main

import (
	"fmt"
	"goxy/handler"
	"goxy/network"
	"net"
	"os"
)

func main() {

	localAddr, _ := network.GetLocalInetAddress(network.AddrTypeIpv4)
	localAddr += ":1080"
	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		fmt.Printf("main: failed to listern %s port err={%s}", localAddr, err.Error())
		os.Exit(1)
	}

	fmt.Printf("main: start listering on %s", localAddr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("main: failed to accept connection on %s port err={%s}", localAddr, err.Error())
			os.Exit(1)
		}

		session := handler.MakeSession(conn)
		go session.Run()
	}
}
