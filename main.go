package main

func main() {
	initialize()

	// setup and run server
	s := newServer()
	s.Routes()
	s.Run()
}
