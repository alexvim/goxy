package main

import (
	"fmt"
	"goxy/handler"
	"goxy/network"
	"net"
	"os"
)

func main() {

	localAddr4, _ := network.GetLocalInetAddress(network.AddrTypeIpv4)

	localAddr4 += ":1080"
	listener, err := net.Listen("tcp", localAddr4)
	if err != nil {
		fmt.Printf("main: failed to listern %s port err={%s}", localAddr4, err.Error())
		os.Exit(1)
	}

	fmt.Printf("main: start listering on %s\n", localAddr4)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("main: failed to accept connection on %s port err={%s}", localAddr4, err.Error())
			os.Exit(1)
		}

		session := handler.MakeSession(conn)
		go session.Run()
	}
}
