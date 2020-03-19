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
	// nif will be closed as well
	s.connection.Close()
	s.nif = nil
	s.connection = nil
}

// Run ...
func (s *Session) Run() {

	// auth methods
	buf, err := s.ReadMessage()
	if err != nil {
		fmt.Println("session: failed to read message err=" + err.Error())
		s.Disconnect()
		return
	}

	auth, err := msg.ParseAuthHandshake(buf)
	if err != nil {
		fmt.Println("session: failed to pasre message err=" + err.Error())
		s.Disconnect()
		return
	}

	// error check here
	authRoutine := s.HandleAuth(auth)
	if !authRoutine() {
		fmt.Println("session: failed to do auth")
		s.Disconnect()
		return
	}

	// command
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
	case msg.ConnectCmd:
		s.nif = s.HandleConnect(cmd)
		if s.nif != nil {
			s.nif.Run()
		}
	default:
		fmt.Println("session: wrong message type")
		return
	}
	s.Disconnect()
}

// HandleAuth ...
func (s *Session) HandleAuth(message *msg.AuthRequest) func() bool {

	fmt.Printf("Handle auth request %s\n", message)

	for _, v := range message.Methods {
		switch v {
		case msg.NO_AUTHENTICATION_REQUIRED:
			fmt.Println("Selected auth method NO_AUTHENTICATION_REQUIRED")
			reply := msg.AuthReply{
				Method: msg.NO_AUTHENTICATION_REQUIRED,
			}
			s.SendMessage(reply)

			return func() bool {
				fmt.Println("Selected auth method NO_AUTHENTICATION_REQUIRED complete")
				return true
			}

		case msg.USERNAME_PASSWORD:
			fmt.Println("Selected auth method USERNAME_PASSWORD")
			reply := msg.AuthReply{
				Method: msg.USERNAME_PASSWORD,
			}
			s.SendMessage(reply)

			return func() bool {
				buf, err := s.ReadMessage()
				if err != nil {
					fmt.Println("Selected auth method USERNAME_PASSWORD failed")
					return false
				}

				cmd, err := msg.ParseUnamePasswordAuth(buf)
				if err != nil {
					fmt.Println("Selected auth method USERNAME_PASSWORD failed" + err.Error())
					return false
				}

				// TODO: check it in u/p storage
				if cmd.UserName != "user" || cmd.Password != "pass" {
					fmt.Println("Selected auth method USERNAME_PASSWORD failed")
					s.SendMessage(
						msg.AuthUnamePassReply{
							Ver:    cmd.Ver,
							Status: 1,
						})
					return false
				}
				s.SendMessage(
					msg.AuthUnamePassReply{
						Ver:    cmd.Ver,
						Status: 0,
					})
				fmt.Println("Selected auth method USERNAME_PASSWORD complete")
				return true
			}
		}
	}

	return func() bool {
		fmt.Println("No auth handshake method found")
		return false
	}
}

// HandleConnect ...
func (s *Session) HandleConnect(message *msg.CommandRequest) *Nif {

	fmt.Printf("Handle connect request %s\n", message)

	nif := MakeNif(s.connection, message.DstAddr, message.DstPort)

	addr, port, err := nif.Prepare()
	if err != nil {
		fmt.Println("session: failed do connect err=" + err.Error())
		cr := msg.CommandReply{
			Result:      msg.CommandResultNetworkUnreaschable,
			AddressType: msg.Ip4Address,
			BindAddress: "0.0.0.0",
			BindPort:    0,
		}
		s.SendMessage(cr)
		return nil
	}

	cr := msg.CommandReply{
		Result:      msg.CommandResultSucceeded,
		AddressType: msg.Ip4Address,
		BindAddress: addr,
		BindPort:    port,
	}

	s.SendMessage(cr)

	return nif
}
