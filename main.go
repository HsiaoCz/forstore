package main

import (
	"log"

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
	s1 := makeServer(":3000", ":4000")
	s2 := makeServer(":4000", ":3000")
	go func() {
		if err := s1.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := s2.Start(); err != nil {
		log.Fatal(err)
	}
}

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpTransportOpts := p2p.TCPTransportOps{
		ListenAddr:    listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		// TODO onPeer func
		// OnPeer:        p2p.NOPOnPeerFunc,
	}
	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)
	fileServerOpts := FileServerOpts{
		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootStrapNodes:    nodes,
	}

	s := NewFileServer(fileServerOpts)

	tcpTransport.OnPeer = s.OnPeer

	return s
}
