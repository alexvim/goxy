package msg

import (
	"fmt"
	"goxy/network"
)

type CommandType uint8

const (
	CONNECT       CommandType = 0x01
	BIND                      = 0x02
	UDP_ASSOCIATE             = 0x03
)

type ATYP uint8

const (
	IP_V4ADDRESS ATYP = 0x01
	DOMAINNAME        = 0X03
	IP_V6ADDRESS      = 0x04
)

type CommandResult uint8

const (
	CommandResultSucceeded                     CommandResult = 0x00
	CommandResultGeneralSocksServerFailure                   = 0x01
	CommandResultConnectionNotAllowedByRuleset               = 0x02
	CommandResultNetworkUnreaschable                         = 0x03
	CommandResultHostUnreachable                             = 0x04
	CommandResultConnectionRefused                           = 0x05
	CommandResultTtlExpired                                  = 0x06
	CommandResultCommandNotSupported                         = 0x07
	CommandResultAddressTypeNotSupport                       = 0x08
)

type CommandRequest struct {
	Command     CommandType
	AddressType ATYP
	DstAddr     string
	DstPort     uint16
}

type CommandReply struct {
	Result      CommandResult
	AddressType ATYP
	BindAddress string
	BindPort    uint16
}

func (cr *CommandRequest) GetType() MessageType {
	return CmdReq
}

func (cr CommandReply) Serialize() []byte {
	data := []byte{byte(ProtoclVersion5), byte(cr.Result), byte(Reserved), byte(cr.AddressType)}

	addr := network.AddressIpToBytes(cr.BindAddress)
	if cr.AddressType == DOMAINNAME {
		data = append(data, byte(len(cr.BindAddress)))
	}
	data = append(data, addr...)
	data = append(data, byte(cr.BindPort>>8), byte(cr.BindPort&0x00FF))
	return data
}

func (c CommandRequest) String() string {
	return fmt.Sprintf("{Cmd=%d, AT=%d, DstAddress=%v, DstPort=%d}", c.Command, c.AddressType, c.DstAddr, c.DstPort)
}

func (c CommandReply) String() string {
	return fmt.Sprintf("{Cmd=%d, AT=%d, DstAddress=%v, DstPort=%d}", c.Result, c.AddressType, c.BindAddress, c.BindPort)
}
