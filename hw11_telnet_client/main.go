package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	var timeout time.Duration

	flag.DurationVar(&timeout, "timeout", 10*time.Second, "Timeout for connection in seconds")
	flag.Parse()
	if len(flag.Args()) < 2 {
		fmt.Fprintln(os.Stderr, "...Required args are not provided")
	}
	host := flag.Arg(0)
	port := flag.Arg(1)
	address := net.JoinHostPort(host, port)

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	err := client.Connect()
	if err != nil {
		fmt.Fprintln(os.Stderr, "...Could not establish connection")
		return
	}

	go func() {
		for {
			err := client.Receive()
			if err != nil {
				fmt.Fprintln(os.Stderr, "...Connection closed by peer")
				return
			}
		}
	}()

	go func() {
		for {
			err := client.Send()
			if err != nil {
				fmt.Fprintln(os.Stderr, "...Connection closed by peer")
				return
			}
		}
	}()

	<-client.Done()
	client.Close()
}
