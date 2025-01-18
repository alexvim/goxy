package netutils

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

const (
	AddressTypeIP4 AddressType = iota
	AddressTypeIP6
)

type AddressType uint8

func DiscoveryIfaceToBind(atype AddressType) (string, error) {
	ifaces, _ := net.Interfaces()

	// handle err
	for _, i := range ifaces {
		if i.Flags&(net.FlagUp|^net.FlagLoopback|^net.FlagPointToPoint) == 0 {
			continue
		}

		addresses, err := i.Addrs()
		if err != nil {
			return "", errors.New("failed to get address")
		}

		for _, address := range addresses {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil && atype == AddressTypeIP4 {
					return ipnet.IP.String(), nil
				}

				if ipnet.IP.To16() != nil && atype == AddressTypeIP6 {
					return ipnet.IP.String(), nil
				}
			}
		}
	}

	return "", errors.New("failed to get address")
}

func GetTcpAddrType(address net.Addr) AddressType {
	ip, _ := address.(*net.TCPAddr)
	if ip.IP.To4() != nil {
		return AddressTypeIP4
	}

	return AddressTypeIP6
}

func IpToString(ip []byte) string {
	if len(ip) == 4 {
		return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
	}

	if len(ip) == 16 {
		return fmt.Sprintf("%X:%X:%X:%X:%X:%X:%X:%X",
			binary.BigEndian.Uint16(ip[0:2]),
			binary.BigEndian.Uint16(ip[2:4]),
			binary.BigEndian.Uint16(ip[4:6]),
			binary.BigEndian.Uint16(ip[6:8]),
			binary.BigEndian.Uint16(ip[8:10]),
			binary.BigEndian.Uint16(ip[10:12]),
			binary.BigEndian.Uint16(ip[12:14]),
			binary.BigEndian.Uint16(ip[14:16]),
		)
	}

	return ""
}

func StringToIp(addr string) []byte {
	ip := net.ParseIP(addr)

	if v4 := ip.To4(); v4 != nil {
		return v4
	}

	if v6 := ip.To16(); v6 != nil {
		return v6
	}

	return nil
}

func (addrtype AddressType) String() string {
	switch addrtype {
	case AddressTypeIP4:
		return "tcp4"
	case AddressTypeIP6:
		return "tcp6"
	default:
		return "unknown"
	}
}
