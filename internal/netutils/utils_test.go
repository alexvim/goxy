package netutils

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIp4ToString(t *testing.T) {
	data := []struct {
		bin []byte
		str string
	}{
		{
			bin: []byte{0x23, 0x32, 0x43, 0x10},
			str: "35.50.67.16",
		},
		{
			bin: []byte{0x0A, 0x14, 0x64, 0x2A},
			str: "10.20.100.42",
		},
	}

	for _, d := range data {
		t.Run("test_convert", func(t *testing.T) {
			assert.Equal(t, IpToString(d.bin), d.str)
		})
	}
}

func TestIp6Convert(t *testing.T) {
	ip := []byte{0x23, 0x32, 0x43, 0x23, 0x23, 0x32, 0x43, 0x23, 0x23, 0x32, 0x43, 0x23, 0x23, 0x32, 0x43, 0x23}

	if len(ip) != 16 {
		panic(0)
	}

	if s := IpToString(ip); s != "2332:4323:2332:4323:2332:4323:2332:4323" {
		t.Error(s)
	}
}

func TestGetLocalInetAddress(t *testing.T) {
	t.Skip()

	a, err := DiscoveryIfaceToBind(AddressTypeIP4)

	if err != nil {
		t.Error("Failed to find " + err.Error())
	}

	if net.ParseIP(a).To4() == nil {
		t.Error("Invalid address " + a)
	}
}
