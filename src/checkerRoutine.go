package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultURL = "https://ipecho.net/plain"

func checkerRoutine(taskQueue chan KRequest, resultQueue chan KRequest, timeout *int) {
	for {
		k := <- taskQueue

		startAt := time.Now()
		proxyUrl := &url.URL{Scheme: k.Scheme, Host: fmt.Sprintf("%s:%s", k.Address, k.Port)}

		client := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},

			Timeout: time.Duration(time.Duration(*timeout) * time.Second),
		}

		response, err := client.Get(defaultURL)

		if err == nil {
			k.Good = true
			k.Speed = float64(time.Now().UnixNano() - startAt.UnixNano()) / 1e9

			body, _ := ioutil.ReadAll(response.Body)
			defer response.Body.Close()
			
			if !strings.Contains(string(body), k.Address) {
				k.ExitAddress = string(body)
			} else {
				k.ExitAddress = k.Address
			}
		} else {
			k.Good = false
			k.Error = err.Error()
		}

		resultQueue <- k
	}
}