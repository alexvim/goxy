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
		fmt.Println("main: failed to listern 1080 port err=" + err.Error())
		os.Exit(1)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("main: failed to acceprt connection err=" + err.Error())
			os.Exit(1)
		}

		session := handler.MakeSession(conn)
		go session.Run()
	}
}
