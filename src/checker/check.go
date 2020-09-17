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
	url := url.URL{}
	proxy, err := url.Parse("http://" + ip + ":" + port)
	if err != nil {
		return "", 0, errors.New("Error parsing proxy.")
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxy),
	}
	client := &http.Client{
		Transport: transport,
		Timeout: 5 * time.Second,
	}
	response, err := client.Get("https://api.ipify.org")
	if err != nil {
		return "", 0, errors.New("Error requesting ipify.org.")
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", 0, errors.New("Error reading body.")
	}
	if string(body) != ip {
		return "", 0, errors.New("Invalid proxy.")
	}
	elapsed := time.Now().Sub(start)
	return string(body), int(elapsed / time.Millisecond), nil
}