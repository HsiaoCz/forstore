package main

import (
	"fmt"
	"log"

	"github.com/HsiaoCz/forstore/p2p"
)

func main() {
	tcpOpts := p2p.TCPTransportOps{
		ListenAddr:    ":3001",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        p2p.NOPOnPeerFunc,
	}
	tr := p2p.NewTCPTransport(tcpOpts)

	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("%+v\n", msg)
		}
	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
