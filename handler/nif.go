package handler

import (
	"fmt"
	"goxy/network"
	"net"
)

var port uint16 = 4040

type Nif struct {
	bindAddress string
	bindPort    uint16

	remoteAddress string
	remotrPort    uint16

	inboundDataPort  net.Conn
	outboundDataPort net.Conn
}

func MakeNif(conn net.Conn, ra string, rp uint16) *Nif {
	nif := new(Nif)

	nif.bindAddress, _ = network.GetLocalInetAddress(network.AddrTypeIpv4)
	nif.bindPort = port
	port += 1

	nif.inboundDataPort = conn

	nif.remoteAddress = ra
	nif.remotrPort = rp

	return nif
}

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

func (n *Nif) Run() {

	fmt.Println("nif: Start relaying")

	dataChannelInbound := make(chan *[]byte, 1000)
	// read from client and push to remote
	go readFrom(n.inboundDataPort, dataChannelInbound)
	go writeTo(n.outboundDataPort, dataChannelInbound)

	dataChannelOutound := make(chan *[]byte, 1000)
	// read from remote and push to client
	go readFrom(n.outboundDataPort, dataChannelOutound)
	go writeTo(n.inboundDataPort, dataChannelOutound)
}

func readFrom(conn net.Conn, ch chan<- *[]byte) {

	remoteAddr := conn.LocalAddr().String()

	fmt.Printf("nif: start reading from: %s\n", remoteAddr)

	for {
		buf := make([]byte, 1500)
		if n, err := conn.Read(buf[0:1500]); err == nil && n > 0 {
			fmt.Printf("nif: read %d bytes from %s\n", n, remoteAddr)
			b := buf[0:n]
			ch <- &b
		} else if err != nil {
			fmt.Printf("nif: error {%s} on reading from %s\n", err.Error(), remoteAddr)
			close(ch)
			break
		} else {
			fmt.Printf("nif: read again on %s\n", remoteAddr)
		}
	}
	fmt.Printf("nif: stop reading from: %s\n", remoteAddr)
}

func writeTo(conn net.Conn, ch <-chan *[]byte) {

	remoteAddr := conn.RemoteAddr().String()

	fmt.Printf("nif: start writing to %s\n", remoteAddr)

	for buf, ok := <-ch; ok == true; buf, ok = <-ch {
		n, err := conn.Write(*buf)
		if err != nil {
			fmt.Printf("nif: %s on writing to %s\n", err.Error(), remoteAddr)
			break
		}
		if n > 0 {
			fmt.Printf("nif: %d bytes was written to %s\n", n, remoteAddr)
		}
	}

	fmt.Printf("nif: close write stream to %s\n", remoteAddr)
}

func (n *Nif) Destroy() {
	fmt.Println("nif: destroying...")
	if n.inboundDataPort != nil {
		n.inboundDataPort.Close()
	}
	if n.outboundDataPort != nil {
		n.outboundDataPort.Close()
	}
}
