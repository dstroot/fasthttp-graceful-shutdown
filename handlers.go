package main

import (
	"time"

	"github.com/valyala/fasthttp"
)

// NOTE: handlers hang off the server:
// func (s *server) handler(ctx *fasthttp.RequestCtx) { ... }
// That way handlers can access the dependencies via the s server variable.

var (
	done    = []byte(`{"done": true}`)
	ready   = []byte(`{"ready": true}`)
	healthy = []byte(`{"alive": true}`)
	html    = []byte(`<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <!-- Bootstrap CSS -->
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css" integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">
  </head>
  <body>
    <div class="container">
      <div class="row">
        <div class="col-md-4">
          <h1>Routes:</h1>
          <ul>
            <h4>
            <li><a href="/">/</a></li>
			<li><a href="/long">/long</a></li>
            <li><a href="/readyz">/readyz</a></li>
            <li><a href="/healthz">/heathz</a></li>
            <li><a href="/stats">/stats</a></li>
            </h4>
          </ul>
        </div>
      </div>
    </div>
  </body>
</html>`)
)

func (s *server) index(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/html")
	ctx.SetBody(html)
}

// longRunning is used to test our graceful shutdown
func (s *server) longRunning(ctx *fasthttp.RequestCtx) {
	time.Sleep(2 * time.Second) // simulate long process
	ctx.SetContentType("application/json")
	ctx.SetBody(done)
}

// healthz supports a health probe. It is a simple handler which
// always return response code 200 @ /healthz
func (s *server) healthz(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	ctx.SetBody(healthy)
}

// readyz supports a readiness probe.  For the readiness probe we might
// need to wait for some event (e.g. the database is ready) to be able
// to serve traffic. Returns 200 @ readyz
func (s *server) readyz(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	ctx.SetBody(ready)

	// NOTE: A lot of additional useful info is exposed to request handler:
	// fmt.Fprintf(ctx, "Request method is %q\n", ctx.Method())
	// fmt.Fprintf(ctx, "Request ID is %v\n", ctx.ID())
	// fmt.Fprintf(ctx, "RequestURI is %q\n", ctx.RequestURI())
	// fmt.Fprintf(ctx, "Requested path is %q\n", ctx.Path())
	// fmt.Fprintf(ctx, "Host is %q\n", ctx.Host())
	// fmt.Fprintf(ctx, "Query string is %q\n", ctx.QueryArgs())
	// fmt.Fprintf(ctx, "User-Agent is %q\n", ctx.UserAgent())
	// fmt.Fprintf(ctx, "Connection has been established at %s\n", ctx.ConnTime())
	// fmt.Fprintf(ctx, "Request has been started at %s\n", ctx.Time())
	// fmt.Fprintf(ctx, "Serial request number for the current connection is %d\n", ctx.ConnRequestNum())
	// fmt.Fprintf(ctx, "Your ip is %q\n\n", ctx.RemoteIP())
	// fmt.Fprintf(ctx, "Raw request is:\n---CUT---\n%s\n---CUT---", &ctx.Request)
}
