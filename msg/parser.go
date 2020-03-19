package msg

import (
	"encoding/binary"
	"fmt"
	"goxy/network"
)

func ParseAuthHandshake(buffer []byte) (*AuthRequest, error) {

	var version ProtoclVersion = ProtoclVersion(buffer[0])

	if version != ProtoclVersion5 {
		return nil, fmt.Errorf("malformed packet VER=%d", int(version))
	}

	// 2 means first octectes version and nmeth len
	if int(buffer[1]) == len(buffer[2:]) && (2+int(buffer[1])) == len(buffer) {
		message := new(AuthRequest)
		message.Methods = make([]AuthHandshkeMethod, buffer[1])
		for i, item := range buffer[2:] {
			message.Methods[i] = AuthHandshkeMethod(item)
		}
		return message, nil
	}

	return nil, fmt.Errorf("malformed auth packet")
}

func ParseUnamePasswordAuth(buf []byte) (*AuthUnamePassRequest, error) {

	var pos int = 0

	ver := int(buf[pos])
	if ver != 0x1 {
		return nil, fmt.Errorf("malformed packet ver=%d is not valid", ver)
	}

	pos = pos + 1
	// 2 means first octectes version and nmeth len
	ulen := int(buf[pos])
	// ver + ulen + user len
	if 2+ulen > len(buf) {
		return nil, fmt.Errorf("malformed packet, ulen=%d is not valid", ulen)
	}

	pos = pos + 1
	uname := string(buf[pos : pos+ulen])

	pos += ulen
	plen := int(buf[pos])
	// ver, ulen, uname, plen, pass
	if 1+1+ulen+1+plen != len(buf) {
		return nil, fmt.Errorf("malformed packet, plen=%d is not valid", 1+1+ulen+1+plen)
	}

	pos = pos + 1
	pass := string(buf[pos : pos+plen])

	return &AuthUnamePassRequest{
		Ver:      ver,
		UserName: uname,
		Password: pass,
	}, nil
}

func ParseCommand(buffer []byte) (*CommandRequest, error) {

	var pos int = 0

	var version ProtoclVersion = ProtoclVersion(buffer[pos])
	if version != ProtoclVersion5 {
		return nil, fmt.Errorf("malformed packet packet VER=%d", int(version))
	}
	pos = pos + 1

	var cmd CommandType = CommandType(buffer[pos])
	if cmd == 0 || cmd >= UdpAssosiate {
		return nil, fmt.Errorf("wrong CMD or malformed packet CMD=%d", int(cmd))
	}
	pos = pos + 1

	var rsv RSV = RSV(buffer[pos])
	if rsv != 0 {
		return nil, fmt.Errorf("wrong RSV or malformed packet RSV=%d", int(rsv))
	}
	pos = pos + 1

	var atype Atype = Atype(buffer[pos])
	var address string = ""
	var port uint16 = 0
	pos = pos + 1

	switch atype {
	case Ip4Address:

		if len(buffer[pos:]) != 4+2 {
			return nil, fmt.Errorf("malformed packet, ipv4 address and port not fit into pkt len=%d", len(buffer[4:]))
		}

		address = network.BytesIp4ToString(buffer[pos : pos+4])
		pos = pos + 4

		// port in big endian
		port = binary.BigEndian.Uint16(buffer[pos:])
		pos = pos + 2

	case Ip6Address:

		if len(buffer[pos:]) != 16+2 {
			return nil, fmt.Errorf("malformed packet, ipv6 address and port not fit into pkt len=%d", len(buffer[4:]))
		}

		address = network.BytesIp6ToString(buffer[pos : pos+16])
		pos = pos + 16

		// port in big endian
		port = binary.BigEndian.Uint16(buffer[pos:])
		pos = pos + 2

	case DomainName:

		var dnameLen int = int(buffer[pos])
		if len(buffer[pos:]) != 1+dnameLen+2 {
			return nil, fmt.Errorf("malformed packet, domain name and port not fit pkt len")
		}
		pos = pos + 1

		address = string(buffer[pos : pos+dnameLen])
		pos += dnameLen

		// port in big endian
		port = binary.BigEndian.Uint16(buffer[pos:])
		pos = pos + 2

	default:
		return nil, fmt.Errorf("wrong ATYP or malformed packet ATYP=%d", atype)
	}

	if pos != len(buffer) {
		return nil, fmt.Errorf("malformed packet pos=%d, expected=%d", pos, len(buffer))
	}

	message := new(CommandRequest)

	message.Command = cmd
	message.AddressType = atype
	message.DstAddr = address
	message.DstPort = port

	return message, nil
}
