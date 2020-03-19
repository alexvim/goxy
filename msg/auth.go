package msg

import "fmt"

// AuthHandshkeMethod ...
type AuthHandshkeMethod uint8

const (
	NO_AUTHENTICATION_REQUIRED AuthHandshkeMethod = 0x00
	GSSAPI                                        = 0x01
	USERNAME_PASSWORD                             = 0x02
	NO_ACCEPTABLE_METHODS                         = 0xFF
)

type AuthRequest struct {
	Methods []AuthHandshkeMethod
}

type AuthReply struct {
	Method AuthHandshkeMethod
}

type AuthUnamePassRequest struct {
	Ver      int
	UserName string
	Password string
}

type AuthUnamePassReply struct {
	Ver    int
	Status int
}

func (m AuthReply) Serialize() []byte {
	return []byte{byte(ProtoclVersion5), byte(m.Method)}
}

func (m AuthRequest) String() string {
	return fmt.Sprintf("{NumOfMethods=%d, Methods=%v}", len(m.Methods), m.Methods)
}

func (m AuthReply) String() string {
	return fmt.Sprintf("{Method=%v}", m.Method)
}

func (ar AuthUnamePassReply) Serialize() []byte {
	return []byte{byte(ar.Ver), byte(ar.Status)}
}

func (m AuthUnamePassRequest) String() string {
	return fmt.Sprintf("{Ver=%d, UserName=%s, Passwd=%s}", m.Ver, m.UserName, m.Password)
}

func (m AuthUnamePassReply) String() string {
	return fmt.Sprintf("{Ver=%d, Status=%d}", m.Ver, m.Status)
}
