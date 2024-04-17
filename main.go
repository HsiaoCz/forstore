package main

import (
	"log"

	"github.com/HsiaoCz/forstore/p2p"
)

func main() {
	tr := p2p.NewTCPTransport(":3001")
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
