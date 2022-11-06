package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

type contextKey string

var (
	timeout       time.Duration
	errContextKey contextKey = "err"
)

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "set connection timeout")
}

func main() {
	flag.Parse()

	var wg sync.WaitGroup

	tail := flag.Args()
	if len(tail) < 2 {
		os.Stderr.WriteString("expected at least 2 arguments")
		os.Exit(1)
	}

	srv := tail[0]
	port := tail[1]

	client := NewTelnetClient(net.JoinHostPort(srv, port), timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		os.Stderr.WriteString(fmt.Sprintf("cannot connect to server: %s", err))
		os.Exit(1)
	}

	wg.Add(2)
	ctx, cancel := context.WithCancel(context.Background())
	ctxPtr := &ctx

	defer func() {
		if v := (*ctxPtr).Value(errContextKey); v != nil {
			os.Stderr.WriteString(fmt.Sprintf("cannot communicate with server: %s", v))
			os.Exit(2)
		}
	}()
	go receive(&ctxPtr, cancel, client, &wg)
	go send(&ctxPtr, cancel, client, &wg)
	wg.Wait()
}

func receive(ctx **context.Context, cancel context.CancelFunc, client TelnetClient, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-(**ctx).Done():
			return
		default:
			err := client.Receive()
			if err != nil {
				client.Close()
				c := context.WithValue(**ctx, errContextKey, err)
				*ctx = &c
				cancel()

				return
			}
		}
	}
}

func send(ctx **context.Context, cancel context.CancelFunc, client TelnetClient, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-(**ctx).Done():
			return
		default:
			err := client.Send()
			if err != nil {
				client.Close()
				c := context.WithValue(**ctx, errContextKey, err)
				*ctx = &c
				cancel()

				return
			}
		}
	}
}
