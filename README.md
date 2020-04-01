# ProxyShiva

This is a tool for checking proxies availability

## Building from source

You'll need a Golang compiler installed for this.

```bash
make build
```

## Usage from console
Pipe addresses into shiva using echo or cat command:
```bash
$ echo "192.168.0.2:8080" | ./bin/proxyshiva -v -t 30
$ cat addresses.txt | ./bin/proxyshiva -v -t 30
```
To input several addresses and/or ports, separate them with comma:
```bash
$ echo "192.168.0.1,192.168.0.2:80,8080" | ./bin/proxyshiva
```
To check all addresses and/or ports in range, use dash:
```bash
$ echo "192.168.0.1-192.168.1.0:80-90" | ./bin/proxyshiva
```

## Flags
```
-v 		Verbose output in JSON format
-r 		Randomize check order
-p 		Interactive mode
-t 	15 	Request timeout in seconds
```

## Usage as Node.js library
```nodejs
const proxyShiva = require('./proxyshiva')

async function main() {
	console.log(await proxyShiva.check(['0.0.0.0','1.1.1.1'], ['80', '8080']))
}

main()
```

## To-do:

* Split ranges into smaller ones if too big
* IPv6 proxies support

## License
[MIT](https://choosealicense.com/licenses/mit/)