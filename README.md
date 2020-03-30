# ProxyShiva

This is a tool for checking proxies availability

## To-do:

* ~~JSON output format~~
* ~~Option to only output good proxies~~
* ~~Wrapper to use checker from Node.js~~
* ~~Listen to stdin continuously without exiting~~
* ~~Take timeout as a parameter~~
* SOCKS proxies support
* Wildcard input

## Building from source

You'll need a Golang compiler installed for this.

```bash
make build
```

## Usage from console
```bash
$ echo "0.0.0.0:8080" | ./bin/proxyshiva
```
or give multiple values and output to file 
```bash
$ echo "0.0.0.0,1.1.1.1:80,8080,1080" | ./bin/proxyshiva > good.txt
```
or use interactive mode
```bash
$ cat | ./bin/proxyshiva -persistent
0.0.0.0:8080
```
etc.

## Usage as Node.js library
```nodejs
const proxyShiva = require('./proxyshiva')

async function main() {
	console.log(await proxyShiva.check(['0.0.0.0','1.1.1.1'], ['80', '8080']))
}

main()
```
etc.

## License
[MIT](https://choosealicense.com/licenses/mit/)