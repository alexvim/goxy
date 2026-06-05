package configreader

import (
	"strings"
	"testing"

	_ "embed"

	"github.com/stretchr/testify/assert"
)

//go:embed data/config.default.json
var testData string

func TestJsonConfigReader(t *testing.T) {
	t.Run("pasre config", func(t *testing.T) {
		jsonConfig, err := readFromJson(strings.NewReader(testData))
		assert.NoError(t, err)

		assert.Equal(t, "192.168.0.1", jsonConfig.Proxy.Address)
		assert.Equal(t, 1080, jsonConfig.Proxy.Port)

		assert.Equal(t, "10.0.0.1", jsonConfig.Host.Address)

		assert.Equal(t, "https://doh.opendns.com/dns-query", jsonConfig.Resolver.DoH.URL)
	})
}
