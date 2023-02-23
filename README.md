# Proxy DNS to DNS-over-TLS

This is a basic implementation of a proxy that listens for standard, non-secure, DNS requests 
and establishes a TLS connection with a DNS server that supports DNS-over-TLS 
-currently Cloudflare, but could be configured to use any other secure DNS server- 
and forwards the response to the original client.

This solution allows a client that doesn't support DNS-over-TLS to use a proxy that 
handles the DNS request with a TLS-compatible server.

## Solution

- Proxy starts a server that listens for DNS requests based on TCP -UDP pending-.
- A new connection is established with the DNS-over-TLS server.
    - The proxy starts a TLS handshake with the server.
    - The certificate received from the server is validated: hostname and expiration date.
- The response from the DNS-over-TLS server is parsed a sent back to the client through the original connection.

The process is transparent to the client, who doesn't need to be aware of the proxy configuration.

## Build & Deploy

## Test

Once the application is up and running, the proxy functionality can be tested using a tool like `dig`.
`dig` doesn't support DNS-over-TLS, so it's a good way of simulating a real-world client for this type of application.

Considering the application is running on port 5353 in the host:
```
dig @127.0.0.1 -p 5353 google.com +tcp
```

## Disclaimer

- Due to the proxy implementation, the communication between the client application and proxy 
is not secure and hence is still vulnerable to man-in-the-middle attacks.

## Integration
