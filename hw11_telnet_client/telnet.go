package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
	Done() <-chan struct{}
}

type TelnetClientImpl struct {
	Address string
	Timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
	done    chan struct{}
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	done := make(chan struct{}, 1)
	return &TelnetClientImpl{
		Address: address,
		Timeout: timeout,
		in:      in,
		out:     out,
		done:    done,
	}
}

func (t *TelnetClientImpl) Connect() error {
	conn, err := net.DialTimeout("tcp", t.Address, t.Timeout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "...Failed to connect to %s\n", t.Address)
		return err
	}
	_, err = fmt.Fprintf(os.Stderr, "...Connected to %s\n", t.Address)
	if err != nil {
		return err
	}
	t.conn = conn
	return nil
}

func (t *TelnetClientImpl) Close() error {
	return t.conn.Close()
}

func (t *TelnetClientImpl) Send() error {
	var buf bytes.Buffer
	in := io.TeeReader(t.in, &buf)
	_, err := io.Copy(t.conn, in)
	if err != nil {
		t.done <- struct{}{}
		return err
	}
	_, err = buf.Read(make([]byte, 1))
	if errors.Is(err, io.EOF) {
		fmt.Fprintln(os.Stderr, "...EOF")
		t.done <- struct{}{}
		return err
	}
	return err
}

func (t *TelnetClientImpl) Receive() error {
	_, err := io.Copy(t.out, t.conn)
	if err != nil {
		t.done <- struct{}{}
	}
	return err
}

func (t *TelnetClientImpl) Done() <-chan struct{} {
	return t.done
}
