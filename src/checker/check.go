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

func check(ip string, port string, proxyType string, proxyTimeout int, proxyURL string, proxyText string, proxyNotext string) (int, error) {
	var text = ip
	var notext = ""
	if proxyText != "" {
		text = proxyText
	}
	if proxyNotext != "" {
		notext = proxyNotext
	}
	if proxyType == "socks5" {
		start := time.Now()
		dial := socks.Dial("socks5://" + ip + ":" + port)
		transport := &http.Transport{
			Dial: dial,
		}
		client := &http.Client{
			Transport: transport,
			Timeout:   time.Duration(proxyTimeout) * time.Second,
		}
		response, err := client.Get(proxyURL)
		if err != nil {
			return 0, errors.New("Error requesting " + proxyURL)
		}
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return 0, errors.New("could not read body")
		}
		if notext != "" {
			if strings.Contains(string(body), notext) {
				return 0, errors.New("bad proxy")
			}
		}
		if !strings.Contains(string(body), text) {
			return 0, errors.New("bad proxy")
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
			Timeout:   time.Duration(proxyTimeout) * time.Second,
		}
		response, err := client.Get(proxyURL)
		if err != nil {
			return 0, errors.New("Error requesting " + proxyURL)
		}
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return 0, errors.New("could not read body")
		}
		if notext != "" {
			if strings.Contains(string(body), notext) {
				return 0, errors.New("bad proxy")
			}
		}
		if !strings.Contains(string(body), text) {
			return 0, errors.New("bad proxy")
		}
		elapsed := time.Now().Sub(start)
		return int(elapsed / time.Millisecond), nil
	} else {
		start := time.Now()
		url := url.URL{}
		proxy, err := url.Parse("http://" + ip + ":" + port)
		if err != nil {
			return 0, errors.New("bad proxy")
		}
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
		client := &http.Client{
			Transport: transport,
			Timeout:   time.Duration(proxyTimeout) * time.Second,
		}
		response, err := client.Get(proxyURL)
		if err != nil {
			return 0, errors.New("Error requesting " + proxyURL)
		}
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return 0, errors.New("could not read body")
		}
		if notext != "" {
			if strings.Contains(string(body), notext) {
				return 0, errors.New("bad proxy")
			}
		}
		if !strings.Contains(string(body), text) {
			return 0, errors.New("bad proxy")
		}
		elapsed := time.Now().Sub(start)
		return int(elapsed / time.Millisecond), nil
	}
}
