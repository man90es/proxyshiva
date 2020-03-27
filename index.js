const spawn = require('child_process').spawn

exports.check = (addresses, ports) => {
	let shiva = spawn('./bin/proxyshiva', ['-json'])
	let responseStack = []

	return new Promise((resolve, reject) => {
		shiva.stdout.on('data', (data) => {
			responseStack.push(JSON.parse(data.toString()))

			if(responseStack.length == addresses.length * ports.length) {
				resolve(responseStack)
			}
		})

		shiva.stdin.setEncoding('utf-8')
		shiva.stdin.write(`${addresses.join(',')}:${ports.join(',')}\n`)
		shiva.stdin.end()
	})
}