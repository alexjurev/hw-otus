package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

// TelnetClient represents simple client for interacting with TCP server.
type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type myTelnetClient struct {
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
	timeout time.Duration
	address string
}

var (
	errClientClosed = fmt.Errorf("client close the connection")
	errServerClosed = fmt.Errorf("server close the connection")
)

// NewTelnetClient creates new client.
func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &myTelnetClient{
		in:      in,
		out:     out,
		timeout: timeout,
		address: address,
	}
}

func (t *myTelnetClient) Connect() (err error) {
	t.conn, err = net.DialTimeout("tcp", t.address, t.timeout)
	return err
}

func (t *myTelnetClient) Send() error {
	return pushToFrom(t.conn, t.in, errClientClosed)
}

func (t *myTelnetClient) Receive() error {
	return pushToFrom(t.out, t.conn, errServerClosed)
}

func (t *myTelnetClient) Close() error {
	return t.conn.Close()
}

func pushToFrom(dst io.Writer, src io.Reader, expectedError error) error {
	const bufSize = 128

	buf := make([]byte, bufSize)

	i, err := src.Read(buf)
	if err != nil {
		return expectedError
	}

	_, err = io.WriteString(dst, string(buf[:i]))

	return err
}
