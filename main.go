package main

import (
	"fmt"
	"goxy/handler"
	"net"
	"os"
)

func main() {

	var listener net.Listener
	var conn net.Conn
	var err error

	listener, err = net.Listen("tcp", "10.206.13.110:1080")
	if err != nil {
		fmt.Println("main: failed to listern 1080 port err=" + err.Error())
		os.Exit(1)
	}

	for {
		conn, err = listener.Accept()
		if err != nil {
			fmt.Println("main: failed to acceprt connection err=" + err.Error())
			os.Exit(1)
		}

		session := handler.MakeSession(conn)
		go session.Run()
	}
}
