package socks5

import (
	"fmt"
)

type CommandRequest struct {
	// VER protocol version: X'05'
	Ver ProtoclVersion
	// CMD
	Cmd Command
	// RSV RESERVED
	Rsv Reserved
	// ATYP address type of following address
	AddressType AddressType
	// DST.ADDR desired destination address
	DstAddr string
	// DST.PORT desired destination port
	DstPort uint16
}

type AuthRequest struct {
	// VER protocol version: X'05'
	Ver ProtoclVersion
	// METHODS
	Methods []AuthHandshakeMethod
}

type AuthUserPassRequest struct {
	// VER protocol version: X'05'
	Ver AuthVer
	// The UNAME field contains the username as known to the source operating system
	User string
	// The PASSWD field contains the password association with the given UNAME.
	Password string
}

func (req CommandRequest) String() string {
	return fmt.Sprintf("command_request: cmd=%s, type=%s, dst=%s, port=%d", req.Cmd, req.AddressType, req.DstAddr, req.DstPort)
}

func (req AuthRequest) String() string {
	return fmt.Sprintf("auth_request: num=%d, methods=%v", len(req.Methods), req.Methods)
}

func (req AuthUserPassRequest) String() string {
	return fmt.Sprintf("auth_user/pass_request: ver=%d, user=%s, password=%s", req.Ver, req.User, req.Password)
}
