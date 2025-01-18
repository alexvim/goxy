package socks5

import (
	"errors"
	"fmt"
	"io"
	"log"
)

const (
	defaultBufferSize = 4096
	addressNull       = "0.0.0.0"
	portNul           = 0
)

type DoConnect func(addr string, port uint16) (string, uint16, error)

type Flow struct {
	uuid      string
	conn      io.ReadWriter
	buffer    []byte
	doConnect DoConnect
}

func NewFlow(uuid string, conn io.ReadWriter, doConnect DoConnect) Flow {
	return Flow{
		uuid:      uuid,
		conn:      conn,
		buffer:    make([]byte, defaultBufferSize),
		doConnect: doConnect,
	}
}

func (flow Flow) Run() error {
	log.Printf("%s: start protocol exchange\n", flow)

	// auth methods
	buf, err := flow.readMessage()
	if err != nil {
		log.Printf("%s: failed to read message err=%s\n", flow, err)
		return errors.Join(ErrReadMessage, err)
	}

	auth, err := ParseAuthHandshake(buf)
	if err != nil {
		log.Printf("%s: failed to parse message err=%s\n", flow, err)
		return errors.Join(ErrParseMessage, err)
	}

	// error check here
	if authRoutine := flow.handleAuth(auth); !authRoutine() {
		log.Printf("%s: failed to do auth\n", flow)
		return ErrAuthenticate
	}

	// command
	buf, err = flow.readMessage()
	if err != nil {
		log.Printf("%s: failed to read message err=%s\n", flow, err)
		return errors.Join(ErrReadMessage, err)
	}

	cmd, err := ParseCommand(buf)
	if err != nil {
		log.Printf("%s: failed to pasre buffer err=%s\n", flow, err)
		return errors.Join(ErrParseMessage, err)
	}

	switch cmd.Cmd {
	case Connect:
		if err := flow.handleConnect(cmd); err != nil {
			log.Printf("%s: failed to send connect err=%s\n", flow, err)
			return errors.Join(ErrSocksConnect, err)
		}
	default:
		if err := flow.handleUnsupported(cmd); err != nil {
			log.Printf("%s: failed to send reply err=%s\n", flow, err)
			return errors.Join(ErrSocksConnect, err)
		}
	}

	log.Printf("%s: complete with success\n", flow)

	return nil
}

func (flow Flow) handleAuth(request *AuthRequest) func() bool {
	log.Printf("%s: handle protocol flow\n", flow)

	for _, method := range request.Methods {
		switch method {
		case AuthMethodNoAuthRequired:
			flow.sendMessage(newAuthMethodReply(AuthMethodNoAuthRequired))

			return func() bool {
				log.Printf("%s: execute authenticate method: %s\n", flow, method)
				return true
			}
		case AuthMethodUserAndPassword:
			flow.sendMessage(newAuthMethodReply(AuthMethodUserAndPassword))

			return func() bool {
				log.Printf("%s: execute authenticate method: %s\n", flow, method)

				buf, err := flow.readMessage()
				if err != nil {
					log.Printf("%s: selected auth method USERNAME_PASSWORD failed err=%s\n", flow, err)
					return false
				}

				cmd, err := ParseUnamePasswordAuth(buf)
				if err != nil {
					log.Printf("%s: selected auth method USERNAME_PASSWORD failed err=%s\n", flow, err)
					return false
				}

				// TODO: check it in u/p storage
				if cmd.User != "user" || cmd.Password != "pass" {
					log.Printf("%s: selected auth method USERNAME_PASSWORD failed\n", flow)

					flow.sendMessage(newAuthUserPassReply(cmd.Ver, AuthStatusFailure))

					return false
				}

				flow.sendMessage(newAuthUserPassReply(cmd.Ver, AuthStatusSuccess))

				log.Printf("%s: selected auth method USERNAME_PASSWORD complete\n", flow)

				return true
			}
		}
	}

	flow.sendMessage(newAuthMethodReply(AuthMethodNoAcceptableMethods))

	return func() bool {
		log.Printf("%s: no auth handshake method found\n", flow)
		return false
	}
}

func (flow Flow) handleConnect(msg *CommandRequest) error {
	log.Printf("%s: handle connect command to %s:%d\n", flow, msg.DstAddr, msg.DstPort)

	addr, port, err := flow.doConnect(msg.DstAddr, msg.DstPort)
	if err != nil {
		log.Printf("%s: failed do connect err=%s\n", flow, err)

		reply := newReply(NetworkUnreachable, AddressTypeIP4, addressNull, portNul)

		flow.sendMessage(reply)

		return err
	}

	reply := newReply(Succeeded, AddressTypeIP4, addr, port)

	return flow.sendMessage(reply)
}

func (flow Flow) handleUnsupported(msg *CommandRequest) error {
	log.Printf("%s: unsupported command %s\n", flow, msg)

	reply := newReply(CommandNotSupported, AddressTypeIP4, addressNull, portNul)

	return flow.sendMessage(reply)
}

func (flow Flow) sendMessage(msg any) error {
	serializable, ok := msg.(interface{ Serialize() []byte })
	if !ok {
		panic("flow: failed to send message, no serialize interface")
	}

	if _, err := flow.conn.Write(serializable.Serialize()); err != nil {
		log.Printf("%s: failed to send message err=%s\n", flow, err)
		return err
	}

	return nil
}

func (flow Flow) readMessage() ([]byte, error) {
	n, err := flow.conn.Read(flow.buffer)
	if err != nil || n <= 0 {
		log.Printf("%s: failed to read message err=%s\n", flow, err)
		return nil, ErrSendMessage
	}

	return flow.buffer[0:n], nil
}

func (flow Flow) String() string {
	return fmt.Sprintf("flow[%s]", flow.uuid)
}
