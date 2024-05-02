package main

import (
	"log"
	"time"

	"github.com/HsiaoCz/forstore/p2p"
)

// func main() {
// 	tcpOpts := p2p.TCPTransportOps{
// 		ListenAddr:    ":3001",
// 		HandshakeFunc: p2p.NOPHandshakeFunc,
// 		Decoder:       p2p.DefaultDecoder{},
// 		OnPeer:        p2p.NOPOnPeerFunc,
// 	}
// 	tr := p2p.NewTCPTransport(tcpOpts)

// 	go func() {
// 		for {
// 			msg := <-tr.Consume()
// 			fmt.Printf("%+v\n", msg)
// 		}
// 	}()

// 	if err := tr.ListenAndAccept(); err != nil {
// 		log.Fatal(err)
// 	}

// 	select {}
// }

func main() {

	tcpTransportOpts := p2p.TCPTransportOps{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		// TODO onPeer func
		// OnPeer:        p2p.NOPOnPeerFunc,
	}

	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fileServerOpts := FileServerOpts{
		StorageRoot:       "3000_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
	}

	s := NewFileServer(fileServerOpts)
	go func() {
		time.Sleep(time.Second * 3)
		s.Stop()
	}()

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

}
