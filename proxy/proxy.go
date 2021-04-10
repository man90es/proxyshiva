package proxy

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"inet.af/netaddr"
)

const defaultURL = "https://ipecho.net/plain"

var reservedSubnets = [...]netaddr.IPPrefix{
	netaddr.MustParseIPPrefix("0.0.0.0/8"),
	netaddr.MustParseIPPrefix("10.0.0.0/8"),
	netaddr.MustParseIPPrefix("100.64.0.0/10"),
	netaddr.MustParseIPPrefix("127.0.0.0/8"),
	netaddr.MustParseIPPrefix("169.254.0.0/16"),
	netaddr.MustParseIPPrefix("172.16.0.0/12"),
	netaddr.MustParseIPPrefix("192.0.0.0/24"),
	netaddr.MustParseIPPrefix("192.0.2.0/24"),
	netaddr.MustParseIPPrefix("192.88.99.0/24"),
	netaddr.MustParseIPPrefix("192.168.0.0/16"),
	netaddr.MustParseIPPrefix("198.18.0.0/15"),
	netaddr.MustParseIPPrefix("198.51.100.0/24"),
	netaddr.MustParseIPPrefix("203.0.113.0/24"),
	netaddr.MustParseIPPrefix("224.0.0.0/4"),
	netaddr.MustParseIPPrefix("240.0.0.0/4"),
	netaddr.MustParseIPPrefix("255.255.255.255/32"),
}

type Proxy struct {
	Scheme      string      `json:"scheme"`
	Address     netaddr.IP  `json:"address"`
	Port        uint16      `json:"port"`
	Good        bool        `json:"good"`
	ExitAddress *netaddr.IP `json:"exitAddress,omitempty"`
	Error       string      `json:"error,omitempty"`
	Speed       float64     `json:"speed,omitempty"`
}

func (p Proxy) IsReserved() bool {
	for _, subnet := range reservedSubnets {
		if subnet.Contains(p.Address) {
			return true
		}
	}

	return false
}

func (p Proxy) Check(resultQueue chan<- *Proxy, timeout *int, skipCert *bool) {
	addressString := p.Address.String()
	startAt := time.Now()
	proxyUrl := &url.URL{
		Scheme: p.Scheme,
		Host:   fmt.Sprintf("%v:%v", addressString, p.Port),
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy:           http.ProxyURL(proxyUrl),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: *skipCert},
		},

		Timeout: time.Duration(time.Duration(*timeout) * time.Second),
	}

	if response, err := client.Get(defaultURL); err == nil {
		p.Speed = float64(time.Now().UnixNano()-startAt.UnixNano()) / 1e9
		p.Good = true

		body, _ := ioutil.ReadAll(response.Body)
		defer response.Body.Close()

		if strings.Contains(string(body), addressString) {
			p.ExitAddress = &p.Address
		} else {
			exitAddress, _ := netaddr.ParseIP(string(body))
			p.ExitAddress = &exitAddress
		}
	} else {
		p.Good = false
		p.Error = err.Error()
	}

	resultQueue <- &p
}
