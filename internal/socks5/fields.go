package socks5

import "fmt"

const (
	// X'00'
	ReservedDefault Reserved = 0x00
)

// RESERVED
type Reserved byte

const (
	// protocol version: X'05'
	ProtoclVersion5 ProtoclVersion = 0x05
)

// Protocol version
type ProtoclVersion byte

const (
	Undefined Command = iota
	// CONNECT X'01'
	Connect
	// BIND X'02'
	Bind
	// UDP ASSOCIATE X'03'
	UdpAssosiate
	// Last
	Unknown
)

// CMD
type Command byte

const (
	// IP V4 address: X'01'
	AddressTypeIP4 AddressType = 0b001
	// DOMAINNAME: X'03'
	AddressTypeDomainName AddressType = 0b011
	// IP V6 address: X'04'
	AddressTypeIP6 AddressType = 0b100
)

// ATYP address type of following address
type AddressType byte

const (
	// X'00' succeeded
	Succeeded Rep = iota
	// X'01' general SOCKS server failure
	GeneralSocksServerFailure
	// X'02' connection not allowed by ruleset
	ConnectionNotAllowedByRuleset
	// X'03' Network unreachable
	NetworkUnreachable
	// X'04' Host unreachable
	HostUnreachable
	//  X'05' Connection refused
	ConnectionRefused
	// X'06' TTL expired
	TtlExpired
	// X'07' Command not supported
	CommandNotSupported
	// X'08' Address type not supported
	AddressTypeNotSupport
	// X'09' to X'FF' unassigned
	Unassigned
)

// Rep field
type Rep byte

const (
	// X'01'
	AuthVer0x01 AuthVer = 0x01
)

// The VER field contains the current version of the subnegotiation, which is X'01'
type AuthVer byte

// Auth methods
const (
	AuthMethodNoAuthRequired      AuthHandshakeMethod = 0x00
	AuthMethodGSSAPI              AuthHandshakeMethod = 0x01
	AuthMethodUserAndPassword     AuthHandshakeMethod = 0x02
	AuthMethodNoAcceptableMethods AuthHandshakeMethod = 0xFF
)

type AuthHandshakeMethod byte

const (
	// A STATUS field of X'00' indicates success.
	AuthStatusSuccess AuthStatus = iota
	// `failure' (STATUS value other than X'00') status
	AuthStatusFailure
)

// A STATUS field of X'00' indicates success.
// If the server returns a`failure' (STATUS value other than X'00') status, it MUST close the connection.
type AuthStatus byte

func (rsv Reserved) String() string {
	return fmt.Sprintf("reserved(%d)", rsv)
}

func (pver ProtoclVersion) String() string {
	return fmt.Sprintf("0x%x", int(pver))
}

func (cmd Command) String() string {
	switch cmd {
	case Connect:
		return "connect"
	case Bind:
		return "bind"
	case UdpAssosiate:
		return "udp_associate"
	default:
		return fmt.Sprintf("undefined(0x%x)", int(cmd))
	}
}

func (addrType AddressType) String() string {
	switch addrType {
	case AddressTypeIP4:
		return "ipv4"
	case AddressTypeDomainName:
		return "domain"
	case AddressTypeIP6:
		return "ipv6"
	default:
		return fmt.Sprintf("undefined(0x%x)", int(addrType))
	}
}

func (rep Rep) String() string {
	switch rep {
	case Succeeded:
		return "success"
	case GeneralSocksServerFailure:
		return "general_socks_server_failure"
	case ConnectionNotAllowedByRuleset:
		return "connection_not_allowed_by_ruleset"
	case NetworkUnreachable:
		return "network_unreachable"
	case HostUnreachable:
		return "host_unreachable"
	case ConnectionRefused:
		return "connection_refused"
	case TtlExpired:
		return "ttl_expired"
	case CommandNotSupported:
		return "command_not_supported"
	case AddressTypeNotSupport:
		return "address_type_not_supported"
	default:
		return fmt.Sprintf("undefined(0x%x)", int(rep))
	}
}

func (authVer AuthVer) String() string {
	return fmt.Sprintf("0x%x", int(authVer))
}

func (authMethod AuthHandshakeMethod) String() string {
	switch authMethod {
	case AuthMethodNoAuthRequired:
		return "no_auth_required"
	case AuthMethodGSSAPI:
		return "gssapi"
	case AuthMethodUserAndPassword:
		return "username/password"
	case AuthMethodNoAcceptableMethods:
		return "no_acceptable_methods"
	default:
		return fmt.Sprintf("undefined(%b)", authMethod)
	}
}

func (status AuthStatus) String() string {
	if status == AuthStatusSuccess {
		return "success"
	}

	return fmt.Sprintf("failure(0x%x)", int(status))
}
