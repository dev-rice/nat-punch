package main

import (
	"fmt"
	"net"
	"os"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	ListenAddress string `long:"listen-address" description:"address to start the UDP server on" required:"true"`
}

func main() {
	opts := Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	addr, err := net.ResolveUDPAddr("udp", opts.ListenAddress)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	defer conn.Close()
	fmt.Printf("listening on %v\n", addr.String())

	clientAddrs := []*net.UDPAddr{}
	for {
		buf := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("read error: %v\n", err)
		}
		fmt.Printf("read %d bytes from %s: %s\n", n, addr, string(buf[0:n]))

		// keep track of the client
		clientAddrs = append(clientAddrs, addr)

		// when 2 clients have connected, tell both of them the other's address
		if len(clientAddrs) >= 2 {
			fmt.Printf("sending peer info to %v\n", clientAddrs[0])
			conn.WriteTo([]byte(clientAddrs[1].String()), clientAddrs[0])
			fmt.Printf("sending peer info to %v\n", clientAddrs[1])
			conn.WriteTo([]byte(clientAddrs[0].String()), clientAddrs[1])
			// reset clientAddrs so we can introduce more peers!
			clientAddrs = []*net.UDPAddr{}
		}
	}
}
