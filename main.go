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

func worker(taskQueue <-chan proxy.Proxy, resultQueue chan<- *proxy.Proxy, timeout *int, skipCert *bool) {
	for proxy := range taskQueue {
		proxy.Check(resultQueue, timeout, skipCert)
	}
}

func printer(resultQueue <-chan *proxy.Proxy, JSON *bool) {
	for result := range resultQueue {
		if *JSON { // Print out every result in JSON format
			jr, _ := json.Marshal(*result)
			fmt.Println(string(jr))
		} else if result.Good { // Print out good proxies in short format
			fmt.Printf("%v://%v:%v\n", result.Scheme, result.Address, result.Port)
		}

		wg.Done()
	}
}

func main() {
	// Parse flags
	flagJSON := flag.Bool("json", false, "Output full data in JSON format")
	flagInteractive := flag.Bool("interactive", false, "Don't exit after completing the task and wait for more input")
	flagSkipCert := flag.Bool("skipcert", false, "Skip the TLS certificate verification")
	flagTimeout := flag.Int("timeout", 15, "Request timeout in seconds")
	flagSkipRes := flag.Bool("skipres", false, "Skip reserved IP addresses")
	flagParallel := flag.Int("parallel", 100, "How many requests to make simultaneously")
	flag.Parse()

	// Create communication queues
	taskQueue := make(chan proxy.Proxy, 100)
	resultQueue := make(chan *proxy.Proxy, 100)
	defer close(taskQueue)
	defer close(resultQueue)

	// Spawn a printer goroutine
	go printer(resultQueue, flagJSON)

	// Spawn worker goroutines
	for i := 0; i < *flagParallel; i++ {
		go worker(taskQueue, resultQueue, flagTimeout, flagSkipCert)
	}

	// Scan for input
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if scanner.Scan() {
			for proxy := range inputParser.RequestGenerator(scanner.Text()) {
				// Copy tasks from individual input channels to a single task queue
				if !*flagSkipRes || !proxy.IsReserved() {
					wg.Add(1)
					taskQueue <- proxy
				}
			}
		}

		if !*flagInteractive {
			wg.Wait()
			os.Exit(0)
		}
	}
}
