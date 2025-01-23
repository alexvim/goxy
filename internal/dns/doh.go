package dns

import (
	"context"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"math"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"time"

	dnsmessage "golang.org/x/net/dns/dnsmessage"
)

var (
	ErrDohHostNotFound = errors.New("host not found")
	ErrDohURL          = errors.New("invalid DoH URL")
)

type doh struct {
	dnsURL     string
	resoveType ResolveType
	client     *http.Client
	cache      recordCache
}

func (d doh) Resolve(domain string) (string, error) {
	if ip, err := d.cache.get(domain); err == nil {
		log.Printf("doh: dns found in cache %s -> %s\n", domain, ip)
		return ip, nil
	}

	// skip IP address
	if net.ParseIP(domain) != nil {
		return domain, nil
	}

	dnsWireQuery, err := makeDnsWireQuery(domain, d.resoveType)
	if err != nil {
		log.Printf("doh: failed to create dns wire request err=%s\n", err)
		return "", err
	}

	req, err := makeDoHGet(d.dnsURL, dnsWireQuery)
	if err != nil {
		log.Printf("doh: failed to create DoH request err=%s\n", err)
		return "", err
	}

	resp, err := d.client.Do(req)
	if err != nil {
		log.Printf("doh: failed to send DoH request err=%s\n", err)
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("doh: http request finifhed with error code=%d\n", resp.StatusCode)
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("doh: failed to read DoH response body err=%s\n", err)
		return "", err
	}

	ipvAddrs, ttl, err := parseDnsWireQuery(body)
	if err != nil {
		log.Printf("doh: failed to parse dns wire response err=%s\n", err)
		return "", err
	}

	return d.cache.put(domain, ipvAddrs, time.Duration(ttl)*time.Second), err
}

func makeDoHGet(server string, dnsWireQuery []byte) (*http.Request, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server, nil)
	if err != nil {
		log.Printf("doh: failed to create request err=%s\n", err)
		return nil, err
	}

	req.Header.Add("accept", "application/dns-message")
	req.Header.Add("user-agent", "goxy")

	hquery := req.URL.Query()

	hquery.Add("dns", base64.RawURLEncoding.EncodeToString(dnsWireQuery))

	req.URL.RawQuery = hquery.Encode()

	return req, nil
}

func makeDnsWireQuery(domain string, qtype ResolveType) ([]byte, error) {
	builder := dnsmessage.NewBuilder(nil, dnsmessage.Header{
		RecursionDesired: true,
	})

	dnsq := dnsmessage.Question{
		// make name canonical
		Name: dnsmessage.MustNewName(domain + "."),
		// Only IPv4/IPv6
		Type: dnsmessage.Type(qtype),
		// Only inet
		Class: dnsmessage.ClassINET,
	}

	if err := builder.StartQuestions(); err != nil {
		log.Printf("doh: failed to start dns wire questions section err=%s\n", err)
		return nil, err
	}

	if err := builder.Question(dnsq); err != nil {
		log.Printf("doh: failed to build dns wire question err=%s\n", err)
		return nil, err
	}

	bin, err := builder.Finish()
	if err != nil {
		log.Printf("doh: failed to build dns request err=%s\n", err)
		return nil, err
	}

	return bin, nil
}

// Parse dns message and return only INET and set of ip4/ip6 addresses depands on request
func parseDnsWireQuery(msg []byte) ([]string, uint32, error) {
	parser := dnsmessage.Parser{}

	header, err := parser.Start(msg)
	if err != nil {
		log.Printf("doh: failed to parse dns wire err=%s\n", err)
		return nil, 0, err
	}

	if !header.Response || header.RCode != dnsmessage.RCodeSuccess {
		log.Printf("doh: faild to complete dns query resp=%t code=%d\n", header.Response, header.RCode)
		return nil, 0, err
	}

	if err := parser.SkipAllQuestions(); err != nil {
		log.Printf("doh: failed to skip questoins in dns query err=%s\n", err)
		return nil, 0, err
	}

	answers, err := parser.AllAnswers()
	if err != nil {
		log.Printf("doh: failed to get dns answers err=%s\n", err)
		return nil, 0, err
	}

	ttl := uint32(math.MaxUint32)
	ipAddrs := make([]string, 0)
	for _, answer := range answers {
		// INET only
		if answer.Header.Class != dnsmessage.ClassINET {
			continue
		}

		// Just A or AAAA records are acceptable
		if answer.Header.Type != dnsmessage.TypeA && answer.Header.Type != dnsmessage.TypeAAAA {
			continue
		}

		switch record := answer.Body.(type) {
		case *dnsmessage.AResource:
			ipAddrs = append(ipAddrs, netip.AddrFrom4(record.A).String())
		case *dnsmessage.AAAAResource:
			ipAddrs = append(ipAddrs, netip.AddrFrom16(record.AAAA).String())
		default:
			continue
		}

		ttl = min(ttl, answer.Header.TTL)
	}

	return ipAddrs, ttl, nil
}

func newDoHResolver(dns string, rt ResolveType) (doh, error) {
	log.Printf("doh: create new DoH resolver via %s for %s addresses", dns, rt)

	dnsURL, err := url.JoinPath("https://", dns)
	if err != nil {
		log.Printf("doh: invalid DoH URL provided")
		return doh{}, ErrDohURL
	}

	return doh{
		dnsURL:     dnsURL,
		resoveType: rt,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		cache: newDnsCache(),
	}, nil
}
