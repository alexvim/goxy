package msg

import (
	"encoding/binary"
	"errors"
	"fmt"
	"goxy/network"
	"strconv"
)

func ParseAuthHandshake(buffer []byte) (*AuthRequest, error) {

	var version ProtoclVersion = ProtoclVersion(buffer[0])

	if version != ProtoclVersion5 {
		return nil, errors.New("wrong protocol version or malformed packet  =" + strconv.Itoa(int(version)))
	}

	// 2 means first octectes version and nmeth len
	if int(buffer[1]) == len(buffer[2:]) && (2+int(buffer[1])) == len(buffer) {
		message := new(AuthRequest)
		message.Methods = make([]AuthMethod, buffer[1])
		for i, item := range buffer[2:] {
			message.Methods[i] = AuthMethod(item)
		}
		return message, nil
	}

	return nil, errors.New("malformed auth packet")
}

func ParseUnamePasswordAuth(buffer []byte) (*AuthRequest, error) {

	if int(buffer[0]) != 0x1 {
		return nil, errors.New("wrong protocol version or malformed packet  =" + strconv.Itoa(int(buffer[0])))
	}

	// 2 means first octectes version and nmeth len
	if int(buffer[1]) == len(buffer[2:]) && (2+int(buffer[1])) == len(buffer) {
		message := new(AuthRequest)
		message.Methods = make([]AuthMethod, buffer[1])
		for i, item := range buffer[2:] {
			message.Methods[i] = AuthMethod(item)
		}
		return message, nil
	}
}

func ParseCommand(buffer []byte) (*CommandRequest, error) {

	var pos int = 0

	var version ProtoclVersion = ProtoclVersion(buffer[pos])
	if version != ProtoclVersion5 {
		return nil, errors.New("wrong protocol version or malformed command packet VER=" + strconv.Itoa(int(version)))
	}
	pos += 1

	var cmd CommandType = CommandType(buffer[pos])
	if cmd == 0 || cmd >= UDP_ASSOCIATE {
		return nil, errors.New("wrong CMD or malformed packet CMD=" + strconv.Itoa(int(cmd)))
	}
	pos += 1

	var rsv RSV = RSV(buffer[pos])
	if rsv != 0 {
		return nil, errors.New("wrong RSV or malformed packet RSV=" + strconv.Itoa(int(rsv)))
	}
	pos += 1

	var atype ATYP = ATYP(buffer[pos])
	var address string = ""
	var port uint16 = 0
	pos += 1

	switch atype {
	case IP_V4ADDRESS:

		if len(buffer[pos:]) != 4+2 {
			return nil, errors.New("malformed packet, ipv4 address and port not fit into pkt len=" + strconv.Itoa(len(buffer[4:])))
		}

		address = network.BytesIp4ToString(buffer[pos : pos+4])
		pos += 4

		// port in big endian
		port = binary.BigEndian.Uint16(buffer[pos:])
		pos += 2

	case IP_V6ADDRESS:

		if len(buffer[pos:]) != 16+2 {
			return nil, errors.New("malformed packet, ipv6 address and port not fit into pkt len=" + strconv.Itoa(len(buffer[4:])))
		}

		address = network.BytesIp6ToString(buffer[pos : pos+16])
		pos += 16

		// port in big endian
		port = binary.BigEndian.Uint16(buffer[pos:])
		pos += 2

	case DOMAINNAME:

		var dnameLen int = int(buffer[pos])
		if len(buffer[pos:]) != 1+dnameLen+2 {
			return nil, errors.New("malformed packet, domain name and port not fit pkt len")
		}
		pos += 1

		address = string(buffer[pos : pos+dnameLen])
		pos += dnameLen

		// port in big endian
		port = binary.BigEndian.Uint16(buffer[pos:])
		pos += 2

	default:
		return nil, errors.New("wrong ATYP or malformed packet ATYP=" + strconv.Itoa(int(atype)))
	}

	if pos != len(buffer) {
		return nil, errors.New(fmt.Sprintf("malformed packet pos=%d, expected=%d", pos, len(buffer)))
	}

	message := new(CommandRequest)

	message.Command = cmd
	message.AddressType = atype
	message.DstAddr = address
	message.DstPort = port

	return message, nil
}
