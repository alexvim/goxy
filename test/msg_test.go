package test

import (
	"goxy/msg"
	"goxy/network"
	"testing"
)

func TestCommandReply(t *testing.T) {

	cr := msg.CommandReply{
		Result:      msg.CommandResultNetworkUnreaschable,
		AddressType: msg.IP_V4ADDRESS,
		BindAddress: "10.20.100.42",
		BindPort:    433,
	}

	b := cr.Serialize()

	s := network.BytesIp4ToString(b[4:10])

	if s != cr.BindAddress {
		t.Error(s)
	}

	var p uint16 = uint16(b[9])<<8 | uint16(b[8])
	if p != cr.BindPort {
		t.Error(p)
	}
}
