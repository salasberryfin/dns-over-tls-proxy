FROM golang:alpine3.17

WORKDIR /app

COPY go.mod ./
COPY *.go ./

RUN go build -o /dns-over-tls-proxy

# TCP port: 53, UDP port: 54
EXPOSE 53
EXPOSE 54/udp

CMD [ "/dns-over-tls-proxy" ]
