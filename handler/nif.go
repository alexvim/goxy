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

	dataChannelInbound := make(chan *[]byte, 100)
	// read from client and push to remote
	go readFrom(n.inboundDataPort, dataChannelInbound)
	go writeTo(n.outboundDataPort, dataChannelInbound)

	dataChannelOutound := make(chan *[]byte, 100)
	// read from remote and push to client
	go readFrom(n.outboundDataPort, dataChannelOutound)
	go writeTo(n.inboundDataPort, dataChannelOutound)
}

// Destroy ...
func (n *Nif) Destroy() {
	fmt.Println("nif: destroying...")
	if n.inboundDataPort != nil {
		n.inboundDataPort.Close()
	}
	if n.outboundDataPort != nil {
		n.outboundDataPort.Close()
	}
}

func readFrom(conn net.Conn, ch chan<- *[]byte) {

	remoteAddr := conn.LocalAddr().String()

	fmt.Printf("nif: start reading from: %s\n", remoteAddr)

	var byteRead int = 0
	var err error = nil

	bufferLen := readSize * (cap(ch) + 1)
	buf := make([]byte, bufferLen)
	rindex := 0
	for {
		byteRead, err = conn.Read(buf[rindex : rindex+readSize])
		if err != nil {
			fmt.Printf("nif: error {%s} on reading from %s\n", err.Error(), remoteAddr)
			close(ch)
			break
		}

		if byteRead <= 0 {
			fmt.Printf("nif: read again on %s\n", remoteAddr)
			continue
		}

		fmt.Printf("nif: read %d bytes from %s\n", byteRead, remoteAddr)

		b := buf[rindex : rindex+byteRead]
		ch <- &b

		rindex = rindex + byteRead
		if rindex+readSize > bufferLen {
			rindex = 0
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
