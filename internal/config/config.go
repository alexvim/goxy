package config

import (
	"errors"
	"net"
)

const defaultPort = 1080

var (
	ErrInavlidProxyAddress = errors.New("invalid proxy adderess")
	ErrInavlidHostAddress  = errors.New("invalid host adderess")
)

type Config struct {
	proxyAddr string
	hostAddr  string
	dohURL    string
}

func (cfg Config) ProxyAddress() string {
	return cfg.proxyAddr
}

func (cfg Config) HostAddress() string {
	return cfg.hostAddr
}

func (cfg Config) DohURL() string {
	return cfg.dohURL
}

func (cfg Config) validate() error {
	if _, err := net.ResolveTCPAddr("tcp", cfg.proxyAddr); err != nil {
		return errors.Join(ErrInavlidProxyAddress, err)
	}

	if _, err := net.ResolveIPAddr("ip", cfg.hostAddr); err != nil {
		return errors.Join(ErrInavlidHostAddress, err)
	}

	return nil
}
