This example application is built on fasthttp and fasthttprouter - it is designed to be exceptionally fast, and use very little memory.  It is also meant to be an example of a graceful shutdown with fasthttp.  Your review, feedback and pull requests are welcome.  Cheers!

Note: assumes you have Dep installed.

To run clone the repo and run `dep ensure` then `go run $(ls *.go | grep -v _test.go)` (assuming I write tests at some point)

Resources:

* https://github.com/valyala/fasthttp
* https://github.com/buaazp/fasthttprouter
