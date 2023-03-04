package format

import (
	"io"
	"net"
)

// ParseMessage reads the data received on the connection
func ParseMessage(conn net.Conn) ([]byte, error) {
	buf := make([]byte, 0, 4096)
	text := make([]byte, 1024)
	n, err := conn.Read(text)
	if err != nil && err != io.EOF {
		return nil, err
	}
	buf = append(buf, text[:n]...)

	return buf, nil
}
