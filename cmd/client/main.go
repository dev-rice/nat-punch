package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	IntroducerAddress string `long:"introducer-address" description:"address of the introducer server" required:"true"`
}

func main() {
	opts := Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	clientAddr, err := net.ResolveUDPAddr("udp", "")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	conn, err := net.ListenUDP("udp", clientAddr)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	defer conn.Close()

	// connect to introducer server
	introducerAddr, err := net.ResolveUDPAddr("udp", opts.IntroducerAddress)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	msg := "hello introducer"
	n, err := conn.WriteToUDP([]byte(msg), introducerAddr)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Printf("wrote %d bytes: %s\n", n, msg)

	// wait for introducer server to respond back with peer's address
	// this will happen when another peer connects to the introducer server
	buf := make([]byte, 1024)
	n, rAddr, err := conn.ReadFromUDP(buf)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Printf("read %d bytes from %s: %s\n", n, rAddr, string(buf[0:n]))
	peerAddr, err := net.ResolveUDPAddr("udp", string(buf[0:n]))
	if err != nil {
		fmt.Printf("error resolving peer address: %v\n", err)
	}
	fmt.Printf("introducer introduced me to peer '%v'\n", peerAddr)

	// send your secret number to peer every second
	rand.Seed(time.Now().UnixNano())
	secretNumber := rand.Int31()
	fmt.Printf("my secret number: %d\n", secretNumber)
	go func(peerAddr *net.UDPAddr) {
		ticker := time.NewTicker(1 * time.Second)
		for ; true; <-ticker.C {
			fmt.Printf("sending secret number to %s\n", peerAddr)
			_, err := conn.WriteTo([]byte(fmt.Sprintf("hello. secret number: %d", secretNumber)), peerAddr)
			if err != nil {
				fmt.Printf("error sending secret number: %v\n", err)
			}
		}
	}(peerAddr)

	// listen for udp packets sent (from the peer)
	for {
		buf := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("read error: %v\n", err)
		}
		fmt.Printf("read %d bytes from %s: %s\n", n, addr, string(buf[0:n]))
	}
}
