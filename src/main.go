package main

import (
	"fmt"
	"bufio"
	"log"
	"os"
	"strings"
	"flag"
)

func main() {
	outputJSON := flag.Bool("json", false, "Output in JSON format")
	goodOnly := flag.Bool("good", false, "Only output good proxies")
	persistent := flag.Bool("persistent", false, "Don't exit after completing the task and wait for new input")
	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)

	var data string;
	queue := make(chan KRequest, 100)
	
	for {
		if scanner.Scan() {
			data = scanner.Text()
		}

		if err := scanner.Err(); err != nil {
			log.Println(err)
		}

		uris := strings.Split(strings.Split(data, ":")[0], ",")
		ports := strings.Split(strings.Split(data, ":")[1], ",")

		for _, port := range ports {
			for _, uri := range uris {
				go check(uri, port, queue)
			}
		}

		for i := 0; i < len(ports) * len(uris); i++ {
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