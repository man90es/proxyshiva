const spawn = require('child_process').spawn

const shiva = spawn('./bin/proxyshiva', ['-v', '-p'¸ '-r'])
shiva.stdin.setEncoding('utf-8')

check = (addresses, ports) => {
	let responseStack = []

	return new Promise((resolve, reject) => {
		shiva.stdout.on('data', (data) => {
			responseStack.push(JSON.parse(data.toString()))

			if(responseStack.length == addresses.length * ports.length) {
				resolve(responseStack)
			}
		})

		shiva.stdin.write(`${addresses.join(',')}:${ports.join(',')}\n`)
	})
}