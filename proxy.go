package main

import (
	"crypto/tls"
	"io"
	"log"
	"net"
	"sync"
	"time"

	proxyproto "github.com/pires/go-proxyproto"
)

func proxyConnection(c net.Conn, backendAddr string) (err error) {
	// unwrap TLS Connection
	c = tls.Server(c, getTLSConfig(certFile, keyFile))

	// dial the backend
	upConn, err := net.DialTimeout("tcp", backendAddr, time.Duration(1000)*time.Millisecond)
	if err != nil {
		log.Printf("Failed to dial backend connection %v: %v", backendAddr, err)
		c.Close()
		return
	}
	log.Printf("Initiated new connction to backend: %v %v", upConn.LocalAddr(), upConn.RemoteAddr())

	// send the proxyprotocol header
	h := getProxyProtoHeaderFor(c)
	_, err = h.WriteTo(upConn)
	if err != nil {
		log.Printf("Failed to write proxy protocol header to backend: %s", err)
		c.Close()
		return
	}

	// join the connections
	joinConnections(c, upConn)
	return
}

func getProxyProtoHeaderFor(c net.Conn) *proxyproto.Header {

	sAddr := c.RemoteAddr().(*net.TCPAddr)
	dAddr := c.LocalAddr().(*net.TCPAddr)

	return &proxyproto.Header{
		Version:            1,
		Command:            proxyproto.PROXY,
		TransportProtocol:  proxyproto.TCPv4,
		SourceAddress:      sAddr.IP,
		DestinationAddress: dAddr.IP,
		SourcePort:         uint16(sAddr.Port),
		DestinationPort:    uint16(dAddr.Port),
	}
}

func joinConnections(c1 net.Conn, c2 net.Conn) {
	var wg sync.WaitGroup
	halfJoin := func(dst net.Conn, src net.Conn) {
		defer wg.Done()
		defer dst.Close()
		defer src.Close()
		for {
			n, err := io.Copy(dst, src)
			if e, ok := err.(net.Error); ok {
				if e.Temporary() {
					continue
				}
			}
			log.Printf("Copy from %v to %v failed after %d bytes with error %v", src.RemoteAddr(), dst.RemoteAddr(), n, err)
			return
		}
	}

	log.Printf("Joining connections: %v %v", c1.RemoteAddr(), c2.RemoteAddr())
	wg.Add(2)
	go halfJoin(c1, c2)
	go halfJoin(c2, c1)
	wg.Wait()
}
