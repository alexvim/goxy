package config

import (
	"fmt"
	configreader "goxy/internal/config/internal"
	"goxy/internal/netutils"
	"log"
	"net"
	"strconv"
)

func LoadConfig(args []string) (Config, error) {
	log.Printf("config: create config usgin command line args %v", args)

	cfg := Config{}

	if appargs := configreader.ParseArgs(args); appargs.HasConfigFile() {
		cfg = fromJson(configreader.ReadFromFile(appargs.ConfigFile()))
	} else {
		cfg = fromArgs(appargs)
	}

	if err := cfg.validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func fromArgs(agrs *configreader.CommandLineAgrs) Config {
	log.Println("config: load config from command line args")

	cfg := Config{
		proxyAddr: agrs.ProxyAddr(),
		hostAddr:  agrs.LocalAddr(),
		dohURL:    agrs.DohURL(),
	}

	if len(cfg.hostAddr) == 0 && len(cfg.proxyAddr) == 0 {
		localAddr4, err := netutils.DiscoveryIfaceToBind(netutils.AddressTypeIP4)
		if err != nil {
			log.Printf("server: failed to get local net interface err: %s", err)
			return Config{}
		}

		cfg.proxyAddr = fmt.Sprintf("%s:%d", localAddr4, defaultPort)
		cfg.hostAddr = localAddr4
	}

	return cfg
}

func fromJson(jsonConf configreader.JsonConfig) Config {
	log.Printf("config: load config from json config %v", jsonConf)

	cfg := Config{
		proxyAddr: net.JoinHostPort(jsonConf.Proxy.Address, strconv.Itoa(jsonConf.Proxy.Port)),
		hostAddr:  jsonConf.Host.Address,
		dohURL:    jsonConf.Resolver.DoH.URL,
	}

	if len(cfg.proxyAddr) == 0 && len(cfg.hostAddr) == 0 {
		localAddr4, err := netutils.DiscoveryIfaceToBind(netutils.AddressTypeIP4)
		if err != nil {
			log.Printf("server: failed to get local net interface err: %s", err)
			return Config{}
		}

		cfg.proxyAddr = fmt.Sprintf("%s:%d", localAddr4, defaultPort)
		cfg.hostAddr = localAddr4
	}

	return cfg
}
