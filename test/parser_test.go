package test

import (
	"goxy/msg"
	"testing"
)

func TestAuthCommandPasrer(t *testing.T) {

	var buffer []byte = []byte{0x05, 0x01, 0x00}

	message, err := msg.ParseAuth(buffer)

	if err != nil {
		t.Error(err.Error())
	}

	if message == nil {
		t.Error("message is null")
	}
}

func TestCmdCommandPasrerIpv4(t *testing.T) {

	var buffer []byte = []byte{0x05, 0x01, 0x00, 0x01, 0x23, 0x32, 0x43, 0x10, 0x1F, 0x90}

	message, err := msg.ParseCommand(buffer)

	if err != nil {
		t.Error(err.Error())
	}

	if message == nil {
		t.Error("message is null")
	}

	if message.DstPort != 8080 {
		t.Error(message.DstPort)
	}

	if message.DstAddr != "35.50.35.35" {
		t.Error(message.DstAddr)
	}
}

func TestCmdCommandPasrerIpv6(t *testing.T) {

	var buffer []byte = []byte{0x05, 0x01, 0x00, 0x04, 0x23, 0x32, 0x43, 0x23, 0x23, 0x32, 0x43, 0x23, 0x23, 0x32, 0x43, 0x23, 0x23, 0x32, 0x43, 0x23, 0x1F, 0x90}

	message, err := msg.ParseCommand(buffer)

	if err != nil {
		t.Error(err.Error())
	}

	if message == nil {
		t.Error("message is null")
	}

	if message.DstPort != 8080 {
		t.Error(message.DstPort)
	}

	if message.DstAddr != "2332:4323:2332:4323:2332:4323:2332:4323" {
		t.Error(message.DstAddr)
	}
}

func TestCmdCommandPasrerDomain(t *testing.T) {

	var buffer []byte = []byte{0x05, 0x01, 0x00, 0x03, 0x5, 0x79, 0x61, 0x2e, 0x72, 0x75, 0x1F, 0x90}

	message, err := msg.ParseCommand(buffer)

	if err != nil {
		t.Error(err.Error())
	}

	if message == nil {
		t.Error("message is null")
	}

	if message.DstPort != 8080 {
		t.Error(message.DstPort)
	}
}
