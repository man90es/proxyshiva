package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var defaultURL = "https://ipecho.net/plain"

func check(scheme string, ip string, port string, c chan KRequest, sTimeout int) {
	k := KRequest{
		Scheme: 	scheme,
		Address: 	ip,
		Port: 		port,
	}

	startAt := time.Now()
	host := fmt.Sprintf("%s:%s", ip, port)
	proxyUrl := &url.URL{Scheme: scheme, Host: host}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
		Timeout: time.Duration(time.Duration(sTimeout) * time.Second),
	}

	response, err := client.Get(defaultURL)

	if err == nil {
		k.Good = true
		k.Speed = float64(time.Now().UnixNano() - startAt.UnixNano()) / 1e9

		body, _ := ioutil.ReadAll(response.Body)
		defer response.Body.Close()
		
		if !strings.Contains(string(body), ip) {
			k.ExitAddress = string(body)
		} else {
			k.ExitAddress = ip
		}
	} else {
		k.Good = false
		k.Error = err.Error()
	}

	c <- k
}