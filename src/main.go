package main

import (
	"fmt"
	"bufio"
	"log"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var data string; 
	
	for scanner.Scan() {
		data = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	queue := make(chan KRequest, 100)

	uris := strings.Split(strings.Split(data, ":")[0], ",")
	ports := strings.Split(strings.Split(data, ":")[1], ",")

	for _, port := range ports {
		for _, uri := range uris {
			go check(uri, port, queue)
		}
	}

	for i := 0; i < len(ports) * len(uris); i++ {
		r := <-queue
		if r.PingTime > 0 {
			fmt.Printf("%s %v\n", r.Url, r.PingTime)
		}
	}
}