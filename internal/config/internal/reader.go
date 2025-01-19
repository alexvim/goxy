package configreader

import (
	"encoding/json"
	"errors"
	"flag"
	"io"
	"log"
	"os"
)

var (
	cmdArgsDefaults = commandLineAgrs{
		proxyAddr: cmdArg{
			flag:  "p",
			descr: "socks5 proxy and port in format ip:port to bind",
		},
		localAddr: cmdArg{
			flag:  "l",
			descr: "target network ip address",
		},
		dohURL: cmdArg{
			flag:  "doh",
			descr: "DoH URL",
		},
		jsonConfig: cmdArg{
			flag:  "c",
			descr: "full path to json config location",
		},
	}

	ErrInvalidCmdArgs     = errors.New("invalid command line args")
	ErrConfigFileNotFound = errors.New("config file not found")
	ErrConfigFileParse    = errors.New("config file parse error")
	ErrConfigAddress      = errors.New("invalid address")
)

type Config struct {
	ProxyAddr string `json:"proxy_address"`
	LocalAddr string `json:"local_address"`
	DohURL    string `json:"doh_url"`
}

type cmdArg struct {
	val   string
	flag  string
	descr string
}

type commandLineAgrs struct {
	proxyAddr  cmdArg
	localAddr  cmdArg
	dohURL     cmdArg
	jsonConfig cmdArg
}

func Read(args []string) Config {
	cargs := cmdArgsDefaults

	cmdArgs := &cargs

	cmdArgs.parse(args)

	if len(cmdArgs.jsonConfig.val) > 0 {
		file, err := os.Open(cmdArgs.jsonConfig.val)
		if err != nil {
			log.Printf("cfgreader: failed to open %s err=%s\n", cmdArgs.jsonConfig.val, err)
			return Config{}
		}

		cfg, err := readFromJson(file)
		if err != nil {
			log.Printf("cfgreader: failed to read file err=%s\n", err)
			return Config{}
		}

		return cfg
	}

	return Config{
		LocalAddr: cmdArgs.localAddr.val,
		ProxyAddr: cmdArgs.proxyAddr.val,
		DohURL:    cmdArgs.dohURL.val,
	}
}

func (cla *commandLineAgrs) parse(args []string) {
	flags := flag.NewFlagSet("", flag.ContinueOnError)

	flags.StringVar(&cla.proxyAddr.val, cla.proxyAddr.flag, "", cla.proxyAddr.descr)
	flags.StringVar(&cla.localAddr.val, cla.localAddr.flag, "", cla.localAddr.descr)
	flags.StringVar(&cla.dohURL.val, cla.dohURL.flag, "", cla.dohURL.descr)
	flags.StringVar(&cla.jsonConfig.val, cla.jsonConfig.flag, "", cla.jsonConfig.descr)

	if err := flags.Parse(args); err != nil {
		log.Printf("cfgreader: failed to parse parameters err=%s\n", err)
	}
}

func readFromJson(r io.Reader) (Config, error) {
	cfg := Config{}

	b, err := io.ReadAll(r)
	if err != nil {
		log.Printf("cfgreader: failed reads json err=%s\n", err)
		return cfg, ErrConfigFileParse
	}

	if err := json.Unmarshal(b, &cfg); err != nil {
		return cfg, ErrConfigFileParse
	}

	return cfg, nil
}
