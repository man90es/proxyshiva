# ProxyShiva
[![license](https://img.shields.io/github/license/octoman90/proxyshiva)](https://github.com/octoman90/proxyshiva/blob/master/LICENSE)

ProxyShiva is a concurrent command-line proxy checker application.

## Building from source
Prerequisites: Go >=1.15
```bash
git clone https://github.com/octoman90/proxyshiva.git
cd proxyshiva
go build
```

## Usage
Use Shiva's standard input stream to check proxies:
```bash
echo "socks5://127.0.0.1:9050" | ./proxyshiva
```

Put the output into a file:
```bash
echo "socks5://127.0.0.1:9050" | ./proxyshiva > good.txt
```

Use scheme lists, IP ranges and port ranges to check more proxies at once:
```bash
echo "http,https,socks5://192.168.0.2-192.168.0.10:8080-8089" | ./proxyshiva
```

Put files into the standard input stream for more convenience and flexibility:
```bash
touch addresses.txt
echo "socks5://127.0.0.1:9050" >> addresses.txt
echo "http,https,socks5://192.168.0.2-192.168.0.10:8080-8089" >> addresses.txt
cat addresses.txt | ./proxyshiva
```

## Flags

| Flag         | Description                                                  |
| ------------ | ------------------------------------------------------------ |
| -json        | Output full data in JSON format                              |
| -interactive | Don't exit after completing the task and wait for more input |
| -skipcert    | Skip the TLS certificate verification                        |
| -skipres     | Skip reserved IP addresses                                   |
| -timeout=15  | Request timeout in seconds (15 by default)                   |
