package main

import(
	"bufio"
	"errors"
	"os"
)

var(
	pool = []string{}
)

func fill() error {
	file, err := os.Open("proxies.txt")
 
	if err != nil {
		return errors.New("Error reading proxies.txt.")
	}
 
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
 
	for scanner.Scan() {
		pool = append(pool, scanner.Text())
	}
 
	file.Close()

	return nil
}

func getProxy() (string, error) {
	i := len(pool) - 1

	if(i <= 0) {
		return "", errors.New("Pool is empty.")
	}

	proxy := pool[i]

	pool[i] = pool[len(pool)-1]
	pool[len(pool)-1] = ""
	pool = pool[:len(pool)-1]

	return proxy, nil
}