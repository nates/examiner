# ProxyChecker
A multi-threaded [https | socks5 | socks4] proxy checker made in Go.
## Requirements
* Go (latest)

## Installation
```
go get -u "h12.io/socks"
go build ./src/checker
```

## Setup
Put your unchecked proxies in proxies.txt like below:
```
1.2.3.4:8080
1.2.3.4:80
1.2.3.4:443
```

## Usage
```
./checker.exe -threads=300 -type=[https | socks5 | socks4]
```
