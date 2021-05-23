package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	working []string
	https   []string
	socks4  []string
	socks5  []string

	threads   int    = 100
	timeout   int    = 5
	checkURL  string = "https://api.ipify.org"
	text      string
	notext    string
	proxyType string = "https"
	input     string = "proxies.txt"
	output    string = "working.txt"
)

func main() {
	flag.IntVar(&threads, "threads", threads, "Amount of threads to check with.")
	flag.IntVar(&timeout, "timeout", timeout, "Timeout in seconds.")
	flag.StringVar(&checkURL, "url", checkURL, "URL to check proxies with. (Requires text option)")
	flag.StringVar(&text, "text", text, "If this text is found on the page, the proxy will be marked good. (If left empty it will default to the proxy IP address.)")
	flag.StringVar(&notext, "notext", notext, "If this text is found on the page, the proxy will be marked invalid.")
	flag.StringVar(&proxyType, "type", proxyType, "Type of proxies [https | socks5 | socks4]")
	flag.StringVar(&input, "input", input, "File to check")
	flag.StringVar(&output, "output", output, "File to output proxies")
	flag.Parse()

	fillPool(input)

	types := strings.Split(proxyType, ",")

	if len(types) > 3 {
		log.Println("Invalid amount of types.")
		return
	}

	for i := 0; i < len(types); i++ {
		if types[i] != "https" && types[i] != "socks4" && types[i] != "socks5" {
			log.Println("Invalid type of proxy.")
			return
		}
	}

	if threads <= 0 {
		log.Println("Invalid amount of threads.")
		return
	}

	var wg sync.WaitGroup

	for i := 1; i <= threads; i++ {
		wg.Add(1)
		go worker(i, &wg, types, timeout, checkURL, text, notext)
	}

	wg.Wait()

	if len(types) == 1 {
		file, err := os.Create(output)
		if err != nil {
			log.Println(err)
			return
		}
		_, err = file.WriteString(strings.Join(working, "\n"))
		if err != nil {
			log.Println(err)
			file.Close()
			return
		}
		log.Println("Wrote proxies to " + output)
		err = file.Close()
		if err != nil {
			log.Println(err)
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
				log.Println(err)
				return
			}
			_, err = file.WriteString(strings.Join(array, "\n"))
			if err != nil {
				log.Println(err)
				file.Close()
				return
			}
			log.Println("Wrote proxies to " + filename)
			err = file.Close()
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func worker(id int, wg *sync.WaitGroup, proxyType []string, timeout int, url string, text string, notext string) {
	proxy, err := getProxy()
	if err != nil {
		if err.Error() == "pool is empty" {
			wg.Done()
			return
		}
	}
	proxySplit := strings.Split(proxy, ":")
	if len(proxySplit) != 2 {
		log.Println("["+strconv.Itoa(id)+"]", "Invalid proxy.")
		worker(id, wg, proxyType, timeout, url, text, notext)
		return
	}
	for _, v := range proxyType {
		speed, err := check(proxySplit[0], proxySplit[1], v, timeout, url, text, notext)
		if err != nil {
			log.Println("["+strconv.Itoa(id)+"]", err.Error())
		} else {
			log.Println("["+strconv.Itoa(id)+"]", "Working proxy, Speed: "+strconv.Itoa(speed)+"ms")
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
