package config

import (
	"errors"
	"fmt"
	configreader "goxy/internal/config/internal"
	"goxy/internal/netutils"
	"log"
	"net"
)

const defaultPort = 1080

var (
	ErrInavlidAddress = errors.New("invalid adderess")
)

type Config struct {
	proxyAddr string
	localAddr string
}

func ReadFromArgs(args []string) (Config, error) {
	reader := configreader.Read(args)

	cfg := Config{
		proxyAddr: reader.ProxyAddr,
		localAddr: reader.LocalAddr,
	}

	if len(reader.LocalAddr) == 0 && len(reader.ProxyAddr) == 0 {
		localAddr4, err := netutils.DiscoveryIfaceToBind(netutils.AddressTypeIP4)
		if err != nil {
			log.Printf("server: failed to get local net interface err: %s", err)

			return Config{}, err
		}

		cfg = Config{
			proxyAddr: fmt.Sprintf("%s:%d", localAddr4, defaultPort),
			localAddr: localAddr4,
		}
	}

	return cfg, cfg.validate()
}

func (cfg Config) ProxyAddress() string {
	return cfg.proxyAddr
}

func (cfg Config) LocalAddress() string {
	return cfg.localAddr
}

func (cfg Config) validate() error {
	if _, err := net.ResolveTCPAddr("tcp", cfg.proxyAddr); err != nil {
		return errors.Join(ErrInavlidAddress, err)
	}

	if _, err := net.ResolveIPAddr("ip", cfg.localAddr); err != nil {
		return errors.Join(ErrInavlidAddress, err)
	}

	return nil
}
