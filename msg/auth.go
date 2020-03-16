package msg

import "fmt"

type AuthMethod uint8

const (
	NO_AUTHENTICATION_REQUIRED AuthMethod = 0x00
	GSSAPI                                = 0x01
	USERNAME_PASSWORD                     = 0x02
	NO_ACCEPTABLE_METHODS                 = 0xFF
)

type AuthRequest struct {
	Methods []AuthMethod
}

type AuthReply struct {
	Method AuthMethod
}

func (ar AuthReply) Serialize() []byte {
	return []byte{byte(ProtoclVersion5), byte(ar.Method)}
}

func (m AuthRequest) String() string {
	return fmt.Sprintf("{NumOfMethods=%d, Methods=%v}", len(m.Methods), m.Methods)
}

func (m AuthReply) String() string {
	return fmt.Sprintf("{Method=%v}", m.Method)
}
