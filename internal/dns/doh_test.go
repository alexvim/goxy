package dns

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testDomain    = "www.example.com"
	testServerURL = "dnsserver.example.net/dns-query"
)

func TestDoh(t *testing.T) {

	t.Run("wire_build", func(t *testing.T) {
		const result = "00000100000100000000000003777777076578616d706c6503636f6d0000010001"

		b, err := makeDnsWireQuery(testDomain, ResolveTypeIPv4)

		assert.NoError(t, err)

		assert.Equal(t, result, hex.EncodeToString(b))
	})

	t.Run("wire_parse", func(t *testing.T) {
		const msg = "00008180000100010000000003777777076578616d706c6503636f6d00001c0001c00c001c000100000e7d001020010db8abcd00120001000200030004"

		b, _ := hex.DecodeString(msg)

		_, v6, err := parseDnsWireQuery(b)

		assert.NoError(t, err)

		assert.NotEmpty(t, v6)
	})

	t.Run("request", func(t *testing.T) {
		_, _ = newDoHResolver(testServerURL, ResolveTypeIPv4)
	})

	t.Run("real_request", func(t *testing.T) {
		t.Skip()

		d, err := newDoHResolver("dns.google/dns-query", ResolveTypeIPv4)

		assert.NoError(t, err)

		v4, err := d.Resolve("google.com")

		assert.NoError(t, err)

		assert.NotEmpty(t, v4)
	})
}
