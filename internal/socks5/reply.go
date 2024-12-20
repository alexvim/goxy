package socks5

import (
	"encoding/binary"
	"fmt"
	"goxy/internal/netutils"
)

// The SOCKS request information is sent by the client as soon as it has
// established a connection to the SOCKS server, and completed the
// authentication negotiations.
type CommandReply struct {
	// VER protocol version: X'05'
	ver ProtoclVersion
	// REP Reply field
	Rep Rep
	// RSV RESERVED
	rsv Reserved
	// ATYP address type of following address
	ATyp AddressType
	// BND.ADDR	server bound address
	BndAddr string
	// BND.PORT server bound port in network octet order
	BndPort uint16
}

type AuthMethodReply struct {
	ver    ProtoclVersion
	Method AuthHandshakeMethod
}

type AuthUnamePassReply struct {
	Ver    AuthVer
	Status AuthStatus
}

func newReply(rep Rep, atyp AddressType, addr string, port uint16) CommandReply {
	return CommandReply{
		ver:     ProtoclVersion5,
		Rep:     rep,
		rsv:     ReservedDefault,
		ATyp:    atyp,
		BndAddr: addr,
		BndPort: port,
	}
}

func newAuthMethodReply(method AuthHandshakeMethod) AuthMethodReply {
	return AuthMethodReply{
		ver:    ProtoclVersion5,
		Method: method,
	}
}

func newAuthUserPassReply(ver AuthVer, status AuthStatus) AuthUnamePassReply {
	return AuthUnamePassReply{
		Ver:    ver,
		Status: status,
	}
}

// CommandReply
func (c CommandReply) Serialize() []byte {
	data := []byte{
		byte(c.ver),
		byte(c.Rep),
		byte(c.rsv),
		byte(c.ATyp),
	}

	addr := netutils.StringToIp(c.BndAddr)
	if c.ATyp == AddressTypeDomainName {
		data = append(data, byte(len(c.BndAddr)))
	}

	data = append(data, addr...)

	// port should be hton-ed
	data = binary.BigEndian.AppendUint16(data, c.BndPort)

	return data
}

func (c CommandReply) String() string {
	return fmt.Sprintf("command_reply: cmd=%s, type=%s, dst=%s, port=%d", c.Rep, c.ATyp, c.BndAddr, c.BndPort)
}

// AuthReply
func (m AuthMethodReply) Serialize() []byte {
	return []byte{
		byte(m.ver),
		byte(m.Method),
	}
}

func (m AuthMethodReply) String() string {
	return fmt.Sprintf("auth_reply: method=%s", m.Method)
}

// AuthUnamePassReply
func (ar AuthUnamePassReply) Serialize() []byte {
	return []byte{
		byte(ar.Ver),
		byte(ar.Status),
	}
}

func (m AuthUnamePassReply) String() string {
	return fmt.Sprintf("auth_user/pass_reply: ver=%s, status=%s", m.Ver, m.Status)
}
