# Puprose
goxy is a sock5 proxy with extended features

## Supported SOCKS5 features

### Authentication methods
- No authentication

### Commands
- CONNECT

# Features
## Support host name resolution via DoH

# Build
go build -o goxy cmd/main.go

# Run examples
- goxy -p 192.168.0.1 -l 172.0.0.1
- goxy -p 192.168.0.1 -l 172.0.0.1 -doh doh.opendns.com/dns-query
- goxy -c config.json
