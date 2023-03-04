package network

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/salasberryfin/dns-over-tls-proxy/format"
)

const (
	dnsServerHost     = "1.1.1.1"            // use Cloudflare DNS server
	dnsServerHostName = "cloudflare-dns.com" // hostname to validate the certificate against
	dnsServerPortTLS  = "853"                // DNS over TLS port
)

// validateCert checks that the certificate is valid and has not expired
func validateCert(conn *tls.Conn) error {
	// verify server certificate
	log.Printf("Validating certificate for %s\n", dnsServerHostName)
	err := conn.VerifyHostname(dnsServerHostName)
	if err != nil {
		return err
	}
	expires := conn.ConnectionState().PeerCertificates[0].NotAfter
	if expires.Before(time.Now()) {
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
	resp, err := format.ParseMessage(conn)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
