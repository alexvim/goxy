package socks5

import (
	"encoding/binary"
	"goxy/internal/netutils"
	"log"
)

func ParseCommand(buffer []byte) (*CommandRequest, error) {
	pos := uint32(0)

	if protocol := ProtoclVersion(buffer[0]); protocol != ProtoclVersion5 {
		log.Printf("unsupported protocol, ver=%s\n", protocol)
		return nil, ErrUnsupportedProtocolVersion
	}

	pos = pos + 1

	cmd := Command(buffer[pos])
	if cmd <= Undefined || cmd >= Unknown {
		log.Printf("malformed packet cmd=%s\n", cmd)
		return nil, ErrMalformedPacket
	}

	pos = pos + 1

	if reserved := Reserved(buffer[pos]); reserved != ReservedDefault {
		log.Printf("malformed packet reserved=%s\n", reserved)
		return nil, ErrMalformedPacket
	}

	pos = pos + 1

	var (
		atype   AddressType = AddressType(buffer[pos])
		address string
		port    uint16
	)

	pos = pos + 1

	switch atype {
	case AddressTypeIP4:
		if len(buffer[pos:]) != 4+2 {
			log.Printf("ipv4 address and port does not fit into pkt len=%d\n", len(buffer[4:]))
			return nil, ErrMalformedPacket
		}

		address = netutils.IpToString(buffer[pos : pos+4])

		pos = pos + 4
	case AddressTypeDomainName:
		dnameLen := uint32(buffer[pos])
		if len(buffer[pos:]) != int(1+dnameLen+2) {
			log.Printf("domain and port does not fit into pkt len=%d\n", len(buffer[4:]))
			return nil, ErrMalformedPacket
		}

		pos = pos + 1

		address = string(buffer[pos : pos+dnameLen])

		pos += dnameLen
	case AddressTypeIP6:
		if len(buffer[pos:]) != 16+2 {
			log.Printf("ipv6 address and port does not fit into pkt len=%d\n", len(buffer[4:]))
			return nil, ErrMalformedPacket
		}

		address = netutils.IpToString(buffer[pos : pos+16])

		pos = pos + 16
	default:
		log.Printf("malformed packet atype=%s\n", atype)
		return nil, ErrMalformedPacket
	}

	// port in big endian
	port = binary.BigEndian.Uint16(buffer[pos:])

	pos = pos + 2

	if len(buffer) != int(pos) {
		log.Printf("malformed packet, len and parsed bytes are not equal len=%d, parsed=%d\n", len(buffer), pos)
		return nil, ErrMalformedPacket
	}

	return &CommandRequest{Cmd: cmd, AddressType: atype, DstAddr: address, DstPort: port}, nil
}

func ParseAuthHandshake(buffer []byte) (*AuthRequest, error) {
	if protocol := ProtoclVersion(buffer[0]); protocol != ProtoclVersion5 {
		log.Printf("unsupported protocol, ver=%s\n", protocol)
		return nil, ErrUnsupportedProtocolVersion
	}

	// 2 means first octectes version and nmethods len
	if int(buffer[1]) != len(buffer[2:]) || (2+int(buffer[1])) != len(buffer) {
		log.Printf("malformed auth packet\n")
		return nil, ErrMalformedPacket
	}

	methdos := make([]AuthHandshakeMethod, buffer[1])

	// scan all mnethods
	for i, item := range buffer[2:] {
		methdos[i] = AuthHandshakeMethod(item)
	}

	return &AuthRequest{Methods: methdos}, nil
}

func ParseUnamePasswordAuth(buf []byte) (*AuthUserPassRequest, error) {
	pos := uint32(0)

	ver := AuthVer(buf[pos])
	if ver != AuthVer0x01 {
		log.Printf("malformed packet, auth ver is not valid ver=%d\n", ver)
		return nil, ErrMalformedPacket
	}

	pos = pos + 1

	// 2 means first octectes version and nmeth len
	ulen := uint32(buf[pos])

	// ver + ulen + user len
	if int(2+ulen) > len(buf) {
		log.Printf("malformed packet, ulen=%d is not valid\n", ulen)
		return nil, ErrMalformedPacket
	}

	pos = pos + 1

	uname := string(buf[pos : pos+ulen])

	pos += ulen

	plen := uint32(buf[pos])

	// ver, ulen, uname, plen, pass
	if int(1+1+ulen+1+plen) != len(buf) {
		log.Printf("malformed packet, plen=%d is not valid\n", 1+1+ulen+1+plen)
		return nil, ErrMalformedPacket
	}

	pos = pos + 1

	password := string(buf[pos : pos+plen])

	return &AuthUserPassRequest{Ver: ver, User: uname, Password: password}, nil
}
