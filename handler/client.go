package handler

import (
	"errors"
	"fmt"
	"goxy/msg"
	"net"

	"github.com/google/uuid"
)

type Client struct {
	id         uuid.UUID
	connection net.Conn
}

func MakeClient(conn net.Conn) *Client {

	var err error

	c := new(Client)

	c.id, err = uuid.NewUUID()
	if err != nil {
		fmt.Println("client: failed to generate uuid")
	}

	c.connection = conn
	return c
}

func (c *Client) SendMessage(s msg.Serializeable) bool {

	_, err := c.connection.Write(s.Serialize())
	if err != nil {
		fmt.Printf("client: failed to send message %s\n", err)
		return false
	}

	return true
}

func (c *Client) ReadMessage() ([]byte, error) {

	var buffer []byte = make([]byte, 50)
	n, err := c.connection.Read(buffer)
	if err != nil || n <= 0 {
		return nil, errors.New("failed to read buffer err=" + err.Error())
	}

	return buffer[0:n], nil
}

func (c *Client) Disconnect() {
	if c.connection != nil {
		c.connection.Close()
	}
}
