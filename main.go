package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/octoman90/proxyshiva/inputParser"
	"github.com/octoman90/proxyshiva/proxy"
)

var wg sync.WaitGroup

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
			for proxy := range inputParser.RequestGenerator(scanner.Text()) {
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
