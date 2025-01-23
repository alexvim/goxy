package dns

import (
	dnsmessage "golang.org/x/net/dns/dnsmessage"
)

const (
	ResolveTypeIPv4 ResolveType = ResolveType(dnsmessage.TypeA)
	ResolveTypeIPv6 ResolveType = ResolveType(dnsmessage.TypeAAAA)
)

type ResolveType uint8

type Resolver interface {
	Resolve(domain string) (string, error)
}

type stubResolver struct {
}

func NewDNSResolver(doh string, resType ResolveType) Resolver {
	resolver, err := newDoHResolver(doh, resType)
	if err != nil {
		return stubResolver{}
	}

	return resolver
}

func (stubResolver) Resolve(_ string) (string, error) {
	return "", nil
}

func (rt ResolveType) String() string {
	switch rt {
	case ResolveTypeIPv4:
		return "IPv4"
	case ResolveTypeIPv6:
		return "IPv6"
	default:
		return "undefined"
	}
}
