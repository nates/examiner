package main

import (
	"fmt"
	"strings"
	"regexp"
	"sync"
	"strconv"
	"os"
	"flag"
)

var(
	working = []string{}
)

func main() {
	fillPool()

	threads := flag.Int("threads", 100, "Amount of threads to check with.")
	proxyType := flag.String("type", "https", "Type of proxies [https | socks5 | socks4]")
	flag.Parse()
	if *threads <= 0 {
		fmt.Println("Invalid amount of threads.")
		return
	}
 
	var wg sync.WaitGroup

	for i := 0; i <= *threads; i++ {
		wg.Add(1)
		go worker(i, &wg, proxyType)
	}

	wg.Wait()

	fmt.Println(strconv.Itoa(len(working)) + " working proxies.")
	file, err := os.Create("working.txt")
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
	fmt.Println("Wrote proxies to working.txt")
    err = file.Close()
    if err != nil {
        fmt.Println(err)
        return
    }
}

func worker(id int, wg *sync.WaitGroup, proxyType *string) {
	proxy, err := getProxy()
	if err != nil {
		if(err.Error() == "Pool is empty.") {
			wg.Done()
			return
		}
	}
	regexMatch, err := regexp.MatchString(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}):(\d{1,5})`, proxy)
	if err != nil {
		fmt.Println("[" + strconv.Itoa(id) + "]", err.Error())
		worker(id, wg, proxyType)
		return
	}
	if regexMatch != true {
		fmt.Println("[" + strconv.Itoa(id) + "]", "Invalid IP.")
		worker(id, wg, proxyType)
		return
	}
	_, speed, err := check(strings.Split(proxy, ":")[0], strings.Split(proxy, ":")[1], proxyType)
	if err != nil {
		fmt.Println("[" + strconv.Itoa(id) + "]", err.Error())
		worker(id, wg, proxyType)
		return
	}
	fmt.Println("[" + strconv.Itoa(id) + "]", "Working proxy, Speed: " + strconv.Itoa(speed) + "ms")
	working = append(working, proxy)
	worker(id, wg, proxyType)
	return
}