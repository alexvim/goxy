package socks5

import (
	"goxy/internal/netutils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandReply(t *testing.T) {
	cr := CommandReply{
		Rep:     NetworkUnreachable,
		ATyp:    AddressTypeIP4,
		BndAddr: "10.20.100.42",
		BndPort: 433,
	}

	b := cr.Serialize()

	s := netutils.IpToString(b[4:8])
	assert.Equal(t, cr.BndAddr, s, "binary %+v", b)

	p := uint16(b[8])<<8 | uint16(b[9])
	assert.Equal(t, cr.BndPort, p)
}
