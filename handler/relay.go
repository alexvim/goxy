package handler

import (
	"fmt"
	"net"
)

const queueLength int = 100

type relay struct {
	src       net.Conn
	dst       net.Conn
	relayChan chan *[]byte
	sch       chan<- bool
}

func makeRelay(src net.Conn, dst net.Conn) *relay {
	r := new(relay)
	r.src = src
	r.dst = dst
	r.relayChan = make(chan *[]byte, queueLength)
	return r
}

func (r *relay) run(sch chan<- bool) {

	r.sch = sch
	go r.read(r.src)
	go r.write(r.dst)
}

func (r *relay) read(conn net.Conn) {

	var ch chan<- *[]byte = r.relayChan

	remoteAddr := conn.LocalAddr().String()

	fmt.Printf("nif: start reading from: %s\n", remoteAddr)

	bufferLen := readSize * (cap(ch) + 1)
	buf := make([]byte, bufferLen)
	rindex := 0
	for {
		byteRead, err := conn.Read(buf[rindex : rindex+readSize])
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

		rindex += byteRead
		if rindex+readSize > bufferLen {
			rindex = 0
		}
	}

	fmt.Printf("nif: stop reading from: %s\n", remoteAddr)

	r.sch <- true
}

func (r *relay) write(conn net.Conn) {

	var ch <-chan *[]byte = r.relayChan

	remoteAddr := conn.RemoteAddr().String()

	fmt.Printf("nif: start writing to %s\n", remoteAddr)

	for buf, ok := <-ch; ok == true; buf, ok = <-ch {
		if n, err := conn.Write(*buf); err != nil {
			fmt.Printf("nif: %s on writing to %s\n", err.Error(), remoteAddr)
			break
		} else if n > 0 {
			fmt.Printf("nif: %d bytes was written to %s\n", n, remoteAddr)
		}
	}

	fmt.Printf("nif: close write stream to %s\n", remoteAddr)

	r.sch <- true
}
