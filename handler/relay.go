package handler

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var localReadCounter uint64 = 0
var localWriteCounter uint64 = 0

func printRwCounter() {
	fmt.Printf("relay: r/w counters: r=%d w=%d\n", localReadCounter, localWriteCounter)
}

const (
	queueLength int = 1
	readSize    int = 256
)

var relayGlobalCounter uint64 = 0

type relay struct {
	id        uint64
	src       net.Conn
	dst       net.Conn
	readEof   bool
	writeEof  bool
	relayChan chan *[]byte
	sch       chan<- bool
}

func makeRelay(src net.Conn, dst net.Conn) *relay {
	relayGlobalCounter++
	r := new(relay)
	r.id = relayGlobalCounter
	r.src = src
	r.dst = dst
	r.readEof = false
	r.writeEof = false
	r.relayChan = make(chan *[]byte, queueLength)
	return r
}

func (r *relay) close() {
}

func (r *relay) run(sch chan<- bool) {

	var wg sync.WaitGroup

	wg.Add(2)

	go r.read(r.src, &wg)
	go r.write(r.dst, &wg)

	printRwCounter()

	wg.Wait()

	sch <- true

	printRwCounter()
}

func (r *relay) read(conn net.Conn, wg *sync.WaitGroup) {

	atomic.AddUint64(&localReadCounter, 1)

	defer wg.Done()

	var ch chan<- *[]byte = r.relayChan

	remoteAddr := conn.LocalAddr().String()

	fmt.Printf("relay{%d}: start read stream from: %s\n", r.id, remoteAddr)

	bufferLen := readSize * (cap(ch) + 1)
	buf := make([]byte, bufferLen)
	rindex := 0
	for {

		if r.writeEof {
			fmt.Printf("relay{%d}: error {write stream closed} on reading from %s\n", r.id, remoteAddr)
			break
		}

		conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		byteRead, err := conn.Read(buf[rindex : rindex+readSize])
		if e, ok := err.(net.Error); ok && e.Timeout() {
			continue
		} else if err != nil {
			fmt.Printf("relay{%d}: error {%s} on reading from %s\n", r.id, err.Error(), remoteAddr)
			// send EOF to write
			ch <- nil
			break
		}

		if byteRead <= 0 {
			fmt.Printf("relay{%d}: read again on %s\n", r.id, remoteAddr)
			continue
		}

		b := buf[rindex : rindex+byteRead]
		ch <- &b

		rindex += byteRead
		if rindex+readSize >= bufferLen {
			rindex = 0
		}
	}

	r.readEof = true
	atomic.AddUint64(&localReadCounter, ^uint64(0))
	fmt.Printf("relay{%d}: close read stream from: %s\n", r.id, remoteAddr)
}

func (r *relay) write(conn net.Conn, wg *sync.WaitGroup) {

	atomic.AddUint64(&localWriteCounter, 1)

	defer wg.Done()

	var ch <-chan *[]byte = r.relayChan

	remoteAddr := conn.RemoteAddr().String()

	fmt.Printf("relay{%d}: start write stream to %s\n", r.id, remoteAddr)

	for buf := <-ch; buf != nil; buf = <-ch {
		if n, err := conn.Write(*buf); err != nil {
			fmt.Printf("relay{%d}: error {%s} on writing to %s\n", r.id, err.Error(), remoteAddr)
			break
		} else if n > 0 {
			//fmt.Printf("relay: %d bytes was written to %s\n", n, remoteAddr)
		}
	}

	r.writeEof = true
	atomic.AddUint64(&localWriteCounter, ^uint64(0))
	// make fake read to unqueue data and unclock write to channel write if it was full
	select {
	case _, ok := <-ch:
		if ok {
			fmt.Printf("relay{%d}: error {queue is not drain} on writing to %s\n", r.id, remoteAddr)
		}
	default:
		//
	}
	fmt.Printf("relay{%d}: close write stream to %s\n", r.id, remoteAddr)
}
