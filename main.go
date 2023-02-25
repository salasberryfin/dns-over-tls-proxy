/*
DNS proxy:
	- Listen for "standard" DNS queries
	- Proxy to use DNS-over-TLS
*/
package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

const (
	portTCP           = "53"                 // TCP port
	portUDP           = "54"                 // UDP port
	dnsServerHost     = "1.1.1.1"            // use Cloudflare DNS server
	dnsServerHostName = "cloudflare-dns.com" // hostname to validate the certificate against
	dnsServerPortTLS  = "853"                // DNS over TLS port
)

// parseMessage reads the data received on the connection
func parseMessage(conn net.Conn) ([]byte, error) {
	buf := make([]byte, 0, 4096)
	text := make([]byte, 1024)
	n, err := conn.Read(text)
	if err != nil && err != io.EOF {
		return nil, err
	}
	buf = append(buf, text[:n]...)

	return buf, nil
}

// validateCert checks that the certificate is valid and has not expired
func validateCert(conn *tls.Conn) error {
	// verify server certificate
	log.Printf("Validating certificate for %s\n", dnsServerHostName)
	err := conn.VerifyHostname(dnsServerHostName)
	if err != nil {
		return err
	}
	expires := conn.ConnectionState().PeerCertificates[0].NotAfter
	diff := expires.Sub(time.Now())
	if diff <= 0 {
		return fmt.Errorf("The certitificate expired on %v\n", expires)
	}

	return nil
}

// resolveTLS gets the request content and forwards it to the end DNS server over TLS
func resolveTLS(buf []byte) ([]byte, error) {
	addr := dnsServerHost + ":" + dnsServerPortTLS
	log.Printf("Establishing TLS connection with %s\n", addr)
	conf := &tls.Config{}
	conn, err := tls.Dial("tcp", addr, conf)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	err = validateCert(conn)
	if err != nil {
		return nil, err
	}
	log.Printf("Connection established!\n")
	_, err = conn.Write(buf)
	if err != nil {
		return nil, err
	}
	resp, err := parseMessage(conn)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// main simply runs two concurrent go routines: TCP Listener and UDP
func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go createListenerTCP(&wg)
	go createListenerUDP(&wg)
	wg.Wait()
}
