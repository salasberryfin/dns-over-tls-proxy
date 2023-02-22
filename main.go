/*
DNS proxy:
	- Listen for "regular" DNS requests
	- Proxy to use DNS over TLS

Features:
	- Listen for TCP requests: concurrently
	- Establish TLS connection to DNS-over-TLS end
		- Certificate
	- Receive TCP response
	- Send back to client

Extra:
	- Add same for UDP->TCP->UDP
*/
package main

import (
	"crypto/tls"
	"io"
	"log"
	"net"
)

const (
	protocol         = "tcp"     // protocol the proxy listens to defaults to tcp
	port             = "5353"    // running on 5353 to avoid having to use sudo during development
	dnsServerHost    = "1.1.1.1" // use Cloudflare DNS server
	dnsServerPort    = "53"
	dnsServerPortTLS = "853"
)

// parseResponse reads the data received on the connection
func parseResponse(conn net.Conn) ([]byte, error) {
	buf := make([]byte, 0, 4096)
	text := make([]byte, 256)
	n, err := conn.Read(text)
	if err != nil && err != io.EOF {
		return nil, err
	}
	buf = append(buf, text[:n]...)

	return buf, nil
}

// resolve gets the request content and forwards it to the end DNS server
func resolveTLS(buf []byte) ([]byte, error) {
	addr := dnsServerHost + ":" + dnsServerPortTLS
	log.Printf("Establishing TLS connection with %s\n", addr)
	conf := &tls.Config{}
	conn, err := tls.Dial("tcp", addr, conf)
	if err != nil {
		return nil, err
	}
	log.Printf("Connection established\n")
	_, err = conn.Write(buf)
	if err != nil {
		return nil, err
	}
	resp, err := parseResponse(conn)
	if err != nil {
		return nil, err
	}
	conn.Close()

	return resp, nil
}

// resolve gets the request content and forwards it to the end DNS server
func resolve(buf []byte) ([]byte, error) {
	addr := dnsServerHost + ":" + dnsServerPort
	log.Printf("Connecting to DNS server %s\n", addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	_, err = conn.Write(buf)
	if err != nil {
		return nil, err
	}
	resp, err := parseResponse(conn)
	if err != nil {
		return nil, err
	}
	conn.Close()

	return resp, nil
}

// handler processed the dns requests received from the clients
func handler(conn net.Conn) {
	// check with `kdig -d @localhost -p 5353 +tcp example.com`
	buf := make([]byte, 0, 4096)
	text := make([]byte, 256)
	n, err := conn.Read(text)
	if err != nil && err != io.EOF {
		log.Fatalf("Erorr reading input: %v\n", err)
	}
	buf = append(buf, text[:n]...)
	//resp, err := resolve(buf)
	resp, err := resolveTLS(buf)
	if err != nil {
		log.Fatalf("Failed to request domain resolution: %s\n", err)
	}
	conn.Write(resp)

	conn.Close()
}

// main simply creates a new TCP server to act as proxy
func main() {
	ln, err := net.Listen(protocol, ":"+port)
	if err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
	defer ln.Close()
	log.Printf("Listening on port %s...\n", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handler(conn)
	}
}
