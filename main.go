package main

import (
	"log"

	"github.com/HsiaoCz/forstore/p2p"
)

func main() {
	tcpOpts := p2p.TCPTransportOps{
		ListenAddr:    ":3001",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}
	tr := p2p.NewTCPTransport(tcpOpts)
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
