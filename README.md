# Fail2ban-exporter

[![Build Status](https://travis-ci.com/syberalexis/fail2ban-exporter.svg?branch=master)](https://travis-ci.com/syberalexis/fail2ban-exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/syberalexis/fail2ban-exporter)](https://goreportcard.com/report/github.com/syberalexis/fail2ban-exporter)

This exporter get informations from Fail2ban and can geoloc IPs if feature is enabled.

## Summary

- [Install](#install)
  - [From binary](#from-binary)
  - [From docker](#from-docker)
  - [From sources](#from-sources)
- [Install as a service](#install-as-a-service)
  - [Systemd](#systemd)
- [Dashboards](#dashboards)
- [Help](#help)
- [Metrics example](#metrics-example)

## Install

### From binary

Download binary from [releases page](https://github.com/syberalexis/fail2ban-exporter/releases)

Example :
```bash
curl -L https://github.com/syberalexis/fail2ban-exporter/releases/download/v3.0.0/fail2ban-exporter-3.0.0-linux-amd64 -o /usr/local/bin/fail2ban-exporter
chmod +x /usr/local/bin/fail2ban-exporter
/usr/local/bin/fail2ban-exporter
```

### From docker

```bash
docker pull syberalexis/fail2ban-exporter
docker run -d -p 9901:9901 -v /dev/serial0:/dev/serial0 syberalexis/fail2ban-exporter:1.0.0
```

### From sources

```bash
git clone git@github.com:syberalexis/fail2ban-exporter.git
cd fail2ban-exporter
go build cmd/fail2ban-exporter/main.go -o fail2ban-exporter
./fail2ban-exporter
```

or

```bash
git clone git@github.com:syberalexis/fail2ban-exporter.git
cd fail2ban-exporter
GOOS=linux GOARCH=amd64 VERSION=3.0.0 make clean build
./dist/fail2ban-exporter-3.0.0-linux-amd64
```

## Install as a service

In file `/lib/systemd/system/fail2ban_exporter.service` :
### Systemd
```
[Unit]
Description=Fail2ban Exporter service
Wants=network-online.target
After=network-online.target

[Service]
Type=simple
ExecStart=/usr/local/bin/fail2ban-exporter
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

```bash
systemctl enable fail2ban-exporter
systemctl start fail2ban-exporter
```

## Help

```
usage: fail2ban-exporter --device=DEVICE [<flags>]

Flags:
  -h, --help               Show context-sensitive help (also try --help-long and --help-man).
      --version            Show application version.
      --debug              Enable debug mode.
      --localisation       Enable localisation mode.
```

## Metrics example

```
```