package main

import (
	"net/http"
	"net/url"
	"io/ioutil"
	"errors"
	"time"
)

func check(ip string, port string) (string, int, error) {
	start := time.Now()
	url_i := url.URL{}
	proxy, err := url_i.Parse("http://" + ip + ":" + port)
	if err != nil {
		return "", 0, errors.New("Error parsing proxy.")
	}
	transport := &http.Transport{}    
	transport.Proxy = http.ProxyURL(proxy)
	transport.IdleConnTimeout = 5
	client := &http.Client{
		Transport: transport,
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get("https://api.ipify.org")
	if err != nil {
		return "", 0, errors.New("Error requesting ipify.org.")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0, errors.New("Error reading body.")
	}
	if string(body) != ip {
		return "", 0, errors.New("Invalid proxy.")
	}
	elapsed := time.Now().Sub(start)
	return string(body), int(elapsed / 1e6), nil
}