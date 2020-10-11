package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	working = []string{}
	https   = []string{}
	socks4  = []string{}
	socks5  = []string{}
)

func main() {
	threads := flag.Int("threads", 100, "Amount of threads to check with.")
	timeout := flag.Int("timeout", 5, "Timeout in seconds.")
	url := flag.String("url", "https://api.ipify.org", "URL to check proxies with. (Requires text option)")
	text := flag.String("text", "", "If this text is found on the page, the proxy will be marked good. (If left empty it will default to the proxy IP address.)")
	notext := flag.String("notext", "", "If this text is found on the page, the proxy will be marked invalid.")
	proxyType := flag.String("type", "https", "Type of proxies [https | socks5 | socks4]")
	input := flag.String("input", "proxies.txt", "File to check")
	output := flag.String("output", "working.txt", "File to output proxies")
	flag.Parse()

	fillPool(input)

	types := strings.Split(*proxyType, ",")

	if len(types) > 3 {
		fmt.Println("Invalid amount of types.")
		return
	}

	for i := 0; i < len(types); i++ {
		if types[i] != "https" && types[i] != "socks4" && types[i] != "socks5" {
			fmt.Println("Invalid type of proxy.")
			return
		}
	}

	if *threads <= 0 {
		fmt.Println("Invalid amount of threads.")
		return
	}

	var wg sync.WaitGroup

	for i := 1; i <= *threads; i++ {
		wg.Add(1)
		go worker(i, &wg, types, timeout, url, text, notext)
	}

	wg.Wait()

	if len(types) == 1 {
		file, err := os.Create(*output)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = file.WriteString(strings.Join(working, "\n"))
		if err != nil {
			fmt.Println(err)
			file.Close()
			return
		}
		fmt.Println("Wrote proxies to " + *output)
		err = file.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		for _, v := range types {
			filename := ""
			array := []string{}
			if v == "https" {
				array = https
				filename = "https.txt"
			}
			if v == "socks5" {
				array = socks5
				filename = "socks5.txt"
			}
			if v == "socks4" {
				array = socks4
				filename = "socks4.txt"
			}
			file, err := os.Create(filename)
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = file.WriteString(strings.Join(array, "\n"))
			if err != nil {
				fmt.Println(err)
				file.Close()
				return
			}
			fmt.Println("Wrote proxies to " + filename)
			err = file.Close()
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func worker(id int, wg *sync.WaitGroup, proxyType []string, timeout *int, url *string, text *string, notext *string) {
	proxy, err := getProxy()
	if err != nil {
		if err.Error() == "Pool is empty." {
			wg.Done()
			return
		}
	}
	proxySplit := strings.Split(proxy, ":")
	if len(proxySplit) != 2 {
		fmt.Println("["+strconv.Itoa(id)+"]", "Invalid proxy.")
		worker(id, wg, proxyType, timeout, url, text, notext)
		return
	}
	for _, v := range proxyType {
		speed, err := check(proxySplit[0], proxySplit[1], v, timeout, url, text, notext)
		if err != nil {
			fmt.Println("["+strconv.Itoa(id)+"]", err.Error())
		} else {
			fmt.Println("["+strconv.Itoa(id)+"]", "Working proxy, Speed: "+strconv.Itoa(speed)+"ms")
			if len(proxyType) == 1 {
				working = append(working, proxy)
			} else if v == "https" {
				https = append(https, proxy)
			} else if v == "socks4" {
				socks4 = append(socks4, proxy)
			} else if v == "socks5" {
				socks5 = append(socks5, proxy)
			}
		}
	}
	worker(id, wg, proxyType, timeout, url, text, notext)
	return
}
