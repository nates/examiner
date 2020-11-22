package main

import (
	"bufio"
	"errors"
	"os"
)

var pool = []string{}

func fillPool(path string) error {
	file, err := os.Open(path)

	if err != nil {
		return errors.New("could not open " + path)
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
	if len(pool) == 0 {
		return "", errors.New("pool is empty")
	}

	proxy := pool[len(pool)-1]
	pool[len(pool)-1] = ""
	pool = pool[:len(pool)-1]

	return proxy, nil
}
