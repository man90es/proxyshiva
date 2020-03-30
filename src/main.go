package main

import (
	"fmt"
	"bufio"
	"log"
	"os"
	"flag"
)

func main() {
	outputJSON := flag.Bool("json", false, "Output in JSON format")
	goodOnly := flag.Bool("good", false, "Only output good proxies")
	persistent := flag.Bool("persistent", false, "Don't exit after completing the task and wait for new input")
	timeout := flag.Int("timeout", 15, "Request timeout in seconds")
	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)

	queue := make(chan KRequest, 100)
	
	for {
		data := make([][]string, 2)
		for i := range data {
			data[i] = make([]string, 0)
		}

		if scanner.Scan() {
			data = parseInput(scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			log.Println(err)
		}

		for _, address := range data[0] {
			for _, port := range data[1] {
				go check(address, port, queue, *timeout)
			}
		}

		for i := 0; i < len(data[0]) * len(data[1]); i++ {
			r := <-queue

			if *goodOnly && r.PingTime < 0 {
				continue
			}

			if *outputJSON {
				fmt.Printf("{\"address\": \"%s\", \"good\": %t, \"speed\": %v}\n", r.Url, r.PingTime > 0, r.PingTime)
			} else {
				fmt.Printf("%s %v\n", r.Url, r.PingTime)
			}
		}

		if !*persistent {
			break
		}
	}
}