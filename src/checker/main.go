package main

import (
	"fmt"
	"strings"
	"sync"
	"strconv"
	"os"
	"flag"
)

var(
	working = []string{}
)

func main() {
	threads := flag.Int("threads", 100, "Amount of threads to check with.")
	timeout := flag.Int("timeout", 5, "Timeout in seconds.")
	url := flag.String("url", "https://api.ipify.org", "URL to check proxies with. (Requires text option)")
	text := flag.String("text", "", "If this text is found on the page, the proxy will be marked good. (IF left empty it will default to the proxy IP address.)")
	proxyType := flag.String("type", "https", "Type of proxies [https | socks5 | socks4]")
	input := flag.String("input", "proxies.txt", "File to check")
	output := flag.String("output", "working.txt", "File to output proxies")
	flag.Parse()

	fillPool(input)

	if *threads <= 0 {
		fmt.Println("Invalid amount of threads.")
		return
	}
 
	var wg sync.WaitGroup

	for i := 1; i <= *threads; i++ {
		wg.Add(1)
		go worker(i, &wg, proxyType, timeout, url, text)
	}

	wg.Wait()

	fmt.Println(strconv.Itoa(len(working)) + " working proxies.")
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
}

func worker(id int, wg *sync.WaitGroup, proxyType *string, timeout *int, url *string, text *string) {
	proxy, err := getProxy()
	if err != nil {
		if(err.Error() == "Pool is empty.") {
			wg.Done()
			return
		}
	}
	proxySplit := strings.Split(proxy, ":")
	if len(proxySplit) != 2 {
		fmt.Println("[" + strconv.Itoa(id) + "]", "Invalid proxy.")
		worker(id, wg, proxyType, timeout, url, text)
		return
	}
	speed, err := check(proxySplit[0], proxySplit[1], proxyType, timeout, url, text)
	if err != nil {
		fmt.Println("[" + strconv.Itoa(id) + "]", err.Error())
		worker(id, wg, proxyType, timeout, url, text)
		return
	}
	fmt.Println("[" + strconv.Itoa(id) + "]", "Working proxy, Speed: " + strconv.Itoa(speed) + "ms")
	working = append(working, proxy)
	worker(id, wg, proxyType, timeout, url, text)
	return
}