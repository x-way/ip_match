# ip_match
[![CircleCI](https://circleci.com/gh/x-way/ip_match/tree/master.svg?style=svg)](https://circleci.com/gh/x-way/ip_match/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/x-way/ip_match)](https://goreportcard.com/report/github.com/x-way/ip_match)

Filter IPs/networks by matching against a list of prefixes.

## Installation

```
# go get github.com/x-way/ip_match
```

## Usage

```
# cat filter_prefixes.txt
192.168.128.0/17
2001:db8:1234::/48

# cat iplist.txt
192.168.1.2
192.168.2.3
192.168.128.0/24
192.168.128.123
192.168.160.0/20
3.4.5.6
192.168.134.20
10.10.10.1
2001:db8::1
10.10.100.0/24
10.20.20.123
2001:db8:1234:1234:1234:1234:1234:1234
2001:db8:1234::1234:1234
2001:db8:1233::1234:1234
2001:db8:1234:0:1234:1234::/64
10.20.20.127
10.20.21.0/24

# ip_match -F filter_prefixes.txt iplist.txt
2001:db8:1234:1234:1234:1234:1234:1234/128
2001:db8:1234::1234:1234/128
2001:db8:1234::/64
192.168.128.0/24
192.168.128.123/32
192.168.160.0/20
192.168.134.20/32
```
