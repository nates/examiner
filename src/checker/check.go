package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"h12.io/socks"
)

func check(ip string, port string, proxyType string, proxyTimeout *int, proxyUrl *string, proxyText *string, proxyNotext *string) (int, error) {
	var text = ip
	var notext = ""
	if *proxyText != "" {
		text = *proxyText
	}
	if *proxyNotext != "" {
		notext = *proxyNotext
	}
	if proxyType == "socks5" {
		start := time.Now()
		dial := socks.Dial("socks5://" + ip + ":" + port)
		transport := &http.Transport{
			Dial: dial,
		}
		client := &http.Client{
			Transport: transport,
			Timeout:   time.Duration(*proxyTimeout) * time.Second,
		}
		response, err := client.Get(*proxyUrl)
		if err != nil {
			return 0, errors.New("Error requesting " + *proxyUrl)
		}
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return 0, errors.New("Error reading body.")
		}
		if notext != "" {
			if strings.Contains(string(body), notext) {
				return 0, errors.New("Invalid proxy.")
			}
		}
		if !strings.Contains(string(body), text) {
			return 0, errors.New("Invalid proxy.")
		}
		elapsed := time.Now().Sub(start)
		return int(elapsed / time.Millisecond), nil
	} else if proxyType == "socks4" {
		start := time.Now()
		dial := socks.Dial("socks4://" + ip + ":" + port)
		transport := &http.Transport{
			Dial: dial,
		}
		client := &http.Client{
			Transport: transport,
			Timeout:   time.Duration(*proxyTimeout) * time.Second,
		}
		response, err := client.Get(*proxyUrl)
		if err != nil {
			return 0, errors.New("Error requesting " + *proxyUrl)
		}
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return 0, errors.New("Error reading body.")
		}
		if notext != "" {
			if strings.Contains(string(body), notext) {
				return 0, errors.New("Invalid proxy.")
			}
		}
		if !strings.Contains(string(body), text) {
			return 0, errors.New("Invalid proxy.")
		}
		elapsed := time.Now().Sub(start)
		return int(elapsed / time.Millisecond), nil
	} else {
		start := time.Now()
		url := url.URL{}
		proxy, err := url.Parse("http://" + ip + ":" + port)
		if err != nil {
			return 0, errors.New("Error parsing proxy.")
		}
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
		client := &http.Client{
			Transport: transport,
			Timeout:   time.Duration(*proxyTimeout) * time.Second,
		}
		response, err := client.Get(*proxyUrl)
		if err != nil {
			return 0, errors.New("Error requesting " + *proxyUrl)
		}
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return 0, errors.New("Error reading body.")
		}
		if notext != "" {
			if strings.Contains(string(body), notext) {
				return 0, errors.New("Invalid proxy.")
			}
		}
		if !strings.Contains(string(body), text) {
			return 0, errors.New("Invalid proxy.")
		}
		elapsed := time.Now().Sub(start)
		return int(elapsed / time.Millisecond), nil
	}
}
