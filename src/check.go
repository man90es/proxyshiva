package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type KRequest struct {
	Url 		string
	PingTime 	float64
}

var defaultURL = "https://ipecho.net/plain"

func check(ip string, port string, c chan KRequest) {
	var timeout = time.Duration(15 * time.Second)

	startAt := time.Now()
	host := fmt.Sprintf("%s:%s", ip, port)

	proxyUrl := &url.URL{Host: host}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
		Timeout: timeout,
	}

	response, err := client.Get(defaultURL)

	if err != nil {
		c <- KRequest{
			Url:      host,
			PingTime: float64(-1),
		}
		return
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	delta := time.Now().UnixNano() - startAt.UnixNano()

	if strings.Contains(string(body), ip) {
		c <- KRequest{
			Url:      host,
			PingTime: float64(delta) / 1e9,
		}
	} else {
		c <- KRequest{
			Url:      host,
			PingTime: float64(-1),
		}
	}
}