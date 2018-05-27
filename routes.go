/*
I like to have a single file inside every component called routes.go
where all the routing can live.  This is handy because most code
maintenance starts with a URL and an error report  —  so one glance at
routes.go will direct us where to look.
*/

package main

import (
	"github.com/valyala/fasthttp/expvarhandler"
)

func (s *server) Routes() {

	// Basics
	s.router.GET("/", s.index)
	s.router.GET("/long", s.longRunning)

	// Kubernetes
	s.router.GET("/healthz", s.healthz)
	s.router.GET("/readyz", s.readyz)

	// Monitoring
	s.router.GET("/stats", expvarhandler.ExpvarHandler)

}
