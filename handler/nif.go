package handler

import (
	"fmt"
	"net"
)

const (
	readSize int = 8192
)

// Nif ...
type Nif struct {
	remoteAddress string
	remotrPort    uint16

	inboundDataPort  net.Conn
	outboundDataPort net.Conn
}

// MakeNif ...
func MakeNif(conn net.Conn, ra string, rp uint16) *Nif {
	nif := new(Nif)
	nif.inboundDataPort = conn

	nif.remoteAddress = ra
	nif.remotrPort = rp

	return nif
}

// Prepare ...
func (n *Nif) Prepare() (string, uint16, error) {

	var rfa string = fmt.Sprintf("%s:%d", n.remoteAddress, n.remotrPort)
	var err error = nil

	fmt.Printf("nif: open remote data port to adds=%s\n", rfa)

	n.outboundDataPort, err = net.Dial("tcp", rfa)
	if err != nil {
		n.inboundDataPort.Close()
		fmt.Printf("nif: failed to connect to remote for adds=%s error=%s\n ", rfa, err.Error())
		return "", 0, err
	}

	ip := n.outboundDataPort.LocalAddr().(*net.TCPAddr)

	return ip.IP.String(), uint16(ip.Port), nil
}

// Run ...
func (n *Nif) Run() {

	fmt.Println("nif: Start relaying")

	inboundRelay := makeRelay(n.inboundDataPort, n.outboundDataPort)
	outboundRelay := makeRelay(n.outboundDataPort, n.inboundDataPort)

	sch := make(chan bool)

	// wait for one of relay part is done. This means one part of realy is disconnected
	// and the other one could be closed
	go inboundRelay.run(sch)
	go outboundRelay.run(sch)

	// wait for someone done their task
	<-sch

	// this force to close channel for read and stop oth coroutines
	n.inboundDataPort.Close()
	n.outboundDataPort.Close()

	n.inboundDataPort = nil
	n.outboundDataPort = nil
}
