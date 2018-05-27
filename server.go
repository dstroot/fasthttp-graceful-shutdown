package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
)

type server struct {
	HTTPServer *fasthttp.Server
	router     *fasthttprouter.Router
}

// NewServer creates a new HTTP Server
func newServer() *server {

	// define router
	r := fasthttprouter.New()

	// compression
	h := r.Handler
	if cfg.Compress {
		h = fasthttp.CompressHandler(h)
	}

	return &server{
		HTTPServer: newHTTPServer(h),
		router:     r,
	}
}

// NewServer creates a new HTTP Server
// TODO: configuration should be configurable
func newHTTPServer(h fasthttp.RequestHandler) *fasthttp.Server {
	return &fasthttp.Server{
		Handler:              h,
		ReadTimeout:          5 * time.Second,
		WriteTimeout:         10 * time.Second,
		MaxConnsPerIP:        500,
		MaxRequestsPerConn:   500,
		MaxKeepaliveDuration: 5 * time.Second,
	}
}

// Run starts the HTTP server and performs a graceful shutdown
func (s *server) Run() {
	// NOTE: Package reuseport provides a TCP net.Listener with SO_REUSEPORT support.
	// SO_REUSEPORT allows linear scaling server performance on multi-CPU servers.

	// create a fast listener ;)
	ln, err := reuseport.Listen("tcp4", "localhost:"+cfg.Port)
	if err != nil {
		log.Fatalf("error in reuseport listener: %s", err)
	}

	// create a graceful shutdown listener
	duration := 5 * time.Second
	graceful := NewGracefulListener(ln, duration)

	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("hostname unavailable: %s", err)
	}

	// Error handling
	listenErr := make(chan error, 1)

	/// Run server
	go func() {
		log.Printf("%s - Web server starting on port %v", hostname, graceful.Addr())
		log.Printf("%s - Press Ctrl+C to stop", hostname)
		// listenErr <- s.HTTPServer.ListenAndServe(":" + cfg.Port)
		listenErr <- s.HTTPServer.Serve(graceful)

	}()

	// SIGINT/SIGTERM handling
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	// Handle channels/graceful shutdown
	for {
		select {
		// If server.ListenAndServe() cannot start due to errors such
		// as "port in use" it will return an error.
		case err := <-listenErr:
			if err != nil {
				log.Fatalf("listener error: %s", err)
			}
			os.Exit(0)
		// handle termination signal
		case <-osSignals:
			fmt.Printf("\n")
			log.Printf("%s - Shutdown signal received.\n", hostname)

			// Servers in the process of shutting down should disable KeepAlives
			// FIXME: This causes a data race
			s.HTTPServer.DisableKeepalive = true

			// Attempt the graceful shutdown by closing the listener
			// and completing all inflight requests.
			if err := graceful.Close(); err != nil {
				log.Fatalf("error with graceful close: %s", err)
			}

			log.Printf("%s - Server gracefully stopped.\n", hostname)
		}
	}
}
