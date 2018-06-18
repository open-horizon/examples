// Locally test msg-receive.js

const msgreceive = require('./msgreceive.js')

params = {
	"messages": [
		{ "value": "this is my first msg" },
		{ "value": "this is my 2nd msg" }
	]
}

// const result = msgreceive.main(params)
// console.log("msgreceive.main() result:")
// console.log(result)

msgreceive.main(params).then(function(response){
	console.log("msgreceive.main() result:")
	console.log(response)
}, function(error) {
    console.log(error.message);
});
