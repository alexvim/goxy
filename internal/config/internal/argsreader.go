package configreader

import (
	"flag"
	"log"
	"os"
)

var (
	cmdArgsDefaults = CommandLineAgrs{
		proxyAddr: CmdArg{
			flag:  "p",
			descr: "socks5 proxy and port in format ip:port to bind",
		},
		localAddr: CmdArg{
			flag:  "l",
			descr: "target network ip address",
		},
		dohURL: CmdArg{
			flag:  "doh",
			descr: "DoH URL",
		},
		jsonConfig: CmdArg{
			flag:  "c",
			descr: "full path to json config location",
		},
	}
)

type CmdArg struct {
	val   string
	flag  string
	descr string
}

type CommandLineAgrs struct {
	proxyAddr  CmdArg
	localAddr  CmdArg
	dohURL     CmdArg
	jsonConfig CmdArg
}

func ParseArgs(args []string) *CommandLineAgrs {
	cargs := cmdArgsDefaults

	cmdArgs := &cargs

	cmdArgs.parse(args)

	return cmdArgs
}

func (cla *CommandLineAgrs) ProxyAddr() string {
	return cla.proxyAddr.val
}

func (cla *CommandLineAgrs) LocalAddr() string {
	return cla.localAddr.val
}

func (cla *CommandLineAgrs) DohURL() string {
	return cla.dohURL.val
}

func (cla *CommandLineAgrs) ConfigFile() string {
	return cla.jsonConfig.val
}

func (cla *CommandLineAgrs) HasConfigFile() bool {
	if len(cla.jsonConfig.val) == 0 {
		return false
	}

	if _, err := os.Stat(cla.jsonConfig.val); err == nil {
		return true
	}

	return false
}

func (cla *CommandLineAgrs) parse(args []string) {
	flags := flag.NewFlagSet("", flag.ContinueOnError)

	flags.StringVar(&cla.proxyAddr.val, cla.proxyAddr.flag, "", cla.proxyAddr.descr)
	flags.StringVar(&cla.localAddr.val, cla.localAddr.flag, "", cla.localAddr.descr)
	flags.StringVar(&cla.dohURL.val, cla.dohURL.flag, "", cla.dohURL.descr)
	flags.StringVar(&cla.jsonConfig.val, cla.jsonConfig.flag, "", cla.jsonConfig.descr)

	if err := flags.Parse(args); err != nil {
		log.Printf("cfgreader: failed to parse parameters err=%s\n", err)
	}
}
