package main

import (
	"fmt"
	"bufio"
	"log"
	"os"
	"flag"
	"encoding/json"
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
	schemes := [4]string{"http", "https", "socks4", "socks5"}

	flagV := flag.Bool("v", false, "Verbose output in JSON format")
	flagP := flag.Bool("p", false, "Don't exit after completing the task and wait for more input")
	flagT := flag.Int("t", 15, "Request timeout in seconds")
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
				for _, scheme := range schemes {
					go check(scheme, address, port, queue, *flagT)
				}
			}
		}

		for i := 0; i < len(data[0]) * len(data[1]) * len(schemes); i++ {
			r := <- queue

			if *flagV {
				jr, _ := json.Marshal(r)
				fmt.Println(string(jr))
			} else {
				if r.Good {
					fmt.Printf("%s://%s:%s\n", r.Scheme, r.Address, r.Port)
				}
			}
		}

		if !*flagP {
			break
		}
	}
}