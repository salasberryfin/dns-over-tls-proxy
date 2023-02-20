/*
DNS proxy:
	- Listen for "regular" DNS requests
	- Proxy to use DNS over TLS

Features:
	- Listen for TCP requests: concurrently
	- Parse connection message
	- Establish TLS connection to DNS-over-TLS end
		- Certificate
	- Receive TCP response
	- Send back to client

Extra:
	- Add same for UDP->TCP->UDP
*/
package main

import (
	"io"
	"log"
	"net"
	"strings"
)

const (
	protocol = "tcp"
	port     = "5353" // running on 5353 to avoid having to use sudo during development
)

func parseDNSMessage(buf []byte) {
	/*
		DNS request format:
			- 2: Trans ID
			- 2: Parameters
			- byte 14: start of the `queries` field - non-fixed length
	*/
	dnsName := []string{}
	queryFirstByte := 14
	length := 0
	for {
		current := length + queryFirstByte
		if buf[current] == 0 {
			break
		}
		dnsName = append(dnsName, string(buf[queryFirstByte:buf[current]+byte(queryFirstByte)+1]))
		queryFirstByte += int(buf[current]) + 1
	}

	log.Printf("Domain: %s\n", strings.Join(dnsName, "."))
}

func handleRequest(conn net.Conn) {
	// check with `kdig -d @localhost -p 5353 +tcp example.com`
	buf := make([]byte, 0, 4096)
	text := make([]byte, 256)
	n, err := conn.Read(text)
	if err != nil {
		if err != io.EOF {
			log.Fatalf("Erorr reading input: %v\n", err)
		}
	}
	buf = append(buf, text[:n]...)
	parseDNSMessage(buf)
	conn.Close()
}

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
		go handleRequest(conn)
	}
}
