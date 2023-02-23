FROM golang:alpine3.17

WORKDIR /app

COPY go.mod ./
COPY *.go ./

RUN go build -o /dns-over-tls-proxy

EXPOSE 53

CMD [ "/dns-over-tls-proxy" ]
