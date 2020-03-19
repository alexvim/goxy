package msg

import (
	"fmt"
	"goxy/network"
)

type CommandType uint8

const (
	ConnectCmd   CommandType = 0x01
	BindCmd                  = 0x02
	UdpAssosiate             = 0x03
)

type Atype uint8

const (
	Ip4Address Atype = 0x01
	DomainName       = 0x03
	Ip6Address       = 0x04
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
	AddressType Atype
	DstAddr     string
	DstPort     uint16
}

type CommandReply struct {
	Result      CommandResult
	AddressType Atype
	BindAddress string
	BindPort    uint16
}

// Serialize implemenet serialize to bytes
func (c CommandReply) Serialize() []byte {
	data := []byte{byte(ProtoclVersion5), byte(c.Result), byte(Reserved), byte(c.AddressType)}

	addr := network.AddressIpToBytes(c.BindAddress)
	if c.AddressType == DomainName {
		data = append(data, byte(len(c.BindAddress)))
	}
	data = append(data, addr...)
	data = append(data, byte(c.BindPort>>8), byte(c.BindPort&0x00FF))
	return data
}

func (c CommandReply) String() string {
	return fmt.Sprintf("{Cmd=%d, AddrTtype=%d, DstAddress=%v, DstPort=%d}", c.Result, c.AddressType, c.BindAddress, c.BindPort)
}

func (c CommandRequest) String() string {
	return fmt.Sprintf("{Cmd=%d, AddrTtype=%d, DstAddress=%v, DstPort=%d}", c.Command, c.AddressType, c.DstAddr, c.DstPort)
}
