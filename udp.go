package main

import (
	"log"
	"net"
	"sync"
)

// handlerUDP processes the UDP DNS queries received from the client
func handlerUDP(pc net.PacketConn, buf []byte, addr net.Addr) {
	log.Printf("Received UDP packet\n")
	newBuf := []byte{byte(len(buf))}
	newBuf = append(newBuf, buf...)
	resp, err := resolveTLS(newBuf)
	if err != nil {
		log.Printf("Failed to request domain resolution: %s\n", err)
		return
	}
	pc.WriteTo(resp[2:], addr)
	log.Printf("Server response sent\n")
}

func createListenerUDP(wg *sync.WaitGroup) {
	protocol := "udp"
	pc, err := net.ListenPacket(protocol, ":"+portUDP)
	if err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
	defer wg.Done()
	defer pc.Close()
	log.Printf("Listening on %s port %s...\n", protocol, portUDP)
	for {
		buf := make([]byte, 1024)
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			log.Print(err)
			continue
		}
		go handlerUDP(pc, buf[:n], addr)
	}
}
