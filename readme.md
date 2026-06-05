# Goxy

SOCKS5 proxy with addtional features.
Usually, proxy accept incoming traffic from specific network and forward it to another network. To help in this matter extra parameters are exposed by Goxy.

Addtional featutes list:

- Support host name resolution via DoH

## SOCKS5 standard support

The proxy is trying to comply with the RFCs:

- [RFC1928](https://www.rfc-editor.org/rfc/rfc1928)
- [RFC1929](https://www.rfc-editor.org/rfc/rfc1929)

### Commands

- CONNECT

### Authentication methods

- No authentication

## Build

Two build methods are available, classical `go build` to build binary and building docker image.

- Build binary

    ```shell
    go build -o goxy cmd/main.go
    ```

- Build Docker image

    ```shell
    docker build -t goxy .
    ```

## Configure Goxy

TBD

## Examples

### Run as binary

- `goxy -p 192.168.0.1 -l 172.0.0.1`
- `goxy -p 192.168.0.1 -l 172.0.0.1 -doh doh.opendns.com/dns-query`
- `goxy -c config.json`

### Run in `Docker` container

Buld image as mentioned above and run container as deamon

with parameters

```shell
docker run -d --rm --network host --name goxy-proxy goxy:latest -p 192.168.0.1 -l 172.0.0.1 -doh doh.opendns.com/dns-query
```

with custom config located at `./config/config.json`

```shell
docker run --rm --network host --name goxy-proxy -v ./config/config.json:/etc/goxy/config.json:ro goxy:latest -c /etc/goxy/config.json
```

or start container with restart on system/docker startup (no parameters example)

```shell
docker run -d --restart unless-stopped --rm --network host --name goxy-proxy goxy:latest
```

To stop container execute:

```shell
docker stop goxy-proxy
```

### Run as `systemd` service

#### Arch Linux

Depends on what you what you should place systemd unit file in correct place. Read about `systemd` before do your experiments [ArchLinux Docs](https://wiki.archlinux.org/title/Systemd).

This example put service file with `docker container` launch to be stared on user _alex_ logon

1. Create user `systmed` folder

    ```shell
    mkdir -p ~/.config/systemd/alex
    ```

2. Copy `service` file

    ```shell
    cp script/goxy-proxy-container.service .config/systemd/alex
    ```

3. Reload `systemd` and start `goxy service`

    ```shell
    systemctl --user daemon-reload
    systemctl --user enable goxy-proxy-container.service
    systemctl --user start goxy-proxy-container.service
    ```

#### Ubuntu Linux 22+

TDB
