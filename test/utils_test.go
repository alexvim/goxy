package test

import (
	"goxy/network"
	"net"
	"testing"
)

func TestGetLocalInetAddress(t *testing.T) {

	a, err := network.GetLocalInetAddress(network.AddrTypeIpv4)

	if err != nil {
		t.Error("Failed to find " + err.Error())
	}

	if net.ParseIP(a).To4() == nil {
		t.Error("Invalid address " + a)
	}
}

func TestIp4Convert(t *testing.T) {

	ip := []byte{0x23, 0x32, 0x43, 0x10}

	if s := network.BytesIp4ToString(ip); s != "35.50.67.16" {
		t.Error(s)
	}
}

func TestIp6Convert(t *testing.T) {

	ip := []byte{0x23, 0x32, 0x43, 0x23, 0x23, 0x32, 0x43, 0x23, 0x23, 0x32, 0x43, 0x23, 0x23, 0x32, 0x43, 0x23}

	if len(ip) != 16 {
		panic(0)
	}

	if s := network.BytesIp6ToString(ip); s != "2332:4323:2332:4323:2332:4323:2332:4323" {
		t.Error(s)
	}
}
