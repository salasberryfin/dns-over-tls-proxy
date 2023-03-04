/*
DNS proxy:
	- Listen for "standard" DNS queries
	- Proxy to use DNS-over-TLS
*/
package main

import (
	"sync"

	"github.com/salasberryfin/dns-over-tls-proxy/network"
)

// main simply runs two concurrent go routines: TCP Listener and UDP
func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go network.CreateListenerTCP(&wg)
	go network.CreateListenerUDP(&wg)
	wg.Wait()
}
