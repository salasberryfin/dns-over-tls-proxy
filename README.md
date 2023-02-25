# Proxy DNS to DNS-over-TLS

This is a basic implementation of a proxy that listens for standard, non-secure, DNS requests 
and establishes a TLS connection with a DNS server that supports DNS-over-TLS 
-currently Cloudflare, but could be configured to use any other secure DNS server- 
and forwards the response to the original client.

This solution allows a client that doesn't support DNS-over-TLS to use a proxy that 
handles the DNS request with a TLS-compatible server.

## Overview

- Proxy starts a server that listens for DNS requests.
    - TCP Listener.
    - UDP Packet listener.
- * If the DNS query is sent to the UDP proxy, the data is first formatted by adding the length of the query.
- A new connection is established with the DNS-over-TLS server.
    - The proxy starts a TLS handshake with the server.
    - The certificate received from the server is validated: hostname and expiration date.
- The response from the DNS-over-TLS server is parsed a sent back to the client through the original connection.
- * If the DNS query was originally sent through UDP, data is truncated to create a valid UDP packet. 

The process is transparent to the client, which doesn't need to be aware of the proxy configuration, and the 
server supports concurrent requests from multiple clients.

## Build & Deploy

The proxy application can be built and deployed using `docker compose`:

```
docker compose up
```

- It creates a container from an image built from the `Dockerfile` in the project folder.
- The Go application is compiled and run.
- Ports are exposed in the container for both TCP and UDP:
    - TCP container port: 53
    - UDP container port: 54
- The compose file is configured to make the ports available to the host:
    - TCP host port: 5353
    - UDP host port: 5354

## Test

Once the application is up and running, the proxy functionality can be tested using a tool like `dig`.
`dig` doesn't support DNS-over-TLS, so it's a good way of simulating a real-world client for this type of application.

Considering the application is running with the default `docker compose` configuration:
- Test TCP DNS resolution:
```
dig @127.0.0.1 -p 5353 google.com +tcp
```

- Test UDP DNS resolution:
```
dig @127.0.0.1 -p 5354 google.com
```

The project includes folder `tests` with two Bash scripts that send TCP/UDP requests 
every second that can be used to validate the server concurrency.

## Security Disclaimer

- Due to the proxy implementation, the communication between the client application and proxy 
is not secure and hence is still vulnerable to man-in-the-middle attacks.
- DNS requests from the proxy to the DNS nameserver are encrypted but an eavesdropper can identify 
that DNS-encrypted traffic is going through the network by observing traffic going through port 853.
- Privacy can still be a concern even when using an encrypted connection and the DNS server must be 
owned by a trusted entity.

## Integration

Considering the massive use of the proxy, the deployment needs to be highly scalable.

- Use a Kubernetes cluster for deployment
- Define one pod with two containers:
    - Application: where the Go program runs.
    - Ambassador: proxy sitting between the application and the clients.
- The application is deployed to a container where the Go binary is run.
    - This container exposes ports for TCP and UDP to containers in the pod.
- The ambassador container, which can be an Nginx proxy.
    - This container includes an independent listener for each transport protocol.
- Configure Horizontal Pod Autoscaling for the application.
- Expose the service through a Load Balancer that manages traffic in case of scaling.

## Improvements

- Include protection mechanisms for clients:
    - DNS validation.
    - Blacklisting.
- Improve performance by adding a cache of requested domains.
- Handle invalid requests and error responses.
