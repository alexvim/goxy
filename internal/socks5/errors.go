package socks5

import "errors"

var (
	ErrUnsupportedProtocolVersion = errors.New("unsupported protocol version")
	ErrMalformedPacket            = errors.New("malformed packet")
	ErrSendMessage                = errors.New("error on send socks5 message")
	ErrReadMessage                = errors.New("error on read socks5 message")
	ErrParseMessage               = errors.New("error on parse socks5 message")
	ErrAuthenticate               = errors.New("failed to authenticate")
	ErrSocksConnect               = errors.New("failed to connect")
)
