package sys

import (
	"runtime"
)

var Config config

type config struct {
	ProxyConnectionAddress string
	ProxyExternalAddress   string // could be list
}

func (c *config) Read() {

	if runtime.GOOS == "windows" {
		c.ReadFromEnv()
	}

	if runtime.GOOS == "linux" {
		c.ReadFromJson()
	}
}

func (c *config) ReadFromJson() {

}

func (c *config) ReadFromEnv() {
	/*
		acon := os.Getenv("GOXY_PROXY_CONN_ADDRESS")
		aext := os.Getenv("GOXY_PROXY_EXT_ADDRESS")

		c.ProxyConnectionAddress = net.ParseIP(acon).
		c.ProxyExternalAddress = net.ParseIP(aext)
	*/
}
