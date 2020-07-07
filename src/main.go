package main

import (
	"fmt"
	"bufio"
	"os"
	"flag"
	"encoding/json"
	"time"
)

type KRequest struct {
	Scheme 		string 		`json:"scheme"`
	Address 	string 		`json:"address"`
	Port 		string 		`json:"port"`
	ExitAddress string 		`json:"exitAddress"`
	Good 		bool 		`json:"good"`
	Error 		string 		`json:"error"`
	Speed 		float64 	`json:"speed"`
}

func main() {
	schemes := []string{"http", "https", "socks4", "socks5"}

	flagV := flag.Bool("v", false, "Verbose output in JSON format")
	flagP := flag.Bool("p", false, "Don't exit after completing the task and wait for more input")
	flagT := flag.Int("t", 15, "Request timeout in seconds")
	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)
	resultQueue := make(chan KRequest, 100)
	taskQueue := make(chan KRequest, 100)

	for i := 0; i < 100; i++ {
		go checkerRoutine(taskQueue, resultQueue, flagT)
	}

	for {
		data := make([][]string, 2)
		requestCount := 0

		if scanner.Scan() {
			for i := range data {
				data[i] = make([]string, 0)
			}
			
			requestCount = schedule(scanner.Text(), schemes, taskQueue)
		} else {
			time.Sleep(2 * time.Second)
			continue
		}

		if *flagV {
			for i := 0; i < requestCount; i++ {
				r := <- resultQueue

				jr, _ := json.Marshal(r)
				fmt.Println(string(jr))
			}
		} else {
			for i := 0; i < requestCount; i++ {
				r := <- resultQueue

				if r.Good {
					fmt.Printf("%s://%s:%s\n", r.Scheme, r.Address, r.Port)
				}
			}
		}

		if !*flagP {
			break
		}
	}

	close(resultQueue)
	close(taskQueue)
}