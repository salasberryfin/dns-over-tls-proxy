package network

import (
	"log"
	"net"
	"sync"

	"github.com/salasberryfin/dns-over-tls-proxy/format"
)

const (
	portTCP = "53" // TCP port
)

// handlerTCP processes the TCP DNS queries received from the client
func handlerTCP(conn net.Conn) {
	log.Printf("Received TCP connection\n")
	buf, err := format.ParseMessage(conn)
	if err != nil {
		log.Printf("Failed to format input: %v", err)
		return
	}
	defer conn.Close()
	resp, err := resolveTLS(buf)
	if err != nil {
		log.Printf("Failed to request domain resolution: %s\n", err)
		return
	}
	conn.Write(resp)
	log.Printf("Server response sent\n")
}

// CreateListenerTCP listens for TCP connections and sends response back to client
func CreateListenerTCP(wg *sync.WaitGroup) {
	protocol := "tcp"
	ln, err := net.Listen(protocol, ":"+portTCP)
	if err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
	defer wg.Done()
	defer ln.Close()
	log.Printf("Listening on %s port %s...\n", protocol, portTCP)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handlerTCP(conn)
	}
}
