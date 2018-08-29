package main

import (
	//"crypto/tls"
	"log"
	"net"
	"runtime"

	flag "github.com/spf13/pflag"
)

var bindAddr string
var backendAddr string
var certFile string
var keyFile string

func init() {
	flag.StringVar(&bindAddr, "l", ":443", "bind address")
	flag.StringVar(&backendAddr, "b", ":8082", "backend address")
	flag.StringVar(&certFile, "c", "cert.pem", "TLS certificate path")
	flag.StringVar(&keyFile, "k", "key.pem", "TLS key path")
}

func main() {
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Printf("Handling connections on %s", bindAddr)
	log.Printf("Proxying connections to %s", backendAddr)

	// bind a port to handle TLS connections
	l, err := net.Listen("tcp", bindAddr)
	if err != nil {
		log.Printf("Failed to bind on address %s: %s", bindAddr, err)
	}
	log.Printf("Serving connections on %v", l.Addr())

	for {
		// accept next connection to this frontend
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Failed to accept new connection for %v", conn.RemoteAddr())
			if e, ok := err.(net.Error); ok {
				if e.Temporary() {
					continue
				}
			}
			log.Printf("cannot continue here, terminating due to error: %s", err)
			panic(err)
		}
		log.Printf("Accepted new connection from %v", conn.RemoteAddr())

		// proxy the connection to an backend
		go proxyConnection(conn, backendAddr)
	}
}
