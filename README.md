# ProxyShiva

ProxyShiva is a concurrent command-line proxy checker application.

## Getting the binary
### Prebuilt
You can download prebuilt binaries [here](https://gitlab.com/man90/ProxyShiva/-/pipelines).

### Building from source
Prerequisites: <abbr title="Not tested on other platforms.">GNU/Linux</abbr>, Go >=1.15

Command:
```bash
go build
```

## Usage
Use Shiva's standard input stream to check proxies:
```bash
$ echo "socks5://127.0.0.1:9050" | ./ProxyShiva
```

Put the output into a file:
```bash
$ echo "socks5://127.0.0.1:9050" | ./ProxyShiva > good.txt
```

Use scheme lists, IP ranges and port ranges to check more proxies at once:
```bash
$ echo "http,https,socks5://192.168.0.2-192.168.0.10:8080-8089" | ./ProxyShiva
```

Put files into the standard input stream for more convenience and flexibility:
```bash
$ touch addresses.txt
$ echo "socks5://127.0.0.1:9050" >> addresses.txt
$ echo "http,https,socks5://192.168.0.2-192.168.0.10:8080-8089" >> addresses.txt
$ cat addresses.txt | ./ProxyShiva
```

## Flags

| Flag | Description |
| ------ | ------ |
| -json | Output full data in JSON format |
| -interactive | Don't exit after completing the task and wait for more input |
| -skipcert | Skip the TLS certificate verification |
| -timeout | Request timeout in seconds (15 by default) |

## License
[MIT](https://choosealicense.com/licenses/mit/)
