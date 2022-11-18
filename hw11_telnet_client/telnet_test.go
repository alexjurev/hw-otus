package main

import (
	"bytes"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClientPositive(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go testBasicRun(t, &wg, l)
		go testBasicCheck(t, &wg, l)

		wg.Wait()
	})

	t.Run("server closed connection", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg, clientWg sync.WaitGroup
		wg.Add(2)
		clientWg.Add(1)

		go testServerCloseRun(t, &wg, &clientWg, l)
		go testServerCloseCheck(t, &wg, &clientWg, l)

		wg.Wait()
	})

	t.Run("client closed connection", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg, clientWg sync.WaitGroup
		wg.Add(2)
		clientWg.Add(1)

		go testClientCloseRun(t, &wg, &clientWg, l)
		go testClientCloseCheck(t, &wg, &clientWg, l)

		wg.Wait()
	})
}

func testBasicRun(t *testing.T, wg *sync.WaitGroup, l net.Listener) {
	t.Helper()
	go func() {
		defer wg.Done()

		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		timeout, err := time.ParseDuration("10s")
		require.NoError(t, err)

		client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
		require.NoError(t, client.Connect())
		defer func() { require.NoError(t, client.Close()) }()

		in.WriteString("hello\n")
		err = client.Send()
		require.NoError(t, err)

		err = client.Receive()
		require.NoError(t, err)
		require.Equal(t, "world\n", out.String())
	}()
}

func testBasicCheck(t *testing.T, wg *sync.WaitGroup, l net.Listener) {
	t.Helper()
	defer wg.Done()

	conn, err := l.Accept()
	require.NoError(t, err)
	require.NotNil(t, conn)
	defer func() { require.NoError(t, conn.Close()) }()

	request := make([]byte, 1024)
	n, err := conn.Read(request)
	require.NoError(t, err)
	require.Equal(t, "hello\n", string(request)[:n])

	n, err = conn.Write([]byte("world\n"))
	require.NoError(t, err)
	require.NotEqual(t, 0, n)
}

func testServerCloseRun(t *testing.T, wg *sync.WaitGroup, clientWg *sync.WaitGroup, l net.Listener) {
	t.Helper()
	defer wg.Done()

	in := &bytes.Buffer{}
	out := &bytes.Buffer{}

	timeout, err := time.ParseDuration("10s")
	require.NoError(t, err)

	client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
	require.NoError(t, client.Connect())
	defer func() { require.NoError(t, client.Close()) }()
	in.WriteString("hello\n")
	err = client.Send()
	require.NoError(t, err)
	clientWg.Wait()
	err = client.Receive()
	require.Error(t, err)
	require.Equal(t, errServerClosed, err)
}

func testServerCloseCheck(t *testing.T, wg *sync.WaitGroup, clientWg *sync.WaitGroup, l net.Listener) {
	t.Helper()
	defer func() {
		wg.Done()
		clientWg.Done()
	}()

	conn, err := l.Accept()
	require.NoError(t, err)
	require.NotNil(t, conn)
	defer func() { require.NoError(t, conn.Close()) }()
}

func testClientCloseRun(t *testing.T, wg *sync.WaitGroup, clientWg *sync.WaitGroup, l net.Listener) {
	t.Helper()
	defer wg.Done()
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	timeout, err := time.ParseDuration("10s")
	require.NoError(t, err)

	client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
	require.NoError(t, client.Connect())
	defer func() { require.NoError(t, client.Close()) }()

	in.WriteString("^D")
	err = client.Send()
	require.NoError(t, err)

	clientWg.Wait()

	err = client.Send()
	require.Error(t, err)
	require.Equal(t, errClientClosed, err)

	err = client.Receive()
	require.Error(t, err)
	require.Equal(t, errServerClosed, err)
}

func testClientCloseCheck(t *testing.T, wg *sync.WaitGroup, clientWg *sync.WaitGroup, l net.Listener) {
	t.Helper()
	defer func() {
		wg.Done()
		clientWg.Done()
	}()

	conn, err := l.Accept()
	require.NoError(t, err)
	require.NotNil(t, conn)
	defer func() { require.NoError(t, conn.Close()) }()

	request := make([]byte, 1024)
	_, err = conn.Read(request)
	require.NoError(t, err)
}
