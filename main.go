package main

import (
	"fmt"
	"bufio"
	"os"
	"flag"
	"time"
	"strings"
	"strconv"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"crypto/tls"

	"inet.af/netaddr"
)

const defaultURL = "https://ipecho.net/plain"

type proxy struct {
	Scheme 		string 		`json:"scheme"`
	Address 	netaddr.IP 	`json:"address"`
	Port 		uint16 		`json:"port"`
	Good 		bool 		`json:"good"`
	ExitAddress *netaddr.IP `json:"exitAddress,omitempty"`
	Error 		string 		`json:"error,omitempty"`
	Speed 		float64 	`json:"speed,omitempty"`
}

func (p proxy) check(resultQueue chan<- *proxy, timeout *int, skipCert *bool) {
	addressString := p.Address.String()
	startAt := time.Now()
	proxyUrl := &url.URL{
		Scheme: p.Scheme,
		Host: fmt.Sprintf("%v:%v", addressString, p.Port),
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: *skipCert},
		},

		Timeout: time.Duration(time.Duration(*timeout) * time.Second),
	}

	if response, err := client.Get(defaultURL); err == nil {
		p.Speed = float64(time.Now().UnixNano() - startAt.UnixNano()) / 1e9
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

func requestGenerator(in string) chan proxy {
	out := make(chan proxy)

	schemeEndIndex := strings.Index(in, "://")
	addressStartIndex := schemeEndIndex + 3
	addressEndIndex := strings.LastIndex(in, ":")
	portStartIndex := addressEndIndex + 1

	// Parse schemes
	schemes := strings.Split(in[:schemeEndIndex], ",")

	// Parse the address range
	addressRange, _ := netaddr.ParseIPRange(in[addressStartIndex:addressEndIndex])
	if addressRange.String() == "zero IPRange" { // Handle single input
		addressRange.From, _ = netaddr.ParseIP(in[addressStartIndex:addressEndIndex])
		addressRange.To = addressRange.From
	}

	// Parse the port range
	var portRange [2]uint16
	if portRangeStr := strings.Split(in[portStartIndex:], "-"); len(portRangeStr) == 2 {
		sP, _ := strconv.Atoi(portRangeStr[0])
		eP, _ := strconv.Atoi(portRangeStr[1])

		portRange[0] = uint16(sP)
		portRange[1] = uint16(eP)
	} else {
		p, _ := strconv.Atoi(in[portStartIndex:])

		portRange[0] = uint16(p)
		portRange[1] = portRange[0]
	}

	go func() {
		defer close(out)

		for _, scheme := range schemes { // Rotate schemes
			for port := portRange[0]; port <= portRange[1]; port++ { // Rotate ports
				for address := addressRange.From; address.Less(addressRange.To) || address == addressRange.To; address = address.Next() { // Rotate IPs
					out <- proxy{
						Scheme: scheme,
						Address: address,
						Port: port,
					}
				}
			}
		}
	}()

	return out
}

func main() {
	flagJSON := flag.Bool("json", false, "Output full data in JSON format")
	flagInteractive := flag.Bool("interactive", false, "Don't exit after completing the task and wait for more input")
	flagSkipCert := flag.Bool("skipcert", false, "Skip the TLS certificate verification")
	flagTimeout := flag.Int("timeout", 15, "Request timeout in seconds")
	flag.Parse()

	var waiting int
	resultQueue := make(chan *proxy)
	defer close(resultQueue)

	scanner := bufio.NewScanner(os.Stdin)

	for {
		if scanner.Scan() {
			for p := range requestGenerator(scanner.Text()) {
				waiting++
				go p.check(resultQueue, flagTimeout, flagSkipCert)
			}
		}

		for ; waiting > 0; waiting-- {
			result := <- resultQueue

			if *flagJSON {
				jr, _ := json.Marshal(*result)
				fmt.Println(string(jr))
			} else {
				if result.Good {
					fmt.Printf("%v://%v:%v\n", result.Scheme, result.Address, result.Port)
				}
			}
		}

		if !*flagInteractive {
			os.Exit(0)
		}
	}
}
