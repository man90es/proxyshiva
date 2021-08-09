package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sync"

	"github.com/octoman90/proxyshiva/inputParser"
	"github.com/octoman90/proxyshiva/proxy"
)

var wg sync.WaitGroup

func validateScanned(str string) bool {
	ipRegex := `((?:(?:25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(?:25[0-5]|2[0-4]\d|[01]?\d\d?))`
	protocolRegex := `(http|https|socks4|socks5)`
	portRegex := `\d+`

	regex := regexp.MustCompile(`` +
		// Match protocol or protocol list
		`^(` + protocolRegex + `(,` + protocolRegex + `)*)+` +

		// Match "://"
		`:\/\/` +

		// Match IP or IP range
		ipRegex + `(-` + ipRegex + `)?` +

		// Match ":"
		`:` +

		// Match port or port range
		`(` + portRegex + `)(-` + portRegex + `)?$`,
	)

	return regex.MatchString(str)
}

func main() {
	flagJSON := flag.Bool("json", false, "Output full data in JSON format")
	flagInteractive := flag.Bool("interactive", false, "Don't exit after completing the task and wait for more input")
	flagSkipCert := flag.Bool("skipcert", false, "Skip the TLS certificate verification")
	flagTimeout := flag.Int("timeout", 15, "Request timeout in seconds")
	flagSkipRes := flag.Bool("skipres", false, "Skip reserved IP addresses")
	flag.Parse()

	resultQueue := make(chan *proxy.Proxy)
	defer close(resultQueue)

	// Receive and print out completed checks
	go func() {
		for result := range resultQueue {
			if *flagJSON { // Print out every result in JSON format
				jr, _ := json.Marshal(*result)
				fmt.Println(string(jr))
			} else if result.Good { // Print out good proxies in short format
				fmt.Printf("%v://%v:%v\n", result.Scheme, result.Address, result.Port)
			}

			wg.Done()
		}
	}()

	// Scan for input
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if scanner.Scan() {
			scanned := scanner.Text()
			if valid := validateScanned(scanned); !valid {
				continue
			}

			for proxy := range inputParser.RequestGenerator(scanned) {
				if !*flagSkipRes || !proxy.IsReserved() {
					wg.Add(1)
					go proxy.Check(resultQueue, flagTimeout, flagSkipCert)
				}
			}
		}

		if !*flagInteractive {
			wg.Wait()
			os.Exit(0)
		}
	}
}
