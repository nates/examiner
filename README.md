# 💻 examiner
A multi-threaded proxy checker made in Go. 
## 🧳 Requirements
* Go (latest)

## 🏗️ Building
```
go get -u "h12.io/socks"
go build ./src/examiner
```

## 🕹️ Usage
Put your proxies into proxies.txt separated by each line.
```
./examiner.exe -threads=300 -type=[https | socks5 | socks4] -url=https://google.com/ -text=google -timeout=5
```
