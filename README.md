# ProxyShiva

This is a tool for checking proxies availability

## To-do:

* ~~JSON output format~~
* ~~Option to only output good proxies~~
* Wrapper to use checker from Node.js
* Listen to stdin continuously without exiting
* SOCKS proxies support
* Take timeout as a parameter
* SOCKS proxies support
* Wildcard input

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