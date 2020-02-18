package network

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

type AddrType uint8

const (
	AddrTypeIpv4 AddrType = 0x0
	AddrTypeIpv6          = 0x1
)

func GetLocalInetAddress(atype AddrType) (string, error) {

	var ip net.IP = nil

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
				if ipnet.IP.To4() != nil {
					ip = ipnet.IP
				}
			}
		}

		if ip != nil {
			break
		}
	}

	return ip.String(), nil
}

func BytesIp4ToString(ip []byte) string {
	return fmt.Sprintf("%v.%v.%v.%v", int(ip[0]), int(ip[1]), int(ip[2]), int(ip[3]))
}

func BytesIp6ToString(ip []byte) string {
	return fmt.Sprintf("%X:%X:%X:%X:%X:%X:%X:%X",
		binary.BigEndian.Uint16(ip[0:2]),
		binary.BigEndian.Uint16(ip[2:4]),
		binary.BigEndian.Uint16(ip[4:6]),
		binary.BigEndian.Uint16(ip[6:8]),
		binary.BigEndian.Uint16(ip[8:10]),
		binary.BigEndian.Uint16(ip[10:12]),
		binary.BigEndian.Uint16(ip[12:14]),
		binary.BigEndian.Uint16(ip[14:16]))
}

func AddressIpToBytes(ip string) []byte {
	return net.ParseIP(ip).To4()
}
