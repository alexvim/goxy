package handler

import (
	"fmt"
	"goxy/msg"
	"net"
)

// Session ...
type Session struct {
	client *Client
	nif    *Nif
}

// MakeSession ...
func MakeSession(conn net.Conn) *Session {
	b := new(Session)
	b.client = MakeClient(conn)
	fmt.Printf("Make new session for client: uuid=%s, ip=%s\n", b.client.id, b.client.connection.RemoteAddr().String())
	return b
}

// Run ...
func (s *Session) Run() {

	//defer b.Destroy()

	buf, err := s.client.ReadMessage()
	if err != nil {
		fmt.Println("session: failed to read message err=" + err.Error())
		return
	}

	auth, err := msg.ParseAuth(buf)
	if err != nil {
		fmt.Println("session: failed to pasre message err=" + err.Error())
		return
	}

	// error check here
	s.HandleAuth(auth)

	buf, err = s.client.ReadMessage()
	if err != nil {
		fmt.Println("session: failed to read message err=" + err.Error())
		return
	}

	cmd, err := msg.ParseCommand(buf)
	if err != nil {
		fmt.Println("session: failed to pasre buffer err=" + err.Error())
		return
	}

	switch cmd.Command {
	case msg.CONNECT:
		s.HandleConnect(cmd)
	default:
		fmt.Println("session: wrong message type")
		return
	}
}

// HandleAuth ...
func (s *Session) HandleAuth(message *msg.AuthRequest) {

	fmt.Printf("Handle auth request %s\n", message)

	for _, v := range message.Methods {
		switch v {
		case msg.NO_AUTHENTICATION_REQUIRED:
			fmt.Println("Selected auth method NO_AUTHENTICATION_REQUIRED")

			reply := msg.AuthReply{
				Method: msg.NO_AUTHENTICATION_REQUIRED,
			}
			s.client.SendMessage(reply)
		}
	}
}

// HandleConnect ...
func (s *Session) HandleConnect(message *msg.CommandRequest) {

	fmt.Printf("Handle connect request %s\n", message)

	nif := MakeNif(s.client.connection, message.DstAddr, message.DstPort)

	addr, port, err := nif.Prepare()
	if err != nil {
		fmt.Println("session: failed do connect err=" + err.Error())
		cr := msg.CommandReply{
			Result:      msg.CommandResultNetworkUnreaschable,
			AddressType: msg.IP_V4ADDRESS,
			BindAddress: "0.0.0.0",
			BindPort:    0,
		}
		s.client.SendMessage(cr)
		s.client.Disconnect()
		return
	}

	s.nif = nif

	go s.nif.Run()

	cr := msg.CommandReply{
		Result:      msg.CommandResultSucceeded,
		AddressType: msg.IP_V4ADDRESS,
		BindAddress: addr,
		BindPort:    port,
	}

	s.client.SendMessage(cr)
}
