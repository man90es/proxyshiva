# ProxyShiva

This is a tool for checking proxies availability

## To-do:

* JSON output format
* Wrapper to use checker from Node.js
* Listening to stdin continuously without exiting
* SOCKS proxies support

## Building from source

You'll need a Golang compiler installed for this.

```bash
make build
```

## Usage from console
```bash
echo "0.0.0.0:8080" | ./bin/proxyshiva
```
or
```bash
echo "0.0.0.0,1.1.1.1:80,8080,1080" | ./bin/proxyshiva > good.txt
```
etc.

## License
[MIT](https://choosealicense.com/licenses/mit/)