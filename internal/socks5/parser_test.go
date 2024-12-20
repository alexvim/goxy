package socks5

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCmdCommandPasrer(t *testing.T) {
	data := []struct {
		bin  []byte
		port uint16
		dst  string
	}{
		{
			bin:  []byte{0x05, 0x01, 0x00, 0x01, 0x23, 0x32, 0x43, 0x10, 0x1F, 0x90},
			port: 8080,
			dst:  "35.50.67.16",
		},
		{
			bin:  []byte{0x05, 0x01, 0x00, 0x04, 0x23, 0x32, 0x43, 0x23, 0x23, 0x32, 0x43, 0x23, 0x23, 0x32, 0x43, 0x23, 0x23, 0x32, 0x43, 0x23, 0x1F, 0x90},
			port: 8080,
			dst:  "2332:4323:2332:4323:2332:4323:2332:4323",
		},
		{
			bin:  []byte{0x05, 0x01, 0x00, 0x03, 0x5, 0x78, 0x2e, 0x63, 0x6F, 0x6D, 0x1F, 0x90},
			port: 8080,
			dst:  "x.com",
		},
	}

	for _, test := range data {
		t.Run("Cmd Parse Test", func(t *testing.T) {
			t.Parallel()

			message, err := ParseCommand(test.bin)

			assert.NotNil(t, message)

			assert.NoError(t, err)

			assert.Equal(t, test.dst, message.DstAddr)

			assert.Equal(t, test.port, message.DstPort)
		})
	}
}
