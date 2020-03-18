package handler

import (
	"errors"
	"fmt"
	"goxy/msg"
	"net"
)

// Session ...
type Session struct {
	connection net.Conn
	nif        *Nif
}

// MakeSession ...
func MakeSession(conn net.Conn) *Session {
	s := new(Session)
	s.connection = conn
	s.nif = nil
	fmt.Printf("Make new session for client: ip=%s\n", s.connection.RemoteAddr().String())
	return s
}

// SendMessage ...
func (s *Session) SendMessage(m msg.Serializeable) bool {

	_, err := s.connection.Write(m.Serialize())
	if err != nil {
		fmt.Printf("client: failed to send message %s\n", err)
		return false
	}

	return true
}

// ReadMessage ...
func (s *Session) ReadMessage() ([]byte, error) {

	var buffer []byte = make([]byte, 50)
	n, err := s.connection.Read(buffer)
	if err != nil || n <= 0 {
		return nil, errors.New("failed to read buffer err=" + err.Error())
	}

	return buffer[0:n], nil
}

// Disconnect ...
func (s *Session) Disconnect() {
	s.connection.Close()
}

// Run ...
func (s *Session) Run() {

	buf, err := s.ReadMessage()
	if err != nil {
		fmt.Println("session: failed to read message err=" + err.Error())
		s.Disconnect()
		return
	}

	auth, err := msg.ParseAuth(buf)
	if err != nil {
		fmt.Println("session: failed to pasre message err=" + err.Error())
		s.Disconnect()
		return
	}

	// error check here
	s.HandleAuth(auth)

	buf, err = s.ReadMessage()
	if err != nil {
		fmt.Println("session: failed to read message err=" + err.Error())
		s.Disconnect()
		return
	}

	cmd, err := msg.ParseCommand(buf)
	if err != nil {
		fmt.Println("session: failed to pasre buffer err=" + err.Error())
		s.Disconnect()
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
			s.SendMessage(reply)
		}
	}
}

// HandleConnect ...
func (s *Session) HandleConnect(message *msg.CommandRequest) {

	fmt.Printf("Handle connect request %s\n", message)

	nif := MakeNif(s.connection, message.DstAddr, message.DstPort)

	addr, port, err := nif.Prepare()
	if err != nil {
		fmt.Println("session: failed do connect err=" + err.Error())
		cr := msg.CommandReply{
			Result:      msg.CommandResultNetworkUnreaschable,
			AddressType: msg.IP_V4ADDRESS,
			BindAddress: "0.0.0.0",
			BindPort:    0,
		}
		s.SendMessage(cr)
		s.Disconnect()
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

	s.SendMessage(cr)
}
